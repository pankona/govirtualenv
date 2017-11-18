#!/bin/bash

# We need these commands are installed
GOVENV_GIT=$(which git)

# URL link to golang repository
GOVENV_GO_REPOSITORY_URL="https://go.googlesource.com/go"

# URL link to govenv command repository
GOVENV_REPOSITORY_PATH="github.com/necomeshi/govenv"

# A master repository name of golang 
GOVENV_GO_REPOSITORY_NAME="golang"

# A bootstrap golang version direcotry 
GOVENV_BOOTSTRAP_V1_4_DIR="bootstrap"

# Location of information ifle
GOVENV_INFORMATION_FILE=${HOME}/.govenvrc

# Location of management directory
GOVENV_MNGMNTDIR_PATH=${HOME}/.govenv

# A GOROOTs direcotry path
GOVENV_GOROOTS_PATH=${GOVENV_MNGMNTDIR_PATH}/goroots

# A govenv tools path
GOVENV_PATH=${GOVENV_MNGMNTDIR_PATH}/tools

# A golang master repostiroy path
GOVENV_GOLANG_MASTER_PATH=${GOVENV_GOROOTS_PATH}/${GOVENV_GO_REPOSITORY_NAME}

# A golang bootstrap direcotry path
GOVENV_GOLANG_BOOTSTRAP_PATH=${GOVENV_GOROOTS_PATH}/${GOVENV_BOOTSTRAP_V1_4_DIR}

# Supress stdout
function _pushd() {
    command pushd "$@" > /dev/null 2>&1
}
function _popd() {
    command pushd "$@" > /dev/null 2>&1
}

function CreateInformationFile() {
    local filename=${1}
    local gitPath=${2}
    local mngmtDir=${3}

    cat > ${filename} << _EOF_
{
    "managementDirectoryPath" : "${mngmtDir}",
    "gitPath" : "${gitPath}"
}
_EOF_

    return ${?}
}

function Install() {
    # Find git executable 
    if [ -z "${GOVENV_GIT}" ]; then
        echo "'git' command is not found in PATH"
        exit 255
    fi

    # Create setup information file
    # This file is used by govenv command
    CreateInformationFile ${GOVENV_INFORMATION_FILE} ${GOVENV_GIT} ${GOVENV_MNGMNTDIR_PATH}
    if [ ${?} -ne 0 ]; then
        echo "Error: Cannot create information file"
        exit 1
    fi

    # Remove already existed management directory
    if [ -e ${GOVENV_MNGMNTDIR_PATH} ]; then
        rm -rf ${GOVENV_MNGMNTDIR_PATH}
    fi

    # Create management direcotries
    mkdir -p ${GOVENV_GOROOTS_PATH}
    if [ ${?} -ne 0 ]; then
        echo "Cannot create ${GOVENV_GOROOTS_PATH}"
        exit 1
    fi

    # Clone git repository under the management directory
    echo -n "[1/4] Cloning golang git repository: "
    _pushd ${GOVENV_GOROOTS_PATH}
        ${GOVENV_GIT} clone ${GOVENV_GO_REPOSITORY_URL} ${GOVENV_GO_REPOSITORY_NAME} > /dev/null 2>&1
        result=${?}
    _popd
    if [ ${result} -eq 0 ]; then
        echo "OK"
    else
        echo "Cannot clone go repository"
        exit 2
    fi

    # Export latest go1.4 source
    echo -n "[2/4] Exporting bootstrap builder for golang: "
    _pushd ${GOVENV_GOLANG_MASTER_PATH}
        ${GOVENV_GIT} checkout -b release-branch.v1.4 remotes/origin/release-branch.go1.4 > /dev/null 2>&1 &&
        ${GOVENV_GIT} checkout-index -a --prefix=../${GOVENV_BOOTSTRAP_V1_4_DIR}/ > /dev/null 2>&1
        result=${?}
    _popd
    if [ ${result} -eq 0 ]; then
        echo "Dnoe"
    else
        echo "Cannot checkout go version 1.4"
        exit 3
    fi

    # Build bootstrap
    echo -n "[3/4] Building bootstrap builder: "
    _pushd ${GOVENV_GOLANG_BOOTSTRAP_PATH}/src
        /bin/bash ./all.bash > /dev/null 2>&1
        result=${?}
    _popd
    if [ ${result} -eq 0 ]; then
        echo "OK"
    else
        echo "Cannot build go version 1.4"
        exit 3
    fi

    # Install govenv
    echo -n "[4/4] Getting govenv tools: "
    GOROOT=${GOVENV_GOLANG_BOOTSTRAP_PATH}
    GOPATH=${GOVENV_PATH}
    ${GOROOT}/bin/go install ${GOVENV_REPOSITORY_PATH}
}

Install