package truco

import (
	"encoding/gob"
	"encoding/json"
)

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
