package ui

import (
	"bytes"
	"testing"

	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"
)

func TestConsoleDiffPrinter_PrintDiff(t *testing.T) {
	tests := map[string]struct {
		givenPrefix string
		givenDiff   string
		expectedOut string
	}{
		"GivenDiffWithUnchangedAddedRemovedLines_WhenPrefixGiven_ThenExpectHeader": {
			givenPrefix: "Diff: github.com/ccremer/greposync",
			givenDiff:   "normal line\n+added line\n-removed line",
			expectedOut: "\x1b[0m\x1b[0m\nnormal line\n\x1b[32m+added line\x1b[0m\n\x1b[32m\x1b[0m\x1b[31m-removed line\x1b[0m\n\x1b[31m\x1b[0m",
		},
		"GivenNoDiff_WhenDiffEmpty_ThenExpectNoOutput": {
			givenPrefix: "Prefix",
			givenDiff:   "",
			expectedOut: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := NewConsoleDiffPrinter()
			buf := &bytes.Buffer{}
			color.SetOutput(buf)
			c.PrintDiff(tt.givenPrefix, tt.givenDiff)

			assert.Contains(t, buf.String(), tt.expectedOut)
		})
	}
}
