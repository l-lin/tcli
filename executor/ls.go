package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
)

type ls struct {
	executor
}

func (l ls) Execute(args []string) {
	if len(args) == 0 {
		l.execute("")
	}
	for _, arg := range args {
		l.execute(arg)
	}
}

func (l ls) execute(arg string) {
	pathResolver := trello.NewPathResolver(l.session)
	p, err := pathResolver.Resolve(arg)
	if err != nil {
		fmt.Fprintf(l.stderr, "%v\n", err)
		return
	}

	if p.BoardName == "" {
		l.renderBoards()
		return
	}

	var board *trello.Board
	if board, err = l.tr.FindBoard(p.BoardName); err != nil || board == nil {
		fmt.Fprintf(l.stderr, "no board found with name '%s'\n", p.BoardName)
		return
	}

	if p.ListName == "" {
		l.renderLists(*board)
		return
	}

	var list *trello.List
	if list, err = l.tr.FindList(board.ID, p.ListName); err != nil || list == nil {
		fmt.Fprintf(l.stderr, "no list found with name '%s'\n", p.ListName)
		return
	}

	if p.CardName == "" {
		l.renderCards(*list)
		return
	}

	var card *trello.Card
	if card, err = l.tr.FindCard(list.ID, p.CardName); err != nil || card == nil {
		fmt.Fprintf(l.stderr, "no card found with name '%s'\n", p.CardName)
		return
	}

	if p.CommentID == "" {
		l.renderComments(*card)
		return
	}

	var comment *trello.Comment
	if comment, err = l.tr.FindComment(card.ID, p.CommentID); err != nil || comment == nil {
		fmt.Fprintf(l.stderr, "no comment found with id '%s'\n", p.CommentID)
		return
	}
	l.renderComment(*comment)
}

func (l ls) renderBoards() {
	boards, err := l.tr.FindBoards()
	if err != nil {
		fmt.Fprintf(l.stderr, "could not fetch boards: %v\n", err)
	} else {
		fmt.Fprintf(l.stdout, "%s\n", l.r.RenderBoards(boards))
	}
}

func (l ls) renderLists(board trello.Board) {
	lists, err := l.tr.FindLists(board.ID)
	if err != nil {
		fmt.Fprintf(l.stderr, "could not fetch lists for board '%s': %v\n", board.Name, err)
	} else {
		fmt.Fprintf(l.stdout, "%s\n", l.r.RenderLists(lists))
	}
}

func (l ls) renderCards(list trello.List) {
	cards, err := l.tr.FindCards(list.ID)
	if err != nil {
		fmt.Fprintf(l.stderr, "could not fetch cards for list '%s': %v\n", list.Name, err)
	} else {
		fmt.Fprintf(l.stdout, "%s\n", l.r.RenderCards(cards))
	}
}

func (l ls) renderCard(card trello.Card) {
	fmt.Fprintf(l.stdout, "%s\n", l.r.RenderCards(trello.Cards{card}))
}

func (l ls) renderComments(card trello.Card) {
	comments, err := l.tr.FindComments(card.ID)
	if err != nil {
		fmt.Fprintf(l.stderr, "could not fetch comments for card '%s': %v\n", card.Name, err)
	} else {
		fmt.Fprintf(l.stdout, "%s\n", l.r.RenderComments(comments))
	}
}

func (l ls) renderComment(comment trello.Comment) {
	fmt.Fprintf(l.stdout, "%s\n", l.r.RenderComment(comment))
}
