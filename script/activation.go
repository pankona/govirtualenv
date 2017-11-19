package script

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const scriptTemplate = `#!/bin/bash
function deactivate() {
	if [ -n "${GOVENV_OLD_PATH}" ]; then
		PATH=${GOVENV_OLD_PATH}
		export PATH
		unset GOVENV_OLD_PATH
	fi

	if [ -n "${GOVENV_OLD_PS1}" ]; then
		PS1=${GOVENV_OLD_PS1}
		export PS1
		unset GOVENV_OLD_PS1
	fi

	if [ -n "${GOROOT}" ]; then
		unset GOROOT
	fi

	if [ -n "${GOPATH}" ]; then
		unset GOPATH
	fi

	unset GOVENV_PROJECT
	unset GOVENV_ENABLE
}

# Activation
if [ -z "${GOVENV_ENABLE+x}" ]; then
	GOVENV_OLD_PATH=${PATH}
	GOVENV_OLD_PS1=${PS1}

	GOROOT={{.GoRoot}}
	GOPATH={{.GoPath}}
	GOVENV_PROJECT={{.ProjectName}}

	PATH=${GOROOT}/bin:${PATH}
	PS1="(${GOVENV_PROJECT}) ${PS1}"

	export PS1 PATH GOROOT GOPATH GOVENV_ENABLE=1 GOVENV_PROJECT
else
	echo "Already in the project. First do 'deactivate'"
fi
`

// GoEnv is environment variables
type GoEnv struct {
	GoRoot      string
	GoPath      string
	ProjectName string
}

// New creates empty GoEnv structure
func New(projectName, goRoot, goPath string) *GoEnv {
	return &GoEnv{GoRoot: goRoot, GoPath: goPath, ProjectName: projectName}
}

// NewFromEnvironmentVariables collects virtual
// environment information from environment variables
func NewFromEnvironmentVariables() (envVal *GoEnv, err error) {
	var found bool
	envVal = &GoEnv{}

	if envVal.GoRoot, found = os.LookupEnv("GOROOT"); !found {
		return nil, fmt.Errorf("GOROOT is not set")
	}

	if envVal.GoPath, found = os.LookupEnv("GOPATH"); !found {
		return nil, fmt.Errorf("GOROOT is not set")
	}

	if envVal.ProjectName, found = os.LookupEnv("GOVENV_PROJECT"); !found {
		return nil, fmt.Errorf("GOROOT is not set")
	}

	return envVal, nil
}

// CreateScript create new activation script
// under the GOPATH/config.GovenvProjectScriptDir
func (goEnv *GoEnv) CreateScript(dest string) (err error) {
	tmpl := template.New("activate")

	if _, err := tmpl.Parse(scriptTemplate); err != nil {
		return err
	}

	fout, err := os.OpenFile(filepath.Join(dest, tmpl.Name()), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fout.Close()

	if err := tmpl.Execute(fout, goEnv); err != nil {
		return err
	}

	return nil
}

// IsActivated checks is activated virtual environment or not.
func IsActivated() bool {
	_, ok := os.LookupEnv("GOVENV_ENABLE")
	return ok
}
