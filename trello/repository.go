//go:generate mockgen -source repository.go -package trello -destination repository_mock.go
package trello

type Repository interface {
	GetBoards() (Boards, error)
	FindBoard(name string) (*Board, error)
	GetLists(idBoard string) (Lists, error)
	FindList(idBoard string, name string) (*List, error)
	GetCards(idList string) (Cards, error)
	FindCard(idList string, name string) (*Card, error)
}
