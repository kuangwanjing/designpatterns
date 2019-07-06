package main

/*
import (
	"fmt"
	"net"
	"os"
)

func main() {
	ips, err := net.LookupIP("redis.default.svc.cluster.local")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Printf("redis.default.svc.cluster.local. IN A %s\n", ip.String())
	}
}
*/

import (
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func Lock(servers []string, key string, value string, expire time.Duration) bool {
	startTime := time.Now()
	serversCount := len(servers)
	successfulCount := 0
	ch := make(chan bool)
	for i := 0; i < serversCount; i++ {
		go func(index int) {
			client := redis.NewClient(&redis.Options{
				Addr:         servers[index],
				Password:     "", // no password set
				DB:           0,  // use default DB
				WriteTimeout: 100 * time.Millisecond,
			})
			set, err := client.SetNX(key, value, expire).Result()
			if err == nil && set {
				ch <- true
			}
			ch <- false
		}(i)
	}
	for i := 0; i < serversCount; i++ {
		rst := <-ch
		if rst {
			successfulCount += 1
		}
	}
	if successfulCount <= serversCount/2 {
		Unlock(servers, key, value)
		return false
	}
	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	if elapsed >= expire {
		return false
	}
	return true
}

func Unlock(servers []string, key string, value string) {
	serversCount := len(servers)
	ch := make(chan bool)
	for i := 0; i < serversCount; i++ {
		go func(index int) {
			client := redis.NewClient(&redis.Options{
				Addr:         servers[index],
				Password:     "", // no password set
				DB:           0,  // use default DB
				ReadTimeout:  100 * time.Millisecond,
				WriteTimeout: 100 * time.Millisecond,
			})
			v, err := client.Get(key).Result()
			if err == nil && v == value {
				client.Del(key)
			}
			ch <- true
		}(i)
	}
	for i := 0; i < serversCount; i++ {
		<-ch
	}
}

func main() {

	redisServers := []string{
		"localhost:7001",
		"localhost:7002",
		"localhost:7003",
		"localhost:7004",
		"localhost:7005",
	}

	//rst := Lock(redisServers, "lock5", "1", 60*time.Second)
	//println(rst)

	count := 0
	routines := 100
	var wg sync.WaitGroup
	wg.Add(routines)

	lockKey := "inc_7"

	for i := 0; i < routines; i++ {
		go func(index int) {
			val := strconv.Itoa(index)
			var sleepTime time.Duration
			for {
				sleepTime = time.Duration(rand.Intn(500)) * time.Millisecond
				lock := Lock(redisServers, lockKey, val, 60*time.Second)
				println(i, "get lock ?", lock)
				if lock {
					count += 1
					Unlock(redisServers, lockKey, val)
					break
				} else {
					time.Sleep(sleepTime)
					sleepTime *= 2
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	fmt.Println(count)
}
