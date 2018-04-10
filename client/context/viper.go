package context

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"

	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/core"
)

func NewCoreContextFromViper() core.CoreContext {
	nodeURI := viper.GetString(client.FlagNode)
	var rpc rpcclient.Client
	if nodeURI != "" {
		rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
	}
	chainID := viper.GetString(client.FlagChainID)
	// if chain ID is not specified manually, read chain ID from genesis file if present
	if chainID == "" {
		cfg, err := tcmd.ParseConfig()
		if err == nil {
			genesisFile := cfg.GenesisFile()
			bz, err := ioutil.ReadFile(genesisFile)
			if err == nil {
				var doc tmtypes.GenesisDoc
				err = json.Unmarshal(bz, &doc)
				if err == nil {
					chainID = doc.ChainID
				}
			}
		}
	}
	return core.CoreContext{
		ChainID:         chainID,
		Height:          viper.GetInt64(client.FlagHeight),
		TrustNode:       viper.GetBool(client.FlagTrustNode),
		FromAddressName: viper.GetString(client.FlagName),
		NodeURI:         nodeURI,
		Sequence:        viper.GetInt64(client.FlagSequence),
		Client:          rpc,
	}
}
