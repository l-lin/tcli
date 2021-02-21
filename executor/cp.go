package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
)

type cp struct {
	executor
}

func (c cp) Execute(args []string) (currentBoard *trello.Board, currentList *trello.List) {
	currentBoard = c.currentBoard
	currentList = c.currentList
	if len(args) == 0 {
		fmt.Fprintf(c.stderr, "missing card source operand\n")
		return
	}
	if len(args) == 1 {
		fmt.Fprintf(c.stderr, "missing destination operand\n")
		return
	}
	if len(args) > 2 {
		fmt.Fprintf(c.stderr, "command only accepts two arguments\n")
		return
	}

	sourceCard, err := c.getCardFromArg(args[0])
	if err != nil {
		fmt.Fprintf(c.stderr, "%s\n", err.Error())
		return
	}
	destList, destCard, err := c.getListAndCardName(args[1])
	if err != nil {
		fmt.Fprintf(c.stderr, "%s\n", err.Error())
		return
	}
	var createCard trello.CreateCard
	createCard = trello.NewCreateCard(*sourceCard)
	createCard.IDList = destList.ID
	if destCard != "" {
		createCard.Name = destCard
	}
	if _, err = c.tr.CreateCard(createCard); err != nil {
		fmt.Fprintf(c.stderr, "could not copy card '%s': %v\n", sourceCard.Name, err)
		return
	}
	return
}

func (c cp) getListAndCardName(arg string) (*trello.List, string, error) {
	pathResolver := trello.NewPathResolver(c.currentBoard, c.currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		return nil, "", err
	}
	if boardName == "" || listName == "" {
		return nil, "", fmt.Errorf("invalid path")
	}
	list, err := c.getList(boardName, listName)
	return list, cardName, err
}
