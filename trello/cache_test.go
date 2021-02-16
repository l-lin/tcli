package trello

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestCacheInMemory_GetBoards(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := NewMockRepository(ctrl)
	expected := Boards{
		{ID: "board 1", Name: "board"},
		{ID: "board 2", Name: "another board"},
	}
	r.EXPECT().
		GetBoards().
		Return(expected, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual1, err1 := cr.GetBoards()
	actual2, err2 := cr.GetBoards()

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
		GetBoards().
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

func TestCacheInMemory_GetLists(t *testing.T) {
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
		GetLists(idBoard).
		Return(expected, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual1, err1 := cr.GetLists(idBoard)
	actual2, err2 := cr.GetLists(idBoard)

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
		GetLists(idBoard).
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

func TestCacheInMemory_GetCards(t *testing.T) {
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
		GetCards(idList).
		Return(expected, nil).
		Times(1)
	cr := NewCacheInMemory(r)

	// WHEN
	actual1, err1 := cr.GetCards(idList)
	actual2, err2 := cr.GetCards(idList)

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
		GetCards(idList).
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
