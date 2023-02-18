package command

import (
	"fmt"
	"github.com/cocov-ci/coverage-reporter/formats"
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/tracking"
	"github.com/levigross/grequests"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"strings"
)

func Submit(ctx *cli.Context) error {
	log := zap.L()

	token := getToken(ctx)
	runMeta, err := meta.ReadMetadata(token)
	if err != nil {
		return err
	}

	currentFiles, err := tracking.FilesOn(runMeta.Pwd)
	if err != nil {
		return err
	}

	var toPush map[string]string
	if runMeta.Manual {
		toPush, err = formats.LoadPartials(runMeta)
		if err != nil {
			log.Debug("Finished aggregating partials", zap.Int("files_covered", len(toPush)))
		}
	} else {
		diffs := tracking.DiffFiles(runMeta.Files, currentFiles)
		if ctx.Bool("multi") {
			toPush, err = formats.AutoFindAll(diffs, runMeta)

			if err != nil {
				log.Debug("Finished scanning all generated files", zap.Int("files_covered", len(toPush)))
			}
		} else {
			toPush, err = formats.AutoFindOne(diffs, runMeta)
			if err != nil {
				log.Debug("Finished scanning a single generated file", zap.Int("files_covered", len(toPush)))
			}
		}
	}

	if err != nil {
		log.Error("Failed preparing data for submission", zap.Error(err))
		return err
	}

	log.Debug("Now submitting data")

	url := strings.TrimPrefix(ctx.String("url"), "/")
	resp, err := grequests.Post(fmt.Sprintf("%s/v1/reports", url), &grequests.RequestOptions{
		JSON: map[string]any{
			"data":       toPush,
			"commit_sha": runMeta.Sha,
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
