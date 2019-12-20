package truco

import (
	"fmt"
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
const debuggingMode = true

func main() {
	p := nuevaPartida(a20, jugadores[:2])

	// mientras no se haya acabado la partida; e.i.,
	for p.getMaxPuntaje() < p.puntuacion.toInt() {
		p.ronda = nuevaRonda(p.cantJugadores)
		p.ronda.Print()
		// checkeo correspondiente unicamente a la `primera` `mano`
		p.ronda.checkFlorDelMano()
		// mientras no se haya acabado la ronda actual; e.i.:

		for p.ronda.enJuego() {
			p.esperandoJugada()
		}
	}
	// termina la primera mano -> ya no es posible el envido
	p.ronda.envido.deshabilitar()
	fmt.Println("BYE BYE")
}
