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

	pathResolver := trello.NewPathResolver(c.session)
	p, err := pathResolver.Resolve(arg)
	if err != nil {
		fmt.Fprintf(c.stderr, "%v\n", err)
		return
	}

	board, err := c.tr.FindBoard(p.BoardName)
	if err != nil {
		fmt.Fprintf(c.stderr, "no board found with name '%s'\n", p.BoardName)
		return
	}

	if p.ListName == "" {
		fmt.Fprintf(c.stdout, "%s\n", c.r.RenderBoard(*board))
		return
	}

	var list *trello.List
	if list, err = c.tr.FindList(board.ID, p.ListName); err != nil || list == nil {
		fmt.Fprintf(c.stderr, "no list found with name '%s'\n", p.ListName)
		return
	}

	if p.CardName == "" {
		fmt.Fprintf(c.stdout, "%s\n", c.r.RenderList(*list))
		return
	}

	var card *trello.Card
	if card, err = c.tr.FindCard(list.ID, p.CardName); err != nil || card == nil {
		fmt.Fprintf(c.stderr, "no card found with name '%s'\n", p.CardName)
		return
	}

	if p.CommentID == "" {
		fmt.Fprintf(c.stdout, "%s\n", c.r.RenderCard(*card))
		return
	}

	var comment *trello.Comment
	if comment, err = c.tr.FindComment(card.ID, p.CommentID); err != nil || comment == nil {
		fmt.Fprintf(c.stderr, "no comment found with id '%s'\n", p.CommentID)
		return
	}
	fmt.Fprintf(c.stdout, "%s\n", c.r.RenderComment(*comment))
}
