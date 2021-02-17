package trello

import (
	"encoding/json"
	"fmt"
	"github.com/l-lin/tcli/conf"
	wrappedhttp "github.com/l-lin/tcli/http"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NewHttpRepository(c conf.Conf, debug bool) Repository {
	return HttpRepository{Conf: c, client: wrappedhttp.NewClient(debug)}
}

type HttpRepository struct {
	conf.Conf
	client *wrappedhttp.Client
}

func (h HttpRepository) GetBoards() (Boards, error) {
	v := h.buildQueries("id,name,shortLink,shortUrl,dateLastActivity")
	u := fmt.Sprintf("%s/members/me/boards?%v", h.BaseURL, v.Encode())

	var boards Boards
	if err := h.get(u, &boards); err != nil {
		return nil, err
	}
	return boards, nil
}

func (h HttpRepository) FindBoard(name string) (*Board, error) {
	boards, err := h.GetBoards()
	if err != nil {
		return nil, err
	}
	for _, board := range boards {
		if board.Name == name {
			return &board, nil
		}
	}
	return nil, nil
}

func (h HttpRepository) GetLists(idBoard string) (Lists, error) {
	v := h.buildQueries("id,name,idBoard")
	u := fmt.Sprintf("%s/boards/%s/lists?%v", h.BaseURL, idBoard, v.Encode())

	var lists Lists
	if err := h.get(u, &lists); err != nil {
		return nil, err
	}
	return lists, nil
}

func (h HttpRepository) FindList(idBoard string, name string) (*List, error) {
	// maybe use adequate API instead of getting all lists?
	// https://developer.atlassian.com/cloud/trello/rest/api-group-boards/#api-boards-id-lists-filter-get
	lists, err := h.GetLists(idBoard)
	if err != nil {
		return nil, err
	}
	for _, list := range lists {
		if list.Name == name {
			return &list, nil
		}
	}
	return nil, nil
}

func (h HttpRepository) GetCards(idList string) (Cards, error) {
	v := h.buildQueries("id,name,desc,idBoard,idList,labels,closed,shortLink,shortUrl")
	u := fmt.Sprintf("%s/lists/%s/cards?%v", h.BaseURL, idList, v.Encode())

	var cards Cards
	if err := h.get(u, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}

func (h HttpRepository) FindCard(idList string, name string) (*Card, error) {
	// maybe use adequate API instead of getting all cards?
	// https://developer.atlassian.com/cloud/trello/rest/api-group-boards/#api-boards-id-cards-filter-get
	cards, err := h.GetCards(idList)
	if err != nil {
		return nil, err
	}
	for _, card := range cards {
		if card.Name == name {
			return &card, nil
		}
	}
	return nil, nil
}

func (h HttpRepository) UpdateCard(updateCard UpdateCard) (*Card, error) {
	v := h.buildQueries("")
	u := fmt.Sprintf("%s/cards/%s?%v", h.BaseURL, updateCard.ID, v.Encode())
	var card Card
	if err := h.put(u, updateCard, &card); err != nil {
		return nil, err
	}
	return &card, nil
}

func (h HttpRepository) get(url string, ret interface{}) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	response, err := h.client.DoOnlyOk(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, ret)
}

func (h HttpRepository) put(url string, reqBody interface{}, ret interface{}) error {
	b, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("PUT", url, strings.NewReader(string(b)))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := h.client.DoOnlyOk(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	var respBody []byte
	respBody, err = io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, ret)
}

func (h HttpRepository) buildQueries(fields string) url.Values {
	v := url.Values{}
	v.Set("key", h.ApiKey)
	v.Set("token", h.AccessToken)
	if fields != "" {
		v.Set("fields", fields)
	}
	return v
}
