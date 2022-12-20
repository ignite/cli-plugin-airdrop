package snapshot

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"
)

func TestNewAccount(t *testing.T) {
	var (
		address = sdk.AccAddress(rand.Str(10)).String()
		got     = newAccount(address)
	)
	require.Equal(t, address, got.Address)
	require.Equal(t, math.ZeroInt(), got.Staked)
	require.Equal(t, math.ZeroInt(), got.UnbondingStake)
	require.Equal(t, sdk.NewCoins(), got.LiquidBalances)
	require.Equal(t, sdk.NewCoins(), got.Bonded)
}

func TestAccounts_getAccount(t *testing.T) {
	sampleAccounts := Accounts{
		"address_1": {
			Address:        "address_1",
			Staked:         math.NewInt(10),
			UnbondingStake: math.NewInt(10),
			LiquidBalances: sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(10))),
			Bonded:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(10))),
		},
		"address_2": {
			Address:        "address_2",
			Staked:         math.NewInt(12),
			UnbondingStake: math.NewInt(12),
			LiquidBalances: sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(12))),
			Bonded:         sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(12))),
		},
	}
	tests := []struct {
		name    string
		a       Accounts
		address string
		want    Account
	}{
		{
			name:    "already exist address 1",
			a:       sampleAccounts,
			address: "address_1",
			want:    sampleAccounts["address_1"],
		},
		{
			name:    "already exist address 2",
			a:       sampleAccounts,
			address: "address_2",
			want:    sampleAccounts["address_2"],
		},
		{
			name:    "not exist address",
			a:       sampleAccounts,
			address: "address_3",
			want: Account{
				Address:        "address_3",
				Staked:         math.ZeroInt(),
				UnbondingStake: math.ZeroInt(),
				LiquidBalances: sdk.NewCoins(),
				Bonded:         sdk.NewCoins(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.getAccount(tt.address)
			require.Equal(t, tt.want, got)
		})
	}
}
