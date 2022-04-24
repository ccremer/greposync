package repositorystore

import (
	"testing"

	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging/loggingtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	giturls "github.com/whilp/git-urls"
)

func TestTestRepositoryStore_FetchGitRepositories(t *testing.T) {
	tests := map[string]struct {
		prepare       func(t *testing.T, s *TestRepositoryStore)
		expectedList  []*domain.GitRepository
		expectedError string
	}{
		"GivenTestDir1_WhenParentContainsFiles_ThenExpectOnlyDirs": {
			prepare: func(t *testing.T, s *TestRepositoryStore) {
				s.ParentDir = "testdata/testcase1"
				s.TestOutputRootDir = "testdata/.tests/testcase1"
			},
			expectedList: []*domain.GitRepository{
				{RootDir: "testdata/.tests/testcase1/fakerepo1", URL: testURL(t, "file://testdata/.tests/testcase1/fakerepo1")},
				{RootDir: "testdata/.tests/testcase1/fakerepo2", URL: testURL(t, "file://testdata/.tests/testcase1/fakerepo2")},
			},
		},
		"GivenTestDir2_WhenParentContainsHiddenDirs_ThenSkipHiddenDir": {
			prepare: func(t *testing.T, s *TestRepositoryStore) {
				s.ParentDir = "testdata/testcase2"
				s.TestOutputRootDir = "testdata/.tests/testcase2"
			},
			expectedList: []*domain.GitRepository{
				{RootDir: "testdata/.tests/testcase2/fakerepo2", URL: testURL(t, "file://testdata/.tests/testcase2/fakerepo2")},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewTestRepositoryStore(NewRepositoryStoreInstrumentation(loggingtest.NewTestingLogger(t)))
			tt.prepare(t, s)
			result, err := s.FetchGitRepositories()
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError, "expected fetch error")
				return
			}
			require.NoError(t, err, "unexpected fetch error")
			assert.Equal(t, tt.expectedList, result)
		})
	}
}

func TestTestRepositoryStore_Diff(t *testing.T) {
	tests := map[string]struct {
		prepare         func(t *testing.T, s *TestRepositoryStore)
		givenRepository *domain.GitRepository
		expectedDiff    string
		expectedError   string
	}{
		"GivenTestDir1_WhenContentIsChanged_ThenExpectDiff": {
			prepare: func(t *testing.T, s *TestRepositoryStore) {
				s.ParentDir = "testdata/diff1"
				s.TestOutputRootDir = "testdata/diff1/actual"
			},
			givenRepository: &domain.GitRepository{
				RootDir: domain.NewFilePath("testdata", "diff1", "actual", "fakerepo"),
				URL:     testURL(t, "fakerepo"),
			},
			expectedDiff: "--- actual:testdata/diff1/actual/fakerepo/file.txt\n+++ expected:testdata/diff1/fakerepo/file.txt\n@@ -1 +1 @@\n-Actual Line\n+Expected Line\n",
		},
		"GivenTestDir2_WhenContentIsSame_ThenExpectNoDiffWithoutError": {
			prepare: func(t *testing.T, s *TestRepositoryStore) {
				s.ParentDir = "testdata/diff2"
				s.TestOutputRootDir = "testdata/diff2/actual"
			},
			givenRepository: &domain.GitRepository{
				RootDir: domain.NewFilePath("testdata", "diff2", "fakerepo"),
				URL:     testURL(t, "fakerepo"),
			},
			expectedDiff: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewTestRepositoryStore(NewRepositoryStoreInstrumentation(loggingtest.NewTestingLogger(t)))
			tt.prepare(t, s)
			result, err := s.Diff(tt.givenRepository, domain.DiffOptions{})
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError, "expected diff error")
				return
			}
			require.NoError(t, err, "unexpected diff error")
			assert.Contains(t, result, tt.expectedDiff)
		})
	}
}

func testURL(t *testing.T, url string) *domain.GitURL {
	parsed, err := giturls.Parse(url)
	require.NoError(t, err)
	return domain.FromURL(parsed)
}
