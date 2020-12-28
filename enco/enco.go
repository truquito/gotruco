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
	Error CodMsg = iota
	ByeBye
	DiceSonBuenas
	CantarFlor
	CantarContraFlor
	CantarContraFlorAlResto
	TocarEnvido
	TocarRealEnvido
	TocarFaltaEnvido
	GritarTruco
	GritarReTruco
	GritarVale4
	NoQuiero
	ConFlorMeAchico
	QuieroTruco
	QuieroEnvite
	SigTurno
	SigTurnoPosMano
	DiceTengo
	DiceSonMejores
	NuevaPartida
	NuevaRonda
	TirarCarta
	SumaPts
	Mazo
	TimeOut
	Info // <- se debe ir
	ElEnvidoEstaPrimero
	Abandono
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
		Info,
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
		DiceSonMejores:

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
