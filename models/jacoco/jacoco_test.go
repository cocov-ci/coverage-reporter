package jacoco

import (
	"encoding/xml"
	"github.com/cocov-ci/coverage-reporter/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestJacocoUnmarshal(t *testing.T) {
	root := test_helpers.GitRoot(t)
	path := filepath.Join(root, "formats/samples/jacoco/small/build/reports/jacoco/aggregate/jacocoTestReport.xml")
	file, err := os.Open(path)
	require.NoError(t, err)
	defer func(file *os.File) { _ = file.Close() }(file)

	cov := Coverage{}
	err = xml.NewDecoder(file).Decode(&cov)
	require.NoError(t, err)

	assert.Len(t, cov.Packages, 1)
	assert.Equal(t, "com/company/project/batch/parsing", cov.Packages[0].Name)

	assert.Len(t, cov.Packages[0].SourceFiles, 2)
	assert.Equal(t, "FileWithNoExecutableCode.kt", cov.Packages[0].SourceFiles[0].Name)
	assert.Len(t, cov.Packages[0].SourceFiles[0].Lines, 0)

	assert.Equal(t, "ParsedBatchWriter.kt", cov.Packages[0].SourceFiles[1].Name)
	assert.Len(t, cov.Packages[0].SourceFiles[1].Lines, 7)
	assert.Equal(t, 7, cov.Packages[0].SourceFiles[1].Lines[0].Number)

	pairs := [][]int{
		{7, 9},
		{9, 5},
		{10, 4},
		{11, 3},
		{12, 1},
		{15, 3},
		{16, 1},
	}

	for i, v := range pairs {
		assert.Equal(t, v[0], cov.Packages[0].SourceFiles[1].Lines[i].Number)
		assert.Equal(t, v[1], cov.Packages[0].SourceFiles[1].Lines[i].Hits)
	}
}
