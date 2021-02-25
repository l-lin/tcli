package executor

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestMv_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	board := trello.Board{ID: "board 1", Name: "board"}
	list1 := trello.List{ID: "list 1", Name: "list"}
	list2 := trello.List{ID: "list 2", Name: "another-list"}
	card := trello.Card{ID: "card 1", Name: "card", IDList: list1.ID}
	type given struct {
		args                  []string
		buildTrelloRepository func() trello.Repository
	}
	type expected struct {
		stdout string
		stderr string
	}
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"/> mv /board/list/card /board/another-list": {
			given: given{
				args: []string{"/board/list/card", "/board/another-list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil).
						Times(2)
					tr.EXPECT().
						FindList(board.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindList(board.ID, list2.Name).
						Return(&list2, nil)
					tr.EXPECT().
						FindCard(list1.ID, card.Name).
						Return(&card, nil)
					updatedCard := trello.NewUpdateCard(card)
					updatedCard.IDList = list2.ID
					tr.EXPECT().
						UpdateCard(updatedCard).
						Return(nil, nil)
					return tr
				},
			},
			expected: expected{},
		},
		"/> mv /board/list/card /board/list/new-card-name": {
			given: given{
				args: []string{"/board/list/card", "/board/list/new-card-name"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil).
						Times(2)
					tr.EXPECT().
						FindList(board.ID, list1.Name).
						Return(&list1, nil).
						Times(2)
					tr.EXPECT().
						FindCard(list1.ID, card.Name).
						Return(&card, nil)
					updatedCard := trello.NewUpdateCard(card)
					updatedCard.Name = "new-card-name"
					tr.EXPECT().
						UpdateCard(updatedCard).
						Return(nil, nil)
					return tr
				},
			},
			expected: expected{},
		},
		// ERRORS
		"mv": {
			given: given{
				args: []string{},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "missing card source operand\n",
			},
		},
		"mv /board/list/card": {
			given: given{
				args: []string{"/board/list/card"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "missing list destination operand\n",
			},
		},
		"mv 1 2 3": {
			given: given{
				args: []string{"1", "2", "3"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "command only accepts two arguments\n",
			},
		},
		"mv /unknown-board/list/card /board/another-lise (board not found)": {
			given: given{
				args: []string{"/unknown-board/list/card", "/board/another-list"},
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
		"mv /board/unknown-list/card /board/another-list (list not found)": {
			given: given{
				args: []string{"/board/unknown-list/card", "/board/another-list"},
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
		"mv /board/list/unknown-card /board/another-lise (card not found)": {
			given: given{
				args: []string{"/board/list/unknown-card", "/board/another-list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, "unknown-card").
						Return(nil, errors.New("not found"))
					return tr
				},
			},
			expected: expected{
				stderr: "no card found with name 'unknown-card'\n",
			},
		},
		"mv /../.. /foo (invalid path at first arg)": {
			given: given{
				args: []string{"/../..", "/foo"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"mv /board/list/card /../.. (invalid path at second arg)": {
			given: given{
				args: []string{"/board/list/card", "/../.."},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, card.Name).
						Return(&card, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"mv / /board/another-list (empty board name at 1st argument)": {
			given: given{
				args: []string{"/", "/board/another-list"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"mv /board/list/card / (empty board name at 2nd argument)": {
			given: given{
				args: []string{"/board/list/card", "/"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, card.Name).
						Return(&card, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"mv /board /board/another-list (empty list name at 1st argument)": {
			given: given{
				args: []string{"/board", "/board/another-list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"mv /board/list/card /board (empty list name at 2nd argument)": {
			given: given{
				args: []string{"/board/list/card", "/board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil).
						Times(2)
					tr.EXPECT().
						FindList(board.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, card.Name).
						Return(&card, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"mv /board/list /board/another-list (empty card name)": {
			given: given{
				args: []string{"/board/list/", "/board/another-list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list1.Name).
						Return(&list1, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"mv /board/list/card /board/another-list (error when updating card)": {
			given: given{
				args: []string{"/board/list/card", "/board/another-list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil).
						Times(2)
					tr.EXPECT().
						FindList(board.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindList(board.ID, list2.Name).
						Return(&list2, nil)
					tr.EXPECT().
						FindCard(list1.ID, card.Name).
						Return(&card, nil)
					updatedCard := trello.NewUpdateCard(card)
					updatedCard.IDList = list2.ID
					tr.EXPECT().
						UpdateCard(updatedCard).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
			},
			expected: expected{
				stderr: "could not update card: unexpected error\n",
			},
		},
		"mv /board/list/card/comment /board/another-list": {
			given: given{
				args: []string{"/board/list/card/comment", "/board/another-list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, card.Name).
						Return(&card, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "cannot move comments\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}

			r := mv{
				executor: executor{
					tr:      tt.given.buildTrelloRepository(),
					session: &trello.Session{},
					stdout:  &stdoutBuf,
					stderr:  &stderrBuf,
				},
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
