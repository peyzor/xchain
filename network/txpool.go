package network

import (
	"fmt"
	"sort"

	"github.com/peyzor/xchain/core"
	"github.com/peyzor/xchain/types"
)

type TxMapSorter struct {
	transactions []*core.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	transactions := make([]*core.Transaction, len(txMap))

	i := 0
	for _, val := range txMap {
		transactions[i] = val
		i++
	}

	s := &TxMapSorter{transactions}
	sort.Sort(s)
	return s
}

func (s *TxMapSorter) Len() int {
	return len(s.transactions)
}

func (s *TxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}

func (s *TxMapSorter) Less(i, j int) bool {
	return s.transactions[i].FirstSeen() < s.transactions[j].FirstSeen()
}

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

func (p *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(p.transactions)
	return s.transactions
}

func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	if p.Has(hash) {
		return fmt.Errorf("transaction (%s) already exists in mempool", hash)
	}

	p.transactions[hash] = tx
	return nil
}

func (p *TxPool) Has(hash types.Hash) bool {
	_, ok := p.transactions[hash]
	return ok
}

func (p *TxPool) Len() int {
	return len(p.transactions)
}

func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*core.Transaction)
}
