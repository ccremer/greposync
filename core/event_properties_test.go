package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventSource_ImplementsEventProperty(t *testing.T) {
	assert.Implements(t, (*EventProperty)(nil), new(EventSource))
}

func TestEventResult_ImplementsEventProperty(t *testing.T) {
	assert.Implements(t, (*EventProperty)(nil), new(EventResult))
}

func TestBoolProperty(t *testing.T) {
	key := "key"
	tests := map[string]struct {
		givenProperties EventProperties
		givenDefault    bool
		expectedResult  bool
	}{
		"GivenEmptyProperties_WhenDefaultIsFalse_ThenExpectFalse": {
			givenProperties: nil,
			givenDefault:    false,
			expectedResult:  false,
		},
		"GivenEmptyProperties_WhenDefaultIsTrue_ThenExpectTrue": {
			givenProperties: nil,
			givenDefault:    true,
			expectedResult:  true,
		},
		"GivenNonExistingProperty_WhenDefaultIsTrue_ThenExpectTrue": {
			givenProperties: map[string]interface{}{},
			givenDefault:    true,
			expectedResult:  true,
		},
		"GivenExistingProperty_WhenValueIsBoolean_ThenExpectValue": {
			givenProperties: map[string]interface{}{key: true},
			givenDefault:    false,
			expectedResult:  true,
		},
		"GivenExistingProperty_WhenValueIsConvertibleToBool_ThenExpectValue": {
			givenProperties: map[string]interface{}{key: "true"},
			givenDefault:    false,
			expectedResult:  true,
		},
		"GivenExistingProperty_WhenValueIsZero_ThenExpectFalse": {
			givenProperties: map[string]interface{}{key: 0},
			givenDefault:    true,
			expectedResult:  false,
		},
		"GivenExistingProperty_WhenValueIsOne_ThenExpectTrue": {
			givenProperties: map[string]interface{}{key: 1},
			givenDefault:    false,
			expectedResult:  true,
		},
		"GivenExistingProperty_WhenValueIsAnyOtherInt_ThenExpectDefault": {
			givenProperties: map[string]interface{}{key: 34},
			givenDefault:    true,
			expectedResult:  true,
		},
		"GivenExistingProperty_WhenValueIsAnotherInt_ThenExpectDefault": {
			// we don't need to support all different kind of integers, plain int is enough.
			// Users can still access the raw properties if they need to.
			givenProperties: map[string]interface{}{key: int64(1)},
			givenDefault:    false,
			expectedResult:  false,
		},
		"GivenExistingProperty_WhenValueIsNotConvertibleToBool_ThenExpectValue": {
			givenProperties: map[string]interface{}{key: "invalid"},
			givenDefault:    true,
			expectedResult:  true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := extractBoolProperty(tt.givenProperties, key, tt.givenDefault)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestEventSource_StringProperty(t *testing.T) {
	key := "key"
	defaultValue := "default"
	value := "value"
	tests := map[string]struct {
		givenProperties EventProperties
		givenDefault    string
		expectedResult  string
	}{
		"GivenEmptyProperties_WhenDefaultIsGiven_ThenExpectDefault": {
			givenProperties: nil,
			givenDefault:    defaultValue,
			expectedResult:  defaultValue,
		},
		"GivenEmptyProperties_WhenDefaultIsEmpty_ThenExpectEmptyString": {
			givenProperties: nil,
			givenDefault:    "",
			expectedResult:  "",
		},
		"GivenNonExistingProperties_WhenDefaultIsGiven_ThenExpectDefault": {
			givenProperties: map[string]interface{}{},
			givenDefault:    defaultValue,
			expectedResult:  defaultValue,
		},
		"GivenExistingProperties_WhenValueIsString_ThenExpectValue": {
			givenProperties: map[string]interface{}{key: value},
			givenDefault:    defaultValue,
			expectedResult:  value,
		},
		"GivenExistingProperties_WhenValueNotString_ThenExpectDefault": {
			givenProperties: map[string]interface{}{key: 0},
			givenDefault:    defaultValue,
			expectedResult:  defaultValue,
		},
		// Technically int and bool and other types can also be converted to string, but that defeats the purpose.
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := extractStringProperty(tt.givenProperties, key, tt.givenDefault)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
