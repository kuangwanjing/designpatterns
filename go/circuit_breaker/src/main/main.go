package main

import (
	cb "circuitbreaker"
	"fmt"
	"sync"
	//"httpclient"
	"time"
)

func main() {
	/*
		client := new(httpclient.HttpClient)
		client.InitClient("kubia.example.com", "35.190.29.215", 80, "", time.Second)
		body, err := client.Get()
		body, err = client.Get()
		body, err = client.Get()

		println(body)
		fmt.Println(err)
	*/

	breaker := new(cb.CircuitBreaker)
	breaker.Run()

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			makeRequest := breaker.IsAvailable()
			if makeRequest {
				breaker.MakeSuccess()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	time.Sleep(time.Second)
	fmt.Println(breaker.IsAvailable())

	fmt.Println(breaker)
}
