package conf

import (
	"os"
)

const (
	defaultTrelloApiBaseURL = "https://trello.com/1"
	defaultEditor           = "editor"
	defaultFormat           = "yaml"
	defaultPrompt           = false
)

var allFormats = []string{"yaml", "toml"}

// Conf of the application
type Conf struct {
	Trello      `yaml:"trello"`
	Editor      string `yaml:"editor"`
	Format      string `yaml:"format"`
	NeverPrompt bool   `yaml:"never_prompt"`
}

type Trello struct {
	AppName             string `yaml:"-"`
	ApiKey              string `yaml:"api_key"`
	AccessToken         string `yaml:"access_token"`
	BaseURL             string `yaml:"base_url"`
	TrelloDefaultConfig `yaml:"default_config"`
}

type TrelloDefaultConfig struct {
	TrelloDefaultBoard `yaml:"board"`
	TrelloDefaultList  `yaml:"list"`
	Labels             []string `yaml:"labels"`
}

type TrelloDefaultBoard struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type TrelloDefaultList struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

func NewConf() *Conf {
	return &Conf{
		Trello: Trello{
			BaseURL: defaultTrelloApiBaseURL,
		},
		Editor:      defaultEditor,
		Format:      defaultFormat,
		NeverPrompt: defaultPrompt,
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
		askFormat().
		create()
	if err != nil {
		return err
	}
	if err = p.r.Save(c); err != nil {
		return err
	}
	return nil
}
