package main

import (
	"bufio"
	"encoding/json"
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
		fmt.Fprintf(os.Stderr, "err1\n")
		panic(err)
	}
	defer fp.Close()

	fp2, err := os.Open(flag.Args()[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "err2\n")
		panic(err)
	}
	defer fp2.Close()
	reader2 := bufio.NewReaderSize(fp2, 40960000)
	line2, _, err := reader2.ReadLine()
	// カテゴリカル列か "0"なら数字 "1"ならカテゴリカル（string型なので注意）
	isCategoricalCol := strings.Split(string(line2), ",")

	reader := bufio.NewReaderSize(fp, 40960000)
	line, _, err := reader.ReadLine()
	// カラム列
	cols := strings.Split(string(line), ",")
	// カラムごとに辞書用意 stringを数字に置き換え
	mapList := make([]map[string]int, len(cols))
	currentIndex := make([]int, len(cols))

	for i := range currentIndex {
		currentIndex[i] = -5000
	}

	for i := range mapList {
		mapList[i] = make(map[string]int)
	}

	i := 0

	for {
		line, _, err := reader.ReadLine()
		values := strings.Split(string(line), ",")
		for j, v := range values {
			if isCategoricalCol[j] == "0" {
				continue
			}
			_, err := strconv.ParseFloat(v, 64)
			if err == nil {
				continue
			}

			if _, ok := mapList[j][v]; !ok {
				mapList[j][v] = currentIndex[j]
				currentIndex[j] = currentIndex[j] - 1
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

		//値の出力
		isStart := true
		for j := range values {
			if !isStart {
				fmt.Print(",")
			} else {
				isStart = false
			}
			if isCategoricalCol[j] == "0" {
				fmt.Print(values[j])
			} else {
				_, err := strconv.ParseFloat(values[j], 64)
				if err == nil {
					fmt.Print(values[j])
				} else {
					fmt.Print(mapList[j][values[j]])
				}
			}

		}
		fmt.Print("\n")
	}

	file, err := os.Create(fmt.Sprintf("map_%s.json", flag.Args()[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "err\n")
		return
	}
	defer file.Close()

	obj, err := json.Marshal(mapList)
	fmt.Fprintf(file, string(obj))

}
