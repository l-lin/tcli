package trello

func FindBoard(boards Boards, query string) *Board {
	sanitizedQuery := sanitize(query)
	for _, board := range boards {
		if board.TCliID() == sanitizedQuery ||
			board.ID == query ||
			board.ShortLink == query ||
			board.Name == query ||
			board.SanitizedName() == sanitizedQuery {
			return &board
		}
	}
	return nil
}

type Boards []Board
type Board struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShortLink        string `json:"shortLink"`
	ShortURL         string `json:"shortUrl"`
	DateLastActivity string `json:"dateLastActivity"`
}

func (b Board) TCliID() string {
	return toTCliID(b.SanitizedName(), b.ShortLink)
}

func (b Board) SanitizedName() string {
	return sanitize(b.Name)
}
