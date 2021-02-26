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
		then().
		findBoard().
		then().
		findList().
		then().
		findCard()
	if execSource.err != nil {
		fmt.Fprintf(c.stderr, "%s\n", execSource.err)
		return
	}
	sourceCard := execSource.session.Card
	var sourceComment *trello.Comment
	if execSource.p.CommentID != "" {
		sourceComment = execSource.then().findComment().comment
	}

	execDest := start(c.tr).
		resolvePath(c.session, args[1]).
		then().
		findBoard().
		then().
		findList().
		then().
		doOnEmptyCardName(func(session *trello.Session) {
			if sourceComment != nil {
				fmt.Fprintf(c.stderr, "cannot copy comment in list\n")
			} else {
				c.createCard("", session, sourceCard)
			}
		}).
		findCard().
		doOnCard(func(destCard *trello.Card) {
			c.createComment(destCard, sourceComment)
		}).
		doOnCardName(func(cardName string, session *trello.Session) {
			c.createCard(cardName, session, sourceCard)
		})
	if execDest.err != nil {
		fmt.Fprintf(c.stderr, "%s\n", execDest.err)
	}
}

func (c cp) createCard(cardName string, session *trello.Session, sourceCard *trello.Card) {
	var createCard trello.CreateCard
	createCard = trello.NewCreateCard(*sourceCard)
	createCard.IDList = session.List.ID
	if cardName != "" {
		createCard.Name = cardName
	}
	if _, err := c.tr.CreateCard(createCard); err != nil {
		fmt.Fprintf(c.stderr, "could not copy card '%s': %v\n", sourceCard.Name, err)
	}
}

func (c cp) createComment(destCard *trello.Card, comment *trello.Comment) {
	createComment := trello.CreateComment{
		IDCard: destCard.ID,
		Text:   comment.Data.Text,
	}
	if _, err := c.tr.CreateComment(createComment); err != nil {
		fmt.Fprintf(c.stderr, "could not copy comment '%s': %v\n", comment.ID, err)
	}
}
