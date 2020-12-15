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
	// カラム列
	cols := strings.Split(string(line), ",")
	// カラムごとにカテゴリカル変数かフラグ
	isCategorical := make([]bool, len(cols))

	for i := range isCategorical {
		isCategorical[i] = false
	}

	i := 0
	for {
		line, _, err := reader.ReadLine()
		values := strings.Split(string(line), ",")
		for j, v := range values {
			_, err := strconv.ParseFloat(v, 64)
			if err != nil {
				isCategorical[j] = true
			}
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
	for i := range isCategorical {
		if !isStart {
			fmt.Print(",")
		} else {
			isStart = false
		}
		if isCategorical[i] == true {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	}
}
