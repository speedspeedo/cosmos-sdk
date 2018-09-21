package stake

import (
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake/types"
	"github.com/pkg/errors"
)

// InitGenesis sets the pool and parameters for the provided keeper and
// initializes the IntraTxCounter. For each validator in data, it sets that
// validator in the keeper along with manually setting the indexes. In
// addition, it also sets any delegations found in data. Finally, it updates
// the bonded validators.
// Returns final validator set after applying all declaration and delegations
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.Validator, err error) {
	keeper.SetPool(ctx, data.Pool)
	keeper.SetNewParams(ctx, data.Params)
	keeper.InitIntraTxCounter(ctx)

	for i, validator := range data.Validators {
		validator.BondIntraTxCounter = int16(i) // set the intra-tx counter to the order the validators are presented
		keeper.SetValidator(ctx, validator)

		if validator.Tokens.IsZero() {
			return res, errors.Errorf("genesis validator cannot have zero pool shares, validator: %v", validator)
		}
		if validator.DelegatorShares.IsZero() {
			return res, errors.Errorf("genesis validator cannot have zero delegator shares, validator: %v", validator)
		}

		// Manually set indexes for the first time
		keeper.SetValidatorByConsAddr(ctx, validator)
		keeper.SetValidatorByPowerIndex(ctx, validator, data.Pool)

		if validator.Status == sdk.Bonded {
			keeper.SetValidatorBondedIndex(ctx, validator)
		}
	}

	for _, bond := range data.Bonds {
		keeper.SetDelegation(ctx, bond)
	}

	keeper.UpdateBondedValidatorsFull(ctx)

	vals := keeper.GetValidatorsBonded(ctx)
	res = make([]abci.Validator, len(vals))
	for i, val := range vals {
		res[i] = sdk.ABCIValidator(val)
	}
	return
}

// WriteGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, validators, and bonds found in
// the keeper.
func WriteGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	pool := keeper.GetPool(ctx)
	params := keeper.GetParams(ctx)
	validators := keeper.GetAllValidators(ctx)
	bonds := keeper.GetAllDelegations(ctx)

	return types.GenesisState{
		Pool:       pool,
		Params:     params,
		Validators: validators,
		Bonds:      bonds,
	}
}

// WriteValidators returns a slice of bonded genesis validators.
func WriteValidators(ctx sdk.Context, keeper Keeper) (vals []tmtypes.GenesisValidator) {
	keeper.IterateValidatorsBonded(ctx, func(_ int64, validator sdk.Validator) (stop bool) {
		vals = append(vals, tmtypes.GenesisValidator{
			PubKey: validator.GetConsPubKey(),
			Power:  validator.GetPower().RoundInt64(),
			Name:   validator.GetMoniker(),
		})

		return false
	})

	return
}
