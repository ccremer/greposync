package labels

import (
	"net/url"
	"testing"

	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/core/corefakes"
	"github.com/ccremer/greposync/printer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelService_createOrUpdateLabels(t *testing.T) {
	labelTests := map[string]struct {
		givenLabels   []core.Label
		expectedErr   bool
		expectedCalls int
	}{
		"NoLabels": {},
		"ActiveLabel": {
			givenLabels: []core.Label{
				newFakeLabel(false),
			},
			expectedCalls: 1,
		},
	}

	for name, tt := range labelTests {
		t.Run(name, func(t *testing.T) {
			gu := createURl(t)
			s := &LabelService{
				log: printer.New(),
			}

			repoFake := createRepoFake(core.GitRepositoryProperties{URL: gu}, tt.givenLabels)
			err := s.createOrUpdateLabels(repoFake)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.expectedCalls > 0 {
				result := tt.givenLabels[0].(*corefakes.FakeLabel).EnsureCallCount()
				assert.Equal(t, tt.expectedCalls, result)
			}
		})
	}
}

func TestLabelService_deleteLabels(t *testing.T) {
	labelTests := map[string]struct {
		givenLabels   []core.Label
		expectedErr   bool
		expectedCalls int
	}{
		"NoLabels": {},
		"DeadLabel": {
			givenLabels: []core.Label{
				newFakeLabel(true),
			},
			expectedCalls: 1,
		},
	}

	for name, tt := range labelTests {
		t.Run(name, func(t *testing.T) {

			gu := createURl(t)
			s := &LabelService{
				log: printer.New(),
			}

			repoFake := createRepoFake(core.GitRepositoryProperties{URL: gu}, tt.givenLabels)
			err := s.deleteLabels(repoFake)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.expectedCalls > 0 {
				result := tt.givenLabels[0].(*corefakes.FakeLabel).DeleteCallCount()
				assert.Equal(t, tt.expectedCalls, result)
			}
		})
	}
}

var filterLabelTests = map[string]struct {
	givenLabels          []core.Label
	expectedActiveLabels []core.Label
	expectedDeadLabels   []core.Label
}{
	"GivenEmptyList": {
		givenLabels:          []core.Label{},
		expectedActiveLabels: []core.Label{},
	},
	"GivenNilList": {
		expectedActiveLabels: []core.Label{},
	},
	"GivenActiveLabel": {
		givenLabels: []core.Label{
			newFakeLabel(false),
		},
		expectedActiveLabels: []core.Label{
			newFakeLabel(false),
		},
	},
	"GivenDeadLabel": {
		givenLabels: []core.Label{
			newFakeLabel(true),
		},
		expectedDeadLabels: []core.Label{
			newFakeLabel(true),
		},
	},
}

func TestLabelService_filterActiveLabels(t *testing.T) {
	for name, tt := range filterLabelTests {
		t.Run(name, func(t *testing.T) {
			result := filterActiveLabels(tt.givenLabels)
			assert.Len(t, result, len(tt.expectedActiveLabels))
			for i, expectedLabel := range tt.expectedActiveLabels {
				assert.Equal(t, expectedLabel.IsInactive(), result[i].IsInactive())
			}
		})
	}
}

func TestLabelService_filterDeadLabels(t *testing.T) {
	for name, tt := range filterLabelTests {
		t.Run(name, func(t *testing.T) {
			result := filterDeadLabels(tt.givenLabels)
			assert.Len(t, result, len(tt.expectedDeadLabels))
			for i, expectedLabel := range tt.expectedDeadLabels {
				assert.Equal(t, expectedLabel.IsInactive(), result[i].IsInactive())
			}
		})
	}
}

func createRepoFake(cfg core.GitRepositoryProperties, labels []core.Label) *corefakes.FakeGitRepository {
	return &corefakes.FakeGitRepository{
		GetConfigStub: func() core.GitRepositoryProperties {
			return cfg
		},
		GetLabelsStub: func() []core.Label {
			return labels
		},
	}
}

func newFakeLabel(delete bool) core.Label {
	return &corefakes.FakeLabel{
		IsInactiveStub: func() bool {
			return delete
		},
	}
}

func createURl(t *testing.T) *core.GitURL {
	u, err := url.Parse("https://github.com/ccremer/greposync")
	require.NoError(t, err)
	gu := core.GitURL(*u)
	return &gu
}
