package spendtracker

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	//	"strings"
	"time"
)

type INGReader struct {
	pdb *PatternDB
}

func (r INGReader) CSVReader(filename string) ([]Transaction, error) {
	var data []Transaction
	dateForm := "02/01/2006"

	file, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return data, err

	}
	defer file.Close()

	csv := csv.NewReader(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for {

		record, err := csv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//first line of ING statements typically contains column headings
		if record[0] == "Date" {
			continue
		}

		if len(record) > 3 {
			t := Transaction{}

			t.Timestamp, err = time.Parse(dateForm, record[0])
			if err != nil {
				fmt.Println(err)
				continue
			}

			//Line item
			t.Line = TrimQuotes(record[1])
			if r.pdb.matchAccountAliases(t.Line) != nil {
				fmt.Fprintf(os.Stderr, "Ignoring aliased pattern %s\n", t.Line)
				continue
			}

			//Credit
			tmpstr := TrimQuotes(record[2])
			if len(tmpstr) == 0 {
				t.Type = Debit
				t.Credit = 0
			} else {

				tmpval, err := strconv.ParseFloat(TrimQuotes(record[2]), 32)
				if err != nil {
					fmt.Println(err)
					continue
				}
				t.Credit = float32(tmpval)
			}

			//Debit
			tmpstr = TrimQuotes(record[3])
			if len(tmpstr) == 0 {
				t.Type = Credit
				t.Debit = 0
			} else {
				tmpval, err := strconv.ParseFloat(TrimQuotes(record[3]), 32)
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
