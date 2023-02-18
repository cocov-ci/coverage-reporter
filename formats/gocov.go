package formats

import (
	"bytes"
	"fmt"
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/report"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/cover"
	"os"
	"path/filepath"
	"strings"
)

type GoCov struct {
	Meta *meta.Metadata
	mods map[string]string
}

func (g *GoCov) SetMeta(meta *meta.Metadata) { g.Meta = meta }

func (g *GoCov) Name() string { return "Golang default coverage format" }

func (g *GoCov) Wants(files []string) *string {
	for _, v := range files {
		if base := filepath.Base(v); base != "c.out" && base != "coverage.out" {
			continue
		}
		b, err := readFilePartial(g.Meta.PathOf(v), 0, 5)
		if err != nil {
			continue
		}
		if bytes.Equal(b, []byte("mode:")) {
			return &v
		}
	}
	return nil
}

func (g *GoCov) Parse(path string) (map[string]string, error) {
	profiles, err := cover.ParseProfiles(g.Meta.PathOf(path))
	if err != nil {
		return nil, err
	}

	if len(profiles) == 0 {
		return map[string]string{}, nil
	}

	allMods, err := g.FindMods()
	if err != nil {
		return nil, err
	}

	if len(allMods) == 0 {
		return nil, fmt.Errorf("modless projects are not supported")
	} else {
		return g.ModBasedCoverage(profiles, allMods)
	}
}

func (g *GoCov) ModBasedCoverage(profiles []*cover.Profile, mods map[string]string) (map[string]string, error) {
	output := map[string]string{}
	for _, v := range profiles {
		mod := g.ModFor(v.FileName, mods)
		if mod == nil {
			return nil, fmt.Errorf("could not locate source for %s", v.FileName)
		}
		baseDir := mods[*mod]
		f := &report.File{
			Path:     strings.TrimPrefix(strings.Replace(v.FileName, *mod, baseDir, 1), g.Meta.Pwd+"/"),
			Coverage: &report.CoverageSet{},
		}

		for _, b := range v.Blocks {
			for i := b.StartLine; i <= b.EndLine; i++ {
				f.Coverage.SetLine(i, report.CoverageKindCovered, b.Count)
			}
		}

		f.Coverage.FillMissingLines(report.CoverageKindNeutral)
		output[f.Path] = f.Coverage.Encode()
	}

	return output, nil
}

func (g *GoCov) ModFor(file string, mods map[string]string) *string {
	components := strings.Split(file, "/")

	for len(components) > 0 {
		components = components[:len(components)-1]
		possibleMod := strings.Join(components, "/")
		_, ok := mods[possibleMod]
		if ok {
			return &possibleMod
		}
	}

	return nil
}

func (g *GoCov) FindMods() (map[string]string, error) {
	if g.mods != nil {
		return g.mods, nil
	}

	mods := map[string]string{}
	for f := range g.Meta.Files {
		if filepath.Base(f) != "go.mod" {
			continue
		}

		path := g.Meta.PathOf(f)
		data, err := os.ReadFile(path)
		if err != nil {
			// TODO: Log?
			continue
		}

		mod, err := modfile.Parse(path, data, nil)
		if err != nil {
			// TODO: Log?
			continue
		}

		mods[mod.Module.Mod.Path] = filepath.Dir(path)
	}

	g.mods = mods
	return mods, nil
}
