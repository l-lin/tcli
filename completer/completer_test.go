package completer

import (
	"errors"
	"fmt"
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
		session               *trello.Session
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

	card1 := trello.Card{Name: "card", ShortLink: "shortLink"}
	card2 := trello.Card{Name: "another card", ShortLink: "another shortLink"}
	cards := trello.Cards{card1, card2}

	comment1 := trello.Comment{
		ID: "comment",
		MemberCreator: trello.CommentMemberCreator{
			Username: "username 1",
		},
		Date: "2014-12-12T11:45:26.371Z",
	}
	comment2 := trello.Comment{
		ID: "another comment",
		MemberCreator: trello.CommentMemberCreator{
			Username: "username 2",
		},
		Date: "2014-12-13T11:45:26.371Z",
	}
	comments := trello.Comments{comment1, comment2}

	notFoundError := errors.New("not found")

	var tests = map[string]struct {
		given    given
		expected []prompt.Suggest
	}{
		"/> ": {
			given: given{
				buildTrelloRepository: func() trello.Repository { return nil },
				session:               &trello.Session{},
			},
			expected: commandSuggestions(),
		},
		"/> unknown": {
			given: given{
				cmd: "unknown",
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		// COMMANDS
		"/> c": {
			given: given{
				cmd:  "c",
				args: []string{""},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{
				{Text: "clear", Description: "clear the terminal screen & cache"},
				{Text: "cd", Description: "change level in the hierarchy"},
				{Text: "cat", Description: "show resource content info"},
				{Text: "cp", Description: "copy resource"},
			},
		},
		// RELATIVE PATHS
		"/board > cat ": {
			given: given{
				cmd:     "cat",
				args:    []string{""},
				session: &trello.Session{Board: &board1},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: list1.TCliID()},
				{Text: list2.TCliID()},
			},
		},
		"/board > cat l": {
			given: given{
				cmd:     "cat",
				args:    []string{"l"},
				session: &trello.Session{Board: &board1},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "l").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: list1.TCliID()},
			},
		},
		"/board > cd unknown": {
			given: given{
				cmd:     "cd",
				args:    []string{"unknown"},
				session: &trello.Session{Board: &board1},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "unknown").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{},
		},
		"/board > cat list/": {
			given: given{
				cmd:     "cat",
				args:    []string{"list/"},
				session: &trello.Session{Board: &board1},
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
				{Text: card1.TCliID()},
				{Text: card2.TCliID()},
			},
		},
		"/board > cat list/c": {
			given: given{
				cmd:     "cat",
				args:    []string{"list/c"},
				session: &trello.Session{Board: &board1},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, "c").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindCards(list1.ID).
						Return(cards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: card1.TCliID()},
			},
		},
		"/board/list > cat ../a": {
			given: given{
				cmd:     "cat",
				args:    []string{"../a"},
				session: &trello.Session{Board: &board1, List: &list1},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "a").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: list2.TCliID()},
			},
		},
		// ABSOLUTE PATHS
		"/board > cd /": {
			given: given{
				cmd:     "cd",
				args:    []string{"/"},
				session: &trello.Session{Board: &board1},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: board1.TCliID()},
				{Text: board2.TCliID()},
			},
		},
		"/board > cd /a": {
			given: given{
				cmd:     "cd",
				args:    []string{"/a"},
				session: &trello.Session{Board: &board1},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("a").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
			},
			expected: []prompt.Suggest{
				{Text: board2.TCliID()},
			},
		},
		// CD
		"/> cd board/list/card/c": {
			given: given{
				cmd:  "cd",
				args: []string{"board/list/card/c"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		"/> cd board/list another-board": {
			given: given{
				cmd:  "cd",
				args: []string{"board/list", "another-board"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		// CP
		"/> cp board/list/card/c": {
			given: given{
				cmd:  "cp",
				args: []string{"board/list/card/c"},
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
					tr.EXPECT().
						FindComments(card1.ID).
						Return(comments, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{
				{Text: comment1.ID, Description: fmt.Sprintf("%s @ %s", comment1.MemberCreator.Username, comment1.Date)},
			},
		},
		"/> cp board/list/card/comment board/list/ano": {
			given: given{
				cmd:  "cp",
				args: []string{"board/list/card/comment", "board/list/ano"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, "ano").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindCards(list1.ID).
						Return(cards, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{
				{Text: card2.TCliID()},
			},
		},
		"/> cp board/list/card/comment board/list/card/ano": {
			given: given{
				cmd:  "cp",
				args: []string{"board/list/card/comment", "board/list/card/ano"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		"/> cp 1 2 3": {
			given: given{
				cmd:  "cp",
				args: []string{"1", "2", "3"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		// MV
		"/> mv bo": {
			given: given{
				cmd:  "mv",
				args: []string{"bo"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("bo").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{
				{Text: board1.TCliID()},
			},
		},
		"/> mv board/list/c": {
			given: given{
				cmd:  "mv",
				args: []string{"board/list/c"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, "c").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindCards(list1.ID).
						Return(cards, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{
				{Text: card1.TCliID()},
			},
		},
		"/> mv board/list/card bo": {
			given: given{
				cmd:  "mv",
				args: []string{"board/list/card", "bo"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("bo").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{
				{Text: board1.TCliID()},
			},
		},
		"/> mv board/list/card board/ano": {
			given: given{
				cmd:  "mv",
				args: []string{"board/list/card", "board/ano"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "ano").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					return tr
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{
				{Text: list2.TCliID()},
			},
		},
		"/> mv board/list/card board/another-list/c": {
			given: given{
				cmd:  "mv",
				args: []string{"board/list/card", "board/another-list/c"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		"/> mv board/list/card board/another-list another-board/another-list": {
			given: given{
				cmd:  "mv",
				args: []string{"board/list/card", "board/another-list", "another-board/another-list"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		"/> mv board/list/card /../..": {
			given: given{
				cmd:  "mv",
				args: []string{"board/list/card", "/../.."},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		// ERRORS
		"/> cat a (server error when finding boards)": {
			given: given{
				cmd:  "cat",
				args: []string{"/a"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("a").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindBoards().
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		"/> cat /board/l (server error when finding lists)": {
			given: given{
				cmd:  "cat",
				args: []string{"/board/l"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "l").
						Return(nil, notFoundError)
					tr.EXPECT().
						FindLists(board1.ID).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		"/> cd /board/list/c (server error when finding cards)": {
			given: given{
				cmd:  "cd",
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
						FindCard(list1.ID, "c").
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindCards(list1.ID).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		"/> cp /board/list/card/c (server error when finding comments)": {
			given: given{
				cmd:  "cp",
				args: []string{"/board/list/card/c"},
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
					tr.EXPECT().
						FindComments(card1.ID).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		"/> cat /../../ (invalid path)": {
			given: given{
				cmd:  "cat",
				args: []string{"/../../"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
		"/> cd /../../ (invalid path)": {
			given: given{
				cmd:  "cd",
				args: []string{"/../../"},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				session: &trello.Session{},
			},
			expected: []prompt.Suggest{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := New(tt.given.buildTrelloRepository(), tt.given.session)
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
