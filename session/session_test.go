package session

import "testing"

func TestSession_Executor(t *testing.T) {
	// TODO
}

func TestSession_Completer(t *testing.T) {
	// TODO
}

func TestSession_LivePrefix(t *testing.T) {
	// TODO
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

func TestGetArg(t *testing.T) {
	var tests = map[string]struct {
		given    string
		expected string
	}{
		"no arg": {
			given:    "",
			expected: "",
		},
		"single arg": {
			given:    "cd",
			expected: "",
		},
		"single arg with space after": {
			given:    "cd ",
			expected: "",
		},
		"two args": {
			given:    "cd foobar",
			expected: "foobar",
		},
		"more than two args": {
			given:    "cd foobar another",
			expected: "foobar another",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := getArg(tt.given)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
