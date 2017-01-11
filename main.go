package main

import (
	"time"
	"os"
	"log"
	"sync"
	"math/rand"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

var counter = 0
var mutex = &sync.Mutex{}

func main() {
	log.Printf("ETCD_ENDPOINT: %v", os.Getenv("ETCD_ENDPOINT"))
	err := run()
	if (err != nil) {
		log.Fatal(err)
	}
}

func run() (error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{os.Getenv("ETCD_ENDPOINT")},
		DialTimeout: 5 * time.Second,
	})
	if (err != nil) {
		return err
	}
	defer client.Close()

	for i := 0; i < 5; i++ {
		go work(client)
	}

	c := make(chan int)
	<-c

	return nil
}

func work(client *clientv3.Client) {
	mutex.Lock()
	workerNumber := counter
	counter++
	mutex.Unlock()

	session, err := concurrency.NewSession(client)
	if (err != nil) {
		log.Fatal(err)
	}

	locker := concurrency.NewLocker(session, "a")
	for {
		log.Printf("worker %v : is waiting for lock", workerNumber)
		locker.Lock()
		log.Printf("worker %v : acquired lock", workerNumber)
		waitCount := rand.Intn(10)
		log.Printf("worker %v : wait for %v * 500 Milliseconds", workerNumber, waitCount)
		for i := 0; i <= waitCount; i++ {
			time.Sleep(500 * time.Millisecond)
			log.Printf("worker %v : is just waiting", workerNumber, waitCount)
		}
		locker.Unlock()
	}
}
