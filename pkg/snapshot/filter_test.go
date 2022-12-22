package snapshot

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"

	"github.com/ignite/cli-plugin-airdrop/pkg/formula"
)

func TestAccounts_getAmount(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	sampleAamounts := Amounts{
		accAddr1: {
			Address:     accAddr1,
			ClaimAmount: math.NewInt(10),
		},
		accAddr2: {
			Address:     accAddr2,
			ClaimAmount: math.NewInt(1000),
		},
	}
	tests := []struct {
		name    string
		a       Amounts
		address string
		want    Amount
	}{
		{
			name:    "already exist address 1",
			a:       sampleAamounts,
			address: accAddr1,
			want:    sampleAamounts[accAddr1],
		},
		{
			name:    "already exist address 2",
			a:       sampleAamounts,
			address: accAddr2,
			want:    sampleAamounts[accAddr2],
		},
		{
			name:    "not exist address",
			a:       sampleAamounts,
			address: accAddr3,
			want: Amount{
				Address:     accAddr3,
				ClaimAmount: math.ZeroInt(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.getAmount(tt.address)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSnapshot_Filter(t *testing.T) {
	type args struct {
		filterType        FilterType
		denom             string
		formula           formula.Value
		excludedAddresses []string
	}
	tests := []struct {
		name     string
		snapshot Snapshot
		args     args
		want     Filter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.snapshot.Filter(tt.args.filterType, tt.args.denom, tt.args.formula, tt.args.excludedAddresses)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFilters_Sum(t *testing.T) {
	tests := []struct {
		name    string
		filters Filters
		want    Filter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filters.Sum()
			require.Equal(t, tt.want, got)
		})
	}
}
