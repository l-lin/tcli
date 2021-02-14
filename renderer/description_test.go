package renderer

import "testing"

func TestPlainDescription_Render(t *testing.T) {
	type expected struct {
		result string
		err    error
	}
	var tests = map[string]struct {
		given    string
		expected expected
	}{
		"empty": {
			given: "",
			expected: expected{
				result: "",
				err:    nil,
			},
		},
		"some description": {
			given: "some description",
			expected: expected{
				result: "some description",
				err:    nil,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			pdr := PlainDescription{}
			actualResult, actualErr := pdr.Render(tt.given)
			if actualResult != tt.expected.result {
				t.Errorf("expected %v, actual %v", tt.expected, actualResult)
			}
			if actualErr != tt.expected.err {
				t.Errorf("expected %v, actual %v", tt.expected.err, actualErr)
			}
		})
	}
}
