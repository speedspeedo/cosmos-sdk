package crisis

import (
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/core/appmodule"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	modulev1 "cosmossdk.io/api/cosmos/crisis/module/v1"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis/client/cli"
	"github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	"github.com/cosmos/cosmos-sdk/x/crisis/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// Module init related flags
const (
	FlagSkipGenesisInvariants = "x-crisis-skip-assert-invariants"
)

// AppModuleBasic defines the basic application module used by the crisis module.
type AppModuleBasic struct{}

// Name returns the crisis module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the crisis module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the crisis
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the crisis module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return types.ValidateGenesis(&data)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the capability module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *gwruntime.ServeMux) {}

// GetTxCmd returns the root tx command for the crisis module.
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns no root query command for the crisis module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command { return nil }

// RegisterInterfaces registers interfaces and implementations of the crisis
// module.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// AppModule implements an application module for the crisis module.
type AppModule struct {
	AppModuleBasic

	// NOTE: We store a reference to the keeper here so that after a module
	// manager is created, the invariants can be properly registered and
	// executed.
	keeper *keeper.Keeper

	skipGenesisInvariants bool
}

// NewAppModule creates a new AppModule object. If initChainAssertInvariants is set,
// we will call keeper.AssertInvariants during InitGenesis (it may take a significant time)
// - which doesn't impact the chain security unless 66+% of validators have a wrongly
// modified genesis file.
func NewAppModule(keeper *keeper.Keeper, skipGenesisInvariants bool) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,

		skipGenesisInvariants: skipGenesisInvariants,
	}
}

// AddModuleInitFlags implements servertypes.ModuleInitFlags interface.
func AddModuleInitFlags(startCmd *cobra.Command) {
	startCmd.Flags().Bool(FlagSkipGenesisInvariants, false, "Skip x/crisis invariants check on startup")
}

// Name returns the crisis module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants performs a no-op.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Deprecated: Route returns the message routing key for the crisis module.
func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// QuerierRoute returns no querier route.
func (AppModule) QuerierRoute() string { return "" }

// LegacyQuerierHandler returns no sdk.Querier.
func (AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier { return nil }

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), am.keeper)
}

// InitGenesis performs genesis initialization for the crisis module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	start := time.Now()
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	telemetry.MeasureSince(start, "InitGenesis", "crisis", "unmarshal")

	am.keeper.InitGenesis(ctx, &genesisState)
	if !am.skipGenesisInvariants {
		am.keeper.AssertInvariants(ctx)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the crisis
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock performs a no-op.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the crisis module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, *am.keeper)
	return []abci.ValidatorUpdate{}
}

// New App Wiring Setup

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(provideModuleBasic, provideModule),
	)
}

type crisisInputs struct {
	depinject.In

	Config     *modulev1.Module
	AppOpts    servertypes.AppOptions `optional:"true"`
	Subspace   paramstypes.Subspace
	BankKeeper types.SupplyKeeper
}

type crisisOutputs struct {
	depinject.Out

	Module       runtime.AppModuleWrapper
	CrisisKeeper *keeper.Keeper
}

func provideModuleBasic() runtime.AppModuleBasicWrapper {
	return runtime.WrapAppModuleBasic(AppModuleBasic{})
}

func provideModule(in crisisInputs) crisisOutputs {
	invalidCheckPeriod := cast.ToUint(in.AppOpts.Get(server.FlagInvCheckPeriod))

	feeCollectorName := in.Config.FeeCollectorName
	if feeCollectorName == "" {
		feeCollectorName = authtypes.FeeCollectorName
	}

	k := keeper.NewKeeper(
		in.Subspace,
		invalidCheckPeriod,
		in.BankKeeper,
		feeCollectorName,
	)

	skipGenesisInvariants := cast.ToBool(in.AppOpts.Get(FlagSkipGenesisInvariants))

	m := NewAppModule(k, skipGenesisInvariants)

	return crisisOutputs{CrisisKeeper: k, Module: runtime.WrapAppModule(m)}
}
