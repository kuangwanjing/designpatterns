package circuitbreaker

import (
	"time"
)

type OpenStateMachine struct {
	sleepTimer time.Duration
}

func (this *OpenStateMachine) Start(stopChan <-chan struct{}, stoppedChan chan struct{},
	stateChan chan CircuitBreakerState) {

	defer func() {
		close(stoppedChan)
	}()

	ticker := time.NewTicker(this.sleepTimer)

	for {
		select {
		case <-ticker.C:
			stateChan <- HalfOpen
		case <-stopChan: // this is a stop signal for this routine
			return
		}
	}
}
