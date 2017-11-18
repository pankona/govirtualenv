package project

import (
	"os"
	"path/filepath"

	"github.com/necomeshi/govirtualenv/config"
	"github.com/necomeshi/govirtualenv/script"
)

// Project is management interface for project directorty
type Project struct {
	path string
}

// Create creates new project directory under the goPath with
// specified go version and return new project struct
func Create(goRootPath, goPath string) (p *Project, err error) {
	projectEnvPath := filepath.Join(goPath, config.GovenvGoProjectDir)
	projectScriptDirPath := filepath.Join(projectEnvPath, config.GovenvGoProjectScriptDir)

	// Create all project data direcotry
	if err := os.MkdirAll(projectScriptDirPath, os.ModePerm); err != nil {
		return nil, err
	}

	// Fill activate script blanks and create activation script under the ENVDIR/bin
	if err := script.CreateScript(projectScriptDirPath, goRootPath, goPath); err != nil {
		return nil, err
	}

	return &Project{path: goPath}, nil
}
