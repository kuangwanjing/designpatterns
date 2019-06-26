package circuitbreaker

import (
	"time"
)

type ClosedStateMachine struct {
	statsTimeSlice time.Duration
	breakThreshold float32
	stats          []int
}

func (this *ClosedStateMachine) Start(stopChan <-chan struct{}, stoppedChan chan struct{},
	statsChan <-chan bool, stateChan chan CircuitBreakerState) {

	defer func() {
		close(stoppedChan)
	}()

	ticker := time.NewTicker(this.statsTimeSlice)
	this.stats[0] = 0 // counter for success
	this.stats[1] = 0 // counter for failure

	for {
		select {
		case <-ticker.C:
			total := this.stats[0] + this.stats[1]
			if total > 0 && float32(this.stats[1])/float32(total) >= this.breakThreshold {
				stateChan <- Open
				// stop the ticker to avoid the new stats window to begin so that another state signal is sent.
				ticker.Stop()
			} else {
				// clear the statistics
				this.stats[0] = 0
				this.stats[1] = 0
			}
		case rst := <-statsChan:
			if rst {
				this.stats[0]++
			} else {
				this.stats[1]++
			}
		case <-stopChan: // this is a stop signal for this routine
			return
		}
	}
}
