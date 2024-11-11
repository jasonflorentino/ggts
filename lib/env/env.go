package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Overload(".env", ".env.local")
	if err != nil {
		log.Fatal("Error loading env files")
	}
}

func NotProd() bool {
	return os.Getenv("GGTS_ENV") != "production"
}
