package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type changeSet struct {
	other     LabelSet
	resultSet LabelSet
}

var labelSetCases = map[string]struct {
	givenSet      LabelSet
	mergeSet      changeSet
	diffSet       changeSet
	hasDuplicates bool
	hasEmptyNames bool
}{
	"NilSet": {
		givenSet: nil,
		mergeSet: changeSet{
			other:     nil,
			resultSet: nil,
		},
		diffSet: changeSet{
			other:     nil,
			resultSet: nil,
		},
	},
	"EmptySet": {
		givenSet: LabelSet{},
		mergeSet: changeSet{
			other:     LabelSet{},
			resultSet: LabelSet{},
		},
		diffSet: changeSet{
			other:     LabelSet{},
			resultSet: LabelSet{},
		},
	},
	"SetWithDistinctLabels": {
		givenSet: LabelSet{
			Label{Name: "foo"},
			Label{Name: "bar"},
		},
		mergeSet: changeSet{
			other: LabelSet{
				Label{Name: "new"},
			},
			resultSet: LabelSet{
				Label{Name: "new"},
				Label{Name: "foo"},
				Label{Name: "bar"},
			},
		},
		diffSet: changeSet{
			other: LabelSet{
				Label{Name: "bar"},
				Label{Name: "foo"},
			},
			resultSet: LabelSet{},
		},
	},
	"SetWithDuplicates": {
		givenSet: LabelSet{
			Label{Name: "foo"},
			Label{Name: "foo"},
		},
		mergeSet: changeSet{
			other: LabelSet{
				Label{Name: "new"},
			},
			resultSet: LabelSet{
				Label{Name: "new"},
				Label{Name: "foo"},
			},
		},
		diffSet: changeSet{
			other: LabelSet{
				Label{Name: "foo"},
			},
			resultSet: LabelSet{},
		},
		hasDuplicates: true,
	},
	"SetWithEmptyNames": {
		givenSet: LabelSet{
			Label{Name: "foo"},
			Label{Name: ""},
		},
		mergeSet: changeSet{
			other: LabelSet{
				Label{Name: ""},
			},
			resultSet: LabelSet{
				Label{Name: ""},
				Label{Name: "foo"},
			},
		},
		diffSet: changeSet{
			other: LabelSet{
				Label{Name: "bar"},
			},
			resultSet: LabelSet{
				Label{Name: "foo"},
				Label{Name: ""},
			},
		},
		hasEmptyNames: true,
	},
}

func TestLabelSet_CheckForDuplicates(t *testing.T) {
	for name, tt := range labelSetCases {
		t.Run(name, func(t *testing.T) {
			result := tt.givenSet.CheckForDuplicates()
			if tt.hasDuplicates {
				assert.Error(t, result)
				return
			}
			assert.NoError(t, result)
		})
	}
}

func TestLabelSet_CheckForEmptyLabelNames(t *testing.T) {
	for name, tt := range labelSetCases {
		t.Run(name, func(t *testing.T) {
			result := tt.givenSet.CheckForEmptyLabelNames()
			if tt.hasEmptyNames {
				assert.Error(t, result)
				return
			}
			assert.NoError(t, result)
		})
	}
}

func TestLabelSet_Merge(t *testing.T) {
	for name, tt := range labelSetCases {
		t.Run(name, func(t *testing.T) {
			result := tt.givenSet.Merge(tt.mergeSet.other)
			assert.Equal(t, tt.mergeSet.resultSet, result)
		})
	}
}

func TestLabelSet_DifferenceOf(t *testing.T) {
	for name, tt := range labelSetCases {
		t.Run(name, func(t *testing.T) {
			result := tt.givenSet.DifferenceOf(tt.diffSet.other)
			assert.Equal(t, tt.diffSet.resultSet, result)
		})
	}
}
