package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
