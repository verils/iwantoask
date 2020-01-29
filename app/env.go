package app

import "os"

const EnvBasePath = "BASE_PATH"

var BasePath = os.Getenv(EnvBasePath)

func PrefixBasePath(path string) string {
	return BasePath + path
}
