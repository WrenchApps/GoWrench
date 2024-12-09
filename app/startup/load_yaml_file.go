package startup

import (
	"io"
	"os"
)

func LoadYamlFile(pathFile string) ([]byte, error) {
	file, err := os.Open(pathFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
