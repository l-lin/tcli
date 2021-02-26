package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/spf13/cobra"
	"os"
)

func NewRMCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm",
		Short: "Archive resource",
		Run:   runRM,
		Args:  cobra.MinimumNArgs(1),
		Example: `
  # archive card 'card'
  tcli rm /board/list/card`,
	}
}

func runRM(_ *cobra.Command, args []string) {
	e := executor.New(*container.Conf, "rm", container.TrelloRepository, container.Renderer, nil, os.Stdout, os.Stderr)
	e.Execute(args)
}
