package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli-plugin-airdrop/pkg/config"
	"github.com/ignite/cli-plugin-airdrop/pkg/genesis"
	"github.com/ignite/cli-plugin-airdrop/pkg/snapshot"
)

func NewAirdrop() *cobra.Command {
	c := &cobra.Command{
		Use:   "airdrop",
		Short: "Utility tool to create snapshots for an airdrop",
	}

	c.AddCommand(
		NewAirdropGenerate(),
		NewAirdropRaw(),
		NewAirdropProcess(),
		NewAirdropGenesis(),
	)

	return c
}

func NewAirdropGenerate() *cobra.Command {
	return &cobra.Command{
		Use:   "generate [airdrop-config] [input-genesis]",
		Short: "Utility tool to create snapshots for an airdrop",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				airdropConfig = args[0]
				inputGenesis  = args[1]
			)

			c, err := config.ParseConfig(airdropConfig)
			if err != nil {
				return err
			}

			genState, err := genesis.GetGenStateFromPath(inputGenesis)
			if err != nil {
				return err
			}

			s, err := snapshot.Generate(genState.AppState)
			if err != nil {
				return err
			}

			filters := make(snapshot.Filters, 0)
			for _, snap := range c.Snapshots {
				filter := s.Filter(snapshot.FilterType(snap.Type), snap.Denom, snap.Formula, snap.Excluded)
				filters = append(filters, filter)
			}
			filter := filters.Sum()

			if err := genState.AddFromClaimRecord(c.AirdropToken, filter.ClaimRecords()); err != nil {
				return err
			}

			// export snapshot json
			genesisJSON, err := json.MarshalIndent(genState, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			cmd.Println(string(genesisJSON))
			return nil
		},
	}
}

func NewAirdropRaw() *cobra.Command {
	return &cobra.Command{
		Use:   "raw [input-genesis]",
		Short: "Generate raw airdrop data based on the input genesis",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			genState, err := genesis.GetGenStateFromPath(args[0])
			if err != nil {
				return err
			}

			s, err := snapshot.Generate(genState.AppState)
			if err != nil {
				return err
			}

			// export snapshot json
			snapshotJSON, err := json.MarshalIndent(s, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			cmd.Println(string(snapshotJSON))
			return nil
		},
	}
}

func NewAirdropProcess() *cobra.Command {
	return &cobra.Command{
		Use:   "process [airdrop-config] [raw-snapshot]",
		Short: "Process the airdrop data based on the config file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				airdropConfig = args[0]
				rawSnapshot   = args[1]
			)

			c, err := config.ParseConfig(airdropConfig)
			if err != nil {
				return err
			}

			s, err := snapshot.ParseSnapshot(rawSnapshot)
			if err != nil {
				return err
			}

			filters := make(snapshot.Filters, 0)
			for _, snap := range c.Snapshots {
				filter := s.Filter(snapshot.FilterType(snap.Type), snap.Denom, snap.Formula, snap.Excluded)
				filters = append(filters, filter)
			}
			filter := filters.Sum()

			// export filter json
			filterJSON, err := json.MarshalIndent(filter, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			cmd.Println(string(filterJSON))
			return nil
		},
	}
}

func NewAirdropGenesis() *cobra.Command {
	return &cobra.Command{
		Use:   "genesis [airdrop-config] [raw-snapshot] [input-genesis]",
		Short: "Generate a genesis based on processed files and airdrop config",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// airdropConfig = args[0]
			// rawSnapshot   = args[1]
			inputGenesis := args[2]

			//c, err := config.ParseConfig(airdropConfig)
			//if err != nil {
			//	return err
			//}

			genState, err := genesis.GetGenStateFromPath(inputGenesis)
			if err != nil {
				return err
			}

			//s, err := snapshot.ParseSnapshot(rawSnapshot)
			//if err != nil {
			//	return err
			//}

			//if err := genState.AddFromClaimRecord(c.AirdropToken, filter.ClaimRecords()); err != nil {
			//	return err
			//}

			// export snapshot json
			genesisJSON, err := json.MarshalIndent(genState, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			cmd.Println(string(genesisJSON))
			return nil
		},
	}
}
