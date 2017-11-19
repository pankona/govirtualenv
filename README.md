# govirtualenv 

**govirtualenv** - Creates a golang virtual environment.

## Description

**govirtualenv** is a tool for creating a virtual environment for golang project, 
which is like python command line tool *virtualenv*.

**govirtualenv** is makes you easy to install/remove golang version,
and it also easy to manage project specific GOPATH and GOROOT.

More information about what this tool can do, see **Usage** section below.


## Prerequistition

**govirtualenv** needs following applications.
Please install these applications before start installation.

- git
- curl (need for installation for Linux/MacOS)

## How to install 

### Install to Linux/MacOS

```bash
$ bash < <(curl -s -S -L https://raw.githubusercontent.com/necomeshi/govirtualenv/master/install.sh)
```

If above command was success, please follows an instruction which was output after the command.

### Install to Windows

Sorry! Currently, this tool supports Linux/MacOS only!
If you like this command, and you want to use this on Windows environment,
please tell me! I will creates a windows version as soon as possible!

## Usages

```bash
Usage: govirtualenv <flags> <subcommand> <subcommand args>

Subcommands:
    available        Show available golang version
    commands         list all command names
    flags            describe all known top-level flags
    help             describe subcommands and their syntax
    install          Install new golang version
    remove           Remove installed golang version
```

### Manage your system's golang version

#### Show available golang version

**govirtualenv available** command shows installable golang versions.

```bash
$ govirtualenv available
go1
go1.0.1
go1.0.2
go1.0.3
go1.1
# Output continue....
weekly.2012-03-13
weekly.2012-03-22
weekly.2012-03-27
```

#### Show installed golang version

**govirtualenv list** command shows already installed golang versions.

```bash
$ govirtualenv list
go1.7.5
```

#### Install new golang version

**govirtualenv install** command installs new golang version.

```bash
# govirtualenv install GOLANGVERSION
$ govirtualenv install go1.9.1
```

#### Remove installed golang version

**govirtualenv remove** command removes installed golang version.

```bash
# govirtualenv remove INSTALLED_GOLANGVERSION
$ govirtualenv remove go1.9.1
```

### Manage a project specific golang environment(GOROOT, GOPATH)

#### Start new project

**govirtualenv startproject** command create project directory and some convinient scripts

```bash
# govirtualenv startproject INSTALLED_GOLANGVERSION PROEJCTNAME
$ govirtualenv startproject go1.9.1 MyProject
Project MyProject created.

$ ls -a MyProject
./        ../       .project/

$ ls -a MyProject/.project/bin
./        ../       activate
```

#### Activate / Deactivate  virutal environment (Linux/MacOS)

After executed **govirtualenv startproject** command, 
following command activates virtual environment.

```bash
# Activate virtual environment. <ProjectName> is the command argument of startproject.
$ source <ProjectName>/.project/bin/activate

# Virtual environment is now activated
(<ProjectName>) $
```

After above command, GOPATH and GOROOT are set to your environment variables automatically.
Also, path `${GOROOT}/bin` is also added to your PATH. 


If you want to deactivate virutal environment, execute

```bash
(ProjectName) $ deactivate
```

#### Change project's GOROOT(golang version)

After activate virutal environment, **govirtualenv change** command 
changes your project GOROOT's to given golang version.

```bash
(ProjectName) $ govirtualenv list
go1.9.5
go1.7.5

(ProjectName) $ govirtualenv change go1.9.5
```

## Author

necomeshi