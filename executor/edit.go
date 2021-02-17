package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
	"io"
)

type edit struct {
	executor
	editor Editor
	stdin  io.ReadCloser
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
	var updatedCard *trello.UpdateCard
	if updatedCard, err = e.editCard(trello.NewUpdateCard(*card)); err != nil {
		fmt.Fprintf(e.stderr, "could not edit card '%s': %v\n", cardName, err)
	}
	prompt := promptui.Prompt{
		Label:     "Do you want to update the card?",
		IsConfirm: true,
		Stdin:     e.stdin,
	}
	if _, err = prompt.Run(); err != nil {
		fmt.Fprintf(e.stdout, "card '%s' not updated\n", cardName)
		return
	}
	if _, err = e.tr.UpdateCard(*updatedCard); err != nil {
		fmt.Fprintf(e.stderr, "could not update card '%s': %v\n", cardName, err)
	}
	return
}

func (e edit) editCard(updateCard trello.UpdateCard) (*trello.UpdateCard, error) {
	cte := newCardToEdit(updateCard)
	var in []byte
	var err error
	if in, err = yaml.Marshal(cte); err != nil {
		return nil, err
	}

	var out []byte
	if out, err = e.editor.Edit(in); err != nil {
		return nil, err
	}

	var editedCard cardToEdit
	if err = yaml.Unmarshal(out, &editedCard); err != nil {
		return nil, err
	}
	result := trello.CopyUpdateCard(updateCard)
	result.Name = editedCard.Name
	result.Description = editedCard.Description
	result.Closed = editedCard.Closed
	return &result, nil
}

type cardToEdit struct {
	Name        string `yaml:"name"`
	Description string `yaml:"desc"`
	Closed      bool   `yaml:"closed"`
}

func newCardToEdit(updateCard trello.UpdateCard) cardToEdit {
	return cardToEdit{
		Name:        updateCard.Name,
		Description: updateCard.Description,
		Closed:      updateCard.Closed,
	}
}
