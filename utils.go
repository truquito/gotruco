package truco

import (
	"fmt"
	"reflect"
)

func chop(str string, l int) string {
	if len(str) <= l {
		return str
	}
	return str[:l]
}

func contains(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func eliminar(manojos []*Manojo, manojo *Manojo) []*Manojo {
	var i int
	// primero encuentro el elemento
	for i = 0; i <= len(manojos); i++ {
		noLoContiene := i == len(manojos)
		if noLoContiene {
			return manojos
		}
		if manojos[i] == manojo {
			break
		}
	}
	manojos[i] = manojos[len(manojos)-1] // Copy last element to index i.
	return manojos[:len(manojos)-1]      // Truncate slice.
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

// obtenerIdx devuelve el ID correspondiente al argumento *Jugador
func obtenerIdx(jugador *Jugador, jugadores []Jugador) (JugadorIdx, error) {
	for i, j := range jugadores {
		if j.Nombre == jugador.Nombre {
			return JugadorIdx(i), nil
		}
	}
	return -1, fmt.Errorf("Jugador `%s` no encontrado", jugador.Nombre)
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
