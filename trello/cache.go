package trello

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

type CacheInMemory struct {
	r Repository
	*Boards
	MapLists map[string]Lists
	MapCards map[string]Cards
}

func NewCacheInMemory(r Repository) Repository {
	return &CacheInMemory{
		r:        r,
		MapLists: map[string]Lists{},
		MapCards: map[string]Cards{},
	}
}

func (c *CacheInMemory) GetBoards() (Boards, error) {
	if c.Boards != nil {
		log.Debug().Msg("fetching boards from cache")
		return *c.Boards, nil
	}
	log.Debug().Msg("fetching boards from remote")
	boards, err := c.r.GetBoards()
	c.Boards = &boards
	return boards, err
}

func (c *CacheInMemory) FindBoard(name string) (*Board, error) {
	boards, err := c.GetBoards()
	if err != nil {
		return nil, err
	}
	for _, board := range boards {
		if board.Name == name {
			return &board, nil
		}
	}
	return nil, fmt.Errorf("no board found with name %s", name)
}

func (c *CacheInMemory) GetLists(idBoard string) (Lists, error) {
	if c.MapLists[idBoard] != nil {
		log.Debug().Str("idBoard", idBoard).Msg("fetching lists from cache")
		return c.MapLists[idBoard], nil
	}
	log.Debug().Str("idBoard", idBoard).Msg("fetching lists from remote")
	lists, err := c.r.GetLists(idBoard)
	c.MapLists[idBoard] = lists
	return lists, err
}

func (c *CacheInMemory) FindList(idBoard string, name string) (*List, error) {
	lists, err := c.GetLists(idBoard)
	if err != nil {
		return nil, err
	}
	for _, list := range lists {
		if list.Name == name {
			return &list, nil
		}
	}
	return nil, fmt.Errorf("no list found with name %s", name)
}

func (c *CacheInMemory) GetCards(idList string) (Cards, error) {
	if c.MapCards[idList] != nil {
		log.Debug().Str("idList", idList).Msg("fetching cards from cache")
		return c.MapCards[idList], nil
	}
	log.Debug().Str("idList", idList).Msg("fetching cards from remote")
	cards, err := c.r.GetCards(idList)
	c.MapCards[idList] = cards
	return cards, err
}

func (c *CacheInMemory) FindCard(idList string, name string) (*Card, error) {
	cards, err := c.GetCards(idList)
	if err != nil {
		return nil, err
	}
	for _, card := range cards {
		if card.Name == name {
			return &card, nil
		}
	}
	return nil, fmt.Errorf("no card found with name %s", name)
}
