package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColor_CheckValue(t *testing.T) {
	tests := map[string]struct {
		givenColor    Color
		expectedError bool
	}{
		"GivenValidColor_ThenAcceptValue": {
			givenColor:    "#FFFFFF",
			expectedError: false,
		},
		"GivenLowercaseColor_ThenExpectError": {
			givenColor:    "#ffffff",
			expectedError: true,
		},
		"GivenColorWithPrefix_ThenExpectError": {
			givenColor:    "FFFFFF",
			expectedError: true,
		},
		"GivenEmptyColor_ThenExpectError": {
			givenColor:    "",
			expectedError: true,
		},
		"GivenArbitraryString_ThenExpectError": {
			givenColor:    "asdf",
			expectedError: true,
		},
		"Given3DigitColor_ThenExpectError": {
			givenColor:    "#FFF",
			expectedError: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.givenColor.CheckValue()
			if tt.expectedError {
				require.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestLabel_IsSameAs(t *testing.T) {
	labelName := "label"
	tests := map[string]struct {
		left           Label
		right          Label
		expectedResult bool
	}{
		"GivenSameNames_ThenExpectTrue": {
			left:           Label{Name: labelName},
			right:          Label{Name: labelName},
			expectedResult: true,
		},
		"GivenDifferentNames_ThenExpectFalse": {
			left:           Label{Name: labelName},
			right:          Label{Name: "different"},
			expectedResult: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.left.IsSameAs(tt.right)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestLabel_IsEqualTo(t *testing.T) {
	labelName := "label"
	description := "desc"
	var color Color = "#FFFFFF"
	tests := map[string]struct {
		left           Label
		right          Label
		expectedResult bool
	}{
		"GivenSameProperties_ThenExpectTrue": {
			left: Label{
				Name:        labelName,
				color:       color,
				Description: description,
			},
			right: Label{
				Name:        labelName,
				color:       color,
				Description: description,
			},
			expectedResult: true,
		},
		"GivenDifferentNames_ThenExpectFalse": {
			left:  Label{Name: labelName},
			right: Label{Name: "different"},
		},
		"GivenDifferentColors_ThenExpectFalse": {
			left:  Label{color: color},
			right: Label{color: ""},
		},
		"GivenDifferentDescriptions_ThenExpectFalse": {
			left:           Label{Description: description},
			right:          Label{Description: ""},
			expectedResult: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.left.IsEqualTo(tt.right)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

var labelSetCases = map[string]struct {
	givenSet      LabelSet
	hasDuplicates bool
	hasEmptyNames bool
}{
	"NilSet": {
		givenSet: nil,
	},
	"EmptySet": {
		givenSet: LabelSet{},
	},
	"SetWithDistinctLabels": {
		givenSet: LabelSet{
			Label{Name: "foo"},
			Label{Name: "bar"},
		},
	},
	"SetWithDuplicates": {
		givenSet: LabelSet{
			Label{Name: "foo"},
			Label{Name: "foo"},
		},
		hasDuplicates: true,
	},
	"SetWithEmptyNames": {
		givenSet: LabelSet{
			Label{Name: "foo"},
			Label{Name: ""},
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
