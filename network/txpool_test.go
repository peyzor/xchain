package network

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/peyzor/xchain/core"
	"github.com/stretchr/testify/assert"
)

func TestTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, p.Len(), 0)
}

func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool()
	txData := []byte("foo")
	tx := core.NewTransaction(txData)
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.Len(), 1)

	txx := core.NewTransaction(txData)
	assert.NotNil(t, p.Add(txx))
	assert.Equal(t, p.Len(), 1)

	p.Flush()
	assert.Equal(t, p.Len(), 0)
}

func TestSortTransactions(t *testing.T) {
	p := NewTxPool()

	txLen := 1000
	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("foo %d", i)))
		tx.SetFirstSeen(int64(i * rand.Intn(10000)))
		assert.Nil(t, p.Add(tx))
	}

	assert.Equal(t, txLen, p.Len())

	transactions := p.Transactions()
	for i := 0; i < len(transactions)-1; i++ {
		assert.True(t, transactions[i].FirstSeen() <= transactions[i+1].FirstSeen())
	}
}
