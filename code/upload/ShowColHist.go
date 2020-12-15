package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
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
	cols := strings.Split(string(line), ",")
	mapList := make([]map[string]struct{}, len(cols))
	for i := range mapList {
		mapList[i] = make(map[string]struct{})
	}

	i := 0
	for {
		line, _, err := reader.ReadLine()
		values := strings.Split(string(line), ",")
		for j, v := range values {
			mapList[j][v] = struct{}{}
		}

		// fmt.Printf("%d\t%d\n", i, len(cols))
		fmt.Fprintf(os.Stderr, "%d\t%d\n", i, len(values))
		i = i + 1
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	isStart := true
	for i := range mapList {
		if !isStart {
			fmt.Print(",")
		} else {
			isStart = false
		}
		fmt.Print(len(mapList[i]))
	}
}
