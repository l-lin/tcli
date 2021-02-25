package trello

import (
	"reflect"
	"testing"
)

func TestNewPathResolver(t *testing.T) {
	var tests = map[string]struct {
		given    *Session
		expected PathResolver
	}{
		"existing board, list, card": {
			given: &Session{
				Board: &Board{Name: "board"},
				List:  &List{Name: "list"},
				Card:  &Card{Name: "card"},
			},
			expected: PathResolver{
				Path: Path{
					BoardName: "board",
					ListName:  "list",
					CardName:  "card",
				},
			},
		},
		"existing board and list": {
			given: &Session{
				Board: &Board{Name: "board"},
				List:  &List{Name: "list"},
			},
			expected: PathResolver{
				Path: Path{
					BoardName: "board",
					ListName:  "list",
				},
			},
		},
		"existing board and non existing list": {
			given: &Session{
				Board: &Board{Name: "board"},
				List:  nil,
			},
			expected: PathResolver{
				Path: Path{
					BoardName: "board",
				},
			},
		},
		"non existing board and non existing list": {
			given: &Session{
				Board: nil,
				List:  nil,
			},
			expected: PathResolver{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := NewPathResolver(tt.given)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestPathResolver_Resolve(t *testing.T) {
	type given struct {
		currentPath  Path
		relativePath string
	}
	type expected struct {
		resolvedPath Path
		err          error
	}
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"full path": {
			given: given{
				currentPath: Path{
					BoardName: "board",
					ListName:  "list",
					CardName:  "card",
				},
				relativePath: "comment",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board",
					ListName:  "list",
					CardName:  "card",
					CommentID: "comment",
				},
				err: nil,
			},
		},
		"empty relative path": {
			given: given{
				currentPath: Path{
					BoardName: "board",
					ListName:  "list",
					CardName:  "card",
				},
				relativePath: "",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board",
					ListName:  "list",
					CardName:  "card",
				},
				err: nil,
			},
		},
		"no list": {
			given: given{
				currentPath: Path{
					BoardName: "board",
				},
				relativePath: "list",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board",
					ListName:  "list",
				},
				err: nil,
			},
		},
		"no board and no list": {
			given: given{
				currentPath:  Path{},
				relativePath: "board/list",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board",
					ListName:  "list",
				},
				err: nil,
			},
		},
		"full path in relativePath": {
			given: given{
				relativePath: "board/list/card",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board",
					ListName:  "list",
					CardName:  "card",
				},
				err: nil,
			},
		},
		"boardName set and rest of the path in relativePath": {
			given: given{
				currentPath: Path{
					BoardName: "board",
				},
				relativePath: "list/card",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board",
					ListName:  "list",
					CardName:  "card",
				},
				err: nil,
			},
		},
		"using .. in relativePath": {
			given: given{
				currentPath: Path{
					BoardName: "board",
					ListName:  "list",
				},
				relativePath: "../list2/card",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board",
					ListName:  "list2",
					CardName:  "card",
				},
				err: nil,
			},
		},
		"using ../.. in relativePath": {
			given: given{
				currentPath: Path{
					BoardName: "board",
					ListName:  "list",
				},
				relativePath: "../../board2/list2/card",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board2",
					ListName:  "list2",
					CardName:  "card",
				},
				err: nil,
			},
		},
		"have / at the end of the relativePath": {
			given: given{
				currentPath: Path{
					BoardName: "board",
					ListName:  "list",
				},
				relativePath: "card/",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board",
					ListName:  "list",
					CardName:  "card",
				},
				err: nil,
			},
		},
		"using absolute path in relativePath": {
			given: given{
				currentPath: Path{
					BoardName: "board",
					ListName:  "list",
				},
				relativePath: "/board2/list2/card",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board2",
					ListName:  "list2",
					CardName:  "card",
				},
				err: nil,
			},
		},
		"have / at the end of the absolute path in relativePath": {
			given: given{
				currentPath: Path{
					BoardName: "board",
					ListName:  "list",
				},
				relativePath: "/board2/list2/card/",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board2",
					ListName:  "list2",
					CardName:  "card",
				},
				err: nil,
			},
		},
		".. in relativePath from /board/list": {
			given: given{
				currentPath: Path{
					BoardName: "board",
					ListName:  "list",
				},
				relativePath: "..",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board",
				},
				err: nil,
			},
		},
		".. in relativePath from /board": {
			given: given{
				currentPath: Path{
					BoardName: "board",
				},
				relativePath: "..",
			},
			expected: expected{
				resolvedPath: Path{},
				err:          nil,
			},
		},
		"board name containing escaped space": {
			given: given{
				currentPath:  Path{},
				relativePath: "/board\\ name/list/card/",
			},
			expected: expected{
				resolvedPath: Path{
					BoardName: "board\\ name",
					ListName:  "list",
					CardName:  "card",
				},
				err: nil,
			},
		},
		// ERRORS --------------------------------------------------
		"more than allowed components": {
			given: given{
				currentPath: Path{
					BoardName: "board",
					ListName:  "list",
					CardName:  "card",
				},
				relativePath: "comment/invalid",
			},
			expected: expected{
				err: invalidPathErr,
			},
		},
		"using .. when already at top level": {
			given: given{
				relativePath: "../invalid",
			},
			expected: expected{
				err: invalidPathErr,
			},
		},
		"using .. too much": {
			given: given{
				currentPath: Path{
					BoardName: "board",
				},
				relativePath: "../../invalid",
			},
			expected: expected{
				err: invalidPathErr,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := PathResolver{
				Path: tt.given.currentPath,
			}
			actual, actualErr := r.Resolve(tt.given.relativePath)
			if actualErr != tt.expected.err {
				t.Errorf("expected %v, actual %v", tt.expected.err, actualErr)
			}
			if !reflect.DeepEqual(actual, tt.expected.resolvedPath) {
				t.Errorf("expected %v, actual %v", tt.expected.resolvedPath, actual)
			}
		})
	}
}
