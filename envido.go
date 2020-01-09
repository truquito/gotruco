package truco

import (
	"bytes"
	"encoding/json"
)

// EstadoEnvido : enum
type EstadoEnvido int

// enums del envido
const (
	DESHABILITADO EstadoEnvido = iota
	NOCANTADOAUN
	ENVIDO
	REALENVIDO
	FALTAENVIDO
)

var toEstadoEnvido = map[string]EstadoEnvido{
	"deshabilitado": DESHABILITADO,
	"noCantadoAun":  NOCANTADOAUN,
	"envido":        ENVIDO,
	"realEnvido":    REALENVIDO,
	"faltaEnvido":   FALTAENVIDO,
}

// toString
func (e EstadoEnvido) String() string {
	estados := []string{
		"deshabilitado",
		"noCantadoAun",
		"envido",
		"realEnvido",
		"faltaEnvido",
	}

	ok := e >= 0 && int(e) < len(toEstadoEnvido)
	if !ok {
		return "Unknown"
	}

	return estados[e]
}

// MarshalJSON marshals the enum as a quoted json string
func (e EstadoEnvido) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(e.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *EstadoEnvido) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*e = toEstadoEnvido[j]
	return nil
}

// Envido :
type Envido struct {
	Puntaje    int          `json:"puntaje"`
	CantadoPor *Jugador     `json:"cantadoPor"`
	Estado     EstadoEnvido `json:"estado"`
}

// estaHabilitado Devuelve `true` si el envido es `tocable`
func (e Envido) estaHabilitado() bool {
	return e.Estado == NOCANTADOAUN || e.Estado == ENVIDO
}

// deshabilitar el envido
func (e *Envido) deshabilitar() {
	e.Estado = DESHABILITADO
}

// EstadoFlor : enum
type EstadoFlor int

// enums de la flor
const (
	DESHABILITADA     EstadoFlor = 0
	NOCANTADA         EstadoFlor = 1
	FLOR              EstadoFlor = 2
	CONTRAFLOR        EstadoFlor = 3
	CONTRAFLORALRESTO EstadoFlor = 4
)
