package spendtracker

import (
	"fmt"
	"sort"
	"time"
)

type Config struct {
	ANZFiles, INGFiles        []string
	PatternsFile, AliasesFile string
	Period                    string //weekly, monthly, yearly
	Aggregation               string //sum, avg

}

type SpendTracker struct {
	cfg Config
	pdb *PatternDB
	//index into the trns array, keyed by period and level 1 tag
	idxTag1 map[IndexTag1Key][]int
	//index into the trns array, keyed by period, level 1 tag, and level 2 tag
	idxTag2 map[IndexTag2Key][]int
	trns    []Transaction
}

type IndexTag1Key struct {
	period, level1tag string
}

type IndexTag2Key struct {
	period, level1tag, level2tag string
}

func TrimQuotes(s string) string {
	first := -1
	last := -1

	r := []rune(s)

	for i := 0; i < len(r); i++ {
		if r[i] == '"' {
			first = i
			break
		}
	}
	for i := len(r) - 1; i > 0; i-- {
		if r[i] == '"' {
			last = i
			break
		}
	}

	if first == -1 || last == -1 {
		return s
	}
	return (string(r[first+1 : last]))
}

func New(cfg Config) *SpendTracker {
	st := SpendTracker{}
	st.cfg = cfg
	st.idxTag1 = make(map[IndexTag1Key][]int)
	st.idxTag2 = make(map[IndexTag2Key][]int)
	if len(cfg.Period) == 0 {
		st.cfg.Period = "monthly"
	}
	if len(cfg.Aggregation) == 0 {
		st.cfg.Aggregation = "sum"
	}
	return (&st)

}

func (st *SpendTracker) LoadPatterns() error {
	var err error
	st.pdb, err = NewDB(st.cfg.PatternsFile, st.cfg.AliasesFile)
	return err
}

/* normalise the timestamp into a string contaiing the year and month. */
func NormaliseDate(timestamp time.Time, period string) string {
	if period == "weekly" {
		year, week := timestamp.ISOWeek()
		return fmt.Sprintf("%4d-%02d", year, week)
	} else if period == "monthly" {
		return fmt.Sprintf("%4d-%02d", timestamp.Year(), int(timestamp.Month()))
	} else if period == "yearly" {
		return fmt.Sprintf("%4d", timestamp.Year())
	}
	return ""
}

func (st *SpendTracker) ReadTransactions() error {

	//read in all the transactions from file
	for _, filename := range st.cfg.ANZFiles {
		tmp_trns, err := ANZReaderFunc(filename)
		fmt.Printf("%s: %d transactions\n", filename, len(tmp_trns))
		if err != nil {
			return err
		}
		st.trns = append(st.trns, tmp_trns...)

	}

	for _, filename := range st.cfg.INGFiles {
		tmp_trns, err := INGReaderFunc(filename)
		fmt.Printf("%s: %d transactions\n", filename, len(tmp_trns))
		if err != nil {
			return err
		}
		st.trns = append(st.trns, tmp_trns...)

	}

	//match and tag transactions
	for i, trn := range st.trns {
		tags := st.pdb.matchTags(trn.Line)
		if tags != nil && len(tags) == 2 {
			st.trns[i].Level1Tag = tags[0]
			st.trns[i].Level2Tag = tags[1]
		} else {
			st.trns[i].Level1Tag = "Unknown"
			st.trns[i].Level1Tag = "Unknown"
		}
	}

	return nil
}

//BuildIndicies creates a internal index of transactions by
//month-normalised date and tags.
func (st *SpendTracker) BuildIndicies() {

	//todo: pregenerate the list of possible periods in stead of getting it from the data
	for i, t := range st.trns {
		dateKey := NormaliseDate(t.Timestamp, st.cfg.Period)
		//fmt.Printf("[%s] period for %v is [%s]\n", st.cfg.Period, t.Timestamp, dateKey)

		k1 := IndexTag1Key{
			period:    dateKey,
			level1tag: t.Level1Tag,
		}
		//fmt.Printf("Adding to idex1: %v [%s] [%s] [%v]\n", k1, k1.period, k1.level1tag, t)
		st.idxTag1[k1] = append(st.idxTag1[k1], i)

		k2 := IndexTag2Key{
			period:    dateKey,
			level1tag: t.Level1Tag,
			level2tag: t.Level2Tag,
		}
		st.idxTag2[k2] = append(st.idxTag2[k2], i)
	}
}

//Periods returns the report column headers
func (st *SpendTracker) Periods() []string {
	var data []string

	tmpkeys := make(map[string]bool)

	for k := range st.idxTag1 {
		tmpkeys[k.period] = true
	}

	for k := range tmpkeys {
		data = append(data, k)
	}
	sort.Strings(data)

	return data
}
