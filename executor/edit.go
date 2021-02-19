package executor

import (
	"fmt"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"github.com/manifoldco/promptui"
	"io"
)

type edit struct {
	executor
	editor       Editor
	stdin        io.ReadCloser
	editRenderer renderer.Edit
}

func (e edit) Execute(arg string) (currentBoard *trello.Board, currentList *trello.List) {
	currentBoard = e.currentBoard
	currentList = e.currentList

	pathResolver := trello.NewPathResolver(e.currentBoard, e.currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		fmt.Fprintf(e.stderr, "%v\n", err)
		return e.currentBoard, e.currentList
	}

	if boardName == "" {
		fmt.Fprintf(e.stderr, "nothing to edit\n")
		return
	}

	var board *trello.Board
	if board, err = e.tr.FindBoard(boardName); err != nil || board == nil {
		fmt.Fprintf(e.stderr, "no board found with name '%s'\n", boardName)
		return
	}

	if listName == "" {
		fmt.Fprintf(e.stderr, "board edition not implemented yet\n")
		return
	}

	var list *trello.List
	if list, err = e.tr.FindList(board.ID, listName); err != nil || list == nil {
		fmt.Fprintf(e.stderr, "no list found with name '%s'\n", listName)
		return
	}

	if cardName == "" {
		fmt.Fprintf(e.stderr, "list edition not implemented yet\n")
		return
	}

	var card *trello.Card
	if card, err = e.tr.FindCard(list.ID, cardName); err != nil || card == nil {
		fmt.Fprintf(e.stderr, "no card found with name '%s'\n", cardName)
		return
	}
	if err = e.editCard(*card); err != nil {
		fmt.Fprintf(e.stderr, "could not edit card '%s': %v\n", cardName, err)
	}
	return
}

func (e edit) editCard(card trello.Card) (err error) {
	var lists trello.Lists
	if lists, err = e.tr.FindLists(card.IDBoard); err != nil {
		return
	}

	cte := trello.NewCardToEdit(card)
	var in []byte
	if in, err = e.editRenderer.Marshal(cte, lists); err != nil {
		return
	}

	var out []byte
	if out, err = e.editor.Edit(in); err != nil {
		return
	}

	var editedCard trello.CardToEdit
	if err = e.editRenderer.Unmarshal(out, &editedCard); err != nil {
		return
	}
	updatedCard := trello.NewUpdateCard(card)
	updatedCard.Name = editedCard.Name
	updatedCard.Description = editedCard.Desc
	updatedCard.Closed = editedCard.Closed
	updatedCard.IDList = editedCard.IDList

	prompt := promptui.Prompt{
		Label:     "Do you want to update the card?",
		IsConfirm: true,
		Stdin:     e.stdin,
	}
	if _, err = prompt.Run(); err != nil {
		fmt.Fprintf(e.stdout, "card '%s' not updated\n", card.Name)
		return nil
	}
	_, err = e.tr.UpdateCard(updatedCard)
	return
}
