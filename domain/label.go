package domain

import (
	"fmt"
	"regexp"
)

type Color string

type Label struct {
	Name        string
	Description string
	color       Color
}

type LabelSet []Label

var colorRegex *regexp.Regexp

func init() {
	colorRegex = regexp.MustCompile("^#[A-F0-9]{6}$")
}

// IsSameAs returns true if each Label.Name is equal.
func (l Label) IsSameAs(label Label) bool {
	return l.Name == label.Name
}

// IsEqualTo returns true if all properties of Label are equal.
func (l Label) IsEqualTo(label Label) bool {
	return l.Name == label.Name && l.Description == label.Description && l.color == label.color
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

func (c Color) String() string {
	return string(c)
}

func (c Color) CheckValue() error {
	if colorRegex.MatchString(c.String()) {
		return nil
	}
	return fmt.Errorf("%w: color value must be 6-digit uppercased hexadecimal with '#' prefix: %s", ErrInvalidArgument, c)
}

func (s LabelSet) CheckForEmptyLabelNames() error {
	for _, label := range s {
		if label.Name == "" {
			return fmt.Errorf("%w: label name cannot be empty", ErrInvalidArgument)
		}
	}
	return nil
}

func (s LabelSet) CheckForDuplicates() error {
	m := make(map[string]int, len(s))
	for _, label := range s {
		if _, exists := m[label.Name]; exists {
			// TODO: Another Error type maybe?
			return fmt.Errorf("%w: label is duplicated", ErrInvalidArgument)
		}
		m[label.Name] = 1
	}
	return nil
}

func (s LabelSet) FindLabelByName(label string) (Label, bool) {
	for _, l := range s {
		if l.Name == label {
			return l, true
		}
	}
	return Label{}, false
}
