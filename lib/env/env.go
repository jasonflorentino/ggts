package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func assertEnv() {
	vars := []string{
		"GGTS_ENV",
		"GGTS_PORT",
	}
	for _, v := range vars {
		if v == "" {
			log.Fatalf("Expected %s to be set in .env", v)
		}
	}
}

func LoadEnv() {
	err := godotenv.Overload(".env", ".env.local")
	if err != nil {
		log.Fatal("Error loading env files")
	}
	assertEnv()
}

func NotProd() bool {
	return os.Getenv("GGTS_ENV") != "production"
}

func Port() string {
	return os.Getenv("GGTS_PORT")
}
