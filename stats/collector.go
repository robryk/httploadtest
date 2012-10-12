package stats

import "fmt"
import "io"
import "sync"
import "math"

type ValueCollector interface {
	Collect(float64)
}

type PrintCollector struct {
	Output io.Writer
	Name   string
}

func (pc PrintCollector) Collect(value float64) {
	fmt.Printf("%s: %f\n", pc.Name, value)
}

type StatsCollector struct {
	count int
	sum float64
	sumSquares float64
	bucketCounts []int

	buckets []float64
	mutex sync.Mutex
}

func NewStatsCollector(buckets []float64) *StatsCollector {
	return &StatsCollector{
		bucketCounts: make([]int, len(buckets) + 1),
		buckets: append([]float64{0}, buckets...),
	}
}

func (sc *StatsCollector) Collect(value float64) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.count++
	sc.sum += value
	sc.sumSquares += value*value

	for i, low := range sc.buckets {
		if value > low {
			sc.bucketCounts[i]++
			break
		}
	}
}

type Stats struct {
	Count int
	Avg float64
	StdDev float64
}

func (sc *StatsCollector) GetStats() Stats {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	avg := sc.sum / float64(sc.count)
	avgSquare := sc.sumSquares / float64(sc.count)

	s := Stats{
		Count: sc.count,
		Avg: avg,
		StdDev: math.Sqrt(avg*avg - avgSquare),
	}

	return s
}
