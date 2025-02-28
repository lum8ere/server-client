package smart_context

import "os"

const (
	LocalEnvName = "local"
)

func IsRealEnv() bool {
	return os.Getenv("LOCAL") == ""
}

func EnvName() string {
	return os.Getenv("ENV_NAME")
}
