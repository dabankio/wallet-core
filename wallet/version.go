package wallet

var (
	version   = "0.1.0"
	buildTime = "no build time set"
	gitHash   = "no git hash set"
)

func GetVersion() string {
	return version
}

func GetBuildTime() string {
	return buildTime
}

func GetGitHash() string {
	return gitHash
}
