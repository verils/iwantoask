package app

import "os"

const (
	Version = "0.1.2"

	EnvBasePath = "WEB_BASE_PATH"

	EnvMysqlHost     = "MYSQL_HOST"
	EnvMysqlUsername = "MYSQL_USERNAME"
	EnvMysqlPassword = "MYSQL_PASSWORD"
)

var BasePath = os.Getenv(EnvBasePath)

func BasePathPrefix(path string) string {
	return BasePath + path
}
