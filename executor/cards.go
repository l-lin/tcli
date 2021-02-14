package executor

import (
	"fmt"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
)

type CardExecutor interface {
	Execute(arg string) (*trello.Board, *trello.List)
}

func NewCardExecutor(cmd string, tr trello.Repository, r renderer.Renderer, currentBoard trello.Board, currentList trello.List) CardExecutor {
	switch cmd {
	case "cd":
		return &cdCardExecutor{
			Executor: Executor{
				tr: tr,
				r:  r,
			},
			currentBoard: currentBoard,
			currentList:  currentList,
		}
	case "ls":
		return &lsCardExecutor{
			Executor: Executor{
				tr: tr,
				r:  r,
			},
			currentBoard: currentBoard,
			currentList:  currentList,
		}
	case "edit":
		return &editCardExecutor{
			Executor: Executor{
				tr: tr,
				r:  r,
			},
			currentBoard: currentBoard,
			currentList:  currentList,
		}
	}
	return nil
}

type cdCardExecutor struct {
	Executor
	currentBoard trello.Board
	currentList  trello.List
}

func (c cdCardExecutor) Execute(arg string) (*trello.Board, *trello.List) {
	if arg == "" {
		log.Debug().Msg("returning to top parent")
		return nil, nil
	}
	if arg == ".." {
		return &c.currentBoard, nil
	}
	log.Warn().
		Str("arg", arg).
		Msg("invalid argument")
	return &c.currentBoard, &c.currentList
}

type lsCardExecutor struct {
	Executor
	currentBoard trello.Board
	currentList  trello.List
}

func (l lsCardExecutor) Execute(arg string) (*trello.Board, *trello.List) {
	if arg == "" {
		cards, err := l.tr.GetCards(l.currentList.ID)
		if err != nil {
			log.Err(err).
				Str("idList", l.currentList.ID).
				Msg("could not fetch cards")
			return &l.currentBoard, &l.currentList
		}
		fmt.Printf(l.r.RenderCards(cards))
		return &l.currentBoard, &l.currentList
	}
	if arg == ".." {
		lists, err := l.tr.GetLists(l.currentBoard.ID)
		if err != nil {
			log.Err(err).
				Str("idBoard", l.currentBoard.ID).
				Msg("could not fetch lists")
			return &l.currentBoard, &l.currentList
		}
		fmt.Printf(l.r.RenderLists(lists))
		return &l.currentBoard, &l.currentList
	}
	card, err := l.tr.FindCard(l.currentList.ID, arg)
	if err != nil {
		log.Error().
			Err(err).
			Str("idList", l.currentList.ID).
			Str("name", arg).
			Msg("no card found")
		return &l.currentBoard, &l.currentList
	}
	fmt.Printf(l.r.RenderCard(*card))
	return &l.currentBoard, &l.currentList
}

type editCardExecutor struct {
	Executor
	currentBoard trello.Board
	currentList  trello.List
}

func (e editCardExecutor) Execute(arg string) (*trello.Board, *trello.List) {
	// TODO
	return &e.currentBoard, &e.currentList
}
