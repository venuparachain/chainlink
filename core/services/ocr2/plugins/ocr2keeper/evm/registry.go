package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	coreTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/patrickmn/go-cache"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_lookup_compatible_interface"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// DefaultUpkeepExpiration decides how long an upkeep info will be valid for. after it expires, a getUpkeepInfo
	// call will be made to the registry to obtain the most recent upkeep info and refresh this cache.
	DefaultUpkeepExpiration = 10 * time.Minute
	// DefaultCooldownExpiration decides how long a Mercury upkeep will be put in cool down for the first time. within
	// 10 minutes, subsequent failures will result in double amount of cool down period.
	DefaultCooldownExpiration = 5 * time.Second
	// DefaultApiErrExpiration decides a running sum of total errors of an upkeep in this 10 minutes window. it is used
	// to decide how long the cool down period will be.
	DefaultApiErrExpiration = 10 * time.Minute
	// CleanupInterval decides when the expired items in cache will be deleted.
	CleanupInterval = 15 * time.Minute
)

var (
	ErrLogReadFailure                = fmt.Errorf("failure reading logs")
	ErrHeadNotAvailable              = fmt.Errorf("head not available")
	ErrRegistryCallFailure           = fmt.Errorf("registry chain call failure")
	ErrBlockKeyNotParsable           = fmt.Errorf("block identifier not parsable")
	ErrUpkeepKeyNotParsable          = fmt.Errorf("upkeep key not parsable")
	ErrInitializationFailure         = fmt.Errorf("failed to initialize registry")
	ErrContextCancelled              = fmt.Errorf("context was cancelled")
	ErrABINotParsable                = fmt.Errorf("error parsing abi")
	ActiveUpkeepIDBatchSize    int64 = 1000
	FetchUpkeepConfigBatchSize       = 10
	separator                        = "|"
	reInitializationDelay            = 15 * time.Minute
	logEventLookback           int64 = 250
)

//go:generate mockery --quiet --name Registry --output ./mocks/ --case=underscore
type Registry interface {
	GetUpkeep(opts *bind.CallOpts, id *big.Int) (keeper_registry_wrapper2_0.UpkeepInfo, error)
	GetState(opts *bind.CallOpts) (keeper_registry_wrapper2_0.GetState, error)
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
	ParseLog(log coreTypes.Log) (generated.AbigenLog, error)
}

//go:generate mockery --quiet --name HttpClient --output ./mocks/ --case=underscore
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LatestBlockGetter interface {
	LatestBlock() int64
}

func NewEVMRegistryServiceV2_0(addr common.Address, client evm.Chain, mc *models.MercuryCredentials, lggr logger.Logger) (*EvmRegistry, error) {
	mercuryLookupCompatibleABI, err := abi.JSON(strings.NewReader(mercury_lookup_compatible_interface.MercuryLookupCompatibleInterfaceABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}
	keeperRegistryABI, err := abi.JSON(strings.NewReader(keeper_registry_wrapper2_0.KeeperRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}

	registry, err := keeper_registry_wrapper2_0.NewKeeperRegistry(addr, client.Client())
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create caller for address and backend", ErrInitializationFailure)
	}

	upkeepInfoCache, cooldownCache, apiErrCache := setupCaches(DefaultUpkeepExpiration, DefaultCooldownExpiration, DefaultApiErrExpiration, CleanupInterval)

	r := &EvmRegistry{
		HeadProvider: HeadProvider{
			ht:     client.HeadTracker(),
			hb:     client.HeadBroadcaster(),
			chHead: make(chan ocr2keepers.BlockKey, 1),
		},
		lggr:       lggr,
		poller:     client.LogPoller(),
		addr:       addr,
		client:     client.Client(),
		txHashes:   make(map[string]bool),
		registry:   registry,
		abi:        keeperRegistryABI,
		upkeeps:    make(map[string]upkeepEntry),
		packer:     &evmRegistryPackerV2_0{abi: keeperRegistryABI},
		headFunc:   func(types.BlockKey) {},
		chLog:      make(chan logpoller.Log, 1000),
		logFilters: newLogFiltersProvider(client.LogPoller()),
		mercury: MercuryConfig{
			cred:          mc,
			abi:           mercuryLookupCompatibleABI,
			upkeepCache:   upkeepInfoCache,
			cooldownCache: cooldownCache,
			apiErrCache:   apiErrCache,
		},
		hc:  http.DefaultClient,
		enc: EVMAutomationEncoder20{},
	}

	if err := r.registerEvents(client.ID().Uint64(), addr); err != nil {
		return nil, fmt.Errorf("logPoller error while registering automation events: %w", err)
	}

	return r, nil
}

func setupCaches(defaultUpkeepExpiration, defaultCooldownExpiration, defaultApiErrExpiration, cleanupInterval time.Duration) (*cache.Cache, *cache.Cache, *cache.Cache) {
	// cache that stores UpkeepInfo for callback during MercuryLookup
	upkeepInfoCache := cache.New(defaultUpkeepExpiration, cleanupInterval)

	// with apiErrCacheExpiration= 10m and cooldownExp= 2^errCount
	// then max cooldown = 2^10 approximately 17m at which point the cooldownExp > apiErrCacheExpiration so the count will get reset
	// cache for Mercurylookup Upkeeps that are on ice due to errors
	cooldownCache := cache.New(defaultCooldownExpiration, cleanupInterval)

	// cache for tracking errors for an Upkeep during MercuryLookup
	apiErrCache := cache.New(defaultApiErrExpiration, cleanupInterval)
	return upkeepInfoCache, cooldownCache, apiErrCache
}

var upkeepStateEvents = []common.Hash{
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegistered{}.Topic(),        // adds new upkeep id to registry
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceived{}.Topic(),          // adds new upkeep id to registry via migration
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepGasLimitSet{}.Topic(),       // unpauses an upkeep
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpaused{}.Topic(),          // updates the gas limit for an upkeep
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepOffchainConfigSet{}.Topic(), // updates upkeep config
}

var upkeepActiveEvents = []common.Hash{
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepPerformed{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryReorgedUpkeepReport{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryInsufficientFundsUpkeepReport{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryStaleUpkeepReport{}.Topic(),
}

type checkResult struct {
	ur  []EVMAutomationUpkeepResult20
	err error
}

type upkeepState int32

const (
	stateActive upkeepState = iota
	stateInactive
)

type upkeepTriggerType int32

const (
	conditionalTrigger upkeepTriggerType = iota
	logTrigger
)

type upkeepEntry struct {
	id              *big.Int
	state           upkeepState
	triggerType     upkeepTriggerType
	performGasLimit uint32
	offchainConfig  []byte
}

type UpkeepFilter func(upkeepEntry) bool

type upkeepFilters []UpkeepFilter

func (uf upkeepFilters) Apply(upkeep upkeepEntry) bool {
	for _, f := range uf {
		if !f(upkeep) {
			return false
		}
	}
	return true
}

func ActiveUpkeepsFilter() UpkeepFilter {
	return func(upkeep upkeepEntry) bool {
		return upkeep.state == stateActive
	}
}

func LogUpkeepsFilter() UpkeepFilter {
	return func(upkeep upkeepEntry) bool {
		return upkeep.triggerType == logTrigger
	}
}

func ConditionalUpkeepsFilter() UpkeepFilter {
	return func(upkeep upkeepEntry) bool {
		return upkeep.triggerType == conditionalTrigger
	}
}

type MercuryConfig struct {
	cred          *models.MercuryCredentials
	abi           abi.ABI
	upkeepCache   *cache.Cache
	cooldownCache *cache.Cache
	apiErrCache   *cache.Cache
}

type EvmRegistry struct {
	HeadProvider
	sync          utils.StartStopOnce
	lggr          logger.Logger
	poller        logpoller.LogPoller
	addr          common.Address
	client        client.Client
	registry      Registry
	abi           abi.ABI
	packer        *evmRegistryPackerV2_0
	chLog         chan logpoller.Log
	reInit        *time.Timer
	mu            sync.RWMutex
	txHashes      map[string]bool
	lastPollBlock int64
	ctx           context.Context
	cancel        context.CancelFunc
	upkeeps       map[string]upkeepEntry
	headFunc      func(types.BlockKey)
	runState      int
	runError      error
	mercury       MercuryConfig
	hc            HttpClient
	enc           EVMAutomationEncoder20

	logFilters *logFiltersProvider
}

// GetActiveUpkeepKeys uses the latest head and map of all active upkeeps to build a
// slice of upkeep keys.
func (r *EvmRegistry) GetActiveUpkeepIDs(ctx context.Context) ([]ocr2keepers.UpkeepIdentifier, error) {
	return r.GetUpkeepIDs(ctx, ActiveUpkeepsFilter(), ConditionalUpkeepsFilter())
}

// GetUpkeepIDs accepts filters to return a customized slice of upkeep.
func (r *EvmRegistry) GetUpkeepIDs(ctx context.Context, filters ...UpkeepFilter) ([]ocr2keepers.UpkeepIdentifier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	uf := upkeepFilters(filters)
	results := make([]ocr2keepers.UpkeepIdentifier, 0)

	for _, upkeep := range r.upkeeps {
		if !uf.Apply(upkeep) {
			continue
		}
		results = append(results, ocr2keepers.UpkeepIdentifier(upkeep.id.String()))
	}

	return results, nil
}

func (r *EvmRegistry) CheckUpkeep(ctx context.Context, mercuryEnabled bool, keys ...ocr2keepers.UpkeepKey) ([]ocr2keepers.UpkeepResult, error) {
	chResult := make(chan checkResult, 1)
	go r.doCheck(ctx, mercuryEnabled, keys, chResult)

	select {
	case rs := <-chResult:
		result := make([]ocr2keepers.UpkeepResult, len(rs.ur))
		for i := range rs.ur {
			result[i] = rs.ur
		}

		return result, rs.err
	case <-ctx.Done():
		// safety on context done to provide an error on context cancellation
		// contract calls through the geth wrappers are a bit of a black box
		// so this safety net ensures contexts are fully respected and contract
		// call functions have a more graceful closure outside the scope of
		// CheckUpkeep needing to return immediately.
		return nil, fmt.Errorf("%w: failed to check upkeep on registry", ErrContextCancelled)
	}
}

func (r *EvmRegistry) Name() string {
	return r.lggr.Name()
}

func (r *EvmRegistry) Start(ctx context.Context) error {
	return r.sync.StartOnce("AutomationRegistry", func() error {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.ctx, r.cancel = context.WithCancel(context.Background())
		r.reInit = time.NewTimer(reInitializationDelay)

		// initialize the upkeep keys; if the reInit timer returns, do it again
		{
			go func(cx context.Context, tmr *time.Timer, lggr logger.Logger, f func() error) {
				err := f()
				if err != nil {
					lggr.Errorf("failed to initialize upkeeps", err)
				}

				for {
					select {
					case <-tmr.C:
						err = f()
						if err != nil {
							lggr.Errorf("failed to re-initialize upkeeps", err)
						}
						tmr.Reset(reInitializationDelay)
					case <-cx.Done():
						return
					}
				}
			}(r.ctx, r.reInit, r.lggr, r.initialize)
		}

		// start polling logs on an interval
		{
			go func(cx context.Context, lggr logger.Logger, f func() error) {
				ticker := time.NewTicker(time.Second)

				for {
					select {
					case <-ticker.C:
						err := f()
						if err != nil {
							lggr.Errorf("failed to poll logs for upkeeps", err)
						}
					case <-cx.Done():
						ticker.Stop()
						return
					}
				}
			}(r.ctx, r.lggr, r.pollLogs)
		}

		// run process to process logs from log channel
		{
			go func(cx context.Context, ch chan logpoller.Log, lggr logger.Logger, f func(logpoller.Log) error) {
				for {
					select {
					case l := <-ch:
						err := f(l)
						if err != nil {
							lggr.Errorf("failed to process log for upkeep", err)
						}
					case <-cx.Done():
						return
					}
				}
			}(r.ctx, r.chLog, r.lggr, r.processUpkeepStateLog)
		}

		r.runState = 1
		return nil
	})
}

func (r *EvmRegistry) Close() error {
	return r.sync.StopOnce("AutomationRegistry", func() error {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.cancel()
		r.runState = 0
		r.runError = nil
		return nil
	})
}

func (r *EvmRegistry) Ready() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.runState == 1 {
		return nil
	}
	return r.sync.Ready()
}

func (r *EvmRegistry) HealthReport() map[string]error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.runState > 1 {
		r.sync.SvcErrBuffer.Append(fmt.Errorf("failed run state: %w", r.runError))
	}
	return map[string]error{r.Name(): r.sync.Healthy()}
}

func (r *EvmRegistry) initialize() error {
	startupCtx, cancel := context.WithTimeout(r.ctx, reInitializationDelay)
	defer cancel()

	upkeeps := make(map[string]upkeepEntry)

	r.lggr.Debugf("Re-initializing active upkeeps list")
	// get active upkeep ids from contract
	ids, err := r.getLatestIDsFromContract(startupCtx)
	if err != nil {
		return fmt.Errorf("failed to get ids from contract: %s", err)
	}

	var offset int
	for offset < len(ids) {
		batch := FetchUpkeepConfigBatchSize
		if len(ids)-offset < batch {
			batch = len(ids) - offset
		}

		entries, err := r.createUpkeepEntries(startupCtx, ids[offset:offset+batch])
		if err != nil {
			return fmt.Errorf("failed to get configs for id batch (length '%d'): %s", batch, err)
		}

		for _, entry := range entries {
			upkeeps[entry.id.String()] = entry
		}

		offset += batch
	}

	r.mu.Lock()
	r.upkeeps = upkeeps
	r.mu.Unlock()

	return nil
}

func (r *EvmRegistry) pollLogs() error {
	var latest int64
	var end int64
	var err error

	if end, err = r.poller.LatestBlock(); err != nil {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}

	r.mu.Lock()
	latest = r.lastPollBlock
	r.lastPollBlock = end
	r.mu.Unlock()

	// if start and end are the same, no polling needs to be done
	if latest == 0 || latest == end {
		return nil
	}

	{
		var logs []logpoller.Log

		if logs, err = r.poller.LogsWithSigs(
			end-logEventLookback,
			end,
			upkeepStateEvents,
			r.addr,
			pg.WithParentCtx(r.ctx),
		); err != nil {
			return fmt.Errorf("%w: %s", ErrLogReadFailure, err)
		}

		for _, log := range logs {
			r.chLog <- log
		}
	}

	return nil
}

func UpkeepFilterName(addr common.Address) string {
	return logpoller.FilterName("EvmRegistry - Upkeep events for", addr.String())
}

func (r *EvmRegistry) registerEvents(chainID uint64, addr common.Address) error {
	// Add log filters for the log poller so that it can poll and find the logs that
	// we need
	err := r.poller.RegisterFilter(logpoller.Filter{
		Name:      UpkeepFilterName(addr),
		EventSigs: append(upkeepStateEvents, upkeepActiveEvents...),
		Addresses: []common.Address{addr},
	})
	if err != nil {
		r.mu.Lock()
		r.mu.Unlock()
	}
	return err
}

func (r *EvmRegistry) processUpkeepStateLog(l logpoller.Log) error {
	hash := l.TxHash.String()
	if _, ok := r.txHashes[hash]; ok {
		return nil
	}
	r.txHashes[hash] = true

	rawLog := l.ToGethLog()
	abilog, err := r.registry.ParseLog(rawLog)
	if err != nil {
		return err
	}

	switch l := abilog.(type) {
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepOffchainConfigSet:
		r.lggr.Debugf("KeeperRegistryUpkeepOffchainConfigSet log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.updateUpkeepConfig(l.Id, l.OffchainConfig)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepCanceled:
		r.lggr.Debugf("KeeperRegistryUpkeepCanceled log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.updateUpkeepState(l.Id, stateInactive)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepPaused:
		r.lggr.Debugf("KeeperRegistryUpkeepPaused log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.updateUpkeepState(l.Id, stateInactive)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegistered:
		r.lggr.Debugf("KeeperRegistryUpkeepRegistered log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.addUpkeep(l.Id, false)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceived:
		r.lggr.Debugf("KeeperRegistryUpkeepReceived log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.addUpkeep(l.Id, false)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpaused:
		r.lggr.Debugf("KeeperRegistryUpkeepUnpaused log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.addUpkeep(l.Id, false)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepGasLimitSet:
		r.lggr.Debugf("KeeperRegistryUpkeepGasLimitSet log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.addUpkeep(l.Id, true)
	}

	return nil
}

func (r *EvmRegistry) addUpkeep(id *big.Int, force bool) {
	r.mu.Lock()
	if r.upkeeps == nil {
		r.upkeeps = make(map[string]upkeepEntry)
	}
	_, ok := r.upkeeps[id.String()]
	r.mu.Unlock()

	if !ok || force {
		entries, err := r.createUpkeepEntries(r.ctx, []*big.Int{id})
		if err != nil {
			r.lggr.Errorf("failed to get upkeep configs during adding active upkeep: %w", err)
			return
		}

		if len(entries) != 1 {
			return
		}
		r.mu.Lock()
		r.upkeeps[id.String()] = entries[0]
		r.mu.Unlock()
		return
	}
	// otherwise, just update the state to active
	r.updateUpkeepState(id, stateActive)
}

func (r *EvmRegistry) updateUpkeepConfig(id *big.Int, upkeepCfg []byte) {
	r.mu.Lock()

	upkeep, ok := r.upkeeps[id.String()]
	if !ok {
		r.mu.Unlock()
		// TODO: TBD if we want to add a new upkeep config if it doesn't exist
		r.lggr.Debugf("received update for upkeep config for id %s that does not exist", id.String())
		return
	}
	r.lggr.Debugw("updating upkeep config", "upkeepID", id.String(), "upkeepCfg", hexutil.Encode(upkeepCfg))
	upkeep.offchainConfig = upkeepCfg
	r.upkeeps[id.String()] = upkeep
	r.mu.Unlock()

	switch upkeep.triggerType {
	case logTrigger:
		if err := r.logFilters.Register(id.String(), upkeepCfg); err != nil {
			r.lggr.Debugw("failed to register log filter", "upkeepID", id.String())
		}
	default:
	}
}

func (r *EvmRegistry) updateUpkeepState(id *big.Int, state upkeepState) {
	r.mu.Lock()
	defer r.mu.Unlock()

	upkeep, ok := r.upkeeps[id.String()]
	if !ok {
		// TODO: TBD if we want to add a new upkeep config if it doesn't exist
		r.lggr.Debugf("received update for upkeep config for id %s that does not exist", id.String())
		return
	}
	upkeep.state = state
}

func (r *EvmRegistry) buildCallOpts(ctx context.Context, block *big.Int) (*bind.CallOpts, error) {
	opts := bind.CallOpts{
		Context:     ctx,
		BlockNumber: nil,
	}

	if block == nil || block.Int64() == 0 {
		if r.LatestBlock() != 0 {
			opts.BlockNumber = big.NewInt(r.LatestBlock())
		}
	} else {
		opts.BlockNumber = block
	}

	return &opts, nil
}

func (r *EvmRegistry) getLatestIDsFromContract(ctx context.Context) ([]*big.Int, error) {
	opts, err := r.buildCallOpts(ctx, nil)
	if err != nil {
		return nil, err
	}

	state, err := r.registry.GetState(opts)
	if err != nil {
		n := "latest"
		if opts.BlockNumber != nil {
			n = fmt.Sprintf("%d", opts.BlockNumber.Int64())
		}

		return nil, fmt.Errorf("%w: failed to get contract state at block number '%s'", err, n)
	}

	ids := make([]*big.Int, 0, int(state.State.NumUpkeeps.Int64()))
	for int64(len(ids)) < state.State.NumUpkeeps.Int64() {
		startIndex := int64(len(ids))
		maxCount := state.State.NumUpkeeps.Int64() - startIndex

		if maxCount == 0 {
			break
		}

		if maxCount > ActiveUpkeepIDBatchSize {
			maxCount = ActiveUpkeepIDBatchSize
		}

		batchIDs, err := r.registry.GetActiveUpkeepIDs(opts, big.NewInt(startIndex), big.NewInt(maxCount))
		if err != nil {
			return nil, fmt.Errorf("%w: failed to get active upkeep IDs from index %d to %d (both inclusive)", err, startIndex, startIndex+maxCount-1)
		}

		ids = append(ids, batchIDs...)
	}

	return ids, nil
}

func (r *EvmRegistry) doCheck(ctx context.Context, mercuryEnabled bool, keys []ocr2keepers.UpkeepKey, chResult chan checkResult) {
	upkeepResults, err := r.checkUpkeeps(ctx, keys)
	if err != nil {
		chResult <- checkResult{
			err: err,
		}
		return
	}

	if mercuryEnabled {
		if r.mercury.cred == nil || !r.mercury.cred.Validate() {
			chResult <- checkResult{
				err: errors.New("mercury credential is empty or not provided but MercuryLookup feature is enabled on registry"),
			}
		}
		upkeepResults, err = r.mercuryLookup(ctx, upkeepResults)
		if err != nil {
			chResult <- checkResult{
				err: err,
			}
			return
		}
	}

	upkeepResults, err = r.simulatePerformUpkeeps(ctx, upkeepResults)
	if err != nil {
		chResult <- checkResult{
			err: err,
		}
		return
	}

	for i, res := range upkeepResults {
		r.mu.RLock()
		up, ok := r.upkeeps[res.ID.String()]
		r.mu.RUnlock()

		if ok {
			upkeepResults[i].ExecuteGas = up.performGasLimit
		}
	}

	chResult <- checkResult{
		ur: upkeepResults,
	}
}

func splitKey(key ocr2keepers.UpkeepKey) (*big.Int, *big.Int, error) {
	var (
		block *big.Int
		id    *big.Int
		ok    bool
	)

	parts := strings.Split(string(key), separator)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("unsplittable key")
	}

	if block, ok = new(big.Int).SetString(parts[0], 10); !ok {
		return nil, nil, fmt.Errorf("could not get block from key")
	}

	if id, ok = new(big.Int).SetString(parts[0], 10); !ok {
		return nil, nil, fmt.Errorf("could not get id from key")
	}

	return block, id, nil
}

// TODO (AUTO-2013): Have better error handling to not return nil results in case of partial errors
func (r *EvmRegistry) checkUpkeeps(ctx context.Context, keys []ocr2keepers.UpkeepKey) ([]EVMAutomationUpkeepResult20, error) {
	var (
		checkReqs    = make([]rpc.BatchElem, len(keys))
		checkResults = make([]*string, len(keys))
	)

	for i, key := range keys {
		block, upkeepId, err := splitKey(key)
		if err != nil {
			return nil, err
		}

		opts, err := r.buildCallOpts(ctx, block)
		if err != nil {
			return nil, err
		}

		payload, err := r.abi.Pack("checkUpkeep", upkeepId)
		if err != nil {
			return nil, err
		}

		var result string
		checkReqs[i] = rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   r.addr.Hex(),
					"data": hexutil.Bytes(payload),
				},
				hexutil.EncodeBig(opts.BlockNumber),
			},
			Result: &result,
		}

		checkResults[i] = &result
	}

	if err := r.client.BatchCallContext(ctx, checkReqs); err != nil {
		return nil, err
	}

	var (
		multiErr error
		results  = make([]EVMAutomationUpkeepResult20, len(keys))
	)

	for i, req := range checkReqs {
		if req.Error != nil {
			r.lggr.Debugf("error encountered for key %s with message '%s' in check", keys[i], req.Error)
			multierr.AppendInto(&multiErr, req.Error)
		} else {
			var err error
			r.lggr.Debugf("UnpackCheckResult key %s checkResult: %s", string(keys[i]), *checkResults[i])
			results[i], err = r.packer.UnpackCheckResult(keys[i], *checkResults[i])
			if err != nil {
				return nil, err
			}
		}
	}

	return results, multiErr
}

// TODO (AUTO-2013): Have better error handling to not return nil results in case of partial errors
func (r *EvmRegistry) simulatePerformUpkeeps(ctx context.Context, checkResults []EVMAutomationUpkeepResult20) ([]EVMAutomationUpkeepResult20, error) {
	var (
		performReqs     = make([]rpc.BatchElem, 0, len(checkResults))
		performResults  = make([]*string, 0, len(checkResults))
		performToKeyIdx = make([]int, 0, len(checkResults))
	)

	for i, checkResult := range checkResults {
		if !checkResult.Eligible {
			continue
		}

		opts, err := r.buildCallOpts(ctx, big.NewInt(int64(checkResult.Block)))
		if err != nil {
			return nil, err
		}

		// Since checkUpkeep is true, simulate perform upkeep to ensure it doesn't revert
		payload, err := r.abi.Pack("simulatePerformUpkeep", checkResult.ID, checkResult.PerformData)
		if err != nil {
			return nil, err
		}

		var result string
		performReqs = append(performReqs, rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   r.addr.Hex(),
					"data": hexutil.Bytes(payload),
				},
				hexutil.EncodeBig(opts.BlockNumber),
			},
			Result: &result,
		})

		performResults = append(performResults, &result)
		performToKeyIdx = append(performToKeyIdx, i)
	}

	if len(performReqs) > 0 {
		if err := r.client.BatchCallContext(ctx, performReqs); err != nil {
			return nil, err
		}
	}

	var multiErr error

	for i, req := range performReqs {
		if req.Error != nil {
			r.lggr.Debugf("error encountered for key %d|%s with message '%s' in simulate perform", checkResults[i].Block, checkResults[i].ID, req.Error)
			multierr.AppendInto(&multiErr, req.Error)
		} else {
			simulatePerformSuccess, err := r.packer.UnpackPerformResult(*performResults[i])
			if err != nil {
				return nil, err
			}

			if !simulatePerformSuccess {
				checkResults[performToKeyIdx[i]].Eligible = false
			}
		}
	}

	return checkResults, multiErr
}

// TODO (AUTO-2013): Have better error handling to not return nil results in case of partial errors
func (r *EvmRegistry) createUpkeepEntries(ctx context.Context, ids []*big.Int) ([]upkeepEntry, error) {
	if len(ids) == 0 {
		return []upkeepEntry{}, nil
	}

	var (
		uReqs    = make([]rpc.BatchElem, len(ids))
		uResults = make([]*string, len(ids))
	)

	for i, id := range ids {
		opts, err := r.buildCallOpts(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get call opts: %s", err)
		}

		payload, err := r.abi.Pack("getUpkeep", id)
		if err != nil {
			return nil, fmt.Errorf("failed to pack id with abi: %s", err)
		}

		var result string
		uReqs[i] = rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   r.addr.Hex(),
					"data": hexutil.Bytes(payload),
				},
				hexutil.EncodeBig(opts.BlockNumber),
			},
			Result: &result,
		}

		uResults[i] = &result
	}

	if err := r.client.BatchCallContext(ctx, uReqs); err != nil {
		return nil, fmt.Errorf("rpc error: %s", err)
	}

	var (
		multiErr error
		results  = make([]upkeepEntry, len(ids))
	)

	for i, req := range uReqs {
		if req.Error != nil {
			r.lggr.Debugf("error encountered for config id %s with message '%s' in get config", ids[i], req.Error)
			multierr.AppendInto(&multiErr, req.Error)
		} else {
			var err error
			results[i], err = r.packer.UnpackUpkeepResult(ids[i], *uResults[i])
			if err != nil {
				return nil, fmt.Errorf("failed to unpack result: %s", err)
			}
		}
	}

	return results, multiErr
}
