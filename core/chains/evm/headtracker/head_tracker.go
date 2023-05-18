package headtracker

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	promCurrentHead = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "head_tracker_current_head",
		Help: "The highest seen head number",
	}, []string{"evmChainID"})

	promOldHead = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_very_old_head",
		Help: "Counter is incremented every time we get a head that is much lower than the highest seen head ('much lower' is defined as a block that is EVM.FinalityDepth or greater below the highest seen head)",
	}, []string{"evmChainID"})
)

// HeadsBufferSize - The buffer is used when heads sampling is disabled, to ensure the callback is run for every head
const HeadsBufferSize = 10

type headTracker struct {
	log             logger.Logger
	headBroadcaster httypes.HeadBroadcaster
	headSaver       httypes.HeadSaver
	mailMon         *utils.MailboxMonitor
	ethClient       evmclient.Client
	chainID         big.Int
	config          Config

	broadcastMB  *utils.Mailbox[*evmtypes.Head]
	headListener httypes.HeadListener
	chStop       utils.StopChan
	wgDone       sync.WaitGroup
	utils.StartStopOnce
}

// NewHeadTracker instantiates a new HeadTracker using HeadSaver to persist new block numbers.
func NewHeadTracker(
	lggr logger.Logger,
	ethClient evmclient.Client,
	config Config,
	headBroadcaster httypes.HeadBroadcaster,
	headSaver httypes.HeadSaver,
	mailMon *utils.MailboxMonitor,
) httypes.HeadTracker {
	chStop := make(chan struct{})
	lggr = lggr.Named("HeadTracker")
	return &headTracker{
		headBroadcaster: headBroadcaster,
		ethClient:       ethClient,
		chainID:         *ethClient.ConfiguredChainID(),
		config:          config,
		log:             lggr,
		broadcastMB:     utils.NewMailbox[*evmtypes.Head](HeadsBufferSize),
		chStop:          chStop,
		headListener:    NewHeadListener(lggr, ethClient, config, chStop),
		headSaver:       headSaver,
		mailMon:         mailMon,
	}
}

// Start starts HeadTracker service.
func (ht *headTracker) Start(ctx context.Context) error {
	return ht.StartOnce("HeadTracker", func() error {
		ht.log.Debugf("Starting HeadTracker with chain id: %v", ht.chainID.Int64())

		initialHead, err := ht.getInitialHead(ctx)
		if err != nil {
			if errors.Is(err, ctx.Err()) {
				return nil
			} else {
				return errors.Wrapf(err, "Error getting initial head")
			}
		} else if initialHead == nil {
			return errors.Wrapf(err, "Got nil initial head")

		}

		// Don't try to load headers before the latest finalized block, given an initial head.
		latestChain, err := ht.headSaver.LoadFromDB(ctx, initialHead)
		if err != nil {
			return err
		}
		if latestChain != nil {
			ht.log.Debugw(
				fmt.Sprintf("HeadTracker: Loaded heads from last block %v with hash %s", config.FriendlyBigInt(latestChain.ToInt()), latestChain.Hash.Hex()),
				"blockNumber", latestChain.Number,
				"blockHash", latestChain.Hash,
			)
		}

		if err := ht.handleNewHead(ctx, initialHead); err != nil {
			return errors.Wrap(err, "error handling initial head")
		}

		ht.wgDone.Add(2)
		go ht.headListener.ListenForNewHeads(ht.handleNewHead, ht.wgDone.Done)
		go ht.broadcastLoop()

		ht.mailMon.Monitor(ht.broadcastMB, "HeadTracker", "Broadcast", ht.chainID.String())

		return nil
	})
}

// Close stops HeadTracker service.
func (ht *headTracker) Close() error {
	return ht.StopOnce("HeadTracker", func() error {
		close(ht.chStop)
		ht.wgDone.Wait()
		return ht.broadcastMB.Close()
	})
}

func (ht *headTracker) Name() string {
	return ht.log.Name()
}

func (ht *headTracker) HealthReport() map[string]error {
	report := map[string]error{
		ht.Name(): ht.StartStopOnce.Healthy(),
	}
	maps.Copy(report, ht.headListener.HealthReport())
	return report
}

func (ht *headTracker) LatestChain() *evmtypes.Head {
	return ht.headSaver.LatestChain()
}

// Backfill fetches all missing heads up until the latest finalized
func (ht *headTracker) Backfill(ctx context.Context, head *evmtypes.Head, latestFinalized *evmtypes.Head) (err error) {
	if head.Number == latestFinalized.Number {
		return nil
	}
	mark := time.Now()
	l := ht.log.With("blockNumber", head.Number,
		"n", head.Number-latestFinalized.Number,
		"fromBlockHeight", latestFinalized,
		"toBlockHeight", head.Number-1)
	l.Debug("Starting backfill")

	var reqs []rpc.BatchElem
	for i := head.Number - 1; i >= latestFinalized.Number; i-- {
		existingHead := ht.headSaver.Chain(head.ParentHash)
		if existingHead != nil {
			head = existingHead
			continue
		}
		req := rpc.BatchElem{
			Method: "eth_getHeadByHash",
			Args:   []interface{}{head.ParentHash, true},
			Result: &evmtypes.Head{},
		}
		reqs = append(reqs, req)
	}
	if head.ParentHash != latestFinalized.Hash {
		return errors.New("Backfill failed: a reorg happened during backfill")
	}

	err = ht.fetchAndSaveHeads(ctx, reqs, latestFinalized)
	if err != nil {
		return errors.Wrap(err, "fetchAndSaveHead failed")
	}

	l.Debugw("Finished backfill",
		"time", time.Since(mark))
	return
}

func (ht *headTracker) getInitialHead(ctx context.Context) (*evmtypes.Head, error) {
	head, err := ht.ethClient.HeadByNumber(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch initial head")
	}
	loggerFields := []interface{}{"head", head}
	if head != nil {
		loggerFields = append(loggerFields, "blockNumber", head.Number, "blockHash", head.Hash)
	}
	ht.log.Debugw("Got initial head", loggerFields...)
	return head, nil
}

func (ht *headTracker) handleNewHead(ctx context.Context, head *evmtypes.Head) error {
	if head == nil {
		return errors.New("got nil head from subscription")
	}
	prevHead := ht.headSaver.LatestChain()
	previousFinalized := ht.headSaver.LatestFinalizedHead()

	ht.log.Debugw(fmt.Sprintf("Received new head %v", config.FriendlyBigInt(head.ToInt())),
		"blockHeight", head.ToInt(),
		"blockHash", head.Hash,
		"parentHeadHash", head.ParentHash,
	)

	if head.Number <= previousFinalized.BlockNumber() {
		if head.Number == previousFinalized.BlockNumber() && head.Hash != prevHead.Hash {
			return errors.Wrapf(errors.New("Invariant violation: current finalized block has the same number as received head, but different hash!"),
				"blockNumber", head.Number, "newHash", head.Hash, "finalizedHash", previousFinalized.BlockHash())
		}
		promOldHead.WithLabelValues(ht.chainID.String()).Inc()
		ht.log.Criticalf(`Got very old block with number %d (highest seen was %d and latest finalized %d).This is a problem and either means a very deep re-org occurred,`+
			`one of the RPC nodes has gotten far out of sync`, head.Number, prevHead.Number, previousFinalized.BlockNumber())
		ht.SvcErrBuffer.Append(errors.New("got very old block"))
		return nil
	}

	if head.Number <= prevHead.Number {
		if ht.headSaver.Chain(head.Hash) != nil {
				ht.log.Debugw("Head already in the database", "head", head.Hash.Hex())
				return nil
		}
		ht.log.Debugw("Got out of order head", "blockNum", head.Number, "head", head.Hash.Hex(), "prevHead", prevHead.Number)
		err := ht.headSaver.Save(ctx, head, previousFinalized)
		if ctx.Err() != nil {
			return nil
		} else if err != nil {
			return errors.Wrapf(err, "failed to save head: %#v", head)
		}
		return nil
	}

	if head.Number > prevHead.Number {
		latestFinalized, err := ht.headSaver.LatestFinalizedBlock(ctx, head.Number)
		if err != nil {
			return err
		}
		err = ht.headSaver.Save(ctx, head, latestFinalized)
		if ctx.Err() != nil {
			return nil
		} else if err != nil {
			return errors.Wrapf(err, "failed to save head: %#v", head)
		}
		promCurrentHead.WithLabelValues(ht.chainID.String()).Set(float64(head.Number))

		headWithChain := ht.headSaver.Chain(head.Hash)
		if headWithChain == nil {
			return errors.Errorf("HeadTracker#handleNewHighestHead headWithChain was unexpectedly nil")
		}
		err = ht.Backfill(ctx, head.EarliestInChain(), latestFinalized)
		if err != nil {
			ht.log.Warnw("Unexpected error while backfilling heads", "err", err)
		}

		ht.broadcastMB.Deliver(headWithChain)
	}
	return nil
}

func (ht *headTracker) broadcastLoop() {
	defer ht.wgDone.Done()

	samplingInterval := ht.config.EvmHeadTrackerSamplingInterval()
	if samplingInterval > 0 {
		ht.log.Debugf("Head sampling is enabled - sampling interval is set to: %v", samplingInterval)
		debounceHead := time.NewTicker(samplingInterval)
		defer debounceHead.Stop()
		for {
			select {
			case <-ht.chStop:
				return
			case <-debounceHead.C:
				item := ht.broadcastMB.RetrieveLatestAndClear()
				if item == nil {
					continue
				}
				ht.headBroadcaster.BroadcastNewLongestChain(item)
			}
		}
	} else {
		ht.log.Info("Head sampling is disabled - callback will be called on every head")
		for {
			select {
			case <-ht.chStop:
				return
			case <-ht.broadcastMB.Notify():
				for {
					item, exists := ht.broadcastMB.Retrieve()
					if !exists {
						break
					}
					ht.headBroadcaster.BroadcastNewLongestChain(item)
				}
			}
		}
	}
}

func (ht *headTracker) fetchAndSaveHeads(ctx context.Context, reqs []rpc.BatchElem, latestFinalized *evmtypes.Head) error {
	err := ht.ethClient.BatchCallContext(ctx, reqs)
	if err != nil {
		return errors.Wrap(err, "HeadTracker#batchFetchHeaders error fetching headers with BatchCallContext")
	}

	for _, r := range reqs {
		if r.Error != nil {
			return r.Error
		}
		head, is := r.Result.(*evmtypes.Head)

		if !is {
			return errors.Errorf("expected result to be a %T, got %T", &evmtypes.Head{}, r.Result)
		}
		if head == nil {
			return errors.New("invariant violation: got nil block")
		}
		if head.Hash == (common.Hash{}) {
			return errors.Errorf("missing block hash for block number: %d", head.Number)
		}
		if head.Number < 0 {
			return errors.Errorf("expected block number to be >= to 0, got %d", head.Number)
		}

		err = ht.headSaver.Save(ctx, head, latestFinalized)
		if err != nil {
			return err
		}
	}
	return nil
}

var NullTracker httypes.HeadTracker = &nullTracker{}

type nullTracker struct{}

func (*nullTracker) Start(context.Context) error    { return nil }
func (*nullTracker) Close() error                   { return nil }
func (*nullTracker) Ready() error                   { return nil }
func (*nullTracker) HealthReport() map[string]error { return map[string]error{} }
func (*nullTracker) Name() string                   { return "" }
func (*nullTracker) SetLogLevel(zapcore.Level)      {}
func (*nullTracker) Backfill(ctx context.Context, headWithChain *evmtypes.Head, lastFinalized *evmtypes.Head) (err error) {
	return nil
}
func (*nullTracker) LatestChain() *evmtypes.Head { return nil }
