package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewRMCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm",
		Short: "Archive resource",
		Run:   runRM,
		Example: `
  # archive card 'my-card'
  tcli rm /my-board/my-list/my-card`,
	}
}

func runRM(_ *cobra.Command, args []string) {
	if e := executor.New(*container.Conf, "rm", container.TrelloRepository, container.Renderer, nil, nil); e != nil {
		e.Execute(args)
	} else {
		log.Fatal().
			Stack().
			Str("cmd", "rm").
			Msg("executor not found")
	}
}
