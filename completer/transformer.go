package completer

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/l-lin/tcli/trello"
)

func boardsToSuggestions(boards trello.Boards) []prompt.Suggest {
	suggestions := make([]prompt.Suggest, len(boards))
	for i, board := range boards {
		suggestions[i] = prompt.Suggest{
			Text: fmt.Sprintf("%s", board.TCliID()),
		}
	}
	return suggestions
}

func listsToSuggestions(lists trello.Lists) []prompt.Suggest {
	suggestions := make([]prompt.Suggest, len(lists))
	for i, list := range lists {
		suggestions[i] = prompt.Suggest{
			Text: fmt.Sprintf("%s", list.TCliID()),
		}
	}
	return suggestions
}

func cardsToSuggestions(cards trello.Cards) []prompt.Suggest {
	suggestions := make([]prompt.Suggest, len(cards))
	for i, card := range cards {
		suggestions[i] = prompt.Suggest{
			Text:        fmt.Sprintf("%s", card.TCliID()),
			Description: truncateCardDescription(card.Desc),
		}
	}
	return suggestions
}

func truncateCardDescription(description string) string {
	if len(description) > maxCardDescriptionLength {
		return description[:maxCardDescriptionLength]
	}
	return description
}
