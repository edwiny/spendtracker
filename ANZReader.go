package spendtracker

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ANZReaderFunc(filename string) ([]Transaction, error) {
	var data []Transaction
	dateForm := "02/01/2006"

	file, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return data, err

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		elems := strings.Split(line, ",")
		if len(elems) > 2 {
			t := Transaction{}

			t.Timestamp, err = time.Parse(dateForm, TrimQuotes(elems[0]))
			if err != nil {
				fmt.Println(err)
				continue
			}
			tmpval, err := strconv.ParseFloat(TrimQuotes(elems[1]), 32)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if tmpval > 0 {
				t.Credit = float32(tmpval)
				t.Type = Credit

			} else {
				t.Debit = float32(tmpval)
				t.Type = Debit

			}
			t.Line = TrimQuotes(elems[2])

			data = append(data, t)
		}
	}
	return data, nil
}
