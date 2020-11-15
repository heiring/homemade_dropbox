package fileoperations

import (
	"os"
)

func CreateDir(path string, dirName string) {
	err := os.Mkdir(path+"/"+dirName, 0755)
	if err != nil {
		panic(err)
	}
}
func Remove(dirPath string, filename string) {

	err := os.Remove(dirPath + "/" + filename)
	if err != nil {
		return
	}
}
