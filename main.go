package main

import (
	"fmt"
	cmd2 "github.com/huantt/article-as-code/cmd"
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
	root.AddCommand(cmd2.CollectArticles("collect"))
	root.AddCommand(cmd2.SyncArticles("sync"))
	err := root.Execute()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
