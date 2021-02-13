package trello

type Repository interface {
	GetBoards() (Boards, error)
}
