package deco

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/pdt"
)

// Str .
func Str(m enco.Message) string {
	var str string
	json.Unmarshal(m.Cont, &str)
	return str
}

// Int .
func Int(m enco.Message) int {
	var num int
	json.Unmarshal(m.Cont, &num)
	return num
}

// Autor .
func Autor(p *pdt.Partida, m enco.Message) *pdt.Manojo {
	var jid string
	json.Unmarshal(m.Cont, &jid)

	return p.Manojo(jid)
}

// Tipo1 .
func Tipo1(p *pdt.Partida, m enco.Message) (*pdt.Manojo, int) {
	var t1 enco.Tipo1
	json.Unmarshal(m.Cont, &t1)
	return p.Manojo(t1.Autor), t1.Valor
}

// Tipo2 .
func Tipo2(p *pdt.Partida, m enco.Message) (*pdt.Manojo, pdt.Palo, int) {
	var t2 enco.Tipo2
	json.Unmarshal(m.Cont, &t2)
	palo := pdt.ToPalo[t2.Palo]
	return p.Manojo(t2.Autor), palo, t2.Valor
}

// Tipo3 .
func Tipo3(p *pdt.Partida, m enco.Message) (*pdt.Manojo, enco.Razon, int) {
	var t3 enco.Tipo3
	json.Unmarshal(m.Cont, &t3)
	return p.Manojo(t3.Autor), t3.Razon, t3.Puntos
}

func Tipo4(p *pdt.Partida, m enco.Message) (*pdt.Manojo, enco.Razon) {
	var t1 enco.Tipo4
	json.Unmarshal(m.Cont, &t1)
	return p.Manojo(t1.Autor), t1.Razon
}

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
func Stringify(pkt *enco.Packet, p *pdt.Partida) string {
	s := Parse(p, pkt.Message)
	return strings.Replace(s, `"`, `'`, -1)
}

// Parse parsea un mensaje de salida y retorna su string correspondiente
func Parse(p *pdt.Partida, m enco.Message) string {

	var decoded string

	switch enco.CodMsg(m.Cod) {

	// (string)
	case enco.TError:
		err := Str(m)
		lower := strings.ToLower(err[:1]) + err[1:]
		decoded = fmt.Sprintf("Error, %s", lower)

	case enco.TLaManoResultaParda:
		decoded = fmt.Sprintf(`La Mano resulta parda`)

	case enco.TMazo:
		decoded = fmt.Sprintf(`%s se fue al mazo`, Autor(p, m).Jugador.ID)

	case enco.TByeBye:
		var template, s string

		if p.EsManoAMano() {
			template = "el ganador fue %s"
			s = Autor(p, m).Jugador.ID
		} else {
			template = "gano el equipo %s"
			s = Autor(p, m).Jugador.Equipo.String()
		}

		decoded = fmt.Sprintf("C'est fini! "+template, s)

	case enco.TDiceSonBuenas:
		decoded = fmt.Sprintf(`%s: "son buenas"`, Autor(p, m).Jugador.ID)

	case enco.TAbandono:
		autor := Autor(p, m)
		decoded = fmt.Sprintf(`%s abandono la partida. Gano el equipo %s`,
			autor.Jugador.ID, autor.Jugador.GetEquipoContrario())

	case enco.TCantarFlor:
		decoded = fmt.Sprintf(`%s canta flor`, Autor(p, m).Jugador.ID)

	case enco.TCantarContraFlor:
		decoded = fmt.Sprintf(`%s canta contra-flor`, Autor(p, m).Jugador.ID)

	case enco.TCantarContraFlorAlResto:
		decoded = fmt.Sprintf(`%s canta contra-flor al resto`, Autor(p, m).Jugador.ID)

	case enco.TTocarEnvido:
		decoded = fmt.Sprintf(`%s toca envido`, Autor(p, m).Jugador.ID)

	case enco.TTocarRealEnvido:
		decoded = fmt.Sprintf(`%s toca real envido`, Autor(p, m).Jugador.ID)

	case enco.TTocarFaltaEnvido:
		decoded = fmt.Sprintf(`%s toca falta envido`, Autor(p, m).Jugador.ID)

	case enco.TElEnvidoEstaPrimero:
		decoded = fmt.Sprintf(`%s "el envido esta primero!"`, Autor(p, m).Jugador.ID)

	case enco.TGritarTruco:
		decoded = fmt.Sprintf(`%s grita truco`, Autor(p, m).Jugador.ID)

	case enco.TGritarReTruco:
		decoded = fmt.Sprintf(`%s grita re-truco`, Autor(p, m).Jugador.ID)

	case enco.TGritarVale4:
		decoded = fmt.Sprintf(`%s grita vale-4`, Autor(p, m).Jugador.ID)

	case enco.TNoQuiero:
		decoded = fmt.Sprintf(`%s: "no quiero"`, Autor(p, m).Jugador.ID)

	case enco.TConFlorMeAchico:
		decoded = fmt.Sprintf(`%s: "con flor me achico"`, Autor(p, m).Jugador.ID)

	case enco.TQuieroTruco:
		decoded = fmt.Sprintf(`%s: "quiero"`, Autor(p, m).Jugador.ID)

	case enco.TQuieroEnvite:
		decoded = fmt.Sprintf(`%s: "quiero"`, Autor(p, m).Jugador.ID)

	// (int)
	case enco.TSigTurno:
		decoded = ""

	case enco.TSigTurnoPosMano:
		decoded = ""

	// (string, int)
	case enco.TDiceTengo:
		autor, valor := Tipo1(p, m)
		decoded = fmt.Sprintf(`%s: "tengo %d"`, autor.Jugador.ID, valor)

	case enco.TManoGanada:
		autor, valor := Tipo1(p, m)
		decoded = fmt.Sprintf(`La %s mano la gano el equipo %s gracias a %s`,
			pdt.NumMano(valor).String(), autor.Jugador.Equipo, autor.Jugador.ID)

	case enco.TRondaGanada:
		autor, razon := Tipo4(p, m)
		decoded = fmt.Sprintf(`La ronda ha sido ganada por el equipo %s debido al %s`,
			autor.Jugador.Equipo, enco.Razon(razon).String())

	case enco.TDiceSonMejores:
		autor, valor := Tipo1(p, m)
		decoded = fmt.Sprintf(`%s: "%d son mejores!"`, autor.Jugador.ID, valor)

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
		autor, razon, pts := Tipo3(p, m)
		if p.EsManoAMano() {
			decoded = fmt.Sprintf(`+%d pts para %s por %s`,
				pts, autor.Jugador.ID, Razon2str(string(razon)))
		} else {
			decoded = fmt.Sprintf(`+%d pts para el equipo %s por %s`,
				pts, autor.Jugador.Equipo.String(), Razon2str(string(razon)))
		}

	}

	return decoded
}
