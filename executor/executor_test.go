package executor

import (
	"github.com/l-lin/tcli/conf"
	"testing"
)

func TestNew(t *testing.T) {
	var tests = map[string]struct {
		given       string
		expectedNil bool
	}{
		"cd": {
			given:       "cd",
			expectedNil: false,
		},
		"unknown command": {
			given:       "unknown",
			expectedNil: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := New(conf.Conf{}, tt.given, nil, nil, nil)
			if tt.expectedNil && actual != nil || !tt.expectedNil && actual == nil {
				t.Errorf("expected %v, actual %v", tt.expectedNil, actual == nil)
			}
		})
	}
}
