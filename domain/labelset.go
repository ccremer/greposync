package domain

import "fmt"

type LabelSet []Label

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

// Merge returns a new copy of LabelSet that contains the Label from other if they are missing in the original slice, and replaces existing ones.
// A label to replace is determined by equality of LabelSet.FindLabelByName.
//
// No validation checks are performed.
// The original order is not preserved.
// Duplicates are removed from the result.
func (s LabelSet) Merge(other LabelSet) LabelSet {
	if other == nil {
		return s
	}

	newSet := make(LabelSet, len(other))
	// Copy other set, this is the minimum
	for i := range other {
		newSet[i] = other[i]
	}

	// add the remaining from s, provided they aren't already in the list.
	for i := range s {
		label := s[i]
		_, found := newSet.FindLabelByName(label.Name)
		if !found {
			newSet = append(newSet, label)
		}
	}
	return newSet
}
