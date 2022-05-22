package config

import "os"

type Config struct{}

var App Config

func Env(key string) string {
	return os.Getenv(key)
}
