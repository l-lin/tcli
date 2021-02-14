package renderer

import (
	"github.com/l-lin/tcli/trello"
	"github.com/logrusorgru/aurora/v3"
	"strings"
)

type Labels interface {
	Render(labels trello.Labels) string
}

type PlainLabel struct{}

func (p PlainLabel) Render(labels trello.Labels) string {
	sb := strings.Builder{}
	for _, label := range labels {
		if label.Name != "" {
			sb.WriteString(label.Name)
		} else {
			sb.WriteString(label.Color)
		}
		sb.WriteString(" ")
	}
	return sb.String()
}

type AuroraLabel struct{}

func (ar AuroraLabel) Render(labels trello.Labels) string {
	sb := strings.Builder{}
	for _, label := range labels {
		if label.Name != "" {
			sb.WriteString(aurora.Sprintf(label.Colorize(" %s "), label.Name))
		} else {
			sb.WriteString(label.Colorize("      ").String())
		}
		sb.WriteString(" ")
	}
	return sb.String()
}
