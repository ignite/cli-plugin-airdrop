package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
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
		Types    string    `yaml:"types"`
		Chain    string    `yaml:"chain"`
		Denom    string    `yaml:"denom"`
		Date     time.Time `yaml:"date"`
		Formula  Formula   `yaml:"formula"`
		Excluded []string  `yaml:"excluded"`
	}

	// Formula defines a struct with the fields that are common to all config snapshot formula.
	Formula struct {
		Type  string `yaml:"type"`
		Value uint64 `yaml:"value"`
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
