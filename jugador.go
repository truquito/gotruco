package truco

import (
	"fmt"
)

// JugadorIdx :
type JugadorIdx int
// Jugador :
type Jugador struct{
	nombre 	string
	// idx			JugadorIdx
	equipo 	Equipo
	manojo 	*Manojo
}

// NuevoJugador : Retorna un nuevo jugador,
// con manojo vacio
func NuevoJugador(nombre string) *Jugador{
	return &Jugador{
		nombre: nombre,
		manojo: nil,
	}
}

// Print imprime la info de `jugador`
// y su manojo
func (jugador *Jugador) Print() {
	fmt.Printf("Jugador %s tiene las cartas:\n", jugador.nombre)
	jugador.manojo.Print()
}
