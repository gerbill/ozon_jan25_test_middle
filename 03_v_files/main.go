package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	Run(in, out)
}

type Folder struct {
	Dir            string   `json:"dir"`
	Files          []string `json:"files"`
	Folders        []Folder `json:"folders"`
	hackedFilesNum int
	isHacked       bool
}

func (f *Folder) countHackedFiles(isHacked bool) int {
	if isHacked == false {
		for _, file := range f.Files {
			if strings.HasSuffix(file, ".hack") {
				isHacked = true
				break
			}
		}
	}

	for _, folder := range f.Folders {
		f.hackedFilesNum += folder.countHackedFiles(isHacked)
	}
	if isHacked {
		return len(f.Files) + f.hackedFilesNum
	}
	return f.hackedFilesNum
}

func Run(in *bufio.Reader, out *bufio.Writer) {
	lineCounter := 0
	openBracers := 0
	var jsonBuffer strings.Builder
	for {
		lineCounter++
		line, err := in.ReadString(byte('\n'))
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read line %d. Error: %+v", lineCounter, err)
		}
		if line == "" {
			break
		}

		openBracers += strings.Count(line, "{")

		if openBracers != 0 {
			jsonBuffer.WriteString(line)
		}
		openBracers -= strings.Count(line, "}")

		if openBracers == 0 && jsonBuffer.Len() > 0 {
			jsonStr := jsonBuffer.String()
			jsonBuffer.Reset()

			var folder Folder
			err := json.Unmarshal([]byte(jsonStr), &folder)
			if err != nil {
				log.Fatalf("failed to unmarshal json from jsonStr %s. Error: %+v", jsonStr, err)
			}
			out.WriteString(fmt.Sprintf("%d\n", folder.countHackedFiles(false)))
		}
	}
}

// func countAllHackedFilesOld(folder Folder) int {
// 	folder.countHackedFiles()
// 	allHackedFilesNum := folder.hackedFilesNum
// 	folders := folder.Folders
// 	isHacked := folder.isHacked
// 	for {
// 		for _, folder := range folders {
// 			if isHacked {
// 				allHackedFilesNum += len(folder.Files)
// 			} else {
// 				folder.countHackedFiles()
// 				allHackedFilesNum += folder.hackedFilesNum
// 				isHacked = folder.isHacked
// 			}
// 			folders = folder.Folders
// 		}
// 		if len(folders) == 0 {
// 			break
// 		}
// 	}
// 	return allHackedFilesNum
// }
