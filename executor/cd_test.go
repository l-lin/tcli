package executor

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/trello"
	"reflect"
	"testing"
)

func TestCd_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type given struct {
		args                  []string
		buildTrelloRepository func() trello.Repository
		session               *trello.Session
	}
	type expected struct {
		stderr string
		board  *trello.Board
		list   *trello.List
		card   *trello.Card
	}

	board1 := trello.Board{ID: "board 1", Name: "board"}
	list1 := trello.List{ID: "list 1", Name: "list"}
	list2 := trello.List{ID: "list 2", Name: "another-list"}
	card1 := trello.Card{ID: "card 1", Name: "card"}

	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"/ > cd board": {
			given: given{
				args: []string{board1.Name},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{
				board: &board1,
			},
		},
		"/ > cd board/list": {
			given: given{
				args: []string{"board/list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{
				board: &board1,
				list:  &list1,
			},
		},
		"/board > cd": {
			given: given{
				args: []string{},
				session: &trello.Session{
					Board: &board1,
				},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				board: nil,
				list:  nil,
			},
		},
		"/board > cd ": {
			given: given{
				args: []string{""},
				session: &trello.Session{
					Board: &board1,
				},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				board: nil,
				list:  nil,
			},
		},
		"/board > cd list": {
			given: given{
				args: []string{list1.Name},
				session: &trello.Session{
					Board: &board1,
				},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					return tr
				},
			},
			expected: expected{
				board: &board1,
				list:  &list1,
			},
		},
		"/board > cd ..": {
			given: given{
				args: []string{".."},
				session: &trello.Session{
					Board: &board1,
				},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				board: nil,
				list:  nil,
			},
		},
		"/board/list > cd ../another-list": {
			given: given{
				args: []string{"../" + list2.Name},
				session: &trello.Session{
					Board: &board1,
					List:  &list1,
				},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list2.Name).
						Return(&list2, nil)
					return tr
				},
			},
			expected: expected{
				board: &board1,
				list:  &list2,
			},
		},
		"/ > cd board/list/card": {
			given: given{
				args: []string{"board/list/card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{
				board: &board1,
				list:  &list1,
				card:  &card1,
			},
		},
		// ERRORS
		"invalid path": {
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
		"/ > cd unknown-board": {
			given: given{
				args: []string{"unknown-board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("unknown-board").
						Return(nil, errors.New("not found"))
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "no board found with name 'unknown-board'\n",
			},
		},
		"/ > cd /board/unknown-list": {
			given: given{
				args: []string{"/board/unknown-list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "unknown-list").
						Return(nil, errors.New("not found"))
					return tr
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "no list found with name 'unknown-list'\n",
			},
		},
		"/ > cd ..": {
			given: given{
				args: []string{".."},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: expected{
				stderr: "invalid path\n",
				board:  nil,
				list:   nil,
			},
		},
		"/board/list > cd ../../..": {
			given: given{
				args: []string{"../../.."},
				session: &trello.Session{
					Board: &board1,
					List:  &list1,
				},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "invalid path\n",
				board:  &board1,
				list:   &list1,
			},
		},
		"/board/list > cd board2 board3": {
			given: given{
				args: []string{"board", "board2"},
				session: &trello.Session{
					Board: &board1,
					List:  &list1,
				},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				stderr: "only one argument is accepted\n",
				board:  &board1,
				list:   &list1,
			},
		},
		"/board/list/card > cd comment": {
			given: given{
				args: []string{"comment"},
				session: &trello.Session{
					Board: &board1,
					List:  &list1,
					Card:  &card1,
				},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					return tr
				},
			},
			expected: expected{
				stderr: "cannot cd on comment\n",
				board:  &board1,
				list:   &list1,
				card:   &card1,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			stderrBuf := bytes.Buffer{}
			c := cd{
				executor{
					tr:      tt.given.buildTrelloRepository(),
					session: tt.given.session,
					stderr:  &stderrBuf,
				},
			}
			c.Execute(tt.given.args)
			actualBoard := tt.given.session.Board
			actualList := tt.given.session.List
			actualCard := tt.given.session.Card
			if tt.expected.board != nil && actualBoard == nil || tt.expected.board == nil && actualBoard != nil {
				t.Errorf("expected board %v, actual board %v", tt.expected.board, actualBoard)
				t.FailNow()
			}
			if tt.expected.board != nil {
				if *actualBoard != *tt.expected.board {
					t.Errorf("expected board %v, actual board %v", tt.expected.board, actualBoard)
				}
			}
			if tt.expected.list != nil && actualList == nil || tt.expected.list == nil && actualList != nil {
				t.Errorf("expected list %v, actual list %v", tt.expected.list, actualList)
				t.FailNow()
			}
			if tt.expected.list != nil {
				if *actualList != *tt.expected.list {
					t.Errorf("expected list %v, actual list %v", tt.expected.list, actualList)
				}
			}
			if tt.expected.card != nil && actualCard == nil || tt.expected.card == nil && actualCard != nil {
				t.Errorf("expected card %v, actual card %v", tt.expected.card, actualCard)
				t.FailNow()
			}
			if tt.expected.card != nil {
				if !reflect.DeepEqual(actualCard, tt.expected.card) {
					t.Errorf("expected card %v, actual card %v", tt.expected.card, actualCard)
				}
			}
			actualStderr := stderrBuf.String()
			if actualStderr != tt.expected.stderr {
				t.Errorf("expected stderr %v, actual stderr %v", tt.expected.stderr, actualStderr)
			}
		})
	}
}
