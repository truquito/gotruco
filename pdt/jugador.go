package pdt

// JugadorIdx :
type JugadorIdx int

// Jugador :
type Jugador struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Equipo Equipo `json:"equipo"`
}

// GetEquipoContrario retorna el equipo contrario
func (j Jugador) GetEquipoContrario() Equipo {
	if j.Equipo == Rojo {
		return Azul
	}
	return Rojo
}
