package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
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
	if e := executor.New(*container.Conf, "ls", container.TrelloRepository, container.Renderer, nil, nil); e != nil {
		e.Execute(strings.Join(args, " "))
	} else {
		log.Fatal().
			Stack().
			Str("cmd", "ls").
			Msg("executor not found")
	}
}
