package subcommands

import (
	"context"

	"github.com/eser/aya.is-services/pkg/api/adapters/appcontext"
	"github.com/eser/aya.is-services/pkg/api/adapters/storage"
	"github.com/eser/aya.is-services/pkg/api/business/profiles"
	"github.com/spf13/cobra"
)

func CmdProfilesList() *cobra.Command {
	profilesListCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "list",
		Short: "Lists profiles",
		Long:  "Lists all available profiles registered on the site",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execProfilesList(cmd.Context())
		},
	}

	return profilesListCmd
}

func execProfilesList(ctx context.Context) error {
	appContext, err := appcontext.NewAppContext(ctx)
	if err != nil {
		return err //nolint:wrapcheck
	}

	store, err := storage.NewFromDefault(appContext.Data)
	if err != nil {
		panic(err)
	}

	service := profiles.NewService(appContext.Logger, store)

	records, err := service.List(ctx)
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		appContext.Logger.InfoContext(ctx, "profile entry", "profile", record)
	}

	return nil
}
