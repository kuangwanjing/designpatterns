package circuitbreaker

import (
	"math/rand"
	"time"
)

type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	state           CircuitBreakerState
	closedStateFn   *ClosedStateMachine
	openStateFn     *OpenStateMachine
	halfOpenStateFn *HalfOpenStateMachine
	requestChan     chan bool
	requestedChan   chan bool
	statsChan       chan bool
	stateChan       chan CircuitBreakerState
	tryRate         int
}

func NewBreaker(statsWindow time.Duration, openThreshold float32, openWindow time.Duration, tryRate int) *CircuitBreaker {
	breaker := new(CircuitBreaker)
	breaker.state = Closed
	breaker.closedStateFn = &ClosedStateMachine{statsWindow, openThreshold, []int{0, 0}}
	breaker.openStateFn = &OpenStateMachine{openWindow}
	breaker.halfOpenStateFn = &HalfOpenStateMachine{}
	breaker.requestChan = make(chan bool, 5)
	breaker.requestedChan = make(chan bool, 5)
	breaker.statsChan = make(chan bool, 5)
	breaker.stateChan = make(chan CircuitBreakerState)
	breaker.tryRate = tryRate
	return breaker
}

func (this *CircuitBreaker) Run() {

	go func() {
		// first run the closed state machine
		stopChan := make(chan struct{})
		stoppedChan := make(chan struct{})
		this.startNewState(Closed, stopChan, stoppedChan)

		for {
			select {
			case <-this.requestChan:
				this.requestedChan <- this.check()
			case newState := <-this.stateChan:
				this.stopCurrentState(stopChan, stoppedChan)
				stopChan = make(chan struct{})
				stoppedChan = make(chan struct{})
				this.startNewState(newState, stopChan, stoppedChan)
				this.state = newState
			}
		}
	}()
}

func (this *CircuitBreaker) MakeSuccess() {
	this.statsChan <- true
}

func (this *CircuitBreaker) MakeFailure() {
	this.statsChan <- false
}

func (this *CircuitBreaker) IsAvailable() bool {
	this.requestChan <- true
	return <-this.requestedChan
}

// check where is circuit is available
func (this *CircuitBreaker) check() bool {
	if this.state == Closed {
		return true
	} else if this.state == HalfOpen {
		r := rand.Intn(this.tryRate)
		if r == 0 {
			return true
		}
	}
	return false
}

// stop the current state machine
func (this *CircuitBreaker) stopCurrentState(stopChan, stoppedChan chan struct{}) {
	// send a stopping signal
	close(stopChan)
	// and receive a stopped signal
	<-stoppedChan
}

func (this *CircuitBreaker) startNewState(state CircuitBreakerState, stopChan, stoppedChan chan struct{}) {
	switch state {
	case Closed:
		go this.closedStateFn.Start(stopChan, stoppedChan, this.statsChan, this.stateChan)
	case Open:
		go this.openStateFn.Start(stopChan, stoppedChan, this.stateChan)
	case HalfOpen:
		go this.halfOpenStateFn.Start(stopChan, stoppedChan, this.statsChan, this.stateChan)
	}
}
