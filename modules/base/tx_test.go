package base

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/go-wire/data"

	sdk "github.com/cosmos/cosmos-sdk"
	"github.com/cosmos/cosmos-sdk/stack"
)

func TestEncoding(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	raw := stack.NewRawTx([]byte{0x34, 0xa7})
	// raw2 := stack.NewRawTx([]byte{0x73, 0x86, 0x22})

	cases := []struct {
		Tx sdk.Tx
	}{
		{raw},
		// {NewMultiTx(raw, raw2)},
		{NewChainTx("foobar", 0, raw)},
	}

	for idx, tc := range cases {
		i := strconv.Itoa(idx)
		tx := tc.Tx

		// test json in and out
		js, err := data.ToJSON(tx)
		require.Nil(err, i)
		var jtx sdk.Tx
		err = data.FromJSON(js, &jtx)
		require.Nil(err, i)
		assert.Equal(tx, jtx, i)

		// test wire in and out
		bin, err := data.ToWire(tx)
		require.Nil(err, i)
		var wtx sdk.Tx
		err = data.FromWire(bin, &wtx)
		require.Nil(err, i)
		assert.Equal(tx, wtx, i)
	}
}
