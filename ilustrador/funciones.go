package ilustrador

import (
	"github.com/jpfilevich/canvas"
	"github.com/jpfilevich/truco"
)

func dibujarMarco(canvas canvas.Canvas) {
	marco := templateMarco()
	canvas.Draw(11, 2, marco)
}

func dibujarMuestra(canvas canvas.Canvas, valor int, palo string) {
	carta := templateCarta(valor, palo)
	canvas.Draw(25, 5, carta)
}

func dibujarNombres(canvas canvas.Canvas, manojos []truco.Manojo) {
	for i, manojo := range manojos {
		nombre := manojo.Jugador.Nombre
		if len(nombre) > 10 {
			nombre = nombre[:10]
		}
		area := areasJugadores["nombres"][posicion(i)]
		var nombreCentrado string
		if posicion(i) == f {
			nombreCentrado = area.Right(nombre)
		} else if posicion(i) == c {
			nombreCentrado = area.Left(nombre)
		} else {
			nombreCentrado = area.Center(nombre)
		}
		canvas.DrawAt(area.From, nombreCentrado)
	}
}

func dibujarTiradas(ca canvas.Canvas, manojos []truco.Manojo) {
	var area canvas.Rectangle
	var ultimaTirada string
	// las de a
	ultimaTirada = templateCartaTriple(3, "Co")
	area = areasJugadores["tiradas"][posicion(a)]
	ca.DrawAt(area.From, area.Center(ultimaTirada))

	// las de b
	ultimaTirada = templateCarta(5, "Or")
	area = areasJugadores["tiradas"][posicion(b)]
	ca.DrawAt(area.From, area.Center(ultimaTirada))

	// las de ca
	ultimaTirada = templateCarta(7, "Ba")
	area = areasJugadores["tiradas"][posicion(c)]
	ca.DrawAt(area.From, area.Center(ultimaTirada))

	// las de d
	ultimaTirada = templateCarta(10, "Es")
	area = areasJugadores["tiradas"][posicion(d)]
	ca.DrawAt(area.From, area.Center(ultimaTirada))

	// las de e
	ultimaTirada = templateCarta(12, "Co")
	area = areasJugadores["tiradas"][posicion(e)]
	ca.DrawAt(area.From, area.Center(ultimaTirada))

	// las de f
	ultimaTirada = templateCarta(1, "Es")
	area = areasJugadores["tiradas"][posicion(f)]
	ca.DrawAt(area.From, area.Center(ultimaTirada))
}
