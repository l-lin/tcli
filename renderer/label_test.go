package renderer

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/muesli/termenv"
	"testing"
)

func TestPlainLabelRenderer_Render(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Labels
		expected string
	}{
		"two labels with name": {
			given: trello.Labels{
				trello.Label{
					Name:  "yellow label",
					Color: "yellow",
				},
				trello.Label{
					Name:  "sky label",
					Color: "sky",
				},
			},
			expected: "yellow label sky label ",
		},
		"two labels without name": {
			given: trello.Labels{
				trello.Label{
					Color: "yellow",
				},
				trello.Label{
					Color: "sky",
				},
			},
			expected: "yellow sky ",
		},
		"one labels with name and another without": {
			given: trello.Labels{
				trello.Label{
					Name:  "yellow label",
					Color: "yellow",
				},
				trello.Label{
					Color: "sky",
				},
			},
			expected: "yellow label sky ",
		},
		"no label": {
			given:    trello.Labels{},
			expected: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			pl := PlainLabel{}
			actual := pl.Render(tt.given)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestTermenvLabelRenderer_Render(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Labels
		expected string
	}{
		"two labels with name": {
			given: trello.Labels{
				trello.Label{
					Name:  "yellow label",
					Color: "yellow",
				},
				trello.Label{
					Name:  "sky label",
					Color: "sky",
				},
			},
			expected: fmt.Sprintf("%s %s ",
				termenv.String(" yellow label ").Background(mapColors["yellow"]).Foreground(foregroundColor),
				termenv.String(" sky label ").Background(mapColors["sky"]).Foreground(foregroundColor),
			),
		},
		"two labels without name": {
			given: trello.Labels{
				trello.Label{
					Color: "yellow",
				},
				trello.Label{
					Color: "sky",
				},
			},
			expected: fmt.Sprintf("%s %s ",
				termenv.String("      ").Background(mapColors["yellow"]).Foreground(foregroundColor),
				termenv.String("      ").Background(mapColors["sky"]).Foreground(foregroundColor),
			),
		},
		"one labels with name and another without": {
			given: trello.Labels{
				trello.Label{
					Name:  "yellow label",
					Color: "yellow",
				},
				trello.Label{
					Color: "sky",
				},
			},
			expected: fmt.Sprintf("%s %s ",
				termenv.String(" yellow label ").Background(mapColors["yellow"]).Foreground(foregroundColor),
				termenv.String("      ").Background(mapColors["sky"]).Foreground(foregroundColor),
			),
		},
		"label without registered color": {
			given: trello.Labels{
				trello.Label{
					Name:  "unknown label",
					Color: "unknown",
				},
			},
			expected: fmt.Sprintf("%s ",
				termenv.String(" unknown label ").Background(termenv.ANSIWhite).Foreground(foregroundColor),
			),
		},
		"no labels": {
			given:    trello.Labels{},
			expected: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ar := TermEnvLabel{}
			actual := ar.Render(tt.given)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
