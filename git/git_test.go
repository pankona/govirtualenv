package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/necomeshi/govirtualenv/test"
)

const (
	TestDataDir           = "testdata"
	TestRepositoryDirName = "testrepo"
)

func TestTags(t *testing.T) {
	g := New(test.GetGitCommandPath(), TestRepositoryDirName)

	tags, err := g.Tags()
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range test.GetTestTable() {
		found := false
		for _, tag := range tags {
			if strings.Compare(entry.Tag, tag) == 0 {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected found '%s', not found", entry.Tag)
			for _, tag := range tags {
				t.Errorf("actual '%s'", tag)
			}
		}
	}
}

func TestFetch(t *testing.T) {
	g := New(test.GetGitCommandPath(), TestRepositoryDirName)

	if err := g.Fetch(); err != nil {
		t.Fatal(err)
	}
}

func TestCheckout(t *testing.T) {
	g := New(test.GetGitCommandPath(), TestRepositoryDirName)

	for _, entry := range test.GetTestTable() {
		tagName := strings.Join([]string{"refs", "tags", entry.Tag}, "/")

		if err := g.Checkout(tagName); err != nil {
			t.Errorf("Error checkouting '%s': %s", tagName, err)
			continue
		}

		err := test.CompareDirectory(TestRepositoryDirName, entry)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestCheckoutIndex(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	g := New(test.GetGitCommandPath(), TestRepositoryDirName)

	for _, entry := range test.GetTestTable() {
		tagName := strings.Join([]string{"refs", "tags", entry.Tag}, "/")

		if err := g.Checkout(tagName); err != nil {
			t.Errorf("Error checkouting '%s': %s", tagName, err)
			continue
		}

		dest := filepath.Join(wd, TestDataDir, tagName)

		if err := g.CheckoutIndex(fmt.Sprintf("../%s/%s/", TestDataDir, tagName)); err != nil {
			t.Errorf("Error %s", err)
			continue
		}

		err := test.CompareDirectory(dest, entry)
		if err != nil {
			t.Error(err)
		}
	}
}

func setup() (err error) {
	// Clone test reomote repository
	cmd := exec.Command("git", "clone", test.RemoteRepositoryURL, TestRepositoryDirName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Error %s; %s", err, string(output))
	}

	return nil
}

func teardown() (err error) {
	if err := os.RemoveAll(TestRepositoryDirName); err != nil {
		return err
	}

	if err := os.RemoveAll(TestDataDir); err != nil {
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
