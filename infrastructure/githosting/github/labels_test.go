package github

import (
	"net/url"
	"testing"

	"github.com/ccremer/greposync/domain"
	"github.com/google/go-github/v39/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_hasLabelChanged(t *testing.T) {
	label := "label"
	description := "description"
	color := "ababab"

	tests := map[string]struct {
		givenGhLabel   *github.Label
		givenRepoLabel domain.Label
		expectedResult bool
	}{
		"GivenSameLabel_ThenExpectFalse": {
			givenGhLabel:   newGitHubLabel(label, description, color),
			givenRepoLabel: newDomainLabel(t, label, description, color),
			expectedResult: false,
		},
		"GivenDifferentDescription_ThenExpectTrue": {
			givenGhLabel:   newGitHubLabel(label, description, color),
			givenRepoLabel: newDomainLabel(t, label, "different", color),
			expectedResult: true,
		},
		"GivenDifferentColor_ThenExpectTrue": {
			givenGhLabel:   newGitHubLabel(label, description, color),
			givenRepoLabel: newDomainLabel(t, label, description, "FFFFFF"),
			expectedResult: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &GhRemote{}
			result := p.hasLabelChanged(tt.givenGhLabel, tt.givenRepoLabel)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func newGitHubLabel(name, description, color string) *github.Label {
	return &github.Label{
		Name:        &name,
		Description: &description,
		Color:       &color,
	}
}

func newDomainLabel(t *testing.T, name, description, color string) domain.Label {
	label := domain.Label{
		Name:        name,
		Description: description,
	}
	err := label.SetColor(ColorConverter{}.ConvertToEntity(color))
	require.NoError(t, err)
	return label
}

/*func TestProvider_findMatchingGhLabel(t *testing.T) {
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
			p := &GhRemote{}
			result := p.findMatchingGhLabel(tt.givenGhLabels, tt.givenRepoLabelForComparison)
			if tt.expectedLabelIndex >= 0 {
				assert.Equal(t, tt.givenGhLabels[tt.expectedLabelIndex], result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
*/

func TestGhRemote_updateLabelCache(t *testing.T) {
	u, err := url.Parse("https://github.com/ccremer/greposync")
	require.NoError(t, err)
	gitUrl := domain.FromURL(u)
	labelName := "label"
	givenLabelToUpdate := LabelConverter{}.ConvertFromEntity(domain.Label{
		Name:        labelName,
		Description: "new description",
	})

	tests := map[string]struct {
		givenLabelCache map[*domain.GitURL][]*github.Label
	}{
		"GivenNonExistingKey_WhenAddingLabel_ThenCreateNewList": {
			givenLabelCache: map[*domain.GitURL][]*github.Label{},
		},
		"GivenExistingKey_WhenLabelIsExisting_ThenUpdateInPlace": {
			givenLabelCache: map[*domain.GitURL][]*github.Label{
				gitUrl: {
					newGitHubLabel(labelName, "old description", "ffffff"),
				},
			},
		},
		"GivenExistingKey_WhenLabelIsNonExisting_ThenAppendToList": {
			givenLabelCache: map[*domain.GitURL][]*github.Label{
				gitUrl: {
					newGitHubLabel("another", "description", "ffffff"),
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewRemote()
			r.labelCache = tt.givenLabelCache
			r.updateLabelCache(gitUrl, givenLabelToUpdate)

			result := r.labelCache[gitUrl]
			require.NotEmpty(t, result)
			assert.Contains(t, result, givenLabelToUpdate)
		})
	}
}

func TestGhRemote_removeLabelFromCache(t *testing.T) {
	u, err := url.Parse("https://github.com/ccremer/greposync")
	require.NoError(t, err)
	gitUrl := domain.FromURL(u)
	labelName := "label"
	givenLabelToRemove := LabelConverter{}.ConvertFromEntity(domain.Label{
		Name:        labelName,
		Description: "new description",
	})

	tests := map[string]struct {
		givenLabelCache map[*domain.GitURL][]*github.Label
		expectedLen     int
	}{
		"GivenNonExistingKey_ThenDoNothing": {
			givenLabelCache: map[*domain.GitURL][]*github.Label{},
		},
		"GivenExistingKey_WhenListIsEmpty_ThenDoNothing": {
			givenLabelCache: map[*domain.GitURL][]*github.Label{
				gitUrl: {},
			},
		},
		"GivenExistingKey_WhenListContainsLabel_ThenRemoveThatLabel": {
			givenLabelCache: map[*domain.GitURL][]*github.Label{
				gitUrl: {
					newGitHubLabel(labelName, "old description", "ffffff"),
				},
			},
		},
		"GivenExistingKey_WhenListContainsOtherLabel_ThenKeepOthers": {
			givenLabelCache: map[*domain.GitURL][]*github.Label{
				gitUrl: {
					newGitHubLabel(labelName, "old description", "ffffff"),
					newGitHubLabel("another", "description", "ffffff"),
				},
			},
			expectedLen: 1,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewRemote()
			r.labelCache = tt.givenLabelCache
			r.removeLabelFromCache(gitUrl, givenLabelToRemove)

			result := r.labelCache[gitUrl]
			assert.NotContains(t, result, givenLabelToRemove)
			assert.Len(t, result, tt.expectedLen)
		})
	}
}
