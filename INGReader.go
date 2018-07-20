package spendtracker

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func INGReaderFunc(filename string) ([]Transaction, error) {
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

		//first line of ING statements typically contains column headings
		if elems[0] == "Date" {
			continue
		}

		if len(elems) > 3 {
			t := Transaction{}

			t.Timestamp, err = time.Parse(dateForm, elems[0])
			if err != nil {
				fmt.Println(err)
				continue
			}

			//Line item
			t.Line = TrimQuotes(elems[1])

			//Credit
			tmpstr := TrimQuotes(elems[2])
			if len(tmpstr) == 0 {
				t.Type = Debit
				t.Credit = 0
			} else {

				tmpval, err := strconv.ParseFloat(TrimQuotes(elems[2]), 32)
				if err != nil {
					fmt.Println(err)
					continue
				}
				t.Credit = float32(tmpval)
			}

			//Debit
			tmpstr = TrimQuotes(elems[3])
			if len(tmpstr) == 0 {
				t.Type = Credit
				t.Debit = 0
			} else {
				tmpval, err := strconv.ParseFloat(TrimQuotes(elems[3]), 32)
				if err != nil {
					fmt.Println(err)
					continue
				}
				t.Debit = float32(tmpval)
			}

			data = append(data, t)
		}
	}
	return data, nil
}
