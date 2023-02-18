package command

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cocov-ci/coverage-reporter/formats"
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func Add(ctx *cli.Context) error {
	token := getToken(ctx)
	runMeta, err := meta.ReadMetadata(token)
	if err != nil {
		return err
	}

	log := zap.L()
	if ctx.NArg() == 0 {
		log.Error("The 'add' command requires at least one file")
		return fmt.Errorf("invalid argument length for command 'add'")
	}

	partialsDir := filepath.Join(meta.MetadataDir(token), "partials")
	if _, err := os.Stat(partialsDir); os.IsNotExist(err) {
		if err = os.Mkdir(partialsDir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	for _, f := range ctx.Args().Slice() {
		f = runMeta.PathOf(f)

		log.Info("Attempting to parse", zap.String("path", f))
		parsed, err := formats.TryParse(f, runMeta)
		if err != nil {
			log.Error("Could not parse output", zap.String("path", f), zap.Error(err))
			return err
		}

		sha := sha1.New()
		sha.Write([]byte(f))
		outputName := hex.EncodeToString(sha.Sum(nil))
		encoded, err := json.Marshal(parsed)
		if err != nil {
			log.Error("Failed marshalling result", zap.Error(err))
			return err
		}

		if err = os.WriteFile(filepath.Join(partialsDir, outputName), encoded, 0655); err != nil {
			log.Error("Failed writing partial", zap.Error(err))
			return err
		}
	}

	// If we got here, we managed to store data. Set the manual flag on context
	if !runMeta.Manual {
		runMeta.Manual = true
		return meta.StoreMetadata(token, runMeta)
	}

	return nil
}
