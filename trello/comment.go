package trello

import (
	"sort"
	"time"
)

const dateTimeLayout = "2006-01-02T15:04:05.000Z"

func FindComment(comments Comments, idComment string) *Comment {
	for _, comment := range comments {
		if idComment == comment.ID {
			return &comment
		}
	}
	return nil
}

type Comments []Comment

func (c Comments) SortedByDateDesc() Comments {
	sort.Slice(c, func(i, j int) bool {
		ti, err := time.Parse(dateTimeLayout, c[i].Date)
		if err != nil {
			return false
		}
		tj, err := time.Parse(dateTimeLayout, c[j].Date)
		if err != nil {
			return true
		}
		return ti.After(tj)
	})

	return c
}

type Comment struct {
	ID            string               `json:"id"`
	Date          string               `json:"date"`
	Data          CommentData          `json:"data"`
	MemberCreator CommentMemberCreator `json:"memberCreator"`
}

type CommentData struct {
	Card CommentDataCard `json:"card"`
	Text string          `json:"text"`
}

type CommentDataCard struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortLink string `json:"shortLink"`
}

type CommentMemberCreator struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Initials string `json:"initials"`
	Username string `json:"username"`
}

type CreateComment struct {
	IDCard string `json:"idCard"`
	Text   string `json:"text"`
}

type UpdateComment struct {
	ID     string `json:"id"`
	IDCard string `json:"idCard"`
	Text   string `json:"text"`
}
