package fileoperations

import (
	"bufio"
	"strings"

	"os"
)

func ReadFileLines(dirPath string, filename string) []string {
	file, err := os.Open(dirPath + "/" + filename)
	if err != nil {
		//log.Fatal(err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		//log.Fatal(err)
		return nil
	}

	return lines
}

func ExtractFileName(filepath string) string {
	slc := strings.SplitAfter(filepath, "/")
	filename := slc[len(slc)-1]
	return filename
}

func ExtractChangeName(eventName string, filepath string) string {
	slc := strings.SplitAfter(eventName, filepath+"/")
	changeName := slc[len(slc)-1]
	return changeName
}
