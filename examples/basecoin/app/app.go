package app

import (
	"fmt"
	"os"

	apm "github.com/cosmos/cosmos-sdk/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/abci/server"
	"github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
)

const appName = "BasecoinApp"

type BasecoinApp struct {
	*apm.App
	cdc        *wire.Codec
	multiStore sdk.CommitMultiStore

	// The key to access the substores.
	mainStoreKey *sdk.KVStoreKey
	ibcStoreKey  *sdk.KVStoreKey

	// Additional stores:
	accStore sdk.AccountStore
}

// TODO: This should take in more configuration options.
func NewBasecoinApp() *BasecoinApp {

	// Create and configure app.
	var app = &BasecoinApp{}
	app.initCapKeys() // ./capkeys.go
	app.initStores()  // ./stores.go
	app.initSDKApp()  // ./sdkapp.go
	app.initRoutes()  // ./routes.go

	// TODO: Load genesis
	// TODO: InitChain with validators
	// TODO: Set the genesis accounts

	app.loadStores()

	return app
}

func (app *BasecoinApp) RunForever() {

	// Start the ABCI server
	srv, err := server.NewServer("0.0.0.0:46658", "socket", app)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	srv.Start()

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})

}

// Load the stores.
func (app *BasecoinApp) loadStores() {
	if err := app.LoadLatestVersion(app.mainStoreKey); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
