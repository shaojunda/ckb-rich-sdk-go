package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ququzone/ckb-sdk-go/types"
)

var (
	NotFound = errors.New("not found")
)

// Client for the Nervos RPC API.
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

	// Close close client
	Close()
}
type client struct {
	c *rpc.Client
}

func Dial(url string) (Client, error) {
	return DialContext(context.Background(), url)
}

func DialContext(ctx context.Context, url string) (Client, error) {
	c, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

func NewClient(c *rpc.Client) Client {
	return &client{c}
}

func (cli *client) Close() {
	cli.c.Close()
}

// Chain RPC

func (cli *client) GetTipBlockNumber(ctx context.Context) (uint64, error) {
	var num hexutil.Uint64
	err := cli.c.CallContext(ctx, &num, "get_tip_block_number")
	if err != nil {
		return 0, err
	}
	return uint64(num), err
}

func (cli *client) GetTipHeader(ctx context.Context) (*types.Header, error) {
	var result header
	err := cli.c.CallContext(ctx, &result, "get_tip_header")
	if err != nil {
		return nil, err
	}
	return toHeader(result), err
}

func (cli *client) GetCurrentEpoch(ctx context.Context) (*types.Epoch, error) {
	var result epoch
	err := cli.c.CallContext(ctx, &result, "get_current_epoch")
	if err != nil {
		return nil, err
	}
	return &types.Epoch{
		CompactTarget: uint64(result.CompactTarget),
		Length:        uint64(result.Length),
		Number:        uint64(result.Number),
		StartNumber:   uint64(result.StartNumber),
	}, err
}

func (cli *client) GetEpochByNumber(ctx context.Context, number uint64) (*types.Epoch, error) {
	var result epoch
	err := cli.c.CallContext(ctx, &result, "get_epoch_by_number", hexutil.Uint64(number))
	if err != nil {
		return nil, err
	}
	return &types.Epoch{
		CompactTarget: uint64(result.CompactTarget),
		Length:        uint64(result.Length),
		Number:        uint64(result.Number),
		StartNumber:   uint64(result.StartNumber),
	}, err
}

func (cli *client) GetBlockHash(ctx context.Context, number uint64) (*types.Hash, error) {
	var result types.Hash

	err := cli.c.CallContext(ctx, &result, "get_block_hash", hexutil.Uint64(number))
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (cli *client) GetBlock(ctx context.Context, hash types.Hash) (*types.Block, error) {
	var raw json.RawMessage

	err := cli.c.CallContext(ctx, &raw, "get_block", hash)
	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, NotFound
	}

	var block block
	if err := json.Unmarshal(raw, &block); err != nil {
		return nil, err
	}

	return &types.Block{
		Header:       toHeader(block.Header),
		Proposals:    toUints(block.Proposals),
		Transactions: toTransactions(block.Transactions),
		Uncles:       toUncles(block.Uncles),
	}, nil
}

func (cli *client) GetHeader(ctx context.Context, hash types.Hash) (*types.Header, error) {
	var result header
	err := cli.c.CallContext(ctx, &result, "get_header", hash)
	if err != nil {
		return nil, err
	}
	return toHeader(result), err
}

func (cli *client) GetHeaderByNumber(ctx context.Context, number uint64) (*types.Header, error) {
	var result header
	err := cli.c.CallContext(ctx, &result, "get_header_by_number", hexutil.Uint64(number))
	if err != nil {
		return nil, err
	}
	return toHeader(result), err
}

func (cli *client) GetCellsByLockHash(ctx context.Context, hash types.Hash, from uint64, to uint64) ([]*types.Cell, error) {
	var result []cell
	err := cli.c.CallContext(ctx, &result, "get_cells_by_lock_hash", hash, hexutil.Uint64(from), hexutil.Uint64(to))
	if err != nil {
		return nil, err
	}
	return toCells(result), err
}

func (cli *client) GetLiveCell(ctx context.Context, point *types.OutPoint, withData bool) (*types.CellWithStatus, error) {
	var result cellWithStatus
	err := cli.c.CallContext(ctx, &result, "get_live_cell", outPoint{
		TxHash: point.TxHash,
		Index:  hexutil.Uint64(point.Index),
	}, true)
	if err != nil {
		return nil, err
	}
	return toCellWithStatus(result), err
}

func (cli *client) GetTransaction(ctx context.Context, hash types.Hash) (*types.TransactionWithStatus, error) {
	var result transactionWithStatus
	err := cli.c.CallContext(ctx, &result, "get_transaction", hash)
	if err != nil {
		return nil, err
	}
	return &types.TransactionWithStatus{
		Transaction: toTransaction(result.Transaction),
		TxStatus: &types.TxStatus{
			BlockHash: result.TxStatus.BlockHash,
			Status:    result.TxStatus.Status,
		},
	}, err
}

func (cli *client) GetCellbaseOutputCapacityDetails(ctx context.Context, hash types.Hash) (*types.BlockReward, error) {
	var result blockReward
	err := cli.c.CallContext(ctx, &result, "get_cellbase_output_capacity_details", hash)
	if err != nil {
		return nil, err
	}

	return &types.BlockReward{
		Primary:        (*big.Int)(&result.Primary),
		ProposalReward: (*big.Int)(&result.ProposalReward),
		Secondary:      (*big.Int)(&result.Secondary),
		Total:          (*big.Int)(&result.Total),
		TxFee:          (*big.Int)(&result.TxFee),
	}, err
}

func (cli *client) GetBlockByNumber(ctx context.Context, number uint64) (*types.Block, error) {
	var raw json.RawMessage

	err := cli.c.CallContext(ctx, &raw, "get_block_by_number", hexutil.Uint64(number))
	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, NotFound
	}

	var block block
	if err := json.Unmarshal(raw, &block); err != nil {
		return nil, err
	}

	return &types.Block{
		Header:       toHeader(block.Header),
		Proposals:    toUints(block.Proposals),
		Transactions: toTransactions(block.Transactions),
		Uncles:       toUncles(block.Uncles),
	}, nil
}
