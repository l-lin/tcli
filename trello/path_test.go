package trello

import (
	"testing"
)

func TestNewPathResolver(t *testing.T) {
	type given struct {
		currentBoard *Board
		currentList  *List
	}
	var tests = map[string]struct {
		given    given
		expected PathResolver
	}{
		"existing board and list": {
			given: given{
				currentBoard: &Board{Name: "board"},
				currentList:  &List{Name: "list"},
			},
			expected: PathResolver{
				currentBoardName: "board",
				currentListName:  "list",
			},
		},
		"existing board and non existing list": {
			given: given{
				currentBoard: &Board{Name: "board"},
				currentList:  nil,
			},
			expected: PathResolver{
				currentBoardName: "board",
				currentListName:  "",
			},
		},
		"non existing board and non existing list": {
			given: given{
				currentBoard: nil,
				currentList:  nil,
			},
			expected: PathResolver{
				currentBoardName: "",
				currentListName:  "",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := NewPathResolver(tt.given.currentBoard, tt.given.currentList)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestPathResolver_Resolve(t *testing.T) {
	type given struct {
		currentBoard string
		currentList  string
		relativePath string
	}
	type expected struct {
		boardName string
		listName  string
		cardName  string
		err       error
	}
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"full path": {
			given: given{
				currentBoard: "board",
				currentList:  "list",
				relativePath: "card",
			},
			expected: expected{
				boardName: "board",
				listName:  "list",
				cardName:  "card",
				err:       nil,
			},
		},
		"empty relative path": {
			given: given{
				currentBoard: "board",
				currentList:  "list",
				relativePath: "",
			},
			expected: expected{
				boardName: "board",
				listName:  "list",
				cardName:  "",
				err:       nil,
			},
		},
		"no list": {
			given: given{
				currentBoard: "board",
				currentList:  "",
				relativePath: "list",
			},
			expected: expected{
				boardName: "board",
				listName:  "list",
				cardName:  "",
				err:       nil,
			},
		},
		"no board and no list": {
			given: given{
				currentBoard: "",
				currentList:  "",
				relativePath: "board/list",
			},
			expected: expected{
				boardName: "board",
				listName:  "list",
				cardName:  "",
				err:       nil,
			},
		},
		"full path in relativePath": {
			given: given{
				currentBoard: "",
				currentList:  "",
				relativePath: "board/list/card",
			},
			expected: expected{
				boardName: "board",
				listName:  "list",
				cardName:  "card",
				err:       nil,
			},
		},
		"boardName set and rest of the path in relativePath": {
			given: given{
				currentBoard: "board",
				currentList:  "",
				relativePath: "list/card",
			},
			expected: expected{
				boardName: "board",
				listName:  "list",
				cardName:  "card",
				err:       nil,
			},
		},
		"using .. in relativePath": {
			given: given{
				currentBoard: "board",
				currentList:  "list",
				relativePath: "../list2/card",
			},
			expected: expected{
				boardName: "board",
				listName:  "list2",
				cardName:  "card",
				err:       nil,
			},
		},
		"using ../.. in relativePath": {
			given: given{
				currentBoard: "board",
				currentList:  "list",
				relativePath: "../../board2/list2/card",
			},
			expected: expected{
				boardName: "board2",
				listName:  "list2",
				cardName:  "card",
				err:       nil,
			},
		},
		"have / at the end of the relativePath": {
			given: given{
				currentBoard: "board",
				currentList:  "list",
				relativePath: "card/",
			},
			expected: expected{
				boardName: "board",
				listName:  "list",
				cardName:  "card",
				err:       nil,
			},
		},
		"using absolute path in relativePath": {
			given: given{
				currentBoard: "board",
				currentList:  "list",
				relativePath: "/board2/list2/card",
			},
			expected: expected{
				boardName: "board2",
				listName:  "list2",
				cardName:  "card",
				err:       nil,
			},
		},
		"have / at the end of the absolute path in relativePath": {
			given: given{
				currentBoard: "board",
				currentList:  "list",
				relativePath: "/board2/list2/card/",
			},
			expected: expected{
				boardName: "board2",
				listName:  "list2",
				cardName:  "card",
				err:       nil,
			},
		},
		".. in relativePath from /board/list": {
			given: given{
				currentBoard: "board",
				currentList:  "list",
				relativePath: "..",
			},
			expected: expected{
				boardName: "board",
				listName:  "",
				cardName:  "",
				err:       nil,
			},
		},
		".. in relativePath from /board": {
			given: given{
				currentBoard: "board",
				currentList:  "",
				relativePath: "..",
			},
			expected: expected{
				boardName: "",
				listName:  "",
				cardName:  "",
				err:       nil,
			},
		},
		// ERRORS --------------------------------------------------
		"more than 4 components": {
			given: given{
				currentBoard: "board",
				currentList:  "list",
				relativePath: "card/invalid",
			},
			expected: expected{
				boardName: "",
				listName:  "",
				cardName:  "",
				err:       invalidPathErr,
			},
		},
		"using .. when already at top level": {
			given: given{
				currentBoard: "",
				currentList:  "",
				relativePath: "../invalid",
			},
			expected: expected{
				boardName: "",
				listName:  "",
				cardName:  "",
				err:       invalidPathErr,
			},
		},
		"using .. too much": {
			given: given{
				currentBoard: "board",
				currentList:  "",
				relativePath: "../../invalid",
			},
			expected: expected{
				boardName: "",
				listName:  "",
				cardName:  "",
				err:       invalidPathErr,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := PathResolver{
				currentBoardName: tt.given.currentBoard,
				currentListName:  tt.given.currentList,
			}
			actualBoardName, actualListName, actualCardName, actualErr := r.Resolve(tt.given.relativePath)
			if actualBoardName != tt.expected.boardName {
				t.Errorf("expected %v, actual %v", tt.expected.boardName, actualBoardName)
			}
			if actualListName != tt.expected.listName {
				t.Errorf("expected %v, actual %v", tt.expected.listName, actualListName)
			}
			if actualCardName != tt.expected.cardName {
				t.Errorf("expected %v, actual %v", tt.expected.cardName, actualCardName)
			}
			if actualErr != tt.expected.err {
				t.Errorf("expected %v, actual %v", tt.expected.err, actualErr)
			}
		})
	}
}
