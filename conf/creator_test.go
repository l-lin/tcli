package conf

import (
	"testing"
)

func TestValidateNotEmpty(t *testing.T) {
	var tests = map[string]struct {
		given            string
		expectedHasError bool
	}{
		"valid content": {
			given:            "foobar",
			expectedHasError: false,
		},
		"only spaces": {
			given:            "  ",
			expectedHasError: true,
		},
		"only tabs": {
			given:            "\t",
			expectedHasError: true,
		},
		"empty string": {
			given:            "",
			expectedHasError: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := validateNotEmpty(tt.given)
			if tt.expectedHasError && actual == nil {
				t.Errorf("expected error")
			}
			if !tt.expectedHasError && actual != nil {
				t.Errorf("expected valid content")
			}
		})
	}
}
