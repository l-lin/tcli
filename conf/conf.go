package conf

import (
	"os"
)

const (
	defaultTrelloApiBaseURL = "https://trello.com/1/"
)

// Conf of the application
type Conf struct {
	Trello `json:"trello"`
}

type Trello struct {
	ApiKey      string `yaml:"api_key"`
	AccessToken string `yaml:"access_token"`
	AppName     string `yaml:"-"`
	BaseURL     string `yaml:"-"`
}

func NewConf() *Conf {
	return &Conf{
		Trello{
			BaseURL: defaultTrelloApiBaseURL,
		},
	}
}

func (c *Conf) areMandatoryFieldsFilled() bool {
	return c.Trello.ApiKey != "" && c.Trello.AccessToken != ""
}

type Provider interface {
	Init() error
	Get() *Conf
}

func NewProvider(r Repository, trelloDevKey, trelloAppName string) Provider {
	return &provider{
		r:             r,
		trelloDevKey:  trelloDevKey,
		trelloAppName: trelloAppName,
	}
}

type provider struct {
	r             Repository
	trelloDevKey  string
	trelloAppName string
}

func (p *provider) Init() error {
	if err := p.r.Init(); err != nil {
		return p.createIfNotExists()
	}
	if err := p.createIfNotExists(); err != nil {
		return err
	}
	return nil
}

func (p *provider) Get() *Conf {
	return p.r.Get()
}

func (p *provider) createIfNotExists() error {
	c := p.Get()
	if c.areMandatoryFieldsFilled() {
		return nil
	}
	confCreator := &creator{
		Conf:   NewConf(),
		stdin:  os.Stdin,
		stdout: os.Stdout,
	}
	var err error
	c, err = confCreator.
		setTrelloApiKey(p.trelloDevKey).
		setTrelloAppName(p.trelloAppName).
		askTrelloAccessToken().
		create()
	if err != nil {
		return err
	}
	if err = p.r.Save(c); err != nil {
		return err
	}
	return nil
}
