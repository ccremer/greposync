package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type changeSet struct {
	toMerge   LabelSet
	resultSet LabelSet
}

var labelSetCases = map[string]struct {
	givenSet      LabelSet
	changeSet     changeSet
	hasDuplicates bool
	hasEmptyNames bool
}{
	"NilSet": {
		givenSet: nil,
		changeSet: changeSet{
			toMerge:   nil,
			resultSet: nil,
		},
	},
	"EmptySet": {
		givenSet: LabelSet{},
		changeSet: changeSet{
			toMerge:   LabelSet{},
			resultSet: LabelSet{},
		},
	},
	"SetWithDistinctLabels": {
		givenSet: LabelSet{
			Label{Name: "foo"},
			Label{Name: "bar"},
		},
		changeSet: changeSet{
			toMerge: LabelSet{
				Label{Name: "new"},
			},
			resultSet: LabelSet{
				Label{Name: "new"},
				Label{Name: "foo"},
				Label{Name: "bar"},
			},
		},
	},
	"SetWithDuplicates": {
		givenSet: LabelSet{
			Label{Name: "foo"},
			Label{Name: "foo"},
		},
		changeSet: changeSet{
			toMerge: LabelSet{
				Label{Name: "new"},
			},
			resultSet: LabelSet{
				Label{Name: "new"},
				Label{Name: "foo"},
			},
		},
		hasDuplicates: true,
	},
	"SetWithEmptyNames": {
		givenSet: LabelSet{
			Label{Name: "foo"},
			Label{Name: ""},
		},
		changeSet: changeSet{
			toMerge: LabelSet{
				Label{Name: ""},
			},
			resultSet: LabelSet{
				Label{Name: ""},
				Label{Name: "foo"},
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
			result := tt.givenSet.Merge(tt.changeSet.toMerge)
			assert.Equal(t, tt.changeSet.resultSet, result)
		})
	}
}
