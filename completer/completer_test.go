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
		cmd                   string
		args                  []string
		buildTrelloRepository func() trello.Repository
	}
	board1 := trello.Board{Name: "board", ID: "board 1", ShortLink: "shortLink"}
	board2 := trello.Board{Name: "another board", ID: "board 2", ShortLink: "another shortLink"}
	boards := trello.Boards{board1, board2}

	list1 := trello.List{Name: "list", ID: "list 1"}
	list2 := trello.List{Name: "another list", ID: "list 2"}
	lists := trello.Lists{list1, list2}

	cards := trello.Cards{
		{Name: "card", ShortLink: "shortLink"},
		{Name: "another card", ShortLink: "another shortLink"},
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
				cmd:  "c",
				args: []string{""},
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
				cmd:  "cd",
				args: []string{""},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "board[shortLink]"},
				{Text: `another\ board[another shortLink]`},
			},
		},
		"typing 'cd b'": {
			given: given{
				cmd:  "cd",
				args: []string{"b"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("b").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "board[shortLink]"},
			},
		},
		"typing 'cd board/'": {
			given: given{
				cmd:  "cd",
				args: []string{"board/"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: list1.SanitizedName()},
				{Text: `another\ list`},
			},
		},
		"typing 'cd board/l'": {
			given: given{
				cmd:  "cd",
				args: []string{"board/l"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "l").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: list1.SanitizedName()},
			},
		},
		"typing 'cd board/list/c'": {
			given: given{
				cmd:  "cd",
				args: []string{"board/list/c"},
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
			expected: []prompt.Suggest{},
		},
		"typing 'cat board/list/'": {
			given: given{
				cmd:  "cat",
				args: []string{"board/list/"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCards(list1.ID).
						Return(cards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "card[shortLink]"},
				{Text: `another\ card[another shortLink]`},
			},
		},
		"typing 'cat board/list/ca'": {
			given: given{
				cmd:  "cat",
				args: []string{"board/list/ca"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCards(list1.ID).
						Return(cards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "card[shortLink]"},
			},
		},
		// WITH CURRENT BOARD
		"has current board & typing 'cd '": {
			given: given{
				cmd:          "cd",
				args:         []string{""},
				currentBoard: &board1,
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: list1.SanitizedName()},
				{Text: list2.SanitizedName()},
			},
		},
		"has current board & typing 'cd l'": {
			given: given{
				cmd:          "cd",
				args:         []string{"l"},
				currentBoard: &board1,
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "l").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: list1.Name},
			},
		},
		"has current board & current list & typing 'cat '": {
			given: given{
				cmd:          "cat",
				args:         []string{""},
				currentBoard: &board1,
				currentList:  &list1,
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCards(list1.ID).
						Return(cards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "card[shortLink]"},
				{Text: `another\ card[another shortLink]`},
			},
		},
		"has current board & current list & typing 'cat c'": {
			given: given{
				cmd:          "cat",
				args:         []string{"c"},
				currentBoard: &board1,
				currentList:  &list1,
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCards(list1.ID).
						Return(cards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "card[shortLink]"},
			},
		},
		"has current board & typing 'cd ../a'": {
			given: given{
				cmd:          "cd",
				args:         []string{"../a"},
				currentBoard: &board1,
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("a").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: `another\ board[another shortLink]`},
			},
		},
		"has current board & current list & typing 'cd ../a'": {
			given: given{
				cmd:          "cd",
				args:         []string{"../a"},
				currentBoard: &board1,
				currentList:  &list1,
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "a").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: `another\ list`},
			},
		},
		// ABSOLUTE PATHS
		"typing 'cd /'": {
			given: given{
				cmd:  "cd",
				args: []string{"/"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: "board[shortLink]"},
				{Text: `another\ board[another shortLink]`},
			},
		},
		"has current board & typing 'cd /a'": {
			given: given{
				cmd:          "cd",
				args:         []string{"/a"},
				currentBoard: &board1,
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("a").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: `another\ board[another shortLink]`},
			},
		},
		// ERRORS
		"server error when finding boards": {
			given: given{
				cmd:  "cd",
				args: []string{"/a"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("a").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindBoards().
						Return(nil, errors.New("unexpected error"))
					return tr
				},
			},
			expected: []prompt.Suggest{},
		},
		"server error when finding lists": {
			given: given{
				cmd:  "cd",
				args: []string{"/board/l"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "l").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists(board1.ID).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
			},
			expected: []prompt.Suggest{},
		},
		"server error when finding cards": {
			given: given{
				cmd:  "cat",
				args: []string{"/board/list/c"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCards(list1.ID).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
			},
			expected: []prompt.Suggest{},
		},
		"invalid path": {
			given: given{
				cmd:  "cd",
				args: []string{"/../../"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
			},
			expected: []prompt.Suggest{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := New(tt.given.buildTrelloRepository(), tt.given.currentBoard, tt.given.currentList)
			actual := c.Complete(tt.given.cmd, tt.given.args)
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
