package executor

import (
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
)

type Executor interface {
	Execute(arg string) (*trello.Board, *trello.List)
}

type executor struct {
	tr           trello.Repository
	r            renderer.Renderer
	currentBoard *trello.Board
	currentList  *trello.List
}

func New(cmd string, tr trello.Repository, r renderer.Renderer, currentBoard *trello.Board, currentList *trello.List) Executor {
	for _, factory := range Factories {
		if factory.Cmd == cmd {
			return factory.Create(tr, r, currentBoard, currentList)
		}
	}
	return nil
}
