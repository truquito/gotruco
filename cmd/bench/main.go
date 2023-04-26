package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/filevich/truco/pdt"
)

var (
	n      = 2 // <-- num. of players
	azules = []string{"Alice", "Ariana", "Annie"}
	rojos  = []string{"Bob", "Ben", "Bill"}
)

func worker(

	id int,
	start time.Time,
	totalRunningTime time.Duration,
	c chan<- int,

) {
	// telemetry
	var (
		last_snapshot []byte   = nil
		actions       []string = nil
	)

	t := 0

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(string(last_snapshot))
			fmt.Println(actions)
			fmt.Println("Recovered in f", r)
		}
	}()

	for time.Since(start) < totalRunningTime {
		p, _ := pdt.NuevaPartida(20, azules[:n>>1], rojos[:n>>1])
		last_snapshot, _ = p.MarshalJSON()
		actions = []string{}
		for !p.Terminada() {
			// elijo una al azar
			chis := pdt.MetaChis(p, false)
			rmix, raix := pdt.Random_action_chis(chis)
			a := chis[rmix][raix]
			// la guardo
			actions = append(actions, a.String())
			// la ejecuto
			pkts := a.Hacer(p)
			if pdt.IsDone(pkts) {
				last_snapshot, _ = p.MarshalJSON()
				actions = []string{}
			}
		}
		t++
	}

	c <- t
}

func main() {

	var wg sync.WaitGroup
	t := 3
	wg.Add(t)
	start := time.Now()
	totalRunTime := time.Minute * 10
	c := make(chan int, t)

	for i := 1; i <= t; i++ {
		i := i
		go func() {
			defer wg.Done()
			worker(i, start, totalRunTime, c)
		}()
	}

	wg.Wait()
	close(c)

	sum := 0
	for x := range c {
		sum += x
	}

	// In unbuffered channel writing to channel will not happen until there must
	// be some receiver which is waiting to receive the data
	// no espero a que terminen los goroutines, sino no pueden escribir
	// the waiting for the channel output (ie, "subscribe") must occur BEFORE the
	// write attemps
	// sum := 0
	// for {
	// 	val, ok := <-c
	// 	if ok {
	// 		sum += val
	// 	} else {
	// 		break
	// 	}
	// }
	// close(c)

	fmt.Println("total", sum, time.Since(start).Round(time.Second))
}
