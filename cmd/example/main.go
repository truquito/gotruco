package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/truquito/truco"
	"github.com/truquito/truco/deco"
	"github.com/truquito/truco/enco"
)

var ioCh chan string = make(chan string, 1)

func handleIO() {
	reader := bufio.NewReader(os.Stdin)
	readLn := func(prefix string) string {
		fmt.Print(prefix)
		cmd, _ := reader.ReadString('\n')
		return strings.TrimSuffix(cmd, "\n")
	}

	for {
		cmd := readLn("")
		ioCh <- cmd
	}
}

func main() {

	logfilePath := flag.String("log", "/tmp/truco_logs/", "Logs directory path")
	seconds := flag.Int("timeout", 60, "Timeout per turn (in seconds)")
	n := flag.Int("n", 2, "Number of players (2, 4 or 6)")
	flag.Parse()

	os.MkdirAll(*logfilePath, os.ModePerm)
	logfile := newLogFile(*logfilePath)

	azules := []string{"Alice", "Ariana", "Annie"}
	rojos := []string{"Bob", "Ben", "Bill"}
	timeout := time.Duration(*seconds) * time.Second
	p, _ := truco.NuevoJuego(20, azules[:*n>>1], rojos[:*n>>1], 4, true, timeout)

	pJSON, _ := p.MarshalJSON()
	logfile.Write(string(pJSON))

	fmt.Println(p)

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

	// hago una gorutine (y channel para avisar) para el io
	go handleIO()

	for {
		select {
		// canal de entrada del usuario
		case cmd := <-ioCh:
			if cmd == "dump" {
				data, _ := json.Marshal(p)
				fmt.Println(string(data))
			} else {
				logfile.Write(cmd)
				err := p.Cmd(cmd)
				if err != nil {
					fmt.Println("<< " + err.Error())
				}
				for _, pkt := range p.Consumir() {
					fmt.Println(deco.Stringify(&pkt, p.Partida))
				}
				fmt.Println(p)
			}
		// canal de error detectado por parte del simulador
		case <-p.ErrCh:
			// el error deberia estar aca
			for _, pkt := range p.Consumir() {
				fmt.Println(pkt.Message.Cod(), deco.Stringify(&pkt, p.Partida))
			}
			// de momento, el unico error posible
			if p.Expirado() {
				m, _ := p.Err.Message.(enco.TimeOut)
				fmt.Printf("el juego terminó debido a que `%s` no realizó niguna jugada en %s.\n", m, p.DurTurno)
			}
			// fmt.Printf(">> ")
		}

		if p.Terminado() {
			return
			// si es modo bucle, entonces que no salga del for sino que
			// cree un juego nuevo
		}
	}

}
