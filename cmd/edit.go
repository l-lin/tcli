package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/spf13/cobra"
)

func NewEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "edit resource content",
		Run:   runEdit,
		Args:  cobra.ExactArgs(1),
		Example: `
  # edit card
  tcli edit /board/list/card`,
	}
}

func runEdit(_ *cobra.Command, args []string) {
	e := executor.New(*container.Conf, "edit", container.TrelloRepository, container.Renderer, nil)
	e.Execute(args)
}
