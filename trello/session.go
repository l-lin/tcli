package trello

type Session struct {
	Board *Board
	List  *List
	Card  *Card
}
