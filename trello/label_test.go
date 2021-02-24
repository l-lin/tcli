package trello

import (
	"reflect"
	"testing"
)

var colors = []string{
	"black",
	"blue",
	"green",
	"lime",
	"orange",
	"pink",
	"purple",
	"red",
	"sky",
	"yellow",
}

func TestLabels_String(t *testing.T) {
	var tests = map[string]struct {
		given    Labels
		expected string
	}{
		"3 labels": {
			given: Labels{
				{ID: "label 1", Name: "label name 1", Color: "red"},
				{ID: "label 2", Name: "label name 2", Color: "sky"},
				{ID: "label 3", Name: "label name 3", Color: "black"},
			},
			expected: "label 1,label 2,label 3",
		},
		"1 label": {
			given: Labels{
				{ID: "label 1", Name: "label name 1", Color: "red"},
			},
			expected: "label 1",
		},
		"no label": {
			given:    Labels{},
			expected: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.String()
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestLabels_Slice(t *testing.T) {
	var tests = map[string]struct {
		given    Labels
		expected []string
	}{
		"3 labels": {
			given: Labels{
				{ID: "id red", Name: "name red", Color: "red"},
				{ID: "id sky", Name: "name sky", Color: "sky"},
				{ID: "id black", Color: "black"},
			},
			expected: []string{"red [name red]", "sky [name sky]", "black"},
		},
		"1 label": {
			given: Labels{
				{ID: "id label red", Name: "name red", Color: "red"},
			},
			expected: []string{"red [name red]"},
		},
		"no label": {
			given:    Labels{},
			expected: []string{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.ToSliceTCliColors()
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestLabels_FilterBy(t *testing.T) {
	labels := Labels{}
	for _, color := range colors {
		labels = append(labels, Label{Color: color, ID: "id " + color, Name: "name " + color})
	}
	var tests = map[string]struct {
		given    []string
		expected Labels
	}{
		"3 labels and filtering by id, TCliColor and color": {
			given: []string{"red", "green [name green]", "id sky"},
			expected: Labels{
				{Color: "red", ID: "id red", Name: "name red"},
				{Color: "green", ID: "id green", Name: "name green"},
				{Color: "sky", ID: "id sky", Name: "name sky"},
			},
		},
		"1 label and filtering by color": {
			given: []string{"black"},
			expected: Labels{
				{Color: "black", ID: "id black", Name: "name black"},
			},
		},
		"no label": {
			given:    []string{},
			expected: Labels{},
		},
		"nonexistent color": {
			given: []string{"unknown", "red"},
			expected: Labels{
				{Color: "red", ID: "id red", Name: "name red"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := labels.FilterBy(tt.given, LabelFilterOr(
				LabelFilterByID,
				LabelFilterByTCliColor,
				LabelFilterByColor,
			))
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestLabels_IDLabelsInString(t *testing.T) {
	var tests = map[string]struct {
		given    Labels
		expected string
	}{
		"3 labels": {
			given: Labels{
				{Color: "red", ID: "id red"},
				{Color: "green", ID: "id green"},
				{Color: "sky", ID: "id sky"},
			},
			expected: "id red,id green,id sky",
		},
		"1 labels": {
			given: Labels{
				{Color: "red", ID: "id red"},
			},
			expected: "id red",
		},
		"no labels": {
			given:    Labels{},
			expected: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.IDLabelsInString()
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestLabels_FilterBy_Then_IDLabelsInString(t *testing.T) {
	labels := Labels{}
	for _, color := range colors {
		labels = append(labels, Label{Color: color, ID: "id " + color, Name: "name " + color})
	}
	var tests = map[string]struct {
		given    []string
		expected string
	}{
		"3 labels": {
			given:    []string{"red", "green [name green]", "id sky"},
			expected: "id red,id green,id sky",
		},
		"1 label": {
			given:    []string{"black"},
			expected: "id black",
		},
		"no label": {
			given:    []string{},
			expected: "",
		},
		"nonexistent color": {
			given:    []string{"unknown", "red"},
			expected: "id red",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := labels.FilterBy(tt.given, LabelFilterOr(
				LabelFilterByID,
				LabelFilterByTCliColor,
				LabelFilterByColor,
			)).IDLabelsInString()
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
