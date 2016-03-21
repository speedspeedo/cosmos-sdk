package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tendermint/basecoin/tests"
	"github.com/tendermint/basecoin/types"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-rpc/client"
	"github.com/tendermint/go-rpc/types"
	"github.com/tendermint/go-wire"
	_ "github.com/tendermint/tendermint/rpc/core/types" // Register RPCResponse > Result types
)

func main() {
	ws := rpcclient.NewWSClient("127.0.0.1:46657", "/websocket")
	_, err := ws.Start()
	if err != nil {
		Exit(err.Error())
	}
	var counter = 0

	// Read a bunch of responses
	go func() {
		for {
			res, ok := <-ws.ResultsCh
			if !ok {
				break
			}
			fmt.Println(counter, "res:", Blue(string(res)))
		}
	}()

	// Get the root account
	root := tests.PrivAccountFromSecret("root")
	sequence := int(0)
	// Make a bunch of PrivAccounts
	privAccounts := tests.RandAccounts(1000, 1000000, 0)
	privAccountSequences := make(map[string]int)

	// Send coins to each account
	for i := 0; i < len(privAccounts); i++ {
		privAccount := privAccounts[i]
		tx := &types.SendTx{
			Inputs: []types.TxInput{
				types.TxInput{
					Address:  root.Account.PubKey.Address(),
					PubKey:   root.Account.PubKey, // TODO is this needed?
					Amount:   1000002,
					Sequence: sequence,
				},
			},
			Outputs: []types.TxOutput{
				types.TxOutput{
					Address: privAccount.Account.PubKey.Address(),
					Amount:  1000000,
				},
			},
		}
		sequence += 1

		// Sign request
		signBytes := wire.BinaryBytes(tx)
		sig := root.PrivKey.Sign(signBytes)
		tx.Inputs[0].Signature = sig
		//fmt.Println("tx:", tx)

		// Write request
		txBytes := wire.BinaryBytes(tx)
		request := rpctypes.NewRPCRequest("fakeid", "broadcast_tx_sync", Arr(txBytes))
		reqBytes := wire.JSONBytes(request)
		//fmt.Print(".")
		err := ws.WriteMessage(websocket.TextMessage, reqBytes)
		if err != nil {
			Exit("writing websocket request: " + err.Error())
		}
	}

	// Now send coins between these accounts
	for {
		counter += 1
		time.Sleep(time.Millisecond * 10)

		randA := RandInt() % len(privAccounts)
		randB := RandInt() % len(privAccounts)
		if randA == randB {
			continue
		}

		privAccountA := privAccounts[randA]
		privAccountASequence := privAccountSequences[privAccountA.Account.PubKey.KeyString()]
		privAccountSequences[privAccountA.Account.PubKey.KeyString()] = privAccountASequence + 1
		privAccountB := privAccounts[randB]

		tx := &types.SendTx{
			Inputs: []types.TxInput{
				types.TxInput{
					Address:  privAccountA.Account.PubKey.Address(),
					PubKey:   privAccountA.Account.PubKey,
					Amount:   3,
					Sequence: privAccountASequence,
				},
			},
			Outputs: []types.TxOutput{
				types.TxOutput{
					Address: privAccountB.Account.PubKey.Address(),
					Amount:  1,
				},
			},
		}

		// Sign request
		signBytes := wire.BinaryBytes(tx)
		sig := privAccountA.PrivKey.Sign(signBytes)
		tx.Inputs[0].Signature = sig
		//fmt.Println("tx:", tx)

		// Write request
		txBytes := wire.BinaryBytes(tx)
		request := rpctypes.NewRPCRequest("fakeid", "broadcast_tx_sync", Arr(txBytes))
		reqBytes := wire.JSONBytes(request)
		//fmt.Print(".")
		err := ws.WriteMessage(websocket.TextMessage, reqBytes)
		if err != nil {
			Exit("writing websocket request: " + err.Error())
		}
	}

	ws.Stop()
}
