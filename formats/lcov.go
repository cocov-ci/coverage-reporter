package formats

import (
	"bytes"
	"fmt"
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/report"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Lcov struct {
	meta *meta.Metadata
}

func (l *Lcov) Wants(files []string) *string {
	for _, f := range files {
		if ext := filepath.Ext(f); ext != ".lcov" && ext != ".info" {
			continue
		}

		data, err := readFilePartial(l.meta.PathOf(f), 0, 3)
		if err != nil {
			// TODO: Log?
			continue
		}

		if !bytes.Equal(data, []byte("TN:")) {
			continue
		}

		return &f
	}

	return nil
}

func (l *Lcov) Name() string { return "Lcov/Gcov" }

func (l *Lcov) Parse(path string) (map[string]string, error) {
	f, err := os.Open(l.meta.PathOf(path))
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	result := map[string]string{}
	var currentFile *report.File
	flushCurrentFile := func() {
		if currentFile != nil {
			currentFile.Coverage.FillMissingLines(report.CoverageKindNeutral)
			result[currentFile.Path] = currentFile.Coverage.Encode()
			currentFile = nil
		}
	}

	lr := bufferedLineReader(f)
	for {
		line, err := lr.NextLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		comps := strings.SplitN(line, ":", 2)
		switch comps[0] {
		case "TN":
			// Test Name: noop
		case "SF":
			// Source File
			flushCurrentFile()
			path := comps[1]
			if filepath.IsAbs(path) {
				if strings.HasPrefix(path, l.meta.Pwd+"/") {
					path = strings.TrimPrefix(path, l.meta.Pwd+"/")
				} else if prfx := findFilePathPrefix(path, l.meta.Pwd); prfx != nil {
					path = strings.TrimPrefix(path, *prfx)
				}
			}

			currentFile = &report.File{
				Path:     path,
				Coverage: &report.CoverageSet{},
			}

			s, err := os.Stat(l.meta.PathOf(path))
			if err != nil {
				return nil, fmt.Errorf("invalid coverage report: could not stat `%s': %w", path, err)
			}
			if s.IsDir() {
				return nil, fmt.Errorf("invalid coverage report: `%s': is a directory", path)
			}
		case "FN":
			// Function Name: noop
		case "FNDA":
			// FuNctions coverage information
		case "FNF":
			// FuNctions Found: noop
		case "FNH":
			// FuNctions Hit: noop
		case "BRDA":
			// BRanch coverage information
		case "BRF":
			// BRanches Found
		case "BRH":
			// BRanches Hit
		case "DA":
			// coverage information
			dataComps := strings.Split(comps[1], ",")
			rawLineNumber, rawHits := dataComps[0], dataComps[1]
			lineNumber, err := strconv.ParseInt(rawLineNumber, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("corrupt coverage report: failed parsing line number on line `%s': %w", line, err)
			}
			hits, err := strconv.ParseInt(rawHits, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("corrupt coverage report: failed parsing hits on line `%s': %w", line, err)
			}

			currentFile.Coverage.SetLine(int(lineNumber), report.CoverageKindCovered, int(hits))
		case "LH":
			// number of lines with a non-zero execution count: noop
		case "LF":
			// number of instrumented lines: noop
		case "end_of_record":
			if currentFile == nil {
				return nil, fmt.Errorf("corrupt coverage report: found unexpected end_of_record")
			}
			flushCurrentFile()
		default:
			return nil, fmt.Errorf("corrupt coverage report: Unexpected line `%s'", line)
		}
	}

	flushCurrentFile()
	return result, nil
}

func (l *Lcov) SetMeta(meta *meta.Metadata) { l.meta = meta }
