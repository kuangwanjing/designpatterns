package circuitbreaker

type HalfOpenStateMachine struct {
}

func (this *HalfOpenStateMachine) Start(stopChan <-chan struct{}, stoppedChan chan struct{},
	statsChan <-chan bool, stateChan chan CircuitBreakerState) {

	defer func() {
		close(stoppedChan)
	}()

	for {
		select {
		case rst := <-statsChan:
			if rst {
				stateChan <- Closed
			} else {
				stateChan <- Open
			}
		case <-stopChan: // this is a stop signal for this routine
			return
		}
	}
}
