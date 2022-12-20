package snapshot

import (
	"cosmossdk.io/math"
	"encoding/json"
	"github.com/ignite/cli-plugin-airdrop/pkg/encode"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type (
	Snapshot struct {
		NumberAccounts uint64             `json:"num_accounts"`
		Accounts       map[string]Account `json:"accounts"`
	}

	// Account provide fields of snapshot per account
	// It is the simplified struct we are presenting in this 'balances from state export' snapshot for people.
	Account struct {
		Address             string               `json:"address"`
		LiquidBalances      sdk.Coins            `json:"liquid_balance"`
		Staked              math.Int             `json:"staked"`
		UnbondingStake      math.Int             `json:"unbonding_stake"`
		Bonded              sdk.Coins            `json:"bonded"`
		BondedBySelectPools map[uint64]sdk.Coins `json:"bonded_by_select_pools"`
		TotalBalances       sdk.Coins            `json:"total_balances"`
	}
)

// newAccount returns a new account.
func newAccount(address string) Account {
	return Account{
		Address:        address,
		LiquidBalances: sdk.Coins{},
		Staked:         math.ZeroInt(),
		UnbondingStake: math.ZeroInt(),
		Bonded:         sdk.Coins{},
	}
}

func Generate(genState map[string]json.RawMessage) (Snapshot, error) {
	// Produce the map of address to total atom balance, both staked and UnbondingStake
	//codec := encode.MakeEncodingConfig()
	marshaler := encode.Codec()
	snapshotAccs := make(map[string]Account)
	var bankGenesis banktypes.GenesisState
	if len(genState[banktypes.ModuleName]) > 0 {
		err := marshaler.UnmarshalJSON(genState[banktypes.ModuleName], &bankGenesis)
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
		err := marshaler.UnmarshalJSON(genState[stakingtypes.ModuleName], &stakingGenesis)
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

		unbondingValue := math.NewInt(0)
		for _, entry := range unbonding.Entries {
			unbondingValue = unbondingValue.Add(entry.Balance)
		}

		acc.UnbondingStake = acc.UnbondingStake.Add(unbondingValue)
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

	// convert balances to underlying coins and sum up balances to total balance
	//for addr, account := range snapshotAccs {
	// All pool shares are in liquid balances OR bonded balances (locked),
	// therefore underlyingCoinsForSelectPools on liquidBalances + bondedBalances
	// will include everything that is in one of those two pools.
	//account.BondedBySelectPools = underlyingCoinsForSelectPools(
	//	account.LiquidBalances.Add(account.Bonded...), pools, selectBondedPoolIDs)
	//account.LiquidBalances = underlyingCoins(account.LiquidBalances, pools)
	//account.Bonded = underlyingCoins(account.Bonded, pools)
	//account.TotalBalances = sdk.NewCoins().
	//	Add(account.LiquidBalances...).
	//	Add(sdk.NewCoin(appparams.BaseCoinUnit, account.Staked)).
	//	Add(sdk.NewCoin(appparams.BaseCoinUnit, account.UnbondingStake)).
	//	Add(account.Bonded...)
	//snapshotAccs[addr] = account
	//}

	snapshot := Snapshot{
		NumberAccounts: uint64(len(snapshotAccs)),
		Accounts:       snapshotAccs,
	}

	return snapshot, nil
}
