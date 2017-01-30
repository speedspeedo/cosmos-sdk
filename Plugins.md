# Basecoin Plugins

Basecoin is an extensible cryptocurrency module.
Each Basecoin account contains a ED25519 public key,
a balance in many different coin denominations,
and a strictly increasing sequence number for replay protection (like in Ethereum).
Accounts are serialized and stored in a merkle tree using the account's address as the key,
where the address is the RIPEMD160 hash of the public key.

Sending tokens around is done via the `SendTx`, which takes a list of inputs and a list of outputs,
and transfers all the tokens listed in the inputs from their corresponding accounts to the accounts listed in the output.
The `SendTx` is structured as follows:

```
type SendTx struct {
  Gas     int64      `json:"gas"` // Gas
  Fee     Coin       `json:"fee"` // Fee
  Inputs  []TxInput  `json:"inputs"`
  Outputs []TxOutput `json:"outputs"`
}

type TxInput struct {
  Address   []byte           `json:"address"`   // Hash of the PubKey
  Coins     Coins            `json:"coins"`     //
  Sequence  int              `json:"sequence"`  // Must be 1 greater than the last committed TxInput
  Signature crypto.Signature `json:"signature"` // Depends on the PubKey type and the whole Tx
  PubKey    crypto.PubKey    `json:"pub_key"`   // Is present iff Sequence == 0
}

type TxOutput struct {
  Address []byte `json:"address"` // Hash of the PubKey
  Coins   Coins  `json:"coins"`   //
}

type Coins []Coin

type Coin struct {
  Denom  string `json:"denom"`
  Amount int64  `json:"amount"`
}

```

Note it also includes a field for `Gas` and `Fee`. The `Gas` limits the total amount of computation that can be done by the transaction,
while the `Fee` refers to the total amount paid in fees. This is slightly different from Ethereum's concept of `Gas` and `GasPrice`,
where `Fee = Gas x GasPrice`. In Basecoin, the `Gas` and `Fee` are independent.


Basecoin also defines another transaction type, the `AppTx`:

```
type AppTx struct {
  Gas   int64   `json:"gas"`   // Gas
  Fee   Coin    `json:"fee"`   // Fee
  Name  string  `json:"type"`  // Which plugin
  Input TxInput `json:"input"`
  Data  []byte  `json:"data"`
}
```

The `AppTx` enables arbitrary additional functionality through the use of plugins.
A plugin is simply a Go package that implements the `Plugin` interface:

```
type Plugin interface {

  // Name of this plugin, should be short.
  Name() string

  // Run a transaction from ABCI DeliverTx
  RunTx(store KVStore, ctx CallContext, txBytes []byte) (res abci.Result)

  // Other ABCI message handlers
  SetOption(store KVStore, key string, value string) (log string)
  InitChain(store KVStore, vals []*abci.Validator)
  BeginBlock(store KVStore, height uint64)
  EndBlock(store KVStore, height uint64) []*abci.Validator
}

type CallContext struct {
  CallerAddress []byte   // Caller's Address (hash of PubKey)
  CallerAccount *Account // Caller's Account, w/ fee & TxInputs deducted
  Coins         Coins    // The coins that the caller wishes to spend, excluding fees
}
```

The workhorse of the plugin is `RunTx`, which is called when an `AppTx` is processed.
The `Name` field in the `AppTx` refers to the plugin name, and the `Data` field of the `AppTx` is
forward to the `RunTx` function.

You can look at some example plugins in the [basecoin repo](https://github.com/tendermint/basecoin/tree/develop/plugins).

If you want to see how you can write a plugin in your own repo, and make use of all the basecoin tooling, cli, etc. please take a look at the [mintcoin example](https://github.com/tendermint/basecoin-examples/tree/master/mintcoin) for inspiration, not just the plugin itself, but also the `cmd/mintcoin` directory to create the custom command.
