package executor

import (
	"github.com/l-lin/tcli/conf"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"io"
)

type Executor interface {
	Execute(arg string) (*trello.Board, *trello.List)
}

type executor struct {
	tr           trello.Repository
	r            renderer.Renderer
	currentBoard *trello.Board
	currentList  *trello.List
	stdout       io.Writer
	stderr       io.Writer
}

func New(conf conf.Conf, cmd string, tr trello.Repository, r renderer.Renderer, currentBoard *trello.Board, currentList *trello.List) Executor {
	for _, factory := range Factories {
		if factory.Cmd == cmd {
			return factory.Create(conf, tr, r, currentBoard, currentList)
		}
	}
	return nil
}
