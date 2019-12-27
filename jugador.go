package truco

import (
	"fmt"
)

// JugadorIdx :
type JugadorIdx int

// Jugador :
type Jugador struct {
	nombre string
	// idx			JugadorIdx
	equipo Equipo
}

// Print imprime la info de `jugador`
// y su manojo
func (jugador *Jugador) Print() {
	fmt.Printf("%s:\n", jugador.nombre)
}
