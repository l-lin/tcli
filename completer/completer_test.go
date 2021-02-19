package completer

import (
	"errors"
	"github.com/c-bata/go-prompt"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/executor"
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestCompleter_Complete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type given struct {
		currentBoard          *trello.Board
		currentList           *trello.List
		cmd, arg              string
		buildTrelloRepository func() trello.Repository
	}
	var tests = map[string]struct {
		given    given
		expected []prompt.Suggest
	}{
		"empty text": {
			given: given{
				buildTrelloRepository: func() trello.Repository { return nil },
			},
			expected: commandSuggestions(),
		},
		"typing 'unknown'": {
			given: given{
				cmd: "unknown",
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: []prompt.Suggest{},
		},
		// COMMANDS
		"typing 'c'": {
			given: given{
				cmd: "c",
				arg: "",
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: []prompt.Suggest{
				{Text: "cd", Description: "change level in the hierarchy"},
				{Text: "cat", Description: "show resource content info"},
			},
		},
		// RELATIVE PATHS
		"typing 'cd '": {
			given: given{
				cmd: "cd",
				arg: "",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(trello.Boards{
							{Name: "board", ShortLink: "shortLink"},
							{Name: "another board", ShortLink: "another shortLink"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "board [shortLink]"},
				{Text: "another board [another shortLink]"},
			},
		},
		"typing 'cd b'": {
			given: given{
				cmd: "cd",
				arg: "b",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("b").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(trello.Boards{
							{Name: "board", ShortLink: "shortLink"},
							{Name: "another board", ShortLink: "another shortLink"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "board [shortLink]"},
			},
		},
		"typing 'cd board/'": {
			given: given{
				cmd: "cd",
				arg: "board/",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{Name: "board", ID: "board 1"}, nil)
					tr.EXPECT().
						FindList("board 1", "").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists("board 1").
						Return(trello.Lists{
							{Name: "list"},
							{Name: "another list"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "list"},
				{Text: "another list"},
			},
		},
		"typing 'cd board/l'": {
			given: given{
				cmd: "cd",
				arg: "board/l",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{Name: "board", ID: "board 1"}, nil)
					tr.EXPECT().
						FindList("board 1", "l").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists("board 1").
						Return(trello.Lists{
							{Name: "list"},
							{Name: "another list"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "list"},
			},
		},
		"typing 'cd board/list/'": {
			given: given{
				cmd: "cd",
				arg: "board/list/",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{Name: "board", ID: "board 1"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{Name: "list", ID: "list 1"}, nil)
					tr.EXPECT().
						FindCards("list 1").
						Return(trello.Cards{
							{Name: "card", ShortLink: "shortLink"},
							{Name: "another card", ShortLink: "another shortLink"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "card [shortLink]"},
				{Text: "another card [another shortLink]"},
			},
		},
		"typing 'cd board/list/ca'": {
			given: given{
				cmd: "cd",
				arg: "board/list/ca",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{Name: "board", ID: "board 1"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{Name: "list", ID: "list 1"}, nil)
					tr.EXPECT().
						FindCards("list 1").
						Return(trello.Cards{
							{Name: "card", ShortLink: "shortLink"},
							{Name: "another card", ShortLink: "another shortLink"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "card [shortLink]"},
			},
		},
		// WITH CURRENT BOARD
		"has current board & typing 'cd '": {
			given: given{
				cmd:          "cd",
				arg:          "",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{Name: "board", ID: "board 1"}, nil)
					tr.EXPECT().
						FindList("board 1", "").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists("board 1").
						Return(trello.Lists{
							{Name: "list"},
							{Name: "another list"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "list"},
				{Text: "another list"},
			},
		},
		"has current board & typing 'cd l'": {
			given: given{
				cmd:          "cd",
				arg:          "l",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{Name: "board", ID: "board 1"}, nil)
					tr.EXPECT().
						FindList("board 1", "l").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists("board 1").
						Return(trello.Lists{
							{Name: "list"},
							{Name: "another list"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "list"},
			},
		},
		"has current board & current list & typing 'cd '": {
			given: given{
				cmd:          "cd",
				arg:          "",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
				currentList:  &trello.List{ID: "list 1", Name: "list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{Name: "board", ID: "board 1"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					tr.EXPECT().
						FindCards("list 1").
						Return(trello.Cards{
							{Name: "card", ShortLink: "shortLink"},
							{Name: "another card", ShortLink: "another shortLink"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "card [shortLink]"},
				{Text: "another card [another shortLink]"},
			},
		},
		"has current board & current list & typing 'cd c'": {
			given: given{
				cmd:          "cd",
				arg:          "c",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
				currentList:  &trello.List{ID: "list 1", Name: "list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{Name: "board", ID: "board 1"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					tr.EXPECT().
						FindCards("list 1").
						Return(trello.Cards{
							{Name: "card", ShortLink: "shortLink"},
							{Name: "another card", ShortLink: "another short link"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "card [shortLink]"},
			},
		},
		"has current board & typing 'cd ../a'": {
			given: given{
				cmd:          "cd",
				arg:          "../a",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("a").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(trello.Boards{
							{Name: "board", ShortLink: "shortLink"},
							{Name: "another board", ShortLink: "another shortLink"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "another board [another shortLink]"},
			},
		},
		"has current board & current list & typing 'cd ../a'": {
			given: given{
				cmd:          "cd",
				arg:          "../a",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
				currentList:  &trello.List{ID: "list 1", Name: "list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{Name: "board", ID: "board 1"}, nil)
					tr.EXPECT().
						FindList("board 1", "a").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists("board 1").
						Return(trello.Lists{
							{Name: "list"},
							{Name: "another list"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "another list"},
			},
		},
		// ABSOLUTE PATHS
		"typing 'cd /'": {
			given: given{
				cmd: "cd",
				arg: "/",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(trello.Boards{
							{Name: "board", ShortLink: "shortLink"},
							{Name: "another board", ShortLink: "another shortLink"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "board [shortLink]"},
				{Text: "another board [another shortLink]"},
			},
		},
		"has current board & typing 'cd /a'": {
			given: given{
				cmd:          "cd",
				arg:          "/a",
				currentBoard: &trello.Board{ID: "board 1", Name: "board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("a").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(trello.Boards{
							{Name: "board", ShortLink: "shortLink"},
							{Name: "another board", ShortLink: "another shortLink"},
						}, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "another board [another shortLink]"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := New(tt.given.buildTrelloRepository(), tt.given.currentBoard, tt.given.currentList)
			actual := c.Complete(tt.given.cmd, tt.given.arg)
			if len(actual) != len(tt.expected) {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
				t.FailNow()
			}
			for i := 0; i < len(actual); i++ {
				if actual[i] != tt.expected[i] {
					t.Errorf("%d: expected %v, actual %v", i, tt.expected[i], actual[i])
				}
			}
		})
	}
}

func commandSuggestions() []prompt.Suggest {
	suggestions := make([]prompt.Suggest, len(executor.Factories))
	for i, factory := range executor.Factories {
		suggestions[i] = prompt.Suggest{
			Text:        factory.Cmd,
			Description: factory.Description,
		}
	}
	return suggestions
}
