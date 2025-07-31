package startup

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"wrench/app"
)

func LoadYamlFile(pathFile string) (map[string][]byte, error) {

	byteResul := make(map[string][]byte)
	paths := strings.Split(pathFile, ",")

	for _, path := range paths {

		log.Printf("Loading config file %s", path)
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		byteResul[path] = data
	}

	return byteResul, nil
}

func GetFileConfigPath() (string, error) {
	configPath := os.Getenv(app.ENV_PATH_FILE_CONFIG)
	configFolder := os.Getenv(app.ENV_PATH_FOLDER_CONFIG)

	if len(configPath) > 0 && len(configFolder) > 0 {
		return "", errors.New("should use ENV_PATH_FILE_CONFIG or ENV_PATH_FOLDER_CONFIG not both")
	}

	if len(configPath) == 0 && len(configFolder) == 0 {
		configPath = "../../configApp.yaml"
	}

	if len(configPath) > 0 {
		return configPath, nil
	} else if len(configFolder) > 0 {
		return getFilesFromFolder(configFolder), nil
	}

	return "", errors.New("can't find files to load")
}

func getFilesFromFolder(rootFolder string) string {
	var pathConcate string

	err := filepath.WalkDir(rootFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories if you want only files
		if d.IsDir() {
			return nil
		}

		if len(pathConcate) > 0 {
			pathConcate += fmt.Sprintf(",%s", path)
		} else {
			pathConcate = path
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return pathConcate
}
