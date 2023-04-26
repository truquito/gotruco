package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/filevich/truco/pdt"
)

var (
	// flags
	n  = flag.Int("n", 2, "a string")
	rt = flag.Int("runtime", 3, "total runtime in secs.")
)

func init() {
	flag.Parse()
}

func worker(

	totalRunningTime time.Duration,

) {
	var (
		azules = []string{"Alice", "Ariana", "Annie"}
		rojos  = []string{"Bob", "Ben", "Bill"}
	)

	var (
		last_snapshot []byte   = nil
		actions       []string = nil
	)

	start := time.Now()
	t := 0

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(string(last_snapshot))
			fmt.Println(actions)
			fmt.Println("Recovered in f", r)
		}
	}()

	for time.Since(start) < totalRunningTime {
		p, _ := pdt.NuevaPartida(20, azules[:*n>>1], rojos[:*n>>1])
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

	fmt.Print(t)
}

func main() {
	totalRunTime := time.Second * time.Duration(*rt)
	worker(totalRunTime)
}
