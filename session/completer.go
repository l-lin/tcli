package session

import (
	"github.com/c-bata/go-prompt"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
)

func completerAtBoardsLevel(cmd string, trelloRepository trello.Repository) func(prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		switch cmd {
		case "cd", "ls":
			suggestions, err := getBoardsSuggestions(trelloRepository)
			if err != nil {
				log.Debug().Err(err).Msg("could not fetch boards")
				return []prompt.Suggest{}
			}
			return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
		}
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "cd", Description: "move in the hierarchy"},
			{Text: "ls", Description: "show boards"},
		}, d.GetWordBeforeCursor(), true)
	}
}

func getBoardsSuggestions(trelloRepository trello.Repository) ([]prompt.Suggest, error) {
	boards, err := trelloRepository.GetBoards()
	if err != nil {
		return nil, err
	}
	suggestions := make([]prompt.Suggest, len(boards))
	for i, board := range boards {
		suggestions[i] = prompt.Suggest{
			Text: board.Name,
		}
	}
	return suggestions, nil
}

func completerAtListsLevel(cmd string, idBoard string, trelloRepository trello.Repository) func(prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		switch cmd {
		case "cd", "ls":
			suggestions, err := getListsSuggestions(trelloRepository, idBoard)
			if err != nil {
				log.Debug().Err(err).Msg("could not fetch lists")
				return []prompt.Suggest{}
			}
			suggestions = append(suggestions, prompt.Suggest{
				Text: "..",
			})
			return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
		}
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "cd", Description: "move in the hierarchy"},
			{Text: "ls", Description: "show lists"},
		}, d.GetWordBeforeCursor(), true)
	}
}

func getListsSuggestions(trelloRepository trello.Repository, idBoard string) ([]prompt.Suggest, error) {
	lists, err := trelloRepository.GetLists(idBoard)
	if err != nil {
		return nil, err
	}
	suggestions := make([]prompt.Suggest, len(lists))
	for i, list := range lists {
		suggestions[i] = prompt.Suggest{
			Text: list.Name,
		}
	}
	return suggestions, nil
}

func completerAtCardsLevel(cmd string, idList string, trelloRepository trello.Repository) func(prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		suggestions, err := getCardsSuggestions(trelloRepository, idList)
		if err != nil {
			log.Debug().Err(err).Msg("could not fetch cards")
			return []prompt.Suggest{}
		}
		switch cmd {
		case "cd":
			return prompt.FilterHasPrefix([]prompt.Suggest{{Text: ".."}}, d.GetWordBeforeCursor(), true)
		case "ls":
			suggestions = append(suggestions, prompt.Suggest{
				Text: "..",
			})
			return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
		case "edit":
			return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
		}
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "cd", Description: "move in the hierarchy"},
			{Text: "ls", Description: "show cards"},
			{Text: "edit", Description: "edit card"},
		}, d.GetWordBeforeCursor(), true)
	}
}

func getCardsSuggestions(trelloRepository trello.Repository, idList string) ([]prompt.Suggest, error) {
	cards, err := trelloRepository.GetCards(idList)
	if err != nil {
		return nil, err
	}
	suggestions := make([]prompt.Suggest, len(cards))
	for i, card := range cards {
		suggestions[i] = prompt.Suggest{
			Text:        card.Name,
			Description: truncateCardDescription(card.Description),
		}
	}
	return suggestions, nil
}

func truncateCardDescription(description string) string {
	if len(description) > maxCardDescriptionLength {
		return description[:maxCardDescriptionLength]
	}
	return description
}
