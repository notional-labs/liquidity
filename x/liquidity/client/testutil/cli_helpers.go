package testutil

import (
	"fmt"

	"cosmossdk.io/simapp"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	paramscli "github.com/cosmos/cosmos-sdk/x/params/client/cli"
	"github.com/gravity-devs/liquidity/v3/app"

	liquiditycli "github.com/gravity-devs/liquidity/v3/x/liquidity/client/cli"

	dbm "github.com/cometbft/cometbft-db"
)

// NewConfig returns config that defines the necessary testing requirements
// used to bootstrap and start an in-process local testing network.
func NewConfig(dbm *dbm.MemDB) network.Config {
	cfg := network.DefaultConfig(simapp.NewTestNetworkFixture)
	cfg.AppConstructor = func(val network.ValidatorI) servertypes.Application {
		// we build a unique app instance for every validator here
		app := app.NewLiquidityApp(
			val.GetCtx().Logger, dbm, nil, true,
			app.EmptyAppOptions{},
			baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			baseapp.SetChainID(cfg.ChainID),
		)
		return app
	}           // the ABCI application constructor
	cfg.GenesisState = app.ModuleBasics.DefaultGenesis(cfg.Codec) // liquidity genesis state to provide
	return cfg
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
}

// MsgCreatePoolExec creates a transaction for creating liquidity pool.
func MsgCreatePoolExec(clientCtx client.Context, from, poolID, depositCoins string,
	_ ...string,
) (testutil.BufferWriter, error) {
	args := append([]string{
		poolID,
		depositCoins,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "1000000"),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewCreatePoolCmd(), args)
}

// MsgDepositWithinBatchExec creates a transaction to deposit new amounts to the pool.
func MsgDepositWithinBatchExec(clientCtx client.Context, from, poolID, depositCoins string,
	_ ...string,
) (testutil.BufferWriter, error) {
	args := append([]string{
		poolID,
		depositCoins,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "1000000"),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewDepositWithinBatchCmd(), args)
}

// MsgWithdrawWithinBatchExec creates a transaction to withraw pool coin amount from the pool.
func MsgWithdrawWithinBatchExec(clientCtx client.Context, from, poolID, poolCoin string,
	_ ...string,
) (testutil.BufferWriter, error) {
	args := append([]string{
		poolID,
		poolCoin,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "1000000"),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewWithdrawWithinBatchCmd(), args)
}

// MsgSwapWithinBatchExec creates a transaction to swap coins in the pool.
func MsgSwapWithinBatchExec(clientCtx client.Context, from, poolID, swapTypeID,
	offerCoin, demandCoinDenom, orderPrice, swapFeeRate string, _ ...string,
) (testutil.BufferWriter, error) {
	args := append([]string{
		poolID,
		swapTypeID,
		offerCoin,
		demandCoinDenom,
		orderPrice,
		swapFeeRate,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "1000000"),
	}, commonArgs...)

	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, liquiditycli.NewSwapWithinBatchCmd(), args)
}

// MsgParamChangeProposalExec creates a transaction for submitting param change proposal
func MsgParamChangeProposalExec(clientCtx client.Context, from string, file string) (testutil.BufferWriter, error) {
	args := append([]string{
		file,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "1000000"),
	}, commonArgs...)

	paramChangeCmd := paramscli.NewSubmitParamChangeProposalTxCmd()
	flags.AddTxFlagsToCmd(paramChangeCmd)

	return clitestutil.ExecTestCLICmd(clientCtx, paramChangeCmd, args)
}

// MsgVote votes for a proposal
func MsgVote(clientCtx client.Context, from, id, vote string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := append([]string{
		id,
		vote,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "1000000"),
	}, commonArgs...)

	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdWeightedVote(), args)
}
