package repositorystore

import (
	"fmt"
	"os"
	"strings"

	"github.com/ccremer/greposync/domain"
)

func (s *RepositoryStore) Commit(repository *domain.GitRepository, options domain.CommitOptions) error {
	f, err := os.CreateTemp("", "COMMIT_MSG_")
	if err != nil {
		return fmt.Errorf("failed to create temporary commit message file: %w", err)
	}
	defer s.deleteFileHandler(f)

	// Write commit message
	_, err = fmt.Fprint(f, options.Message)
	if err != nil {
		return err
	}

	args := []string{"commit", "-a", "-F", f.Name()}

	// Try to figure out if amend makes sense
	if options.Amend {
		if hasCommits, err := HasCommitsBetween(repository, repository.DefaultBranch, repository.CommitBranch); err != nil {
			return err
		} else if hasCommits {
			args = append(args, "--amend")
		}
	}

	// Commit
	out, stderr, err := execGitCommand(repository.RootDir, s.instrumentation.logGitArguments(repository, 0, args))
	if err != nil {
		s.instrumentation.logInfo(repository, out)
		return mergeWithStdErr(err, stderr)
	}
	s.instrumentation.logDebugInfo(repository, out)
	return nil
}

func (s *RepositoryStore) Add(repository *domain.GitRepository) error {
	out, stderr, err := execGitCommand(repository.RootDir, s.instrumentation.logGitArguments(repository, 0, []string{"add", "-A"}))
	if err != nil {
		return mergeWithStdErr(err, stderr)
	}
	s.instrumentation.logDebugInfo(repository, out)
	return nil
}

func (s *RepositoryStore) Diff(repository *domain.GitRepository, options domain.DiffOptions) (string, error) {
	args := []string{"diff", "HEAD~1"}
	if options.WorkDirToHEAD {
		args = []string{"diff", "HEAD"}
	}
	out, stderr, err := execGitCommand(repository.RootDir, args)
	if err != nil {
		if strings.Contains(stderr, "ambiguous argument 'HEAD~1': unknown revision or path not in the working tree.") {
			s.instrumentation.logInfo(repository, "This is the first commit, no diff available.")
			return "", nil
		}
		return "", mergeWithStdErr(err, stderr)
	}
	return out, nil
}

func (s *RepositoryStore) IsDirty(repository *domain.GitRepository) bool {
	out, stderr, err := execGitCommand(repository.RootDir, []string{"status", "--short"})
	if err != nil {
		s.instrumentation.logInfo(repository, stderr)
		return true
	}
	if out == "" {
		s.instrumentation.logInfo(repository, "Nothing to commit, working tree clean")
		return false
	}
	return true
}
