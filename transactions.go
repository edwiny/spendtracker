package spendtracker

import (
	"fmt"
	"time"
)

type TrnType int

const (
	Debit  = 0
	Credit = 1
)

type Transaction struct {
	Timestamp    time.Time
	Credit       float32
	Debit        float32
	Type         TrnType
	Line         string
	YearAndMonth string
	Level1Tag    string
	Level2Tag    string
}

func (t Transaction) String() string {

	return fmt.Sprintf("%s: %s cr %0.2f dt %0.2f [%s] tag1: %s, tag2: %s",
		t.Timestamp,
		t.YearAndMonth,
		t.Credit,
		t.Debit,
		t.Line,
		t.Level1Tag,
		t.Level2Tag)

}
