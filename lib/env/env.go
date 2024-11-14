package env

import (
	"log"
	"os"
	"strings"
)

var env map[string]string

func LoadEnv() {
	const FILE_NAME string = ".env"
	const KV_SEP string = "="
	env = make(map[string]string)
	file, err := os.ReadFile(FILE_NAME)
	if err != nil {
		log.Fatal("Error loading env files")
	}
	lines := strings.Split(string(file), "\n")
	for i, line := range lines {
		kv := strings.Split(line, KV_SEP)
		if len(kv) > 2 {
			log.Fatalf("env line %d splits to more than 2 items: %v", i, kv)
		}
		if strings.HasPrefix(kv[1], "\"") {
			log.Fatalf("quoted vals not supported: env line %d - %v", i, kv[1])
		}
		keyParts := strings.Split(kv[1], " ") // anything after the first space is ignored
		env[kv[0]] = keyParts[0]
	}
	assertEnv()
}

func assertEnv() {
	vars := []string{
		"GGTS_ENV",
		"GGTS_PORT",
	}
	for _, v := range vars {
		if v := env[v]; v == "" {
			log.Fatalf("Expected %s to be set in .env", v)
		}
	}
}

func IsProd() bool {
	return env["GGTS_ENV"] == "production"
}

func NotProd() bool {
	return env["GGTS_ENV"] != "production"
}

func Port() string {
	return env["GGTS_PORT"]
}
