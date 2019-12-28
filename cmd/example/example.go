package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jpfilevich/truco"
)

func main() {

	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andrés"}, []string{"Roro", "Renzo", "Richard"})
	p.Ronda.Print()

	for p.NoAcabada() {

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\n**Ingresar comando: ")
		cmd, _ := reader.ReadString('\n')
		res := p.SetSigJugada(cmd)

		fmt.Println(res)
	}

}
