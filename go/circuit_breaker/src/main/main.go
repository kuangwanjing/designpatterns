package main

import (
	cb "circuitbreaker"
	"fmt"
	"sync"
	"time"
	//"httpclient"
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

	//testAlwaysClosed()
	testCloseToOpenToHalfToClose()
}

func testAlwaysClosed() {
	breaker := cb.NewBreaker(500*time.Millisecond, 0.3, time.Second, 100)
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
}

func testCloseToOpenToHalfToClose() {
	breaker := cb.NewBreaker(500*time.Millisecond, 0.3, 9*time.Second, 5)
	breaker.Run()

	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func(index int) {
			makeRequest := breaker.IsAvailable()
			if makeRequest {
				if index < 700 {
					breaker.MakeSuccess()
				} else {
					breaker.MakeFailure()
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	time.Sleep(time.Second)
	fmt.Println(breaker.IsAvailable())
	time.Sleep(10 * time.Second)
	for i := 0; i < 10; i++ {
		fmt.Println(breaker.IsAvailable())
	}
	breaker.MakeSuccess()
	time.Sleep(time.Second)
	fmt.Println(breaker.IsAvailable())
}
