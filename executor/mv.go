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
		then().
		findBoard().
		then().
		findList().
		then().
		findCard()
	if execSource.err != nil {
		fmt.Fprintf(m.stderr, "%s\n", execSource.err)
		return
	}
	if execSource.p.CommentID != "" {
		fmt.Fprintf(m.stderr, "cannot move comments\n")
		return
	}
	sourceCard := execSource.session.Card

	execDest := start(m.tr).
		resolvePath(m.session, args[1]).
		then().
		findBoard().
		then().
		findList().
		doOnList(func(list *trello.List) {
			m.move("", list.ID, sourceCard)
		}).
		then().
		doOnCardName(func(cardName string, session *trello.Session) {
			m.move(cardName, session.List.ID, sourceCard)
		})
	if execDest.err != nil {
		fmt.Fprintf(m.stderr, "%s\n", execDest.err)
	}
}

func (m mv) move(cardName, idList string, sourceCard *trello.Card) {
	updateCard := trello.NewUpdateCard(*sourceCard)
	updateCard.IDList = idList
	if cardName != "" {
		updateCard.Name = cardName
	}
	if _, err := m.tr.UpdateCard(updateCard); err != nil {
		fmt.Fprintf(m.stderr, "could not update card: %v\n", err)
	}
}
