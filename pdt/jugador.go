package pdt

// JIX :
type JIX int

// Jugador :
type Jugador struct {
	ID     string `json:"id"`
	jix    JIX    `json:"-"`
	Equipo Equipo `json:"equipo"`
}

// GetEquipoContrario retorna el equipo contrario
func (j Jugador) GetEquipoContrario() Equipo {
	if j.Equipo == Rojo {
		return Azul
	}
	return Rojo
}
