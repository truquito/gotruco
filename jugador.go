package truco

// JugadorIdx :
type JugadorIdx int

// Jugador :
type Jugador struct {
	id     string
	nombre string
	equipo Equipo
}

func (j Jugador) getEquipoContrario() Equipo {
	if j.equipo == Rojo {
		return Azul
	}
	return Rojo
}
