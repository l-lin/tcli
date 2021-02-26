package trello

type ReactionSummaries []ReactionSummary
type ReactionSummary struct {
	ID         string        `json:"id"`
	IDReaction string        `json:"idReaction"`
	IDEmoji    string        `json:"idEmoji"`
	Count      int           `json:"count"`
	Emoji      ReactionEmoji `json:"emoji"`
}

type ReactionEmoji struct {
	Name      string `json:"name"`
	Native    string `json:"native"`
	ShortName string `json:"shortName"`
	Unified   string `json:"unified"`
}
