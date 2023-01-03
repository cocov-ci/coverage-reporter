package formats

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/models/jacoco"
	"github.com/cocov-ci/coverage-reporter/report"
	"os"
	"path/filepath"
	"strings"
)

type Jacoco struct {
	meta *meta.Metadata
}

func (j *Jacoco) Wants(diffs map[string]string) *string {
	jacocoHeader := []byte("-//JACOCO//DTD")
	reportTag := []byte("<report")
	sessionInfoTag := []byte("<sessioninfo")
	wants := [][]byte{jacocoHeader, reportTag, sessionInfoTag}

	buf := make([]byte, 1024)
diffLoop:
	for f := range diffs {
		if filepath.Ext(f) != ".xml" {
			continue
		}
		// Read the first 1024 bytes, and look for header, report tag and
		// session info. Discard otherwise.
		l, err := readFilePartialBuf(j.meta.PathOf(f), 0, buf)
		if err != nil {
			// TODO: Log?
			continue
		}
		for _, w := range wants {
			if !bytes.Contains(buf[:l], w) {
				continue diffLoop
			}
		}

		return &f
	}

	return nil
}

func (j *Jacoco) Name() string { return "JaCoCo" }

func (j *Jacoco) Parse(path string) (map[string]string, error) {
	xmlFile, err := os.Open(j.meta.PathOf(path))
	if err != nil {
		return nil, err
	}
	defer func(xmlFile *os.File) { _ = xmlFile.Close() }(xmlFile)

	var cov jacoco.Coverage
	err = xml.NewDecoder(xmlFile).Decode(&cov)
	if err != nil {
		return nil, err
	}

	if cov.SessionInfo.Start == "" {
		// Assume this is an invalid report
		return nil, fmt.Errorf("invalid or corrupt jacoco report missing sessioninfo")
	}

	output := map[string]string{}

	for _, pkg := range cov.Packages {
		for _, sf := range pkg.SourceFiles {
			path, lineCount := j.filePath(pkg.Name, sf.Name)
			if path == nil {
				return nil, fmt.Errorf("could not found source file for %s in package %s", sf.Name, pkg.Name)
			}

			f := report.File{
				Path:     *path,
				Coverage: &report.CoverageSet{},
			}

			f.Coverage.SetLine(lineCount, report.CoverageKindNeutral, 0)

			for _, l := range sf.Lines {
				f.Coverage.SetLine(l.Number, report.CoverageKindCovered, l.Hits)
			}

			f.Coverage.FillMissingLines(report.CoverageKindNeutral)

			output[f.Path] = f.Coverage.Encode()
		}
	}

	return output, nil
}

func (j *Jacoco) filePath(pkg, file string) (*string, int) {
	builtPath := append(strings.Split(pkg, "/"), file)
	expectedPath := filepath.Join(builtPath...)
	for f := range j.meta.Files {
		if !strings.HasSuffix(f, expectedPath) {
			continue
		}

		data, err := os.ReadFile(j.meta.PathOf(f))
		if err != nil {
			continue
		}
		return &f, bytes.Count(data, []byte{'\n'})
	}

	return nil, 0
}

func (j *Jacoco) SetMeta(meta *meta.Metadata) {
	j.meta = meta
}
