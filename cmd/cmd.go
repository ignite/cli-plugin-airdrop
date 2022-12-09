package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewAirdrop() *cobra.Command {
	return &cobra.Command{
		Use:   "airdrop",
		Short: "Command description...",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hello from airdrop command")
			return nil
		},
	}
}
