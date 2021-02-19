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
	mapLabels map[string]Labels // <idBoard, Labels>
	mapLists  map[string]Lists  // <idBoard, Lists>
	mapCards  map[string]Cards  // <idList, Cards>
}

func NewCacheInMemory(r Repository) Repository {
	return &CacheInMemory{
		r:         r,
		mapLabels: map[string]Labels{},
		mapLists:  map[string]Lists{},
		mapCards:  map[string]Cards{},
	}
}

func (c *CacheInMemory) FindBoards() (Boards, error) {
	if c.Boards != nil {
		log.Debug().Msg("fetching boards from cache")
		return *c.Boards, nil
	}
	log.Debug().Msg("fetching boards from remote")
	boards, err := c.r.FindBoards()
	c.Boards = &boards
	return boards, err
}

func (c *CacheInMemory) FindBoard(query string) (*Board, error) {
	boards, err := c.FindBoards()
	if err != nil {
		return nil, err
	}
	if board := FindBoard(boards, query); board != nil {
		return board, nil
	}
	return nil, fmt.Errorf("no board found with query %s", query)
}

func (c *CacheInMemory) FindLabels(idBoard string) (Labels, error) {
	if c.mapLabels[idBoard] != nil {
		log.Debug().Msg("fetching labels from cache")
		return c.mapLabels[idBoard], nil
	}
	log.Debug().Str("idBoard", idBoard).Msg("fetching labels from remote")
	labels, err := c.r.FindLabels(idBoard)
	c.mapLabels[idBoard] = labels
	return labels, err
}

func (c *CacheInMemory) FindLists(idBoard string) (Lists, error) {
	if c.mapLists[idBoard] != nil {
		log.Debug().Str("idBoard", idBoard).Msg("fetching lists from cache")
		return c.mapLists[idBoard], nil
	}
	log.Debug().Str("idBoard", idBoard).Msg("fetching lists from remote")
	lists, err := c.r.FindLists(idBoard)
	c.mapLists[idBoard] = lists
	return lists, err
}

func (c *CacheInMemory) FindList(idBoard string, query string) (*List, error) {
	lists, err := c.FindLists(idBoard)
	if err != nil {
		return nil, err
	}
	if list := FindList(lists, query); list != nil {
		return list, nil
	}
	return nil, fmt.Errorf("no list found with query %s", query)
}

func (c *CacheInMemory) FindCards(idList string) (Cards, error) {
	if c.mapCards[idList] != nil {
		log.Debug().Str("idList", idList).Msg("fetching cards from cache")
		return c.mapCards[idList], nil
	}
	log.Debug().Str("idList", idList).Msg("fetching cards from remote")
	cards, err := c.r.FindCards(idList)
	c.mapCards[idList] = cards
	return cards, err
}

func (c *CacheInMemory) FindCard(idList string, query string) (*Card, error) {
	cards, err := c.FindCards(idList)
	if err != nil {
		return nil, err
	}
	if card := FindCard(cards, query); card != nil {
		return card, nil
	}
	return nil, fmt.Errorf("no card found with query %s", query)
}

func (c *CacheInMemory) CreateCard(createCard CreateCard) (*Card, error) {
	card, err := c.r.CreateCard(createCard)
	if err != nil {
		return nil, err
	}

	// add card to cache
	c.mapCards[createCard.IDList] = append(c.mapCards[createCard.IDList], *card)
	return card, nil
}

func (c *CacheInMemory) UpdateCard(updateCard UpdateCard) (*Card, error) {
	card, err := c.r.UpdateCard(updateCard)
	if err != nil {
		return nil, err
	}

	cardIndex := c.findCardIndex(updateCard.IDList, updateCard.ID)
	if card.Closed {
		c.removeCard(updateCard.IDList, cardIndex)
	} else {
		c.mapCards[updateCard.IDList][cardIndex] = *card
	}
	return card, nil
}

func (c *CacheInMemory) findCardIndex(idList, cardID string) int {
	for i, cachedCard := range c.mapCards[idList] {
		if cachedCard.ID == cardID {
			return i
		}
	}
	return -1
}

func (c *CacheInMemory) removeCard(idList string, cardIndex int) {
	if cardIndex == -1 {
		return
	}
	cards, found := c.mapCards[idList]
	if !found {
		return
	}
	if cardIndex >= len(cards) {
		return
	}
	// swap
	cards[len(cards)-1], cards[cardIndex] = cards[cardIndex], cards[len(cards)-1]
	// then remove last element of slice
	c.mapCards[idList] = cards[:len(cards)-1]
}
