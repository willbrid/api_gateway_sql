package file

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

type Buffer struct {
	StartLine int
	EndLine   int
	Lines     [][]string
}

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

func ReadCSVInBuffer(file multipart.File, bufferSize int) ([]Buffer, error) {
	var (
		buffers []Buffer
	)

	reader := csv.NewReader(file)
	numLine := 0

	for {
		buffer := Buffer{
			StartLine: bufferSize*numLine + 1,
			Lines:     make([][]string, 0, bufferSize),
		}

		for i := 0; i < bufferSize; i++ {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			buffer.Lines = append(buffer.Lines, strings.Split(record[0], ";"))
		}

		if len(buffer.Lines) == 0 {
			break
		}

		buffer.EndLine = buffer.StartLine + len(buffer.Lines) - 1
		buffers = append(buffers, buffer)
		numLine++
	}

	return buffers, nil
}
