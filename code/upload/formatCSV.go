package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
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
	sort.StringSlice(cols).Sort()

	col_index_map := map[string]int{}
	for i, col := range cols {
		col_index_map[col] = i
	}
	newLine := make([]string, len(cols))

	dataCSV, err := os.Open(flag.Args()[1])
	if err != nil {
		panic(err)
	}
	defer dataCSV.Close()
	labelCSV, err := os.Open(flag.Args()[2])
	if err != nil {
		panic(err)
	}
	defer labelCSV.Close()

	fmt.Println(strings.Join(cols, ","))

	reader1 := bufio.NewReaderSize(dataCSV, 40960000)
	reader2 := bufio.NewReaderSize(labelCSV, 40960000)
	i := 0
	for {
		dataLine, _, err := reader1.ReadLine()
		labelLine, _, err := reader2.ReadLine()

		dataList := strings.Split(string(dataLine), ",")
		labelList := strings.Split(string(labelLine), ",")
		for j, _ := range dataList {
			newLine[col_index_map[labelList[j]]] = dataList[j]
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

	// isStart := true
	// for _, k := range cols {
	// 	if !isStart {
	// 		fmt.Print(",")
	// 	} else {
	// 		isStart = false
	// 	}
	// 	fmt.Print(k)
	// }

}
