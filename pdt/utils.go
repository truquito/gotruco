package pdt

import (
	"reflect"
)

// Contains .
func Contains(slice interface{}, item interface{}) bool {
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

// Eliminar .
func Eliminar(manojos []*Manojo, manojo *Manojo) []*Manojo {
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

func maxOf3(cartas [3]*Carta) int {
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

// A es el conjunto de acciones posibles para el manojo `m`
/*
Gritos
	Truco    // 1/2
	Re-truco // 2/3
	Vale 4   // 3/4

Toques
	Envido
	Real envido
	Falta envido

Cantos
	Flor                 // 2pts (tanto o el-primero)
	Contra flor          // 3 pts
	Contra flor al resto // 4 pts

	// Con flor me achico ~ quiero
	// Con flor quiero ~ no quiero

Respuestas
	Quiero
	No quiero

*/

// 3 tiradas + 12 jugadas = 15 acciones
type A [15]bool // por default arranca en `false` todos

func (A A) String() string {
	s := ""

	codigos := []string{
		// cartas
		"primera",
		"segunda",
		"tercera",
		// envite
		"envido",
		"real-envido",
		"falta-envido",
		"flor",
		"contra-flor",
		"contra-flor-al-resto",
		// truco
		"truco",
		"re-truco",
		"vale-4",
		// respuestas
		"quiero",
		"no-Quiero",
		"mazo",
	}

	for i, v := range A {
		if v {
			s += codigos[i] + ", "
		}
	}

	if len(s) > 0 {
		s = s[:len(s)-2]
	}

	return s
}

func (p *PartidaDT) A(m *Manojo) A {

	var A [15]bool

	// tirada de cartas
	for i, c := range m.Cartas {
		j := TirarCarta{Manojo: m, Carta: *c}
		_, ok := j.Ok(p)
		A[i] = ok
		// msg := enco.Msg(enco.TirarCarta, m.Jugador.ID, int(j.Carta.Palo), j.Carta.Valor)
		// A = append(A, msg)
	}

	// ijugada debe tener metodo ToCod
	js := []IJugada{
		// TirarCarta{},

		// envite
		TocarEnvido{m},
		TocarRealEnvido{m},
		TocarFaltaEnvido{m},
		CantarFlor{m},
		CantarContraFlor{m},
		CantarContraFlorAlResto{m},
		// { CantarConFlorMeAchico{m}, enco.new },

		// truco
		GritarTruco{m},
		GritarReTruco{m},
		GritarVale4{m},

		// respuestas
		ResponderQuiero{m},
		ResponderNoQuiero{m},

		// mazo
		IrseAlMazo{m},
	}

	for i, j := range js {
		_, ok := j.Ok(p)
		A[i+3] = ok
	}

	return A
}
