package truco

import (
	"fmt"
	"strconv"

	"github.com/jpfilevich/canvas"
)

type posicion int

const (
	a posicion = iota
	b
	c
	d
	e
	f
)

type templates struct{}

func (t templates) marco() string {
	marco := canvas.Raw(`
╔══════════════════════════════╗
║                              ║
║                              ║
║                              ║
║                              ║
║                              ║
║                              ║
║                              ║
╚══════════════════════════════╝
	`)
	return marco
}

func (t templates) carta(valor int, palo string) string {
	carta := "┌xx┐" + "\n"
	carta += "│PP│" + "\n"
	carta += "└──┘"

	// numero
	numStr := strconv.Itoa(valor)
	if valor <= 9 {
		numStr = numStr + "─"
	}
	carta = canvas.Replace("xx", numStr, carta)
	carta = canvas.Replace("PP", palo[:2], carta)

	return carta
}

func (t templates) cartaDoble(valor int, palo string) string {
	cartaDoble := canvas.Raw(`
┌xx┐┐
│PP││
└──┘┘
`)

	// numero
	numStr := strconv.Itoa(valor)
	if valor <= 9 {
		numStr = numStr + "─"
	}
	cartaDoble = canvas.Replace("xx", numStr, cartaDoble)
	cartaDoble = canvas.Replace("PP", palo[:2], cartaDoble)

	return cartaDoble
}

func (t templates) cartaTriple(valor int, palo string) string {
	cartaTriple := canvas.Raw(`
┌xx┐┐┐
│PP│││
└──┘┘┘
`)

	// numero
	numStr := strconv.Itoa(valor)
	if valor <= 9 {
		numStr = numStr + "─"
	}
	cartaTriple = canvas.Replace("xx", numStr, cartaTriple)
	cartaTriple = canvas.Replace("PP", palo[:2], cartaTriple)

	return cartaTriple
}

type impresora struct {
	canvas         canvas.Canvas
	areasJugadores map[string](map[posicion]canvas.Rectangle)
	otrasAreas     map[string]canvas.Rectangle
	templates
}

func (pr impresora) dibujarMarco() {
	marco := pr.templates.marco()
	pr.canvas.Draw(11, 2, marco)
}

func (pr impresora) dibujarMuestra(muestra Carta) {
	carta := pr.templates.carta(muestra.Valor, muestra.Palo.String())
	pr.canvas.Draw(25, 5, carta)
}

func (pr impresora) dibujarNombres(manojos []Manojo) {
	for i, manojo := range manojos {
		nombre := manojo.Jugador.Nombre
		if len(nombre) > 10 {
			nombre = nombre[:10]
		}
		area := pr.areasJugadores["nombres"][posicion(i)]
		var nombreCentrado string
		if posicion(i) == f {
			nombreCentrado = area.Right(nombre)
		} else if posicion(i) == c {
			nombreCentrado = area.Left(nombre)
		} else {
			nombreCentrado = area.Center(nombre)
		}
		pr.canvas.DrawAt(area.From, nombreCentrado)
	}
}

func (pr impresora) dibujarTiradas(manojos []Manojo) {
	var area canvas.Rectangle
	var ultimaTirada string
	// las de a
	ultimaTirada = pr.templates.cartaTriple(3, "Co")
	area = pr.areasJugadores["tiradas"][posicion(a)]
	pr.canvas.DrawAt(area.From, area.Center(ultimaTirada))

	// las de b
	ultimaTirada = pr.templates.carta(5, "Or")
	area = pr.areasJugadores["tiradas"][posicion(b)]
	pr.canvas.DrawAt(area.From, area.Center(ultimaTirada))

	// las de ca
	ultimaTirada = pr.templates.carta(7, "Ba")
	area = pr.areasJugadores["tiradas"][posicion(c)]
	pr.canvas.DrawAt(area.From, area.Center(ultimaTirada))

	// las de d
	ultimaTirada = pr.templates.carta(10, "Es")
	area = pr.areasJugadores["tiradas"][posicion(d)]
	pr.canvas.DrawAt(area.From, area.Center(ultimaTirada))

	// las de e
	ultimaTirada = pr.templates.carta(12, "Co")
	area = pr.areasJugadores["tiradas"][posicion(e)]
	pr.canvas.DrawAt(area.From, area.Center(ultimaTirada))

	// las de f
	ultimaTirada = pr.templates.carta(1, "Es")
	area = pr.areasJugadores["tiradas"][posicion(f)]
	pr.canvas.DrawAt(area.From, area.Center(ultimaTirada))
}

func (pr impresora) Print(p *Partida) {
	pr.dibujarMarco()
	pr.dibujarMuestra(p.Ronda.Muestra)
	pr.dibujarNombres(p.Ronda.Manojos)
	pr.dibujarTiradas(p.Ronda.Manojos)

	render := pr.canvas.Render()
	fmt.Printf(render)
}

func nuevaImpresora() impresora {
	return impresora{
		canvas: canvas.NewCanvas(),
		areasJugadores: map[string](map[posicion]canvas.Rectangle){
			"nombres": map[posicion]canvas.Rectangle{
				a: canvas.Rectangle{
					From: canvas.Point{X: 15, Y: 11},
					To:   canvas.Point{X: 25, Y: 11},
				},
				b: canvas.Rectangle{
					From: canvas.Point{X: 29, Y: 11},
					To:   canvas.Point{X: 39, Y: 11},
				},
				c: canvas.Rectangle{
					From: canvas.Point{X: 44, Y: 6},
					To:   canvas.Point{X: 53, Y: 6},
				},
				d: canvas.Rectangle{
					From: canvas.Point{X: 29, Y: 1},
					To:   canvas.Point{X: 39, Y: 1},
				},
				e: canvas.Rectangle{
					From: canvas.Point{X: 15, Y: 1},
					To:   canvas.Point{X: 25, Y: 1},
				},
				f: canvas.Rectangle{
					From: canvas.Point{X: 0, Y: 6},
					To:   canvas.Point{X: 9, Y: 6},
				},
			},
			"tiradas": map[posicion]canvas.Rectangle{
				a: canvas.Rectangle{
					From: canvas.Point{X: 19, Y: 7},
					To:   canvas.Point{X: 24, Y: 9},
				},
				b: canvas.Rectangle{
					From: canvas.Point{X: 29, Y: 7},
					To:   canvas.Point{X: 34, Y: 9},
				},
				c: canvas.Rectangle{
					From: canvas.Point{X: 35, Y: 5},
					To:   canvas.Point{X: 40, Y: 7},
				},
				d: canvas.Rectangle{
					From: canvas.Point{X: 29, Y: 3},
					To:   canvas.Point{X: 34, Y: 5},
				},
				e: canvas.Rectangle{
					From: canvas.Point{X: 19, Y: 3},
					To:   canvas.Point{X: 24, Y: 5},
				},
				f: canvas.Rectangle{
					From: canvas.Point{X: 13, Y: 5},
					To:   canvas.Point{X: 18, Y: 7},
				},
			},
			"tooltips": map[posicion]canvas.Rectangle{
				a: canvas.Rectangle{
					From: canvas.Point{X: 15, Y: 12},
					To:   canvas.Point{X: 24, Y: 12},
				},
				b: canvas.Rectangle{
					From: canvas.Point{X: 29, Y: 12},
					To:   canvas.Point{X: 38, Y: 12},
				},
				c: canvas.Rectangle{
					From: canvas.Point{X: 44, Y: 5},
					To:   canvas.Point{X: 53, Y: 5},
				},
				d: canvas.Rectangle{
					From: canvas.Point{X: 29, Y: 0},
					To:   canvas.Point{X: 38, Y: 0},
				},
				e: canvas.Rectangle{
					From: canvas.Point{X: 15, Y: 0},
					To:   canvas.Point{X: 24, Y: 0},
				},
				f: canvas.Rectangle{
					From: canvas.Point{X: 0, Y: 5},
					To:   canvas.Point{X: 5, Y: 9},
				},
			},
		},
		otrasAreas: map[string]canvas.Rectangle{
			"exteriorMesa": canvas.Rectangle{
				From: canvas.Point{X: 11, Y: 2},
				To:   canvas.Point{X: 42, Y: 10},
			},
			"interiorMesa": canvas.Rectangle{
				From: canvas.Point{X: 12, Y: 3},
				To:   canvas.Point{X: 41, Y: 9},
			},
		},
	}
}
