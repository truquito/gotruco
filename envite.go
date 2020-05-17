package truco

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

// toString
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

// Envite :
type Envite struct {
	Estado                        EstadoEnvite `json:"estado"`
	Puntaje                       int          `json:"puntaje"`
	CantadoPor                    *Manojo      `json:"cantadoPor"`
	JugadoresConFlor              []*Manojo    `json:"jugadoresConFlor"`
	JugadoresConFlorQueNoCantaron []*Manojo    `json:"jugadoresConFlorQueNoCantaron"`
}
