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

func prepareJacoco(t *testing.T) *Jacoco {
	root := test_helpers.GitRoot(t)
	path := filepath.Join(root, "formats", "samples", "jacoco", "small")
	files, err := tracking.FilesOn(path)
	require.NoError(t, err)

	return &Jacoco{
		meta: &meta.Metadata{
			Files: files,
			Pwd:   path,
		},
	}
}

func TestJacoco_Wants(t *testing.T) {
	cov := prepareJacoco(t)
	assert.Nil(t, cov.Wants([]string{
		"random-file.xml",
	}))

	assert.Equal(t, "build/reports/jacoco/aggregate/jacocoTestReport.xml", *cov.Wants([]string{
		"build/reports/jacoco/aggregate/jacocoTestReport.xml",
	}))
}

func TestJacoco_Parse(t *testing.T) {
	cov := prepareJacoco(t)
	p, err := cov.Parse("build/reports/jacoco/aggregate/jacocoTestReport.xml")
	require.NoError(t, err)
	assert.Len(t, p, 2)
	assert.Equal(t, "AAAAAAAA", p["app/src/main/kotlin/com/company/project/batch/parsing/FileWithNoExecutableCode.kt"])
	assert.Equal(t, "AAAAAAAAOQA1HjQeMx4xAAAzHjEA", p["app/src/main/kotlin/com/company/project/batch/parsing/ParsedBatchWriter.kt"])
}
