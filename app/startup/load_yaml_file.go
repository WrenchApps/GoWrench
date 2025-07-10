package startup

import (
	"io"
	"os"
	"strings"
	"wrench/app"
)

func LoadYamlFile(pathFile string) (map[string][]byte, error) {

	byteResul := make(map[string][]byte)
	paths := strings.Split(pathFile, ",")

	for _, path := range paths {

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

func GetFileConfigPath() string {
	configPath := os.Getenv(app.ENV_PATH_FILE_CONFIG)
	if len(configPath) == 0 {
		configPath = "../../configApp.yaml"
	}
	return configPath
}
