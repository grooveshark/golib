package agg

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"text/tabwriter"
	"time"
)

type aggMsg struct {
	name string
	n    float64
}

type printMsg struct {
	div   float64
	retCh chan bool
}

type stats struct {
	start, end time.Time
	ls         []float64
}

var inCh = make(chan *aggMsg)
var printCh = make(chan *printMsg)

func init() {
	go spin()
}

func sorted(lsa []float64) []float64 {
	lsb := make([]float64, len(lsa))
	copy(lsb, lsa)
	sort.Float64s(lsb)
	return lsb
}

func median(ls []float64) float64 {
	return ls[len(ls)/2]
}

func average(ls []float64) float64 {
	var tot float64
	for i := range ls {
		tot += ls[i]
	}
	return tot / float64(len(ls))
}

func genStats(s *stats, div float64) (elapsed, rate, min, max, med, avg float64) {
	lss := sorted(s.ls)
	min = lss[0] / div
	max = lss[len(lss)-1] / div
	avg = average(lss) / div
	med = median(lss) / div
	elapsed = s.end.Sub(s.start).Seconds()
	rate = float64(len(s.ls)) / elapsed
	return
}

func spin() {
	m := map[string]*stats{}
	for {
		select {
		case msg := <-inCh:
			s, ok := m[msg.name]
			if !ok {
				s = &stats{
					ls:    make([]float64, 0, 1024),
					start: time.Now(),
				}
				m[msg.name] = s
			}
			s.ls = append(s.ls, msg.n)
		case msg := <-printCh:
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
			fmt.Println("\n--- aggregator stats ---\n")
			fmt.Fprintf(w, "\ttotal events\telapsed (s)\trate (events/s)\tmedian\tavg\tmin/max\n")
			fmt.Fprintf(w, "\t---\t---\t---\t---\t---\t---\n")
			for n, s := range m {
				s.end = time.Now()
				elapsed, rate, min, max, med, avg := genStats(s, msg.div)
				fmt.Fprintf(
					w,
					"%s\t| %d\t%f\t%f\t%f\t%f\t%f/%f\n",
					n, len(s.ls), elapsed, rate, med, avg, min, max,
				)
			}
			w.Flush()
			msg.retCh <- true
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
	msg := printMsg{div, make(chan bool)}
	printCh <- &msg
	<-msg.retCh
}

// Creates a signal interrupt so that upon a Ctrl-C (as well as some others)
// Print(div) will be called and then the process will be exited
func CreateInterrupt(div float64) {
	go func() {
		log.Println("Waiting for signal")
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		<-c
		go func() {
			<-c
			os.Exit(1)
		}()
		log.Println("Got SIG")
		Print(div)
		os.Exit(0)
	}()
}
