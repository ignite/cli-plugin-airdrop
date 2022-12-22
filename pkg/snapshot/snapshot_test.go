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
}

func TestAccounts_getAccount(t *testing.T) {
	sampleAccounts := Accounts{
		"address_1": {
			Address:        "address_1",
			Staked:         math.NewInt(10),
			UnbondingStake: math.NewInt(10),
			LiquidBalances: sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(10))),
		},
		"address_2": {
			Address:        "address_2",
			Staked:         math.NewInt(12),
			UnbondingStake: math.NewInt(12),
			LiquidBalances: sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(12))),
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

func TestAccounts_ExcludeAddress(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	tests := []struct {
		name    string
		accs    Accounts
		address string
		want    Accounts
	}{
		{
			name: "exclude an address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			address: accAddr1,
			want: Accounts{
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
		},
		{
			name: "exclude an non exist address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			address: sdk.AccAddress(rand.Str(32)).String(),
			want: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
		},
		{
			name: "exclude the last address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
			},
			address: accAddr1,
			want:    Accounts{},
		},
		{
			name:    "exclude an address from a empty list",
			accs:    Accounts{},
			address: accAddr1,
			want:    Accounts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.accs.ExcludeAddress(tt.address)
			require.EqualValues(t, tt.want, tt.accs)
		})
	}
}

func TestAccounts_ExcludeAddresses(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	tests := []struct {
		name      string
		accs      Accounts
		addresses []string
		want      Accounts
	}{
		{
			name: "exclude one address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			addresses: []string{accAddr1},
			want: Accounts{
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
		},
		{
			name: "exclude two address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			addresses: []string{accAddr1, accAddr2},
			want: Accounts{
				accAddr3: {Address: accAddr3},
			},
		},
		{
			name: "exclude all address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			addresses: []string{accAddr1, accAddr2, accAddr3},
			want:      Accounts{},
		},
		{
			name: "exclude all address and a non exiting",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			addresses: []string{accAddr1, accAddr2, accAddr3, sdk.AccAddress(rand.Str(32)).String()},
			want:      Accounts{},
		},
		{
			name: "exclude an non exist address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			addresses: []string{sdk.AccAddress(rand.Str(32)).String()},
			want: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
		},
		{
			name:      "exclude an address from a empty list",
			accs:      Accounts{},
			addresses: []string{sdk.AccAddress(rand.Str(32)).String()},
			want:      Accounts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.accs.ExcludeAddresses(tt.addresses...)
			require.EqualValues(t, tt.want, tt.accs)
		})
	}
}

func TestAccounts_FilterDenom(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	tests := []struct {
		name  string
		accs  Accounts
		denom string
		want  Accounts
	}{
		{
			name: "filter uatom denom",
			accs: Accounts{
				accAddr1: {
					Address: accAddr1,
					LiquidBalances: sdk.NewCoins(
						sdk.NewCoin("uatom", math.NewInt(1000)),
						sdk.NewCoin("token", math.NewInt(2000)),
						sdk.NewCoin("stake", math.NewInt(3000)),
					),
				},
				accAddr2: {
					Address: accAddr2,
					LiquidBalances: sdk.NewCoins(
						sdk.NewCoin("stake", math.NewInt(3000)),
					),
				},
				accAddr3: {
					Address: accAddr3,
					LiquidBalances: sdk.NewCoins(
						sdk.NewCoin("uatom", math.NewInt(1000)),
						sdk.NewCoin("stake", math.NewInt(3000)),
					),
				},
			},
			denom: "uatom",
			want: Accounts{
				accAddr1: {
					Address: accAddr1,
					LiquidBalances: sdk.NewCoins(
						sdk.NewCoin("uatom", math.NewInt(1000)),
					),
				},
				accAddr2: {
					Address:        accAddr2,
					LiquidBalances: sdk.NewCoins(),
				},
				accAddr3: {
					Address: accAddr3,
					LiquidBalances: sdk.NewCoins(
						sdk.NewCoin("uatom", math.NewInt(1000)),
					),
				},
			},
		},
		{
			name: "filter non exiting denom",
			accs: Accounts{
				accAddr1: {
					Address: accAddr1,
					LiquidBalances: sdk.NewCoins(
						sdk.NewCoin("uatom", math.NewInt(1000)),
						sdk.NewCoin("token", math.NewInt(2000)),
						sdk.NewCoin("stake", math.NewInt(3000)),
					),
				},
				accAddr2: {
					Address: accAddr2,
					LiquidBalances: sdk.NewCoins(
						sdk.NewCoin("stake", math.NewInt(3000)),
					),
				},
			},
			denom: "void",
			want: Accounts{
				accAddr1: {
					Address:        accAddr1,
					LiquidBalances: sdk.NewCoins(),
				},
				accAddr2: {
					Address:        accAddr2,
					LiquidBalances: sdk.NewCoins(),
				},
			},
		},
		{
			name: "filter empty balances",
			accs: Accounts{
				accAddr1: {
					Address:        accAddr1,
					LiquidBalances: sdk.NewCoins(),
				},
				accAddr2: {
					Address: accAddr2,
				},
			},
			denom: "uatom",
			want: Accounts{
				accAddr1: {
					Address:        accAddr1,
					LiquidBalances: sdk.NewCoins(),
				},
				accAddr2: {
					Address:        accAddr2,
					LiquidBalances: sdk.NewCoins(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.accs.FilterDenom(tt.denom)
			require.EqualValues(t, tt.want, tt.accs)
		})
	}
}

func TestAccount_TotalStake(t *testing.T) {
	tests := []struct {
		name           string
		staked         math.Int
		unbondingStake math.Int
		want           math.Int
	}{
		{
			name:           "staked and unbounding stake",
			unbondingStake: math.NewInt(100),
			staked:         math.NewInt(999),
			want:           math.NewInt(1099),
		},
		{
			name:           "zero staked",
			unbondingStake: math.NewInt(0),
			staked:         math.NewInt(999),
			want:           math.NewInt(999),
		},
		{
			name:           "zero unbounding stake",
			unbondingStake: math.NewInt(100),
			staked:         math.NewInt(0),
			want:           math.NewInt(100),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Account{
				Staked:         tt.staked,
				UnbondingStake: tt.unbondingStake,
			}
			got := a.TotalStake()
			require.Equal(t, tt.want, got)
		})
	}
}
