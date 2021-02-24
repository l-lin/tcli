package trello

import "testing"

func TestSanitize(t *testing.T) {
	var tests = map[string]struct {
		given    string
		expected string
	}{
		"string containing spaces": {
			given:    "foo bar",
			expected: `foo\ bar`,
		},
		"string without spaces": {
			given:    "foobar",
			expected: "foobar",
		},
		"empty string": {
			given:    "",
			expected: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := sanitize(tt.given)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestGetPos(t *testing.T) {
	var tests = map[string]struct {
		given    string
		expected interface{}
	}{
		"top": {
			given:    "top",
			expected: "top",
		},
		"bottom": {
			given:    "bottom",
			expected: "bottom",
		},
		"int number": {
			given:    "1234",
			expected: float64(1234),
		},
		"float number": {
			given:    "1234.56",
			expected: 1234.56,
		},
		"unknown value": {
			given:    "unknown",
			expected: "unknown",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := getPos(tt.given)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestToTCliID(t *testing.T) {
	type given struct {
		name, uid string
	}
	var tests = map[string]struct {
		given    given
		expected string
	}{
		"card": {
			given: given{
				name: "card",
				uid:  "shortLink",
			},
			expected: "card[shortLink]",
		},
		"list": {
			given: given{
				name: "list",
				uid:  "uuid",
			},
			expected: "list[uuid]",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := toTCliID(tt.given.name, tt.given.uid)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
