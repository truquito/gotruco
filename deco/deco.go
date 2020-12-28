package deco

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/filevich/truco/out"
	"github.com/filevich/truco/pdt"
)

// Manojo .
func Manojo(p *pdt.PartidaDT, id string) *pdt.Manojo {
	manojo, _ := p.Ronda.GetManojoByStr(id)
	return manojo
}

// Str .
func Str(m *out.Message) string {
	var str string
	json.Unmarshal(m.Cont, &str)
	return str
}

// Int .
func Int(m *out.Message) int {
	var num int
	json.Unmarshal(m.Cont, &num)
	return num
}

// Autor .
func Autor(p *pdt.PartidaDT, m *out.Message) *pdt.Manojo {
	id := out.ParseStr(m)
	return Manojo(p, id)
}

// Tipo1 .
func Tipo1(p *pdt.PartidaDT, m *out.Message) (*pdt.Manojo, int) {
	var t1 out.Tipo1
	json.Unmarshal(m.Cont, &t1)
	return Manojo(p, t1.Autor), t1.Valor
}

// Tipo2 .
func Tipo2(p *pdt.PartidaDT, m *out.Message) (*pdt.Manojo, pdt.Palo, int) {
	var t2 out.Tipo2
	json.Unmarshal(m.Cont, &t2)
	palo := pdt.Palo(t2.Palo)
	return Manojo(p, t2.Autor), palo, t2.Valor
}

// Tipo3 .
func Tipo3(p *pdt.PartidaDT, m *out.Message) (*pdt.Manojo, int, int) {
	var t3 out.Tipo3
	json.Unmarshal(m.Cont, &t3)
	return Manojo(p, t3.Autor), t3.Razon, t3.Puntos
}

// Razon2str retorna el string correspondiente a `r`
func Razon2str(r int) string {
	var str string
	switch out.Razon(r) {
	case out.EnvidoGanado:
		str = "el envido ganado"
	case out.RealEnvidoGanado:
		str = "el real envido ganado"
	case out.FaltaEnvidoGanado:
		str = "la falta envido ganada"
	case out.EnviteNoQuerido:
		str = "el envido no querido"
	case out.FlorAchicada:
		str = "la flor achicada"
	case out.LaUnicaFlor:
		str = "ser la unica flor"
	case out.LasFlores:
		str = "todas las flores"
	case out.LaFlorMasAlta:
		str = "tener la flor mas alta"
	case out.ContraFlorGanada:
		str = "la contra-flor ganada"
	case out.ContraFlorAlRestoGanada:
		str = "la contra-flor al resto ganada"
	case out.TrucoNoQuerido:
		str = "el truco no querido"
	case out.TrucoQuerido:
		str = "el truco ganado"
	}
	return str
}

// Stringify parsea un pkt
// de momento solo su contenido (el msg)
func Stringify(pkt *out.Packet, p *pdt.PartidaDT) string {
	s := Parse(p, pkt.Message)
	return strings.Replace(s, `"`, `'`, -1)
}

// Parse parsea un mensaje de salida y retorna su string correspondiente
func Parse(p *pdt.PartidaDT, m *out.Message) string {

	var decoded string

	switch out.CodMsg(m.Cod) {

	// (string)
	case out.Error:
		err := Str(m)
		lower := strings.ToLower(err[:1]) + err[1:]
		decoded = fmt.Sprintf("Error, %s", lower)

	case out.Info:
		err := Str(m)
		lower := strings.ToLower(err[:1]) + err[1:]
		decoded = fmt.Sprintf("Info, %s", lower)

	case out.Mazo:
		decoded = fmt.Sprintf(`%s se fue al mazo`, Autor(p, m).Jugador.Nombre)

	case out.ByeBye:
		var template, s string

		if p.EsManoAMano() {
			template = "el ganador fue %s"
			s = Autor(p, m).Jugador.Nombre
		} else {
			template = "gano el equipo %s"
			s = Autor(p, m).Jugador.Equipo.String()
		}

		decoded = fmt.Sprintf("C'est fini! "+template, s)

	case out.DiceSonBuenas:
		decoded = fmt.Sprintf(`%s: "son buenas"`, Autor(p, m).Jugador.Nombre)

	case out.CantarFlor:
		decoded = fmt.Sprintf(`%s canta flor`, Autor(p, m).Jugador.Nombre)

	case out.CantarContraFlor:
		decoded = fmt.Sprintf(`%s canta contra-flor`, Autor(p, m).Jugador.Nombre)

	case out.CantarContraFlorAlResto:
		decoded = fmt.Sprintf(`%s canta contra-flor al resto`, Autor(p, m).Jugador.Nombre)

	case out.TocarEnvido:
		decoded = fmt.Sprintf(`%s toca envido`, Autor(p, m).Jugador.Nombre)

	case out.TocarRealEnvido:
		decoded = fmt.Sprintf(`%s toca real envido`, Autor(p, m).Jugador.Nombre)

	case out.TocarFaltaEnvido:
		decoded = fmt.Sprintf(`%s toca falta envido`, Autor(p, m).Jugador.Nombre)

	case out.GritarTruco:
		decoded = fmt.Sprintf(`%s grita truco`, Autor(p, m).Jugador.Nombre)

	case out.GritarReTruco:
		decoded = fmt.Sprintf(`%s grita re-truco`, Autor(p, m).Jugador.Nombre)

	case out.GritarVale4:
		decoded = fmt.Sprintf(`%s grita vale-4`, Autor(p, m).Jugador.Nombre)

	case out.NoQuiero:
		decoded = fmt.Sprintf(`%s: "no quiero"`, Autor(p, m).Jugador.Nombre)

	case out.ConFlorMeAchico:
		decoded = fmt.Sprintf(`%s: "con flor me achico"`, Autor(p, m).Jugador.Nombre)

	case out.QuieroTruco:
		decoded = fmt.Sprintf(`%s: "quiero"`, Autor(p, m).Jugador.Nombre)

	case out.QuieroEnvite:
		decoded = fmt.Sprintf(`%s: "quiero"`, Autor(p, m).Jugador.Nombre)

	// (int)
	case out.SigTurno:
		decoded = ""

	case out.SigTurnoPosMano:
		decoded = ""

	// (string, int)
	case out.DiceTengo:
		autor, valor := Tipo1(p, m)
		decoded = fmt.Sprintf(`%s: "tengo %d"`, autor.Jugador.Nombre, valor)

	case out.DiceSonMejores:
		autor, valor := Tipo1(p, m)
		decoded = fmt.Sprintf(`%s: "%d son mejores!"`, autor.Jugador.Nombre, valor)

	// (partida)
	case out.NuevaPartida:
		decoded = ""

	case out.NuevaRonda:
		decoded = "Empieza nueva ronda"

	// (string, palo, valor)
	case out.TirarCarta:
		// autor, palo, valor := Tipo2(p, m)
		// fmt.Printf("detectado codigo %d\n", out.CodMsg(m.Cod))
		// decoded = fmt.Sprintf(`%s tira %d de %s`, autor.Jugador.Nombre, valor, palo.String())
		decoded = ""

	// (string, string, int)
	case out.SumaPts:
		autor, razon, pts := Tipo3(p, m)
		if p.EsManoAMano() {
			decoded = fmt.Sprintf(`+%d pts para %s por %s`,
				pts, autor.Jugador.Nombre, Razon2str(razon))
		} else {
			decoded = fmt.Sprintf(`+%d pts para el equipo %s por %s`,
				pts, autor.Jugador.Equipo.String(), Razon2str(razon))
		}

	}

	return decoded
}
