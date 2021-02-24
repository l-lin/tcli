package trello

import (
	"fmt"
	"strings"
)

type Labels []Label

func (l Labels) String() string {
	if len(l) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i := 0; i < len(l); i++ {
		sb.WriteString(l[i].ID)
		if i < len(l)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

type LabelFilter func(s string, l Label) bool

var LabelFilterByID = func(s string, l Label) bool {
	return s == l.ID
}
var LabelFilterByTCliColor = func(s string, l Label) bool {
	return s == l.ToTCliColor()
}
var LabelFilterByColor = func(s string, l Label) bool {
	return s == l.Color
}
var LabelFilterOr = func(filters ...LabelFilter) LabelFilter {
	return func(s string, l Label) bool {
		for _, filter := range filters {
			if filter(s, l) {
				return true
			}
		}
		return false
	}
}

func (l Labels) FilterBy(labels []string, filter LabelFilter) Labels {
	filtered := Labels{}
	for _, s := range labels {
		for _, label := range l {
			if filter(s, label) {
				filtered = append(filtered, label)
				continue
			}
		}
	}
	return filtered
}

func (l Labels) IDLabelsInString() string {
	var idLabels []string
	for _, label := range l {
		idLabels = append(idLabels, label.ID)
	}
	return strings.Join(idLabels, ",")
}

func (l Labels) ToSliceTCliColors() []string {
	s := make([]string, len(l))
	for i := 0; i < len(l); i++ {
		s[i] = l[i].ToTCliColor()
	}
	return s
}

type Label struct {
	ID      string `json:"id"`
	IDBoard string `json:"idBoard"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

func (l Label) ToTCliColor() string {
	if l.Name == "" {
		return l.Color
	}
	return fmt.Sprintf("%s [%s]", l.Color, l.Name)
}
