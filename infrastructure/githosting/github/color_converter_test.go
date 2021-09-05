package github

import (
	"testing"

	"github.com/ccremer/greposync/domain"
	"github.com/stretchr/testify/assert"
)

func TestColorConverter_ConvertToEntity(t *testing.T) {
	tests := map[string]struct {
		givenColor     string
		expectedResult domain.Color
	}{
		"GivenEmptyString_ThenReturnEmpty": {
			givenColor:     "",
			expectedResult: "",
		},
		"GivenInvalidColor_ThenReturnEmpty": {
			givenColor:     "invalid",
			expectedResult: "",
		},
		"GivenValidColor_WhenLowerCaseWithoutPrefix_ThenReturnUppercaseWithPrefix": {
			// This is actually how GitHub handles colors
			givenColor:     "ababab",
			expectedResult: "#ABABAB",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			converter := ColorConverter{}
			result := converter.ConvertToEntity(tt.givenColor)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestColorConverter_ConvertFromEntity(t *testing.T) {
	tests := map[string]struct {
		givenColor     domain.Color
		expectedResult string
	}{
		"GivenEmptyString_ThenReturnEmpty": {
			givenColor:     "",
			expectedResult: "",
		},
		"GivenInvalidColor_ThenReturnEmpty": {
			givenColor:     "invalid",
			expectedResult: "",
		},
		"GivenValidColor_WhenLowerCaseWithoutPrefix_ThenReturnUppercaseWithPrefix": {
			// This is actually how GitHub handles colors
			givenColor:     "#ABABAB",
			expectedResult: "ababab",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			converter := ColorConverter{}
			result := converter.ConvertFromEntity(tt.givenColor)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
