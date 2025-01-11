package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	Run(in, out)
}

type Packet struct {
	ExpectedIntsNum  int
	ProvidedNumbers  string
	ConvertedNumbers []int
	RawData          string
	RawNumbers       []int
}

func (p *Packet) IsValid() bool {
	if strings.HasPrefix(p.RawData, " ") || strings.HasSuffix(p.RawData, " ") {
		return false
	}
	for _, rawNumber := range strings.Split(p.ProvidedNumbers, " ") {
		number, err := strconv.Atoi(rawNumber)
		if err != nil {
			log.Fatalf("[validation] failed to convert rawNumber %s to int. Error: %+v. Packet: %+v", rawNumber, err, p)
		}
		p.ConvertedNumbers = append(p.ConvertedNumbers, number)
	}
	sort.Slice(p.ConvertedNumbers, func(i, j int) bool {
		return p.ConvertedNumbers[i] < p.ConvertedNumbers[j]
	})
	rawDataSlice := strings.Split(p.RawData, " ")
	if len(rawDataSlice) != p.ExpectedIntsNum {
		return false
	}
	for _, rawNumber := range rawDataSlice {
		if rawNumber[0] == '0' || strings.HasPrefix(rawNumber, "-0") {
			return false
		}
		number, err := strconv.Atoi(rawNumber)
		if err != nil {
			return false
		}
		p.RawNumbers = append(p.RawNumbers, number)
	}
	for i, number := range p.RawNumbers {
		if number != p.ConvertedNumbers[i] {
			return false
		}
	}
	return true
}

func Run(in *bufio.Reader, out *bufio.Writer) {
	lineNumber := 0
	var p Packet
	for {
		lineNumber++
		line, err := in.ReadString(byte('\n'))
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read line. Error: %+v", err)
		}
		if line == "" {
			break
		}
		line = line[:len(line)-1]
		if lineNumber == 1 {
			continue
		}
		l := lineNumber - 2

		switch l % 3 {
		case 0:
			expectedIntsNum, err := strconv.Atoi(line)
			if err != nil {
				log.Fatalf("[stage 1] failed to convert line %s to int. Error: %+v", line, err)
			}
			p = Packet{
				ExpectedIntsNum:  expectedIntsNum,
				ConvertedNumbers: make([]int, 0),
				RawNumbers:       make([]int, 0),
			}
			continue
		case 1:
			p.ProvidedNumbers = line
			continue
		case 2:
			p.RawData = line
			if p.IsValid() {
				fmt.Fprint(out, "yes\n")
			} else {
				fmt.Fprint(out, "no\n")
			}
			continue
		}
	}
}
