package config

const (
	// GovenvEnvironmentFile is goenv's (this program's) environment file name.
	GovenvEnvironmentFile = ".govenvrc"

	// GovenvManagementDir is goenv's (this program's) workspace directory name.
	GovenvManagementDir = ".govenv"

	// GovenvGoRootsDirName is directory name which has golang version directory inside.
	GovenvGoRootsDirName = "goroots"

	// GovenvGoMasterDir is golang git repository directory name.
	GovenvGoMasterDir = "golang"

	// GovenvGoBootstrapDir is golang bootstrap builder directory name.
	GovenvGoBootstrapDir = "bootstrapper"

	// GovenvGoProjectDir is project specific data directory name
	GovenvGoProjectDir = ".project"

	// GovenvGoProjectScriptDir is project specific scripts (activate) directory name
	GovenvGoProjectScriptDir = "bin"
)
