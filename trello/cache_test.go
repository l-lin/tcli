package trello

import (
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
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
	if !reflect.DeepEqual(expected, actual1) || !reflect.DeepEqual(expected, actual2) {
		t.Errorf("expected %v, actual1 %v, actual2 %v", expected, actual1, actual2)
	}
}

func TestCacheInMemory_FindBoard(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
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
	})
	t.Run("no board found", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		boards := Boards{
			{ID: "board 1", Name: "board"},
			{ID: "board 2", Name: "another board"},
		}
		r := NewMockRepository(ctrl)
		r.EXPECT().
			FindBoards().
			Return(boards, nil)
		cr := NewCacheInMemory(r)

		// WHEN
		actual, err := cr.FindBoard("unknown board")

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil returned value, got %v", actual)
		}
	})
	t.Run("error when finding boards", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		r.EXPECT().
			FindBoards().
			Return(nil, errors.New("unexpected error"))
		cr := NewCacheInMemory(r)

		// WHEN
		actual, err := cr.FindBoard("board")

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil returned value, got %v", actual)
		}
	})
}

func TestCacheInMemory_FindLabels(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := NewMockRepository(ctrl)
	idBoard := "board 1"
	expected := Labels{
		{ID: "label 1", IDBoard: idBoard, Color: "red", Name: "label name 1"},
		{ID: "label 2", IDBoard: idBoard, Color: "sky", Name: "label name 2"},
	}
	r.EXPECT().
		FindLabels(idBoard).
		Return(expected, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual1, err1 := cr.FindLabels(idBoard)
	actual2, err2 := cr.FindLabels(idBoard)

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
	if !reflect.DeepEqual(expected, actual1) || !reflect.DeepEqual(expected, actual2) {
		t.Errorf("expected %v, actual1 %v, actual2 %v", expected, actual1, actual2)
	}
}

func TestCacheInMemory_FindList(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
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
	})
	t.Run("list not found", func(t *testing.T) {
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
		actual, err := cr.FindList(idBoard, "unknown list")

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil returned value, actual %v", actual)
		}
	})
	t.Run("error when finding lists", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		idBoard := "board 1"
		r.EXPECT().
			FindLists(idBoard).
			Return(nil, errors.New("unexpected error"))
		cr := NewCacheInMemory(r)

		// WHEN
		actual, err := cr.FindList(idBoard, "list")

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil returned value, actual %v", actual)
		}
	})
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
	if !reflect.DeepEqual(expected, actual1) || !reflect.DeepEqual(expected, actual2) {
		t.Errorf("expected %v, actual1 %v, actual2 %v", expected, actual1, actual2)
	}
}

func TestCacheInMemory_FindCard(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
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
	})
	t.Run("card not found", func(t *testing.T) {
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
		actual, err := cr.FindCard(idList, "unknown card")

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil returned value, actual %v", actual)
		}
	})
	t.Run("error when finding cards", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		idList := "list 1"
		r.EXPECT().
			FindCards(idList).
			Return(nil, errors.New("unexpected error"))
		cr := NewCacheInMemory(r)

		// WHEN
		actual, err := cr.FindCard(idList, "card")

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil returned value, actual %v", actual)
		}
	})
}

func TestCacheInMemory_CreateCard(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		expected := Card{ID: "card 1", Name: "card", IDList: "list 1"}
		r.EXPECT().
			CreateCard(NewCreateCard(expected)).
			Return(&expected, nil)
		cr := &CacheInMemory{
			r:                 r,
			mapListsByIDBoard: map[string]Lists{},
			mapCardsByIDList: map[string]Cards{
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
		cardInCache := cr.mapCardsByIDList[expected.IDList][1]
		if cardInCache.ID != expected.ID && cardInCache.Name != expected.Name &&
			actual.ID != expected.ID && actual.Name != expected.Name {
			t.Errorf("expected %v, actual %v, card in cache %v", expected, actual, cardInCache)
		}
	})
	t.Run("error when creating card", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		card := Card{ID: "card 1", Name: "card", IDList: "list 1"}
		r.EXPECT().
			CreateCard(NewCreateCard(card)).
			Return(nil, errors.New("unexpected error"))
		cr := &CacheInMemory{
			r:                 r,
			mapListsByIDBoard: map[string]Lists{},
			mapCardsByIDList: map[string]Cards{
				card.IDList: {
					{ID: "card 2", Name: "another card", IDList: card.IDList},
				},
			},
		}

		// WHEN
		actual, err := cr.CreateCard(NewCreateCard(card))

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil card returned, got %v", actual)
			t.FailNow()
		}
	})
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
			r:                 r,
			mapListsByIDBoard: map[string]Lists{},
			mapCardsByIDList: map[string]Cards{
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
		cardInCache := cr.mapCardsByIDList[cardToUpdate.IDList][1]
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
			r:                 r,
			mapListsByIDBoard: map[string]Lists{},
			mapCardsByIDList: map[string]Cards{
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
		if len(cr.mapCardsByIDList[cardToUpdate.IDList]) != 1 {
			t.Errorf("expected updated card removed from cache, got %v", cr.mapCardsByIDList[cardToUpdate.IDList])
		}
	})
	t.Run("error when updating card", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		cardToUpdate := Card{ID: "card 1", Name: "card", IDList: "list 1"}
		updateCard := NewUpdateCard(cardToUpdate)
		updateCard.Closed = true
		r.EXPECT().
			UpdateCard(updateCard).
			Return(nil, errors.New("unexpected error"))
		cr := &CacheInMemory{
			r:                 r,
			mapListsByIDBoard: map[string]Lists{},
			mapCardsByIDList: map[string]Cards{
				cardToUpdate.IDList: {
					{ID: "card 2", Name: "another card", IDList: cardToUpdate.IDList},
					cardToUpdate,
				},
			},
		}

		// WHEN
		actual, err := cr.UpdateCard(updateCard)

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil card returned, actual %v", actual)
			t.FailNow()
		}
		if len(cr.mapCardsByIDList[cardToUpdate.IDList]) != 2 {
			t.Errorf("expected cache not modified, got %v", cr.mapCardsByIDList[cardToUpdate.IDList])
		}
	})
}

func TestCacheInMemory_findCardIndex(t *testing.T) {
	cr := &CacheInMemory{
		mapCardsByIDList: map[string]Cards{
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
		"remove card at index 0": {
			given: given{
				mapCards: map[string]Cards{
					"list 1": {
						{ID: "card 1"},
						{ID: "card 2"},
						{ID: "card 3"},
					},
				},
				idList:    "list 1",
				cardIndex: 0,
			},
			expected: expected{
				mapCards: map[string]Cards{
					"list 1": {
						{ID: "card 2"},
						{ID: "card 3"},
					},
				},
			},
		},
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
		"index greater than mapCardsByIDList": {
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
			cr := &CacheInMemory{mapCardsByIDList: tt.given.mapCards}
			cr.removeCard(tt.given.idList, tt.given.cardIndex)
			actual := cr.mapCardsByIDList
			if !reflect.DeepEqual(tt.expected.mapCards, actual) {
				t.Errorf("expected %v, actual %v", tt.expected.mapCards, actual)
			}
		})
	}
}

func TestCacheInMemory_FindComments(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := NewMockRepository(ctrl)
	expected := Comments{
		{ID: "card 1", Date: "2021-02-02T16:18:41.228Z", Data: CommentData{Text: "text comment 1"}},
		{ID: "card 2", Date: "2021-02-02T18:17:41.228Z", Data: CommentData{Text: "text comment 2"}},
	}
	idCard := "card 1"
	r.EXPECT().
		FindComments(idCard).
		Return(expected, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual1, err1 := cr.FindComments(idCard)
	actual2, err2 := cr.FindComments(idCard)

	// THEN
	if err1 != nil || err2 != nil {
		t.Error("expected no error")
	}
	if !reflect.DeepEqual(expected, actual1) || !reflect.DeepEqual(expected, actual2) {
		t.Errorf("expected %v, actual1 %v, actual2 %v", expected, actual1, actual2)
	}
}

func TestCacheInMemory_FindComment(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		comments := Comments{
			{ID: "comment 1"},
			{ID: "comment 2"},
		}
		idCard := "card 1"
		r.EXPECT().
			FindComments(idCard).
			Return(comments, nil).
			Times(1)
		cr := NewCacheInMemory(r)

		// WHEN
		actual, err := cr.FindComment(idCard, "comment 2")

		// THEN
		if err != nil {
			t.Error("expected no error")
		}
		expected := &comments[1]
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	})
	t.Run("comment not found", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		comments := Comments{
			{ID: "comment 1"},
			{ID: "comment 2"},
		}
		idCard := "card 1"
		r.EXPECT().
			FindComments(idCard).
			Return(comments, nil).
			Times(1)
		cr := NewCacheInMemory(r)

		// WHEN
		actual, err := cr.FindComment(idCard, "unknown comment")

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil returned value, actual %v", actual)
		}
	})
	t.Run("error when finding cards", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		idList := "list 1"
		r.EXPECT().
			FindCards(idList).
			Return(nil, errors.New("unexpected error"))
		cr := NewCacheInMemory(r)

		// WHEN
		actual, err := cr.FindCard(idList, "card")

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil returned value, actual %v", actual)
		}
	})
}

func TestCacheInMemory_CreateComment(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		expected := &Comment{ID: "comment 2"}
		idCard := "card 1"
		createComment := CreateComment{IDCard: idCard, Text: "comment content"}
		r.EXPECT().
			CreateComment(createComment).
			Return(expected, nil)
		cr := &CacheInMemory{
			r: r,
			mapCommentsByIDCard: map[string]Comments{
				idCard: {
					{ID: "comment 1"},
				},
			},
		}

		// WHEN
		actual, err := cr.CreateComment(createComment)

		// THEN
		if err != nil {
			t.Error("expected no error")
		}
		if actual == nil {
			t.Error("expected no nil comment returned")
			t.FailNow()
		}
		commentInCache := cr.mapCommentsByIDCard[idCard][1]
		if !reflect.DeepEqual(expected, actual) || !reflect.DeepEqual(expected, &commentInCache) {
			t.Errorf("expected %v, actual %v, card in cache %v", expected, actual, &commentInCache)
		}
	})
	t.Run("error when creating comment", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		idCard := "card 1"
		createComment := CreateComment{IDCard: idCard, Text: "comment content"}
		r.EXPECT().
			CreateComment(createComment).
			Return(nil, errors.New("unexpected error"))
		cr := &CacheInMemory{
			r: r,
			mapCommentsByIDCard: map[string]Comments{
				idCard: {
					{ID: "comment 1"},
				},
			},
		}

		// WHEN
		actual, err := cr.CreateComment(createComment)

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil card returned, got %v", actual)
			t.FailNow()
		}
	})
}

func TestCacheInMemory_UpdateComment(t *testing.T) {
	t.Run("update comment existing in cache", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		updateComment := UpdateComment{ID: "comment 1", IDCard: "card 1", Text: "updated comment"}
		expected := &Comment{ID: updateComment.ID, Data: CommentData{Text: updateComment.Text}}
		r.EXPECT().
			UpdateComment(updateComment).
			Return(expected, nil)
		cr := &CacheInMemory{
			r: r,
			mapCommentsByIDCard: map[string]Comments{
				updateComment.IDCard: {
					{ID: "comment 2"},
					{ID: updateComment.ID, Data: CommentData{Text: "comment"}},
				},
			},
		}

		// WHEN
		actual, err := cr.UpdateComment(updateComment)

		// THEN
		if err != nil {
			t.Error("expected no error")
		}
		if actual == nil {
			t.Error("expected no nil card returned")
			t.FailNow()
		}
		commentInCache := cr.mapCommentsByIDCard[updateComment.IDCard][1]
		if !reflect.DeepEqual(expected, &commentInCache) || !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, actual %v, card in cache %v", expected, actual, &commentInCache)
		}
	})

	t.Run("update comment not existing in cache", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		updateComment := UpdateComment{ID: "comment 1", IDCard: "card 1", Text: "updated comment"}
		expected := &Comment{ID: updateComment.ID, Data: CommentData{Text: updateComment.Text}}
		r.EXPECT().
			UpdateComment(updateComment).
			Return(expected, nil)
		cr := &CacheInMemory{
			r: r,
			mapCommentsByIDCard: map[string]Comments{
				updateComment.IDCard: {
					{ID: "comment 2"},
				},
			},
		}

		// WHEN
		actual, err := cr.UpdateComment(updateComment)

		// THEN
		if err != nil {
			t.Error("expected no error")
		}
		if actual == nil {
			t.Error("expected no nil card returned")
			t.FailNow()
		}
		commentInCache := cr.mapCommentsByIDCard[updateComment.IDCard][1]
		if !reflect.DeepEqual(expected, &commentInCache) || !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, actual %v, card in cache %v", expected, actual, &commentInCache)
		}
	})

	t.Run("error when updating comment", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := NewMockRepository(ctrl)
		updateComment := UpdateComment{ID: "comment 1", IDCard: "card 1", Text: "updated comment"}
		r.EXPECT().
			UpdateComment(updateComment).
			Return(nil, errors.New("unexpected error"))
		cr := &CacheInMemory{
			r: r,
			mapCommentsByIDCard: map[string]Comments{
				updateComment.IDCard: {
					{ID: "comment 2"},
					{ID: updateComment.ID, Data: CommentData{Text: "comment"}},
				},
			},
		}

		// WHEN
		actual, err := cr.UpdateComment(updateComment)

		// THEN
		if err == nil {
			t.Error("expected error")
		}
		if actual != nil {
			t.Errorf("expected nil comment returned, actual %v", actual)
			t.FailNow()
		}
		if len(cr.mapCommentsByIDCard[updateComment.IDCard]) != 2 {
			t.Errorf("expected cache not modified, got %v", cr.mapCardsByIDList[updateComment.IDCard])
		}
	})
}
