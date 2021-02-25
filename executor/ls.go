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
	if err := start(l.tr).
		resolvePath(l.session, arg).
		doOnEmptyBoardName(func() {
			l.renderBoards()
		}).
		thenFindBoard().
		doOnEmptyListName(func(session *trello.Session) {
			l.renderLists(*session.Board)
		}).
		thenFindList().
		doOnEmptyCardName(func(session *trello.Session) {
			l.renderCards(*session.List)
		}).
		thenFindCard().
		doOnEmptyCommentID(func(session *trello.Session) {
			l.renderComments(*session.Card)
		}).
		thenFindComment().
		andDoOnComment(func(comment *trello.Comment) {
			l.renderComment(*comment)
		}); err != nil {
		fmt.Fprintf(l.stderr, "%s\n", err)
	}
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
