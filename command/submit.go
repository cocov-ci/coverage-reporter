package command

import (
	"bytes"
	"fmt"
	"github.com/cocov-ci/coverage-reporter/formats"
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/tracking"
	"github.com/levigross/grequests"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strings"
)

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

func Submit(ctx *cli.Context) error {
	log := zap.L()

	token := getToken(ctx)
	runMeta, err := meta.ReadMetadata(token)
	if err != nil {
		return err
	}

	commit, err := ensureCommit(runMeta.Pwd, ctx)
	if err != nil {
		return err
	}

	currentFiles, err := tracking.FilesOn(runMeta.Pwd)
	if err != nil {
		return err
	}

	diffs := tracking.DiffFiles(runMeta.Files, currentFiles)
	toPush, err := formats.AutoFind(diffs, runMeta)
	if err != nil {
		return err
	}

	log.Debug("Now submitting data")

	url := strings.TrimPrefix(ctx.String("url"), "/")
	resp, err := grequests.Post(fmt.Sprintf("%s/v1/reports", url), &grequests.RequestOptions{
		JSON: map[string]any{
			"data":       toPush,
			"commit_sha": commit,
		},
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "token " + token,
		},
	})

	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API over %s responded with HTTP %d: %s", url, resp.StatusCode, resp.String())
	}

	log.Debug("Submit successful.")
	return nil
}
