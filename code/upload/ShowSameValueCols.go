package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type SafeStdout struct {
	sync.Mutex
}

func (m *SafeStdout) Write(p []byte) (n int, err error) {
	m.Lock()
	defer m.Unlock()
	return os.Stdout.Write(p)
}

func main() {
	var fp *os.File
	var err error
	flag.Parse()
	safeStdout := new(SafeStdout)

	fp, err = os.Open(flag.Args()[0])
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 40960000)
	line, _, err := reader.ReadLine()
	cols := strings.Split(string(line), ",")

	table := [][]string{}
	i := 0

	for {
		line, _, err = reader.ReadLine()
		if len(string(line)) == 0 {
			if err == io.EOF {
				break
			}
			continue
		}
		values := strings.Split(string(line), ",")
		vl := make([]string, len(values))
		if len(values) != len(cols) {
			fmt.Fprintln(os.Stderr, "error")
			fmt.Fprintln(os.Stderr, string(line))
			return
		}
		copy(vl, values)
		table = append(table, vl)

		fmt.Fprintf(os.Stderr, "%d\t%d\n", i, len(cols))
		i = i + 1
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	//Generate Matrix
	fmt.Fprintln(os.Stderr, "Generate Matrix")
	corMatrix := make([][]bool, len(cols))
	for i := range corMatrix {
		corMatrix[i] = make([]bool, len(cols))
	}
	fmt.Fprintln(os.Stderr, "Generated Matrix")

	limit := make(chan struct{}, 10)
	var wg sync.WaitGroup
	for i := range cols {
		wg.Add(1)
		f := func(i int) {
			limit <- struct{}{}
			defer wg.Done()
			for j := range cols {
				corMatrix[i][j] = false
				if i <= j {
					break
				}
				for t := range table {
					corMatrix[i][j] = true

					if table[t][i] != table[t][j] {
						corMatrix[i][j] = false
					}
				}
				if corMatrix[i][j] {
					fmt.Fprintf(safeStdout, "%d:%d\n", i, j)
				}
				fmt.Fprintf(os.Stderr, "%d\t%d\n", i, j)
			}
			<-limit
		}
		go f(i)
		time.Sleep(20 * time.Millisecond)
	}
}
