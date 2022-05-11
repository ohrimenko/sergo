package config

import "os"

type data struct{}

var Config data

// Config func to get env value
func Env(key string) string {
	return os.Getenv(key)
}
