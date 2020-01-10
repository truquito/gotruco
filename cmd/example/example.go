package main

import (
	"github.com/jpfilevich/truco"
)

func main() {

	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andrés"}, []string{"Roro", "Renzo", "Richard"})
	p.Print()

	// for p.NoAcabada() {

	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Printf("\n**Ingresar comando: ")
	// 	cmd, _ := reader.ReadString('\n')
	// 	res := p.SetSigJugada(cmd)

	// 	fmt.Println(res)
	// }

}
