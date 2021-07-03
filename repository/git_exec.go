package repository

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	pipeline "github.com/ccremer/go-command-pipeline"
)

var (
	GitBin = "git"
)

func (s *Service) execGitCommand(args ...string) (string, string, error) {
	cmd := exec.Command(GitBin, args...)
	if s.DirExists(s.Config.Dir) {
		cmd.Dir = s.Config.Dir
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func (s *Service) logArgs(args ...string) []string {
	s.p.InfoF("%s %s", GitBin, strings.Join(args, " "))
	return args
}

func (s *Service) toResult(err error, msg string) pipeline.Result {
	return pipeline.Result{
		Err: fmt.Errorf("%w: %s", err, msg),
	}
}
