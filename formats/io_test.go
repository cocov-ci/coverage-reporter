package formats

import (
	"bytes"
	"github.com/cocov-ci/coverage-reporter/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadPartialZeroOffset(t *testing.T) {
	root := test_helpers.GitRoot(t)
	file := filepath.Join(root, "go.mod")
	expected := "module"
	d, err := readFilePartial(file, 0, len(expected))
	require.NoError(t, err)
	assert.Equal(t, expected, string(d))
}

func TestReadPartialNonZeroOffset(t *testing.T) {
	root := test_helpers.GitRoot(t)
	file := filepath.Join(root, "go.mod")
	d, err := readFilePartial(file, 7, 37)
	require.NoError(t, err)
	assert.Equal(t, "github.com/cocov-ci/coverage-reporter", string(d))
}

func TestLineReader_NextLine(t *testing.T) {
	data := bytes.NewReader([]byte("a\nb\nc\nd\ne\nf"))
	expected := strings.Split("a\nb\nc\nd\ne\nf", "\n")
	reader := bufferedLineReader(data)
	for _, v := range expected {
		l, err := reader.NextLine()
		assert.NoError(t, err)
		assert.Equal(t, v, l)
	}
	l, err := reader.NextLine()
	assert.Zero(t, l)
	assert.ErrorIs(t, err, io.EOF)
}
