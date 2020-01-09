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
	Puntaje    int          `json:"puntaje"`
	CantadoPor *Jugador     `json:"cantadoPor"`
	Estado     EstadoEnvido `json:"estado"`
}

// estaHabilitado Devuelve `true` si el envido es `tocable`
func (e Envido) estaHabilitado() bool {
	return e.Estado == NOCANTADOAUN || e.Estado == ENVIDO
}

// deshabilitar el envido
func (e *Envido) deshabilitar() {
	e.Estado = DESHABILITADO
}

// EstadoFlor : enum
type EstadoFlor int

// enums de la flor
const (
	DESHABILITADA     EstadoFlor = 0
	NOCANTADA         EstadoFlor = 1
	FLOR              EstadoFlor = 2
	CONTRAFLOR        EstadoFlor = 3
	CONTRAFLORALRESTO EstadoFlor = 4
)
