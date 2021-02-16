package completer

import (
	"github.com/c-bata/go-prompt"
	"github.com/l-lin/tcli/executor"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
	"path"
	"strings"
)

const maxCardDescriptionLength = 20

func New(tr trello.Repository, currentBoard *trello.Board, currentList *trello.List) Completer {
	return Completer{
		tr:           tr,
		currentBoard: currentBoard,
		currentList:  currentList,
	}
}

type Completer struct {
	tr           trello.Repository
	currentBoard *trello.Board
	currentList  *trello.List
	suggestions  []prompt.Suggest
}

func (c Completer) Complete(cmd, arg string) []prompt.Suggest {
	if !isKnownCmd(cmd) {
		return c.suggestCommands(cmd)
	}

	pathResolver := trello.NewPathResolver(c.currentBoard, c.currentList)
	boardName, listName, _, err := pathResolver.Resolve(arg)
	if err != nil {
		log.Debug().
			Err(err).
			Str("arg", arg).
			Msg("could not resolve path")
		return []prompt.Suggest{}
	}

	board, suggestions := c.suggestBoards(arg, boardName)
	if suggestions != nil {
		return suggestions
	}

	list, suggestions := c.suggestLists(arg, board, listName)
	if suggestions != nil {
		return suggestions
	}

	return c.suggestCards(arg, list)
}

func (c Completer) suggestCommands(arg string) []prompt.Suggest {
	suggestions := make([]prompt.Suggest, len(executor.Factories))
	for i, factory := range executor.Factories {
		suggestions[i] = prompt.Suggest{
			Text:        factory.Cmd,
			Description: factory.Description,
		}
	}
	return prompt.FilterHasPrefix(suggestions, arg, true)
}

func (c Completer) suggestBoards(arg string, boardName string) (*trello.Board, []prompt.Suggest) {
	board, _ := c.tr.FindBoard(boardName)
	if board == nil {
		boards, err := c.tr.GetBoards()
		if err != nil {
			log.Debug().
				Err(err).
				Msg("could not fetch boards")
			return board, []prompt.Suggest{}
		}
		return board, prompt.FilterHasPrefix(boardsToSuggestions(boards), getBase(arg), true)
	}
	return board, nil
}

func (c Completer) suggestLists(arg string, board *trello.Board, listName string) (*trello.List, []prompt.Suggest) {
	list, _ := c.tr.FindList(board.ID, listName)
	if list == nil {
		lists, err := c.tr.GetLists(board.ID)
		if err != nil {
			log.Debug().
				Err(err).
				Str("idBoard", board.ID).
				Msg("could not find lists")
			return list, []prompt.Suggest{}
		}
		return list, prompt.FilterHasPrefix(listsToSuggestions(lists), getBase(arg), true)
	}
	return list, nil
}

func (c Completer) suggestCards(arg string, list *trello.List) []prompt.Suggest {
	var cards trello.Cards
	var err error
	if cards, err = c.tr.GetCards(list.ID); err != nil {
		log.Debug().
			Err(err).
			Str("idList", list.ID).
			Msg("could not find cards")
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(cardsToSuggestions(cards), getBase(arg), true)
}

func isKnownCmd(cmd string) bool {
	for _, factory := range executor.Factories {
		if factory.Cmd == cmd {
			return true
		}
	}
	return false
}

func getBase(arg string) string {
	if arg == "" || strings.HasSuffix(arg, "/") {
		return ""
	}
	return path.Base(arg)
}