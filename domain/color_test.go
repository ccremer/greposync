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
