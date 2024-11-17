package env

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	gommonlog "github.com/labstack/gommon/log"
)

var env map[string]string
var loc *time.Location

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
		if len(kv) == 0 {
			continue
		}
		if len(kv) == 1 && kv[0] == "" {
			continue
		}
		if len(kv) != 2 {
			log.Fatalf("env line %d splits into %d items: %v", i, len(kv), kv)
		}
		if strings.HasPrefix(kv[1], "\"") {
			log.Fatalf("quoted vals not supported: env line %d - %v", i, kv[1])
		}
		keyParts := strings.Split(kv[1], " ") // anything after the first space is ignored
		env[kv[0]] = keyParts[0]
	}
	assertEnv()
	initLocation()
}

func assertEnv() {
	keys := []string{
		"GGTS_ENV",
		"GGTS_LOGFILE",
		"GGTS_PORT",
		"GGTS_TITLE",
		"GGTS_URL",
	}
	for _, k := range keys {
		if v := env[k]; v == "" {
			log.Fatalf("Expected %s to be set in .env", k)
		}
	}
}

func initLocation() {
	l, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("error loading location %s", err)
	}
	loc = l
}

func Location() *time.Location {
	return loc
}

func IsProd() bool {
	return env["GGTS_ENV"] == "production"
}

func NotProd() bool {
	return env["GGTS_ENV"] != "production"
}

func LogFile() string {
	return env["GGTS_LOGFILE"]
}

func LogLevel() gommonlog.Lvl {
	l := env["GGTS_LOGLEVEL"]
	if l == "" {
		return gommonlog.Lvl(0)
	}
	li, err := strconv.Atoi(l)
	if err != nil {
		return gommonlog.Lvl(0)
	}
	return gommonlog.Lvl(li)
}

func Port() string {
	return env["GGTS_PORT"]
}

func Title() string {
	return env["GGTS_TITLE"]
}

func URL() string {
	return env["GGTS_URL"]
}
