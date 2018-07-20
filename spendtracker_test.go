package spendtracker

import (
	"testing"

	"sort"
)

func testEqStringSlices(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func setup(set string) *SpendTracker {
	c := Config{}
	c.PatternsFile = "testdata/" + set + "/patterns.csv"
	c.AliasesFile = "testdata/" + set + "/aliases.csv"
	c.Period = "monthly"

	c.ANZFiles = []string{"testdata/" + set + "/anz/anz.csv"}

	tracker := New(c)
	return tracker

}

func TestPatternsLoadCorrectly(t *testing.T) {
	tracker := setup("set1")
	expected_level1tags := []string{"Living Expenses", "Non-essential Expenses"}

	err := tracker.LoadPatterns()
	if err != nil {
		t.Error("Patterns failed to load")
	}

	level1tags := tracker.pdb.Level1Tags()
	sort.Strings(level1tags)
	sort.Strings(expected_level1tags)
	if !testEqStringSlices(level1tags, expected_level1tags) {
		t.Error("expected ", expected_level1tags, " got ", level1tags)
	}
}

func TestANZLoadsTransactions(t *testing.T) {
	setup("set1")

	tracker := setup("set1")

	err := tracker.LoadPatterns()
	if err != nil {
		t.Error("Patterns failed to load")
	}
	err = tracker.ReadTransactions()
	if err != nil {
		t.Error("ANZ transactions failed to load")
	}

	found := false
	for _, trn := range tracker.trns {

		if trn.Line == "Tetsuya Restaruant" && trn.Debit == -300 {
			found = true

		}
	}
	if !found {
		t.Error("Did not find expected transaction in loaded transaction list")
	}

}
