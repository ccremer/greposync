package domain

import "fmt"

// LabelSet is a set of Label.
type LabelSet []Label

// CheckForEmptyLabelNames returns an error if there's a Label in the set that is an empty string.
func (s LabelSet) CheckForEmptyLabelNames() error {
	for _, label := range s {
		if label.Name == "" {
			return fmt.Errorf("%w: label name cannot be empty", ErrInvalidArgument)
		}
	}
	return nil
}

// CheckForDuplicates returns an error if two or more Label have the same Label.Name.
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

// FindLabelByName returns the Label by given Name, if there is one matching.
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

// Without returns a new LabelSet that contain only the labels that do not exist in other set.
// A label is not included in the result if the name matches.
//
// No validation checks are performed.
// The original order is preserved.
func (s LabelSet) Without(other LabelSet) LabelSet {
	if other == nil {
		return nil
	}

	newSet := make(LabelSet, 0)
	for i := range s {
		label := s[i]
		_, found := other.FindLabelByName(label.Name)
		if !found {
			newSet = append(newSet, label)
		}
	}
	return newSet
}
