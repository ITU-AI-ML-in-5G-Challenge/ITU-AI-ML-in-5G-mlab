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
	dic := map[string]struct{}{}
	flag.Parse()

	fp, err = os.Open(flag.Args()[0])
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 40960000)
	i := 0
	for {
		line, _, err := reader.ReadLine()
		cols := strings.Split(string(line), ",")
		for _, v := range cols {
			dic[v] = struct{}{}
		}

		// fmt.Printf("%d\t%d\n", i, len(cols))
		fmt.Fprintf(os.Stderr, "%d\t%d\n", i, len(cols))
		i = i + 1
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	isStart := true
	for k, _ := range dic {
		if !isStart {
			fmt.Print(",")
		} else {
			isStart = false
		}
		fmt.Print(k)
	}
}
