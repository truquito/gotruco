package truco

// EstadoEnvido : enum
type EstadoEnvido int

// enums del envido
const (
	DESHABILITADO EstadoEnvido = 0
	NOCANTADOAUN  EstadoEnvido = 1
	ENVIDO        EstadoEnvido = 2
	REALENVIDO    EstadoEnvido = 3
	FALTAENVIDO   EstadoEnvido = 4
)

// Envido :
type Envido struct {
	puntaje    int
	cantadoPor *Jugador
	estado     EstadoEnvido
}

// estaHabilitado Devuelve `true` si el envido es `tocable`
func (e Envido) estaHabilitado() bool {
	return e.estado == NOCANTADOAUN || e.estado == ENVIDO
}

// deshabilitar el envido
func (e *Envido) deshabilitar() {
	e.estado = DESHABILITADO
}

// EstadoFlor : enum
type EstadoFlor int

// enums de la flor
const (
	NOCANTADA         EstadoFlor = 0
	FLOR              EstadoFlor = 1
	CONTRAFLOR        EstadoFlor = 2
	CONTRAFLORALRESTO EstadoFlor = 3
)