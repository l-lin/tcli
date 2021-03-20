package renderer

import (
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestInTable_RenderBoards(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Boards
		expected string
	}{
		"two boards": {
			given: trello.Boards{
				{
					ID:               "1",
					Name:             "Board 1",
					ShortURL:         "https://trello.com/b/azerty",
					DateLastActivity: "2021-02-04T14:19:25.229Z",
				},
				{
					ID:               "2",
					Name:             "Board 2",
					ShortURL:         "https://trello.com/b/popo",
					DateLastActivity: "2021-02-08T21:02:58.117Z",
				},
			},
			expected: `Name       Last activity date
----       ------------------
Board 1    2021-02-04T14:19:25.229Z
Board 2    2021-02-08T21:02:58.117Z
`,
		},
		"no board": {
			given: trello.Boards{},
			expected: `Name    Last activity date
----    ------------------
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer(PlainLabel{}, PlainDescription{})
			actual := r.RenderBoards(tt.given)
			if actual != tt.expected {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected, actual)
			}
		})
	}
}

func TestInTable_RenderBoard(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Board
		expected string
	}{
		"existing board": {
			given: trello.Board{
				ID:               "1",
				Name:             "Board 1",
				ShortLink:        "azerty",
				ShortURL:         "https://trello.com/b/azerty",
				DateLastActivity: "2021-02-04T14:19:25.229Z",
			},
			expected: `ID:                    1
Short link:            azerty
Short URL:             https://trello.com/b/azerty
Name:                  Board 1
Last activity date:    2021-02-04T14:19:25.229Z
`,
		},
		"no board": {
			given: trello.Board{},
			expected: `ID:                    
Short link:            
Short URL:             
Name:                  
Last activity date:    
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer(PlainLabel{}, PlainDescription{})
			actual := r.RenderBoard(tt.given)
			if actual != tt.expected {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected, actual)
			}
		})
	}
}

func TestInTable_RenderLists(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Lists
		expected string
	}{
		"two lists": {
			given: trello.Lists{
				{ID: "1", Name: "List 1"},
				{ID: "2", Name: "List 2"},
			},
			expected: `Name
----
List 1
List 2
`,
		},
		"no list": {
			given: trello.Lists{},
			expected: `Name
----
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer(PlainLabel{}, PlainDescription{})
			actual := r.RenderLists(tt.given)
			if actual != tt.expected {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected, actual)
			}
		})
	}
}

func TestInTable_RenderList(t *testing.T) {
	var tests = map[string]struct {
		given    trello.List
		expected string
	}{
		"existing list": {
			given: trello.List{ID: "1", Name: "List 1"},
			expected: `ID:      1
Name:    List 1
`,
		},
		"no list": {
			given: trello.List{},
			expected: `ID:      
Name:    
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer(PlainLabel{}, PlainDescription{})
			actual := r.RenderList(tt.given)
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
						ID:       "1",
						ShortURL: "https://trello.com/c/abcd1",
						Name:     "Card 1",
						Pos:      10,
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
						ID:       "2",
						ShortURL: "https://trello.com/c/abcd2",
						Name:     "Card 2",
						Pos:      20,
						Labels: trello.Labels{
							trello.Label{
								ID:    "20",
								Name:  "Label 20",
								Color: "black",
							},
							trello.Label{
								ID:    "21",
								Color: "sky",
							},
						},
					},
				}
			},
			expected: `Name      Position    Labels
----      --------    ------
Card 1    10          Label 10 Label 11 
Card 2    20          Label 20 sky 
`,
		},
		"two cards without label": {
			given: func() trello.Cards {
				return trello.Cards{
					trello.Card{
						ID:       "1",
						ShortURL: "https://trello.com/c/abcd1",
						Name:     "Card 1",
						Pos:      10,
						Labels:   trello.Labels{},
					},
					trello.Card{
						ID:       "2",
						ShortURL: "https://trello.com/c/abcd2",
						Name:     "Card 2",
						Pos:      20,
						Labels:   trello.Labels{},
					},
				}
			},
			expected: `Name      Position    Labels
----      --------    ------
Card 1    10          
Card 2    20          
`,
		},
		"display cards by position order": {
			given: func() trello.Cards {
				return trello.Cards{
					trello.Card{
						ID:       "1",
						ShortURL: "https://trello.com/c/abcd1",
						Name:     "Card 1",
						Labels:   trello.Labels{},
						Pos:      10,
					},
					trello.Card{
						ID:       "2",
						ShortURL: "https://trello.com/c/abcd2",
						Name:     "Card 2",
						Labels:   trello.Labels{},
						Pos:      1,
					},
				}
			},
			expected: `Name      Position    Labels
----      --------    ------
Card 2    1           
Card 1    10          
`,
		},
		"no card": {
			given: func() trello.Cards {
				return trello.Cards{}
			},
			expected: `Name    Position    Labels
----    --------    ------
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer(PlainLabel{}, PlainDescription{})
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
		"card with labels": {
			given: trello.Card{
				ID:        "1",
				Name:      "Card 1",
				Pos:       1234,
				ShortLink: "abcd1234",
				ShortURL:  "https://trello.com/c/abcd1234",
				Desc: `# Card title

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
			expected: `ID:             1
Name:           Card 1
Position:       1234
Short link:     abcd1234
Short URL:      https://trello.com/c/abcd1234
Labels:         Label 10 Label 11 
Description:    
# Card title

> some context

Here are some markdown contents
`,
		},
		"card without label": {
			given: trello.Card{
				ID:        "2",
				Name:      "Card 2",
				Pos:       1234,
				ShortLink: "abcd1234",
				ShortURL:  "https://trello.com/c/abcd1234",
				Desc: `# Card title

> some context

Here are some markdown contents`,
				Labels: trello.Labels{},
			},
			expected: `ID:             2
Name:           Card 2
Position:       1234
Short link:     abcd1234
Short URL:      https://trello.com/c/abcd1234
Labels:         
Description:    
# Card title

> some context

Here are some markdown contents
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer(PlainLabel{}, PlainDescription{})
			actual := r.RenderCard(tt.given)
			if actual != tt.expected {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected, actual)
			}
		})
	}
}

func TestInTable_RenderComments(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Comments
		expected string
	}{
		"3 comments": {
			given: trello.Comments{
				{
					ID:   "comment 1",
					Date: "2014-11-12T11:45:26.371Z",
					MemberCreator: trello.CommentMemberCreator{
						Username: "user 1",
					},
					Data: trello.CommentData{
						Text: `# Comment title

> some context

Here are some markdown content.`,
					},
				},
				{
					ID:   "comment 2",
					Date: "2014-11-13T11:45:26.371Z",
					MemberCreator: trello.CommentMemberCreator{
						Username: "user 2",
					},
					Data: trello.CommentData{
						Text: "NO",
					},
				},
				{
					ID:   "comment 3",
					Date: "2014-11-14T11:45:26.371Z",
					MemberCreator: trello.CommentMemberCreator{
						Username: "user 1",
					},
					Data: trello.CommentData{
						Text: "But why?",
					},
				},
			},
			expected: `
user 1 @ 2014-11-14T11:45:26.371Z [comment 3]
---------------------------------------------
But why?

user 2 @ 2014-11-13T11:45:26.371Z [comment 2]
---------------------------------------------
NO

user 1 @ 2014-11-12T11:45:26.371Z [comment 1]
---------------------------------------------
# Comment title

> some context

Here are some markdown content.
`,
		},
		"no comment": {
			given:    trello.Comments{},
			expected: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer(PlainLabel{}, PlainDescription{})
			actual := r.RenderComments(tt.given)
			if actual != tt.expected {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected, actual)
			}
		})
	}
}

func TestInTable_RenderComment(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Comment
		expected string
	}{
		"long text": {
			given: trello.Comment{
				ID:   "comment 1",
				Date: "2014-11-12T11:45:26.371Z",
				MemberCreator: trello.CommentMemberCreator{
					Username: "user 1",
				},
				Data: trello.CommentData{
					Text: `# Comment title

> some context

Here are some markdown content.`,
				},
			},
			expected: `
user 1 @ 2014-11-12T11:45:26.371Z [comment 1]
---------------------------------------------
# Comment title

> some context

Here are some markdown content.
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewInTableRenderer(PlainLabel{}, PlainDescription{})
			actual := r.RenderComment(tt.given)
			if actual != tt.expected {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected, actual)
			}
		})
	}
}
