package types

import (
	"math/big"
)

type ScriptHashType string
type DepType string
type TransactionStatus string

const (
	Data ScriptHashType = "data"
	Type ScriptHashType = "type"

	Code     DepType = "code"
	DepGroup DepType = "dep_group"

	Pending   TransactionStatus = "pending"
	Proposed  TransactionStatus = "proposed"
	Committed TransactionStatus = "committed"
)

type Epoch struct {
	CompactTarget uint64 `json:"compact_target"`
	Length        uint64 `json:"length"`
	Number        uint64 `json:"number"`
	StartNumber   uint64 `json:"start_number"`
}

type Header struct {
	CompactTarget    uint64   `json:"compact_target"`
	Dao              Hash     `json:"dao"`
	Epoch            uint64   `json:"epoch"`
	Hash             Hash     `json:"hash"`
	Nonce            *big.Int `json:"nonce"`
	Number           uint64   `json:"number"`
	ParentHash       Hash     `json:"parent_hash"`
	ProposalsHash    Hash     `json:"proposals_hash"`
	Timestamp        uint64   `json:"timestamp"`
	TransactionsRoot Hash     `json:"transactions_root"`
	UnclesHash       Hash     `json:"uncles_hash"`
	Version          uint     `json:"version"`
}

type OutPoint struct {
	TxHash Hash   `json:"tx_hash"`
	Index  uint64 `json:"index"`
}

type CellDep struct {
	OutPoint *OutPoint `json:"out_point"`
	DepType  DepType   `json:"dep_type"`
}

type Script struct {
	CodeHash Hash           `json:"code_hash"`
	HashType ScriptHashType `json:"hash_type"`
	Args     []byte         `json:"args"`
}

type CellInput struct {
	Since          uint64    `json:"since"`
	PreviousOutput *OutPoint `json:"previous_output"`
}

type CellOutput struct {
	Capacity *big.Int `json:"capacity"`
	Lock     *Script  `json:"lock"`
	Type     *Script  `json:"type"`
}

type Transaction struct {
	Version     uint          `json:"version"`
	Hash        Hash          `json:"hash"`
	CellDeps    []*CellDep    `json:"cell_deps"`
	HeaderDeps  []Hash        `json:"header_deps"`
	Inputs      []*CellInput  `json:"inputs"`
	Outputs     []*CellOutput `json:"outputs"`
	OutputsData [][]byte      `json:"outputs_data"`
	Witnesses   [][]byte      `json:"witnesses"`
}

type UncleBlock struct {
	Header    *Header `json:"header"`
	Proposals []uint  `json:"proposals"`
}

type Block struct {
	Header       *Header        `json:"header"`
	Proposals    []uint         `json:"proposals"`
	Transactions []*Transaction `json:"transactions"`
	Uncles       []*UncleBlock  `json:"uncles"`
}

type Cell struct {
	BlockHash Hash      `json:"block_hash"`
	Capacity  *big.Int  `json:"capacity"`
	Lock      *Script   `json:"lock"`
	OutPoint  *OutPoint `json:"out_point"`
	Type      *Script   `json:"type"`
}

type CellData struct {
	Content []byte `json:"content"`
	Hash    Hash   `json:"hash"`
}

type CellInfo struct {
	Data   *CellData   `json:"data"`
	Output *CellOutput `json:"output"`
}

type CellWithStatus struct {
	Cell   *CellInfo `json:"cell"`
	Status string    `json:"status"`
}

type TxStatus struct {
	BlockHash *Hash             `json:"block_hash"`
	Status    TransactionStatus `json:"status"`
}

type TransactionWithStatus struct {
	Transaction *Transaction `json:"transaction"`
	TxStatus    *TxStatus    `json:"tx_status"`
}

type BlockReward struct {
	Primary        *big.Int `json:"primary"`
	ProposalReward *big.Int `json:"proposal_reward"`
	Secondary      *big.Int `json:"secondary"`
	Total          *big.Int `json:"total"`
	TxFee          *big.Int `json:"tx_fee"`
}
