package circuit_breaker_test

import (
	cb "circuitbreaker"
	"sync"
	"testing"
)

func BenchmarkStructChan(b *testing.B) {
	breaker := new(cb.CircuitBreaker)
	breaker.Run()

	var wg sync.WaitGroup
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func() {
			makeRequest := breaker.IsAvailable()
			if makeRequest {
				breaker.MakeSuccess()
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
