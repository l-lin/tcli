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

func New(tr trello.Repository, session *trello.Session) Completer {
	return Completer{
		tr:      tr,
		session: session,
	}
}

// Completer used to provide the content of the auto-completion for go-prompt
type Completer struct {
	tr          trello.Repository
	session     *trello.Session
	suggestions []prompt.Suggest
}

func (c Completer) Complete(cmd string, args []string) []prompt.Suggest {
	if !isKnownCmd(cmd) {
		return c.suggestCommands(cmd)
	}

	switch cmd {
	case "cd":
		return c.suggestForCD(args)
	case "cp":
		return c.suggestForCP(args)
	case "mv":
		return c.suggestForMV(args)
	}
	return c.suggestBoardsAndListsAndCardsAndCommentIDs(args)
}

func (c Completer) suggestForCD(args []string) []prompt.Suggest {
	if len(args) > 1 {
		return []prompt.Suggest{}
	}
	return c.suggestBoardsAndListsAndCards(args)
}

func (c Completer) suggestForCP(args []string) []prompt.Suggest {
	if len(args) < 2 {
		return c.suggestBoardsAndListsAndCardsAndCommentIDs(args)
	}
	if len(args) > 2 {
		return []prompt.Suggest{}
	}
	return c.suggestBoardsAndListsAndCards(args)
}

func (c Completer) suggestForMV(args []string) []prompt.Suggest {
	if len(args) < 2 {
		return c.suggestBoardsAndListsAndCards(args)
	}
	if len(args) > 2 {
		return []prompt.Suggest{}
	}
	return c.suggestBoardsAndLists(args)
}

func (c Completer) suggestBoardsAndLists(args []string) []prompt.Suggest {
	arg, p, err := c.resolvePath(args)
	if err != nil {
		log.Debug().
			Err(err).
			Str("arg", arg).
			Msg("could not resolve path")
		return []prompt.Suggest{}
	}
	if p.CardName != "" {
		return []prompt.Suggest{}
	}

	board, suggestions := c.suggestBoards(arg, p.BoardName)
	if suggestions != nil {
		return suggestions
	}

	_, suggestions = c.suggestLists(arg, board, p.ListName)
	return suggestions
}

func (c Completer) suggestBoardsAndListsAndCards(args []string) []prompt.Suggest {
	arg, p, err := c.resolvePath(args)
	if err != nil {
		log.Debug().
			Err(err).
			Str("arg", arg).
			Msg("could not resolve path")
		return []prompt.Suggest{}
	}

	if p.CommentID != "" {
		return []prompt.Suggest{}
	}

	board, suggestions := c.suggestBoards(arg, p.BoardName)
	if suggestions != nil {
		return suggestions
	}

	list, suggestions := c.suggestLists(arg, board, p.ListName)
	if suggestions != nil {
		return suggestions
	}

	_, suggestions = c.suggestCards(arg, list, p.CardName)
	return suggestions
}

func (c Completer) suggestBoardsAndListsAndCardsAndCommentIDs(args []string) []prompt.Suggest {
	arg, p, err := c.resolvePath(args)
	if err != nil {
		log.Debug().
			Err(err).
			Str("arg", arg).
			Msg("could not resolve path")
		return []prompt.Suggest{}
	}

	board, suggestions := c.suggestBoards(arg, p.BoardName)
	if suggestions != nil {
		return suggestions
	}

	list, suggestions := c.suggestLists(arg, board, p.ListName)
	if suggestions != nil {
		return suggestions
	}

	card, suggestions := c.suggestCards(arg, list, p.CardName)
	if suggestions != nil {
		return suggestions
	}
	return c.suggestComments(arg, card)
}

func (c Completer) resolvePath(args []string) (arg string, p trello.Path, err error) {
	arg = ""
	if len(args) > 0 {
		arg = args[len(args)-1]
	}
	pathResolver := trello.NewPathResolver(c.session)
	p, err = pathResolver.Resolve(arg)
	return
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
	var board *trello.Board
	if boardName != "" {
		board, _ = c.tr.FindBoard(boardName)
	}
	if board == nil {
		boards, err := c.tr.FindBoards()
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
	var list *trello.List
	if listName != "" {
		list, _ = c.tr.FindList(board.ID, listName)
	}
	if list == nil {
		lists, err := c.tr.FindLists(board.ID)
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

func (c Completer) suggestCards(arg string, list *trello.List, cardName string) (*trello.Card, []prompt.Suggest) {
	var card *trello.Card
	if cardName != "" {
		card, _ = c.tr.FindCard(list.ID, cardName)
	}
	if card == nil {
		cards, err := c.tr.FindCards(list.ID)
		if err != nil {
			log.Debug().
				Err(err).
				Str("idList", list.ID).
				Msg("could not find cards")
			return card, []prompt.Suggest{}
		}
		return card, prompt.FilterHasPrefix(cardsToSuggestions(cards), getBase(arg), true)
	}
	return card, nil
}

func (c Completer) suggestComments(arg string, card *trello.Card) []prompt.Suggest {
	var comments trello.Comments
	var err error
	if comments, err = c.tr.FindComments(card.ID); err != nil {
		log.Debug().
			Err(err).
			Str("idCard", card.ID).
			Msg("could not find comments")
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(commentsToSuggestions(comments), getBase(arg), true)
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
