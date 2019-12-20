package truco

// Resultado : Enum para el resultado de toda mano
type Resultado int

// 3 opciones (ya sea para 2, 4 o 6 jugadores)
const (
	GanoRojo      Resultado = 0
	GanoAzul      Resultado = 1
	Empardada Resultado = 2
)

// NumMano : Enum para el numero de la mano en juego
type NumMano int

// 3 opciones: Primera, seguna o tercera mano
const (
	primera NumMano = 0
	segunda NumMano = 1
	tercera	NumMano = 2
)

// Mano :
type Mano struct {
	repartidor 	JugadorIdx
	resultado  	Resultado
	ganador			JugadorIdx
}