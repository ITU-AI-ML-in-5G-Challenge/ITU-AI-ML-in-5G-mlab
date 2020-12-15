package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	var fp *os.File
	var err error
	flag.Parse()

	fp, err = os.Open(flag.Args()[0])
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 40960000)
	line, _, err := reader.ReadLine()
	hists := strings.Split(string(line), ",")

	dataCSV, err := os.Open(flag.Args()[1])
	if err != nil {
		panic(err)
	}
	defer dataCSV.Close()

	reader1 := bufio.NewReaderSize(dataCSV, 40960000)

	i := 0
	for {
		dataLine, _, err := reader1.ReadLine()

		dataList := strings.Split(string(dataLine), ",")
		newLine := []string{}
		for j, v := range dataList {
			n, _ := strconv.Atoi(hists[j])
			if n > 1 {
				newLine = append(newLine, v)
			}
		}

		// fmt.Printf("%d\t%d\n", i, len(cols))
		fmt.Fprintf(os.Stderr, "%d\t%d\n", i, len(dataList))
		fmt.Println(strings.Join(newLine, ","))
		i = i + 1
		if err == io.EOF {
			fmt.Fprintf(os.Stderr, "end")
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			panic(err)
		}
	}
}
