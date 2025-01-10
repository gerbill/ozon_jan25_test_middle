package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	Run(in, out)
}

func Run(in *bufio.Reader, out *bufio.Writer) {
	reader := bufio.NewReader(in)
	lineCounter := 0
	// items := 0
	for {
		lineCounter++
		line, err := reader.ReadString(byte('\n'))
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read line %d. Error: %+v", lineCounter, err)
		}
		if line == "" {
			break
		}
		if lineCounter == 1 {
			continue
		}
		line = strings.ReplaceAll(line, "\n", "")
		if len(line) == 1 {
			out.WriteString("0\n")
			continue
		}

		digits := make([]int, len(line))
		for i, char := range line {
			digit, err := strconv.Atoi(string(char))
			if err != nil {
				log.Fatalf("failed to convert char %c to int in line %s. Error: %+v", char, line, err)
			}
			digits[i] = digit
		}

		removeDigitIndex := len(line) - 1
		for i, d := range digits[:len(line)-1] {
			if d < digits[i+1] {
				removeDigitIndex = i
				break
			}
		}
		out.WriteString(line[:removeDigitIndex] + line[removeDigitIndex+1:] + "\n")
	}
}
