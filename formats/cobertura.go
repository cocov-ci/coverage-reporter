package formats

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/models/cobertura"
	"github.com/cocov-ci/coverage-reporter/report"
	"os"
	"path/filepath"
	"strings"
)

type Cobertura struct {
	meta    *meta.Metadata
	sources []string
}

func (c *Cobertura) Wants(diffs map[string]string) *string {
	coverageBytes := []byte("<coverage")
	lineRateBytes := []byte("line-rate=")

	for f := range diffs {
		if filepath.Base(f) != "cobertura.xml" {
			continue
		}

		fileData, err := os.ReadFile(c.meta.PathOf(f))
		if err != nil {
			// TODO: log?
			continue
		}

		if bytes.Contains(fileData, coverageBytes) &&
			bytes.Contains(fileData, lineRateBytes) {
			return &f
		}
	}
	return nil
}

func (c *Cobertura) Name() string { return "Cobertura" }

func (c *Cobertura) Parse(path string) (map[string]string, error) {
	xmlFile, err := os.Open(c.meta.PathOf(path))
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) { _ = f.Close() }(xmlFile)

	cov := cobertura.Coverage{}
	if err = xml.NewDecoder(xmlFile).Decode(&cov); err != nil {
		return nil, err
	}

	if !cov.Valid() {
		return nil, fmt.Errorf("could not successfully parse cobertura output")
	}
	c.sources = cov.Sources
	classesPerFile := map[string][]cobertura.Class{}

	for _, p := range cov.Packages {
		for _, pc := range p.Classes {
			f := c.realFilePath(pc.Filename)
			if f == "" {
				return nil, fmt.Errorf("could not obtain real path for file %s", pc.Filename)
			}
			cp := classesPerFile[f]
			cp = append(cp, pc)
			classesPerFile[f] = cp
		}
	}

	output := map[string]string{}

	for fPath, clSlice := range classesPerFile {
		f := report.File{
			Path:     fPath,
			Coverage: &report.CoverageSet{},
		}

		for _, cl := range clSlice {
			for _, l := range cl.Lines {
				f.Coverage.SetLine(l.Number, report.CoverageKindCovered, l.Hits)
			}
		}

		f.Coverage.FillMissingLines(report.CoverageKindNeutral)
		output[f.Path] = f.Coverage.Encode()
	}

	return output, nil
}

func (c *Cobertura) realFilePath(partial string) string {
	for _, s := range c.sources {
		path := filepath.Join(s, partial)

		if prfx := findFilePathPrefix(path, c.meta.Pwd); prfx != nil {
			path = strings.TrimPrefix(path, *prfx)
		} else if strings.HasPrefix(path, c.meta.Pwd+"/") {
			path = strings.TrimPrefix(path, c.meta.Pwd+"/")
		}

		if stat, err := os.Stat(filepath.Join(c.meta.Pwd, path)); err == nil && !stat.IsDir() {
			return path
		}
	}

	return ""
}

func (c *Cobertura) SetMeta(meta *meta.Metadata) {
	c.meta = meta
}
