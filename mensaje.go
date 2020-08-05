package truco

import (
	"encoding/gob"
	"encoding/json"
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

// IContenido ...
type IContenido interface {
	// hacer(p *Partida)
	// getAutor() *Manojo
}

// Msg ..
type Msg struct {
	Tipo string `json:"tipo"`
	Nota string `json:"nota,omitempty"`
	Cont json.RawMessage
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

// ContAutor ..
type ContAutor struct {
	Autor string `json:"autor"` // perspectiva
}

func init() {
	// registrar las structs
	gob.Register(ContNuevaRonda{})
	gob.Register(ContSumPts{})
	gob.Register(ContTirarCarta{})
	gob.Register(ContAutor{})
}
