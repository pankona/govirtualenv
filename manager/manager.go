package manager

import (
	"path/filepath"

	"github.com/necomeshi/govirtualenv/config"
)

// Manager is govenv management interface
type Manager struct {
	managementDirectoryPath string
	gitCommandPath          string
}

// New creates new management interface structure
func New(gitCommandPath, managementDirectoryPath string) (m *Manager) {
	m = &Manager{
		gitCommandPath:          gitCommandPath,
		managementDirectoryPath: managementDirectoryPath,
	}

	return m
}

// GetGoRootsManager creates new goroots directory management interface structure
func (m *Manager) GetGoRootsManager() (g *goRootsManager) {
	goRootsPath := filepath.Join(m.managementDirectoryPath, config.GovenvGoRootsDirName)

	return newGoRootsManager(m.gitCommandPath, goRootsPath)
}
