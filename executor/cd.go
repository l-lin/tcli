package executor

import (
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
)

type cd struct {
	executor
}

func (c cd) Execute(arg string) (*trello.Board, *trello.List) {
	if arg == "" {
		log.Debug().Msg("returning to top parent")
		return nil, nil
	}

	pathResolver := trello.NewPathResolver(c.currentBoard, c.currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		log.Debug().Err(err).Str("arg", arg).Msg("could not resolve path")
		return c.currentBoard, c.currentList
	}

	if cardName != "" {
		log.Error().Str("cardName", cardName).Msg("cannot cd on card")
		return c.currentBoard, c.currentList
	}

	if boardName == "" {
		return nil, nil
	}

	var board *trello.Board
	if board, err = c.tr.FindBoard(boardName); err != nil || board == nil {
		log.Debug().
			Err(err).
			Str("boardName", boardName).
			Msg("could not find board")
		return c.currentBoard, c.currentList
	}

	if listName == "" {
		return board, nil
	}

	var list *trello.List
	if list, err = c.tr.FindList(board.ID, listName); err != nil || list == nil {
		log.Debug().
			Err(err).
			Str("idBoard", board.ID).
			Str("listName", listName).
			Msg("could not find list")
		return board, c.currentList
	}
	return board, list
}
