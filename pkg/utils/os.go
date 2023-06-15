package utils

import (
	"bufio"
	"os"
	"path/filepath"
)

func GetFilesData(directoryPath string) (map[string][]byte, error) {
	fileData := make(map[string][]byte)

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Open file for reading
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Read file data using bufio.Scanner
		scanner := bufio.NewScanner(file)
		var data []byte
		for scanner.Scan() {
			data = append(data, scanner.Bytes()...)
		}
		if err = scanner.Err(); err != nil {
			return err
		}

		// Store file name and data in the map
		fileName := filepath.Base(path)
		fileData[fileName] = data

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileData, nil
}
