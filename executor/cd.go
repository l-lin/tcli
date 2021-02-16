package executor

import (
	"fmt"
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
		fmt.Fprintf(c.errOutput, "%v\n", err)
		return c.currentBoard, c.currentList
	}

	if boardName == "" {
		return nil, nil
	}

	var board *trello.Board
	if board, err = c.tr.FindBoard(boardName); err != nil || board == nil {
		fmt.Fprintf(c.errOutput, "no board found with name '%s'\n", boardName)
		return c.currentBoard, c.currentList
	}

	if listName == "" {
		return board, nil
	}

	var list *trello.List
	if list, err = c.tr.FindList(board.ID, listName); err != nil || list == nil {
		fmt.Fprintf(c.errOutput, "no list found with name '%s'\n", listName)
		return board, c.currentList
	}

	if cardName != "" {
		fmt.Fprintf(c.errOutput, "cannot cd on card\n")
		return c.currentBoard, c.currentList
	}

	return board, list
}
