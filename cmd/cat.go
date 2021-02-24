package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/spf13/cobra"
)

func NewCatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cat",
		Short: "Show resource content info",
		Run:   runCat,
		Args:  cobra.MinimumNArgs(1),
		Example: `
  # show 'board' info
  tcli cat /board

  # show 'list' info from 'board'
  tcli cat /board/list

  # show 'card' info'
  tcli cat /board/list/card`,
	}
}

func runCat(_ *cobra.Command, args []string) {
	e := executor.New(*container.Conf, "cat", container.TrelloRepository, container.Renderer, nil)
	e.Execute(args)
}
