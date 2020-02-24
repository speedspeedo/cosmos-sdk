package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestRedelegateFromUnbondedValidator(t *testing.T) {
	ctx, _, bk, keeper, _ := CreateTestInput(t, false, 0)
	startTokens := sdk.TokensFromConsensusPower(30)
	startCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), startTokens))

	// add bonded tokens to pool for delegations
	notBondedPool := keeper.GetNotBondedPool(ctx)
	oldNotBonded := bk.GetAllBalances(ctx, notBondedPool.GetAddress())
	err := bk.SetBalances(ctx, notBondedPool.GetAddress(), oldNotBonded.Add(startCoins...))
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	//create a validator with a self-delegation
	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})

	valTokens := sdk.TokensFromConsensusPower(10)
	validator, issuedShares := validator.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())
	validator = TestingUpdateValidator(keeper, ctx, validator, true)
	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
	keeper.SetDelegation(ctx, selfDelegation)

	// create a second delegation to this validator
	keeper.DeleteValidatorByPowerIndex(ctx, validator)
	delTokens := sdk.TokensFromConsensusPower(10)
	validator, issuedShares = validator.AddTokensFromDel(delTokens)
	require.Equal(t, delTokens, issuedShares.RoundInt())
	validator = TestingUpdateValidator(keeper, ctx, validator, true)
	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
	keeper.SetDelegation(ctx, delegation)

	// create a second validator
	validator2 := types.NewValidator(addrVals[1], PKs[1], types.Description{})
	validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())
	validator2 = TestingUpdateValidator(keeper, ctx, validator2, true)
	require.Equal(t, sdk.Bonded, validator2.Status)

	ctx = ctx.WithBlockHeight(10)
	ctx = ctx.WithBlockTime(time.Unix(333, 0))

	// unbond the all self-delegation to put validator in unbonding state
	_, err = keeper.Undelegate(ctx, val0AccAddr, addrVals[0], delTokens.ToDec())
	require.NoError(t, err)

	// end block
	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))

	validator, found := keeper.GetValidator(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, ctx.BlockHeight(), validator.UnbondingHeight)
	params := keeper.GetParams(ctx)
	require.True(t, ctx.BlockHeader().Time.Add(params.UnbondingTime).Equal(validator.UnbondingTime))

	// unbond the validator
	keeper.unbondingToUnbonded(ctx, validator)

	// redelegate some of the delegation's shares
	redelegationTokens := sdk.TokensFromConsensusPower(6)
	_, err = keeper.BeginRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1], redelegationTokens.ToDec())
	require.NoError(t, err)

	// no red should have been found
	red, found := keeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.False(t, found, "%v", red)
}
