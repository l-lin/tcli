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
	destList, destCardName, err := c.getListAndCardNameFromArg(args[1])
	if err != nil {
		fmt.Fprintf(c.stderr, "%s\n", err.Error())
		return
	}
	var createCard trello.CreateCard
	createCard = trello.NewCreateCard(*sourceCard)
	createCard.IDList = destList.ID
	if destCardName != "" {
		createCard.Name = destCardName
	}
	if _, err = c.tr.CreateCard(createCard); err != nil {
		fmt.Fprintf(c.stderr, "could not copy card '%s': %v\n", sourceCard.Name, err)
		return
	}
	return
}
