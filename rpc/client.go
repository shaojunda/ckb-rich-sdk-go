package rpc

import (
	"context"
	"github.com/ququzone/ckb-rich-sdk-go/indexer"
	"github.com/ququzone/ckb-sdk-go/rpc"

	"github.com/ququzone/ckb-sdk-go/types"
)

type Client interface {
	////// Chain
	// GetTipBlockNumber returns the number of blocks in the longest blockchain.
	GetTipBlockNumber(ctx context.Context) (uint64, error)

	// GetTipHeader returns the information about the tip header of the longest.
	GetTipHeader(ctx context.Context) (*types.Header, error)

	// GetCurrentEpoch returns the information about the current epoch.
	GetCurrentEpoch(ctx context.Context) (*types.Epoch, error)

	// GetEpochByNumber return the information corresponding the given epoch number.
	GetEpochByNumber(ctx context.Context, number uint64) (*types.Epoch, error)

	// GetBlockHash returns the hash of a block in the best-block-chain by block number; block of No.0 is the genesis block.
	GetBlockHash(ctx context.Context, number uint64) (*types.Hash, error)

	// GetBlock returns the information about a block by hash.
	GetBlock(ctx context.Context, hash types.Hash) (*types.Block, error)

	// GetHeader returns the information about a block header by hash.
	GetHeader(ctx context.Context, hash types.Hash) (*types.Header, error)

	// GetHeaderByNumber returns the information about a block header by block number.
	GetHeaderByNumber(ctx context.Context, number uint64) (*types.Header, error)

	// GetCellsByLockHash returns the information about cells collection by the hash of lock script.
	GetCellsByLockHash(ctx context.Context, hash types.Hash, from uint64, to uint64) ([]*types.Cell, error)

	// GetLiveCell returns the information about a cell by out_point if it is live.
	// If second with_data argument set to true, will return cell data and data_hash if it is live.
	GetLiveCell(ctx context.Context, outPoint *types.OutPoint, withData bool) (*types.CellWithStatus, error)

	// GetTransaction returns the information about a transaction requested by transaction hash.
	GetTransaction(ctx context.Context, hash types.Hash) (*types.TransactionWithStatus, error)

	// GetCellbaseOutputCapacityDetails returns each component of the created CKB in this block's cellbase,
	// which is issued to a block N - 1 - ProposalWindow.farthest, where this block's height is N.
	GetCellbaseOutputCapacityDetails(ctx context.Context, hash types.Hash) (*types.BlockReward, error)

	// GetBlockByNumber get block by number
	GetBlockByNumber(ctx context.Context, number uint64) (*types.Block, error)

	////// Experiment
	// DryRunTransaction dry run transaction and return the execution cycles.
	// This method will not check the transaction validity,
	// but only run the lock script and type script and then return the execution cycles.
	// Used to debug transaction scripts and query how many cycles the scripts consume.
	DryRunTransaction(ctx context.Context, transaction *types.Transaction) (*types.DryRunTransactionResult, error)

	// CalculateDaoMaximumWithdraw calculate the maximum withdraw one can get, given a referenced DAO cell, and a withdraw block hash.
	CalculateDaoMaximumWithdraw(ctx context.Context, point *types.OutPoint, hash types.Hash) (uint64, error)

	// EstimateFeeRate Estimate a fee rate (capacity/KB) for a transaction that to be committed in expect blocks.
	EstimateFeeRate(ctx context.Context, blocks uint64) (*types.EstimateFeeRateResult, error)

	////// Indexer
	// IndexLockHash create index for live cells and transactions by the hash of lock script.
	IndexLockHash(ctx context.Context, lockHash types.Hash, indexFrom uint64) (*types.LockHashIndexState, error)

	// GetLockHashIndexStates Get lock hash index states.
	GetLockHashIndexStates(ctx context.Context) ([]*types.LockHashIndexState, error)

	// GetLiveCellsByLockHash returns the live cells collection by the hash of lock script.
	GetLiveCellsByLockHash(ctx context.Context, lockHash types.Hash, page uint, per uint, reverseOrder bool) ([]*types.LiveCell, error)

	// GetTransactionsByLockHash returns the transactions collection by the hash of lock script.
	// Returns empty array when the lock_hash has not been indexed yet.
	GetTransactionsByLockHash(ctx context.Context, lockHash types.Hash, page uint, per uint, reverseOrder bool) ([]*types.CellTransaction, error)

	// DeindexLockHash Remove index for live cells and transactions by the hash of lock script.
	DeindexLockHash(ctx context.Context, lockHash types.Hash) error

	////// Net
	// LocalNodeInfo returns the local node information.
	LocalNodeInfo(ctx context.Context) (*types.Node, error)

	// GetPeers returns the connected peers information.
	GetPeers(ctx context.Context) ([]*types.Node, error)

	// GetBannedAddresses returns all banned IPs/Subnets.
	GetBannedAddresses(ctx context.Context) ([]*types.BannedAddress, error)

	// SetBan insert or delete an IP/Subnet from the banned list
	SetBan(ctx context.Context, address string, command string, banTime uint64, absolute bool, reason string) error

	////// Pool
	// SendTransaction send new transaction into transaction pool.
	SendTransaction(ctx context.Context, tx *types.Transaction) (*types.Hash, error)

	// SendTransactionNoneValidation send new transaction into transaction pool skipping outputs validation.
	SendTransactionNoneValidation(ctx context.Context, tx *types.Transaction) (*types.Hash, error)

	// TxPoolInfo return the transaction pool information
	TxPoolInfo(ctx context.Context) (*types.TxPoolInfo, error)

	////// Stats
	// GetBlockchainInfo return state info of blockchain
	GetBlockchainInfo(ctx context.Context) (*types.BlockchainInfo, error)

	////// Batch
	BatchTransactions(ctx context.Context, batch []types.BatchTransactionItem) error

	///// ckb-indexer
	//GetTip returns the latest height processed by indexer
	GetTip(ctx context.Context) (*indexer.TipHeader, error)

	// GetCells returns the live cells collection by the lock or type script.
	GetCells(ctx context.Context, searchKey *indexer.SearchKey, order indexer.SearchOrder, limit uint64, afterCursor string) (*indexer.LiveCells, error)

	// GetTransactions returns the transactions collection by the lock or type script.
	GetTransactions(ctx context.Context, searchKey *indexer.SearchKey, order indexer.SearchOrder, limit uint64, afterCursor string) (*indexer.Transactions, error)

	// Close close client
	Close()
}

type client struct {
	ckb     rpc.Client
	indexer indexer.Client
}

func Dial(ckbUrl string, indexUrl string) (Client, error) {
	ckb, err := rpc.Dial(ckbUrl)
	if err != nil {
		return nil, err
	}
	index, err := indexer.Dial(indexUrl)
	if err != nil {
		return nil, err
	}

	return &client{
		ckb:     ckb,
		indexer: index,
	}, nil
}

func (cli *client) Close() {
	cli.ckb.Close()
	cli.indexer.Close()
}

func (cli *client) GetTipBlockNumber(ctx context.Context) (uint64, error) {
	return cli.ckb.GetTipBlockNumber(ctx)
}

func (cli *client) GetTipHeader(ctx context.Context) (*types.Header, error) {
	return cli.ckb.GetTipHeader(ctx)
}

func (cli *client) GetCurrentEpoch(ctx context.Context) (*types.Epoch, error) {
	return cli.ckb.GetCurrentEpoch(ctx)
}

func (cli *client) GetEpochByNumber(ctx context.Context, number uint64) (*types.Epoch, error) {
	return cli.ckb.GetEpochByNumber(ctx, number)
}

func (cli *client) GetBlockHash(ctx context.Context, number uint64) (*types.Hash, error) {
	return cli.ckb.GetBlockHash(ctx, number)
}

func (cli *client) GetBlock(ctx context.Context, hash types.Hash) (*types.Block, error) {
	return cli.ckb.GetBlock(ctx, hash)
}

func (cli *client) GetHeader(ctx context.Context, hash types.Hash) (*types.Header, error) {
	return cli.ckb.GetHeader(ctx, hash)
}

func (cli *client) GetHeaderByNumber(ctx context.Context, number uint64) (*types.Header, error) {
	return cli.ckb.GetHeaderByNumber(ctx, number)
}

func (cli *client) GetCellsByLockHash(ctx context.Context, hash types.Hash, from uint64, to uint64) ([]*types.Cell, error) {
	return cli.ckb.GetCellsByLockHash(ctx, hash, from, to)
}

func (cli *client) GetLiveCell(ctx context.Context, outPoint *types.OutPoint, withData bool) (*types.CellWithStatus, error) {
	return cli.ckb.GetLiveCell(ctx, outPoint, withData)
}

func (cli *client) GetTransaction(ctx context.Context, hash types.Hash) (*types.TransactionWithStatus, error) {
	return cli.ckb.GetTransaction(ctx, hash)
}

func (cli *client) GetCellbaseOutputCapacityDetails(ctx context.Context, hash types.Hash) (*types.BlockReward, error) {
	return cli.ckb.GetCellbaseOutputCapacityDetails(ctx, hash)
}

func (cli *client) GetBlockByNumber(ctx context.Context, number uint64) (*types.Block, error) {
	return cli.ckb.GetBlockByNumber(ctx, number)
}

func (cli *client) DryRunTransaction(ctx context.Context, transaction *types.Transaction) (*types.DryRunTransactionResult, error) {
	return cli.ckb.DryRunTransaction(ctx, transaction)
}

func (cli *client) CalculateDaoMaximumWithdraw(ctx context.Context, point *types.OutPoint, hash types.Hash) (uint64, error) {
	return cli.ckb.CalculateDaoMaximumWithdraw(ctx, point, hash)
}

func (cli *client) EstimateFeeRate(ctx context.Context, blocks uint64) (*types.EstimateFeeRateResult, error) {
	return cli.ckb.EstimateFeeRate(ctx, blocks)
}

func (cli *client) IndexLockHash(ctx context.Context, lockHash types.Hash, indexFrom uint64) (*types.LockHashIndexState, error) {
	return cli.ckb.IndexLockHash(ctx, lockHash, indexFrom)
}

func (cli *client) GetLockHashIndexStates(ctx context.Context) ([]*types.LockHashIndexState, error) {
	return cli.ckb.GetLockHashIndexStates(ctx)
}

func (cli *client) GetLiveCellsByLockHash(ctx context.Context, lockHash types.Hash, page uint, per uint, reverseOrder bool) ([]*types.LiveCell, error) {
	return cli.ckb.GetLiveCellsByLockHash(ctx, lockHash, page, per, reverseOrder)
}

func (cli *client) GetTransactionsByLockHash(ctx context.Context, lockHash types.Hash, page uint, per uint, reverseOrder bool) ([]*types.CellTransaction, error) {
	return cli.ckb.GetTransactionsByLockHash(ctx, lockHash, page, per, reverseOrder)
}

func (cli *client) DeindexLockHash(ctx context.Context, lockHash types.Hash) error {
	return cli.ckb.DeindexLockHash(ctx, lockHash)
}

func (cli *client) LocalNodeInfo(ctx context.Context) (*types.Node, error) {
	return cli.ckb.LocalNodeInfo(ctx)
}

func (cli *client) GetPeers(ctx context.Context) ([]*types.Node, error) {
	return cli.ckb.GetPeers(ctx)
}

func (cli *client) GetBannedAddresses(ctx context.Context) ([]*types.BannedAddress, error) {
	return cli.ckb.GetBannedAddresses(ctx)
}

func (cli *client) SetBan(ctx context.Context, address string, command string, banTime uint64, absolute bool, reason string) error {
	return cli.ckb.SetBan(ctx, address, command, banTime, absolute, reason)
}

func (cli *client) SendTransaction(ctx context.Context, tx *types.Transaction) (*types.Hash, error) {
	return cli.ckb.SendTransaction(ctx, tx)
}

func (cli *client) SendTransactionNoneValidation(ctx context.Context, tx *types.Transaction) (*types.Hash, error) {
	return cli.ckb.SendTransactionNoneValidation(ctx, tx)
}

func (cli *client) TxPoolInfo(ctx context.Context) (*types.TxPoolInfo, error) {
	return cli.ckb.TxPoolInfo(ctx)
}

func (cli *client) GetBlockchainInfo(ctx context.Context) (*types.BlockchainInfo, error) {
	return cli.ckb.GetBlockchainInfo(ctx)
}

func (cli *client) BatchTransactions(ctx context.Context, batch []types.BatchTransactionItem) error {
	return cli.ckb.BatchTransactions(ctx, batch)
}

func (cli *client) GetTip(ctx context.Context) (*indexer.TipHeader, error) {
	return cli.indexer.GetTip(ctx)
}

func (cli *client) GetCells(ctx context.Context, searchKey *indexer.SearchKey, order indexer.SearchOrder, limit uint64, afterCursor string) (*indexer.LiveCells, error) {
	return cli.indexer.GetCells(ctx, searchKey, order, limit, afterCursor)
}

func (cli *client) GetTransactions(ctx context.Context, searchKey *indexer.SearchKey, order indexer.SearchOrder, limit uint64, afterCursor string) (*indexer.Transactions, error) {
	return cli.indexer.GetTransactions(ctx, searchKey, order, limit, afterCursor)
}
