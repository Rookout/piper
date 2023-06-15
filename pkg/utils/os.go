package utils

import (
	"os"
	"path/filepath"
)

func GetFilesData(directory string) (map[string][]byte, error) {
	fileData := make(map[string][]byte)
	var data []byte

	path, dirList, err := GetFilesInLinkDirectory(directory)
	if err != nil {
		return nil, err
	}
	for _, dir := range dirList {
		data, err = GetFileData(*path + "/" + dir)
		if err != nil {
			return nil, err
		}
		fileData[dir] = data
	}

	return fileData, nil
}

func GetFileData(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetFilesInLinkDirectory(linkPath string) (*string, []string, error) {
	realPath, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		return nil, nil, err
	}

	var fileNames []string
	err = filepath.WalkDir(realPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileNames = append(fileNames, d.Name())
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &realPath, fileNames, nil
}
