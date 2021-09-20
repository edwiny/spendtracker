package main

import (
	"flag"
	"fmt"
	st "github.com/edwiny/spendtracker"
	"io/ioutil"
	"log"
	"strings"
)

func expandDir(dirname, suffix string) []string {
	var paths []string
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), suffix) {
			continue
		}
		paths = append(paths, dirname+"/"+file.Name())

	}
	return paths

}

func main() {

	var optionANZDir, optionINGDir, optionPatternFile, optionAliasesFile string
	var optionPeriod string

	flag.StringVar(&optionANZDir, "inputdir-anz", "", "directory containing statement files from ANZ")
	flag.StringVar(&optionINGDir, "inputdir-ing", "", "directory containing statement files from ING")
	flag.StringVar(&optionPatternFile, "patternfile", "patterns.csv", "name of file that maps patterns to tags")
	flag.StringVar(&optionAliasesFile, "aliasesfile", "aliases.csv", "name of file for identifying accounts that belong to you")
	flag.StringVar(&optionPeriod, "period", "monthly", "period to aggregate to. One of weekly, monthly, yearly")

	flag.Parse()

	c := st.Config{}
	c.PatternsFile = optionPatternFile
	c.AliasesFile = optionAliasesFile
	c.Period = optionPeriod

	if len(optionANZDir) > 0 {
		c.ANZFiles = expandDir(optionANZDir, ".csv")

	}

	if len(optionINGDir) > 0 {
		c.INGFiles = expandDir(optionINGDir, ".csv")
	}

	tracker := st.New(c)
	err := tracker.LoadPatterns()
	if err != nil {
		log.Fatal(err)
	}

	err = tracker.ReadTransactions()
	if err != nil {
		log.Fatal(err)
	}
	tracker.BuildIndicies()

	rep := tracker.BuildReport()
	fmt.Printf("%s", rep.PrintReport("csv"))

}
