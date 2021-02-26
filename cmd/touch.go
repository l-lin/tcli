package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/spf13/cobra"
	"os"
)

func NewTouchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "touch",
		Short: "create new resource",
		Run:   runTouch,
		Args:  cobra.MinimumNArgs(1),
		Example: `
  # create new card with name 'card'
  tcli touch /board/list/card`,
	}
}

func runTouch(_ *cobra.Command, args []string) {
	e := executor.New(*container.Conf, "touch", container.TrelloRepository, container.Renderer, nil, os.Stdout, os.Stderr)
	e.Execute(args)
}
