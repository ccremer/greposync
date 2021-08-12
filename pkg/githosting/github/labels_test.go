package github

import (
	"testing"

	"github.com/ccremer/greposync/core"
	"github.com/google/go-github/v38/github"
	"github.com/stretchr/testify/assert"
)

func TestLabelImpl_ImplementsInterface(t *testing.T) {
	assert.Implements(t, (*core.Label)(nil), new(LabelImpl))
}

func TestProvider_hasLabelChanged(t *testing.T) {
	label := "label"
	description := "description"
	color := "ABABAB"

	tests := map[string]struct {
		givenGhLabel   github.Label
		givenRepoLabel LabelImpl
		expectedResult bool
	}{
		"GivenSameLabel_ThenExpectFalse": {
			givenGhLabel: *newLabel(label, color, description),
			givenRepoLabel: LabelImpl{
				Name:        label,
				Description: description,
				Color:       color,
			},
			expectedResult: false,
		},
		"GivenDifferentDescription_ThenExpectTrue": {
			givenGhLabel: *newLabel(label, color, description),
			givenRepoLabel: LabelImpl{
				Name:        label,
				Description: "different",
				Color:       color,
			},
			expectedResult: true,
		},
		"GivenDifferentColor_ThenExpectTrue": {
			givenGhLabel: *newLabel(label, color, description),
			givenRepoLabel: LabelImpl{
				Name:        label,
				Description: description,
				Color:       "FFFFFF",
			},
			expectedResult: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &Remote{}
			result := p.hasLabelChanged(&tt.givenGhLabel, &tt.givenRepoLabel)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestProvider_findMatchingGhLabel(t *testing.T) {
	tests := map[string]struct {
		givenGhLabels               []*github.Label
		givenRepoLabelForComparison *LabelImpl
		expectedLabelIndex          int
	}{
		"GivenNilList_ThenExpectNil": {
			givenGhLabels:      nil,
			expectedLabelIndex: -1,
		},
		"GivenEmptyList_ThenExpectNil": {
			givenGhLabels:      []*github.Label{},
			expectedLabelIndex: -1,
		},
		"GivenListWithMatchingLabel_ThenExpectLabel": {
			givenGhLabels: []*github.Label{
				newLabel("label", "ABABAB", "desc"),
			},
			givenRepoLabelForComparison: &LabelImpl{Name: "label"},
			expectedLabelIndex:          0,
		},
		"GivenListWithMatchingLabels_ThenExpectSecond": {
			givenGhLabels: []*github.Label{
				newLabel("label1", "ABABAB", "desc"),
				newLabel("label2", "ABABAB", "desc"),
			},
			givenRepoLabelForComparison: &LabelImpl{Name: "label2"},
			expectedLabelIndex:          1,
		},
		"GivenListWithNonMatchingLabel_ThenExpectNil": {
			givenGhLabels: []*github.Label{
				newLabel("label", "ABABAB", "desc"),
			},
			givenRepoLabelForComparison: &LabelImpl{Name: "different"},
			expectedLabelIndex:          -1,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &Remote{}
			result := p.findMatchingGhLabel(tt.givenGhLabels, tt.givenRepoLabelForComparison)
			if tt.expectedLabelIndex >= 0 {
				assert.Equal(t, tt.givenGhLabels[tt.expectedLabelIndex], result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func newLabel(label string, color string, description string) *github.Label {
	return &github.Label{
		Name:        &label,
		Color:       &color,
		Description: &description,
	}
}
