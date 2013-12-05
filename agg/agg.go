package agg

import (
	"fmt"
	"sort"
	"text/tabwriter"
	"os"
)

type aggMsg struct {
	name string
	n float64
}

var inCh = make(chan *aggMsg)
var printCh = make(chan float64)

func init() {
	go spin()
}

func sorted(lsa []float64) []float64 {
	lsb := make([]float64, len(lsa))
	copy(lsb,lsa)
	sort.Float64s(lsb)
	return lsb
}

func median(ls []float64) float64 {
	return ls[len(ls) / 2]
}

func average(ls []float64) float64 {
	var tot float64
	for i := range ls {
		tot += ls[i]
	}
	return tot / float64(len(ls))
}

func stats(ls []float64, div float64) (min, max, med, avg float64) {
	lss := sorted(ls)
	min = lss[0] / div
	max = lss[len(lss)-1] / div
	avg = average(lss) / div
	med = median(lss) / div
	return
}

func spin() {
	m := map[string][]float64{}
	for {
		select {
		case msg := <- inCh:
			if _, ok := m[msg.name]; !ok {
				m[msg.name] = make([]float64, 0, 1024)
			}
			m[msg.name] = append(m[msg.name], msg.n)
		case div := <- printCh:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
			fmt.Println("--- aggregator stats ---")
			for n, ls := range m {
				min, max, med, avg := stats(ls, div)
				fmt.Fprintf(
					w,
					"%s\ttotal events: %d\tmedian: %f\tavg: %f\tmin/max: %f/%f\n",
					n, len(ls), med, avg, min, max,
				)
			}
			w.Flush()
		}
	}

}

// Agg sends the given value and adds it to the statistics for the given name
func Agg(name string, n float64) {
	inCh <- &aggMsg{name, n}
}

// Prints the current aggregation stats to stdout, dividing each by the given
// float. The dividing is so you can change the units that your statistics are
// being shown in, put in 1 if you want them as they were aggregated. Use 0 if
// you want your program to panic.
func Print(div float64) {
	printCh <- div
}
