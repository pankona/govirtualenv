package test

import (
	"os"
	"path/filepath"
)

const (
	// PathToGitCommand is path to git command
	PathToGitCommand = "/usr/bin/git"
	// RemoteRepositoryURL is URL to repository for tests
	RemoteRepositoryURL = "https://github.com/necomeshi/govenvtest"
)

type TestTableEntry struct {
	Tag       string
	Exist     bool
	Exported  bool
	Filenames []string
}

func GetTestTable() []TestTableEntry {
	var TestTable = []TestTableEntry{
		{Tag: "v1.0", Exist: true, Exported: false, Filenames: []string{"README.md"}},
		{Tag: "v1.1", Exist: true, Exported: true, Filenames: []string{"README.md", "File1"}},
		{Tag: "v2.0", Exist: true, Exported: true, Filenames: []string{"README.md", "File1", "File2"}},
		{Tag: "v2.1", Exist: true, Exported: false, Filenames: []string{"README.md", "File1", "File2", "File3"}},
		{Tag: "v3.0", Exist: true, Exported: false, Filenames: []string{
			"README.md", "File1", "File2", "File3", filepath.Join("src", "all.bash")}},
		{Tag: "bootstrap", Exist: true, Exported: false, Filenames: []string{
			"README.md", "File1", "File2", "File3", filepath.Join("src", "all.bash")}},
	}
	return TestTable
}

func GetGitCommandPath() string {
	return PathToGitCommand
}

func CompareDirectory(pathToDir string, expected TestTableEntry) (err error) {

	for _, filename := range expected.Filenames {
		path := filepath.Join(pathToDir, filename)

		rootDir, err := os.Open(path)
		if err != nil {
			return err
		}
		rootDir.Close()
	}
	return nil
}
