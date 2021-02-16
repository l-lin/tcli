package executor

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestCd_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type given struct {
		arg                   string
		buildTrelloRepository func() trello.Repository
		currentBoard          *trello.Board
		currentList           *trello.List
	}
	type expected struct {
		errOutput string
		board     *trello.Board
		list      *trello.List
	}

	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"/ > cd board": {
			given: given{
				arg: "board",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					return tr
				},
			},
			expected: expected{
				board: &trello.Board{ID: "board 1", Name: "board"},
			},
		},
		"/ > cd board/list": {
			given: given{
				arg: "board/list",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					return tr
				},
			},
			expected: expected{
				board: &trello.Board{ID: "board 1", Name: "board"},
				list:  &trello.List{ID: "list 1", Name: "list"},
			},
		},
		"/board > cd ": {
			given: given{
				arg:          "",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
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
				arg:          "list",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					return tr
				},
			},
			expected: expected{
				board: &trello.Board{ID: "board 1", Name: "board"},
				list:  &trello.List{ID: "list 1", Name: "list"},
			},
		},
		"/board > cd ..": {
			given: given{
				arg:          "..",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
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
				arg:          "../another-list",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
				currentList:  &trello.List{ID: "list 1", Name: "list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "another-list").
						Return(&trello.List{ID: "list 2", Name: "another-list"}, nil)
					return tr
				},
			},
			expected: expected{
				board: &trello.Board{ID: "board 1", Name: "board"},
				list:  &trello.List{ID: "list 2", Name: "another-list"},
			},
		},
		// ERRORS
		"/ > cd board/list/card": {
			given: given{
				arg: "board/list/card",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					return tr
				},
			},
			expected: expected{
				errOutput: "cannot cd on card\n",
			},
		},
		"/ > cd ..": {
			given: given{
				arg: "..",
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				errOutput: "invalid path\n",
				board:     nil,
				list:      nil,
			},
		},
		"/board/list > cd ../../..": {
			given: given{
				arg:          "../../..",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
				currentList:  &trello.List{ID: "list 1", Name: "list"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: expected{
				errOutput: "invalid path\n",
				board:     &trello.Board{ID: "board 1", Name: "board"},
				list:      &trello.List{ID: "list 1", Name: "list"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			errBuf := bytes.Buffer{}
			c := cd{
				executor{
					tr:           tt.given.buildTrelloRepository(),
					currentBoard: tt.given.currentBoard,
					currentList:  tt.given.currentList,
					errOutput:    &errBuf,
				},
			}
			actualBoard, actualList := c.Execute(tt.given.arg)
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
			actualErrOutput := errBuf.String()
			if actualErrOutput != tt.expected.errOutput {
				t.Errorf("expected errOutput %v, actual errOutput %v", tt.expected.errOutput, actualErrOutput)
			}
		})
	}
}
