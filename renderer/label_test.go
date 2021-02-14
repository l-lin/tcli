package renderer

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/logrusorgru/aurora/v3"
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

func TestAuroraLabelRenderer_Render(t *testing.T) {
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
				aurora.Sprintf(aurora.BgBrightYellow(" %s "), "yellow label"),
				aurora.Sprintf(aurora.BgBrightCyan(" %s "), "sky label"),
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
				aurora.Sprintf(aurora.BgBrightYellow(" %s "), "    "),
				aurora.Sprintf(aurora.BgBrightCyan(" %s "), "    "),
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
				aurora.Sprintf(aurora.BgBrightYellow(" %s "), "yellow label"),
				aurora.Sprintf(aurora.BgBrightCyan(" %s "), "    "),
			),
		},
		"no labels": {
			given:    trello.Labels{},
			expected: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ar := AuroraLabel{}
			actual := ar.Render(tt.given)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
