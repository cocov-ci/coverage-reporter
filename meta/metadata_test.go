package meta

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestMetadataDir(t *testing.T) {
	tmpdir := os.TempDir()
	metadata := MetadataDir("foobar")
	assert.Equal(t, filepath.Join(tmpdir, "cocov-foobar"), metadata)
}

func TestMetadataFilePath(t *testing.T) {
	tmpdir := os.TempDir()
	metadata := MetadataFilePath("foobar")
	assert.Equal(t, filepath.Join(tmpdir, "cocov-foobar", "meta.json"), metadata)
}

func TestReadMetadata(t *testing.T) {
	data := `{"files":{"a":"b"},"pwd":"c"}`
	targetDir := filepath.Join(os.TempDir(), "cocov-foobar")
	targetFile := filepath.Join(targetDir, "meta.json")

	require.NoError(t, os.MkdirAll(targetDir, 0755))
	t.Cleanup(func() {
		_ = os.RemoveAll(targetDir)
	})
	require.NoError(t, os.WriteFile(targetFile, []byte(data), 0655))

	meta, err := ReadMetadata("foobar")
	assert.NoError(t, err)

	assert.Equal(t, map[string]string{"a": "b"}, meta.Files)
	assert.Equal(t, "c", meta.Pwd)
}
