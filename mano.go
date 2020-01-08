package truco

// Resultado : Enum para el resultado de toda mano
type Resultado int

// 3 opciones (ya sea para 2, 4 o 6 jugadores)
const (
	GanoRojo  Resultado = 0
	GanoAzul  Resultado = 1
	Empardada Resultado = 2
)

// NumMano : Enum para el numero de la mano en juego
type NumMano int

// 3 opciones: Primera, seguna o tercera mano
const (
	primera NumMano = iota
	segunda
	tercera
)

func (n NumMano) toInt() int {
	switch n {
	case primera:
		return 1
	case segunda:
		return 2
	default:
		return 3
	}
}

// Mano :
type Mano struct {
	resultado Resultado
	ganador   *Manojo
	// en cada mano los jugadores van a tirar hata 1 carta
	cartasTiradas []tirarCarta
}

func (m *Mano) agregarCarta(cartaTirada tirarCarta) {
	m.cartasTiradas = append(m.cartasTiradas, cartaTirada)
}
