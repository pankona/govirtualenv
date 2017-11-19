package project

import (
	"os"
	"path/filepath"

	"github.com/necomeshi/govirtualenv/config"
	"github.com/necomeshi/govirtualenv/script"
)

// Project is management interface for project directorty
type Project struct {
	env *script.GoEnv
}

// IsActivated checks is activated virtual environment or not.
func IsActivated() bool {
	return script.IsActivated()
}

// Create creates new project directory under the goPath with
// specified go version and return new project struct
func Create(name, goRootPath, goPath string) (p *Project, err error) {
	projectPath := filepath.Join(goPath, config.GovenvGoProjectDir)
	env := script.New(name, goRootPath, goPath)
	dest := filepath.Join(projectPath, config.GovenvGoProjectScriptDir)

	// Create project data direcotry
	if err := os.MkdirAll(projectPath, os.ModePerm); err != nil {
		return nil, err
	}

	// Create scripts direcotry
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return nil, err
	}

	// Fill activate script blanks and create activation script under the ENVDIR/bin
	if err := env.CreateScript(dest); err != nil {
		return nil, err
	}

	return &Project{env: env}, nil
}

// CreateFromCurrentEnvironment creates Project from current environment variable
func CreateFromCurrentEnvironment() (p *Project, err error) {
	p = &Project{}

	p.env, err = script.NewFromEnvironmentVariables()

	return p, err
}

// SetGoRoot set goRootPath to GOROOT
func (pj *Project) SetGoRoot(goRootPath string) {
	pj.env.GoRoot = goRootPath
}

// Configure configurate golang project
func (pj *Project) Configure() (err error) {
	projectPath := filepath.Join(pj.env.GoPath, config.GovenvGoProjectDir)
	dest := filepath.Join(projectPath, config.GovenvGoProjectScriptDir)

	return pj.env.CreateScript(dest)
}
