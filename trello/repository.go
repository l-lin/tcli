//go:generate mockgen -source repository.go -package trello -destination repository_mock.go
package trello

type Repository interface {
	GetBoards() (Boards, error)
	FindBoard(query string) (*Board, error)
	GetLists(idBoard string) (Lists, error)
	FindList(idBoard string, query string) (*List, error)
	GetCards(idList string) (Cards, error)
	FindCard(idList string, query string) (*Card, error)
	UpdateCard(updateCard UpdateCard) (*Card, error)
}
