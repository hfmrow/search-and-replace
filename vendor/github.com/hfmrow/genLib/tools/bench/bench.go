// bench.go

/*
  Measuring lapse (may be multiples) between operations.

Usage:
		bench := BenchNew(true) // true mean we want display results at Stop().
		bench.Lapse("first lapse")
		// Doing something ...
		bench.Lapse("Nth lapse")
		// Doing another thing ...
		bench.Stop                 // Display results
		fmt.Println(bench.Average) // Mean time
		fmt.Println(bench.NumTime) // numeric sliced version of results ...
*/

package bench

import (
	"fmt"
	"time"
)

type Bench struct {
	lapse   []time.Time
	label   []string
	totalNs int64
	Results []string
	Average string
	Display bool
	NumTime []numeric
}

func BenchNew(showResults bool) (bench *Bench) {
	bench = new(Bench)
	bench.Display = showResults
	return
}

type numeric struct {
	Min, Sec, Ms, Ns int64
}

func (b *Bench) Lapse(label ...string) {
	b.lapse = append(b.lapse, time.Now())
	if len(label) == 0 {
		label = append(label, fmt.Sprintf("%d", len(b.lapse)))
	}
	b.label = append(b.label, label...)
}

func (b *Bench) Stop() {
	b.Lapse("Total")
	lapseCount := len(b.lapse) - 1
	b.Results = b.Results[:0]
	if lapseCount > 1 {
		for idx := 0; idx < len(b.lapse)-1; idx++ {
			b.calculateLapse(idx, b.lapse[idx], b.lapse[idx+1])
		}
	}
	b.calculateLapse(lapseCount, b.lapse[0], b.lapse[lapseCount])

	m, s, ms, ns := b.getMSmsnano(b.totalNs / int64(len(b.lapse)))
	b.Average = fmt.Sprintf("%v m, %v s, %v ms, %v ns", m, s, ms, ns)

	b.totalNs = 0
	b.label = b.label[:0]
	b.lapse = b.lapse[:0]
}

func (b *Bench) getMSmsnano(diff int64) (min, sec, ms, ns int64) {
	min = (diff / 1000000000) / 60
	sec = diff/1000000000 - (min * 60)
	ms = (diff / 1000000) - (sec * 1000)
	ns = diff - ((diff / 1000000) * 1000000)
	return min, sec, ms, ns
}

func (b *Bench) calculateLapse(count int, start, stop time.Time) {
	diff := stop.Sub(start).Nanoseconds()
	m, s, ms, ns := b.getMSmsnano(diff)
	b.NumTime = append(b.NumTime, numeric{Min: m, Sec: s, Ms: ms, Ns: ns})
	b.Results = append(b.Results, fmt.Sprintf("%s: %v m, %v s, %v ms, %v ns",
		b.label[count], m, s, ms, ns))
	b.totalNs += diff
	if b.Display {
		fmt.Println(b.Results[len(b.Results)-1])
	}
}
