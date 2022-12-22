package snapshot

import (
	"cosmossdk.io/math"
	claimtypes "github.com/ignite/modules/x/claim/types"

	"github.com/ignite/cli-plugin-airdrop/pkg/formula"
)

type (
	// Filter provide a filter with all airdrop balances
	Filter map[string]claimtypes.ClaimRecord

	// FilterType represents a Filter type
	FilterType string

	// Filters represents an array of Filter's
	Filters []Filter
)

const (
	// Staking filter type staking
	Staking = "staking"
	// Liquidity filter type liquidity
	Liquidity = "liquidity"
)

// ClaimRecords return a list of claim records
func (f Filter) ClaimRecords() []claimtypes.ClaimRecord {
	result := make([]claimtypes.ClaimRecord, 0)
	for _, filter := range f {
		result = append(result, filter)
	}
	return result
}

// Sum sum all filters into one
func (f Filters) Sum() Filter {
	result := make(Filter)
	for _, filter := range f {
		for _, amount := range filter {
			resultAmount := result.getAmount(amount.Address)
			resultAmount.Claimable = resultAmount.Claimable.Add(amount.Claimable)
			result[amount.Address] = resultAmount
		}
	}
	return result
}

// getAccount get an existing account or generate a new one
func (f Filter) getAmount(address string) claimtypes.ClaimRecord {
	acc, ok := f[address]
	if ok {
		return acc
	}
	return claimtypes.ClaimRecord{
		Address:   address,
		Claimable: math.NewInt(0),
	}
}

// Filter filters a snapshot based on the filter type,
// denom and excluded address, and apply the formula
func (s Snapshot) Filter(
	filterType FilterType,
	denom string,
	formula formula.Value,
	excludedAddresses []string,
) Filter {
	if len(excludedAddresses) > 0 {
		s.Accounts.excludeAddresses(excludedAddresses...)
	}
	s.Accounts.filterDenom(denom)

	filter := make(Filter)
	for address, account := range s.Accounts {
		// TODO FIXME for the liquidity model
		amount := account.balanceAmount()
		if filterType == Staking {
			amount = account.balanceAmount()
		}
		claimAmount := formula.Calculate(amount, account.Staked)
		filter[address] = claimtypes.ClaimRecord{
			Address:   address,
			Claimable: claimAmount,
		}
	}
	return filter
}
