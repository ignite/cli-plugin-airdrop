package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/ignite/cli-plugin-airdrop/pkg/formula"
)

type (
	// Config defines a struct with the fields that are common to all config.
	Config struct {
		AirdropToken string     `yaml:"airdrop_token"`
		DustWallet   uint64     `yaml:"dust_wallet"`
		Snapshots    []Snapshot `yaml:"snapshots,omitempty"`
	}

	// Snapshot defines a struct with the fields that are common to all config snapshot.
	Snapshot struct {
		Type     string        `yaml:"type"`
		Denom    string        `yaml:"denom"`
		Formula  formula.Value `yaml:"formula"`
		Excluded []string      `yaml:"excluded"`
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
