package spendtracker

import "fmt"

type Report struct {
	headings []string
	rows     [][]string
}

func (st *SpendTracker) sum(list []int) float32 {
	var val float32
	val = 0
	for _, v := range list {
		val += st.trns[v].Debit
		val += st.trns[v].Credit

	}
	return val
}

func (st *SpendTracker) BuildReport() Report {
	var r Report
	r.headings = st.Periods()

	for _, tag1 := range st.pdb.Level1Tags() {

		for _, tag2 := range st.pdb.Level2Tags(tag1) {

			row := make([]string, len(r.headings)+2)
			row[0] = tag1
			row[1] = tag2

			for coli, date := range r.headings {
				key := IndexTag2Key{
					period:    date,
					level1tag: tag1,
					level2tag: tag2,
				}

				trnlist, ok := st.idxTag2[key]

				val := float32(0.0)
				if ok {
					val = st.sum(trnlist)
				}
				row[coli+2] = fmt.Sprintf("%0.2f", val)
			}
			r.rows = append(r.rows, row)
		}
	}
	return r
}

func (r *Report) PrintReport(format string) string {
	var rep string

	rep = printCSVRow(r.headings)
	for _, line := range r.rows {
		rep = rep + printCSVRow(line)
	}
	return rep
}

func printCSVRow(elems []string) string {
	var ret string
	for i, v := range elems {
		if i == 0 {
			ret = v
		} else {
			ret = ret + "," + v
		}
	}
	ret = ret + "\n"
	return ret
}
