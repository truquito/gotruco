package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/filevich/truco"
	"github.com/filevich/truco/deco"
	"github.com/filevich/truco/enco"
)

var reader = bufio.NewReader(os.Stdin)

func readLn(prefix string) string {
	fmt.Printf(prefix)
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

	logfile := newLogFile("/home/jp/Workspace/_tmp/truco_logs/")

	p, out, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	pJSON, _ := p.MarshalJSON()
	logfile.Write(string(pJSON))

	fmt.Println(p)
	enco.Consume(out, func(pkt *enco.Packet) {
		fmt.Print(deco.Stringify(pkt, &p.PartidaDT))
	})

	// hago una gorutine (y channel para avisar) para el io
	go handleIO()

	for {
		select {
		case cmd := <-ioCh:
			logfile.Write(cmd)
			err := p.Cmd(cmd)
			if err != nil {
				fmt.Println("<< " + err.Error())
			}
			enco.Consume(out, func(pkt *enco.Packet) {
				fmt.Print(deco.Stringify(pkt, &p.PartidaDT))
			})
			fmt.Println(p)
		case <-p.ErrCh:
			enco.Consume(out, func(pkt *enco.Packet) {
				fmt.Print(deco.Stringify(pkt, &p.PartidaDT))
			})
			fmt.Printf(">> ")
		}

		if p.Terminada() {
			break
		}
	}

}
