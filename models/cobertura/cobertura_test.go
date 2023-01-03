package cobertura

import (
	"encoding/xml"
	"github.com/cocov-ci/coverage-reporter/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestCoberturaUnmarshal(t *testing.T) {
	root := test_helpers.GitRoot(t)
	path := filepath.Join(root, "formats", "samples", "cobertura", "coverage", "cobertura.xml")
	file, err := os.Open(path)
	require.NoError(t, err)
	defer func(file *os.File) { _ = file.Close() }(file)

	cov := Coverage{}
	err = xml.NewDecoder(file).Decode(&cov)
	require.NoError(t, err)

	assert.Equal(t, []string{"/app"}, cov.Sources)
	assert.Equal(t, "src", cov.Packages[0].Name)
	assert.Equal(t, "src/config.js", cov.Packages[0].Classes[0].Filename)

	assert.Equal(t, 2, cov.Packages[0].Classes[0].Lines[0].Hits)
	assert.Equal(t, 1, cov.Packages[0].Classes[0].Lines[0].Number)

	assert.Equal(t, 2, cov.Packages[0].Classes[0].Lines[1].Hits)
	assert.Equal(t, 7, cov.Packages[0].Classes[0].Lines[1].Number)
}
