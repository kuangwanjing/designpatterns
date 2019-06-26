package circuit_breaker_test

import (
	cb "circuitbreaker"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestMinorFailures(t *testing.T) {

	breaker := cb.NewBreaker(500*time.Millisecond, 0.3, time.Second, 100)
	breaker.Run()

	N := 100000
	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func() {
			makeRequest := breaker.IsAvailable()
			if makeRequest {
				sleepTime := rand.Intn(50) + 50
				time.Sleep(time.Duration(sleepTime) * time.Millisecond)
				r := rand.Intn(10)
				// the failure rate is 0.1, smaller than the threshold
				if r == 0 {
					breaker.MakeFailure()
				} else {
					breaker.MakeSuccess()
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	if !breaker.IsAvailable() {
		t.Errorf("breaker should be closed, but open")
	}
}

func TestLargeAmountFailures(t *testing.T) {

	breaker := cb.NewBreaker(500*time.Millisecond, 0.05, 5*time.Second, 100)
	breaker.Run()

	N := 1000
	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func() {
			makeRequest := breaker.IsAvailable()
			if makeRequest {
				sleepTime := rand.Intn(1000) + 100
				time.Sleep(time.Duration(sleepTime) * time.Millisecond)
				r := rand.Intn(10)
				// the failure rate is 0.2, greater than the threshold
				if r == 0 || r == 1 {
					breaker.MakeFailure()
				} else {
					breaker.MakeSuccess()
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	time.Sleep(time.Second)

	if breaker.IsAvailable() {
		t.Errorf("breaker should be open")
	}
}
