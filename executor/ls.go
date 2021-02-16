package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
)

type ls struct {
	executor
}

func (l ls) Execute(arg string) (currentBoard *trello.Board, currentList *trello.List) {
	currentBoard = l.currentBoard
	currentList = l.currentList

	pathResolver := trello.NewPathResolver(currentBoard, currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		log.Debug().
			Err(err).
			Str("arg", arg).
			Msg("could not resolve path")
		return
	}

	if boardName == "" {
		l.renderBoards()
		return
	}

	var board *trello.Board
	if board, err = l.tr.FindBoard(boardName); err != nil || board == nil {
		fmt.Fprintf(l.errOutput, "no board found with name '%s'", boardName)
		return
	}

	if listName == "" {
		l.renderLists(*board)
		return
	}

	var list *trello.List
	if list, err = l.tr.FindList(board.ID, listName); err != nil || list == nil {
		fmt.Fprintf(l.errOutput, "no list found with name '%s'", listName)
		return
	}

	if cardName == "" {
		l.renderCards(*list)
		return
	}

	var card *trello.Card
	if card, err = l.tr.FindCard(list.ID, cardName); err != nil || card == nil {
		fmt.Fprintf(l.errOutput, "no card found with name '%s'", cardName)
		return
	}
	l.renderCard(*card)
	return
}

func (l ls) renderBoards() {
	boards, err := l.tr.GetBoards()
	if err != nil {
		fmt.Fprintf(l.errOutput, "could not fetch boards: %v", err)
	} else {
		fmt.Fprintf(l.output, "%s", l.r.RenderBoards(boards))
	}
}

func (l ls) renderLists(board trello.Board) {
	lists, err := l.tr.GetLists(board.ID)
	if err != nil {
		fmt.Fprintf(l.errOutput, "could not fetch lists for board '%s': %v", board.Name, err)
	} else {
		fmt.Fprintf(l.output, "%s", l.r.RenderLists(lists))
	}
}

func (l ls) renderCards(list trello.List) {
	cards, err := l.tr.GetCards(list.ID)
	if err != nil {
		fmt.Fprintf(l.errOutput, "could not fetch cards for list '%s': %v", list.Name, err)
	} else {
		fmt.Fprintf(l.output, "%s", l.r.RenderCards(cards))
	}
}

func (l ls) renderCard(card trello.Card) {
	fmt.Fprintf(l.output, "%s", l.r.RenderCards(trello.Cards{card}))
}
