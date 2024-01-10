package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var (
	p = flag.Int("p", 2, "total subprocesses")
	n = flag.Int("n", 2, "player per match (2, 4 or 6)")
	t = flag.Int("t", 100, "total runtime in secs.")
	// v = flag.Bool("v", false, "verbose (or silent) mode")
)

func init() {
	flag.Parse()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(*p)
	start := time.Now()
	c := make(chan int, *p)
	cmd := fmt.Sprintf("go run cmd/bench/self-play/*.go -n %d -t %d", *n, *t)

	for i := 1; i <= *p; i++ {
		go func() {
			defer wg.Done()
			out, err := exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				panic(err)
			}
			res, _ := strconv.Atoi(string(out))
			c <- res
		}()
	}

	wg.Wait()
	close(c)

	sum := 0
	for x := range c {
		sum += x
	}

	fmt.Printf("total: %d time: %s, procs: %d players: %d runtime: %d\n",
		sum,
		time.Since(start).Round(time.Second),
		*p,
		*n,
		*t)
}
