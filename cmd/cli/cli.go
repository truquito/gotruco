package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jpfilevich/truco"
)

var reader = bufio.NewReader(os.Stdin)

func readLn(prefix string) string {
	fmt.Printf(prefix)
	cmd, _ := reader.ReadString('\n')
	return strings.TrimSuffix(cmd, "\n")
}

func handleIO() {
	for {
		cmd := readLn(">> ")
		ioCh <- cmd
	}
}

var ioCh chan string = make(chan string, 1)

func main() {

	logfile := newLogFile("/home/juan/Workspace/_tmp/truco_logs/")

	// p, _ := truco.NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})
	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	logfile.Write(string(p.MarshalJSON()))

	p.Print()
	truco.out.Consume(p.Stdout)

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
			// truco.out.Consume(p.Stdout)
			p.Print()
		case <-p.ErrCh:
			truco.out.Consume(p.Stdout)
			fmt.Printf(">> ")
		}

		if p.Terminada() {
			break
		}
	}

}
