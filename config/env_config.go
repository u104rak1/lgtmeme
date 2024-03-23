package config

import (
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	if os.Getenv("ECHO_MODE") == "local" {
		if err := godotenv.Load(".env.local"); err != nil {
			panic(err)
		}
		if err := godotenv.Load(".env.local.private"); err != nil {
			panic(err)
		}
	}
}
