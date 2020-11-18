package fileoperations

import (
	"os"
)

func CreateDir(dirPath string) {
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		panic(err)
	}
}
func Remove(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		return
	}
}
