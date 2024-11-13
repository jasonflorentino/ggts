package env

import (
	"log"
	"os"
)

func assertEnv() {
	vars := []string{
		"GGTS_ENV",
		"GGTS_PORT",
	}
	for _, v := range vars {
		if v := os.Getenv(v); v == "" {
			log.Fatalf("Expected %s to be set in .env", v)
		}
	}
}

func LoadEnv() {
	// This did not work on my linux build running on Ubuntu 22.04
	// with a go1.23.3.linux-amd64 install, but maybe it's just
	// just because I'm always doing this at midnight ><
	//
	// const FILE_NAME string = ".env"
	// const KV_SEP string = "="
	// file, err := os.ReadFile(FILE_NAME)
	// if err != nil {
	// 	log.Fatal("Error loading env files")
	// }
	// lines := strings.Split(string(file), "\n")
	// for i, line := range lines {
	// 	kv := strings.Split(line, KV_SEP)
	// 	if len(kv) > 2 {
	// 		log.Fatalf("env line %d splits to more than 2 items: %v", i, kv)
	// 	}
	// 	if strings.HasPrefix(kv[1], "\"") {
	// 		log.Fatalf("quoted vals not supported: env line %d - %v", i, kv[1])
	// 	}
	// 	keyParts := strings.Split(kv[1], " ") // anything after the first space is ignored
	// 	os.Setenv(kv[0], keyParts[0])
	// }
	assertEnv()
}

func IsProd() bool {
	return os.Getenv("GGTS_ENV") == "production"
}

func NotProd() bool {
	return os.Getenv("GGTS_ENV") != "production"
}

func Port() string {
	return os.Getenv("GGTS_PORT")
}
