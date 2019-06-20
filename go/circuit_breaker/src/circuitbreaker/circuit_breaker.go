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
}

func (this *CircuitBreaker) init() {
	this.state = Closed
	this.closedStateFn = &ClosedStateMachine{time.Second, 0.3, []int{0, 0}}
	this.openStateFn = &OpenStateMachine{500 * time.Millisecond}
	this.halfOpenStateFn = &HalfOpenStateMachine{}
	this.requestChan = make(chan bool, 10)
	this.requestedChan = make(chan bool, 10)
	this.statsChan = make(chan bool)
	this.stateChan = make(chan CircuitBreakerState)
}

func (this *CircuitBreaker) Run() {

	this.init()

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
		// the chance is 1 in 10
		r := rand.Intn(10)
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
