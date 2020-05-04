package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ququzone/ckb-rich-sdk-go/indexer"
	"github.com/ququzone/ckb-rich-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/types"
)

func main() {
	c, err := rpc.Dial("http://localhost:8117/rpc", "http://localhost:8117/indexer")
	if err != nil {
		log.Fatalf("dial rpc error: %v", err)
	}

	blockNumber, _ := c.GetTipBlockNumber(context.Background())
	fmt.Println(blockNumber)

	args, _ := hex.DecodeString("c2baa1d5b45a3ad6452b9c98ad8e2cc52e5123c7")
	searchKey := &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
			HashType: types.HashTypeType,
			Args:     args,
		},
		ScriptType: indexer.ScriptTypeLock,
	}

	// first page
	liveCells, _ := c.GetCells(context.Background(), searchKey, indexer.SearchOrderAsc, 100, "")

	fmt.Println(liveCells.LastCursor)
	fmt.Println(len(liveCells.Objects))
	fmt.Println(liveCells.Objects[0].OutPoint.TxHash)
	fmt.Println(liveCells.Objects[0].OutPoint.Index)

	// search next page
	liveCells, _ = c.GetCells(context.Background(), searchKey, indexer.SearchOrderAsc, 100, liveCells.LastCursor)
	fmt.Println(liveCells.LastCursor)
	fmt.Println(len(liveCells.Objects))
	fmt.Println(liveCells.Objects[0].OutPoint.TxHash)
	fmt.Println(liveCells.Objects[0].OutPoint.Index)

	fmt.Println("-------------------------------------")
	transactions, _ := c.GetTransactions(context.Background(), searchKey, indexer.SearchOrderAsc, 100, "")
	fmt.Println(transactions.LastCursor)
	fmt.Println(len(transactions.Objects))
	fmt.Println(transactions.Objects[0].TxHash)
	fmt.Println(transactions.Objects[0].IoType)
	transactions, _ = c.GetTransactions(context.Background(), searchKey, indexer.SearchOrderAsc, 100, transactions.LastCursor)
	fmt.Println(transactions.LastCursor)
	fmt.Println(len(transactions.Objects))
	fmt.Println(transactions.Objects[0].TxHash)
	fmt.Println(transactions.Objects[0].IoType)
}
