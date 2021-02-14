package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

func NewLSListsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lists",
		Short: "display lists info",
		Run:   runLSLists,
		Args:  cobra.MinimumNArgs(1),
	}
}

func runLSLists(_ *cobra.Command, args []string) {
	boardName := strings.Join(args, " ")
	board, err := container.TrelloRepository.FindBoard(boardName)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("boardName", boardName).
			Msg("could not find board")
	}
	e := executor.NewListsExecutor("ls", container.TrelloRepository, container.Renderer, *board)
	e.Execute("")
}
