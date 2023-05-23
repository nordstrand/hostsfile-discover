package main

import (
	"os"
	"strconv"
)

type Config struct{}

var CONFIG = Config{}

func (Config) TLD() string {
	env_TLD := os.Getenv("TLD")

	if env_TLD != "" {
		return env_TLD
	} else {
		return "homestead"
	}
}

func (Config) PORT() int {
	env_PORT := os.Getenv("PORT")
	port, err := strconv.Atoi(env_PORT)
	if env_PORT != "" && err == nil {
		return port
	} else {
		return 8000
	}
}

func (Config) HOSTS_FILE_PATH() string {
	env_HOSTS_FILE_PATH := os.Getenv("HOSTS_FILE_PATH")

	if env_HOSTS_FILE_PATH != "" {
		return env_HOSTS_FILE_PATH
	} else {
		return "../../hosts"
	}
}
