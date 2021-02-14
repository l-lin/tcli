package renderer

type Description interface {
	Render(description string) (string, error)
}

type PlainDescription struct{}

func (p PlainDescription) Render(description string) (string, error) {
	return description, nil
}
