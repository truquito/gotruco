package pdt

import (
	"bytes"
	"encoding/json"
)

// EstadoTruco : enum
type EstadoTruco int

// enums del truco
const (
	// el truco se "grita"
	NOGRITADOAUN EstadoTruco = iota
	TRUCO
	TRUCOQUERIDO
	RETRUCO
	RETRUCOQUERIDO
	VALE4
	VALE4QUERIDO
)

var toEstadoTruco = map[string]EstadoTruco{
	"noGritadoAun":   NOGRITADOAUN,
	"truco":          TRUCO,
	"trucoQuerido":   TRUCOQUERIDO,
	"reTruco":        RETRUCO,
	"reTrucoQuerido": RETRUCOQUERIDO,
	"vale4":          VALE4,
	"vale4Querido":   VALE4QUERIDO,
}

func (e EstadoTruco) esTrucoRespondible() bool {
	return e == TRUCO || e == RETRUCO || e == VALE4
}

func (e EstadoTruco) String() string {
	estados := []string{
		"noGritadoAun",
		"truco",
		"trucoQuerido",
		"reTruco",
		"reTrucoQuerido",
		"vale4",
		"vale4Querido",
	}

	ok := e >= 0 || int(e) < len(toEstadoTruco)
	if !ok {
		return "Unknown"
	}

	return estados[e]
}

// MarshalJSON marshals the enum as a quoted json string
func (e EstadoTruco) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(e.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *EstadoTruco) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*e = toEstadoTruco[j]
	return nil
}

// Truco :
type Truco struct {
	CantadoPor string      `json:"cantadoPor"`
	Estado     EstadoTruco `json:"estado"`
}
