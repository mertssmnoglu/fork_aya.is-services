package subcommands

import (
	"context"
	"fmt"

	"github.com/eser/ajan/lib"
	"github.com/spf13/cobra"
)

func CmdId() *cobra.Command {
	var flagCount int

	idCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "id",
		Short: "Generates id",
		Long:  "Generates id",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execId(cmd.Context(), flagCount)
		},
	}

	idCmd.Flags().IntVarP(&flagCount, "count", "n", 1, "count of ids will be generated")

	return idCmd
}

func execId(_ context.Context, count int) error {
	for range count {
		id := lib.IdsGenerateUnique()

		fmt.Println(id) //nolint:forbidigo
	}

	return nil
}
