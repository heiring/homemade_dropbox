package etc

import "strings"

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
