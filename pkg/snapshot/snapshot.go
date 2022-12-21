package snapshot

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// Snapshot provide a snapshot with all genesis accounts
	Snapshot struct {
		NumberAccounts uint64   `json:"num_accounts"`
		Accounts       Accounts `json:"accounts"`
	}

	// Account provide fields of snapshot per account
	// It is the simplified struct we are presenting
	// in this 'balances from state export' snapshot for people.
	Account struct {
		Address        string    `json:"address"`
		Staked         math.Int  `json:"staked"`
		UnbondingStake math.Int  `json:"unbonding_stake"`
		LiquidBalances sdk.Coins `json:"liquid_balance"`
	}

	// Accounts represents a map of snapshot accounts
	Accounts map[string]Account
)

// newAccount returns a new account.
func newAccount(address string) Account {
	return Account{
		Address:        address,
		Staked:         math.ZeroInt(),
		UnbondingStake: math.ZeroInt(),
		LiquidBalances: sdk.NewCoins(),
	}
}

// getAccount get an existing account or generate a new one
func (a Accounts) getAccount(address string) Account {
	acc, ok := a[address]
	if ok {
		return acc
	}
	return newAccount(address)
}

// ExcludeAddress exclude an address from the accounts
func (a Accounts) ExcludeAddress(address string) {
	for accAddress := range a {
		if accAddress == address {
			delete(a, accAddress)
		}
	}
}

// ExcludeAddresses exclude an address list from the accounts
func (a Accounts) ExcludeAddresses(addresses ...string) {
	for _, address := range addresses {
		a.ExcludeAddress(address)
	}
}

// FilterDenom filter balance by denom
func (a Accounts) FilterDenom(denom string) {
	for address, account := range a {
		found, liquidBalance := account.LiquidBalances.Find(denom)
		if found {
			account.LiquidBalances = sdk.NewCoins(liquidBalance)
		} else {
			account.LiquidBalances = sdk.NewCoins()
		}
		a[address] = account
	}
}
