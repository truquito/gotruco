package truco

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
)

/*
 * -----------------------------------
 * | Cod. | Tipos de mensaje         |
 * |------|--------------------------|
 * |  00  | Error					 |
 * |  01  | Nueva-Partida			 |
 * |  02  | Nueva-Ronda				 |
 * |  03  | Tirar-Carta				 |
 * |  04  | Sumar-Puntos			 |
 * |  05  | Fin-Partida				 |
 * |  06  | Info					 |
 * |=================================|
 * |  07  | Toca-Envido				 |
 * |  08  | Toca-RealEnvido			 |
 * |  09  | Toca-FaltaEnvido		 |
 * |---------------------------------|
 * |  10  | Canta-Flor				 |
 * |  11  | Canta-ContraFlor		 |
 * |  12  | Canta-ContraFlorAlResto	 |
 * |  13  | Responde-ConFlorMeAchico |
 * |---------------------------------|
 * |  14  | Responde-Quiero			 |
 * |  15  | Responde-NoQuiero		 |
 * |---------------------------------|
 * |  16  | Mazo					 |
 * |---------------------------------|
 * |  17  | Grita-Truco				 |
 * |  18  | Grita-ReTruco			 |
 * |  19  | Grita-Vale4				 |
 * -----------------------------------
 */

// Pkt ..
type Pkt struct {
	Dest []string // eso no debe tener JSON porque no lo va a usar
	Msg
}

func (pkt *Pkt) String() string {
	return fmt.Sprintf("[%v] %s", strings.Join(pkt.Dest, ":"), pkt.Msg.String())
}

// Msg ..
type Msg struct {
	Tipo string `json:"tipo"`
	Nota string `json:"nota,omitempty"`
	Cont json.RawMessage
}

// GetAutor dado el contenido de un msg retorna su autor
func GetAutor(cont json.RawMessage) string {
	var autor string
	json.Unmarshal(cont, &autor)
	return autor
}

// GetConTirada parsea cont a ConTirada
func GetConTirada(cont json.RawMessage) ContTirarCarta {
	var conTirada ContTirarCarta
	json.Unmarshal(cont, &conTirada)
	return conTirada
}

// GetContSumPts parsea cont a ConSumPts
func GetContSumPts(cont json.RawMessage) ContSumPts {
	var conSumPts ContSumPts
	json.Unmarshal(cont, &conSumPts)
	return conSumPts
}

func (m Msg) String() string {

	switch m.Tipo {

	case "Error":
		return fmt.Sprintf("Error: %s", m.Nota)

	case "Info":
		return fmt.Sprintf("Info: %s", m.Nota)

	case "Toca-Envido":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s toca envido", autor)

	case "Toca-RealEnvido":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s toca envido", autor)

	case "Toca-FaltaEnvido":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s toca envido", autor)

	case "Canta-Flor":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s canta flor", autor)

	case "Canta-ContraFlor":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s canta contra-flor", autor)

	case "Canta-ContraFlorAlResto":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s canta contra-flor-al-resto", autor)

	case "Responde-ConFlorMeAchico":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s dice con-flor-me-achico", autor)

	case "Responde-Quiero":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s dice quiero", autor)

	case "Responde-NoQuiero":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s dice no-quiero", autor)

	case "Tirar-Carta":
		// autor := GetAutor(m.Cont)
		// carta := getCarta(m.Cont)
		return fmt.Sprintf(m.Nota)

	case "Sumar-Puntos":
		// autor := GetAutor(m.Cont)
		// carta := getCarta(m.Cont)
		return fmt.Sprintf(m.Nota)

	case "Mazo":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s se va al mazo", autor)

	case "Nueva-Partida":
		return fmt.Sprintf("arranca nueva partida")

	case "Nueva-Ronda":
		return fmt.Sprintf("arranca nueva ronda")

	case "Fin-Partida":
		return fmt.Sprintf(m.Nota)

	case "Grita-Truco":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s grita truco", autor)

	case "Grita-ReTruco":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s grita re-truco", autor)

	case "Grita-Vale4":
		autor := GetAutor(m.Cont)
		return fmt.Sprintf("%s grita vale4", autor)

	default:
		return "???"
	}
}

// ToJSON ..
func (m *Msg) ToJSON() json.RawMessage {
	mJSON, _ := json.Marshal(m)
	return mJSON
}

/* Contenidos */

// ContNuevaRonda ..
type ContNuevaRonda struct {
	Pers *PartidaDT `json:"pers"` // perspectiva
}

// ContSumPts ..
type ContSumPts struct {
	Pts    int    `json:"pts"`
	Equipo string `json:"equipo"`
}

// ToJSON retorna la struct en bytes
func (csp ContSumPts) ToJSON() json.RawMessage {
	cspJSON, _ := json.Marshal(csp)
	return cspJSON
}

// ContTirarCarta ..
type ContTirarCarta struct {
	Autor string `json:"autor"` // perspectiva
	Carta Carta  `json:"carta"` // perspectiva
}

// ToJSON retorna la struct en bytes
func (ctc ContTirarCarta) ToJSON() json.RawMessage {
	ctcJSON, _ := json.Marshal(ctc)
	return ctcJSON
}

func init() {
	// registrar las structs
	gob.Register(ContNuevaRonda{})
	gob.Register(ContSumPts{})
	gob.Register(ContTirarCarta{})
}
