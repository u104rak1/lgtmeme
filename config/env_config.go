package config

import (
	"os"

	"github.com/joho/godotenv"
)

func NewEnv() {
	if os.Getenv("ECHO_MODE") == "local" {
		if err := godotenv.Load(".env.local"); err != nil {
			panic(err)
		}
		if err := godotenv.Load(".env.local.secret"); err != nil {
			panic(err)
		}
	}
}
