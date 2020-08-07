package out

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Pkt ..
type Pkt struct {
	Dest *[]string
	*Msg
}

func (pkt *Pkt) String() string {
	return fmt.Sprintf("[%v] %s", strings.Join(*pkt.Dest, ":"), pkt.Msg.String())
}

// Msg .
type Msg struct {
	Cod  int `json:"cod"`
	Cont json.RawMessage
}

func (m Msg) String() string {
	return "el mensaje aqui"
}

// TipoMsg ..
type TipoMsg int

// asd
const (
	Error TipoMsg = iota
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
	TimeOut
)

// Razon ..
type Razon int

// asd
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

func pkt(dest *[]string, m *Msg) *Pkt {
	return &Pkt{
		Dest: dest,
		Msg:  m,
	}
}

func dest(ds ...string) *[]string {
	return &ds
}

func msg(t TipoMsg, data ...interface{}) *Msg {

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
		QuieroEnvite:

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

	return &Msg{
		Cod:  int(t),
		Cont: cont,
	}
}
