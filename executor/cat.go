package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
)

type cat struct {
	executor
}

func (c cat) Execute(args []string) {
	if len(args) != 0 {
		for _, arg := range args {
			c.execute(arg)
		}
	}
}

func (c cat) execute(arg string) {
	if arg == "" {
		return
	}

	if err := start(c.tr).
		resolvePath(c.session, arg).
		then().
		findBoard().
		doOnBoard(func(board *trello.Board) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderBoard(*board))
		}).
		then().
		findList().
		doOnList(func(list *trello.List) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderList(*list))
		}).
		then().
		findCard().
		doOnCard(func(card *trello.Card) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderCard(*card))
		}).
		then().
		findComment().
		doOnComment(func(comment *trello.Comment) {
			c.renderComment(*comment)
		}).err; err != nil {
		fmt.Fprintf(c.stderr, "%s\n", err)
	}
}

func (c cat) renderComment(comment trello.Comment) {
	reactionSummaries, err := c.tr.FindReactionSummaries(comment.ID)
	if err != nil {
		fmt.Fprintf(c.stderr, "could not fetch reaction summaries for comment '%s': %v\n", comment.ID, err)
	} else {
		fmt.Fprintf(c.stdout, "%s\n", c.r.RenderComment(comment, reactionSummaries))
	}
}
