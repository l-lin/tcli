package trello

import (
	"encoding/json"
	"fmt"
	"github.com/l-lin/tcli/conf"
	wrappedhttp "github.com/l-lin/tcli/http"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
)

func NewHttpRepository(c conf.Conf, debug bool) Repository {
	return HttpRepository{Conf: c, client: wrappedhttp.NewClient(debug)}
}

type HttpRepository struct {
	conf.Conf
	client *wrappedhttp.Client
}

func (h HttpRepository) GetBoards() (Boards, error) {
	v := url.Values{}
	v.Set("key", h.ApiKey)
	v.Set("token", h.AccessToken)
	u := fmt.Sprintf("%s/members/me/boards?%v", h.BaseURL, v.Encode())

	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("url", u).Msg("getting user")
	response, err := h.client.DoOnlyOk(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	var boards Boards
	if err = json.Unmarshal(body, &boards); err != nil {
		return nil, err
	}
	return boards, nil
}
