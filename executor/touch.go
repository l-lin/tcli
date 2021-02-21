package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
)

type touch struct {
	executor
}

func (t touch) Execute(args []string) (currentBoard *trello.Board, currentList *trello.List) {
	currentBoard = t.currentBoard
	currentList = t.currentList
	if len(args) == 0 {
		fmt.Fprintf(t.stderr, "missing card operand\n")
		return
	}
	for _, arg := range args {
		t.execute(arg)
	}
	return
}

func (t touch) execute(arg string) (currentBoard *trello.Board, currentList *trello.List) {
	currentBoard = t.currentBoard
	currentList = t.currentList

	if arg == "" {
		fmt.Fprintf(t.stderr, "missing card operand\n")
		return
	}
	pathResolver := trello.NewPathResolver(t.currentBoard, t.currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		fmt.Fprintf(t.stderr, "%v\n", err)
		return t.currentBoard, t.currentList
	}

	if boardName == "" {
		fmt.Fprintf(t.stderr, "nothing to create\n")
		return
	}

	var board *trello.Board
	if board, err = t.tr.FindBoard(boardName); err != nil || board == nil {
		fmt.Fprintf(t.stderr, "no board found with name '%s'\n", boardName)
		return
	}

	if listName == "" {
		fmt.Fprintf(t.stderr, "board creation not implemented yet\n")
		return
	}

	var list *trello.List
	if list, err = t.tr.FindList(board.ID, listName); err != nil || list == nil {
		fmt.Fprintf(t.stderr, "no list found with name '%s'\n", listName)
		return
	}

	if cardName == "" {
		fmt.Fprintf(t.stderr, "list creation not implemented yet\n")
		return
	}

	createCard := trello.CreateCard{
		Name:   cardName,
		IDList: list.ID,
	}
	if _, err = t.tr.CreateCard(createCard); err != nil {
		fmt.Fprintf(t.stderr, "could not create card '%s': %v\n", cardName, err)
	}
	return
}
