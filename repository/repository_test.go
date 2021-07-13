package repository

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Regex(t *testing.T) {
	tests := map[string]struct {
		givenIncludeRegex string
		givenExcludeRegex string
		expectedErr       bool
		expectedSkip      bool
	}{
		"GivenEmptyRegex_WhenMatching_ThenReturnFalse": {},
		"GivenInvalidRegex_WhenCompiling_ExpectError": {
			givenIncludeRegex: "/\\",
			expectedErr:       true,
		},
		"GivenIncludeRegex_WhenMatching_ThenReturnFalse": {
			givenIncludeRegex: "repository",
		},
		"GivenExcludeRegex_WhenMatching_ThenReturnTrue": {
			givenExcludeRegex: "repository",
			expectedSkip:      true,
		},
		"GivenExcludeAndIncludeRegex_WhenMatching_ThenReturnTrue": {
			givenIncludeRegex: "repository",
			givenExcludeRegex: "repository",
			expectedSkip:      true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ir, er, err := compileRegex(tt.givenIncludeRegex, tt.givenExcludeRegex)
			if tt.expectedErr {
				require.Error(t, err)
				return
			}
			u, err := url.Parse("github.com/namespace/repository")
			require.NoError(t, err)
			result := skipRepository(u, ir, er)
			assert.Equal(t, tt.expectedSkip, result)
		})
	}
}
