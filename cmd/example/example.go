package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jpfilevich/truco"
)

var reader = bufio.NewReader(os.Stdin)

func main() {
	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andrés"}, []string{"Roro", "Renzo", "Richard"})
	p.Print()

	for p.NoAcabada() {
		fmt.Printf("\n>> ")
		cmd := readLn()
		res := p.SetSigJugada(cmd)
		p.Print()
		if res != nil {
			fmt.Println(res)
		}
	}
}

func readLn() string {
	cmd, _ := reader.ReadString('\n')
	return strings.TrimSuffix(cmd, "\n")
}
