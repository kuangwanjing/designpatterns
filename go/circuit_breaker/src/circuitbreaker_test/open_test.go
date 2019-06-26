package circuit_breaker_test

import (
	cb "circuitbreaker"
	"sync"
	"testing"
	"time"
)

func TestTurnIntoHalfOpen(t *testing.T) {
	breaker := cb.NewBreaker(500*time.Millisecond, 0.1, 3*time.Second, 2)
	breaker.Run()

	N := 100
	var wg sync.WaitGroup
	wg.Add(N)

	// after that the breaker should be open.
	for i := 0; i < N; i++ {
		go func() {
			makeRequest := breaker.IsAvailable()
			if makeRequest {
				time.Sleep(300 * time.Millisecond)
				breaker.MakeFailure()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	time.Sleep(3500 * time.Millisecond)

	allTrue := true
	allFalse := true

	for i := 0; i < 50; i++ {
		if breaker.IsAvailable() {
			allTrue = false
		} else {
			allFalse = false
		}
	}

	if allTrue && !allFalse || !allTrue && allFalse {
		t.Errorf("breaker is not half-open")
	}

	breaker.MakeSuccess()

	time.Sleep(time.Second)

	if !breaker.IsAvailable() {
		t.Errorf("breaker should have turned back to closed")
	}
}
