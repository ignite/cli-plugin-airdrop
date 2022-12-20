package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/ignite/cli-plugin-airdrop/cmd"
)

func main() {
	rootCmd := cmd.NewAirdrop()
	if err := svrcmd.Execute(rootCmd, "", ""); err != nil {
		os.Exit(1)
	}
}
