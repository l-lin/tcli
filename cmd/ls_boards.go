package cmd

import (
	"github.com/l-lin/tcli/executor"
	"github.com/spf13/cobra"
)

func NewLSBoardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "boards",
		Short: "display boards info",
		Run:   runLSBoards,
	}
}

func runLSBoards(_ *cobra.Command, _ []string) {
	e := executor.NewBoardsExecutor("ls", container.TrelloRepository, container.Renderer)
	e.Execute("")
}
