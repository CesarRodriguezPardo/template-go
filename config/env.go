package config

import (
	"errors"
	"log"
	"os"
	"strings"
)

func GetEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}

func CheckMissingEnv(listedVars []string) {
	missing := []string{}

	for _, v := range listedVars {
		_, set := os.LookupEnv(v)
		if !set {
			missing = append(missing, v)
		}
	}

	if len(missing) != 0 {
		log.Fatal("Missing variables: "+strings.Join(missing, ", "), errors.New("Missing variables"))
	}
}
