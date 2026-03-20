package utils

import (
	"os"
)

// getEnvWithDefault: funcion que retorna la variable de entorno encontrada, si no, una default.
func GetEnvWithDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
