package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jpfilevich/truco"
	"github.com/jpfilevich/truco/deco"
	"github.com/jpfilevich/truco/out"
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

// Print imprime los mensajes
func Print(p *truco.Partida) out.Consumer {
	return func(m *out.Packet) {
		if s := deco.Parse(p, m.Message); s != "" {
			fmt.Println(s)
		}
	}
}

func main() {

	logfile := newLogFile("/home/juan/Workspace/_tmp/truco_logs/")

	// p, _ := truco.NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})
	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	pJSON, _ := p.MarshalJSON()
	logfile.Write(string(pJSON))

	p.Print()
	out.Consume(p.Stdout, out.Print)

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
			// consumo el channel de output
			out.Consume(p.Stdout, Print(p))
			p.Print()
		case <-p.ErrCh:
			out.Consume(p.Stdout, Print(p))
			fmt.Printf(">> ")
		}

		if p.Terminada() {
			break
		}
	}

}
