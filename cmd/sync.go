package cmd

import (
	"context"
	"github.com/huantt/acc/handler/synchronizer"
	"github.com/huantt/acc/impl/forem"
	"github.com/huantt/acc/impl/hashnode"
	forem2 "github.com/huantt/acc/pkg/forem"
	"github.com/spf13/cobra"
)

func SyncArticles(use string) *cobra.Command {
	var articleFolder string
	var destination string
	var username string
	var authToken string
	var maxRps int

	command := &cobra.Command{
		Use: use,
		Run: func(cmd *cobra.Command, args []string) {
			handler := synchronizer.NewHandler()

			// Register article services
			devtoService := forem.NewAuthenticatedService(forem2.DevToEndpoint, maxRps, username, authToken)
			handler.RegisterService("dev.to", devtoService)
			hashnodeService := hashnode.NewService(username, authToken)
			handler.RegisterService("hashnode.dev", hashnodeService)

			// Sync
			err := handler.Sync(context.Background(), articleFolder, destination)
			if err != nil {
				panic(err)
			}
		},
	}

	command.Flags().StringVarP(&articleFolder, "article-folder", "f", "data/articles", "Article folder")
	command.Flags().StringVarP(&destination, "destination", "d", "", "Destination: dev.to/hashnode.dev")
	command.Flags().StringVarP(&username, "username", "u", "", "Username")
	command.Flags().StringVarP(&authToken, "auth-token", "a", "", "Auth token")
	command.Flags().IntVar(&maxRps, "rps", 5, "Limit concurrent requests")

	if err := command.MarkFlagRequired("auth-token"); err != nil {
		panic(err)
	}
	if err := command.MarkFlagRequired("destination"); err != nil {
		panic(err)
	}
	if err := command.MarkFlagRequired("username"); err != nil {
		panic(err)
	}
	if err := command.MarkFlagRequired("username"); err != nil {
		panic(err)
	}

	return command
}
