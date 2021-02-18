package conf

import (
	"os"
)

const (
	defaultTrelloApiBaseURL = "https://trello.com/1"
)

// Conf of the application
type Conf struct {
	Trello `yaml:"trello"`
}

type Trello struct {
	AppName     string `yaml:"-"`
	ApiKey      string `yaml:"api_key"`
	AccessToken string `yaml:"access_token"`
	BaseURL     string `yaml:"base_url"`
}

func NewConf() *Conf {
	return &Conf{
		Trello{
			BaseURL: defaultTrelloApiBaseURL,
		},
	}
}

func (c *Conf) areMandatoryFieldsFilled() bool {
	return c.Trello.ApiKey != "" && c.Trello.AccessToken != "" && c.Trello.BaseURL != ""
}

type Provider interface {
	Init() error
	Get() *Conf
}

func NewProvider(r Repository) Provider {
	return &provider{r: r}
}

type provider struct {
	r Repository
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
		askTrelloAppName().
		askTrelloApiKey().
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
