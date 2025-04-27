package utils

import "os"

func GetSaveDir() string {
	if homedir, err := os.UserHomeDir(); err == nil {
		return homedir + "/.citylyf/saves"
	}
	return ""
}

func GetDirFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}
