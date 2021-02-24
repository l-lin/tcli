package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/spf13/cobra"
)

func NewCPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cat",
		Short: "Show resource content info",
		Run:   runCP,
		Args:  cobra.ExactArgs(2),
		Example: `
  # copy card on same list
  tcli cp /board/list/card /board/list

  # copy card on same list with another name
  tcli cp /board/list/card /board/list/new-card-name

  # copy card on different list
  tcli cp /board/list/card /board/another-list/new-card-name`,
	}
}

func runCP(_ *cobra.Command, args []string) {
	e := executor.New(*container.Conf, "cp", container.TrelloRepository, container.Renderer, nil)
	e.Execute(args)
}
