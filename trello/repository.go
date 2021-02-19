//go:generate mockgen -source repository.go -package trello -destination repository_mock.go
package trello

type Repository interface {
	FindBoards() (Boards, error)
	FindBoard(query string) (*Board, error)
	FindLists(idBoard string) (Lists, error)
	FindList(idBoard string, query string) (*List, error)
	FindCards(idList string) (Cards, error)
	FindCard(idList string, query string) (*Card, error)
	UpdateCard(updateCard UpdateCard) (*Card, error)
}
