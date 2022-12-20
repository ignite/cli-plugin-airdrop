package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestParseConfig(t *testing.T) {
	sampleConfig := Config{
		AirdropToken: "ufoo",
		DustWallet:   1,
		Snapshots: []Snapshot{
			{
				Types: "staking",
				Chain: "cosmos-hub",
				Denom: "uatom",
				Date:  time.Now(),
				Formula: Formula{
					Type:  "quadratic",
					Value: 2,
				},
				Excluded: []string{"cosmos1aqn8ynvr3jmq67879qulzrwhchq5dtrvh6h4er"},
			},
			{
				Types: "liquidity",
				Chain: "cosmos-hub",
				Denom: "uatom",
				Date:  time.Now(),
				Formula: Formula{
					Type:  "quadratic",
					Value: 10,
				},
				Excluded: []string{"cosmos1aqn8ynvr3jmq67879qulzrwhchq5dtrvh6h4er"},
			},
		},
	}
	yamlData, err := yaml.Marshal(&sampleConfig)
	require.NoError(t, err)
	sampleConfigPath := filepath.Join(t.TempDir(), "config.yml")
	err = os.WriteFile(sampleConfigPath, yamlData, 0o644)

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfig(tt.filename)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, tt.err, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want.AirdropToken, got.AirdropToken)
			require.Equal(t, tt.want.DustWallet, got.DustWallet)
			for i, wantSnapshot := range tt.want.Snapshots {
				require.Equal(t, wantSnapshot.Chain, got.Snapshots[i].Chain)
				require.Equal(t, wantSnapshot.Denom, got.Snapshots[i].Denom)
				require.Equal(t, wantSnapshot.Formula, got.Snapshots[i].Formula)
				require.Equal(t, wantSnapshot.Date.Unix(), got.Snapshots[i].Date.Unix())
				require.Equal(t, wantSnapshot.Types, got.Snapshots[i].Types)
				require.EqualValues(t, wantSnapshot.Excluded, got.Snapshots[i].Excluded)

			}
		})
	}
}
