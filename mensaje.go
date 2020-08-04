package truco

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
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
func (m *Msg) ToJSON() string {
	mJSON, _ := json.Marshal(m)
	return string(mJSON)
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

// ToBytes retorna la struct en bytes
func (csp ContSumPts) ToBytes() []byte {
	return []byte(fmt.Sprintf("%v", csp))
}

// ContTirarCarta ..
type ContTirarCarta struct {
	Autor string `json:"autor"` // perspectiva
	Carta Carta  `json:"carta"` // perspectiva
}

// ToBytes retorna la struct en bytes
func (ctc ContTirarCarta) ToBytes() []byte {
	return []byte(fmt.Sprintf("%v", ctc))
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
