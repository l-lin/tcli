package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
)

type mv struct {
	executor
}

func (m mv) Execute(args []string) (currentBoard *trello.Board, currentList *trello.List) {
	currentBoard = m.currentBoard
	currentList = m.currentList
	if len(args) == 0 {
		fmt.Fprintf(m.stderr, "missing card source operand\n")
		return
	}
	if len(args) == 1 {
		fmt.Fprintf(m.stderr, "missing list destination operand\n")
		return
	}
	if len(args) > 2 {
		fmt.Fprintf(m.stderr, "command only accepts two arguments\n")
		return
	}

	card, err := m.getCard(args[0])
	if err != nil {
		fmt.Fprintf(m.stderr, "%s\n", err.Error())
		return
	}
	list, err := m.getList(args[1])
	if err != nil {
		fmt.Fprintf(m.stderr, "%s\n", err.Error())
		return
	}
	updateCard := trello.NewUpdateCard(*card)
	updateCard.IDList = list.ID
	if _, err = m.tr.UpdateCard(updateCard); err != nil {
		fmt.Fprintf(m.stderr, "could not update card: %v\n", err)
	}
	return
}

func (m mv) getCard(arg string) (*trello.Card, error) {
	pathResolver := trello.NewPathResolver(m.currentBoard, m.currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		return nil, err
	}

	if boardName == "" || listName == "" || cardName == "" {
		return nil, fmt.Errorf("invalid path")
	}

	var list *trello.List
	if list, err = m.getList(arg); err != nil {
		return nil, err
	}

	var card *trello.Card
	if card, err = m.tr.FindCard(list.ID, cardName); err != nil || card == nil {
		return nil, fmt.Errorf("no card found with name '%s'", cardName)
	}
	return card, nil
}

func (m mv) getList(arg string) (*trello.List, error) {
	pathResolver := trello.NewPathResolver(m.currentBoard, m.currentList)
	boardName, listName, _, err := pathResolver.Resolve(arg)
	if err != nil {
		return nil, err
	}

	if boardName == "" || listName == "" {
		return nil, fmt.Errorf("invalid path")
	}

	var board *trello.Board
	if board, err = m.tr.FindBoard(boardName); err != nil || board == nil {
		return nil, fmt.Errorf("no board found with name '%s'", boardName)
	}

	var list *trello.List
	if list, err = m.tr.FindList(board.ID, listName); err != nil || list == nil {
		return nil, fmt.Errorf("no list found with name '%s'", listName)
	}
	return list, nil
}
