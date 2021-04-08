package enco

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Packet ..
type Packet struct {
	Destination *[]string
	*Message
}

func (pkt *Packet) String() string {
	return fmt.Sprintf("[%v] ", strings.Join(*pkt.Destination, ":"))
}

// Message .
type Message struct {
	Cod  int `json:"cod"`
	Cont json.RawMessage
}

// CodMsg ..
type CodMsg int

// Tipos de Mensajes
const (
	Error CodMsg = iota // 0
	ByeBye // 1
	DiceSonBuenas // 2
	CantarFlor // 3
	CantarContraFlor // 4
	CantarContraFlorAlResto // 5
	TocarEnvido // 6
	TocarRealEnvido // 7
	TocarFaltaEnvido // 8
	GritarTruco // 9
	GritarReTruco // 10
	GritarVale4 // 11
	NoQuiero // 12
	ConFlorMeAchico // 13
	QuieroTruco // 14
	QuieroEnvite // 15
	SigTurno // 16
	SigTurnoPosMano // 17
	DiceTengo // 18
	DiceSonMejores // 19
	NuevaPartida // 20
	NuevaRonda // 21
	TirarCarta // 22
	SumaPts // 23
	Mazo // 24
	TimeOut // 25
	ElEnvidoEstaPrimero // 26
	Abandono // 27
	LaManoResultaParda // 28
	ManoGanada // 29
	RondaGanada // 30
)

// Razon ..
type Razon int

// Razon por la que se suman puntos
const (
	EnvidoGanado Razon = iota
	RealEnvidoGanado
	FaltaEnvidoGanado

	EnviteNoQuerido
	FlorAchicada

	LaUnicaFlor
	LasFlores
	LaFlorMasAlta
	ContraFlorGanada
	ContraFlorAlRestoGanada

	TrucoNoQuerido
	TrucoQuerido
)

func (r Razon) String() string {
	var s string
	switch r {
	case EnvidoGanado:
		s = "envido"
	case RealEnvidoGanado:
		s = "real envido"
	case FaltaEnvidoGanado:
		s = "falta envido"
	case EnviteNoQuerido:
		s = "envite no querido"
	case FlorAchicada:
		s = "flor achicada"
	case LaUnicaFlor:
		s = "la unica flor"
	case LasFlores:
		s = "las flores"
	case LaFlorMasAlta:
		s = "la flor mas alta"
	case ContraFlorGanada:
		s = "contra flor ganada"
	case ContraFlorAlRestoGanada:
		s = "contra flor al resto"
	case TrucoNoQuerido:
		s = "truco no querido"
	case TrucoQuerido:
		s = "truco"
	}
	return s
}

// Tipo1 .
type Tipo1 struct {
	Autor string `json:"autor"`
	Valor int    `json:"valor"`
}

// Tipo2 .
type Tipo2 struct {
	Autor string `json:"autor"`
	Palo  int    `json:"palo"`
	Valor int    `json:"valor"`
}

// Tipo3 .
type Tipo3 struct {
	Autor  string `json:"autor"`
	Razon  int    `json:"razon"`
	Puntos int    `json:"puntos"`
}

// Pkt .
func Pkt(dest *[]string, m *Message) *Packet {
	return &Packet{
		Destination: dest,
		Message:     m,
	}
}

// Dest .
func Dest(ds ...string) *[]string {
	return &ds
}

// ParseStr dado un `Message` cuyo `Cont` es de tipo `string`
// lo extrae y retorna
func ParseStr(m *Message) string {
	var str string
	json.Unmarshal(m.Cont, &str)
	return str
}

// ParseInt dado un `Message` cuyo `Cont` es de tipo `int`
// lo extrae y retorna
func ParseInt(m *Message) int {
	var num int
	json.Unmarshal(m.Cont, &num)
	return num
}

// Msg .
func Msg(t CodMsg, data ...interface{}) *Message {

	var cont json.RawMessage

	switch t {
	case // (nil)
		LaManoResultaParda:

		cont = nil
	case // (string)
		Error,
		ByeBye,
		DiceSonBuenas,
		CantarFlor,
		CantarContraFlor,
		CantarContraFlorAlResto,
		TocarEnvido,
		TocarRealEnvido,
		TocarFaltaEnvido,
		GritarTruco,
		GritarReTruco,
		GritarVale4,
		NoQuiero,
		ConFlorMeAchico,
		QuieroTruco,
		QuieroEnvite,
		Mazo,
		ElEnvidoEstaPrimero,
		Abandono:

		bs, err := json.Marshal(data[0])

		if err != nil {
			log.Panic(err)
		}

		cont = bs

	case // (int)
		SigTurno,
		SigTurnoPosMano:

		bs, err := json.Marshal(data[0])

		if err != nil {
			log.Panic(err)
		}

		cont = bs

	case // (string, int)
		DiceTengo,
		DiceSonMejores,
		ManoGanada,
		RondaGanada:

		autor := data[0].(string)
		valor := data[1].(int)

		bs, _ := json.Marshal(&Tipo1{autor, valor})

		cont = bs

	case // (partida)
		NuevaPartida,
		NuevaRonda:

		pJSON, _ := data[0].(json.Marshaler).MarshalJSON()
		cont = pJSON

	case // (string, palo, valor)
		TirarCarta:

		autor := data[0].(string)
		palo := data[1].(int)
		valor := data[2].(int)

		bs, _ := json.Marshal(&Tipo2{autor, palo, valor})

		cont = bs

	case // (string, string, int)
		SumaPts:

		autor := data[0].(string)
		razon := int(data[1].(Razon))
		pts := data[2].(int)

		bs, _ := json.Marshal(&Tipo3{autor, razon, pts})

		cont = bs

	default:
		cont = nil

	}

	return &Message{
		Cod:  int(t),
		Cont: cont,
	}
}
