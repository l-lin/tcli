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
		log.Debug().
			Err(err).
			Str("boardName", boardName).
			Msg("could not find board")
		return
	}

	if listName == "" {
		l.renderLists(*board)
		return
	}

	var list *trello.List
	if list, err = l.tr.FindList(board.ID, listName); err != nil || list == nil {
		log.Debug().
			Err(err).
			Str("idBoard", board.ID).
			Str("name", arg).
			Msg("no list found")
		return
	}

	if cardName == "" {
		l.renderCards(*list)
		return
	}

	var card *trello.Card
	if card, err = l.tr.FindCard(list.ID, cardName); err != nil || card == nil {
		log.Debug().
			Err(err).
			Str("idList", list.ID).
			Str("name", arg).
			Msg("no list found")
		return
	}
	l.renderCard(*card)
	return
}

func (l ls) renderBoards() {
	boards, err := l.tr.GetBoards()
	if err != nil {
		log.Debug().
			Err(err).
			Msg("could not fetch boards")
	} else {
		fmt.Printf("%s", l.r.RenderBoards(boards))
	}
}

func (l ls) renderLists(board trello.Board) {
	lists, err := l.tr.GetLists(board.ID)
	if err != nil {
		log.Debug().
			Err(err).
			Str("idBoard", board.ID).
			Msg("could not fetch lists")
	} else {
		fmt.Printf("%s", l.r.RenderLists(lists))
	}
}

func (l ls) renderCards(list trello.List) {
	cards, err := l.tr.GetCards(list.ID)
	if err != nil {
		log.Debug().
			Err(err).
			Str("idList", list.ID).
			Msg("could not fetch cards")
	} else {
		fmt.Printf("%s", l.r.RenderCards(cards))
	}
}

func (l ls) renderCard(card trello.Card) {
	fmt.Printf("%s", l.r.RenderCards(trello.Cards{card}))
}
