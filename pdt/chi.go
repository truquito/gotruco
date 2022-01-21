package pdt

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

func (A A) ToJugada(p *Partida, mix, aix int) IJugada {
	esCarta := aix < 3

	m := &p.Ronda.Manojos[mix]

	if esCarta {
		return TirarCarta{Manojo: m, Carta: *m.Cartas[aix]}
	}

	var jugada IJugada

	switch aix {
	case 3: // envido
		jugada = TocarEnvido{Manojo: m}
	case 4: // real-envido
		jugada = TocarRealEnvido{Manojo: m}
	case 5: // falta-envido
		jugada = TocarFaltaEnvido{Manojo: m}
	case 6: // flor
		jugada = CantarFlor{Manojo: m}
	case 7: // contra-flor
		jugada = CantarContraFlor{Manojo: m}
	case 8: // contra-flor-al-resto
		jugada = CantarContraFlorAlResto{Manojo: m}
	case 9: // truco
		jugada = GritarTruco{Manojo: m}
	case 10: // re-truco
		jugada = GritarReTruco{Manojo: m}
	case 11: // vale-4
		jugada = GritarVale4{Manojo: m}
	case 12: // quiero
		jugada = ResponderQuiero{Manojo: m}
	case 13: // no-Quiero
		jugada = ResponderNoQuiero{Manojo: m}
	case 14: // mazo
		jugada = IrseAlMazo{Manojo: m}

	}

	return jugada
}

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
func GetA(p *Partida, m *Manojo) A {

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
		TocarEnvido{Manojo: m},
		TocarRealEnvido{Manojo: m},
		TocarFaltaEnvido{Manojo: m},
		CantarFlor{Manojo: m},
		CantarContraFlor{Manojo: m},
		CantarContraFlorAlResto{Manojo: m},
		// { CantarConFlorMeAchico{Manojo: m}, enco.new },

		// truco
		GritarTruco{Manojo: m},
		GritarReTruco{Manojo: m},
		GritarVale4{Manojo: m},

		// respuestas
		ResponderQuiero{Manojo: m},
		ResponderNoQuiero{Manojo: m},

		// mazo
		IrseAlMazo{Manojo: m},
	}

	for i, j := range js {
		_, ok := j.Ok(p)
		A[i+3] = ok
	}

	return A
}

// Retorna TODAS las jugadas posibles de cada jugador
func GetAA(p *Partida) []A {
	n := len(p.Ronda.Manojos)
	res := make([]A, n)
	for i := range p.Ronda.Manojos {
		res[i] = GetA(p, &p.Ronda.Manojos[i])
	}
	return res
}

// Retorna todas las acciones posibles para un jugador `m` dado
func Chi(p *Partida, m *Manojo) []IJugada {

	chi := make([]IJugada, 0, 15)

	// tirada de cartas
	for _, c := range m.Cartas {
		j := TirarCarta{Manojo: m, Carta: *c}
		_, ok := j.Ok(p)
		if ok {
			chi = append(chi, j)
		}
	}

	// ijugada debe tener metodo ToCod
	js := []IJugada{
		// TirarCarta{},

		// envite
		TocarEnvido{Manojo: m},
		TocarRealEnvido{Manojo: m},
		TocarFaltaEnvido{Manojo: m},
		CantarFlor{Manojo: m},
		CantarContraFlor{Manojo: m},
		CantarContraFlorAlResto{Manojo: m},
		// { CantarConFlorMeAchico{Manojo: m}, enco.new },

		// truco
		GritarTruco{Manojo: m},
		GritarReTruco{Manojo: m},
		GritarVale4{Manojo: m},

		// respuestas
		ResponderQuiero{Manojo: m},
		ResponderNoQuiero{Manojo: m},

		// mazo
		IrseAlMazo{Manojo: m},
	}

	for _, j := range js {
		_, ok := j.Ok(p)
		if ok {
			chi = append(chi, j)
		}
	}

	return chi
}
