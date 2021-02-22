package renderer

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/muesli/termenv"
	"strings"
)

var (
	foregroundColor = termenv.ANSIBlack
	mapColors       = map[string]termenv.Color{
		"black":  termenv.ANSIBrightBlack,
		"blue":   termenv.ANSIBrightBlue,
		"green":  termenv.ANSIBrightGreen,
		"lime":   termenv.ANSIGreen,
		"orange": termenv.ANSIYellow,
		"pink":   termenv.ANSIMagenta,
		"purple": termenv.ANSIBrightMagenta,
		"red":    termenv.ANSIBrightRed,
		"sky":    termenv.ANSIBrightCyan,
		"yellow": termenv.ANSIBrightYellow,
	}
)

// Labels rendering labels
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

type TermEnvLabel struct{}

func (t TermEnvLabel) Render(labels trello.Labels) string {
	sb := strings.Builder{}
	for _, label := range labels {
		labelName := ""
		if label.Name != "" {
			labelName = fmt.Sprintf(" %s ", label.Name)
		} else {
			labelName = "      "
		}
		sb.WriteString(
			termenv.String(labelName).
				Background(getTermEnvColor(label.Color)).
				Foreground(foregroundColor).
				String(),
		)
		sb.WriteString(" ")
	}
	return sb.String()
}

func getTermEnvColor(color string) termenv.Color {
	if c := mapColors[color]; c != nil {
		return c
	}
	return termenv.ANSIWhite
}
