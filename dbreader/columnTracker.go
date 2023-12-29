package dbreader

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

type ColumnTracker struct {
	path string
	keys map[string]string `yaml:"keys"`
}

func NewColumnTracker(path string) (*ColumnTracker, error) {
	ct := &ColumnTracker{path: filepath.Join(path, "/.tracking-column.yml"), keys: make(map[string]string)}
	err := ct.load()
	if err != nil {
		return nil, err
	} // Load existing keys from the file
	return ct, nil
}

func (t *ColumnTracker) GetTrackingColumn(name string, hash ...string) string {
	hashedName := hashName(name, hash...)
	if value, ok := t.keys[hashedName]; ok {
		return value
	}
	return "" // or return an error if you prefer
}

func (t *ColumnTracker) SetTrackingColumn(name string, value string, hash ...string) error {
	hashedName := hashName(name, hash...)
	t.keys[hashedName] = value
	return t.save()
}

// load reads the tracking columns from the file.
func (t *ColumnTracker) load() error {
	if _, err := os.Stat(t.path); os.IsNotExist(err) {
		t.keys = make(map[string]string)

		return nil // File does not exist, nothing to load
	}

	data, err := os.ReadFile(t.path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &t.keys)
}

// save writes the tracking columns to the file.
func (t *ColumnTracker) save() error {
	data, err := yaml.Marshal(t.keys)
	if err != nil {
		return err
	}

	return os.WriteFile(t.path, data, 0644) // 0644 provides read and write permissions to the owner and read-only to others.
}

func hashName(name string, hash ...string) string {
	h := md5.Sum([]byte(strings.Join(hash, "")))

	// Convert the byte slice to a hexadecimal string
	hexString := fmt.Sprintf("%x", h)
	return name + "-" + hexString
}
