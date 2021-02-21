package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
)

type ls struct {
	executor
}

func (l ls) Execute(args []string) (currentBoard *trello.Board, currentList *trello.List) {
	currentBoard = l.currentBoard
	currentList = l.currentList
	if len(args) == 0 {
		l.execute("")
	}
	for _, arg := range args {
		l.execute(arg)
	}
	return
}

func (l ls) execute(arg string) (currentBoard *trello.Board, currentList *trello.List) {
	currentBoard = l.currentBoard
	currentList = l.currentList

	pathResolver := trello.NewPathResolver(currentBoard, currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		fmt.Fprintf(l.stderr, "%v\n", err)
		return
	}

	if boardName == "" {
		l.renderBoards()
		return
	}

	var board *trello.Board
	if board, err = l.tr.FindBoard(boardName); err != nil || board == nil {
		fmt.Fprintf(l.stderr, "no board found with name '%s'\n", boardName)
		return
	}

	if listName == "" {
		l.renderLists(*board)
		return
	}

	var list *trello.List
	if list, err = l.tr.FindList(board.ID, listName); err != nil || list == nil {
		fmt.Fprintf(l.stderr, "no list found with name '%s'\n", listName)
		return
	}

	if cardName == "" {
		l.renderCards(*list)
		return
	}

	var card *trello.Card
	if card, err = l.tr.FindCard(list.ID, cardName); err != nil || card == nil {
		fmt.Fprintf(l.stderr, "no card found with name '%s'\n", cardName)
		return
	}
	l.renderCard(*card)
	return
}

func (l ls) renderBoards() {
	boards, err := l.tr.FindBoards()
	if err != nil {
		fmt.Fprintf(l.stderr, "could not fetch boards: %v\n", err)
	} else {
		fmt.Fprintf(l.stdout, "%s\n", l.r.RenderBoards(boards))
	}
}

func (l ls) renderLists(board trello.Board) {
	lists, err := l.tr.FindLists(board.ID)
	if err != nil {
		fmt.Fprintf(l.stderr, "could not fetch lists for board '%s': %v\n", board.Name, err)
	} else {
		fmt.Fprintf(l.stdout, "%s\n", l.r.RenderLists(lists))
	}
}

func (l ls) renderCards(list trello.List) {
	cards, err := l.tr.FindCards(list.ID)
	if err != nil {
		fmt.Fprintf(l.stderr, "could not fetch cards for list '%s': %v\n", list.Name, err)
	} else {
		fmt.Fprintf(l.stdout, "%s\n", l.r.RenderCards(cards))
	}
}

func (l ls) renderCard(card trello.Card) {
	fmt.Fprintf(l.stdout, "%s\n", l.r.RenderCards(trello.Cards{card}))
}
