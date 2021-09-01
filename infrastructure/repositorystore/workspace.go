package repositorystore

import (
	"errors"
	"strings"

	"github.com/ccremer/greposync/domain"
)

func (s *RepositoryStore) Clone(repository *domain.GitRepository) error {
	if repository.RootDir.DirExists() {
		return errors.New("clone exists already")
	}
	dir := repository.RootDir.String()
	gitURL := repository.URL

	// Don't want to expose credentials in the log, so we're not calling logArgs().
	s.log.InfoF("%s %s", GitBinary, strings.Join([]string{"clone", gitURL.Redacted(), dir}, " "))

	out, stderr, err := execGitCommand(repository.RootDir, []string{"clone", gitURL.String(), dir})
	if err != nil {
		return mergeWithStdErr(err, stderr)
	}
	s.log.PrintF(out)
	return nil
}

func (s *RepositoryStore) Checkout(repository *domain.GitRepository) error {
	args := []string{"checkout"}
	if localExists, err := hasLocalBranch(repository, repository.CommitBranch); err != nil {
		return err
	} else if !localExists {
		// Checkout to new branch
		args = append(args, "-t", "-b")
	}
	args = append(args, repository.CommitBranch)

	out, stderr, err := execGitCommand(repository.RootDir, s.logArgs(args))
	if err != nil {
		return mergeWithStdErr(err, stderr)
	}
	s.log.DebugF(out)
	return nil
}

func (s *RepositoryStore) Fetch(repository *domain.GitRepository) error {
	out, stderr, err := execGitCommand(repository.RootDir, s.logArgs([]string{"fetch"}))
	if err != nil {
		return mergeWithStdErr(err, stderr)
	}
	if out != "" {
		s.log.InfoF(out)
	}
	return nil
}

func (s *RepositoryStore) Reset(repository *domain.GitRepository) error {
	out, stderr, err := execGitCommand(repository.RootDir, s.logArgs([]string{"reset", "--hard"}))
	if err != nil {
		return mergeWithStdErr(err, stderr)
	}
	s.log.DebugF(out)
	return nil
}

func (s *RepositoryStore) Pull(repository *domain.GitRepository) error {
	exists, err := hasRemoteBranch(repository, repository.CommitBranch)
	if err != nil {
		return err
	}
	if exists {
		out, stderr, err := execGitCommand(repository.RootDir, s.logArgs([]string{"pull"}))
		if err != nil {
			return mergeWithStdErr(err, stderr)
		}
		s.log.DebugF(out)
	}
	return nil
}

func (s *RepositoryStore) Push(repository *domain.GitRepository, options domain.PushOptions) error {
	args := []string{"push", "origin", repository.CommitBranch}
	if options.Force {
		args = append(args, "--force")
	}
	out, stderr, err := execGitCommand(repository.RootDir, s.logArgs(args))
	if err != nil {
		return mergeWithStdErr(err, stderr)
	}
	s.log.DebugF(out)
	return nil
}
