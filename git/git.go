package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

// Git is git command interface
type Git struct {
	Path       string
	Repository string
	gitDir     string
}

// New creates new git command interface
func New(commandPath, repository string) (g *Git) {
	g = &Git{Path: commandPath, Repository: repository}

	g.gitDir = filepath.Join(repository, ".git")

	return g
}

func (g *Git) execute(options ...string) (output []byte, err error) {

	cmd := exec.Command(g.Path, "--git-dir", g.gitDir, "--work-tree", g.Repository)

	for _, opt := range options {
		cmd.Args = append(cmd.Args, opt)
	}

	if output, err = cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("Error %s; %s", string(output), err)
	}

	return output, nil
}

// Fetch fetch repository
func (g *Git) Fetch() (err error) {
	_, err = g.execute("fetch")
	return err
}

// Tags returns available tags
func (g *Git) Tags() (tags []string, err error) {

	var output []byte
	if output, err = g.execute("tag", "-l"); err != nil {
		return tags, err
	}

	reader := bufio.NewReader(bytes.NewReader(output))

	for {
		s, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		tags = append(tags, strings.Trim(s, "\n"))
	}

	return tags, nil
}

// Checkout checkouts source by tag name
func (g *Git) Checkout(branch string) (err error) {
	if _, err := g.execute("checkout", branch); err != nil {
		return err
	}
	return nil
}

// CheckoutIndex export indexed files to destination
func (g *Git) CheckoutIndex(dest string) (err error) {
	if !strings.HasSuffix(dest, "/") {
		dest += "/"
	}

	if _, err := g.execute("checkout-index", "-a", "--prefix", dest); err != nil {
		return err
	}

	return nil
}
