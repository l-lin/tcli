package executor

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/trello"
	"io"
	"testing"
)

func TestRm_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	board := trello.Board{ID: "board 1", Name: "board"}
	list := trello.List{ID: "list 1", Name: "list"}
	card := trello.Card{ID: "card 1", Name: "card"}
	updatedCard := trello.NewUpdateCard(card)
	updatedCard.Closed = true
	card2 := trello.Card{ID: "card 2", Name: "another card"}
	updatedCard2 := trello.NewUpdateCard(card2)
	updatedCard2.Closed = true
	comment := trello.Comment{ID: "comment", Data: trello.CommentData{Card: trello.CommentDataCard{ID: card.ID}}}

	type given struct {
		args                  []string
		buildTrelloRepository func() trello.Repository
		stdin                 io.ReadCloser
	}
	type expected struct {
		stdout string
		stderr string
	}
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"rm": {
			given: given{
				args:                  []string{},
				buildTrelloRepository: func() trello.Repository { return nil },
			},
			expected: expected{
				stderr: "missing card operand\n",
			},
		},
		"rm ": {
			given: given{
				args:                  []string{""},
				buildTrelloRepository: func() trello.Repository { return nil },
			},
			expected: expected{
				stderr: "missing card operand\n",
			},
		},
		"rm /board/list/card (user accepts to archive)": {
			given: given{
				args: []string{"/board/list/card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						FindCard(list.ID, card.Name).
						Return(&card, nil)
					tr.EXPECT().
						UpdateCard(updatedCard).
						Return(nil, nil)
					return tr
				},
				stdin: acceptStdin(),
			},
			expected: expected{},
		},
		"rm /board/list/card (user refuses to archive)": {
			given: given{
				args: []string{"/board/list/card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						FindCard(list.ID, card.Name).
						Return(&card, nil)
					return tr
				},
				stdin: refuseStdin(),
			},
			expected: expected{},
		},
		"rm /board/list/card/comment (user accepts to delete)": {
			given: given{
				args: []string{"/board/list/card/comment"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						FindCard(list.ID, card.Name).
						Return(&card, nil)
					tr.EXPECT().
						FindComment(card.ID, comment.ID).
						Return(&comment, nil)
					tr.EXPECT().
						DeleteComment(card.ID, comment.ID).
						Return(nil)
					return tr
				},
				stdin: acceptStdin(),
			},
			expected: expected{},
		},
		"rm /board/list/card/comment (user refuses to delete)": {
			given: given{
				args: []string{"/board/list/card/comment"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						FindCard(list.ID, card.Name).
						Return(&card, nil)
					tr.EXPECT().
						FindComment(card.ID, comment.ID).
						Return(&comment, nil)
					return tr
				},
				stdin: refuseStdin(),
			},
			expected: expected{},
		},
		"rm /board/list/*": {
			given: given{
				args: []string{"/board/list/*"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						FindCards(list.ID).
						Return(trello.Cards{card, card2}, nil)
					tr.EXPECT().
						ArchiveAllCards(list.ID).
						Return(nil)
					return tr
				},
				stdin: acceptStdin(),
			},
		},
		// ERRORS
		"rm /../..": {
			given: given{
				args: []string{"/../.."},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"rm /": {
			given: given{
				args: []string{"/"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "nothing to archive\n",
			},
		},
		"rm /unknown-board": {
			given: given{
				args: []string{"/unknown-board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("unknown-board").
						Return(nil, errors.New("not found"))
					return tr
				},
			},
			expected: expected{
				stderr: "no board found with name 'unknown-board'\n",
			},
		},
		"rm /board": {
			given: given{
				args: []string{"/board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "board archiving not implemented yet\n",
			},
		},
		"rm /board/unknown-list": {
			given: given{
				args: []string{"/board/unknown-list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, "unknown-list").
						Return(nil, errors.New("not found"))
					return tr
				},
			},
			expected: expected{
				stderr: "no list found with name 'unknown-list'\n",
			},
		},
		"rm /board/list/": {
			given: given{
				args: []string{"/board/list/"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "list archiving not implemented yet\n",
			},
		},
		"rm /board/list/unknown-card": {
			given: given{
				args: []string{"/board/list/unknown-card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						FindCard(list.ID, "unknown-card").
						Return(nil, errors.New("not found"))
					return tr
				},
			},
			expected: expected{
				stderr: "no card found with name 'unknown-card'\n",
			},
		},
		"rm /board/list/card (error when archiving card)": {
			given: given{
				args: []string{"/board/list/card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						FindCard(list.ID, card.Name).
						Return(&card, nil)
					tr.EXPECT().
						UpdateCard(updatedCard).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				stdin: acceptStdin(),
			},
			expected: expected{
				stderr: fmt.Sprintf("could not archive card '%s': unexpected error\n", updatedCard.Name),
			},
		},
		"rm /board/list/card/comment (error when deleting comment)": {
			given: given{
				args: []string{"/board/list/card/comment"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						FindCard(list.ID, card.Name).
						Return(&card, nil)
					tr.EXPECT().
						FindComment(card.ID, comment.ID).
						Return(&comment, nil)
					tr.EXPECT().
						DeleteComment(card.ID, comment.ID).
						Return(errors.New("unexpected error"))
					return tr
				},
				stdin: acceptStdin(),
			},
			expected: expected{
				stderr: fmt.Sprintf("could not delete comment '%s': unexpected error\n", comment.ID),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}

			r := rm{
				executor: executor{
					tr:      tt.given.buildTrelloRepository(),
					session: &trello.Session{},
					stdout:  &stdoutBuf,
					stderr:  &stderrBuf,
				},
				stdin: tt.given.stdin,
			}
			r.Execute(tt.given.args)

			actualStdout := stdoutBuf.String()
			if actualStdout != tt.expected.stdout {
				t.Errorf("expected stdout %v, actual stdout %v", tt.expected.stdout, actualStdout)
			}
			actualStderr := stderrBuf.String()
			if actualStderr != tt.expected.stderr {
				t.Errorf("expected stderr %v, actual stderr %v", tt.expected.stderr, actualStderr)
			}
		})
	}
}
