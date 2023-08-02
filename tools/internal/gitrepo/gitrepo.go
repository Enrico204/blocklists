// Package gitrepo handles basic Git operations in a directory.
//
// Currently, the package is using the command line `git`. However, you should not rely on this implementation detail,
// as hopefully it will change in the future.
package gitrepo

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strings"
	"time"
)

const GitOperationsTimeout = 60 * time.Second

var ErrEmptyRepo = errors.New("empty repository")

type GitRepo struct {
	url    string
	dir    string
	logger *zap.SugaredLogger
}

func (r *GitRepo) Dir() string {
	return r.dir
}

func (r *GitRepo) Exists() (bool, error) {
	finfo, err := os.Stat(r.dir)
	switch {
	case err != nil && os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	case !finfo.IsDir():
		return false, errors.New("expected repository directory, other type found instead")
	}
	return true, nil
}

func (r *GitRepo) Delete() error {
	return os.RemoveAll(r.dir)
}

func (r *GitRepo) Clone() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GitOperationsTimeout)
	defer cancel()

	args := []string{
		"clone",
		"-c", "core.compression=1",
		"--depth", "1",
		r.url,
		r.dir,
	}
	cmd := exec.CommandContext(ctx, "git", args...) // nolint:gosec
	cmd.Env = append(os.Environ(), "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -F /etc/ssh/ssh_config")

	var stderrbuf strings.Builder
	cmd.Stderr = &stderrbuf

	stdout, err := cmd.Output()
	if err != nil {
		var err2 *exec.ExitError
		if errors.As(err, &err2) {
			stdout = append(stdout, err2.Stderr...)
		}
	}
	return StripANSI(string(stdout) + stderrbuf.String()), err
}

func (r *GitRepo) Reset(commitID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GitOperationsTimeout)
	defer cancel()
	logger := r.logger.With("git-commit", commitID)

	cmd := exec.CommandContext(ctx, "git", "reset", "--hard", commitID)
	cmd.Dir = r.dir
	cmd.Env = append(os.Environ(), "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -F /etc/ssh/ssh_config")
	stdout, err := cmd.Output()
	if err != nil {
		var err2 *exec.ExitError
		if errors.As(err, &err2) {
			logger.Errorw("git reset error: ", string(err2.Stderr))
		} else {
			logger.Errorw("git reset error", "err", err)
		}
	} else {
		logger.Debugw("git reset", "output", string(stdout))
	}
	return StripANSI(string(stdout)), err
}

func (r *GitRepo) Clean() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GitOperationsTimeout)
	defer cancel()
	logger := r.logger

	cmd := exec.CommandContext(ctx, "git", "clean", "-f", "-d", "-x")
	cmd.Dir = r.dir
	cmd.Env = append(os.Environ(), "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -F /etc/ssh/ssh_config")
	stdout, err := cmd.Output()
	if err != nil {
		var err2 *exec.ExitError
		if errors.As(err, &err2) {
			logger.Errorw("git clean error: ", string(err2.Stderr))
		} else {
			logger.Errorw("git clean", "err", err)
		}
	} else {
		logger.Debugw("git clean", "output", string(stdout))
	}
	return StripANSI(string(stdout)), err
}

func (r *GitRepo) Pull() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GitOperationsTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "pull", "-p", "--all", "--ff-only")
	cmd.Dir = r.dir
	cmd.Env = append(os.Environ(), "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -F /etc/ssh/ssh_config")
	stdout, err := cmd.Output()
	if err != nil {
		var err2 *exec.ExitError
		if errors.As(err, &err2) {
			stdout = append(stdout, err2.Stderr...)
		}
	}
	return StripANSI(string(stdout)), err
}

func (r *GitRepo) GetHash() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GitOperationsTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--verify", "HEAD", "--short")
	cmd.Dir = r.dir
	cmd.Env = append(os.Environ(), "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -F /etc/ssh/ssh_config")
	stdout, err := cmd.Output()
	if err != nil {
		var err2 *exec.ExitError
		if errors.As(err, &err2) {
			if err2.ExitCode() == 128 {
				// Empty repository
				return "", ErrEmptyRepo
			}
			stdout = append(stdout, err2.Stderr...)
		}
	}
	return strings.Trim(StripANSI(string(stdout)), " \n\r"), err
}
