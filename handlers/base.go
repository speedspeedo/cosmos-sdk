package handlers

import (
	crypto "github.com/tendermint/go-crypto"

	"github.com/tendermint/basecoin"
	"github.com/tendermint/basecoin/errors"
	"github.com/tendermint/basecoin/types"
)

type SignedHandler struct {
	AllowMultiSig bool
	Inner         basecoin.Handler
}

func (h SignedHandler) Next() basecoin.Handler {
	return h.Inner
}

var _ basecoin.Handler = SignedHandler{}

type Signed interface {
	basecoin.TxLayer
	Signers() ([]crypto.PubKey, error)
}

func (h SignedHandler) CheckTx(ctx basecoin.Context, store types.KVStore, tx basecoin.Tx) (res basecoin.Result, err error) {
	var sigs []crypto.PubKey

	stx, ok := tx.Unwrap().(Signed)
	if !ok {
		return res, errors.Unauthorized()
	}

	sigs, err = stx.Signers()
	if err != nil {
		return res, err
	}

	// add the signers to the context and continue
	ctx2 := ctx.AddSigners(sigs...)
	return h.Next().CheckTx(ctx2, store, stx.Next())
}

func (h SignedHandler) DeliverTx(ctx basecoin.Context, store types.KVStore, tx basecoin.Tx) (res basecoin.Result, err error) {
	var sigs []crypto.PubKey

	stx, ok := tx.Unwrap().(Signed)
	if !ok {
		return res, errors.Unauthorized()
	}

	sigs, err = stx.Signers()
	if err != nil {
		return res, err
	}

	// add the signers to the context and continue
	ctx2 := ctx.AddSigners(sigs...)
	return h.Next().DeliverTx(ctx2, store, stx.Next())
}
