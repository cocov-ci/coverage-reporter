package formats

import (
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/test_helpers"
	"github.com/cocov-ci/coverage-reporter/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func prepareCobertura(t *testing.T) *Cobertura {
	root := test_helpers.GitRoot(t)
	path := filepath.Join(root, "formats", "samples", "cobertura")
	files, err := tracking.FilesOn(path)
	require.NoError(t, err)

	return &Cobertura{
		meta: &meta.Metadata{
			Files: files,
			Pwd:   path,
		},
		sources: nil,
	}
}

func TestCobertura_Wants(t *testing.T) {
	cov := prepareCobertura(t)
	assert.Nil(t, cov.Wants(map[string]string{
		"cobeturas.xml": "",
	}))

	assert.Equal(t, "coverage/cobertura.xml", *cov.Wants(map[string]string{
		"coverage/cobertura.xml": "",
	}))
}

func TestCobertura_Parse(t *testing.T) {
	cov := prepareCobertura(t)
	p, err := cov.Parse("coverage/cobertura.xml")
	require.NoError(t, err)
	assert.Len(t, p, 1)
	assert.Equal(t, "MgAAAAAAMg==", p["src/config.js"])
}
