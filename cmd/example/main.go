package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/filevich/truco"
	"github.com/filevich/truco/deco"
)

var reader = bufio.NewReader(os.Stdin)

func readLn(prefix string) string {
	fmt.Print(prefix)
	cmd, _ := reader.ReadString('\n')
	return strings.TrimSuffix(cmd, "\n")
}

func handleIO() {
	for {
		cmd := readLn("")
		ioCh <- cmd
	}
}

var ioCh chan string = make(chan string, 1)

func main() {

	logfile_path := "/home/jp/Workspace/_tmp/truco_logs/"
	os.MkdirAll(logfile_path, os.ModePerm)
	logfile := newLogFile(logfile_path)

	n := 2 // <-- num. of players
	azules := []string{"Alice", "Ariana", "Annie"}
	rojos := []string{"Bob", "Ben", "Bill"}
	p, _ := truco.NuevoJuego(20, azules[:n>>1], rojos[:n>>1], true)

	pJSON, _ := p.MarshalJSON()
	logfile.Write(string(pJSON))

	fmt.Println(p)

	for _, pkt := range p.Consume() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

	// hago una gorutine (y channel para avisar) para el io
	go handleIO()

	for {
		select {
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
				for _, pkt := range p.Consume() {
					fmt.Println(deco.Stringify(&pkt, p.Partida))
				}
				fmt.Println(p)
			}
		case <-p.ErrCh:
			for _, pkt := range p.Consume() {
				fmt.Println(deco.Stringify(&pkt, p.Partida))
			}
			fmt.Printf(">> ")
		}

		if p.Terminada() {
			break
		}
	}

}
