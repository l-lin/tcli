package renderer

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/logrusorgru/aurora/v3"
	"testing"
)

func TestInTable_RenderBoards(t *testing.T) {
	type given func() trello.Boards
	var tests = map[string]struct {
		given    given
		expected string
	}{
		"two boards": {
			given: func() trello.Boards {
				return trello.Boards{
					trello.Board{
						ID:               "1",
						Name:             "Board 1",
						ShortURL:         "https://trello.com/b/azerty",
						DateLastActivity: "2021-02-04T14:19:25.229Z",
					},
					trello.Board{
						ID:               "2",
						Name:             "Board 2",
						ShortURL:         "https://trello.com/b/popo",
						DateLastActivity: "2021-02-08T21:02:58.117Z",
					},
				}
			},
			expected: `Name       ID    Short URL                      Last activity date
----       --    ---------                      ------------------
Board 1    1     https://trello.com/b/azerty    2021-02-04T14:19:25.229Z
Board 2    2     https://trello.com/b/popo      2021-02-08T21:02:58.117Z
`,
		},
		"no board": {
			given: func() trello.Boards {
				return trello.Boards{}
			},
			expected: `Name    ID    Short URL    Last activity date
----    --    ---------    ------------------
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer()
			actual := r.RenderBoards(tt.given())
			if actual != tt.expected {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected, actual)
			}
		})
	}
}

func TestInTable_RenderLists(t *testing.T) {
	type given func() trello.Lists
	var tests = map[string]struct {
		given    given
		expected string
	}{
		"two lists": {
			given: func() trello.Lists {
				return trello.Lists{
					trello.List{
						ID:   "1",
						Name: "List 1",
					},
					trello.List{
						ID:   "2",
						Name: "List 2",
					},
				}
			},
			expected: `Name      ID
----      --
List 1    1
List 2    2
`,
		},
		"no list": {
			given: func() trello.Lists {
				return trello.Lists{}
			},
			expected: `Name    ID
----    --
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer()
			actual := r.RenderLists(tt.given())
			if actual != tt.expected {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected, actual)
			}
		})
	}
}

func TestInTable_RenderCards(t *testing.T) {
	type given func() trello.Cards
	var tests = map[string]struct {
		given    given
		expected string
	}{
		"two cards with labels": {
			given: func() trello.Cards {
				return trello.Cards{
					trello.Card{
						ID:   "1",
						Name: "Card 1",
						Labels: trello.Labels{
							trello.Label{
								ID:    "10",
								Name:  "Label 10",
								Color: "green",
							},
							trello.Label{
								ID:    "11",
								Name:  "Label 11",
								Color: "yellow",
							},
						},
					},
					trello.Card{
						ID:   "2",
						Name: "Card 2",
						Labels: trello.Labels{
							trello.Label{
								ID:    "20",
								Name:  "Label 20",
								Color: "black",
							},
							trello.Label{
								ID:    "21",
								Name:  "Label 21",
								Color: "sky",
							},
						},
					},
				}
			},
			expected: aurora.Sprintf(`Name      ID    Labels
----      --    ------
Card 1    1     %s%s
Card 2    2     %s%s
`,
				aurora.BrightGreen("██ Label 10 "),
				aurora.BrightYellow("██ Label 11 "),
				aurora.BrightBlack("██ Label 20 "),
				aurora.BrightCyan("██ Label 21 "),
			),
		},
		"two cards without label": {
			given: func() trello.Cards {
				return trello.Cards{
					trello.Card{
						ID:     "1",
						Name:   "Card 1",
						Labels: trello.Labels{},
					},
					trello.Card{
						ID:     "2",
						Name:   "Card 2",
						Labels: trello.Labels{},
					},
				}
			},
			expected: aurora.Sprintf(`Name      ID    Labels
----      --    ------
Card 1    1     
Card 2    2     
`),
		},
		"no card": {
			given: func() trello.Cards {
				return trello.Cards{}
			},
			expected: `Name    ID    Labels
----    --    ------
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer()
			actual := r.RenderCards(tt.given())
			if actual != tt.expected {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected, actual)
			}
		})
	}
}

func TestInTable_RenderCard(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Card
		expected string
	}{
		"card with labe": {
			given: trello.Card{
				ID:   "1",
				Name: "Card 1",
				Description: `# Card title

> some context

Here are some markdown contents`,
				Labels: trello.Labels{
					trello.Label{
						ID:    "10",
						Name:  "Label 10",
						Color: "green",
					},
					trello.Label{
						ID:    "11",
						Name:  "Label 11",
						Color: "yellow",
					},
				},
			},
			expected: fmt.Sprintf(`ID:             1
Name:           Card 1
Labels:         %s%s
Description:    # Card title
                
                > some context
                
                Here are some markdown contents
`,
				aurora.BrightGreen("██ Label 10 "),
				aurora.BrightYellow("██ Label 11 "),
			),
		},
		"card without label": {
			given: trello.Card{
				ID:   "2",
				Name: "Card 2",
				Description: `# Card title

> some context

Here are some markdown contents`,
				Labels: trello.Labels{},
			},
			expected: `ID:             2
Name:           Card 2
Labels:         
Description:    # Card title
                
                > some context
                
                Here are some markdown contents
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer()
			actual := r.RenderCard(tt.given)
			if actual != tt.expected {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected, actual)
			}
		})
	}
}
