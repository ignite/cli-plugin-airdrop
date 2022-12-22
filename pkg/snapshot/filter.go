package snapshot

import (
	"cosmossdk.io/math"

	"github.com/ignite/cli-plugin-airdrop/pkg/formula"
)

type (
	// Filter provide a filter with all airdrop balances
	Filter struct {
		NumberAccounts uint64  `json:"num_accounts" yaml:"num_accounts"`
		Amounts        Amounts `json:"accounts" yaml:"accounts"`
	}

	// Amount provide balance fields of filter Account
	Amount struct {
		Address     string   `json:"address" yaml:"address"`
		DustAmount  math.Int `json:"dust_amount" yaml:"dust_amount"`
		ClaimAmount math.Int `json:"claim_amount" yaml:"claim_amount"`
	}

	// Amounts represents a map of filter Amounts
	Amounts map[string]Amount

	// FilterType represents a Filter type
	FilterType uint64
)

const (
	// FilterStaking filter type staking
	FilterStaking FilterType = iota
	// FilterLiquidity filter type liquidity
	FilterLiquidity
)

// Filter filters a snapshot based on the filter type,
// denom and excluded address, and apply the formula
func (s Snapshot) Filter(
	filterType FilterType,
	denom string,
	dustAmount math.Int,
	formula formula.Value,
	excludedAddresses []string,
) Filter {
	if len(excludedAddresses) > 0 {
		s.Accounts.ExcludeAddresses(excludedAddresses...)
	}
	if len(denom) > 0 {
		s.Accounts.FilterDenom(denom)
	}

	amounts := make(Amounts)
	for address, account := range s.Accounts {
		claimAmount := formula.Calculate(account.BalanceAmount(), account.Staked)
		amounts[address] = Amount{
			Address:     address,
			DustAmount:  dustAmount,
			ClaimAmount: claimAmount,
		}
	}
	return Filter{
		NumberAccounts: uint64(len(amounts)),
		Amounts:        amounts,
	}
}
