package app

import "os"

const EnvBasePath = "BASE_PATH"

var BasePath = os.Getenv(EnvBasePath)

func BasePathPrefix(path string) string {
	return BasePath + path
}
