package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/spf13/cobra"
)

func NewMVCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mv",
		Short: "Move resource",
		Run:   runMV,
		Args:  cobra.ExactArgs(2),
		Example: `
  # move card 'card' to 'target-list'
  tcli mv /source-board/source-list/source-card /target-board/target-list`,
	}
}

func runMV(_ *cobra.Command, args []string) {
	e := executor.New(*container.Conf, "mv", container.TrelloRepository, container.Renderer, nil, nil)
	e.Execute(args)
}
