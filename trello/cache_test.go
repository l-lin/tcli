package trello

import (
	"github.com/golang/mock/gomock"
	"sort"
	"testing"
)

func TestCacheInMemory_FindBoards(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := NewMockRepository(ctrl)
	expected := Boards{
		{ID: "board 1", Name: "board"},
		{ID: "board 2", Name: "another board"},
	}
	r.EXPECT().
		FindBoards().
		Return(expected, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual1, err1 := cr.FindBoards()
	actual2, err2 := cr.FindBoards()

	// THEN
	if err1 != nil || err2 != nil {
		t.Error("expected no error")
	}
	if len(actual1) != len(expected) || len(actual2) != len(expected) {
		t.Errorf("expected %v, actual1 %v, actual2 %v", expected, actual1, actual2)
		t.FailNow()
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != actual1[i] || expected[i] != actual2[i] {
			t.Errorf("%d: expected %v, actual1 %v, actual2 %v", i, expected[i], actual1[i], actual2[i])
		}
	}
}

func TestCacheInMemory_FindBoard(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := NewMockRepository(ctrl)
	boards := Boards{
		{ID: "board 1", Name: "board"},
		{ID: "board 2", Name: "another board"},
	}
	r.EXPECT().
		FindBoards().
		Return(boards, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual, err := cr.FindBoard("board")

	// THEN
	if err != nil {
		t.Error("expected no error")
	}
	expected := boards[0]
	if actual == nil || *actual != expected {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func TestCacheInMemory_FindLists(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := NewMockRepository(ctrl)
	expected := Lists{
		{ID: "list 1", Name: "list"},
		{ID: "list 2", Name: "another list"},
	}
	idBoard := "board 1"
	r.EXPECT().
		FindLists(idBoard).
		Return(expected, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual1, err1 := cr.FindLists(idBoard)
	actual2, err2 := cr.FindLists(idBoard)

	// THEN
	if err1 != nil || err2 != nil {
		t.Error("expected no error")
	}
	if len(actual1) != len(expected) || len(actual2) != len(expected) {
		t.Errorf("expected %v, actual1 %v, actual2 %v", expected, actual1, actual2)
		t.FailNow()
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != actual1[i] || expected[i] != actual2[i] {
			t.Errorf("%d: expected %v, actual1 %v, actual2 %v", i, expected[i], actual1[i], actual2[i])
		}
	}
}

func TestCacheInMemory_FindList(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := NewMockRepository(ctrl)
	lists := Lists{
		{ID: "list 1", Name: "list"},
		{ID: "list 2", Name: "another list"},
	}
	idBoard := "board 1"
	r.EXPECT().
		FindLists(idBoard).
		Return(lists, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual, err := cr.FindList(idBoard, "list")

	// THEN
	if err != nil {
		t.Error("expected no error")
	}
	expected := lists[0]
	if actual == nil || *actual != expected {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func TestCacheInMemory_FindCards(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := NewMockRepository(ctrl)
	expected := Cards{
		{ID: "card 1", Name: "card"},
		{ID: "card 2", Name: "another card"},
	}
	idList := "list 1"
	r.EXPECT().
		FindCards(idList).
		Return(expected, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual1, err1 := cr.FindCards(idList)
	actual2, err2 := cr.FindCards(idList)

	// THEN
	if err1 != nil || err2 != nil {
		t.Error("expected no error")
	}
	if len(actual1) != len(expected) || len(actual2) != len(expected) {
		t.Errorf("expected %v, actual1 %v, actual2 %v", expected, actual1, actual2)
		t.FailNow()
	}
	for i := 0; i < len(expected); i++ {
		if expected[i].ID != actual1[i].ID || expected[i].ID != actual2[i].ID {
			t.Errorf("%d: expected %v, actual1 %v, actual2 %v", i, expected[i], actual1[i], actual2[i])
		}
	}
}

func TestCacheInMemory_FindCard(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := NewMockRepository(ctrl)
	cards := Cards{
		{ID: "card 1", Name: "card"},
		{ID: "card 2", Name: "another card"},
	}
	idList := "list 1"
	r.EXPECT().
		FindCards(idList).
		Return(cards, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual, err := cr.FindCard(idList, "card")

	// THEN
	if err != nil {
		t.Error("expected no error")
	}
	expected := cards[0]
	if actual == nil || actual.ID != expected.ID {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func TestCacheInMemory_CreateCard(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := NewMockRepository(ctrl)
	expected := Card{ID: "card 1", Name: "card", IDList: "list 1"}
	r.EXPECT().
		CreateCard(NewCreateCard(expected)).
		Return(&expected, nil)
	cr := &CacheInMemory{
		r:        r,
		mapLists: map[string]Lists{},
		mapCards: map[string]Cards{
			expected.IDList: {
				{ID: "card 2", Name: "another card", IDList: expected.IDList},
			},
		},
	}

	// WHEN
	actual, err := cr.CreateCard(NewCreateCard(expected))

	// THEN
	if err != nil {
		t.Error("expected no error")
	}
	if actual == nil {
		t.Error("expected no nil card returned")
		t.FailNow()
	}
	cardInCache := cr.mapCards[expected.IDList][1]
	if cardInCache.ID != expected.ID && cardInCache.Name != expected.Name &&
		actual.ID != expected.ID && actual.Name != expected.Name {
		t.Errorf("expected %v, actual %v, card in cache %v", expected, actual, cardInCache)
	}
}

func TestCacheInMemory_UpdateCard(t *testing.T) {
	t.Run("update card", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		cardToUpdate := Card{ID: "card 1", Name: "card", IDList: "list 1"}
		expected := Card{ID: "card 1", Name: "updated card", IDList: "list 1"}
		r.EXPECT().
			UpdateCard(NewUpdateCard(cardToUpdate)).
			Return(&expected, nil)
		cr := &CacheInMemory{
			r:        r,
			mapLists: map[string]Lists{},
			mapCards: map[string]Cards{
				cardToUpdate.IDList: {
					{ID: "card 2", Name: "another card", IDList: cardToUpdate.IDList},
					cardToUpdate,
				},
			},
		}

		// WHEN
		actual, err := cr.UpdateCard(NewUpdateCard(cardToUpdate))

		// THEN
		if err != nil {
			t.Error("expected no error")
		}
		if actual == nil {
			t.Error("expected no nil card returned")
			t.FailNow()
		}
		cardInCache := cr.mapCards[cardToUpdate.IDList][1]
		if cardInCache.ID != expected.ID && cardInCache.Name != expected.Name &&
			actual.ID != expected.ID && actual.Name != expected.Name {
			t.Errorf("expected %v, actual %v, card in cache %v", expected, actual, cardInCache)
		}
	})

	t.Run("update card to closed", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		cardToUpdate := Card{ID: "card 1", Name: "card", IDList: "list 1"}
		expected := Card{ID: "card 1", Name: "updated card", IDList: "list 1", Closed: true}
		updateCard := NewUpdateCard(cardToUpdate)
		updateCard.Closed = true
		r.EXPECT().
			UpdateCard(updateCard).
			Return(&expected, nil)
		cr := &CacheInMemory{
			r:        r,
			mapLists: map[string]Lists{},
			mapCards: map[string]Cards{
				cardToUpdate.IDList: {
					{ID: "card 2", Name: "another card", IDList: cardToUpdate.IDList},
					cardToUpdate,
				},
			},
		}

		// WHEN
		actual, err := cr.UpdateCard(updateCard)

		// THEN
		if err != nil {
			t.Error("expected no error")
		}
		if actual == nil {
			t.Error("expected no nil card returned")
			t.FailNow()
		}
		if len(cr.mapCards[cardToUpdate.IDList]) != 1 {
			t.Errorf("expected updated card removed from cache, got %v", cr.mapCards[cardToUpdate.IDList])
		}
	})
}

func TestCacheInMemory_findCardIndex(t *testing.T) {
	cr := &CacheInMemory{
		mapCards: map[string]Cards{
			"list 1": {
				{ID: "card 1"},
				{ID: "card 2"},
				{ID: "card 3"},
			},
		},
	}

	type given struct {
		idList string
		idCard string
	}
	var tests = map[string]struct {
		given    given
		expected int
	}{
		"card found on index 1": {
			given: given{
				idList: "list 1",
				idCard: "card 2",
			},
			expected: 1,
		},
		"card not found list 1": {
			given: given{
				idList: "list 1",
				idCard: "card unknown",
			},
			expected: -1,
		},
		"list not found": {
			given: given{
				idList: "list unknown",
			},
			expected: -1,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := cr.findCardIndex(tt.given.idList, tt.given.idCard)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestCacheInMemory_removeCard(t *testing.T) {
	type given struct {
		mapCards  map[string]Cards
		idList    string
		cardIndex int
	}
	type expected struct {
		mapCards map[string]Cards
	}
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"remove card at index 1": {
			given: given{
				mapCards: map[string]Cards{
					"list 1": {
						{ID: "card 1"},
						{ID: "card 2"},
						{ID: "card 3"},
					},
				},
				idList:    "list 1",
				cardIndex: 1,
			},
			expected: expected{
				mapCards: map[string]Cards{
					"list 1": {
						{ID: "card 1"},
						{ID: "card 3"},
					},
				},
			},
		},
		"invalid index": {
			given: given{
				mapCards: map[string]Cards{
					"list 1": {
						{ID: "card 1"},
						{ID: "card 2"},
						{ID: "card 3"},
					},
				},
				idList:    "list 1",
				cardIndex: -1,
			},
			expected: expected{
				mapCards: map[string]Cards{
					"list 1": {
						{ID: "card 1"},
						{ID: "card 2"},
						{ID: "card 3"},
					},
				},
			},
		},
		"index greater than mapCards": {
			given: given{
				mapCards: map[string]Cards{
					"list 1": {
						{ID: "card 1"},
						{ID: "card 2"},
						{ID: "card 3"},
					},
				},
				idList:    "list 1",
				cardIndex: 3,
			},
			expected: expected{
				mapCards: map[string]Cards{
					"list 1": {
						{ID: "card 1"},
						{ID: "card 2"},
						{ID: "card 3"},
					},
				},
			},
		},
		"idList not found in cache": {
			given: given{
				mapCards: map[string]Cards{
					"list 1": {
						{ID: "card 1"},
						{ID: "card 2"},
						{ID: "card 3"},
					},
				},
				idList:    "list unknown",
				cardIndex: 1,
			},
			expected: expected{
				mapCards: map[string]Cards{
					"list 1": {
						{ID: "card 1"},
						{ID: "card 2"},
						{ID: "card 3"},
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cr := &CacheInMemory{mapCards: tt.given.mapCards}
			cr.removeCard(tt.given.idList, tt.given.cardIndex)
			actual := cr.mapCards
			if len(actual) != len(tt.expected.mapCards) {
				t.Errorf("expected %v, actual %v", tt.expected.mapCards, actual)
				t.FailNow()
			}
			for idList, cards := range tt.expected.mapCards {
				if len(actual[idList]) != len(cards) {
					t.Errorf("expected %v, actual %v", cards, actual[idList])
					t.FailNow()
				}
				sort.Slice(actual[idList], func(i, j int) bool {
					return actual[idList][i].ID < actual[idList][j].ID
				})
				for i := 0; i < len(cards); i++ {
					if actual[idList][i].ID != cards[i].ID {
						t.Errorf("%d: expected %v, actual %v", i, cards[i], actual[idList][i])
					}
				}
			}
		})
	}
}
