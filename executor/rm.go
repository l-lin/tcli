package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/manifoldco/promptui"
	"io"
)

type rm struct {
	executor
	stdin io.ReadCloser
}

func (r rm) Execute(arg string) (currentBoard *trello.Board, currentList *trello.List) {
	currentBoard = r.currentBoard
	currentList = r.currentList

	if arg == "" {
		fmt.Fprintf(r.stderr, "missing card operand")
		return
	}
	pathResolver := trello.NewPathResolver(r.currentBoard, r.currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		fmt.Fprintf(r.stderr, "%v\n", err)
		return r.currentBoard, r.currentList
	}

	if boardName == "" {
		fmt.Fprintf(r.stderr, "nothing to archive\n")
		return
	}

	var board *trello.Board
	if board, err = r.tr.FindBoard(boardName); err != nil || board == nil {
		fmt.Fprintf(r.stderr, "no board found with name '%s'\n", boardName)
		return
	}

	if listName == "" {
		fmt.Fprintf(r.stderr, "board archiving not implemented yet\n")
		return
	}

	var list *trello.List
	if list, err = r.tr.FindList(board.ID, listName); err != nil || list == nil {
		fmt.Fprintf(r.stderr, "no list found with name '%s'\n", listName)
		return
	}

	if cardName == "" {
		fmt.Fprintf(r.stderr, "list archiving not implemented yet\n")
		return
	}

	var card *trello.Card
	if card, err = r.tr.FindCard(list.ID, cardName); err != nil || card == nil {
		fmt.Fprintf(r.stderr, "no card found with name '%s'\n", cardName)
		return
	}
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Archive card '%s'?", card.Name),
		IsConfirm: true,
		Stdin:     r.stdin,
	}
	if _, err = prompt.Run(); err != nil {
		return
	}
	updatedCard := trello.NewUpdateCard(*card)
	updatedCard.Closed = true
	_, err = r.tr.UpdateCard(updatedCard)
	return
}
