package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/spf13/cobra"
	"os"
)

func NewClearCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Clear cache",
		Run:   runClear,
	}
}

func runClear(_ *cobra.Command, args []string) {
	e := executor.New(*container.Conf, "clear", container.TrelloRepository, container.Renderer, nil, os.Stdout, os.Stderr)
	e.Execute(args)
}
