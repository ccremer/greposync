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
	out, stderr, err := execGitCommand(repository.RootDir, s.logArgs(args))
	if err != nil {
		s.log.InfoF(out)
		return mergeWithStdErr(err, stderr)
	}
	s.log.DebugF(out)
	return nil
}

func (s *RepositoryStore) Add(repository *domain.GitRepository) error {
	out, stderr, err := execGitCommand(repository.RootDir, s.logArgs([]string{"add", "-A"}))
	if err != nil {
		return mergeWithStdErr(err, stderr)
	}
	s.log.DebugF(out)
	return nil
}

func (s *RepositoryStore) Diff(repository *domain.GitRepository) (string, error) {
	out, stderr, err := execGitCommand(repository.RootDir, []string{"diff", "HEAD~1"})
	if err != nil {
		if strings.Contains(stderr, "ambiguous argument 'HEAD~1': unknown revision or path not in the working tree.") {
			s.log.InfoF("This is the first commit, no diff available.")
			return "", nil
		}
		return "", mergeWithStdErr(err, stderr)
	}
	return out, nil
}

func (s *RepositoryStore) IsDirty(repository *domain.GitRepository) bool {
	out, stderr, err := execGitCommand(repository.RootDir, []string{"status", "--short"})
	if err != nil {
		s.log.WarnF("Could not determine working tree status: %s: %w", stderr, err)
		return true
	}
	if out == "" {
		s.log.InfoF("Nothing to commit, working tree clean")
		return false
	}
	return true
}
