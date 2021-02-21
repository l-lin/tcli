package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewTouchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "touch",
		Short: "create new resource",
		Run:   runTouch,
		Example: `
  # create new card with name 'my-card'
  tcli touch /my-board/my-list/my-card`,
	}
}

func runTouch(_ *cobra.Command, args []string) {
	if e := executor.New(*container.Conf, "touch", container.TrelloRepository, container.Renderer, nil, nil); e != nil {
		e.Execute(args)
	} else {
		log.Fatal().
			Stack().
			Str("cmd", "touch").
			Msg("executor not found")
	}
}
