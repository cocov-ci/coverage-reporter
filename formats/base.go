package formats

import (
	"encoding/json"
	"fmt"
	"github.com/cocov-ci/coverage-reporter/meta"
	"go.uber.org/zap"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
)

type BaseFormat interface {
	Wants(files []string) *string
	Name() string
	Parse(path string) (map[string]string, error)
	SetMeta(meta *meta.Metadata)
}

var covers = []reflect.Type{
	reflect.TypeOf(&GoCov{}),
	reflect.TypeOf(&SimpleCov{}),
	reflect.TypeOf(&Lcov{}),
}

func makeCover(kind reflect.Type) BaseFormat {
	return reflect.New(kind.Elem()).Interface().(BaseFormat)
}

func TryParse(path string, meta *meta.Metadata) (map[string]string, error) {
	log := zap.L()
	log.Debug("TryParse called", zap.String("path", path))

	for _, kind := range covers {
		handler := makeCover(kind)
		log := log.With(zap.String("handler", handler.Name()))
		log.Debug("Trying handler", zap.String("name", handler.Name()))
		handler.SetMeta(meta)
		which := handler.Wants([]string{path})
		if which != nil {
			log.Debug("Handler indicates it may be able to handle report", zap.Stringp("report_path", which))
			if data, err := handler.Parse(*which); err == nil {
				log.Debug("Successful parse from handler. Returning it.")
				return data, nil
			} else {
				log.Debug("Handler failed to process report. Trying next.", zap.Stringp("report_path", which), zap.Error(err))
				continue
			}
		}
	}

	return nil, fmt.Errorf("could not find a handler for %s", path)
}

func filesFrom(diffs map[string]string) []string {
	out := make([]string, 0, len(diffs))
	for k := range diffs {
		out = append(out, k)
	}
	return out
}

func AutoFindAll(diffs map[string]string, meta *meta.Metadata) (map[string]string, error) {
	foundOne := false
	result := map[string]string{}

	for f := range diffs {
		parsed, err := TryParse(f, meta)
		if err != nil {
			continue
		}
		foundOne = true
		for k, v := range parsed {
			result[k] = v
		}
	}

	if !foundOne {
		return nil, fmt.Errorf("could not auto-detect coverage data")
	}

	return result, nil
}

func AutoFindOne(diffs map[string]string, meta *meta.Metadata) (map[string]string, error) {
	log := zap.L()
	log.Debug("AutoFind called")
	files := filesFrom(diffs)

	for _, kind := range covers {
		handler := makeCover(kind)
		log := log.With(zap.String("handler", handler.Name()))
		log.Debug("Trying handler")
		handler.SetMeta(meta)
		which := handler.Wants(files)
		if which != nil {
			log.Debug("Handler indicates it may be able to handle report", zap.Stringp("report_path", which))
			if data, err := handler.Parse(*which); err == nil {
				log.Debug("Successful parse from handler. Returning it.")
				return data, nil
			} else {
				log.Debug("Handler failed to process report. Trying next.", zap.Stringp("report_path", which), zap.Error(err))
				continue
			}
		}
	}

	return nil, fmt.Errorf("could not auto-detect coverage data")
}

func LoadPartials(runMeta *meta.Metadata) (map[string]string, error) {
	partialsDir := filepath.Join(meta.MetadataDir(runMeta.Token), "partials")
	result := map[string]string{}

	err := filepath.Walk(partialsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		var partial map[string]string
		if err = json.Unmarshal(data, &partial); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}

		for k, v := range partial {
			result[k] = v
		}

		return nil
	})

	return result, err
}
