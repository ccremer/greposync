package domain

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var GitBinary = "git"

type GitRepository struct {
	RootDir      Path
	URL          *GitURL
	PullRequests []PullRequest
	Labels       LabelSet

	CommitBranch  string
	DefaultBranch string
}

func NewGitRepository(u *GitURL, labels LabelSet) (*GitRepository, error) {
	r := &GitRepository{
		URL: u,
	}
	if err := r.validateLabels(labels); hasFailed(err) {
		return r, err
	}
	r.Labels = labels
	return r, nil
}

type CommitOptions struct {
	Amend         bool
	CommitMessage string
}

func (r *GitRepository) validateLabels(labels LabelSet) error {
	return firstOf(labels.CheckForEmptyLabelNames(), labels.CheckForDuplicates())
}

func (r *GitRepository) SetLabels(labels LabelSet) error {
	if err := r.validateLabels(labels); hasFailed(err) {
		return err
	}
	r.Labels = labels
	return nil
}

func (r *GitRepository) Commit(logger GitLogger, options CommitOptions) error {
	f, err := os.CreateTemp("", "COMMIT_MSG_")
	if err != nil {
		return fmt.Errorf("failed to create temporary commit message file: %w", err)
	}
	defer r.deleteFileHandler(f)

	// Write commit message
	_, err = fmt.Fprint(f, options.CommitMessage)
	if err != nil {
		return err
	}

	args := []string{"commit", "-a", "-F", f.Name()}

	// Try to figure out if amend makes sense
	if options.Amend {
		if hasCommits, err := r.HasCommitsBetween(r.DefaultBranch, r.CommitBranch); err != nil {
			return err
		} else if hasCommits {
			args = append(args, "--amend")
		}
	}

	// Commit
	return r.executeGitWithLogger(logger, args...)
}

func (r *GitRepository) Push(logger GitLogger, force bool) error {
	args := []string{"push", "origin", r.CommitBranch}
	if force {
		args = append(args, "--force")
	}
	return r.executeGitWithLogger(logger, args...)
}

func (r *GitRepository) HasRemoteBranch(branch string) (bool, error) {
	out, stderr, err := r.execGitCommand([]string{"branch", "-r", "--list"})
	if err != nil {
		return false, errors.New(stderr)
	}
	branches := strings.Split(out, "\n")
	for _, line := range branches {
		if strings.Contains(strings.TrimSpace(line), branch) {
			return true, nil
		}
	}
	return false, nil
}

// HasCommitsBetween returns true if there are commits in the given revision range.
// If headBranch is empty, "HEAD" is used.
// Returns ErrInvalidArgument if rootBranch is empty.
// Returns errors in all other Git failures.
func (r *GitRepository) HasCommitsBetween(rootBranch, headBranch string) (bool, error) {
	if rootBranch == "" {
		return false, fmt.Errorf("%w: rootBranch cannot be empty", ErrInvalidArgument)
	}
	out, _, err := r.execGitCommand([]string{"log", fmt.Sprintf("%s..%s", rootBranch, headBranch), "--oneline"})
	return out != "", err
}

func (r *GitRepository) executeGitWithLogger(logger GitLogger, args ...string) error {
	if !r.RootDir.DirExists() {
		return fmt.Errorf("dir doesnt exist: %s", r.RootDir.String())
	}
	logger.LogArgs(args)
	out, stderr, err := r.execGitCommand(args)
	logger.LogOutput(out, stderr)
	return err
}

func (r *GitRepository) deleteFileHandler(file *os.File) {
	_ = file.Close()
	_ = os.Remove(file.Name())
}
