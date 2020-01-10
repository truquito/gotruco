package main

import (
	"github.com/jpfilevich/truco"
	"github.com/jpfilevich/truco/ilustrador"
)

func main() {

	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andrés"}, []string{"Roro", "Renzo", "Richard"})

	ilustrador.Imprimir(*p)

	p.Print()

	// for p.NoAcabada() {

	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Printf("\n**Ingresar comando: ")
	// 	cmd, _ := reader.ReadString('\n')
	// 	res := p.SetSigJugada(cmd)

	// 	fmt.Println(res)
	// }

}
