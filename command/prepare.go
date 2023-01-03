package command

import (
	"encoding/json"
	"fmt"
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/tracking"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
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

	metaEntry := &meta.Metadata{
		Files: files,
		Pwd:   pwd,
	}
	data, err := json.Marshal(metaEntry)

	if os.Getenv("COCOV_REPORTER_DEBUG_METADATA") == "true" {
		log.Debug("Obtained metadata directory", zap.String("path", meta.MetadataDir(token)))
	}

	if err != nil {
		return err
	}

	if err = os.RemoveAll(meta.MetadataDir(token)); err != nil {
		return err
	}

	if err = os.MkdirAll(meta.MetadataDir(token), 0750); err != nil {
		return err
	}

	target := meta.MetadataFilePath(token)
	return os.WriteFile(target, data, 0655)
}
