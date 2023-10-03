package cmd

import (
	"context"
	"github.com/huantt/acc/src/handler/collector"
	"github.com/huantt/acc/src/impl/forem"
	foremsrv "github.com/huantt/acc/src/pkg/forem"
	"github.com/spf13/cobra"
	"log/slog"
)

func CollectArticles(use string) *cobra.Command {
	var articleFolder string
	var username string
	var maxRps int

	command := &cobra.Command{
		Use: use,
		Run: func(cmd *cobra.Command, args []string) {
			devToService := forem.NewService(foremsrv.DevToEndpoint, maxRps)
			handler := collector.NewHandler(devToService, articleFolder)
			err := handler.Collect(context.Background(), username)
			if err != nil {
				panic(err)
			}
			slog.Info("Collected!")
		},
	}

	command.Flags().StringVarP(&articleFolder, "article-folder", "f", "data/articles", "Article folder")
	command.Flags().StringVarP(&username, "username", "u", "", "Username")
	command.Flags().IntVar(&maxRps, "rps", 5, "Limit concurrent requests")
	err := command.MarkFlagRequired("username")
	if err != nil {
		panic(err)
	}
	err = command.MarkFlagRequired("username")
	if err != nil {
		panic(err)
	}

	return command
}
