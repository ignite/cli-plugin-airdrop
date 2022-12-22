package snapshot

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// Snapshot provide a snapshot with all genesis Accounts
	Snapshot struct {
		NumberAccounts uint64   `json:"num_accounts" yaml:"num_accounts"`
		Accounts       Accounts `json:"accounts" yaml:"accounts"`
	}

	// Account provide fields of snapshot per account
	// It is the simplified struct we are presenting
	// in this 'balances from state export' snapshot for people.
	Account struct {
		Address        string    `json:"address" yaml:"address"`
		Staked         math.Int  `json:"staked" yaml:"staked"`
		UnbondingStake math.Int  `json:"unbonding_stake" yaml:"unbonding_stake"`
		LiquidBalances sdk.Coins `json:"liquid_balance" yaml:"liquid_balance"`
	}

	// Accounts represents a map of snapshot Accounts
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

// TotalStake returns a sum of stake and unbounding stake
func (a Account) TotalStake() math.Int {
	if a.Staked.IsNil() {
		return a.UnbondingStake
	}
	if a.UnbondingStake.IsNil() {
		return a.Staked
	}
	return a.Staked.Add(a.UnbondingStake)
}

// BalanceAmount returns a sum of all denom balances
func (a Account) BalanceAmount() math.Int {
	amount := math.ZeroInt()
	for _, coin := range a.LiquidBalances {
		amount = amount.Add(coin.Amount)
	}
	return amount
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
