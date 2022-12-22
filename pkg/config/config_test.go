package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/ignite/cli-plugin-airdrop/pkg/formula"
)

func TestParseConfig(t *testing.T) {
	sampleConfig := Config{
		AirdropToken: "ufoo",
		DustWallet:   1,
		Snapshots: []Snapshot{
			{
				Type:  "staking",
				Denom: "uatom",
				Formula: formula.Value{
					Type:  "quadratic",
					Value: math.NewInt(2),
				},
				Excluded: []string{"cosmos1aqn8ynvr3jmq67879qulzrwhchq5dtrvh6h4er"},
			},
			{
				Type:  "liquidity",
				Denom: "uatom",
				Formula: formula.Value{
					Type:  "quadratic",
					Value: math.NewInt(10),
				},
				Excluded: []string{"cosmos1aqn8ynvr3jmq67879qulzrwhchq5dtrvh6h4er"},
			},
		},
	}
	yamlData, err := yaml.Marshal(&sampleConfig)
	require.NoError(t, err)
	sampleConfigPath := filepath.Join(t.TempDir(), "config.yml")
	err = os.WriteFile(sampleConfigPath, yamlData, 0o644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filename string
		want     Config
		err      error
	}{
		{
			name:     "valid config file",
			filename: sampleConfigPath,
			want:     sampleConfig,
		},
		{
			name:     "valid config file",
			filename: "invalid_file_path",
			want:     sampleConfig,
			err:      errors.New("open invalid_file_path: no such file or directory"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfig(tt.filename)
			if tt.err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want.AirdropToken, got.AirdropToken)
			require.Equal(t, tt.want.DustWallet, got.DustWallet)
			for i, wantSnapshot := range tt.want.Snapshots {
				require.Equal(t, wantSnapshot.Denom, got.Snapshots[i].Denom)
				require.Equal(t, wantSnapshot.Formula, got.Snapshots[i].Formula)
				require.Equal(t, wantSnapshot.Type, got.Snapshots[i].Type)
				require.EqualValues(t, wantSnapshot.Excluded, got.Snapshots[i].Excluded)

			}
		})
	}
}
