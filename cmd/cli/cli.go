package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jpfilevich/truco"
)

var reader = bufio.NewReader(os.Stdin)

func readLn() string {
	fmt.Printf(">> ")
	cmd, _ := reader.ReadString('\n')
	return strings.TrimSuffix(cmd, "\n")
}

func handleIO() {
	for {
		cmd := readLn()
		ioCh <- cmd
	}
}

func read(buff *bytes.Buffer) (*truco.Msg, error) {
	e := new(truco.Msg)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func consume(buff *bytes.Buffer) {
	for {
		e, err := read(buff)
		if err == io.EOF {
			break
		}
		fmt.Println(*e)
	}
}

var ioCh chan string = make(chan string, 1)

func main() {

	logfile := newLogFile("/home/juan/Workspace/_tmp/truco_logs/")

	// p, _ := truco.NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})
	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	logfile.Write(p.ToJSON())

	p.Print()
	consume(p.Stdout)

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
			consume(p.Stdout)
			p.Print()
		case <-p.ErrCh:
			consume(p.Stdout)
			fmt.Printf(">> ")
		}

		if p.Terminada() {
			break
		}
	}

}
