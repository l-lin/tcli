package executor

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestTouch_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type given struct {
		args                  []string
		buildTrelloRepository func() trello.Repository
		session               *trello.Session
	}
	type expected struct {
		stdout string
		stderr string
	}

	board := trello.Board{ID: "board 1", Name: "board"}
	list := trello.List{ID: "list 1", Name: "list"}
	card := trello.Card{ID: "card 1", Name: "card"}
	createCard1 := trello.CreateCard{
		Name:   "card",
		IDList: list.ID,
	}
	createCard2 := trello.CreateCard{
		Name:   "another-card",
		IDList: list.ID,
	}
	createComment := trello.CreateComment{
		IDCard: card.ID,
		Text:   "comment",
	}

	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"/> touch": {
			given: given{
				args: []string{},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "missing card operand\n",
			},
		},
		"/> touch ": {
			given: given{
				args: []string{""},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "missing card operand\n",
			},
		},
		"/> touch /board/list/card": {
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
						CreateCard(createCard1).
						Return(nil, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{},
		},
		"/board/list> touch card": {
			given: given{
				args: []string{"card"},
				session: &trello.Session{
					Board: &board,
					List:  &list,
				},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						CreateCard(createCard1).
						Return(nil, nil)
					return tr
				},
			},
			expected: expected{},
		},
		"/> touch /board/list/card /board/list/another-card": {
			given: given{
				args: []string{"/board/list/card", "/board/list/another-card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil).
						Times(2)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil).
						Times(2)
					tr.EXPECT().
						CreateCard(createCard1).
						Return(nil, nil)
					tr.EXPECT().
						CreateCard(createCard2).
						Return(nil, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{},
		},
		"/> touch /board/list/card/comment": {
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
						CreateComment(createComment).
						Return(nil, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{},
		},
		// ERRORS
		"/> touch / (no board name)": {
			given: given{
				args: []string{"/"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "nothing to create\n",
			},
		},
		"/> touch /board (no list name)": {
			given: given{
				args: []string{"/board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "board creation not implemented yet\n",
			},
		},
		"/> touch /board/list/card (board not found)": {
			given: given{
				args: []string{"/board/list/card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(nil, errors.New("not found"))
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "no board found with name 'board'\n",
			},
		},
		"/> touch /board/list (no card name)": {
			given: given{
				args: []string{"/board/list"},
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
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "list creation not implemented yet\n",
			},
		},
		"/> touch /board/list/card (list not found)": {
			given: given{
				args: []string{"/board/list/card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(nil, errors.New("not found"))
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "no list found with name 'list'\n",
			},
		},
		"/> touch /../.. (invalid path)": {
			given: given{
				args: []string{"/../.."},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"/> touch /board/list/card (error when creating card)": {
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
						CreateCard(createCard1).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: fmt.Sprintf("could not create card '%s': unexpected error\n", createCard1.Name),
			},
		},
		"/> touch /board/list/card/comment (error when creating comment)": {
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
						CreateComment(createComment).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: fmt.Sprintf("could not create comment '%s': unexpected error\n", createComment.Text),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}
			to := touch{
				executor{
					tr:      tt.given.buildTrelloRepository(),
					session: tt.given.session,
					stdout:  &stdoutBuf,
					stderr:  &stderrBuf,
				}}
			to.Execute(tt.given.args)

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
