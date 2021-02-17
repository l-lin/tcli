package executor

import (
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"os"
)

var Factories = []Factory{
	{
		Cmd:         "help",
		Description: "display help",
		Create: func(tr trello.Repository, r renderer.Renderer, currentBoard *trello.Board, currentList *trello.List) Executor {
			return &help{
				stdout:       os.Stdout,
				currentBoard: currentBoard,
				currentList:  currentList,
			}
		},
	},
	{
		Cmd:         "exit",
		Description: "exit CLI",
		Create: func(tr trello.Repository, r renderer.Renderer, currentBoard *trello.Board, currentList *trello.List) Executor {
			return &exit{}
		},
	},
	{
		Cmd:         "cd",
		Description: "change level in the hierarchy",
		Create: func(tr trello.Repository, r renderer.Renderer, currentBoard *trello.Board, currentList *trello.List) Executor {
			return &cd{executor{
				tr:           tr,
				r:            r,
				currentBoard: currentBoard,
				currentList:  currentList,
				stdout:       os.Stdout,
				stderr:       os.Stderr,
			}}
		},
	},
	{
		Cmd:         "ls",
		Description: "list resource content",
		Create: func(tr trello.Repository, r renderer.Renderer, currentBoard *trello.Board, currentList *trello.List) Executor {
			return &ls{executor{
				tr:           tr,
				r:            r,
				currentBoard: currentBoard,
				currentList:  currentList,
				stdout:       os.Stdout,
				stderr:       os.Stderr,
			}}
		},
	},
	{
		Cmd:         "cat",
		Description: "show resource content info",
		Create: func(tr trello.Repository, r renderer.Renderer, currentBoard *trello.Board, currentList *trello.List) Executor {
			return &cat{executor{
				tr:           tr,
				r:            r,
				currentBoard: currentBoard,
				currentList:  currentList,
				stdout:       os.Stdout,
				stderr:       os.Stderr,
			}}
		},
	},
	{
		Cmd:         "edit",
		Description: "edit resource content",
		Create: func(tr trello.Repository, r renderer.Renderer, currentBoard *trello.Board, currentList *trello.List) Executor {
			return &edit{
				executor: executor{
					tr:           tr,
					r:            r,
					currentBoard: currentBoard,
					currentList:  currentList,
					stdout:       os.Stdout,
					stderr:       os.Stderr,
				},
				stdin:  os.Stdin,
				editor: NewOsEditor(),
			}
		},
	},
}

type Factory struct {
	Cmd         string
	Description string
	Create      func(tr trello.Repository, r renderer.Renderer, currentBoard *trello.Board, currentList *trello.List) Executor
}
