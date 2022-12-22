package formula

import "cosmossdk.io/math"

const (
	Quadratic Type = "quadratic"
)

type (
	// Value defines a struct for the formula type
	Value struct {
		Type  Type     `yaml:"type"`
		Value math.Int `yaml:"value"`
	}
	// Type defines a formula type
	Type string
)

func (v Value) Calculate(value math.Int) math.Int {
	switch v.Type {
	case Quadratic:
		return value.Mul(v.Value)
	}
	return math.ZeroInt()
}
