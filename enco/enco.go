package enco

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type IMessage interface {
	// json.Marshaler //
	Cod() CodMsg
}

// Packet ..
type Packet struct {
	Destination []string `json:"destination"`
	Message     Message  `json:"message"`
}

func (pkt *Packet) String() string {
	return fmt.Sprintf("[%v] ", strings.Join(pkt.Destination, ":"))
}

// Message .
type Message struct {
	Cod  CodMsg          `json:"cod"`
	Cont json.RawMessage `json:"cont"`
}

// CodMsg ..
type CodMsg string

// Tipos de Mensajes
const (
	TError                   CodMsg = "Error"
	TByeBye                  CodMsg = "ByeBye"
	TDiceSonBuenas           CodMsg = "DiceSonBuenas"
	TCantarFlor              CodMsg = "CantarFlor"
	TCantarContraFlor        CodMsg = "CantarContraFlor"
	TCantarContraFlorAlResto CodMsg = "CantarContraFlorAlResto"
	TTocarEnvido             CodMsg = "TocarEnvido"
	TTocarRealEnvido         CodMsg = "TocarRealEnvido"
	TTocarFaltaEnvido        CodMsg = "TocarFaltaEnvido"
	TGritarTruco             CodMsg = "GritarTruco"
	TGritarReTruco           CodMsg = "GritarReTruco"
	TGritarVale4             CodMsg = "GritarVale4"
	TNoQuiero                CodMsg = "NoQuiero"
	TConFlorMeAchico         CodMsg = "ConFlorMeAchico"
	TQuieroTruco             CodMsg = "QuieroTruco"
	TQuieroEnvite            CodMsg = "QuieroEnvite"
	TSigTurno                CodMsg = "SigTurno"
	TSigTurnoPosMano         CodMsg = "SigTurnoPosMano"
	TDiceTengo               CodMsg = "DiceTengo"
	TDiceSonMejores          CodMsg = "DiceSonMejores"
	TNuevaPartida            CodMsg = "NuevaPartida"
	TNuevaRonda              CodMsg = "NuevaRonda"
	TTirarCarta              CodMsg = "TirarCarta"
	TSumaPts                 CodMsg = "SumaPts"
	TMazo                    CodMsg = "Mazo"
	TTimeOut                 CodMsg = "TimeOut"
	TElEnvidoEstaPrimero     CodMsg = "ElEnvidoEstaPrimero"
	TAbandono                CodMsg = "Abandono"
	TLaManoResultaParda      CodMsg = "LaManoResultaParda"
	TManoGanada              CodMsg = "ManoGanada"
	TRondaGanada             CodMsg = "RondaGanada"
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
	Palo  string `json:"palo"`
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

// Pkt Packet maker
func Pkt(dest []string, m Message) *Packet {
	return &Packet{
		Destination: dest,
		Message:     m,
	}
}

// Dest Dest maker
func Dest(ds ...string) []string {
	return ds
}

// Msg Message maker
func Msg(t CodMsg, data ...interface{}) Message {

	var cont json.RawMessage

	switch t {
	case // (nil)
		TLaManoResultaParda:

		var s *string = nil
		cont, _ = json.Marshal(s)
	case // (string)
		TError,
		TByeBye,
		TDiceSonBuenas,
		TCantarFlor,
		TCantarContraFlor,
		TCantarContraFlorAlResto,
		TTocarEnvido,
		TTocarRealEnvido,
		TTocarFaltaEnvido,
		TGritarTruco,
		TGritarReTruco,
		TGritarVale4,
		TNoQuiero,
		TConFlorMeAchico,
		TQuieroTruco,
		TQuieroEnvite,
		TMazo,
		TElEnvidoEstaPrimero,
		TAbandono:

		bs, err := json.Marshal(data[0])

		if err != nil {
			log.Panic(err)
		}

		cont = bs

	case // (int)
		TSigTurno,
		TSigTurnoPosMano:

		bs, err := json.Marshal(data[0])

		if err != nil {
			log.Panic(err)
		}

		cont = bs

	case // (string, int)
		TDiceTengo,
		TDiceSonMejores,
		TManoGanada:

		autor := data[0].(string)
		valor := data[1].(int)

		bs, _ := json.Marshal(&Tipo1{autor, valor})

		cont = bs

	case // (string, string)
		TRondaGanada:

		autor := data[0].(string)
		razon := data[1].(Razon)

		bs, _ := json.Marshal(&Tipo4{autor, razon})

		cont = bs

	case // (partida)
		TNuevaPartida,
		TNuevaRonda:

		pJSON, _ := data[0].(json.Marshaler).MarshalJSON()
		cont = pJSON

	case // (string, palo, valor)
		TTirarCarta:

		autor := data[0].(string)
		palo := data[1].(string)
		valor := data[2].(int)

		bs, _ := json.Marshal(&Tipo2{autor, palo, valor})

		cont = bs

	case // (string, string, int)
		TSumaPts:

		autor := data[0].(string)
		razon := data[1].(Razon)
		pts := data[2].(int)

		bs, _ := json.Marshal(&Tipo3{autor, razon, pts})

		cont = bs

	default:
		cont = nil

	}

	return Message{
		Cod:  t,
		Cont: cont,
	}
}

//
//
//
//
//
//

type LaManoResultaParda struct{}

func (m LaManoResultaParda) Cod() CodMsg {
	return TLaManoResultaParda
}

type Error string

func (m Error) Cod() CodMsg {
	return TError
}

type ByeBye string

func (m ByeBye) Cod() CodMsg {
	return TByeBye
}

type DiceSonBuenas string

func (m DiceSonBuenas) Cod() CodMsg {
	return TDiceSonBuenas
}

type CantarFlor string

func (m CantarFlor) Cod() CodMsg {
	return TCantarFlor
}

type CantarContraFlor string

func (m CantarContraFlor) Cod() CodMsg {
	return TCantarContraFlor
}

type CantarContraFlorAlResto string

func (m CantarContraFlorAlResto) Cod() CodMsg {
	return TCantarContraFlorAlResto
}

type TocarEnvido string

func (m TocarEnvido) Cod() CodMsg {
	return TTocarEnvido
}

type TocarRealEnvido string

func (m TocarRealEnvido) Cod() CodMsg {
	return TTocarRealEnvido
}

type TocarFaltaEnvido string

func (m TocarFaltaEnvido) Cod() CodMsg {
	return TTocarFaltaEnvido
}

type GritarTruco string

func (m GritarTruco) Cod() CodMsg {
	return TGritarTruco
}

type GritarReTruco string

func (m GritarReTruco) Cod() CodMsg {
	return TGritarReTruco
}

type GritarVale4 string

func (m GritarVale4) Cod() CodMsg {
	return TGritarVale4
}

type NoQuiero string

func (m NoQuiero) Cod() CodMsg {
	return TNoQuiero
}

type ConFlorMeAchico string

func (m ConFlorMeAchico) Cod() CodMsg {
	return TConFlorMeAchico
}

type QuieroTruco string

func (m QuieroTruco) Cod() CodMsg {
	return TQuieroTruco
}

type QuieroEnvite string

func (m QuieroEnvite) Cod() CodMsg {
	return TQuieroEnvite
}

type Mazo string

func (m Mazo) Cod() CodMsg {
	return TMazo
}

type ElEnvidoEstaPrimero string

func (m ElEnvidoEstaPrimero) Cod() CodMsg {
	return TElEnvidoEstaPrimero
}

type Abandono string

func (m Abandono) Cod() CodMsg {
	return TAbandono
}

type SigTurno int

func (m SigTurno) Cod() CodMsg {
	return TSigTurno
}

type SigTurnoPosMano int

func (m SigTurnoPosMano) Cod() CodMsg {
	return TSigTurnoPosMano
}

type A struct {
	Autor string `json:"autor"`
	Valor int    `json:"valor"`
}

type DiceTengo A

func (m DiceTengo) Cod() CodMsg {
	return TDiceTengo
}

type DiceSonMejores A

func (m DiceSonMejores) Cod() CodMsg {
	return TDiceSonMejores
}

type ManoGanada A

func (m ManoGanada) Cod() CodMsg {
	return TManoGanada
}

// (string, string)
type RondaGanada struct {
	Autor string `json:"autor"`
	Razon Razon  `json:"razon"`
}

func (m RondaGanada) Cod() CodMsg {
	return TRondaGanada
}

// (partida) <- uso json.Marshaler para evitar ciclos de importacion
type NuevaPartida struct {
	Perspectiva interface{}
}

func (m NuevaPartida) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Perspectiva)
}

func (m NuevaPartida) Cod() CodMsg {
	return TNuevaPartida
}

// type NuevaRonda json.Marshaler
type NuevaRonda struct {
	Perspectiva interface{}
}

func (m NuevaRonda) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Perspectiva)
}

func (m NuevaRonda) Cod() CodMsg {
	return TNuevaRonda
}

// (string, palo, valor)
type TirarCarta struct {
	Autor string `json:"autor"`
	Palo  string `json:"palo"`
	Valor int    `json:"valor"`
}

func (m TirarCarta) Cod() CodMsg {
	return TTirarCarta
}

// (string, string, int)
type SumaPts struct {
	Autor  string `json:"autor"`
	Razon  Razon  `json:"razon"`
	Puntos int    `json:"puntos"`
}

func (m SumaPts) Cod() CodMsg {
	return TSumaPts
}

//
//
//
//

func Msg2(t CodMsg, data ...interface{}) IMessage {
	var m IMessage = nil

	switch t {
	case TLaManoResultaParda:
		m = LaManoResultaParda{}
	case TError:
		v, _ := data[0].(string)
		m = Error(v)
	case TByeBye:
		v, _ := data[0].(string)
		m = ByeBye(v)
	case TDiceSonBuenas:
		v, _ := data[0].(string)
		m = DiceSonBuenas(v)
	case TCantarFlor:
		v, _ := data[0].(string)
		m = CantarFlor(v)
	case TCantarContraFlor:
		v, _ := data[0].(string)
		m = CantarContraFlor(v)
	case TCantarContraFlorAlResto:
		v, _ := data[0].(string)
		m = CantarContraFlorAlResto(v)
	case TTocarEnvido:
		v, _ := data[0].(string)
		m = TocarEnvido(v)
	case TTocarRealEnvido:
		v, _ := data[0].(string)
		m = TocarRealEnvido(v)
	case TTocarFaltaEnvido:
		v, _ := data[0].(string)
		m = TocarFaltaEnvido(v)
	case TGritarTruco:
		v, _ := data[0].(string)
		m = GritarTruco(v)
	case TGritarReTruco:
		v, _ := data[0].(string)
		m = GritarReTruco(v)
	case TGritarVale4:
		v, _ := data[0].(string)
		m = GritarVale4(v)
	case TNoQuiero:
		v, _ := data[0].(string)
		m = NoQuiero(v)
	case TConFlorMeAchico:
		v, _ := data[0].(string)
		m = ConFlorMeAchico(v)
	case TQuieroTruco:
		v, _ := data[0].(string)
		m = QuieroTruco(v)
	case TQuieroEnvite:
		v, _ := data[0].(string)
		m = QuieroEnvite(v)
	case TMazo:
		v, _ := data[0].(string)
		m = Mazo(v)
	case TElEnvidoEstaPrimero:
		v, _ := data[0].(string)
		m = ElEnvidoEstaPrimero(v)
	case TAbandono:
		v, _ := data[0].(string)
		m = Abandono(v)
	case TSigTurno:
		v, _ := data[0].(int)
		m = SigTurno(v)
	case TSigTurnoPosMano:
		v, _ := data[0].(int)
		m = SigTurnoPosMano(v)
	case TDiceTengo:
		autor := data[0].(string)
		valor := data[1].(int)
		m = DiceTengo{autor, valor}
	case TDiceSonMejores:
		autor := data[0].(string)
		valor := data[1].(int)
		m = DiceSonMejores{autor, valor}
	case TManoGanada:
		autor := data[0].(string)
		valor := data[1].(int)
		m = ManoGanada{autor, valor}
	case TRondaGanada:
		autor := data[0].(string)
		razon := data[1].(Razon)
		m = RondaGanada{autor, razon}
	case TNuevaPartida:
		m = NuevaPartida{data[0]}
	case TNuevaRonda:
		m = NuevaRonda{data[0]}
	case TTirarCarta:
		autor := data[0].(string)
		palo := data[1].(string)
		valor := data[2].(int)
		m = TirarCarta{autor, palo, valor}
	case TSumaPts:
		autor := data[0].(string)
		razon := data[1].(Razon)
		pts := data[2].(int)
		m = SumaPts{autor, razon, pts}
	}

	return m
}
