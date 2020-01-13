package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jpfilevich/truco"
)

var reader = bufio.NewReader(os.Stdin)

func readLn() string {
	cmd, _ := reader.ReadString('\n')
	return strings.TrimSuffix(cmd, "\n")
}

func main() {
	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	p.Print()

	for p.NoAcabada() {
		fmt.Printf("\n>> ")
		cmd := readLn()
		res := p.SetSigJugada(cmd)
		p.Esperar()
		p.Print()
		if res != nil {
			fmt.Println(res)
		}
	}
}
