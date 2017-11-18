package manager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/necomeshi/govirtualenv/config"
	"github.com/necomeshi/govirtualenv/git"
)

// goRootsManager is interface for management goroots directory
type goRootsManager struct {
	path       string
	gitCommand *git.Git
}

func newGoRootsManager(
	gitCommandPath, goRootsDirectoryPath string) (g *goRootsManager) {

	g = &goRootsManager{}
	g.path = goRootsDirectoryPath
	g.gitCommand = git.New(gitCommandPath, filepath.Join(goRootsDirectoryPath, config.GovenvGoMasterDir))

	return g
}

// LookupVersion returns installed GOROOT path of golang version
func (g *goRootsManager) LookupInstalledVersion(version string) (path string, found bool, err error) {
	versions, err := g.Installed()
	if err != nil {
		return "", false, err
	}

	for _, v := range versions {
		if v == version {
			return filepath.Join(g.path, version), true, nil
		}
	}
	return "", false, nil
}

// Installed returns installed golang version directory's path
func (g *goRootsManager) Installed() (versions []string, err error) {
	file, err := os.OpenFile(g.path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	dirs, err := file.Readdir(0)
	if err != nil {
		return nil, err
	}

	for _, dinfo := range dirs {
		if dinfo.IsDir() {
			if strings.Compare(config.GovenvGoBootstrapDir, dinfo.Name()) != 0 &&
				strings.Compare(config.GovenvGoMasterDir, dinfo.Name()) != 0 {

				versions = append(versions, dinfo.Name())
			}
		}
	}
	return versions, nil
}

// Available returns installable versions list (which not installed yet)
func (g *goRootsManager) Available() (versions []string, err error) {
	tags, err := g.gitCommand.Tags()
	if err != nil {
		return versions, err
	}

	for _, tagName := range tags {
		items := strings.Split(tagName, "/")
		versions = append(versions, items[len(items)-1])
	}

	return tags, err
}

// Install download/build new golang version under the path
func (g *goRootsManager) Install(version string) (err error) {
	found := false
	installable, err := g.Available()
	if err != nil {
		return err
	}

	for _, v := range installable {
		if v == version {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("No released version %s found", version)
	}

	if err := g.gitCommand.Checkout(version); err != nil {
		return err
	}

	if err := g.gitCommand.CheckoutIndex(fmt.Sprintf("../%s/", version)); err != nil {
		return err
	}

	bootstrap := filepath.Join(g.path, config.GovenvGoBootstrapDir)
	src := filepath.Join(g.path, version, "src")

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Chdir(src); err != nil {
		return err
	}
	defer os.Chdir(wd)

	buildCommand := exec.Command("/bin/bash", "all.bash")
	buildCommand.Env = append(os.Environ(), fmt.Sprintf("GOROOT_BOOTSTRAP=%s", bootstrap))

	if err := buildCommand.Run(); err != nil {
		return err
	}

	return nil
}

// Remove removes installed golang directory
func (g *goRootsManager) Remove(version string) (err error) {

	path, found, err := g.LookupInstalledVersion(version)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("Version %s not installed", version)
	}

	return os.RemoveAll(path)
}
