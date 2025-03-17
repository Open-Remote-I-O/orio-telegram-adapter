package config

import "os"

// NOTE: investigate later on if function needs to be generic
func GetEnvOrDefault(env, default_val string) string {
	if env_val, ok := os.LookupEnv(env); ok {
		return env_val
	}
	return default_val
}
