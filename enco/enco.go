package enco

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Packet ..
type Packet struct {
	Destination *[]string `json:"destination"`
	Message     *Message  `json:"message"`
}

func (pkt *Packet) String() string {
	return fmt.Sprintf("[%v] ", strings.Join(*pkt.Destination, ":"))
}

// Message .
type Message struct {
	Cod  string          `json:"cod"`
	Cont json.RawMessage `json:"cont"`
}

// CodMsg ..
type CodMsg string

// Tipos de Mensajes
const (
	Error                   CodMsg = "Error"
	ByeBye                  CodMsg = "ByeBye"
	DiceSonBuenas           CodMsg = "DiceSonBuenas"
	CantarFlor              CodMsg = "CantarFlor"
	CantarContraFlor        CodMsg = "CantarContraFlor"
	CantarContraFlorAlResto CodMsg = "CantarContraFlorAlResto"
	TocarEnvido             CodMsg = "TocarEnvido"
	TocarRealEnvido         CodMsg = "TocarRealEnvido"
	TocarFaltaEnvido        CodMsg = "TocarFaltaEnvido"
	GritarTruco             CodMsg = "GritarTruco"
	GritarReTruco           CodMsg = "GritarReTruco"
	GritarVale4             CodMsg = "GritarVale4"
	NoQuiero                CodMsg = "NoQuiero"
	ConFlorMeAchico         CodMsg = "ConFlorMeAchico"
	QuieroTruco             CodMsg = "QuieroTruco"
	QuieroEnvite            CodMsg = "QuieroEnvite"
	SigTurno                CodMsg = "SigTurno"
	SigTurnoPosMano         CodMsg = "SigTurnoPosMano"
	DiceTengo               CodMsg = "DiceTengo"
	DiceSonMejores          CodMsg = "DiceSonMejores"
	NuevaPartida            CodMsg = "NuevaPartida"
	NuevaRonda              CodMsg = "NuevaRonda"
	TirarCarta              CodMsg = "TirarCarta"
	SumaPts                 CodMsg = "SumaPts"
	Mazo                    CodMsg = "Mazo"
	TimeOut                 CodMsg = "TimeOut"
	ElEnvidoEstaPrimero     CodMsg = "ElEnvidoEstaPrimero"
	Abandono                CodMsg = "Abandono"
	LaManoResultaParda      CodMsg = "LaManoResultaParda"
	ManoGanada              CodMsg = "ManoGanada"
	RondaGanada             CodMsg = "RondaGanada"
)

// Razon ..
type Razon string

// Razon por la que se suman puntos
const (
	EnvidoGanado      Razon = "EnvidoGanado"
	RealEnvidoGanado  Razon = "RealEnvidoGanado"
	FaltaEnvidoGanado Razon = "FaltaEnvidoGanado"

	EnviteNoQuerido Razon = "EnviteNoQuerido"
	FlorAchicada    Razon = "FlorAchicada"

	LaUnicaFlor             Razon = "LaUnicaFlor"
	LasFlores               Razon = "LasFlores"
	LaFlorMasAlta           Razon = "LaFlorMasAlta"
	ContraFlorGanada        Razon = "ContraFlorGanada"
	ContraFlorAlRestoGanada Razon = "ContraFlorAlRestoGanada"

	TrucoNoQuerido Razon = "TrucoNoQuerido"
	TrucoQuerido   Razon = "TrucoQuerido"

	SeFueronAlMazo Razon = "SeFueronAlMazo"
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
	Razon  Razon  `json:"razon"`
	Puntos int    `json:"puntos"`
}

type Tipo4 struct {
	Autor string `json:"autor"`
	Razon Razon  `json:"razon"`
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

		var s *string = nil
		cont, _ = json.Marshal(s)
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
		ManoGanada:

		autor := data[0].(string)
		valor := data[1].(int)

		bs, _ := json.Marshal(&Tipo1{autor, valor})

		cont = bs

	case // (string, string)
		RondaGanada:

		autor := data[0].(string)
		razon := data[1].(Razon)

		bs, _ := json.Marshal(&Tipo4{autor, razon})

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
		razon := data[1].(Razon)
		pts := data[2].(int)

		bs, _ := json.Marshal(&Tipo3{autor, razon, pts})

		cont = bs

	default:
		cont = nil

	}

	return &Message{
		Cod:  string(t),
		Cont: cont,
	}
}
