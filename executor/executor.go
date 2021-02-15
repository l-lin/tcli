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
	e := executor{
		tr:           tr,
		r:            r,
		currentBoard: currentBoard,
		currentList:  currentList,
	}
	switch cmd {
	case "cd":
		return &cd{executor: e}
	case "ls":
		return &ls{executor: e}
	}
	return nil
}
