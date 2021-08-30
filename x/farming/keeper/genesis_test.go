package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	farmingapp "github.com/tendermint/farming/app"
	"github.com/tendermint/farming/x/farming/types"
)

func TestExportGenesis(t *testing.T) {
	app := farmingapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: mustParseRFC3339("2021-08-27T00:00:00Z")})

	creator := farmingapp.AddTestAddrs(app, ctx, 1, sdk.ZeroInt())[0]
	err := farmingapp.FundAccount(app.BankKeeper, ctx, creator, sdk.NewCoins(
		sdk.NewInt64Coin(denom3, 1_000_000_000_000_000)))
	require.NoError(t, err)

	app.FarmingKeeper.SetPlan(ctx, types.NewFixedAmountPlan(
		types.NewBasePlan(
			1,
			"testPlan1",
			types.PlanTypePrivate,
			creator.String(),
			creator.String(),
			sdk.NewDecCoins(sdk.NewDecCoinFromDec(denom1, sdk.NewDec(1))),
			mustParseRFC3339("2021-08-02T00:00:00Z"),
			mustParseRFC3339("2021-09-30T00:00:00Z"),
		),
		sdk.NewCoins(sdk.NewInt64Coin(denom3, 1000000)),
	))

	addrs := farmingapp.AddTestAddrs(app, ctx, 100000, sdk.ZeroInt())
	for _, addr := range addrs {
		err := farmingapp.FundAccount(app.BankKeeper, ctx, addr, initialBalances)
		require.NoError(t, err)

		_, err = app.FarmingKeeper.Stake(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin(denom1, 1000000)))
		require.NoError(t, err)
	}

	app.FarmingKeeper.ProcessQueuedCoins(ctx)
	err = app.FarmingKeeper.DistributeRewards(ctx)
	require.NoError(t, err)

	rewards := app.FarmingKeeper.GetRewardsByFarmer(ctx, addrs[0])
	require.NotEmpty(t, rewards)
	fmt.Println(rewards)

	app.Commit()

	res, err := app.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err)

	genDoc := &tmtypes.GenesisDoc{
		ChainID:  "localnet",
		AppState: res.AppState,
	}

	err = genDoc.ValidateAndComplete()
	require.NoError(t, err)

	genDoc.ConsensusParams.Block.MaxBytes = 200000
	genDoc.ConsensusParams.Block.MaxGas = 40000000
	genDoc.ConsensusParams.Evidence.MaxBytes = 50000

	err = genDoc.SaveAs("genesis.json")
	require.NoError(t, err)
}
