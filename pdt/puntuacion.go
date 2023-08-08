package pdt

// Puntuacion : Enum para el puntaje maximo de la partida
type Puntuacion int

// hasta 15 pts, 20 pts, 30 pts o 40 pts
const (
	A20 Puntuacion = 20
	A30 Puntuacion = 30
	A40 Puntuacion = 40
)

func (pt Puntuacion) toInt() int {
	return int(pt)
}
