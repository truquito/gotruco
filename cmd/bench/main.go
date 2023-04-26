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
	// worker
	n  = flag.Int("n", 2, "a string")
	rt = flag.Int("runtime", 10, "total runtime in secs.")
)

func init() {
	flag.Parse()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(*p)
	start := time.Now()
	c := make(chan int, *p)
	cmd := fmt.Sprintf("go run cmd/bench/self-play/*.go -n=%d -runtime=%d", *n, *rt)

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

	fmt.Println("total", sum, time.Since(start).Round(time.Second), "procs", *rt)
}
