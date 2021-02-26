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
	mapLabelsByIDBoard  map[string]Labels   // <idBoard, Labels>
	mapListsByIDBoard   map[string]Lists    // <idBoard, Lists>
	mapCardsByIDList    map[string]Cards    // <idList, Cards>
	mapCommentsByIDCard map[string]Comments // <idCard, Comments>
}

func NewCacheInMemory(r Repository) Repository {
	return &CacheInMemory{
		r:                   r,
		mapLabelsByIDBoard:  map[string]Labels{},
		mapListsByIDBoard:   map[string]Lists{},
		mapCardsByIDList:    map[string]Cards{},
		mapCommentsByIDCard: map[string]Comments{},
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
	if c.mapLabelsByIDBoard[idBoard] != nil {
		log.Debug().Msg("fetching labels from cache")
		return c.mapLabelsByIDBoard[idBoard], nil
	}
	log.Debug().Str("idBoard", idBoard).Msg("fetching labels from remote")
	labels, err := c.r.FindLabels(idBoard)
	c.mapLabelsByIDBoard[idBoard] = labels
	return labels, err
}

func (c *CacheInMemory) FindLists(idBoard string) (Lists, error) {
	if c.mapListsByIDBoard[idBoard] != nil {
		log.Debug().Str("idBoard", idBoard).Msg("fetching lists from cache")
		return c.mapListsByIDBoard[idBoard], nil
	}
	log.Debug().Str("idBoard", idBoard).Msg("fetching lists from remote")
	lists, err := c.r.FindLists(idBoard)
	c.mapListsByIDBoard[idBoard] = lists
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
	if c.mapCardsByIDList[idList] != nil {
		log.Debug().Str("idList", idList).Msg("fetching cards from cache")
		return c.mapCardsByIDList[idList], nil
	}
	log.Debug().Str("idList", idList).Msg("fetching cards from remote")
	cards, err := c.r.FindCards(idList)
	c.mapCardsByIDList[idList] = cards
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
	c.mapCardsByIDList[createCard.IDList] = append(c.mapCardsByIDList[createCard.IDList], *card)
	return card, nil
}

func (c *CacheInMemory) UpdateCard(updateCard UpdateCard) (*Card, error) {
	card, err := c.r.UpdateCard(updateCard)
	if err != nil {
		return nil, err
	}

	cardIndex := c.findCardIndex(updateCard.IDList, updateCard.ID)
	if cardIndex == -1 {
		// card may have been moved, so clear cache completely like a brute
		// we can find a smarter way, but well, performance wise, it's still acceptable...
		c.mapCardsByIDList = map[string]Cards{}
	} else {
		if card.Closed {
			c.removeCard(updateCard.IDList, cardIndex)
		} else {
			c.mapCardsByIDList[updateCard.IDList][cardIndex] = *card
		}
	}
	return card, nil
}

func (c *CacheInMemory) FindComments(idCard string) (Comments, error) {
	if c.mapCommentsByIDCard[idCard] != nil {
		log.Debug().Str("idCard", idCard).Msg("fetching comments from cache")
		return c.mapCommentsByIDCard[idCard], nil
	}
	log.Debug().Str("idCard", idCard).Msg("fetching comments from remote")
	comments, err := c.r.FindComments(idCard)
	c.mapCommentsByIDCard[idCard] = comments
	return comments, err
}

func (c *CacheInMemory) FindComment(idCard string, idComment string) (*Comment, error) {
	comments, err := c.FindComments(idCard)
	if err != nil {
		return nil, err
	}
	if comment := FindComment(comments, idComment); comment != nil {
		return comment, nil
	}
	return nil, fmt.Errorf("no comment found with id %s", idComment)
}

func (c *CacheInMemory) CreateComment(createComment CreateComment) (*Comment, error) {
	comment, err := c.r.CreateComment(createComment)
	if err != nil {
		return nil, err
	}

	// add comment to cache
	c.mapCommentsByIDCard[createComment.IDCard] = append(c.mapCommentsByIDCard[createComment.IDCard], *comment)
	return comment, nil
}

func (c *CacheInMemory) UpdateComment(updateComment UpdateComment) (*Comment, error) {
	comment, err := c.r.UpdateComment(updateComment)
	if err != nil {
		return nil, err
	}

	commentIndex := c.findCommentIndex(updateComment.IDCard, updateComment.ID)
	if commentIndex != -1 {
		c.mapCommentsByIDCard[updateComment.IDCard][commentIndex] = *comment
	} else {
		c.mapCommentsByIDCard[updateComment.IDCard] = append(c.mapCommentsByIDCard[updateComment.IDCard], *comment)
	}
	return comment, nil
}

func (c *CacheInMemory) DeleteComment(idCard, idComment string) error {
	if err := c.r.DeleteComment(idCard, idComment); err != nil {
		return err
	}

	commentIndex := c.findCommentIndex(idCard, idComment)
	c.removeComment(idCard, commentIndex)
	return nil
}

func (c *CacheInMemory) findCardIndex(idList, idCard string) int {
	for i, cachedCard := range c.mapCardsByIDList[idList] {
		if cachedCard.ID == idCard {
			return i
		}
	}
	return -1
}

func (c *CacheInMemory) findCommentIndex(idCard, idComment string) int {
	for i, cachedComment := range c.mapCommentsByIDCard[idCard] {
		if cachedComment.ID == idComment {
			return i
		}
	}
	return -1
}

func (c *CacheInMemory) removeCard(idList string, cardIndex int) {
	if cardIndex == -1 {
		return
	}
	cards, found := c.mapCardsByIDList[idList]
	if !found {
		return
	}
	if cardIndex >= len(cards) {
		return
	}
	c.mapCardsByIDList[idList] = append(cards[:cardIndex], cards[cardIndex+1:]...)
}

func (c *CacheInMemory) removeComment(idCard string, commentIndex int) {
	if commentIndex == -1 {
		return
	}
	comments, found := c.mapCommentsByIDCard[idCard]
	if !found {
		return
	}
	if commentIndex >= len(comments) {
		return
	}
	c.mapCommentsByIDCard[idCard] = append(comments[:commentIndex], comments[commentIndex+1:]...)
}
