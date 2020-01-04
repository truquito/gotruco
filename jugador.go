package truco

// JugadorIdx :
type JugadorIdx int

// Jugador :
type Jugador struct {
	nombre string
	// idx			JugadorIdx
	equipo Equipo
}

func (j Jugador) getEquipoContrario() Equipo {
	if j.equipo == Rojo {
		return Azul
	}
	return Rojo
}
