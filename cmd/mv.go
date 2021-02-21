package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewMVCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mv",
		Short: "Move resource",
		Run:   runMV,
		Example: `
  # archive card 'my-card'
  tcli mv /source-board/source-list/source-card /target-board/target-list`,
	}
}

func runMV(_ *cobra.Command, args []string) {
	if e := executor.New(*container.Conf, "mv", container.TrelloRepository, container.Renderer, nil, nil); e != nil {
		e.Execute(args)
	} else {
		log.Fatal().
			Stack().
			Str("cmd", "mv").
			Msg("executor not found")
	}
}
