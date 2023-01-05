package meta

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Metadata struct {
	Files map[string]string `json:"files,omitempty"`
	Pwd   string            `json:"pwd,omitempty"`
	Sha   string            `json:"sha,omitempty"`
}

func (m Metadata) PathOf(relative string) string {
	if filepath.IsAbs(relative) {
		return relative
	}

	return filepath.Join(m.Pwd, relative)
}

func MetadataDir(token string) string {
	return filepath.Join(os.TempDir(), "cocov-"+token)
}

func MetadataFilePath(token string) string {
	return filepath.Join(MetadataDir(token), "meta.json")
}

func ReadMetadata(token string) (*Metadata, error) {
	data, err := os.ReadFile(MetadataFilePath(token))
	if err != nil {
		return nil, err
	}

	var meta Metadata
	if err = json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	return &meta, err
}
