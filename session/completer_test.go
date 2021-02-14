package session

import "testing"

func TestCompleterAtBoardsLevel(t *testing.T) {
	// TODO
}

func TestCompleterAtListsLevel(t *testing.T) {
	// TODO
}

func TestCompleterAtCardsLevel(t *testing.T) {
	// TODO
}

func TestTruncateCardDescription(t *testing.T) {
	var tests = map[string]struct {
		given    string
		expected string
	}{
		"empty description": {
			given:    "",
			expected: "",
		},
		"long description": {
			given:    "long description that exceed the threshold",
			expected: "long description tha",
		},
		"short description": {
			given:    "short desc",
			expected: "short desc",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := truncateCardDescription(tt.given)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
