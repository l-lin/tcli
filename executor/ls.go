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
		then().
		doOnEmptyBoardName(func() {
			l.renderBoards()
		}).
		findBoard().
		doOnBoard(func(board *trello.Board) {
			l.renderLists(*board)
		}).
		then().
		findList().
		doOnList(func(list *trello.List) {
			l.renderCards(*list)
		}).
		then().
		findCard().
		doOnCard(func(card *trello.Card) {
			l.renderComments(*card)
		}).
		then().
		findComment().
		doOnComment(func(comment *trello.Comment) {
			l.renderComment(*comment)
		}).err; err != nil {
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

func (l ls) renderComments(card trello.Card) {
	comments, err := l.tr.FindComments(card.ID)
	if err != nil {
		fmt.Fprintf(l.stderr, "could not fetch comments for card '%s': %v\n", card.Name, err)
	} else {
		mapReactionSummariesByIDComment := make(map[string]trello.ReactionSummaries)
		for _, comment := range comments {
			var reactionSummaries trello.ReactionSummaries
			reactionSummaries, err = l.tr.FindReactionSummaries(comment.ID)
			if err != nil {
				fmt.Fprintf(l.stderr, "could not fetch reaction summaries for comment '%s': %v\n", comment.ID, err)
				return
			}
			mapReactionSummariesByIDComment[comment.ID] = reactionSummaries
		}
		fmt.Fprintf(l.stdout, "%s\n", l.r.RenderComments(comments, mapReactionSummariesByIDComment))
	}
}

func (l ls) renderComment(comment trello.Comment) {
	reactionSummaries, err := l.tr.FindReactionSummaries(comment.ID)
	if err != nil {
		fmt.Fprintf(l.stderr, "could not fetch reaction summaries for comment '%s': %v\n", comment.ID, err)
	} else {
		fmt.Fprintf(l.stdout, "%s\n", l.r.RenderComment(comment, reactionSummaries))
	}
}
