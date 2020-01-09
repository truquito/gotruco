package truco

// JugadorIdx :
type JugadorIdx int

// Jugador :
type Jugador struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Equipo Equipo `json:"equipo"`
}

func (j Jugador) getEquipoContrario() Equipo {
	if j.Equipo == Rojo {
		return Azul
	}
	return Rojo
}
