#!/bin/bash

# We need these commands are installed
GOVENV_GIT=$(which git)

# URL link to golang repository
GOVENV_GOLANG_REPOSITORY_URL="https://go.googlesource.com/go"

# URL link to govenv command repository
GOVENV_REPOSITORY_PATH="github.com/necomeshi/govirtualenv"

# A master repository name of golang 
GOVENV_GOLANG_MASTER_NAME="golang"

# A bootstrap golang version direcotry 
GOVENV_BOOTSTRAP_DIRNAME="bootstrapper"

# A required golang version tag
GOVENV_GOLANG_REQUIRED_VERSION_TAG="go1.9.2"

# Location of information ifle
GOVENV_INFORMATION_FILE=${HOME}/.govenvrc

# Location of management directory
GOVENV_MNGMNTDIR_PATH=${HOME}/.govenv

# A GOROOTs direcotry path
GOVENV_GOROOTS_PATH=${GOVENV_MNGMNTDIR_PATH}/goroots

# A govenv tools path
GOVENV_INSTALL_PATH=${GOVENV_MNGMNTDIR_PATH}/tools

# A golang master repostiroy path
GOVENV_GOLANG_MASTER_PATH=${GOVENV_GOROOTS_PATH}/${GOVENV_GOLANG_MASTER_NAME}

# A golang bootstrap direcotry path
GOVENV_GOLANG_BOOTSTRAP_PATH=${GOVENV_GOROOTS_PATH}/${GOVENV_BOOTSTRAP_DIRNAME}

# A reqired golang version path
GOVENV_GOLANG_REQUIRED_VERSION_PATH=${GOVENV_GOROOTS_PATH}/${GOVENV_GOLANG_REQUIRED_VERSION_TAG}



# Supress stdout
function _pushd() {
    command pushd "$@" > /dev/null 2>&1
}
function _popd() {
    command pushd "$@" > /dev/null 2>&1
}

function PrintUsage() {
    echo "Thank you for installing govirtualenv!"
    echo ""
    echo "Please add a path to '${GOVENV_INSTALL_PATH}/bin' to your PATH."
    echo "More information is provided by exectuting following command."
    echo ""
    echo " $ govirtualenv --help"
    echo ""
}

function CheckPrerequisition() {
    # Find git executable 
    if [ -z "${GOVENV_GIT}" ]; then
        echo "'git' command is not found in PATH"
        return 1
    fi
    return 0
}

function InstallGolang() {
    # Create management direcotries
    mkdir -p ${GOVENV_GOROOTS_PATH}
    if [ ${?} -ne 0 ]; then
        echo "Cannot create ${GOVENV_GOROOTS_PATH}"
        return 1
    fi

    # Clone git repository under the management directory
    _pushd ${GOVENV_GOROOTS_PATH}
        ${GOVENV_GIT} clone ${GOVENV_GOLANG_REPOSITORY_URL} \
                    ${GOVENV_GOLANG_MASTER_NAME} > /dev/null 2>&1
        result=${?}
    _popd
    if [ ${result} -ne 0 ]; then
        echo "Cannot clone golang repository"
        return 2
    fi

    # Export latest go1.4 source
    _pushd ${GOVENV_GOLANG_MASTER_PATH}
        ${GOVENV_GIT} checkout remotes/origin/release-branch.go1.4 > /dev/null 2>&1 &&
        ${GOVENV_GIT} checkout-index -a --prefix=../${GOVENV_BOOTSTRAP_DIRNAME}/ > /dev/null 2>&1
        result=${?}

        # Maybe this command shows an error for anonymous branch is deleted
        ${GOVENV_GIT} checkout master > /dev/null 2>&1
    _popd
    if [ ${result} -ne 0 ]; then
        echo "Cannot checkout go version 1.4"
        return 3
    fi

    # Build bootstrap
    _pushd ${GOVENV_GOLANG_BOOTSTRAP_PATH}/src
        /bin/bash ./all.bash > /dev/null 2>&1
        result=${?}
    _popd
    if [ ${result} -ne 0 ]; then
        echo "Cannot build go version 1.4"
        return 4
    fi

    return 0
}

function InstallRequiredVersion() {
    # Install required golang version
    _pushd ${GOVENV_GOLANG_MASTER_PATH}
        ${GOVENV_GIT} checkout ${GOVENV_GOLANG_REQUIRED_VERSION_TAG} > /dev/null 2>&1 &&
        ${GOVENV_GIT} checkout-index -a --prefix=../${GOVENV_GOLANG_REQUIRED_VERSION_TAG}/ > /dev/null 2>&1
        result=${?}

        # Maybe this command shows an error for anonymous branch is deleted
        ${GOVENV_GIT} checkout master > /dev/null 2>&1
    _popd
    if [ ${result} -ne 0 ]; then
        echo "Cannot checkout go version ${GOVENV_GOLANG_REQUIRED_VERSION_TAG}"
        return 1
    fi

    _pushd ${GOVENV_GOLANG_REQUIRED_VERSION_PATH}/src
        GOROOT_BOOTSTRAP=${GOVENV_GOLANG_BOOTSTRAP_PATH} /bin/bash ./all.bash > /dev/null 2>&1
        result=${?}
    _popd
    if [ ${result} -ne 0 ]; then
        echo "Cannot build go version ${GOVENV_GOLANG_REQUIRED_VERSION_TAG}"
        return 2
    fi

    return 0
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

# Install govirtualenv
function InstallGovirtualenv() {

    mkdir -p ${GOVENV_INSTALL_PATH}

    err=$(GOROOT=${GOVENV_GOLANG_REQUIRED_VERSION_PATH} GOPATH=${GOVENV_INSTALL_PATH} \
        ${GOVENV_GOLANG_REQUIRED_VERSION_PATH}/bin/go get ${GOVENV_REPOSITORY_PATH} 2>&1)
    if [ ${?} -ne 0 ]; then
        echo "Cannot go get ${GOVENV_REPOSITORY_PATH}; ${err}"
        return 1
    fi

    err=$(GOROOT=${GOVENV_GOLANG_REQUIRED_VERSION_PATH} GOPATH=${GOVENV_INSTALL_PATH} \
        ${GOVENV_GOLANG_REQUIRED_VERSION_PATH}/bin/go install ${GOVENV_REPOSITORY_PATH} 2>&1)
    if [ ${?} -ne 0 ]; then
        echo "Cannot go install govirtualenv ${err}"
        return 1
    fi

    # Create setup information file
    # This file is used by govenv command
    CreateInformationFile ${GOVENV_INFORMATION_FILE} ${GOVENV_GIT} ${GOVENV_MNGMNTDIR_PATH}
    if [ ${?} -ne 0 ]; then
        echo "Cannot create information file"
        return 1
    fi
    return 0
}

#
# Install newly installs govirtualenv 
#
function Install() {
    # Remove already existing management directory
    if [ -e ${GOVENV_MNGMNTDIR_PATH} ]; then
        rm -rf ${GOVENV_MNGMNTDIR_PATH}
    fi

    echo -n "[1/4] Checking the prerequistition:"
    err=$(CheckPrerequisition)
    if [ ${?} -eq 0 ]; then
        echo "OK"
    else
        echo "Error: ${err}"
        exit 1
    fi

    echo -n "[2/4] Installing golang bootstrap builder:"
    err=$(InstallGolang)
    if [ ${?} -eq 0 ]; then
        echo "OK"
    else
        echo "Error: ${err}"
        exit 2
    fi

    echo -n "[3/4] Installing golang govirtualenv builder version:"
    err=$(InstallGolang)
    if [ ${?} -eq 0 ]; then
        echo "OK"
    else
        echo "Error: ${err}"
        exit 3
    fi

    echo -n "[4/4] Installing govirtualenv:"
    err=$(InstallGovirtualenv)
    if [ ${?} -eq 0 ]; then
        echo "OK"
    else
        echo "Error: ${err}"
        exit 4
    fi

    PrintUsage

    exit 0
}

Install