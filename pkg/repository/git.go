package repository

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

var (
	GitBin = "git"
)

func (s *Repository) execGitCommand(args ...string) (string, string, error) {
	cmd := exec.Command(GitBin, args...)
	if s.dirExists(s.GitConfig.Dir) {
		cmd.Dir = s.GitConfig.Dir
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func (s *Repository) logArgs(args ...string) []string {
	s.log.InfoF("%s %s", GitBin, strings.Join(args, " "))
	return args
}

func (s *Repository) dirExists(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}
