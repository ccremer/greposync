package repositorystore

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ccremer/greposync/domain"
)

var GitBinary = "git"

func execGitCommand(rootDir domain.Path, args []string) (stdOut, stdErr string, cmdErr error) {
	cmd := exec.Command(GitBinary, args...)
	if rootDir.DirExists() {
		cmd.Dir = rootDir.String()
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func mergeWithStdErr(err error, stderr string) error {
	return fmt.Errorf("%w: %s", err, stderr)
}

func hasRemoteBranch(repository *domain.GitRepository, branch string) (bool, error) {
	out, stderr, err := execGitCommand(repository.RootDir, []string{"branch", "-r", "--list"})
	return parseBranch(err, stderr, out, branch)
}

func hasLocalBranch(repository *domain.GitRepository, branch string) (bool, error) {
	out, stderr, err := execGitCommand(repository.RootDir, []string{"branch", "--list"})
	return parseBranch(err, stderr, out, branch)
}

func parseBranch(err error, stderr string, out string, branch string) (bool, error) {
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
func HasCommitsBetween(repository *domain.GitRepository, rootBranch, headBranch string) (bool, error) {
	if rootBranch == "" {
		return false, fmt.Errorf("%w: rootBranch cannot be empty", domain.ErrInvalidArgument)
	}
	out, _, err := execGitCommand(repository.RootDir, []string{"log", fmt.Sprintf("%s..%s", rootBranch, headBranch), "--oneline"})
	return out != "", err
}

// GetDefaultBranch returns the name of the default branch in origin.
// Returns an error if either Git command failed or if no default branch could be detected.
func GetDefaultBranch(repository *domain.GitRepository) (string, error) {
	out, stderr, err := execGitCommand(repository.RootDir, []string{"remote", "show", "origin"})
	if err != nil {
		return "master", mergeWithStdErr(err, stderr)
	}
	lines := strings.Split(out, "\n")
	head := "HEAD branch: "
	for _, line := range lines {
		str := strings.TrimSpace(line)
		if strings.Contains(str, head) {
			return strings.TrimPrefix(str, head), nil
		}
	}
	return "master", fmt.Errorf("no default branch determined")
}
