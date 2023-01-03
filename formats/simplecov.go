package formats

import (
	"encoding/json"
	"fmt"
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/report"
	"os"
	"path/filepath"
	"strings"
)

type simplecovMeta struct {
	Meta struct {
		SimpleCovVersion string `json:"simplecov_version"`
	} `json:"meta"`
}

type simplecovEntry struct {
	Lines []any `json:"lines"`
}

type simplecovFile struct {
	Coverage map[string]simplecovEntry `json:"coverage"`
}

type SimpleCov struct {
	meta *meta.Metadata
}

func (s *SimpleCov) Wants(diffs map[string]string) *string {
	for f := range diffs {
		if filepath.Base(f) != "coverage.json" {
			continue
		}
		var sMeta simplecovMeta
		data, err := os.ReadFile(s.meta.PathOf(f))
		if err != nil {
			// TODO: Log?
			continue
		}
		if err = json.Unmarshal(data, &sMeta); err != nil {
			// TODO: Log?
			continue
		}
		if sMeta.Meta.SimpleCovVersion == "" {
			continue
		}

		return &f
	}

	return nil
}

func (s *SimpleCov) Name() string { return "SimpleCov" }

func (s *SimpleCov) Parse(path string) (map[string]string, error) {
	data, err := os.ReadFile(s.meta.PathOf(path))
	if err != nil {
		return nil, err
	}

	var coverage simplecovFile
	if err = json.Unmarshal(data, &coverage); err != nil {
		return nil, err
	}

	output := map[string]string{}

	if len(coverage.Coverage) == 0 {
		return output, nil
	}

	var examplePath string
	for f := range coverage.Coverage {
		examplePath = f
		break
	}
	prefix := s.extractPrefix(examplePath)
	if prefix == nil {
		prefix = &s.meta.Pwd
	}

	for f, d := range coverage.Coverage {
		filePath := strings.TrimPrefix(f, *prefix)
		rf := &report.File{
			Path:     filePath,
			Coverage: &report.CoverageSet{},
		}

		for i, v := range d.Lines {
			line := i + 1
			switch l := v.(type) {
			case string:
				if l == "ignored" {
					rf.Coverage.SetLine(line, report.CoverageKindIgnore, 0)
				} else {
					return nil, fmt.Errorf("unexpected string value `%s' on coverage information for %s, line offset %d", l, filePath, i)
				}
			case nil:
				rf.Coverage.SetLine(line, report.CoverageKindNeutral, 0)
			case float64:
				rf.Coverage.SetLine(line, report.CoverageKindCovered, int(l))
			default:
				return nil, fmt.Errorf("unexpected type %T on coverage information for %s, line offset %d", v, filePath, i)
			}
		}

		output[filePath] = rf.Coverage.Encode()
	}

	return output, nil
}

func (s *SimpleCov) extractPrefix(path string) *string {
	return findFilePathPrefix(path, s.meta.Pwd)
}

func (s *SimpleCov) SetMeta(meta *meta.Metadata) { s.meta = meta }
