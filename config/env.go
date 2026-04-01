package config

import (
	"fmt"
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

func CheckMissingEnv(listedVars []string) error {
	missing := []string{}

	for _, v := range listedVars {
		_, set := os.LookupEnv(v)
		if !set {
			missing = append(missing, v)
		}
	}

	if len(missing) != 0 {
		return fmt.Errorf("missing variables: " + strings.Join(missing, ", "))
	}

	return nil
}
