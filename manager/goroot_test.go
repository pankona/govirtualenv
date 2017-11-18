package manager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/necomeshi/govirtualenv/config"
	"github.com/necomeshi/govirtualenv/test"
)

const (
	TestDataDirName = "testdata" // This is GOROOTs directory for test
)

func TestLoockupInstalledVersion(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	m := newGoRootsManager("git", filepath.Join(wd, TestDataDirName))

	for _, entry := range test.GetTestTable() {
		path, found, err := m.LookupInstalledVersion(entry.Tag)
		if err != nil {
			t.Error(err)
			continue
		}

		if entry.Exported {
			if found {
				dest := filepath.Join(wd, TestDataDirName, entry.Tag)

				if strings.Compare(path, dest) != 0 {
					t.Errorf("expected %s, actual %s", dest, path)
				}
			} else {
				t.Errorf("Expected '%s' should be exist, but not found", entry.Tag)
			}
		} else {
			if found {
				t.Errorf("Expected '%s' is not exist", entry.Tag)
			}
		}
	}
}

func TestInstalled(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	m := newGoRootsManager("git", filepath.Join(wd, TestDataDirName))

	versions, err := m.Installed()
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range test.GetTestTable() {
		found := false
		for _, v := range versions {
			if strings.Compare(v, entry.Tag) == 0 {
				found = true
				break
			}
		}

		if entry.Exported {
			// entry.tag was exported, so it should be found
			if !found {
				t.Errorf("expected %s found", entry.Tag)
			}
		} else {
			// entry.tag was not exported yet, so it should NOT be found
			if found {
				t.Errorf("expected %s is not found", entry.Tag)
			}
		}
	}
}

func TestSuccessAvailable(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	m := newGoRootsManager("git", filepath.Join(wd, TestDataDirName))

	versions, err := m.Available()
	if err != nil {
		t.Fatal(err)
	}

	// May version found in tags
	for _, entry := range test.GetTestTable() {
		found := false
		for _, v := range versions {
			if entry.Tag == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Version '%s' expected found in Available()", entry.Tag)
		}
	}

	for _, v := range versions {
		found := false
		for _, entry := range test.GetTestTable() {
			if entry.Tag == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Version '%s' expected found in TestTable", v)
		}
	}
}

func TestRemove(t *testing.T) {
	// Note This must need.
	// Above test remove directories which other tests
	// expect that directories should be exist
	// So we need to recraete these removed directories
	defer func() {
		teardown()
		setup()
	}()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	m := newGoRootsManager("git", filepath.Join(wd, TestDataDirName))

	for _, entry := range test.GetTestTable() {
		err := m.Remove(entry.Tag)

		if entry.Exported {
			if err != nil {
				t.Error(err)
			}
		} else {
			if err == nil {
				t.Errorf("expected error happen but not occur %s", entry.Tag)
			}
		}
	}
}

func TestInstall(t *testing.T) {
	defer func() {
		teardown()
		setup()
	}()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	m := newGoRootsManager("git", filepath.Join(wd, TestDataDirName))

	if err := m.Install("v3.0"); err != nil {
		t.Error(err)
	}
}

func setup() (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	goRootsPath := filepath.Join(wd, TestDataDirName)
	golangDirPath := filepath.Join(goRootsPath, config.GovenvGoMasterDir)

	// Clone test reomote repository
	cmd := exec.Command("git", "clone", test.RemoteRepositoryURL, golangDirPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Error %s; %s", err, string(output))
	}

	gitDir := filepath.Join(golangDirPath, ".git")
	workTree := golangDirPath

	for _, entry := range test.GetTestTable() {
		if entry.Exported {
			tag := strings.Join([]string{"refs", "tags", entry.Tag}, "/")

			cmd = exec.Command("git", "--git-dir", gitDir, "--work-tree", workTree, "checkout", tag)
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("Error %s; %s", err, string(output))
			}

			cmd = exec.Command(
				"git", "--git-dir", gitDir, "--work-tree", workTree,
				"checkout-index", "-a", "--prefix", "../"+entry.Tag+"/")

			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("Error %s; %s", err, string(output))
			}
		}
	}
	return nil
}

func teardown() (err error) {
	if err := os.RemoveAll(TestDataDirName); err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		panic(err)
	}

	ret := m.Run()

	if err := teardown(); err != nil {
		panic(err)
	}

	os.Exit(ret)
}
