package spendtracker

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)


func (p *PatternDB) matchTags(line string) []string {

	for i := 0; i < len(p.patterns); i++ {
		if p.patterns[i].re.MatchString(line) {
			return (p.patterns[i].tags)
		}
	}
	return nil
}

func LoadAliases(aliasesFile string) ([]*regexp.Regexp, error) {


	file, err := os.Open(aliasesFile)
	var returnData []*regexp.Regexp

	if err != nil {
		return nil, err

	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		re, err := regexp.Compile("(?i)" + line)

		if err != nil {
			return nil, errors.New("Invalid regex in " + aliasesFile + ": " + line)
		}

		returnData = append(returnData, re)

	}
	return returnData, nil
}

func IsMyAccount(aliases []*regexp.Regexp) bool {
	for re := range aliases {


	}

}