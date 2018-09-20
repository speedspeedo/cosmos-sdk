package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Allocate fees handles distribution of the collected fees
func (k Keeper) AllocateFees(ctx sdk.Context) {

	// get the proposer of this block
	proposerConsAddr := k.GetProposerConsAddr(ctx)
	proserValidator := k.stakeKeeper.GetValidatorFromConsAddr(ctx, proposerConsAddr)
	proposerDist := k.GetValidatorDistInfo(ctx, proserValidator.GetOperator())

	// get the fees which have been getting collected through all the
	// transactions in the block
	feesCollected := k.feeCollectionKeeper.GetCollectedFees(ctx)
	feesCollectedDec := types.NewDecCoins(feesCollected)

	// allocated rewards to proposer
	bondedTokens := k.stakeKeeper.TotalPower(ctx)
	sumPowerPrecommitValidators := sdk.NewDec(1) // XXX TODO actually calculate this
	proposerMultiplier := sdk.NewDecWithPrec(1, 2).Add(sdk.NewDecWithPrec(4, 2).Mul(
		sumPowerPrecommitValidators).Quo(bondedTokens))
	proposerReward := feesCollectedDec.Mul(proposerMultiplier)

	// apply commission
	commission := proposerReward.Mul(proserValidator.GetCommission())
	remaining := proposerReward.Mul(sdk.OneDec().Sub(proserValidator.GetCommission()))
	proposerDist.PoolCommission = proposerDist.PoolCommission.Plus(commission)
	proposerDist.Pool = proposerDist.Pool.Plus(remaining)

	// allocate community funding
	communityTax := k.GetCommunityTax(ctx)
	communityFunding := feesCollectedDec.Mul(communityTax)
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Plus(communityFunding)

	// set the global pool within the distribution module
	poolReceived := feesCollectedDec.Mul(sdk.OneDec().Sub(proposerMultiplier).Sub(communityTax))
	feePool.Pool = feePool.Pool.Plus(poolReceived)

	k.SetValidatorDistInfo(ctx, proposerDist)
	k.SetFeePool(ctx, feePool)

	// clear the now distributed fees
	k.feeCollectionKeeper.ClearCollectedFees(ctx)
}
