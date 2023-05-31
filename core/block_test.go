package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/peyzor/xchain/types"
	"github.com/stretchr/testify/assert"
)

func TestHeader_Encode_Decode(t *testing.T) {
	hEncode := &Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Timestamp: uint64(time.Now().UnixNano()),
		Height:    10,
		Nonce:     12345,
	}
	hDecode := &Header{}

	buf := &bytes.Buffer{}
	assert.Nil(t, hEncode.EncodeBinary(buf))
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, hEncode, hDecode)
}

func TestBlock_Encode_Decode(t *testing.T) {
	bEncode := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: uint64(time.Now().UnixNano()),
			Height:    10,
			Nonce:     12345,
		},
		Transactions: nil,
	}
	bDecode := &Block{}

	buf := &bytes.Buffer{}
	assert.Nil(t, bEncode.EncodeBinary(buf))
	assert.Nil(t, bDecode.DecodeBinary(buf))
	assert.Equal(t, bEncode, bDecode)
}

func TestBlockHash(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: uint64(time.Now().UnixNano()),
			Height:    10,
			Nonce:     12345,
		},
		Transactions: []Transaction{},
	}

	h := b.Hash()
	assert.False(t, h.IsZero())
}
