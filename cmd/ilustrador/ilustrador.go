package ilustrador

import (
	"fmt"
)

// Instante de una partida
type Instante struct {
	Jugadores []string
	Tiradas   []Carta
	Turno     int
	Mano      int
}

// Imprimir un instante
func Imprimir(instante Instante) {
	lienzo := nuevoLienzo()

	dibujarMarco(lienzo)
	dibujarMuestra(lienzo, 12, "Es")
	dibujarNombres(lienzo, instante.Jugadores)
	dibujarTiradas(lienzo, instante.Tiradas)

	render := lienzo.renderizar()
	fmt.Printf(render)
}
