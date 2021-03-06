package circuitbreaker

import (
	"time"
)

type OpenStateMachine struct {
	sleepTimer time.Duration
}

func (this *OpenStateMachine) Start(stopChan <-chan struct{}, stoppedChan chan struct{},
	statsChan <-chan bool, stateChan chan CircuitBreakerState) {

	defer func() {
		close(stoppedChan)
	}()

	ticker := time.NewTicker(this.sleepTimer)

	for {
		select {
		case <-ticker.C:
			stateChan <- HalfOpen
			// stop the ticker to avoid the new stats window to begin so that another state signal is sent.
			ticker.Stop()
		case <-statsChan:
			continue // do nothing
		case <-stopChan: // this is a stop signal for this routine
			return
		}
	}
}
