package domain

// Label is a Value object containing the properties of labels in a Git hosting service.
type Label struct {
	// Name is the label name
	Name        string
	// Description adds additional details to the label.
	Description string
	color       Color
}

// GetColor returns the color of the Label.
func (l Label) GetColor() Color {
	return l.color
}

// SetColor sets the color of the Label.
// If Color.CheckValue fails, then that error is returned.
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
