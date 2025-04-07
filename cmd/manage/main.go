package main

import (
	"github.com/eser/acik.io/cmd/manage/subcommands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "manaeg",
		Short: "acik.io CLI for site management",
		Long:  `acik.io CLI provides various functionalities for site management including reporting and administration.`,
	}

	rootCmd.AddCommand(subcommands.CmdHealthCheck())

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
