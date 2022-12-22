package formula

import (
	"math/big"

	"cosmossdk.io/math"
)

const (
	// Quadratic represents a quadratic airdrop type
	Quadratic Type = "quadratic"
)

type (
	// Value defines a struct for the formula type
	Value struct {
		Type   Type     `json:"type" yaml:"type"`
		Value  math.Int `json:"value" yaml:"value"`
		Ignore math.Int `json:"ignore" yaml:"ignore"`
	}
	// Type defines a formula type
	Type string
)

// Calculate calculates the airdrop amount base on the formula type
// and parameters, total amount, staked amount and the balance
func (v Value) Calculate(amount, staked math.Int) math.Int {
	switch v.Type {
	case Quadratic:
		stakedPercent := staked.Quo(amount)
		base := math.NewIntFromBigInt(big.NewInt(0).Sqrt(amount.BigInt()))
		bonus := base.Mul(v.Value).Mul(stakedPercent)
		airdrop := base.Add(bonus)
		if airdrop.LTE(v.Ignore) {
			return math.ZeroInt()
		}
		return airdrop
	}
	return math.ZeroInt()
}
