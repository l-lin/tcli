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
		arg                   string
		buildTrelloRepository func() trello.Repository
		currentBoard          *trello.Board
		currentList           *trello.List
	}
	type expected struct {
		stdout string
		stderr string
	}

	board := trello.Board{ID: "board 1", Name: "board"}
	list := trello.List{ID: "list 1", Name: "list"}
	createCard := trello.CreateCard{
		Name:   "card",
		IDList: list.ID,
	}

	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"no arg": {
			given: given{
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "missing card operand",
			},
		},
		"/: create card from absolute path": {
			given: given{
				arg: "/board/list/card",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						CreateCard(createCard).
						Return(nil, nil)
					return tr
				},
			},
			expected: expected{},
		},
		"/board/list: create card from relative path": {
			given: given{
				arg:          "card",
				currentBoard: &board,
				currentList:  &list,
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						CreateCard(createCard).
						Return(nil, nil)
					return tr
				},
			},
			expected: expected{},
		},
		// ERRORS
		"no board name": {
			given: given{
				arg: "/",
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "nothing to create\n",
			},
		},
		"no list name": {
			given: given{
				arg: "/board",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "board creation not implemented yet\n",
			},
		},
		"board not found": {
			given: given{
				arg: "/board/list/card",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(nil, errors.New("not found"))
					return tr
				},
			},
			expected: expected{
				stderr: "no board found with name 'board'\n",
			},
		},
		"no card name": {
			given: given{
				arg: "/board/list",
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
				stderr: "list creation not implemented yet\n",
			},
		},
		"list not found": {
			given: given{
				arg: "/board/list/card",
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
			},
			expected: expected{
				stderr: "no list found with name 'list'\n",
			},
		},
		"invalid path": {
			given: given{
				arg: "/../..",
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"error when creating card": {
			given: given{
				arg: "/board/list/card",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						CreateCard(createCard).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
			},
			expected: expected{
				stderr: fmt.Sprintf("could not create card '%s': unexpected error\n", createCard.Name),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}
			to := touch{
				executor{
					tr:           tt.given.buildTrelloRepository(),
					currentBoard: tt.given.currentBoard,
					currentList:  tt.given.currentList,
					stdout:       &stdoutBuf,
					stderr:       &stderrBuf,
				}}
			to.Execute(tt.given.arg)

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
