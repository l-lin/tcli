package renderer

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/l-lin/tcli/trello"
	"html/template"
)

type EditInPrettyToml struct{}

func NewEditInPrettyToml() Edit {
	return EditInPrettyToml{}
}

func (e EditInPrettyToml) MarshalCardToCreate(ctc trello.CardToCreate, lists trello.Lists, labels trello.Labels) ([]byte, error) {
	t := `
{{- /* ---------------- NAME ---------------- */ -}}
name = "{{.Card.Name}}"
{{/* ---------------- LISTS ---------------- */ -}}
# available lists:
{{- if .Lists -}}
{{range $list := .Lists}}
# {{$list.ID}}: {{$list.Name}}
{{- end -}}
{{end}}
idList = "{{.Card.IDList}}"
{{/* ---------------- POSITION ---------------- */ -}}
# the position of the card in its list: "top", "bottom" or a positive float
pos = "bottom"
{{/* ---------------- LABELS ---------------- */ -}}
# available labels (use color or ID):
{{- if .Labels -}}
{{range $label := .Labels}}
# {{$label.ID}}: {{$label.Color}}{{if $label.Name}} [{{$label.Name}}]{{end }}
{{- end -}}
{{end}}
labels = []
{{/* ---------------- DESCRIPTION ---------------- */ -}}
desc = '''
'''
`
	tpl := template.Must(template.New("create-card").Parse(t))
	tplParams := struct {
		Card   trello.CardToCreate
		Lists  trello.Lists
		Labels trello.Labels
	}{
		Card:   ctc,
		Lists:  lists,
		Labels: labels,
	}
	w := bytes.NewBufferString("")
	if err := tpl.Execute(w, tplParams); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (e EditInPrettyToml) MarshalCardToEdit(cte trello.CardToEdit, lists trello.Lists, labels trello.Labels) ([]byte, error) {
	t := `
{{- /* ---------------- NAME ---------------- */ -}}
name = "{{.Card.Name}}"
{{/* ---------------- CLOSED ---------------- */ -}}
# whether the card should be archived (closed: true)
closed = {{.Card.Closed}}
{{/* ---------------- LISTS ---------------- */ -}}
# available lists:
{{- if .Lists -}}
{{range $list := .Lists}}
# {{$list.ID}}: {{$list.Name}}
{{- end -}}
{{end}}
idList = "{{.Card.IDList}}"
{{/* ---------------- POSITION ---------------- */ -}}
# the position of the card in its list: "top", "bottom" or a positive float
pos = "{{.Card.Pos}}"
{{/* ---------------- LABELS ---------------- */ -}}
# available labels (use color or ID):
{{- if .Labels -}}
{{range $label := .Labels}}
# {{$label.ID}}: {{$label.Color}}{{if $label.Name}} [{{$label.Name}}]{{end }}
{{- end -}}
{{end}}
labels = [
{{- if .Card.Labels -}}
  {{- range $i, $label := .Card.Labels -}}
    {{if $i}},{{end}}"{{$label}}"
  {{- end -}}
{{- end -}}
]
{{/* ---------------- DESCRIPTION ---------------- */ -}}
desc = '''
{{htmlSafe .CardDescription}}
'''
`
	tpl := template.Must(template.New("edit-card").Funcs(template.FuncMap{
		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},
	}).Parse(t))
	tplParams := struct {
		Card            trello.CardToEdit
		Lists           trello.Lists
		Labels          trello.Labels
		CardDescription string
	}{
		Card:            cte,
		Lists:           lists,
		CardDescription: cte.Desc,
		Labels:          labels,
	}
	w := bytes.NewBufferString("")
	if err := tpl.Execute(w, tplParams); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (e EditInPrettyToml) Unmarshal(bytes []byte, i interface{}) error {
	return toml.Unmarshal(bytes, i)
}

func (e EditInPrettyToml) GetFileType() string {
	return "toml"
}
