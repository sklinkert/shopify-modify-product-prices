package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func checkErr(err error, line string) {
	if err != nil {
		log.Panicf("%v (line=%q)", err, line)
	}
}

func main() {
	var path = os.Args[1]
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer file.Close()

	csvLines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	var newLines [][]string
	for i, line := range csvLines {
		if len(line) != 49 {
			log.Fatalf("Unexpected number of fields: %d", len(line))
		}

		if i == 0 {
			newLines = append(newLines, line)
			continue
		}
		handle := line[0]

		if line[20] == "" {
			log.Printf("Empty price: %+v", line)
			newLines = append(newLines, line)
			continue
		}
		price, err := strconv.ParseFloat(line[20], 64)
		checkErr(err, strings.Join(line, ","))

		if line[21] == "" {
			line[21] = line[20]
		}
		priceCompareAt, err := strconv.ParseFloat(line[21], 64)
		checkErr(err, strings.Join(line, ","))

		log.Printf("%s,%f,%f", handle, price, priceCompareAt)

		if priceCompareAt < price {
			priceCompareAt = price * 1.2
		} else {
			priceCompareAt = priceCompareAt * 1.1
		}

		line[20] = fmt.Sprintf("%.2f", price)
		line[21] = fmt.Sprintf("%.2f", priceCompareAt)

		newLines = append(newLines, line)
	}

	// Write newLines as CSV to a file
	newFile, err := os.Create("products-updated.csv")
	checkErr(err, "")
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	err = writer.WriteAll(newLines)
	checkErr(err, "")
}
