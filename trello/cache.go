package trello

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

// CacheInMemory is a decorator that caches the results of the proxified Repository
// in memory
type CacheInMemory struct {
	r Repository
	*Boards
	mapLists map[string]Lists
	mapCards map[string]Cards
}

func NewCacheInMemory(r Repository) Repository {
	return &CacheInMemory{
		r:        r,
		mapLists: map[string]Lists{},
		mapCards: map[string]Cards{},
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

func (c *CacheInMemory) FindBoard(query string) (*Board, error) {
	boards, err := c.GetBoards()
	if err != nil {
		return nil, err
	}
	if board := FindBoard(boards, query); board != nil {
		return board, nil
	}
	return nil, fmt.Errorf("no board found with query %s", query)
}

func (c *CacheInMemory) GetLists(idBoard string) (Lists, error) {
	if c.mapLists[idBoard] != nil {
		log.Debug().Str("idBoard", idBoard).Msg("fetching lists from cache")
		return c.mapLists[idBoard], nil
	}
	log.Debug().Str("idBoard", idBoard).Msg("fetching lists from remote")
	lists, err := c.r.GetLists(idBoard)
	c.mapLists[idBoard] = lists
	return lists, err
}

func (c *CacheInMemory) FindList(idBoard string, query string) (*List, error) {
	lists, err := c.GetLists(idBoard)
	if err != nil {
		return nil, err
	}
	if list := FindList(lists, query); list != nil {
		return list, nil
	}
	return nil, fmt.Errorf("no list found with query %s", query)
}

func (c *CacheInMemory) GetCards(idList string) (Cards, error) {
	if c.mapCards[idList] != nil {
		log.Debug().Str("idList", idList).Msg("fetching cards from cache")
		return c.mapCards[idList], nil
	}
	log.Debug().Str("idList", idList).Msg("fetching cards from remote")
	cards, err := c.r.GetCards(idList)
	c.mapCards[idList] = cards
	return cards, err
}

func (c *CacheInMemory) FindCard(idList string, query string) (*Card, error) {
	cards, err := c.GetCards(idList)
	if err != nil {
		return nil, err
	}
	if card := FindCard(cards, query); card != nil {
		return card, nil
	}
	return nil, fmt.Errorf("no card found with query %s", query)
}

func (c *CacheInMemory) UpdateCard(updateCard UpdateCard) (*Card, error) {
	card, err := c.r.UpdateCard(updateCard)
	if err != nil {
		return nil, err
	}

	// replace card
	for i, cachedCard := range c.mapCards[updateCard.IDList] {
		if cachedCard.ID == card.ID {
			c.mapCards[updateCard.IDList][i] = *card
			break
		}
	}
	return card, nil
}
