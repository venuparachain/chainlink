package headtracker

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils/mathutil"
)

type headSaver struct {
	orm    ORM
	config Config
	logger logger.Logger
	heads  Heads
	client client.Client
}

func NewHeadSaver(lggr logger.Logger, orm ORM, config Config, client client.Client) httypes.HeadSaver {
	return &headSaver{
		orm:    orm,
		config: config,
		logger: lggr.Named("HeadSaver"),
		heads:  NewHeads(),
		client: client,
	}
}

func (hs *headSaver) Save(ctx context.Context, head *evmtypes.Head, latestFinalized *evmtypes.Head) error {
	previousFinalized := hs.heads.LatestFinalized()
	if previousFinalized == nil {
		// This shouldn't happen but it's best if we check anyway.
		return errors.New("Couldn't find previous finalized head!")
	}
	if latestFinalized.Number == previousFinalized.Number && latestFinalized.Hash != previousFinalized.Hash {
		return errors.Wrapf(errors.New("Invariant violation: latest finalized block has the same number as the previous one, but different hash!"),
			"blockNumber", latestFinalized.Number, "previousHash", previousFinalized.Hash, "newHash", latestFinalized.Hash)
	}
	if latestFinalized.Number < previousFinalized.Number {
		return errors.Wrapf(errors.New("Invariant violation: latest finalized block number is lower than the previous one!"),
			"latestFinalized", latestFinalized.Number, "previousFinalized", previousFinalized.Number)
	}

	if err := hs.orm.IdempotentInsertHead(ctx, head); err != nil {
		return err
	}

	hs.heads.AddHeads(latestFinalized, head)
	if hs.heads.Count() > int(hs.config.EvmHeadTrackerTotalHeadsLimit()) {
		hs.logger.Warnw("chain larger than EvmHeadTrackerTotalHeadsLimit. In memory heads exceed limit.",
			"headsCount", hs.heads.Count(), "evmHeadTrackerTotalHeadsLimit", hs.config.EvmHeadTrackerTotalHeadsLimit())
	}

	return hs.orm.TrimOldHeads(ctx, uint(latestFinalized.Number))
}

func (hs *headSaver) LoadFromDB(ctx context.Context, initialHead *evmtypes.Head) (chain *evmtypes.Head, err error) {
	latestFinalized, err := hs.LatestFinalizedBlock(ctx, initialHead.Number)
	if err != nil {
		return nil, err
	}

	heads, err := hs.orm.LatestHeads(ctx, uint(latestFinalized.Number))
	if err != nil || len(heads) == 0 {
		return nil, err
	}

	hs.heads.AddHeads(latestFinalized, heads...)
	if hs.heads.Count() > int(hs.config.EvmHeadTrackerTotalHeadsLimit()) {
		hs.logger.Warnw("chain larger than EvmHeadTrackerTotalHeadsLimit. In memory heads exceed limit.",
			"headsCount", hs.heads.Count(), "evmHeadTrackerTotalHeadsLimit", hs.config.EvmHeadTrackerTotalHeadsLimit())
	}
	return hs.heads.LatestHead(), nil
}

func (hs *headSaver) LatestHeadFromDB(ctx context.Context) (head *evmtypes.Head, err error) {
	return hs.orm.LatestHead(ctx)
}

func (hs *headSaver) LatestChain() *evmtypes.Head {
	head := hs.heads.LatestHead()
	if head == nil {
		return nil
	}
	if head.ChainLength() < hs.config.EvmFinalityDepth() {
		hs.logger.Debugw("chain shorter than EvmFinalityDepth", "chainLen", head.ChainLength(), "evmFinalityDepth", hs.config.EvmFinalityDepth())
	}
	return head
}

func (hs *headSaver) Chain(hash common.Hash) *evmtypes.Head {
	return hs.heads.HeadByHash(hash)
}

func (hs *headSaver) LatestFinalizedHead() *evmtypes.Head {
	return hs.heads.LatestFinalized()
}

func (hs *headSaver) LatestFinalizedBlock(ctx context.Context, currentHeadNumber int64) (*evmtypes.Head, error) {
	if !hs.config.EvmFinalityTag() {
		// If heads is not empty, calculate the latest finalized head from the highest head saved, otherwise just use the current one
		var newFinalized int64
		previousFinalized := hs.heads.LatestFinalized()
		if previousFinalized == nil {
			newFinalized = currentHeadNumber - int64(hs.config.EvmFinalityDepth()-1)
		} else {
			newFinalized = mathutil.Max(previousFinalized.Number, currentHeadNumber-int64(hs.config.EvmFinalityDepth())-1)
		}
		if newFinalized == previousFinalized.Number {
			return previousFinalized, nil
		}
		return hs.client.HeadByNumber(ctx, big.NewInt(mathutil.Max(newFinalized, 0)))
	}
	return hs.client.LatestBlockByType(ctx, "finalized")
}

var NullSaver httypes.HeadSaver = &nullSaver{}

type nullSaver struct{}

func (*nullSaver) Save(ctx context.Context, head *evmtypes.Head, latestFinalized *evmtypes.Head) error {
	return nil
}
func (*nullSaver) LoadFromDB(ctx context.Context, initialHead *evmtypes.Head) (*evmtypes.Head, error) {
	return nil, nil
}
func (*nullSaver) LatestHeadFromDB(ctx context.Context) (*evmtypes.Head, error) { return nil, nil }
func (*nullSaver) LatestChain() *evmtypes.Head                                  { return nil }
func (*nullSaver) Chain(hash common.Hash) *evmtypes.Head                        { return nil }
func (*nullSaver) LatestFinalizedHead() *evmtypes.Head                          { return nil }
func (*nullSaver) LatestFinalizedBlock(ctx context.Context, currentHeadNumber int64) (*evmtypes.Head, error) {
	return nil, nil
}
