package renderer

import (
	"bytes"
	"github.com/l-lin/tcli/trello"
	"gopkg.in/yaml.v2"
	"html/template"
	"strings"
)

type EditInYaml struct{}

func NewEditInYaml() Edit {
	return EditInYaml{}
}

func (e EditInYaml) MarshalCardToEdit(cte trello.CardToEdit, _ trello.Lists) ([]byte, error) {
	return yaml.Marshal(cte)
}

func (e EditInYaml) MarshalCardToCreate(create trello.CardToCreate, _ trello.Lists) ([]byte, error) {
	return yaml.Marshal(create)
}

func (e EditInYaml) Unmarshal(in []byte, v interface{}) error {
	return yaml.Unmarshal(in, v)
}

func NewEditInPrettyYaml() Edit {
	return EditInPrettyYaml{}
}

type EditInPrettyYaml struct{}

func (e EditInPrettyYaml) MarshalCardToEdit(cte trello.CardToEdit, boardLists trello.Lists) ([]byte, error) {
	t := `name: "{{.Card.Name}}"
# whether the card should be archived (closed: true)
closed: {{.Card.Closed}}
# available board lists:{{if .Lists}}
{{range $list := .Lists}}# {{$list.ID}}: {{$list.Name}}
{{end}}{{end}}idList: "{{.Card.IDList}}"
# the position of the card in its list: "top", "bottom" or a positive float
pos: {{.Card.Pos}}
desc: |-
{{htmlSafe .CardDescription}}`
	tpl := template.Must(template.New("edit-card").Funcs(template.FuncMap{
		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},
	}).Parse(t))
	tplParams := struct {
		Card            trello.CardToEdit
		Lists           trello.Lists
		CardDescription string
	}{
		Card:            cte,
		Lists:           boardLists,
		CardDescription: e.transformDescription(cte.Desc),
	}
	w := bytes.NewBufferString("")
	if err := tpl.Execute(w, tplParams); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (e EditInPrettyYaml) MarshalCardToCreate(ctc trello.CardToCreate, boardLists trello.Lists) ([]byte, error) {
	t := `name: "{{.Card.Name}}"
# available board lists:{{if .Lists}}
{{range $list := .Lists}}# {{$list.ID}}: {{$list.Name}}
{{end}}{{end}}idList: "{{.Card.IDList}}"
# the position of the card in its list: "top", "bottom" or a positive float
pos: "bottom"
desc: |-
  `
	tpl := template.Must(template.New("create-card").Parse(t))
	tplParams := struct {
		Card  trello.CardToCreate
		Lists trello.Lists
	}{
		Card:  ctc,
		Lists: boardLists,
	}
	w := bytes.NewBufferString("")
	if err := tpl.Execute(w, tplParams); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (e EditInPrettyYaml) Unmarshal(in []byte, v interface{}) error {
	return yaml.Unmarshal(in, v)
}

func (e EditInPrettyYaml) transformDescription(desc string) string {
	sb := strings.Builder{}
	for _, s := range strings.Split(desc, "\n") {
		if s != "" {
			sb.WriteString("  ")
			sb.WriteString(s)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
