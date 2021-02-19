package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

func NewEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "edit resource content",
		Run:   runEdit,
		Example: `
  # edit card
  tcli ls /my-board/my-list/my-card`,
	}
}

func runEdit(_ *cobra.Command, args []string) {
	if e := executor.New(*container.Conf, "edit", container.TrelloRepository, container.Renderer, nil, nil); e != nil {
		e.Execute(strings.Join(args, " "))
	} else {
		log.Fatal().
			Stack().
			Str("cmd", "edit").
			Msg("executor not found")
	}
}
