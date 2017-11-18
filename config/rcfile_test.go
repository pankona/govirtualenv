package config

import (
	"os"
	"path/filepath"
	"testing"
)

const TestDataDirectory = "testdata"

func TestReadRCFileSuccess(t *testing.T) {
	path := filepath.Join(TestDataDirectory, "rcfile1")

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	if rc, err := ReadRCFile(file); err == nil {
		if rc.GitPath != "/usr/bin/git" {
			t.Errorf("expected %s, actual %s", "/usr/bin/git", rc.GitPath)
		}

		if rc.ManagementDirectoryPath != "/home/govenv/.govenv" {
			t.Errorf("expected %s, actual %s", "/home/govenv/.govenv", rc.ManagementDirectoryPath)
		}
	} else {
		t.Error(err)
	}
}
