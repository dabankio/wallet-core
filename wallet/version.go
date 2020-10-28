package wallet

import "fmt"

var (
	version   = "0.1.0"
	buildTime = "no build time set"
	gitHash   = "no git hash set"

	// Debug 模式下会打印关键过程
	Debug = false
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

func debugPrint(msg ...interface{}) {
	if Debug {
		fmt.Println("[sdk-debug]", msg)
	}
}
