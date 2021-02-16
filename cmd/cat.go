package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

func NewCatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cat",
		Short: "Show resource content info",
		Run:   runCat,
		Args:  cobra.MinimumNArgs(1),
		Example: `
  # show 'my-board' info
  tcli cat /my-board

  # show 'my-list' info from 'my-board'
  tcli cat /my-board/my-list

  # show 'my-card' info'
  tcli cat /my-board/my-list/my-card`,
	}
}

func runCat(_ *cobra.Command, args []string) {
	if e := executor.New("cat", container.TrelloRepository, container.Renderer, nil, nil); e != nil {
		e.Execute(strings.Join(args, " "))
	} else {
		log.Fatal().
			Str("cmd", "cat").
			Msg("executor not found")
	}
}
