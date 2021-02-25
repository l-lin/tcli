package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
)

type cp struct {
	executor
}

func (c cp) Execute(args []string) {
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

	execSource := start(c.tr).
		resolvePath(c.session, args[0]).
		thenFindBoard().
		thenFindList().
		thenFindCard()
	if execSource.err != nil {
		fmt.Fprintf(c.stderr, "%s\n", execSource.err)
		return
	}
	sourceCard := execSource.session.Card

	execDest := start(c.tr).
		resolvePath(c.session, args[1]).
		thenFindBoard().
		thenFindList().
		doOnCardName(func(cardName string, session *trello.Session) {
			var createCard trello.CreateCard
			createCard = trello.NewCreateCard(*sourceCard)
			createCard.IDList = session.List.ID
			if cardName != "" {
				createCard.Name = cardName
			}
			if _, err := c.tr.CreateCard(createCard); err != nil {
				fmt.Fprintf(c.stderr, "could not copy card '%s': %v\n", sourceCard.Name, err)
				return
			}
		})
	if execDest.err != nil {
		fmt.Fprintf(c.stderr, "%s\n", execDest.err)
	} else if !execDest.isFinished {
		fmt.Fprintf(c.stderr, "comment copy not supported yet\n")
	}
}
