package trello

func FindList(lists Lists, query string) *List {
	sanitizedQuery := sanitize(query)
	for _, list := range lists {
		if list.TCliID() == sanitizedQuery ||
			list.ID == query ||
			list.Name == query ||
			list.SanitizedName() == sanitizedQuery {
			return &list
		}
	}
	return nil
}

type Lists []List
type List struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IDBoard string `json:"idBoard"`
}

func (l List) TCliID() string {
	return toTCliID(l.SanitizedName(), l.ID)
}

func (l List) SanitizedName() string {
	return sanitize(l.Name)
}
