package executor

import (
	"fmt"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
)

type BoardsExecutor interface {
	Execute(arg string) *trello.Board
}

func NewBoardsExecutor(cmd string, tr trello.Repository, r renderer.Renderer) BoardsExecutor {
	switch cmd {
	case "cd":
		return &cdBoardsExecutor{
			Executor{
				tr: tr,
				r:  r,
			},
		}
	case "ls":
		return &lsBoardsExecutor{Executor{
			tr: tr,
			r:  r,
		}}
	}
	return nil
}

type cdBoardsExecutor struct {
	Executor
}

func (c cdBoardsExecutor) Execute(arg string) *trello.Board {
	if arg == "" {
		return nil
	}
	board, err := c.tr.FindBoard(arg)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", arg).
			Msg("no board found")
		return nil
	}
	return board
}

type lsBoardsExecutor struct {
	Executor
}

func (l lsBoardsExecutor) Execute(arg string) *trello.Board {
	if arg == "" {
		boards, err := l.tr.GetBoards()
		if err != nil {
			log.Err(err).Msg("could not fetch boards")
			return nil
		}
		fmt.Printf("%s", l.r.RenderBoards(boards))
		return nil
	}
	board, err := l.tr.FindBoard(arg)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", arg).
			Msg("no board found")
		return nil
	}
	var lists trello.Lists
	lists, err = l.tr.GetLists(board.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("idBoard", board.ID).
			Str("name", arg).
			Msg("could not find lists")
		return nil
	}
	fmt.Printf("%s", l.r.RenderLists(lists))
	return nil
}
