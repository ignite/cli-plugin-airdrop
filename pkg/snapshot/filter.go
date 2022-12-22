package snapshot

import (
	"cosmossdk.io/math"

	"github.com/ignite/cli-plugin-airdrop/pkg/formula"
)

type (
	// Filter provide a filter with all airdrop balances
	Filter struct {
		NumberAmounts uint64  `json:"num_accounts" yaml:"num_accounts"`
		Amounts       Amounts `json:"accounts" yaml:"accounts"`
	}

	// Amount provide balance fields of filter Account
	Amount struct {
		Address     string   `json:"address" yaml:"address"`
		ClaimAmount math.Int `json:"claim_amount" yaml:"claim_amount"`
	}

	// Amounts represents a map of filter Amounts
	Amounts map[string]Amount

	// FilterType represents a Filter type
	FilterType string

	// Filters represents an array of Filter's
	Filters []Filter
)

const (
	// FilterStaking filter type staking
	FilterStaking = "staking"
	// FilterLiquidity filter type liquidity
	FilterLiquidity = "liquidity"
)

// Sum sum all filters into one
func (f Filters) Sum() (result Filter) {
	result.Amounts = make(Amounts)
	for _, filter := range f {
		for _, amount := range filter.Amounts {
			resultAmount := result.Amounts.getAmount(amount.Address)
			resultAmount.ClaimAmount = resultAmount.ClaimAmount.Add(amount.ClaimAmount)
			result.Amounts[amount.Address] = resultAmount
		}
	}
	result.NumberAmounts = uint64(len(result.Amounts))
	return
}

// getAccount get an existing account or generate a new one
func (a Amounts) getAmount(address string) Amount {
	acc, ok := a[address]
	if ok {
		return acc
	}
	return Amount{
		Address:     address,
		ClaimAmount: math.NewInt(0),
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
	if denom != "" {
		s.Accounts.filterDenom(denom)
	}

	amounts := make(Amounts)
	for address, account := range s.Accounts {
		claimAmount := formula.Calculate(account.balanceAmount(), account.Staked)
		amounts[address] = Amount{
			Address:     address,
			ClaimAmount: claimAmount,
		}
	}
	return Filter{
		NumberAmounts: uint64(len(amounts)),
		Amounts:       amounts,
	}
}
