package prompt

import (
	"github.com/l-lin/tcli/trello"
	"reflect"
	"testing"
)

func TestPrompt_LivePrefix(t *testing.T) {
	var tests = map[string]struct {
		given    *trello.Session
		expected string
	}{
		"/board/list": {
			given: &trello.Session{
				CurrentBoard: &trello.Board{Name: "board"},
				CurrentList:  &trello.List{Name: "list"},
			},
			expected: "/board/list> ",
		},
		"/board": {
			given: &trello.Session{
				CurrentBoard: &trello.Board{Name: "board"},
			},
			expected: "/board> ",
		},
		"/": {
			given:    &trello.Session{},
			expected: "/> ",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := Prompt{
				Session: tt.given,
			}
			actual, _ := s.LivePrefix()
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestGetCmd(t *testing.T) {
	type expected struct {
		cmd   string
		found bool
	}
	var tests = map[string]struct {
		given    string
		expected expected
	}{
		"no args": {
			given: "",
			expected: expected{
				cmd:   "",
				found: false,
			},
		},
		"single arg": {
			given: "cd",
			expected: expected{
				cmd:   "cd",
				found: true,
			},
		},
		"single arg with space after": {
			given: "cd ",
			expected: expected{
				cmd:   "cd",
				found: true,
			},
		},
		"two args": {
			given: "cd foobar",
			expected: expected{
				cmd:   "cd",
				found: true,
			},
		},
		"more than two args": {
			given: "cd foobar another",
			expected: expected{
				cmd:   "cd",
				found: true,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actualCmd, actualFound := getCmd(tt.given)
			if actualCmd != tt.expected.cmd {
				t.Errorf("expected %v, actual %v", tt.expected.cmd, actualCmd)
			}
			if actualFound != tt.expected.found {
				t.Errorf("expected %v, actual %v", tt.expected.found, actualFound)
			}
		})
	}
}

func TestGetArgs(t *testing.T) {
	type expected struct {
		args     []string
		hasError bool
	}
	var tests = map[string]struct {
		given    string
		expected expected
	}{
		"no cmd nor arg": {
			given: "",
			expected: expected{
				args: []string{},
			},
		},
		"single cmd": {
			given: "cd",
			expected: expected{
				args: []string{},
			},
		},
		"single cmd with space after": {
			given: "cd ",
			expected: expected{
				args: []string{},
			},
		},
		"cmd and one arg": {
			given: "cd foobar",
			expected: expected{
				args: []string{"foobar"},
			},
		},
		"cmd and one arg with space at the end": {
			given: "cd foobar ",
			expected: expected{
				args: []string{"foobar", ""},
			},
		},
		"cmd and two args": {
			given: "cd foobar another",
			expected: expected{
				args: []string{"foobar", "another"},
			},
		},
		"cmd and two args with first arg with escape space": {
			given: `cd foo\ bar another`,
			expected: expected{
				args: []string{"foo bar", "another"},
			},
		},
		"more complex example with escape space": {
			given: `cd path\ to\ board/list/card another\ path\ to\ board/list2 "with quotes"`,
			expected: expected{
				args: []string{"path to board/list/card", "another path to board/list2", "with quotes"},
			},
		},
		"containing unicodes": {
			given: `cd TODO/🎉\ DONE`,
			expected: expected{
				args: []string{"TODO/🎉 DONE"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actualArgs, actualErr := getArgs(tt.given)
			if tt.expected.hasError != (actualErr != nil) {
				t.Errorf("expected error %v, actual %v", tt.expected.hasError, actualErr)
				t.FailNow()
			}
			if !reflect.DeepEqual(actualArgs, tt.expected.args) {
				t.Errorf("expected %q, actual %q", tt.expected.args, actualArgs)
			}
		})
	}
}
