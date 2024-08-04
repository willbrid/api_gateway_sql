package file

import (
	"api-gateway-sql/logging"

	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func CreateConfigFileForTesting(configContent string) (string, error) {
	configFile, err := os.CreateTemp("", "testconfig-*.yaml")
	if err != nil {
		return "", fmt.Errorf("unable to create temp file : %s", err.Error())
	}

	_, err = configFile.WriteString(configContent)
	if err != nil {
		return "", fmt.Errorf("unable to write to temp file : %s", err.Error())
	}
	configFile.Close()

	return configFile.Name(), nil
}

func ReadCSVInBuffer(filePath string, bufferSize int) ([][][]string, error) {
	var (
		buffers [][][]string
		file    *os.File
		err     error = nil
	)

	file, err = os.Open(filePath)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		buffer := make([][]string, bufferSize)

		for i := 0; i < bufferSize; i++ {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			buffer = append(buffer, record)
		}

		if len(buffer) == 0 {
			break
		}

		buffers = append(buffers, buffer)
	}

	return buffers, nil
}
