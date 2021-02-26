package executor

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestCp_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	board := trello.Board{ID: "board 1", Name: "board"}
	list1 := trello.List{ID: "list 1", Name: "list"}
	list2 := trello.List{ID: "list 2", Name: "another-list"}
	card1 := trello.Card{ID: "card 1", Name: "card", IDList: list1.ID}
	card2 := trello.Card{ID: "card 2", Name: "another-card", IDList: list1.ID}
	comment := trello.Comment{
		ID: "comment",
		Data: trello.CommentData{
			Card: trello.CommentDataCard{
				ID: card1.ID,
			},
			Text: "comment content",
		},
	}
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
		// CARDS
		"/ > cp /board/list/card /board/another-list": {
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
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					tr.EXPECT().
						FindList(board.ID, list2.Name).
						Return(&list2, nil)
					createCard := trello.NewCreateCard(card1)
					createCard.IDList = list2.ID
					tr.EXPECT().
						CreateCard(createCard).
						Return(nil, nil)
					return tr
				},
			},
			expected: expected{},
		},
		"/ > cp /board/list/card /board/list/another-card": {
			given: given{
				args: []string{"/board/list/card", "/board/list/another-card"},
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
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					tr.EXPECT().
						FindCard(list1.ID, card2.Name).
						Return(nil, cardNotFoundError(card2.ID))
					createCard := trello.NewCreateCard(card1)
					createCard.Name = card2.Name
					tr.EXPECT().
						CreateCard(createCard).
						Return(nil, nil)
					return tr
				},
			},
			expected: expected{},
		},
		// COMMENTS
		"/ > cp /board/list/card/comment /board/comment/another-card": {
			given: given{
				args: []string{"/board/list/card/comment", "/board/list/another-card"},
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
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					tr.EXPECT().
						FindComment(card1.ID, comment.ID).
						Return(&comment, nil)
					tr.EXPECT().
						FindCard(list1.ID, card2.Name).
						Return(&card2, nil)
					tr.EXPECT().
						CreateComment(trello.CreateComment{
							IDCard: card2.ID,
							Text:   comment.Data.Text,
						}).
						Return(nil, nil)
					return tr
				},
			},
			expected: expected{},
		},
		// ERRORS
		"/ > cp": {
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
		"/ > cp /board/list/card": {
			given: given{
				args: []string{"/board/list/card"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "missing destination operand\n",
			},
		},
		"/ > cp 1 2 3": {
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
		"/ > cp /unknown-board/list/card /board/another-list (board not found)": {
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
		"/ > cp /board/unknown-list/card /board/another-list (list not found)": {
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
		"/ > cp /board/list/unknown-card /board/another-list (card not found)": {
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
		"/ > cp /../.. /foo (invalid path at first arg)": {
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
		"/ > cp /board/list/card /../.. (invalid path at second arg)": {
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
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"/ > cp / /board/another-list": {
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
		"/ > cp /board/list/card /": {
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
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"/ > cp /board/ /board/another-list": {
			given: given{
				args: []string{"/board/", "/board/another-list"},
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
		"/ > cp /board/list/card /board": {
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
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"/ > cp /board/list/ /board/another-list": {
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
		"/ > cp /board/list/card /board/another-list (error when creating card)": {
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
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					createCard := trello.NewCreateCard(card1)
					createCard.IDList = list2.ID
					tr.EXPECT().
						CreateCard(createCard).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
			},
			expected: expected{
				stderr: "could not copy card 'card': unexpected error\n",
			},
		},
		"/ > cp /board/list/card/comment /board/comment": {
			given: given{
				args: []string{"/board/list/card/comment", "/board/list"},
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
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					tr.EXPECT().
						FindComment(card1.ID, comment.ID).
						Return(&comment, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "cannot copy comment in list\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}

			r := cp{
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
