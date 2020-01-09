package ilustrador

func dibujarMarco(lienzo lienzo) {
	marco := templateMarco()
	lienzo.draw(11, 2, marco)
}

func dibujarMuestra(lienzo lienzo, valor int, palo string) {
	carta := templateCarta(valor, palo)
	lienzo.draw(25, 5, carta)
}

func dibujarNombres(lienzo lienzo, jugadores []string) {
	for i, nombre := range jugadores {
		if len(nombre) > 10 {
			nombre = nombre[:10]
		}
		area := areasJugadores["nombres"][posicion(i)]
		var nombreCentrado string
		if posicion(i) == f {
			nombreCentrado = area.right(nombre)
		} else if posicion(i) == c {
			nombreCentrado = area.left(nombre)
		} else {
			nombreCentrado = area.center(nombre)
		}
		lienzo.drawAt(area.from, nombreCentrado)
	}
}

func dibujarTiradas(lienzo lienzo, tiradas []Carta) {
	var area rectangle
	var ultimaTirada string
	// las de a
	ultimaTirada = templateCartaTriple(3, "Co")
	area = areasJugadores["tiradas"][posicion(a)]
	lienzo.drawAt(area.from, area.center(ultimaTirada))

	// las de b
	ultimaTirada = templateCarta(5, "Or")
	area = areasJugadores["tiradas"][posicion(b)]
	lienzo.drawAt(area.from, area.center(ultimaTirada))

	// las de c
	ultimaTirada = templateCarta(7, "Ba")
	area = areasJugadores["tiradas"][posicion(c)]
	lienzo.drawAt(area.from, area.center(ultimaTirada))

	// las de d
	ultimaTirada = templateCarta(10, "Es")
	area = areasJugadores["tiradas"][posicion(d)]
	lienzo.drawAt(area.from, area.center(ultimaTirada))

	// las de e
	ultimaTirada = templateCarta(12, "Co")
	area = areasJugadores["tiradas"][posicion(e)]
	lienzo.drawAt(area.from, area.center(ultimaTirada))

	// las de f
	ultimaTirada = templateCarta(1, "Es")
	area = areasJugadores["tiradas"][posicion(f)]
	lienzo.drawAt(area.from, area.center(ultimaTirada))
}
