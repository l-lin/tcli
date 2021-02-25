package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
)

type mv struct {
	executor
}

func (m mv) Execute(args []string) {
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

	execSource := start(m.tr).
		resolvePath(m.session, args[0]).
		thenFindBoard().
		thenFindList().
		thenFindCard()
	if execSource.err != nil {
		fmt.Fprintf(m.stderr, "%s\n", execSource.err)
		return
	}
	sourceCard := execSource.session.Card

	execDest := start(m.tr).
		resolvePath(m.session, args[1]).
		thenFindBoard().
		thenFindList().
		doOnCardName(func(cardName string, session *trello.Session) {
			updateCard := trello.NewUpdateCard(*sourceCard)
			updateCard.IDList = session.List.ID
			if cardName != "" {
				updateCard.Name = cardName
			}
			if _, err := m.tr.UpdateCard(updateCard); err != nil {
				fmt.Fprintf(m.stderr, "could not update card: %v\n", err)
			}
		})
	if execDest.err != nil {
		fmt.Fprintf(m.stderr, "%s\n", execDest.err)
	} else if !execDest.isFinished {
		fmt.Fprintf(m.stderr, "comment move not supported yet\n")
	}
}
