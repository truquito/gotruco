package pdt

import (
	"bytes"
	"encoding/json"
)

// EstadoEnvite : enum
type EstadoEnvite int

// enums del envite
const (
	DESHABILITADO EstadoEnvite = iota
	NOCANTADOAUN
	ENVIDO
	REALENVIDO
	FALTAENVIDO
	FLOR
	CONTRAFLOR
	CONTRAFLORALRESTO
)

var toEstadoEnvite = map[string]EstadoEnvite{
	"deshabilitado":     DESHABILITADO,
	"noCantadoAun":      NOCANTADOAUN,
	"envido":            ENVIDO,
	"realEnvido":        REALENVIDO,
	"faltaEnvido":       FALTAENVIDO,
	"flor":              FLOR,
	"contraFlor":        CONTRAFLOR,
	"contraFlorAlResto": CONTRAFLORALRESTO,
}

func (e EstadoEnvite) String() string {
	estados := []string{
		"deshabilitado",
		"noCantadoAun",
		"envido",
		"realEnvido",
		"faltaEnvido",
		"flor",
		"contraFlor",
		"contraFlorAlResto",
	}

	ok := e >= 0 && int(e) < len(estados)
	if !ok {
		return "Unknown"
	}

	return estados[e]
}

// MarshalJSON marshals the enum as a quoted json string
func (e EstadoEnvite) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(e.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *EstadoEnvite) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*e = toEstadoEnvite[j]
	return nil
}

func (e *Envite) noCantoFlorAun(j string) bool {
	for _, id := range e.SinCantar {
		if id == j {
			return true
		}
	}
	return false
}

// Elimina a `j` de los jugadores que tienen pendiente cantar flor
func (e *Envite) cantoFlor(j string) {
	// lo elimino
	xs := e.SinCantar // un slice es un puntero

	for i, x := range e.SinCantar {
		if x == j {
			xs[i] = xs[len(xs)-1]
			xs = xs[:len(xs)-1]
			e.SinCantar = xs
			return
		}
	}

}

// Envite :
type Envite struct {
	Estado           EstadoEnvite `json:"estado"`
	Puntaje          int          `json:"puntaje"`
	CantadoPor       string       `json:"cantadoPor"`
	JugadoresConFlor []*Manojo    `json:"-"`
	SinCantar        []string     `json:"sinCantar"`
}
