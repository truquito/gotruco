package deco

import (
	"fmt"
	"strings"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/pdt"
)

// Razon2str retorna el string correspondiente a `r`
func Razon2str(r string) string {
	var str string
	switch enco.Razon(r) {
	case enco.EnvidoGanado:
		str = "el envido ganado"
	case enco.RealEnvidoGanado:
		str = "el real envido ganado"
	case enco.FaltaEnvidoGanado:
		str = "la falta envido ganada"
	case enco.EnviteNoQuerido:
		str = "el envido no querido"
	case enco.FlorAchicada:
		str = "la flor achicada"
	case enco.LaUnicaFlor:
		str = "ser la unica flor"
	case enco.LasFlores:
		str = "todas las flores"
	case enco.LaFlorMasAlta:
		str = "tener la flor mas alta"
	case enco.ContraFlorGanada:
		str = "la contra-flor ganada"
	case enco.ContraFlorAlRestoGanada:
		str = "la contra-flor al resto ganada"
	case enco.TrucoNoQuerido:
		str = "el truco no querido"
	case enco.TrucoQuerido:
		str = "el truco ganado"
	case enco.SeFueronAlMazo:
		str = "se fueron al mazo"
	}
	return str
}

// Stringify parsea un pkt
// de momento solo su contenido (el msg)
func Stringify(pkt *enco.Envelope, p *pdt.Partida) string {
	s := Parse(p, pkt.Message)
	return strings.Replace(s, `"`, `'`, -1)
}

func Parse(p *pdt.Partida, msg enco.IMessage) string {

	var decoded string

	switch enco.CodMsg(msg.Cod()) {

	// (string)
	case enco.TError:
		m, _ := msg.(enco.Error)
		s := string(m)
		lower := strings.ToLower(s[:1]) + s[1:]
		decoded = fmt.Sprintf("Error, %s", lower)

	case enco.TLaManoResultaParda:
		decoded = `La Mano resulta parda`

	case enco.TMazo:
		m, _ := msg.(enco.Mazo)
		decoded = fmt.Sprintf(`%s se fue al mazo`, m)

	case enco.TByeBye:
		var template string
		var s string
		m, _ := msg.(enco.ByeBye)

		if p.EsManoAMano() {
			template = "el ganador fue %s"
			s = string(m)
		} else {
			template = "gano el equipo %s"
			s = p.Manojo(string(m)).Jugador.Equipo.String()
		}

		decoded = fmt.Sprintf("C'est fini! "+template, s)

	case enco.TDiceSonBuenas:
		m, _ := msg.(enco.DiceSonBuenas)
		decoded = fmt.Sprintf(`%s: "son buenas"`, m)

	case enco.TAbandono:
		m, _ := msg.(enco.Abandono)
		autor := p.Manojo(string(m))
		decoded = fmt.Sprintf(`%s abandono la partida. Gano el equipo %s`,
			autor.Jugador.ID, autor.Jugador.GetEquipoContrario())

	case enco.TCantarFlor:
		m, _ := msg.(enco.CantarFlor)
		decoded = fmt.Sprintf(`%s canta flor`, string(m))

	case enco.TCantarContraFlor:
		m, _ := msg.(enco.CantarContraFlor)
		decoded = fmt.Sprintf(`%s canta contra-flor`, string(m))

	case enco.TCantarContraFlorAlResto:
		m, _ := msg.(enco.CantarContraFlorAlResto)
		decoded = fmt.Sprintf(`%s canta contra-flor al resto`, string(m))

	case enco.TTocarEnvido:
		m, _ := msg.(enco.TocarEnvido)
		decoded = fmt.Sprintf(`%s toca envido`, string(m))

	case enco.TTocarRealEnvido:
		m, _ := msg.(enco.TocarRealEnvido)
		decoded = fmt.Sprintf(`%s toca real envido`, string(m))

	case enco.TTocarFaltaEnvido:
		m, _ := msg.(enco.TocarFaltaEnvido)
		decoded = fmt.Sprintf(`%s toca falta envido`, string(m))

	case enco.TElEnvidoEstaPrimero:
		m, _ := msg.(enco.ElEnvidoEstaPrimero)
		decoded = fmt.Sprintf(`%s "el envido esta primero!"`, string(m))

	case enco.TGritarTruco:
		m, _ := msg.(enco.GritarTruco)
		decoded = fmt.Sprintf(`%s grita truco`, string(m))

	case enco.TGritarReTruco:
		m, _ := msg.(enco.GritarReTruco)
		decoded = fmt.Sprintf(`%s grita re-truco`, string(m))

	case enco.TGritarVale4:
		m, _ := msg.(enco.GritarVale4)
		decoded = fmt.Sprintf(`%s grita vale-4`, string(m))

	case enco.TNoQuiero:
		m, _ := msg.(enco.NoQuiero)
		decoded = fmt.Sprintf(`%s: "no quiero"`, string(m))

	case enco.TConFlorMeAchico:
		m, _ := msg.(enco.ConFlorMeAchico)
		decoded = fmt.Sprintf(`%s: "con flor me achico"`, string(m))

	case enco.TQuieroTruco:
		m, _ := msg.(enco.QuieroTruco)
		decoded = fmt.Sprintf(`%s: "quiero"`, string(m))

	case enco.TQuieroEnvite:
		m, _ := msg.(enco.QuieroEnvite)
		decoded = fmt.Sprintf(`%s: "quiero"`, string(m))

	// (int)
	case enco.TSigTurno:
		decoded = ""

	case enco.TSigTurnoPosMano:
		decoded = ""

	// (string, int)
	case enco.TDiceTengo:
		m, _ := msg.(enco.DiceTengo)
		decoded = fmt.Sprintf(`%s: "tengo %d"`, m.Autor, m.Valor)

	case enco.TManoGanada:
		m, _ := msg.(enco.ManoGanada)
		decoded = fmt.Sprintf(`La %s mano la gano el equipo %s gracias a %s`,
			pdt.NumMano(m.Valor).String(), p.Manojo(m.Autor).Jugador.Equipo, m.Autor)

	case enco.TRondaGanada:
		m, _ := msg.(enco.RondaGanada)
		decoded = fmt.Sprintf(`La ronda ha sido ganada por el equipo %s debido al %s`,
			p.Manojo(m.Autor).Jugador.Equipo, enco.Razon(m.Razon).String())

	case enco.TDiceSonMejores:
		m, _ := msg.(enco.DiceSonMejores)
		decoded = fmt.Sprintf(`%s: "%d son mejores!"`, m.Autor, m.Valor)

	// (partida)
	case enco.TNuevaPartida:
		decoded = ""

	case enco.TNuevaRonda:
		decoded = "Empieza nueva ronda"

	// (string, palo, valor)
	case enco.TTirarCarta:
		// autor, palo, valor := Tipo2(p, m)
		// fmt.Printf("detectado codigo %d\n", enco.CodMsg(m.Cod))
		// decoded = fmt.Sprintf(`%s tira %d de %s`, autor.Jugador.ID, valor, palo.String())
		decoded = ""

	// (string, string, int)
	case enco.TSumaPts:
		m, _ := msg.(enco.SumaPts)
		if p.EsManoAMano() {
			decoded = fmt.Sprintf(`+%d pts para %s por %s`,
				m.Puntos, m.Autor, Razon2str(string(m.Razon)))
		} else {
			decoded = fmt.Sprintf(`+%d pts para el equipo %s por %s`,
				m.Puntos, p.Manojo(m.Autor).Jugador.Equipo.String(), Razon2str(string(m.Razon)))
		}

	}

	return decoded
}
