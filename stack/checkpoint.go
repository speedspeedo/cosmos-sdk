package stack

import (
	"github.com/tendermint/basecoin"
	"github.com/tendermint/basecoin/state"
)

//nolint
const (
	NameCheckpoint = "check"
)

// Checkpoint isolates all data store below this
type Checkpoint struct {
	OnCheck   bool
	OnDeliver bool
	PassOption
}

// Name of the module - fulfills Middleware interface
func (Checkpoint) Name() string {
	return NameCheckpoint
}

var _ Middleware = Checkpoint{}

// CheckTx reverts all data changes if there was an error
func (c Checkpoint) CheckTx(ctx basecoin.Context, store state.SimpleDB, tx basecoin.Tx, next basecoin.Checker) (res basecoin.CheckResult, err error) {
	if !c.OnCheck {
		return next.CheckTx(ctx, store, tx)
	}
	ps := store.Checkpoint()
	res, err = next.CheckTx(ctx, ps, tx)
	if err == nil {
		err = store.Commit(ps)
	}
	return res, err
}

// DeliverTx reverts all data changes if there was an error
func (c Checkpoint) DeliverTx(ctx basecoin.Context, store state.SimpleDB, tx basecoin.Tx, next basecoin.Deliver) (res basecoin.DeliverResult, err error) {
	if !c.OnDeliver {
		return next.DeliverTx(ctx, store, tx)
	}
	ps := store.Checkpoint()
	res, err = next.DeliverTx(ctx, ps, tx)
	if err == nil {
		err = store.Commit(ps)
	}
	return res, err
}
