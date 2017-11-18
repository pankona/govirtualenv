package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// RCFile is govenv environment file
type RCFile struct {
	ManagementDirectoryPath string
	GitPath                 string
}

// RCFilePath returns rc file path
func RCFilePath() string {
	if homeDir, exists := os.LookupEnv("HOME"); exists {
		return filepath.Join(homeDir, GovenvEnvironmentFile)
	}
	return filepath.Join(".", GovenvEnvironmentFile)
}

// ReadRCFile reads environment file
func ReadRCFile(file *os.File) (rcfile *RCFile, err error) {
	decoder := json.NewDecoder(file)

	rcfile = &RCFile{}

	if err := decoder.Decode(rcfile); err != nil {
		return nil, err
	}

	return rcfile, nil
}
