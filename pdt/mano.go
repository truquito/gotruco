package pdt

import (
	"bytes"
	"encoding/json"
)

// Resultado : Enum para el resultado de toda mano
type Resultado int

// 3 opciones (ya sea para 2, 4 o 6 jugadores)
const (
	GanoRojo Resultado = iota
	GanoAzul
	Empardada
)

var toResultado = map[string]Resultado{
	"ganoRojo":  GanoRojo,
	"ganoAzul":  GanoAzul,
	"empardada": Empardada,
}

func (r Resultado) String() string {
	resultados := []string{
		"ganoRojo",
		"ganoAzul",
		"empardada",
	}

	ok := r >= 0 || int(r) < len(toResultado)
	if !ok {
		return "Unknown"
	}

	return resultados[r]
}

// MarshalJSON marshals the enum as a quoted json string
func (r Resultado) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(r.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (r *Resultado) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*r = toResultado[j]
	return nil
}

// NumMano : Enum para el numero de la mano en juego
type NumMano int

// 3 opciones: Primera, seguna o Tercera mano
const (
	Primera NumMano = iota
	Segunda
	Tercera
)

// ToInt .
func (n NumMano) ToInt() int {
	switch n {
	case Primera:
		return 1
	case Segunda:
		return 2
	default:
		return 3
	}
}

func (n NumMano) String() string {
	manos := []string{
		"Primera",
		"Segunda",
		"Tercera",
	}

	ok := n >= 0 || int(n) < len(manos)
	if !ok {
		return "Unknown"
	}

	return manos[n]
}

type cartaTirada struct {
	autor *Manojo
	Carta
}

// Mano :
type Mano struct {
	Resultado Resultado `json:"resultado"`
	Ganador   string    `json:"ganador"`
	// en cada mano los jugadores van a tirar hata 1 carta
	CartasTiradas []cartaTirada `json:"cartasTiradas"`
}

func (m *Mano) agregarTirada(cartaTirada cartaTirada) {
	m.CartasTiradas = append(m.CartasTiradas, cartaTirada)
}
