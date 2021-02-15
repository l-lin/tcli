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
		Short: "Display info",
		Run:   runLS,
		Args:  cobra.MinimumNArgs(1),
	}
}

func runLS(_ *cobra.Command, args []string) {
	if e := executor.New("ls", container.TrelloRepository, container.Renderer, nil, nil); e != nil {
		e.Execute(strings.Join(args, " "))
	} else {
		log.Fatal().
			Str("cmd", "ls").
			Msg("executor not found")
	}
}
