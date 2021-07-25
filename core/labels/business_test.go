package labels

import (
	"errors"
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
		givenLabels   []core.GitRepositoryLabel
		expectedErr   bool
		expectedCalls int
	}{
		"NoLabels": {},
		"ActiveLabel": {
			givenLabels: []core.GitRepositoryLabel{
				newFakeLabel("active", false),
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

			repoFake := createRepoFake(core.GitRepositoryConfig{Url: gu}, tt.givenLabels)
			hostingFake := createHostingFake(nil)
			err := s.createOrUpdateLabels(repoFake, hostingFake)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCalls, hostingFake.CreateOrUpdateLabelsForRepoCallCount())
			if tt.expectedCalls > 0 {
				gitUrl, result := hostingFake.CreateOrUpdateLabelsForRepoArgsForCall(0)
				assert.Equal(t, gu, gitUrl)
				assert.Equal(t, tt.givenLabels, result)
			}
		})
	}
}

func TestLabelService_deleteLabels(t *testing.T) {
	labelTests := map[string]struct {
		givenLabels   []core.GitRepositoryLabel
		expectedErr   bool
		expectedCalls int
	}{
		"NoLabels": {},
		"DeadLabel": {
			givenLabels: []core.GitRepositoryLabel{
				newFakeLabel("dead", true),
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

			repoFake := createRepoFake(core.GitRepositoryConfig{Url: gu}, tt.givenLabels)
			hostingFake := createHostingFake(nil)
			err := s.deleteLabels(repoFake, hostingFake)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCalls, hostingFake.DeleteLabelsForRepoCallCount())
			if tt.expectedCalls > 0 {
				gitUrl, result := hostingFake.DeleteLabelsForRepoArgsForCall(0)
				assert.Equal(t, gu, gitUrl)
				assert.Equal(t, tt.givenLabels, result)
			}
		})
	}
}

var filterLabelTests = map[string]struct {
	givenLabels          []core.GitRepositoryLabel
	expectedActiveLabels []core.GitRepositoryLabel
	expectedDeadLabels   []core.GitRepositoryLabel
}{
	"GivenEmptyList": {
		givenLabels:          []core.GitRepositoryLabel{},
		expectedActiveLabels: []core.GitRepositoryLabel{},
	},
	"GivenNilList": {
		expectedActiveLabels: []core.GitRepositoryLabel{},
	},
	"GivenActiveLabel": {
		givenLabels: []core.GitRepositoryLabel{
			newFakeLabel("fake", false),
		},
		expectedActiveLabels: []core.GitRepositoryLabel{
			newFakeLabel("fake", false),
		},
	},
	"GivenDeadLabel": {
		givenLabels: []core.GitRepositoryLabel{
			newFakeLabel("fake", true),
		},
		expectedDeadLabels: []core.GitRepositoryLabel{
			newFakeLabel("fake", true),
		},
	},
}

func TestLabelService_filterActiveLabels(t *testing.T) {
	for name, tt := range filterLabelTests {
		t.Run(name, func(t *testing.T) {
			result := filterActiveLabels(tt.givenLabels)
			assert.Len(t, result, len(tt.expectedActiveLabels))
			for i, expectedLabel := range tt.expectedActiveLabels {
				assert.Equal(t, expectedLabel.GetName(), result[i].GetName())
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
				assert.Equal(t, expectedLabel.GetName(), result[i].GetName())
			}
		})
	}
}

func createHostingFake(returnErr error) *corefakes.FakeGitHostingFacade {
	return &corefakes.FakeGitHostingFacade{
		CreateOrUpdateLabelsForRepoStub: func(gu *core.GitUrl, labels []core.GitRepositoryLabel) error {
			return returnErr
		},
		DeleteLabelsForRepoStub: func(gu *core.GitUrl, labels []core.GitRepositoryLabel) error {
			return returnErr
		},
		InitializeStub: func() error {
			return returnErr
		},
	}
}

func createRepoFake(cfg core.GitRepositoryConfig, labels []core.GitRepositoryLabel) *corefakes.FakeGitRepositoryFacade {
	return &corefakes.FakeGitRepositoryFacade{
		GetConfigStub: func() core.GitRepositoryConfig {
			return cfg
		},
		GetLabelsStub: func() []core.GitRepositoryLabel {
			return labels
		},
	}
}

func newFakeLabel(name string, delete bool) core.GitRepositoryLabel {
	return &corefakes.FakeGitRepositoryLabel{
		GetNameStub: func() string {
			return name
		},
		IsBoundForDeletionStub: func() bool {
			return delete
		},
	}
}

func createURl(t *testing.T) *core.GitUrl {
	u, err := url.Parse("https://github.com/ccremer/greposync")
	require.NoError(t, err)
	gu := core.GitUrl(*u)
	return &gu
}

func TestLabelService_initHostingAPIs(t *testing.T) {
	tests := map[string]struct {
		givenProvider   core.GitHostingProvider
		expectErrString string
		expectedCalls   int
	}{
		"GivenSupportedProviders_ThenInitHostingApi": {
			givenProvider: "provider",
			expectedCalls: 1,
		},
		"GivenSupportedProvider_WhenError_ThenExpectError": {
			givenProvider:   "provider",
			expectedCalls:   1,
			expectErrString: "failed",
		},
		"GivenUnSupportedProviders_ThenIgnore": {
			givenProvider: "unsupported",
			expectedCalls: 0,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var returnErr error
			if tt.expectErrString != "" {
				returnErr = errors.New(tt.expectErrString)
			}
			hostingFake := createHostingFake(returnErr)

			s := &LabelService{
				repoFacades: []core.GitRepositoryFacade{
					createRepoFake(core.GitRepositoryConfig{
						Provider: tt.givenProvider,
					}, nil),
				},
				repoProvider: &corefakes.FakeManagedRepoProvider{
					GetSupportedGitHostingProvidersStub: func() map[core.GitHostingProvider]core.GitHostingFacade {
						providers := map[core.GitHostingProvider]core.GitHostingFacade{
							"provider": hostingFake,
						}
						return providers
					},
				},
				log: printer.New(),
			}
			err := s.initHostingAPIs()
			if tt.expectErrString != "" {
				require.EqualError(t, err, tt.expectErrString)
				assert.Equal(t, tt.expectedCalls, hostingFake.InitializeCallCount())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCalls, hostingFake.InitializeCallCount())
		})
	}
}
