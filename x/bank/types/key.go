package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

const (
	// ModuleName defines the module name
	ModuleName = "bank"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore keys
var (
	SupplyKey           = []byte{0x00}
	DenomMetadataPrefix = []byte{0x1}
	DenomAddressPrefix  = []byte{0x03}

	// BalancesPrefix is the prefix for the account balances store. We use a byte
	// (instead of `[]byte("balances")` to save some disk space).
	BalancesPrefix = []byte{0x02}
)

// DenomMetadataKey returns the denomination metadata key.
func DenomMetadataKey(denom string) []byte {
	d := []byte(denom)
	return append(DenomMetadataPrefix, d...)
}

// AddressFromBalancesStore returns an account address from a balances prefix
// store. The key must not contain the prefix BalancesPrefix as the prefix store
// iterator discards the actual prefix.
//
// If invalid key is passed, AddressFromBalancesStore returns ErrInvalidKey.
func AddressFromBalancesStore(key []byte) (sdk.AccAddress, error) {
	if len(key) == 0 {
		return nil, ErrInvalidKey
	}

	kv.AssertKeyAtLeastLength(key, 1)

	addrLen := key[0]
	bound := int(addrLen)

	if len(key)-1 < bound {
		return nil, ErrInvalidKey
	}

	return key[1 : bound+1], nil
}

// CreateAccountBalancesPrefix creates the prefix for an account's balances.
func CreateAccountBalancesPrefix(addr []byte) []byte {
	return append(BalancesPrefix, address.MustLengthPrefix(addr)...)
}

// CreateDenomAddressPrefix creates a prefix for a reverse index of denomination
// to account balance for that denomination.
func CreateDenomAddressPrefix(denom string) []byte {
	key := append(DenomAddressPrefix, []byte(denom)...)
	return append(key, 0)
}
