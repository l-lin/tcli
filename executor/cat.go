package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
)

type cat struct {
	executor
}

func (c cat) Execute(arg string) (currentBoard *trello.Board, currentList *trello.List) {
	currentBoard = c.currentBoard
	currentList = c.currentList
	if arg == "" {
		return
	}

	pathResolver := trello.NewPathResolver(currentBoard, currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		log.Debug().
			Err(err).
			Str("arg", arg).
			Msg("could not resolve path")
		return
	}

	board, err := c.tr.FindBoard(boardName)
	if err != nil {
		fmt.Fprintf(c.stderr, "no board found with name '%s'\n", boardName)
		return
	}

	if listName == "" {
		fmt.Fprintf(c.stdout, "%s\n", c.r.RenderBoard(*board))
		return
	}

	var list *trello.List
	if list, err = c.tr.FindList(board.ID, listName); err != nil || list == nil {
		fmt.Fprintf(c.stderr, "no list found with name '%s'\n", listName)
		return
	}

	if cardName == "" {
		fmt.Fprintf(c.stdout, "%s\n", c.r.RenderList(*list))
		return
	}

	var card *trello.Card
	if card, err = c.tr.FindCard(list.ID, cardName); err != nil || card == nil {
		fmt.Fprintf(c.stderr, "no card found with name '%s'\n", cardName)
		return
	}
	fmt.Fprintf(c.stdout, "%s\n", c.r.RenderCard(*card))
	return
}
