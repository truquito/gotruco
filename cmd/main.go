package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/jpfilevich/truco"
)

/**
 *
 * ATENCION:
 * --------
 * NO OLVIDARSE DE SETEAR LA FLAG `debuggingMode` A
 * FALSO CUANDO SE CORRA EL MAIN, PORQUE DE NO SER
 * ASI, EL PROGRAMA QUEDA ATRAPADO EN UN BUCLE
 * INFINITO USANDO EL 100% DEL CPU
 *
 */

// debuggin flag
// const debuggingMode = true

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

	// todo: retorna nil porque juan no existe
	// todo: si el jugador no existe se muere el programa

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
