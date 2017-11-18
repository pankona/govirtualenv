package script

import "text/template"
import "os"
import "path/filepath"

const scriptTemplate = `
#!/bin/bash
function deactivate() {
	if [ -n ${GOVENV_OLD_PATH} ]; then
		PATH=${GOVENV_OLD_PATH}
		export PATH
		unset GOVENV_OLD_PATH
	fi

	if [ -n ${GOVENV_OLD_PS1} ]; then
		PS1=${GOVENV_OLD_PS1}
		export PS1
		unset GOVENV_OLD_PS1
	fi

	if [ -n ${GOROOT} ]; then
		unset GOROOT
	fi

	if [ -n ${GOPATH} ]; then
		unset GOPATH
	fi

	unset GOENV_ENABLE
}

# Activation
if [ -z "${GOENV_ENABLE+x}" ]; then
	GOVENV_OLD_PATH=${PATH}
	GOVENV_OLD_PS1=${PS1}

	GOROOT={{.GoRoot}}
	GOPATH={{.GoPath}}

	PATH=${GOROOT}/bin:${PATH}
	PS1="({{.Project}}) ${PS1}"

	export PS1 PATH GOROOT GOPATH GOVENV_ENABLE=1
else
	echo "Already in the project. First do 'deactivate'"
fi
`

type goEnv struct {
	GoRoot  string
	GoPath  string
	Project string
}

// CreateScript create a new activation script under the path
func CreateScript(dest, goroot, gopath string) (err error) {
	tmpl := template.New("activate")

	if _, err := tmpl.Parse(scriptTemplate); err != nil {
		return err
	}

	fout, err := os.OpenFile(filepath.Join(dest, tmpl.Name()), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fout.Close()

	env := goEnv{goroot, gopath, filepath.Base(gopath)}

	if err := tmpl.Execute(fout, env); err != nil {
		return err
	}

	return nil
}
