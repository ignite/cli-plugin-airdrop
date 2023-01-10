package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ignite/cli-plugin-airdrop/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "cli-plugin-airdrop",
	Short: "debug command for CLI airdrop plugin",
}

func main() {
	rootCmd.AddCommand(cmd.NewAirdrop())
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
