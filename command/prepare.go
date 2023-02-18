package command

import (
	"bytes"
	"fmt"
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/models/github_event"
	"github.com/cocov-ci/coverage-reporter/tracking"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Prepare(ctx *cli.Context) error {
	log := zap.L()
	log.Info("Cocov Coverage Reporter")
	log.Info("Copyright (c) 2022-2023 - The Cocov Authors")
	token := getToken(ctx)

	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("reading pwd: %w", err)
	}

	log.Debug("Started hashing process...")
	thence := time.Now()
	files, err := tracking.FilesOn(pwd)
	if err != nil {
		return err
	}
	log.Debug("Hashing finished", zap.Duration("duration", time.Since(thence)), zap.Int("length", len(files)))

	sha := ctx.String("commitish")
	if sha == "" {
		ok, ev, err := github_event.Lookup()
		if err != nil {
			log.Error("GitHub Event is present, but could not be parsed", zap.Error(err))
		}
		if ok && ev.PullRequest != nil && ev.PullRequest.Head != nil {
			sha = ev.PullRequest.Head.Sha
		}
	}

	// If even after all attempts, sha is empty...
	if sha == "" {
		c, err := ensureCommit(pwd, ctx)
		if err != nil {
			log.Error("Failed obtaining commit. All attempts failed. Last error follows.", zap.Error(err))
			return err
		}
		sha = c
	}

	metaEntry := &meta.Metadata{
		Files: files,
		Pwd:   pwd,
		Sha:   sha,
		Token: token,
	}

	if os.Getenv("COCOV_REPORTER_DEBUG_METADATA") == "true" {
		log.Debug("Obtained metadata directory", zap.String("path", meta.MetadataDir(token)))
	}

	if err = os.RemoveAll(meta.MetadataDir(token)); err != nil {
		return err
	}

	if err = os.MkdirAll(meta.MetadataDir(token), 0750); err != nil {
		return err
	}

	return meta.StoreMetadata(token, metaEntry)
}

func ensureCommit(pwd string, ctx *cli.Context) (string, error) {
	sha := ctx.String("commitish")
	sha = strings.TrimSpace(sha)
	if len(sha) != 0 {
		return sha, nil
	}

	cmd, err := exec.LookPath("git")
	if err != nil {
		return "", err
	}

	output := bytes.Buffer{}

	execCmd := exec.Cmd{
		Path:   cmd,
		Args:   []string{cmd, "rev-parse", "HEAD"},
		Env:    os.Environ(),
		Dir:    pwd,
		Stdin:  nil,
		Stdout: &output,
		Stderr: &output,
	}
	if err = execCmd.Start(); err != nil {
		return "", err
	}
	if err = execCmd.Wait(); err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("process `git' exited with status %d: %s", e.ExitCode(), output.String())
		}
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}
