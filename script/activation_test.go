package script

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func compareFile(expectedFile, actualFile *os.File) (err error) {
	expectedReader := bufio.NewReader(expectedFile)
	actualReader := bufio.NewReader(actualFile)

	for {
		expected, err := expectedReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		actual, err := actualReader.ReadString('\n')
		if err != nil {
			return err
		}

		if strings.Compare(expected, actual) != 0 {
			return fmt.Errorf("expected %s, actual %s", expected, actual)
		}
	}

	return nil
}

func TestCreateActivationScript(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	env := New("project", wd, wd)
	if err := env.CreateScript(wd); err != nil {
		t.Fatal(err)
	}

	expectedFile, err := os.OpenFile("testdata", os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer expectedFile.Close()

	actualFile, err := os.OpenFile(filepath.Join(wd, "activate"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer actualFile.Close()

	if err := compareFile(expectedFile, actualFile); err != nil {
		t.Fatal(err)
	}
}

func teardown() (err error) {
	if err := os.RemoveAll("activate"); err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {

	ret := m.Run()

	teardown()

	os.Exit(ret)
}
