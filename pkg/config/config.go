package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/ignite/cli-plugin-airdrop/pkg/formula"
)

type (
	// Config defines a struct with the fields that are common to all config.
	Config struct {
		AirdropToken string     `json:"airdrop_token" yaml:"airdrop_token"`
		DustWallet   uint64     `json:"dust_wallet" yaml:"dust_wallet"`
		Snapshots    []Snapshot `json:"snapshots" yaml:"snapshots"`
	}

	// Snapshot defines a struct with the fields that are common to all config snapshot.
	Snapshot struct {
		Type     string        `json:"type" yaml:"type"`
		Denom    string        `json:"denom" yaml:"denom"`
		Formula  formula.Value `json:"formula" yaml:"formula"`
		Excluded []string      `json:"excluded" yaml:"excluded"`
	}
)

// ParseConfig expects to find and parse a config file.
func ParseConfig(filename string) (c Config, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return c, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&c); err != nil {
		return c, err
	}
	return c, nil
}
