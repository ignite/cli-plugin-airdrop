package snapshot

import (
	"encoding/json"

	"cosmossdk.io/math"

	"github.com/ignite/cli-plugin-airdrop/pkg/encode"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type (
	// Snapshot provide a snapshot with all genesis accounts
	Snapshot struct {
		NumberAccounts uint64             `json:"num_accounts"`
		Accounts       map[string]Account `json:"accounts"`
	}

	// Account provide fields of snapshot per account
	// It is the simplified struct we are presenting
	// in this 'balances from state export' snapshot for people.
	Account struct {
		Address        string    `json:"address"`
		Staked         math.Int  `json:"staked"`
		UnbondingStake math.Int  `json:"unbonding_stake"`
		LiquidBalances sdk.Coins `json:"liquid_balance"`
		Bonded         sdk.Coins `json:"bonded"`
	}
)

// newAccount returns a new account.
func newAccount(address string) Account {
	return Account{
		Address:        address,
		Staked:         math.ZeroInt(),
		UnbondingStake: math.ZeroInt(),
		LiquidBalances: sdk.NewCoins(),
		Bonded:         sdk.NewCoins(),
	}
}

// Generate produce the snapshot of address with the total of atom balance liquid,
// staked, bounded and unbonding stake
func Generate(genState map[string]json.RawMessage) (Snapshot, error) {
	marshaller := encode.Codec()
	snapshotAccs := make(map[string]Account)
	var bankGenesis banktypes.GenesisState
	if len(genState[banktypes.ModuleName]) > 0 {
		err := marshaller.UnmarshalJSON(genState[banktypes.ModuleName], &bankGenesis)
		if err != nil {
			return Snapshot{}, err
		}
	}
	for _, balance := range bankGenesis.Balances {
		address := balance.Address
		acc, ok := snapshotAccs[address]
		if !ok {
			acc = newAccount(address)
		}

		acc.LiquidBalances = balance.Coins
		snapshotAccs[address] = acc
	}

	var stakingGenesis stakingtypes.GenesisState
	if len(genState[stakingtypes.ModuleName]) > 0 {
		err := marshaller.UnmarshalJSON(genState[stakingtypes.ModuleName], &stakingGenesis)
		if err != nil {
			return Snapshot{}, err
		}
	}
	for _, unbonding := range stakingGenesis.UnbondingDelegations {
		address := unbonding.DelegatorAddress
		acc, ok := snapshotAccs[address]
		if !ok {
			acc = newAccount(address)
		}

		unbondingStake := sdk.NewInt(0)
		for _, entry := range unbonding.Entries {
			unbondingStake = unbondingStake.Add(entry.Balance)
		}

		acc.UnbondingStake = acc.UnbondingStake.Add(unbondingStake)
		snapshotAccs[address] = acc
	}

	// Make a map from validator operator address to the v036 validator type
	validators := make(map[string]stakingtypes.Validator)
	for _, validator := range stakingGenesis.Validators {
		validators[validator.OperatorAddress] = validator
	}

	for _, delegation := range stakingGenesis.Delegations {
		address := delegation.DelegatorAddress

		acc, ok := snapshotAccs[address]
		if !ok {
			acc = newAccount(address)
		}

		val := validators[delegation.ValidatorAddress]
		staked := delegation.Shares.MulInt(val.Tokens).Quo(val.DelegatorShares).RoundInt()
		acc.Staked = acc.Staked.Add(staked)

		snapshotAccs[address] = acc
	}

	snapshot := Snapshot{
		NumberAccounts: uint64(len(snapshotAccs)),
		Accounts:       snapshotAccs,
	}
	return snapshot, nil
}
