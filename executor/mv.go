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

	sourceCard, err := m.getCardFromArg(args[0])
	if err != nil {
		fmt.Fprintf(m.stderr, "%s\n", err.Error())
		return
	}
	destList, destCardName, err := m.getListAndCardNameFromArg(args[1])
	if err != nil {
		fmt.Fprintf(m.stderr, "%s\n", err.Error())
		return
	}
	updateCard := trello.NewUpdateCard(*sourceCard)
	updateCard.IDList = destList.ID
	if destCardName != "" {
		updateCard.Name = destCardName
	}
	if _, err = m.tr.UpdateCard(updateCard); err != nil {
		fmt.Fprintf(m.stderr, "could not update card: %v\n", err)
	}
	return
}
