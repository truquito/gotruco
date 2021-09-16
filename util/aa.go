package util

import "github.com/filevich/truco/pdt"

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

// Retorna todas las acciones posibles para un jugador `m` dado
func GetA(p *pdt.PartidaDT, m *pdt.Manojo) A {

	var A [15]bool

	// tirada de cartas
	for i, c := range m.Cartas {
		j := pdt.TirarCarta{Manojo: m, Carta: *c}
		_, ok := j.Ok(p)
		A[i] = ok
		// msg := enco.Msg(enco.TirarCarta, m.Jugador.ID, int(j.Carta.Palo), j.Carta.Valor)
		// A = append(A, msg)
	}

	// ijugada debe tener metodo ToCod
	js := []pdt.IJugada{
		// TirarCarta{},

		// envite
		pdt.TocarEnvido{Manojo: m},
		pdt.TocarRealEnvido{Manojo: m},
		pdt.TocarFaltaEnvido{Manojo: m},
		pdt.CantarFlor{Manojo: m},
		pdt.CantarContraFlor{Manojo: m},
		pdt.CantarContraFlorAlResto{Manojo: m},
		// { CantarConFlorMeAchico{Manojo: m}, enco.new },

		// truco
		pdt.GritarTruco{Manojo: m},
		pdt.GritarReTruco{Manojo: m},
		pdt.GritarVale4{Manojo: m},

		// respuestas
		pdt.ResponderQuiero{Manojo: m},
		pdt.ResponderNoQuiero{Manojo: m},

		// mazo
		pdt.IrseAlMazo{Manojo: m},
	}

	for i, j := range js {
		_, ok := j.Ok(p)
		A[i+3] = ok
	}

	return A
}

// Retorna TODAS las jugadas posibles de cada jugador
func GetAA(p *pdt.PartidaDT) []A {
	res := make([]A, p.CantJugadores)
	for i := range p.Ronda.Manojos {
		res[i] = GetA(p, &p.Ronda.Manojos[i])
	}
	return res
}