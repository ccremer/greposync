package domain

type Label struct {
	Name        string
	Description string
	color       Color
}

func (l Label) GetColor() Color {
	return l.color
}

func (l *Label) SetColor(color Color) error {
	if err := color.CheckValue(); hasFailed(err) {
		return err
	}
	l.color = color
	return nil
}

// IsSameAs returns true if each Label.Name is equal.
func (l Label) IsSameAs(label Label) bool {
	return l.Name == label.Name
}

// IsEqualTo returns true if all properties of Label are equal.
func (l Label) IsEqualTo(label Label) bool {
	return l.Name == label.Name && l.Description == label.Description && l.color == label.color
}
