package fileoperations

import (
	"fmt"
	"os"
)

func CreateFile(filename string) {
	var _, err = os.Stat(filename)

	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
	} else {
		fmt.Println("CreateFile: File already exists!", filename)
		return
	}

}

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
