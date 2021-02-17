package trello

import (
	"github.com/l-lin/tcli/conf"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpRepository_GetBoards(t *testing.T) {
	type given struct {
		tsFn func() *httptest.Server
	}

	var tests = map[string]struct {
		given given
		test  func(actual Boards, err error)
	}{
		"happy path": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`
[{
  "id": "board 1",
  "name": "board"
}, {
  "id": "board 2",
  "name": "another board"
}]`))
					}))
				},
			},
			test: func(actual Boards, err error) {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
					t.FailNow()
				}
				if actual == nil {
					t.Error("expected not nil boards")
					t.FailNow()
				}
				expected := Boards{
					{ID: "board 1", Name: "board"},
					{ID: "board 2", Name: "another board"},
				}
				if len(expected) != len(actual) {
					t.Errorf("expected %v, actual %v", expected, actual)
					t.FailNow()
				}
				for i := 0; i < len(expected); i++ {
					if actual[i] != expected[i] {
						t.Errorf("%d: expected %v, actual %v", i, expected[i], actual[i])
					}
				}
			},
		},
		"server error": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
					}))
				},
			},
			test: func(actual Boards, err error) {
				if err == nil {
					t.Error("expected error")
				}
				if actual != nil {
					t.Error("expected nil boards")
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := tt.given.tsFn()
			repository := NewHttpRepository(conf.Conf{
				Trello: conf.Trello{
					BaseURL: ts.URL,
				},
			}, false)
			tt.test(repository.GetBoards())
		})
	}
}

func TestHttpRepository_FindBoard(t *testing.T) {
	type given struct {
		tsFn func() *httptest.Server
	}

	var tests = map[string]struct {
		given given
		test  func(actual *Board, err error)
	}{
		"happy path": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`
[{
  "id": "board 1",
  "name": "board"
}, {
  "id": "board 2",
  "name": "another board"
}]`))
					}))
				},
			},
			test: func(actual *Board, err error) {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
					t.FailNow()
				}
				if actual == nil {
					t.Error("expected not nil boards")
					t.FailNow()
				}
				expected := &Board{ID: "board 1", Name: "board"}
				if *expected != *actual {
					t.Errorf("expected %v, actual %v", expected, actual)
				}
			},
		},
		"server error": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
					}))
				},
			},
			test: func(actual *Board, err error) {
				if err == nil {
					t.Error("expected error")
				}
				if actual != nil {
					t.Error("expected nil board")
				}
			},
		},
		"no board found": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`
[{
  "id": "board 2",
  "name": "another board"
}]`))
					}))
				},
			},
			test: func(actual *Board, err error) {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
					t.FailNow()
				}
				if actual != nil {
					t.Error("expected nil board")
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := tt.given.tsFn()
			repository := NewHttpRepository(conf.Conf{
				Trello: conf.Trello{
					BaseURL: ts.URL,
				},
			}, false)
			tt.test(repository.FindBoard("board"))
		})
	}
}

func TestHttpRepository_GetLists(t *testing.T) {
	type given struct {
		tsFn func() *httptest.Server
	}

	var tests = map[string]struct {
		given given
		test  func(actual Lists, err error)
	}{
		"happy path": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`
[{
  "id": "list 1",
  "name": "list"
}, {
  "id": "list 2",
  "name": "another list"
}]`))
					}))
				},
			},
			test: func(actual Lists, err error) {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
					t.FailNow()
				}
				if actual == nil {
					t.Error("expected not nil boards")
					t.FailNow()
				}
				expected := Lists{
					{ID: "list 1", Name: "list"},
					{ID: "list 2", Name: "another list"},
				}
				if len(expected) != len(actual) {
					t.Errorf("expected %v, actual %v", expected, actual)
					t.FailNow()
				}
				for i := 0; i < len(expected); i++ {
					if actual[i] != expected[i] {
						t.Errorf("%d: expected %v, actual %v", i, expected[i], actual[i])
					}
				}
			},
		},
		"server error": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
					}))
				},
			},
			test: func(actual Lists, err error) {
				if err == nil {
					t.Error("expected error")
				}
				if actual != nil {
					t.Error("expected nil boards")
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := tt.given.tsFn()
			repository := NewHttpRepository(conf.Conf{
				Trello: conf.Trello{
					BaseURL: ts.URL,
				},
			}, false)
			tt.test(repository.GetLists("board 1"))
		})
	}
}

func TestHttpRepository_FindList(t *testing.T) {
	type given struct {
		tsFn func() *httptest.Server
	}

	var tests = map[string]struct {
		given given
		test  func(actual *List, err error)
	}{
		"happy path": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`
[{
  "id": "list 1",
  "name": "list"
}, {
  "id": "list 2",
  "name": "another list"
}]`))
					}))
				},
			},
			test: func(actual *List, err error) {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
					t.FailNow()
				}
				if actual == nil {
					t.Error("expected not nil boards")
					t.FailNow()
				}
				expected := &List{ID: "list 1", Name: "list"}
				if *expected != *actual {
					t.Errorf("expected %v, actual %v", expected, actual)
				}
			},
		},
		"server error": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
					}))
				},
			},
			test: func(actual *List, err error) {
				if err == nil {
					t.Error("expected error")
				}
				if actual != nil {
					t.Error("expected nil board")
				}
			},
		},
		"no board found": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`
[{
  "id": "list 2",
  "name": "another list"
}]`))
					}))
				},
			},
			test: func(actual *List, err error) {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
					t.FailNow()
				}
				if actual != nil {
					t.Error("expected nil list")
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := tt.given.tsFn()
			repository := NewHttpRepository(conf.Conf{
				Trello: conf.Trello{
					BaseURL: ts.URL,
				},
			}, false)
			tt.test(repository.FindList("board 1", "list"))
		})
	}
}

func TestHttpRepository_GetCards(t *testing.T) {
	type given struct {
		tsFn func() *httptest.Server
	}

	var tests = map[string]struct {
		given given
		test  func(actual Cards, err error)
	}{
		"happy path": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`
[{
  "id": "card 1",
  "name": "card"
}, {
  "id": "card 2",
  "name": "another card"
}]`))
					}))
				},
			},
			test: func(actual Cards, err error) {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
					t.FailNow()
				}
				if actual == nil {
					t.Error("expected not nil boards")
					t.FailNow()
				}
				expected := Cards{
					{ID: "card 1", Name: "card"},
					{ID: "card 2", Name: "another card"},
				}
				if len(expected) != len(actual) {
					t.Errorf("expected %v, actual %v", expected, actual)
					t.FailNow()
				}
				for i := 0; i < len(expected); i++ {
					if actual[i].ID != expected[i].ID || actual[i].Name != expected[i].Name {
						t.Errorf("%d: expected %v, actual %v", i, expected[i], actual[i])
					}
				}
			},
		},
		"server error": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
					}))
				},
			},
			test: func(actual Cards, err error) {
				if err == nil {
					t.Error("expected error")
				}
				if actual != nil {
					t.Error("expected nil boards")
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := tt.given.tsFn()
			repository := NewHttpRepository(conf.Conf{
				Trello: conf.Trello{
					BaseURL: ts.URL,
				},
			}, false)
			tt.test(repository.GetCards("list 1"))
		})
	}
}

func TestHttpRepository_FindCard(t *testing.T) {
	type given struct {
		tsFn func() *httptest.Server
	}

	var tests = map[string]struct {
		given given
		test  func(actual *Card, err error)
	}{
		"happy path": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`
[{
  "id": "card 1",
  "name": "card"
}, {
  "id": "card 2",
  "name": "another card"
}]`))
					}))
				},
			},
			test: func(actual *Card, err error) {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
					t.FailNow()
				}
				if actual == nil {
					t.Error("expected not nil boards")
					t.FailNow()
				}
				expected := &Card{ID: "card 1", Name: "card"}
				if expected.ID != actual.ID || expected.Name != actual.Name {
					t.Errorf("expected %v, actual %v", expected, actual)
				}
			},
		},
		"server error": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
					}))
				},
			},
			test: func(actual *Card, err error) {
				if err == nil {
					t.Error("expected error")
				}
				if actual != nil {
					t.Error("expected nil card")
				}
			},
		},
		"no board found": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`
[{
  "id": "card 2",
  "name": "another card"
}]`))
					}))
				},
			},
			test: func(actual *Card, err error) {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
					t.FailNow()
				}
				if actual != nil {
					t.Error("expected nil board")
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := tt.given.tsFn()
			repository := NewHttpRepository(conf.Conf{
				Trello: conf.Trello{
					BaseURL: ts.URL,
				},
			}, false)
			tt.test(repository.FindCard("list 1", "card"))
		})
	}
}

func TestHttpRepository_UpdateCard(t *testing.T) {
	type given struct {
		tsFn func() *httptest.Server
	}
	type expected struct {
		hasError bool
		card     *Card
	}
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"happy path": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write([]byte(`
{
  "id": "card 1",
  "name": "updated card"
}`))
					}))
				},
			},
			expected: expected{
				hasError: false,
				card:     &Card{ID: "card 1", Name: "updated card"},
			},
		},
		"server error": {
			given: given{
				tsFn: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
					}))
				},
			},
			expected: expected{
				hasError: true,
				card:     nil,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := tt.given.tsFn()
			repository := NewHttpRepository(conf.Conf{
				Trello: conf.Trello{
					BaseURL: ts.URL,
				},
			}, false)
			actual, actualErr := repository.UpdateCard(UpdateCard{ID: "card 1", Name: "card"})
			if tt.expected.hasError && actualErr == nil || !tt.expected.hasError && actualErr != nil {
				t.Errorf("expected err %v, actual err %v", tt.expected.hasError, actualErr != nil)
			}
			if tt.expected.card != nil && actual == nil || tt.expected.card == nil && actual != nil {
				t.Errorf("expected %v, actual %v", tt.expected.card, actual)
			}
			if tt.expected.card != nil {
				if tt.expected.card.ID != actual.ID && tt.expected.card.Name != actual.Name {
					t.Errorf("expected %v, actual %v", tt.expected.card, actual)
				}
			}
		})
	}
}
