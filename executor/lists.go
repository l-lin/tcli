package executor

import (
	"fmt"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
)

type ListsExecutor interface {
	Execute(arg string) (*trello.Board, *trello.List)
}

func NewListsExecutor(cmd string, tr trello.Repository, r renderer.Renderer, currentBoard trello.Board) ListsExecutor {
	switch cmd {
	case "cd":
		return &cdListsExecutor{
			Executor: Executor{
				tr: tr,
				r:  r,
			},
			currentBoard: currentBoard,
		}
	case "ls":
		return &lsListsExecutor{
			Executor: Executor{
				tr: tr,
				r:  r,
			},
			currentBoard: currentBoard,
		}
	}
	return nil
}

type cdListsExecutor struct {
	Executor
	currentBoard trello.Board
}

func (c cdListsExecutor) Execute(arg string) (*trello.Board, *trello.List) {
	if arg == "" {
		log.Debug().Msg("returning to top parent")
		return nil, nil
	}
	if arg == ".." {
		return nil, nil
	}
	list, err := c.tr.FindList(c.currentBoard.ID, arg)
	if err != nil {
		log.Error().Err(err).
			Str("name", arg).
			Str("idBoard", c.currentBoard.ID).
			Msg("no list found")
		return &c.currentBoard, nil
	}
	return &c.currentBoard, list
}

type lsListsExecutor struct {
	Executor
	currentBoard trello.Board
}

func (l lsListsExecutor) Execute(arg string) (*trello.Board, *trello.List) {
	if arg == "" {
		lists, err := l.tr.GetLists(l.currentBoard.ID)
		if err != nil {
			log.Err(err).
				Str("idBoard", l.currentBoard.ID).
				Msg("could not fetch boards")
			return &l.currentBoard, nil
		}
		fmt.Printf("%s", l.r.RenderLists(lists))
		return &l.currentBoard, nil
	}
	list, err := l.tr.FindList(l.currentBoard.ID, arg)
	if err != nil {
		log.Error().
			Err(err).
			Str("idBoard", l.currentBoard.ID).
			Str("name", arg).
			Msg("no list found")
		return &l.currentBoard, nil
	}
	var cards trello.Cards
	cards, err = l.tr.GetCards(list.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("idList", list.ID).
			Msg("could not find cards")
		return &l.currentBoard, nil
	}
	fmt.Printf("%s", l.r.RenderCards(cards))
	return &l.currentBoard, nil
}
