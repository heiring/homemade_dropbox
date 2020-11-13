package fileoperations

import (
	"bufio"
	"fmt"
	"log"
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

	fmt.Println("CreateFile: File created successfully", filename)
}

func WriteSliceToFile(filepath string, filename string, lines []string) {

	file, err := os.OpenFile(filepath+"/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	err = file.Truncate(0)
	if err != nil {
		log.Fatalf("failed truncating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)
	fmt.Print("lines to write: ")
	fmt.Println(lines)

	for _, data := range lines {
		_, _ = datawriter.WriteString(data + "\n")

	}

	datawriter.Flush()
	file.Close()
}

func CreateDir(path string, dirName string) {
	err := os.Mkdir(path+"/"+dirName, 0755)
	if err != nil {
		fmt.Println(err)
		//log
		return
	}
}
func Remove(dirPath string, filename string) {

	err := os.Remove(dirPath + "/" + filename)
	if err != nil {
		//log.Fatal(e)
		fmt.Println("remove error")
		fmt.Println(err)
	}
}
