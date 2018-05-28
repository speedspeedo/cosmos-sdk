package slashing

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

func TestHandleDoubleSign(t *testing.T) {
	ctx, ck, sk, keeper := createTestInput(t)
	addr, val, amt := addrs[0], pks[0], int64(100)
	got := stake.NewHandler(sk)(ctx, newTestMsgDeclareCandidacy(addr, val, amt))
	require.True(t, got.IsOK())
	_ = sk.Tick(ctx)
	require.Equal(t, ck.GetCoins(ctx, addr), sdk.Coins{{sk.GetParams(ctx).BondDenom, initCoins - amt}})
	require.Equal(t, sdk.NewRat(amt), sk.Validator(ctx, addr).GetPower())
	keeper.handleDoubleSign(ctx, 0, 0, val) // double sign less than max age
	require.Equal(t, sdk.NewRat(amt).Mul(sdk.NewRat(19).Quo(sdk.NewRat(20))), sk.Validator(ctx, addr).GetPower())
	ctx = ctx.WithBlockHeader(abci.Header{Time: 300})
	keeper.handleDoubleSign(ctx, 0, 0, val) // double sign past max age
	require.Equal(t, sdk.NewRat(amt).Mul(sdk.NewRat(19).Quo(sdk.NewRat(20))), sk.Validator(ctx, addr).GetPower())
}

func TestHandleAbsentValidator(t *testing.T) {
	ctx, ck, sk, keeper := createTestInput(t)
	addr, val, amt := addrs[0], pks[0], int64(100)
	sh := stake.NewHandler(sk)
	slh := NewHandler(keeper)
	got := sh(ctx, newTestMsgDeclareCandidacy(addr, val, amt))
	require.True(t, got.IsOK())
	_ = sk.Tick(ctx)
	require.Equal(t, ck.GetCoins(ctx, addr), sdk.Coins{{sk.GetParams(ctx).BondDenom, initCoins - amt}})
	require.Equal(t, sdk.NewRat(amt), sk.Validator(ctx, addr).GetPower())
	info, found := keeper.getValidatorSigningInfo(ctx, val.Address())
	require.False(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	require.Equal(t, int64(0), info.SignedBlocksCounter)
	height := int64(0)
	// 1000 blocks OK
	for ; height < 1000; height++ {
		ctx = ctx.WithBlockHeight(height)
		keeper.handleValidatorSignature(ctx, val, true)
	}
	info, found = keeper.getValidatorSigningInfo(ctx, val.Address())
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	require.Equal(t, SignedBlocksWindow, info.SignedBlocksCounter)
	// 50 blocks missed
	for ; height < 1050; height++ {
		ctx = ctx.WithBlockHeight(height)
		keeper.handleValidatorSignature(ctx, val, false)
	}
	info, found = keeper.getValidatorSigningInfo(ctx, val.Address())
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	require.Equal(t, SignedBlocksWindow-50, info.SignedBlocksCounter)
	// should be bonded still
	validator := sk.ValidatorByPubKey(ctx, val)
	require.Equal(t, sdk.Bonded, validator.GetStatus())
	pool := sk.GetPool(ctx)
	require.Equal(t, int64(100), pool.BondedTokens)
	// 51st block missed
	ctx = ctx.WithBlockHeight(height)
	keeper.handleValidatorSignature(ctx, val, false)
	info, found = keeper.getValidatorSigningInfo(ctx, val.Address())
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	require.Equal(t, SignedBlocksWindow-51, info.SignedBlocksCounter)
	height++
	// should have been revoked
	validator = sk.ValidatorByPubKey(ctx, val)
	require.Equal(t, sdk.Unbonded, validator.GetStatus())
	got = slh(ctx, NewMsgUnrevoke(addr))
	require.False(t, got.IsOK()) // should fail prior to jail expiration
	ctx = ctx.WithBlockHeader(abci.Header{Time: int64(86400 * 2)})
	got = slh(ctx, NewMsgUnrevoke(addr))
	require.True(t, got.IsOK()) // should succeed after jail expiration
	validator = sk.ValidatorByPubKey(ctx, val)
	require.Equal(t, sdk.Bonded, validator.GetStatus())
	// should have been slashed
	pool = sk.GetPool(ctx)
	require.Equal(t, int64(99), pool.BondedTokens)
}
