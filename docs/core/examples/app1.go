package app

import (
	"encoding/json"
	"reflect"

	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	bapp "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	app1Name = "App1"
)

func NewApp1(logger log.Logger, db dbm.DB) *bapp.BaseApp {

	// TODO: make this an interface or pass in
	// a TxDecoder instead.
	cdc := wire.NewCodec()

	// Create the base application object.
	app := bapp.NewBaseApp(app1Name, cdc, logger, db)

	// Create a key for accessing the account store.
	keyAccount := sdk.NewKVStoreKey("acc")

	// Determine how transactions are decoded.
	app.SetTxDecoder(txDecoder)

	// Register message routes.
	// Note the handler gets access to the account store.
	app.Router().
		AddRoute("bank", NewApp1Handler(keyAccount))

	// Mount stores and load the latest state.
	app.MountStoresIAVL(keyAccount)
	err := app.LoadLatestVersion(keyAccount)
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}

//------------------------------------------------------------------
// Msg

// MsgSend implements sdk.Msg
var _ sdk.Msg = MsgSend{}

// MsgSend to send coins from Input to Output
type MsgSend struct {
	From   sdk.Address `json:"from"`
	To     sdk.Address `json:"to"`
	Amount sdk.Coins   `json:"amount"`
}

// NewMsgSend
func NewMsgSend(from, to sdk.Address, amt sdk.Coins) MsgSend {
	return MsgSend{from, to, amt}
}

// Implements Msg.
func (msg MsgSend) Type() string { return "bank" }

// Implements Msg. Ensure the addresses are good and the
// amount is positive.
func (msg MsgSend) ValidateBasic() sdk.Error {
	if len(msg.From) == 0 {
		return sdk.ErrInvalidAddress("From address is empty")
	}
	if len(msg.To) == 0 {
		return sdk.ErrInvalidAddress("To address is empty")
	}
	if !msg.Amount.IsPositive() {
		return sdk.ErrInvalidCoins("Amount is not positive")
	}
	return nil
}

// Implements Msg. JSON encode the message.
func (msg MsgSend) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// Implements Msg. Return the signer.
func (msg MsgSend) GetSigners() []sdk.Address {
	return []sdk.Address{msg.From}
}

// Returns the sdk.Tags for the message
func (msg MsgSend) Tags() sdk.Tags {
	return sdk.NewTags("sender", []byte(msg.From.String())).
		AppendTag("receiver", []byte(msg.To.String()))
}

//------------------------------------------------------------------
// Handler for the message

func NewApp1Handler(keyAcc *sdk.KVStoreKey) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, keyAcc, msg)
		default:
			errMsg := "Unrecognized bank Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
// NOTE: msg.From, msg.To, and msg.Amount were already validated
func handleMsgSend(ctx sdk.Context, key *sdk.KVStoreKey, msg MsgSend) sdk.Result {
	// Load the store.
	store := ctx.KVStore(key)

	// Debit from the sender.
	if res := handleFrom(store, msg.From, msg.Amount); !res.IsOK() {
		return res
	}

	// Credit the receiver.
	if res := handleTo(store, msg.To, msg.Amount); !res.IsOK() {
		return res
	}

	// Return a success (Code 0).
	// Add list of key-value pair descriptors ("tags").
	return sdk.Result{
		Tags: msg.Tags(),
	}
}

func handleFrom(store sdk.KVStore, from sdk.Address, amt sdk.Coins) sdk.Result {
	// Get sender account from the store.
	accBytes := store.Get(from)
	if accBytes == nil {
		// Account was not added to store. Return the result of the error.
		return sdk.NewError(2, 101, "Account not added to store").Result()
	}

	// Unmarshal the JSON account bytes.
	var acc app1Account
	err := json.Unmarshal(accBytes, &acc)
	if err != nil {
		// InternalError
		return sdk.ErrInternal("Error when deserializing account").Result()
	}

	// Deduct msg amount from sender account.
	senderCoins := acc.Coins.Minus(amt)

	// If any coin has negative amount, return insufficient coins error.
	if !senderCoins.IsNotNegative() {
		return sdk.ErrInsufficientCoins("Insufficient coins in account").Result()
	}

	// Set acc coins to new amount.
	acc.Coins = senderCoins

	// Encode sender account.
	accBytes, err = json.Marshal(acc)
	if err != nil {
		return sdk.ErrInternal("Account encoding error").Result()
	}

	// Update store with updated sender account
	store.Set(from, accBytes)
	return sdk.Result{}
}

func handleTo(store sdk.KVStore, to sdk.Address, amt sdk.Coins) sdk.Result {
	// Add msg amount to receiver account
	accBytes := store.Get(to)
	var acc app1Account
	if accBytes == nil {
		// Receiver account does not already exist, create a new one.
		acc = app1Account{}
	} else {
		// Receiver account already exists. Retrieve and decode it.
		err := json.Unmarshal(accBytes, &acc)
		if err != nil {
			return sdk.ErrInternal("Account decoding error").Result()
		}
	}

	// Add amount to receiver's old coins
	receiverCoins := acc.Coins.Plus(amt)

	// Update receiver account
	acc.Coins = receiverCoins

	// Encode receiver account
	accBytes, err := json.Marshal(acc)
	if err != nil {
		return sdk.ErrInternal("Account encoding error").Result()
	}

	// Update store with updated receiver account
	store.Set(to, accBytes)
	return sdk.Result{}
}

type app1Account struct {
	Coins sdk.Coins `json:"coins"`
}

//------------------------------------------------------------------
// Tx

// Simple tx to wrap the Msg.
type app1Tx struct {
	MsgSend
}

// This tx only has one Msg.
func (tx app1Tx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{tx.MsgSend}
}

// JSON decode MsgSend.
func txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx app1Tx
	err := json.Unmarshal(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode(err.Error())
	}
	return tx, nil
}
