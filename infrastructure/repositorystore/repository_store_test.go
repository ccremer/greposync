package repositorystore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryStore_skipRepositories(t *testing.T) {
	tests := map[string]struct {
		givenIncludeRegex string
		givenExcludeRegex string
		givenTestString   string
		expectedSkip      bool
	}{
		"GivenEmptyFilters_WhenTestStringEmpty_ThenExpectTrue": {
			expectedSkip: true,
		},
		"GivenEmptyFilters_WhenTestStringHasContent_ThenExpectFalse": {
			givenTestString: "someString",
			expectedSkip:    false,
		},
		"GivenIncludeFilter_ThenExpectFalse": {
			givenTestString:   "someString",
			givenIncludeRegex: "some",
			expectedSkip:      false,
		},
		"GivenExcludeFilter_ThenExpectTrue": {
			givenTestString:   "someString",
			givenExcludeRegex: "some",
			expectedSkip:      true,
		},
		"GivenBothFilters_ThenExpectTrue": {
			givenTestString:   "someString",
			givenIncludeRegex: "String",
			givenExcludeRegex: "some",
			expectedSkip:      true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			includeRegex, excludeRegex, err := compileRegex(tt.givenIncludeRegex, tt.givenExcludeRegex)
			require.NoError(t, err)

			result := skipRepository(tt.givenTestString, includeRegex, excludeRegex)
			assert.Equal(t, tt.expectedSkip, result)
		})
	}
}
