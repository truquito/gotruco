package pdt

import (
	"math/rand"
	"strconv"

	"github.com/filevich/truco/enco"
)

/*

NUEVO SISTEMA:

*/

func IsDone(pkts []enco.Packet) bool {
	for _, pkt := range pkts {
		if pkt.Message.Cod() == enco.TNuevaPartida ||
			pkt.Message.Cod() == enco.TNuevaRonda ||
			pkt.Message.Cod() == enco.TRondaGanada {
			return true
		}
	}
	return false
}

func Random_action_chi(chi []IJugada) (raix int) {
	return rand.Intn(len(chi))
}

func Random_action_chis(chis [][]IJugada) (rmix, raix int) {
	// hago un cambio de variable:
	// tomo en cuenta solo aquellos chi's que tengan al menos una accion habilitada
	// lo almaceno como un slice de mix's
	habilitados := make([]int, 0, len(chis))
	for mix, chi := range chis {
		if len(chi) > 0 {
			habilitados = append(habilitados, mix)
		}
	}

	r_habilitado_ix := rand.Intn(len(habilitados))
	rmix = habilitados[r_habilitado_ix]
	raix = rand.Intn(len(chis[rmix]))

	return rmix, raix
}

// Retorna todas las acciones posibles para un jugador `m` dado
func Chi(p *Partida, m *Manojo) []IJugada {
	return MetaChi(p, m, true)
}

// Retorna TODAS las jugadas posibles de cada jugador
func Chis(p *Partida) [][]IJugada {
	return MetaChis(p, true)
}

// Retorna todas las acciones posibles para un jugador `m` dado
func MetaChi(p *Partida, m *Manojo, allowMazo bool) []IJugada {

	chi := make([]IJugada, 0, 15)

	// tirada de cartas
	for _, c := range m.Cartas {
		j := TirarCarta{JID: m.Jugador.ID, Carta: *c}
		_, ok := j.Ok(p)
		if ok {
			chi = append(chi, j)
		}
	}

	// ijugada debe tener metodo ToCod
	js := []IJugada{
		// TirarCarta{},

		// envite
		TocarEnvido{JID: m.Jugador.ID},
		TocarRealEnvido{JID: m.Jugador.ID},
		TocarFaltaEnvido{JID: m.Jugador.ID},
		CantarFlor{JID: m.Jugador.ID},
		CantarContraFlor{JID: m.Jugador.ID},
		CantarContraFlorAlResto{JID: m.Jugador.ID},
		// { CantarConFlorMeAchico{JID: m.Jugador.ID}, enco.new },

		// truco
		GritarTruco{JID: m.Jugador.ID},
		GritarReTruco{JID: m.Jugador.ID},
		GritarVale4{JID: m.Jugador.ID},

		// respuestas
		ResponderQuiero{JID: m.Jugador.ID},
		ResponderNoQuiero{JID: m.Jugador.ID},

		// mazo
		IrseAlMazo{JID: m.Jugador.ID},
	}

	if !allowMazo {
		js = js[:len(js)-1]
	}

	for _, j := range js {
		_, ok := j.Ok(p)
		if ok {
			chi = append(chi, j)
		}
	}

	return chi
}

// Retorna TODAS las jugadas posibles de cada jugador
func MetaChis(p *Partida, allowMazo bool) [][]IJugada {
	n := len(p.Ronda.Manojos)
	res := make([][]IJugada, n)
	for i := range p.Ronda.Manojos {
		res[i] = MetaChi(p, &p.Ronda.Manojos[i], allowMazo)
	}
	return res
}

/*

DEPRECADO:

*/

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

	// Con flor me achico ~ no quiero
	// Con flor quiero ~ quiero

Respuestas
	Quiero
	No quiero
	Mazo

*/

// 3 tiradas + 12 jugadas = 15 acciones
type A [15]bool // por default arranca en `false` todos

// no tiene sentido que sea un metodo de A
func (a A) ToJugada(p *Partida, mix, aix int) IJugada {
	return ToJugada(p, mix, aix)
}

func (a A) String() string {
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

	for i, v := range a {
		if v {
			s += codigos[i] + ", "
		}
	}

	if len(s) > 0 {
		s = s[:len(s)-2]
	}

	return s
}

func ActionToString(a A, aix int, jix int, p *Partida) string {
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

	if aix <= 2 {
		// este string no lo agarra la regex de p.Cmd(.)
		// return p.Ronda.Manojos[jix].Cartas[aix].String()
		c := p.Ronda.Manojos[jix].Cartas[aix]
		return strconv.Itoa(c.Valor) + " " + c.Palo.String()
	}

	return codigos[aix]
}

// Retorna todas las acciones posibles para un jugador `m` dado
func GetA(p *Partida, m *Manojo) A {

	var A [15]bool

	// tirada de cartas
	for i, c := range m.Cartas {
		j := TirarCarta{JID: m.Jugador.ID, Carta: *c}
		_, ok := j.Ok(p)
		A[i] = ok
		// msg := enco.Msg(enco.TirarCarta, m.Jugador.ID, int(j.Carta.Palo), j.Carta.Valor)
		// A = append(A, msg)
	}

	// ijugada debe tener metodo ToCod
	js := []IJugada{
		// TirarCarta{},

		// envite
		TocarEnvido{JID: m.Jugador.ID},
		TocarRealEnvido{JID: m.Jugador.ID},
		TocarFaltaEnvido{JID: m.Jugador.ID},
		CantarFlor{JID: m.Jugador.ID},
		CantarContraFlor{JID: m.Jugador.ID},
		CantarContraFlorAlResto{JID: m.Jugador.ID},
		// { CantarConFlorMeAchico{JID: m.Jugador.ID}, enco.new },

		// truco
		GritarTruco{JID: m.Jugador.ID},
		GritarReTruco{JID: m.Jugador.ID},
		GritarVale4{JID: m.Jugador.ID},

		// respuestas
		ResponderQuiero{JID: m.Jugador.ID},
		ResponderNoQuiero{JID: m.Jugador.ID},

		// mazo
		IrseAlMazo{JID: m.Jugador.ID},
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

func ToJugada(p *Partida, mix, aix int) IJugada {
	m := &p.Ronda.Manojos[mix]

	if esCarta := aix < 3; esCarta {
		return TirarCarta{JID: m.Jugador.ID, Carta: *m.Cartas[aix]}
	}

	var jugada IJugada

	switch aix {
	case 3: // envido
		jugada = TocarEnvido{JID: m.Jugador.ID}
	case 4: // real-envido
		jugada = TocarRealEnvido{JID: m.Jugador.ID}
	case 5: // falta-envido
		jugada = TocarFaltaEnvido{JID: m.Jugador.ID}
	case 6: // flor
		jugada = CantarFlor{JID: m.Jugador.ID}
	case 7: // contra-flor
		jugada = CantarContraFlor{JID: m.Jugador.ID}
	case 8: // contra-flor-al-resto
		jugada = CantarContraFlorAlResto{JID: m.Jugador.ID}
	case 9: // truco
		jugada = GritarTruco{JID: m.Jugador.ID}
	case 10: // re-truco
		jugada = GritarReTruco{JID: m.Jugador.ID}
	case 11: // vale-4
		jugada = GritarVale4{JID: m.Jugador.ID}
	case 12: // quiero
		jugada = ResponderQuiero{JID: m.Jugador.ID}
	case 13: // no-Quiero
		jugada = ResponderNoQuiero{JID: m.Jugador.ID}
	case 14: // mazo
		jugada = IrseAlMazo{JID: m.Jugador.ID}
	}

	return jugada
}
