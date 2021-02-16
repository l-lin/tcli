// Code generated by MockGen. DO NOT EDIT.
// Source: renderer.go

// Package renderer is a generated GoMock package.
package renderer

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	trello "github.com/l-lin/tcli/trello"
)

// MockRenderer is a mock of Renderer interface.
type MockRenderer struct {
	ctrl     *gomock.Controller
	recorder *MockRendererMockRecorder
}

// MockRendererMockRecorder is the mock recorder for MockRenderer.
type MockRendererMockRecorder struct {
	mock *MockRenderer
}

// NewMockRenderer creates a new mock instance.
func NewMockRenderer(ctrl *gomock.Controller) *MockRenderer {
	mock := &MockRenderer{ctrl: ctrl}
	mock.recorder = &MockRendererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRenderer) EXPECT() *MockRendererMockRecorder {
	return m.recorder
}

// RenderBoard mocks base method.
func (m *MockRenderer) RenderBoard(arg0 trello.Board) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenderBoard", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// RenderBoard indicates an expected call of RenderBoard.
func (mr *MockRendererMockRecorder) RenderBoard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderBoard", reflect.TypeOf((*MockRenderer)(nil).RenderBoard), arg0)
}

// RenderBoards mocks base method.
func (m *MockRenderer) RenderBoards(arg0 trello.Boards) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenderBoards", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// RenderBoards indicates an expected call of RenderBoards.
func (mr *MockRendererMockRecorder) RenderBoards(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderBoards", reflect.TypeOf((*MockRenderer)(nil).RenderBoards), arg0)
}

// RenderCard mocks base method.
func (m *MockRenderer) RenderCard(arg0 trello.Card) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenderCard", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// RenderCard indicates an expected call of RenderCard.
func (mr *MockRendererMockRecorder) RenderCard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderCard", reflect.TypeOf((*MockRenderer)(nil).RenderCard), arg0)
}

// RenderCards mocks base method.
func (m *MockRenderer) RenderCards(arg0 trello.Cards) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenderCards", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// RenderCards indicates an expected call of RenderCards.
func (mr *MockRendererMockRecorder) RenderCards(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderCards", reflect.TypeOf((*MockRenderer)(nil).RenderCards), arg0)
}

// RenderList mocks base method.
func (m *MockRenderer) RenderList(arg0 trello.List) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenderList", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// RenderList indicates an expected call of RenderList.
func (mr *MockRendererMockRecorder) RenderList(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderList", reflect.TypeOf((*MockRenderer)(nil).RenderList), arg0)
}

// RenderLists mocks base method.
func (m *MockRenderer) RenderLists(arg0 trello.Lists) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenderLists", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// RenderLists indicates an expected call of RenderLists.
func (mr *MockRendererMockRecorder) RenderLists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderLists", reflect.TypeOf((*MockRenderer)(nil).RenderLists), arg0)
}