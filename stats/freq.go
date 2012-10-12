package stats

import "time"

type FreqCounter struct {
	sample chan struct{}
	output ValueCollector
}

func (fc FreqCounter) handleCounter() {
	value := 0
	tick := time.Tick(5.0 * time.Second)
	for {
		select {
		case _ = <-fc.sample:
			value++
		case _ = <-tick:
			fc.output.Collect(float64(value) / 5.0)
			value = 0
		}
	}
}

func (fc FreqCounter) Trigger() {
	fc.sample <- struct{}{}
}

func (fc FreqCounter) Start() {
	go fc.handleCounter()
}

func NewFreqCounter(output ValueCollector) FreqCounter {
	return FreqCounter{
		sample: make(chan struct{}),
		output: output,
	}
}
