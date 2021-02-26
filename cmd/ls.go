package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/spf13/cobra"
	"os"
)

func NewLSCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List resource content",
		Run:   runLS,
		Example: `
  # show all boards
  tcli ls
  tcli ls /

  # show 'my-board' lists
  tcli ls /my-board

  # show 'my-list' cards
  tcli ls /my-board/my-list`,
	}
}

func runLS(_ *cobra.Command, args []string) {
	e := executor.New(*container.Conf, "ls", container.TrelloRepository, container.Renderer, nil, os.Stdout, os.Stderr)
	e.Execute(args)
}
