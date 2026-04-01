package utils

import (
	"CesarRodriguezPardo/template-go/config"
	"errors"
	"io"
	"mime/multipart"
	"os"
)

func SaveFile(file multipart.File, filename string) error {
	path := config.GetEnvOrDefault("STORAGE_PATH", "app/storage")

	filepath := path + "/" + filename

	out, err := os.Create(filepath)
	if err != nil {
		return errors.New("Error creando el archivo: " + err.Error())
	}
	defer out.Close()

	write, err := io.Copy(out, file)
	if err != nil {
		return errors.New("Error copiando el archivo: " + err.Error())
	}
	if write == 0 {
		return errors.New("Error copiando el archivo: esta vacio.")
	}

	return nil
}

func FileExists(routeFile string) bool {
	_, err := os.Stat(routeFile)
	return !os.IsNotExist(err)
}

func ReadFile(routeFile string) (*os.File, error) {
	file, err := os.Open(routeFile)
	if err != nil {
		return nil, errors.New("Error leyendo el archivo: " + err.Error())
	}
	return file, nil
}
