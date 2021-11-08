package spendtracker

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Pattern struct {
	regexStr string
	re       *regexp.Regexp
	tags     []string
}

type TagCombo struct {
	Level1Tag, Level2Tag string
}

type PatternDB struct {
	patterns       []Pattern
	accountAliases []Pattern
	tags           map[TagCombo]bool
}

func (p *PatternDB) matchTags(line string) []string {

	for i := 0; i < len(p.patterns); i++ {
		if p.patterns[i].re != nil && p.patterns[i].re.MatchString(line) {
			return (p.patterns[i].tags)
		}
	}
	return nil
}

func (p *PatternDB) matchAccountAliases(line string) []string {

	for i := 0; i < len(p.accountAliases); i++ {
		if p.accountAliases[i].re.MatchString(line) {
			return (p.accountAliases[i].tags)
		}
	}
	return nil
}

func NewDB(patternsFile, aliasesFile string) (*PatternDB, error) {
	pdb := PatternDB{}
	var err error
	pdb.patterns, err = parsePatterns(patternsFile)
	if err != nil {
		return nil, err
	}

	pdb.patterns = append(pdb.patterns, Pattern{
		regexStr: "",
		re: nil,
		tags: []string{UNKNOWN_CATEGORY, UNKNOWN_SUBCATEGORY},
		
	})
	pdb.accountAliases, err = parsePatterns(aliasesFile)
	if err != nil {
		return nil, err
	}



	pdb.tags = make(map[TagCombo]bool)

	//use keys of a map to collect unique tags
	for _, p := range pdb.patterns {
		if len(p.tags) > 1 {
			var t TagCombo
			t.Level1Tag = p.tags[0]
			t.Level2Tag = p.tags[1]

			pdb.tags[t] = true
		}

	}
	return &pdb, nil
}

func (p *PatternDB) Level1Tags() []string {
	var data []string

	keys := make(map[string]bool)

	for k := range p.tags {
		keys[k.Level1Tag] = true
	}

	for k := range keys {
		data = append(data, k)
	}

	return data
}

func (p *PatternDB) Level2Tags(level1Tag string) []string {
	var data []string

	for k := range p.tags {
		if k.Level1Tag == level1Tag {
			data = append(data, k.Level2Tag)
		}

	}
	return data
}

func parsePatterns(filename string) ([]Pattern, error) {
	var data []Pattern

	file, err := os.Open(filename)

	if err != nil {
		return data, err

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		elems := strings.Split(line, ",")

		if(strings.HasPrefix(line, "#") ||
			len(strings.Trim(elems[0], " ")) == 0 ||
			len(strings.Trim(elems[1], " ")) == 0 ) {
			continue
		}

		if len(elems) >= 3 {
			p := Pattern{}
			p.tags = make([]string, 2)
			p.tags[0] = strings.Trim(elems[0], " ")
			p.tags[1] = strings.Trim(elems[1], " ")
			p.regexStr = strings.Trim(elems[2], " ")
			if len(p.regexStr) == 0 {
				continue
			}
				
			
			p.re, err = regexp.Compile("(?i)" + p.regexStr)
			if err != nil {
				return nil, errors.New("Invalid regex: " + p.regexStr)
			}
			data = append(data, p)
			fmt.Fprintf(os.Stderr, "Loading [%s][%s] Regex %s\n", p.tags[0], p.tags[1], p.regexStr)
		}
	}
	return data, nil
}
