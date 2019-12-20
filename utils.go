package truco

import (
	"fmt"
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func maxOf3(cartas [3]Carta) int {
	max := 0
	for _, carta := range cartas {
		if carta.Valor > max {
			max = int(carta.Valor)			
		}
	}
	return max
}

func max(nums ...int) int {
	max := nums[0]
	for _, x := range nums[1:] {
		if x > max {
			max = x
		}
	}
	return max
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// todo: esto es ineficiente
// decodeJugador devuelve el puntero al jugador, dado un string que los identifique
func decodeJugador(codigoJugador string, jugadores []Jugador) (*Jugador, error) {
	for i := range jugadores {
		if jugadores[i].nombre == codigoJugador {
			return &jugadores[i], nil
		}
	}
	return nil, fmt.Errorf("Jugador `%s` no encontrado", codigoJugador)
}

// obtenerIdx devuelve el ID correspondiente al argumento *Jugador
func obtenerIdx(jugador *Jugador, jugadores []Jugador) (JugadorIdx, error) {
	for i, j := range jugadores {
		if j.nombre == jugador.nombre {
			return JugadorIdx(i), nil
		}
	}
	return -1, fmt.Errorf("Jugador `%s` no encontrado", jugador.nombre)
}

// complemento de 'x in A, where A=[0, 1, 2]' devuelve
// *por separado*, el conjunto ordenado 'A - x'
// (A sin el elemento x)
// PRE: 0 <= x <= 2
func complemento(x int) (p, q int) {
	switch x {
	case 0:
		return 1, 2
	case 1:
		return 0, 2
	default: // x == 2
		return 0, 1	
	}
}

// leGanaDeMano devuelve `true` sii
// `i` "le gana de mano" a `j`
// PARA MAS INFO DETALLADA VER DOCUMENTACION
func leGanaDeMano(i, j, mano JugadorIdx, cantJugadores int) bool {
	// cambios de variables
	p := cv(i, mano, cantJugadores)
	q := cv(j, mano, cantJugadores)
	return p < q
}
// ver documentacion cambio de variable
func cv(x, mano JugadorIdx, cantJugadores int) (y JugadorIdx) {
	if x >= mano {
		y = x - mano
	} else {
		c := JugadorIdx(cantJugadores) - mano
		y = x + c
	}
	return y
}

func print(a []string) {
	for _, s := range a {
		fmt.Print(s)
	}
}