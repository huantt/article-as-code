package main

import (
	"fmt"
	"github.com/huantt/acc/src/cmd"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("%s", r))
			os.Exit(1)
		}
	}()

	var loggingLevel = new(slog.LevelVar)
	loggingLevel.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: loggingLevel}))
	slog.SetDefault(logger)

	root := &cobra.Command{}
	root.AddCommand(cmd.CollectArticles("collect-articles"))
	root.AddCommand(cmd.SyncArticles("sync-articles"))
	err := root.Execute()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
