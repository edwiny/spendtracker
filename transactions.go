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

type TransactionKey struct {
	Timestamp    time.Time
	Line         string
	Amount       float32
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

func (t Transaction) key() TransactionKey {
	amount := t.Debit
	if(t.Type == Credit) {
		amount = t.Credit
	}
	return TransactionKey{
		Timestamp: t.Timestamp, 
		Line: t.Line,
		Amount: amount,
	}
}
