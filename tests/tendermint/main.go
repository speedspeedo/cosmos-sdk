package main

func main() {}

// import (
// 	"fmt"
// 	"time"

// 	"github.com/gorilla/websocket"
// 	"github.com/cosmos/cosmos-sdk/types"
// 	wire "github.com/tendermint/go-wire"
// 	_ "github.com/tendermint/tendermint/rpc/core/types" // Register RPCResponse > Result types
// 	"github.com/tendermint/tendermint/rpc/lib/client"
// 	"github.com/tendermint/tendermint/rpc/lib/types"
// 	cmn "github.com/tendermint/tmlibs/common"
// )

// func main() {
// 	// ws := rpcclient.NewWSClient("127.0.0.1:46657", "/websocket")
// 	ws := rpcclient.NewWSClient("192.168.99.100:46657", "/websocket")
// 	chainID := "test_chain_id"

// 	_, err := ws.Start()
// 	if err != nil {
// 		cmn.Exit(err.Error())
// 	}
// 	var counter = 0

// 	// Read a bunch of responses
// 	go func() {
// 		for {
// 			res, ok := <-ws.ResultsCh
// 			if !ok {
// 				break
// 			}
// 			fmt.Println(counter, "res:", cmn.Blue(string(res)))
// 		}
// 	}()

// 	// Get the root account
// 	root := types.PrivAccountFromSecret("test")
// 	sequence := int(0)
// 	// Make a bunch of PrivAccounts
// 	privAccounts := types.RandAccounts(1000, 1000000, 0)
// 	privAccountSequences := make(map[string]int)

// 	// Send coins to each account
// 	for i := 0; i < len(privAccounts); i++ {
// 		privAccount := privAccounts[i]
// 		tx := &types.SendTx{
// 			Inputs: []types.TxInput{
// 				types.TxInput{
// 					Address:  root.Account.PubKey.Address(),
// 					PubKey:   root.Account.PubKey, // TODO is this needed?
// 					Coins:    coin.Coins{{"", 1000002}},
// 					Sequence: sequence,
// 				},
// 			},
// 			Outputs: []types.TxOutput{
// 				types.TxOutput{
// 					Address: privAccount.Account.PubKey.Address(),
// 					Coins:   coin.Coins{{"", 1000000}},
// 				},
// 			},
// 		}
// 		sequence += 1

// 		// Sign request
// 		signBytes := tx.SignBytes(chainID)
// 		sig := root.Sign(signBytes)
// 		tx.Inputs[0].Signature = sig
// 		//fmt.Println("tx:", tx)

// 		// Write request
// 		txBytes := wire.BinaryBytes(struct{ types.Tx }{tx})
// 		request, err := rpctypes.MapToRequest("fakeid", "broadcast_tx_sync", map[string]interface{}{"tx": txBytes})
// 		if err != nil {
// 			cmn.Exit("cannot encode request: " + err.Error())
// 		}
// 		reqBytes := wire.JSONBytes(request)
// 		//fmt.Print(".")
// 		err = ws.WriteMessage(websocket.TextMessage, reqBytes)
// 		if err != nil {
// 			cmn.Exit("writing websocket request: " + err.Error())
// 		}
// 	}

// 	// Now send coins between these accounts
// 	for {
// 		counter += 1
// 		time.Sleep(time.Millisecond * 10)

// 		randA := cmn.RandInt() % len(privAccounts)
// 		randB := cmn.RandInt() % len(privAccounts)
// 		if randA == randB {
// 			continue
// 		}

// 		privAccountA := privAccounts[randA]
// 		privAccountASequence := privAccountSequences[privAccountA.Account.PubKey.KeyString()]
// 		privAccountSequences[privAccountA.Account.PubKey.KeyString()] = privAccountASequence + 1
// 		privAccountB := privAccounts[randB]

// 		tx := &types.SendTx{
// 			Inputs: []types.TxInput{
// 				types.TxInput{
// 					Address:  privAccountA.Account.PubKey.Address(),
// 					PubKey:   privAccountA.Account.PubKey,
// 					Coins:    coin.Coins{{"", 3}},
// 					Sequence: privAccountASequence + 1,
// 				},
// 			},
// 			Outputs: []types.TxOutput{
// 				types.TxOutput{
// 					Address: privAccountB.Account.PubKey.Address(),
// 					Coins:   coin.Coins{{"", 1}},
// 				},
// 			},
// 		}

// 		// Sign request
// 		signBytes := tx.SignBytes(chainID)
// 		sig := privAccountA.Sign(signBytes)
// 		tx.Inputs[0].Signature = sig
// 		//fmt.Println("tx:", tx)

// 		// Write request
// 		txBytes := wire.BinaryBytes(struct{ types.Tx }{tx})
// 		request, err := rpctypes.MapToRequest("fakeid", "broadcast_tx_sync", map[string]interface{}{"tx": txBytes})
// 		if err != nil {
// 			cmn.Exit("cannot encode request: " + err.Error())
// 		}
// 		reqBytes := wire.JSONBytes(request)
// 		//fmt.Print(".")
// 		err = ws.WriteMessage(websocket.TextMessage, reqBytes)
// 		if err != nil {
// 			cmn.Exit("writing websocket request: " + err.Error())
// 		}
// 	}

// 	ws.Stop()
// }
