//go:generate mockgen -source repository.go -package trello -destination repository_mock.go
package trello

// Repository to call perform CRUD operation on Trello resources
// We may want to update this interface to accept channels to support async
type Repository interface {
	FindBoards() (Boards, error)
	FindBoard(query string) (*Board, error)
	FindLabels(idBoard string) (Labels, error)
	FindLists(idBoard string) (Lists, error)
	FindList(idBoard string, query string) (*List, error)
	FindCards(idList string) (Cards, error)
	FindCard(idList string, query string) (*Card, error)
	CreateCard(createCard CreateCard) (*Card, error)
	UpdateCard(updateCard UpdateCard) (*Card, error)
	FindComments(idCard string) (Comments, error)
	FindComment(idCard string, idComment string) (*Comment, error)
	CreateComment(createComment CreateComment) (*Comment, error)
	UpdateComment(updateComment UpdateComment) (*Comment, error)
	DeleteComment(idCard, idComment string) error
	FindReactionSummaries(idComment string) (ReactionSummaries, error)
}
