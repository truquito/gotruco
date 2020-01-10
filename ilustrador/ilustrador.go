package ilustrador

import (
	"fmt"

	"github.com/jpfilevich/canvas"
	"github.com/jpfilevich/truco"
)

// Imprimir un instante
func Imprimir(p truco.Partida) {
	canvas := canvas.NewCanvas()

	dibujarMarco(canvas)
	dibujarMuestra(canvas, 12, "Es")
	dibujarNombres(canvas, p.Ronda.Manojos)
	dibujarTiradas(canvas, p.Ronda.Manojos)

	render := canvas.Render()
	fmt.Printf(render)
}
