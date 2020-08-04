package truco

import (
	"encoding/json"
)

// Pkt ..
type Pkt struct {
	Dest []string // eso no debe tener JSON porque no lo va a usar
	Msg2
}

// IContenido ...
type IContenido interface {
	// hacer(p *Partida)
	// getAutor() *Manojo
}

// Msg2 ..
type Msg2 struct {
	Tipo string     `json:"tipo"`
	Nota string     `json:"nota,omitempty"`
	Cont IContenido `json:"contenido,omitempty"`
}

// ToJSON ..
func (m *Msg2) ToJSON() string {
	mJSON, _ := json.Marshal(m)
	return string(mJSON)
}

/* Contenidos */

// ContNuevaRonda ..
type ContNuevaRonda struct {
	Pers string `json:"pers"` // perspectiva
}

// ContSumPts ..
type ContSumPts struct {
	Pts    int    `json:"pts"`
	Equipo string `json:"equipo"`
}

// ContTirarCarta ..
type ContTirarCarta struct {
	Autor string `json:"autor"` // perspectiva
	Carta Carta  `json:"carta"` // perspectiva
}

// ContAutor ..
type ContAutor struct {
	Autor string `json:"autor"` // perspectiva
}
