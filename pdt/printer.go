package pdt

import (
	"strconv"
	"strings"

	"github.com/jpfilevich/canvas"
)

func chop(str string, l int) string {
	if len(str) <= l {
		return str
	}
	return str[:l]
}

type posicion int

const (
	a posicion = iota
	b
	c
	d
	e
	f
)

// iPrinter Interface para las impresoras de 2, 4 y 6 jugadores
type iPrinter interface {
	dibujarMarco()
	dibujarEstadisticas(p *PartidaDT)
	dibujarMuestra(muestra Carta)
	dibujarNombres(manojos []Manojo, muestra Carta)
	dibujarTiradas(manojos []Manojo)
	dibujarPosesiones(manojos []Manojo)
	dibujarTooltips(r Ronda)
	// Print(p *PartidaDT)
	render() string
}

// abstracta ya que no implementa el metodo dibujarTooltips
type impresora struct {
	canvas         canvas.Canvas
	areasJugadores map[string](map[posicion]canvas.Rectangle)
	otrasAreas     map[string]canvas.Rectangle
	templates
}

type impresora2 struct {
	impresora
}

type impresora4 struct {
	impresora
}

type impresora6 struct {
	impresora
}

func (pr impresora) render() string {
	render := "\n" + pr.canvas.Render()
	return render
}

func (pr impresora) dibujarMarco() {
	marco := pr.templates.marco()
	pr.canvas.DrawAt(pr.otrasAreas["exteriorMesa"].From, marco)
}

func (pr impresora) dibujarEstadisticas(p *PartidaDT) {
	template := pr.templates.estadisticas()
	pr.canvas.DrawAt(pr.otrasAreas["estadisticas"].From, template)

	// NUMERO DE Mano en juego
	numMano := p.Ronda.ManoEnJuego.String()
	pr.canvas.DrawAt(pr.otrasAreas["#Mano"].From, numMano)

	// Mano
	mano := p.Ronda.GetElMano().Jugador.Nombre
	mano = chop(mano, 8)
	pr.canvas.DrawAt(pr.otrasAreas["Mano"].From, mano)

	// Turno
	turno := p.Ronda.GetElTurno().Jugador.Nombre
	turno = chop(turno, 8)
	pr.canvas.DrawAt(pr.otrasAreas["Turno"].From, turno)

	// puntuacion
	puntuacion := strconv.Itoa((int(p.Puntuacion)))
	pr.canvas.DrawAt(pr.otrasAreas["Puntuacion"].From, puntuacion)

	// ROJO
	ptjRojo := strconv.Itoa((int(p.Puntajes[Rojo])))
	pr.canvas.DrawAt(pr.otrasAreas["puntajeRojo"].From, ptjRojo)

	// AZUL
	ptjAzul := strconv.Itoa((int(p.Puntajes[Azul])))
	pr.canvas.DrawAt(pr.otrasAreas["puntajeAzul"].From, ptjAzul)
}

func (pr impresora) dibujarMuestra(muestra Carta) {
	carta := pr.templates.carta(muestra)
	pr.canvas.DrawAt(pr.otrasAreas["muestra"].From, carta)
}

func (pr impresora) dibujarNombres(manojos []Manojo, muestra Carta) {
	for i, manojo := range manojos {
		nombre := manojo.Jugador.Nombre
		// tieneFlor, _ := manojo.tieneFlor(muestra)
		// if tieneFlor {
		// 	nombre = "❀ " + nombre
		// }
		if len(nombre) > 10 {
			nombre = nombre[:10]
		}
		area := pr.areasJugadores["nombres"][posicion(i)]
		var nombreCentrado string
		if posicion(i) == f {
			nombreCentrado = area.Center(nombre)
		} else if posicion(i) == c {
			nombreCentrado = area.Center(nombre)
		} else {
			nombreCentrado = area.Center(nombre)
		}
		pr.canvas.DrawAt(area.From, nombreCentrado)
	}
}

func (pr impresora) dibujarTiradas(manojos []Manojo) {
	var area canvas.Rectangle

	for i := range manojos {
		area = pr.areasJugadores["tiradas"][posicion(i)]
		// necesito saber cuantas tiro
		manojo := manojos[i]
		cantTiradas := manojo.GetCantCartasTiradas()
		carta := manojo.Cartas[manojo.UltimaTirada]
		var tiradas string
		switch cantTiradas {
		case 1:
			tiradas = pr.templates.carta(*carta)
		case 2:
			tiradas = pr.templates.cartaDobleSolapada(*carta)
		case 3:
			tiradas = pr.templates.cartaTripleSolapada(*carta)
		default:
			tiradas = pr.templates.vacio()
		}
		pr.canvas.DrawAt(area.From, area.Center(tiradas))
	}
}

func lasConoce(cartas []*Carta) bool {
	lasConoce := true
	// si hay al menos una carta con nil -> no las conoce
	for _, c := range cartas {
		if c == nil {
			lasConoce = false
			break
		}
	}
	return lasConoce
}

func (pr impresora) dibujarPosesiones(manojos []Manojo) {
	var area canvas.Rectangle

	for i := range manojos {
		area = pr.areasJugadores["posesiones"][posicion(i)]
		manojo := manojos[i]

		var cartasEnPosesion []*Carta
		for j, c := range manojo.Cartas {
			if manojo.CartasNoTiradas[j] {
				cartasEnPosesion = append(cartasEnPosesion, c)
			}
		}

		cantTiradas := manojo.GetCantCartasTiradas()
		cantPosesion := 3 - cantTiradas

		var template string

		if lasConoce(cartasEnPosesion) {
			switch cantPosesion {
			case 1:
				template = pr.templates.carta(*cartasEnPosesion[0])
			case 2:
				template = pr.templates.cartaDobleVisible(cartasEnPosesion)
			case 3:
				template = pr.templates.cartaTripleVisible(cartasEnPosesion)
			default:
				template = pr.templates.vacio()
			}
		} else {
			switch cantPosesion {
			case 1:
				template = pr.templates.cartaOculta(*cartasEnPosesion[0])
			case 2:
				template = pr.templates.cartaDobleOculta(cartasEnPosesion)
			case 3:
				template = pr.templates.cartaTripleOculta(cartasEnPosesion)
			default:
				template = pr.templates.vacio()
			}
		}
		pr.canvas.DrawAt(area.From, area.Center(template))
	}
}

// dibuja: turno y flor
func (pr impresora) dibujarTooltips(r Ronda) {
	turno := int(r.Turno)

	for i, manojo := range r.Manojos {
		tooltip := ""

		// flor
		if lasConoce(manojo.Cartas[:]) {
			tieneFlor, _ := manojo.TieneFlor(r.Muestra)
			if tieneFlor {
				tooltip += "❀"
			}
		}

		// el turno
		esSuTurno := turno == i
		if esSuTurno {
			posicion := posicion(turno)
			switch posicion {
			case a, b:
				tooltip += " ⭡"
			default:
				tooltip += " ⭣"
			}

		}
		tooltip = strings.Trim(tooltip, " ")
		area := pr.areasJugadores["tooltips"][posicion(i)]
		pr.canvas.DrawAt(area.From, area.Center(tooltip))
	}

}

func renderizar(p *PartidaDT) string {

	// como tiene el parametro en Print
	// basta con tener una sola instancia de impresora
	// para imprimir varias instancias de partidas diferentes
	var pr iPrinter
	switch p.CantJugadores {
	case 2:
		pr = nuevaImpresora2()
	case 4:
		pr = nuevaImpresora4()
	case 6:
		pr = nuevaImpresora6()
	}

	pr.dibujarMarco()
	pr.dibujarEstadisticas(p)
	pr.dibujarMuestra(p.Ronda.Muestra)
	pr.dibujarNombres(p.Ronda.Manojos, p.Ronda.Muestra)
	pr.dibujarTiradas(p.Ronda.Manojos)
	pr.dibujarPosesiones(p.Ronda.Manojos)
	pr.dibujarTooltips(p.Ronda)

	return pr.render()
}

/* overrides */
func (pr impresora2) dibujarTooltips(r Ronda) {
	turno := int(r.Turno)

	for i, manojo := range r.Manojos {
		tooltip := ""
		tieneFlor, _ := manojo.TieneFlor(r.Muestra)
		if tieneFlor {
			tooltip += "❀"
		}
		esSuTurno := turno == i
		if esSuTurno {
			tooltip += " ⭣"
		}
		tooltip = strings.Trim(tooltip, " ")
		area := pr.areasJugadores["tooltips"][posicion(i)]
		pr.canvas.DrawAt(area.From, area.Center(tooltip))
	}

}

func nuevaImpresora2() impresora2 {
	return impresora2{

		impresora{canvas: canvas.NewCanvas(75, 19),
			areasJugadores: map[string](map[posicion]canvas.Rectangle){
				"nombres": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 44, Y: 9},
						To:   canvas.Point{X: 53, Y: 9},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 0, Y: 9},
						To:   canvas.Point{X: 9, Y: 9},
					},
				},
				"tiradas": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 35, Y: 8},
						To:   canvas.Point{X: 40, Y: 13},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 13, Y: 8},
						To:   canvas.Point{X: 18, Y: 10},
					},
				},
				"posesiones": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 44, Y: 10},
						To:   canvas.Point{X: 53, Y: 12},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 0, Y: 10},
						To:   canvas.Point{X: 9, Y: 12},
					},
				},
				"tooltips": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 44, Y: 8},
						To:   canvas.Point{X: 53, Y: 8},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 0, Y: 8},
						To:   canvas.Point{X: 5, Y: 12},
					},
				},
			},
			otrasAreas: map[string]canvas.Rectangle{
				"muestra": canvas.Rectangle{
					From: canvas.Point{X: 25, Y: 8},
					To:   canvas.Point{X: 28, Y: 10},
				},
				"exteriorMesa": canvas.Rectangle{
					From: canvas.Point{X: 11, Y: 5},
					To:   canvas.Point{X: 42, Y: 13},
				},
				"interiorMesa": canvas.Rectangle{
					From: canvas.Point{X: 12, Y: 6},
					To:   canvas.Point{X: 41, Y: 12},
				},
				"estadisticas": canvas.Rectangle{
					From: canvas.Point{X: 57, Y: 2},
					To:   canvas.Point{X: 74, Y: 15},
				},
				"#Mano": canvas.Rectangle{
					From: canvas.Point{X: 66, Y: 3},
					To:   canvas.Point{X: 72, Y: 3},
				},
				"Mano": canvas.Rectangle{
					From: canvas.Point{X: 65, Y: 5},
					To:   canvas.Point{X: 72, Y: 5},
				},
				"Turno": canvas.Rectangle{
					From: canvas.Point{X: 66, Y: 7},
					To:   canvas.Point{X: 72, Y: 7},
				},
				"Puntuacion": canvas.Rectangle{
					From: canvas.Point{X: 71, Y: 9},
					To:   canvas.Point{X: 72, Y: 9},
				},
				"puntajeRojo": canvas.Rectangle{
					From: canvas.Point{X: 61, Y: 14},
					To:   canvas.Point{X: 62, Y: 14},
				},
				"puntajeAzul": canvas.Rectangle{
					From: canvas.Point{X: 68, Y: 14},
					To:   canvas.Point{X: 69, Y: 14},
				},
			},
		},
	}
}

func nuevaImpresora4() impresora4 {
	return impresora4{
		impresora{
			canvas: canvas.NewCanvas(75, 19),
			areasJugadores: map[string](map[posicion]canvas.Rectangle){
				"nombres": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 14},
						To:   canvas.Point{X: 24, Y: 14},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 14},
						To:   canvas.Point{X: 38, Y: 14},
					},

					c: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 4},
						To:   canvas.Point{X: 38, Y: 4},
					},
					d: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 4},
						To:   canvas.Point{X: 24, Y: 4},
					},
				},
				"tiradas": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 19, Y: 10},
						To:   canvas.Point{X: 24, Y: 12},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 10},
						To:   canvas.Point{X: 34, Y: 12},
					},

					c: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 6},
						To:   canvas.Point{X: 34, Y: 8},
					},
					d: canvas.Rectangle{
						From: canvas.Point{X: 19, Y: 6},
						To:   canvas.Point{X: 24, Y: 8},
					},
				},
				"posesiones": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 16},
						To:   canvas.Point{X: 24, Y: 18},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 16},
						To:   canvas.Point{X: 38, Y: 18},
					},

					c: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 0},
						To:   canvas.Point{X: 38, Y: 2},
					},
					d: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 0},
						To:   canvas.Point{X: 24, Y: 2},
					},
				},
				"tooltips": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 15},
						To:   canvas.Point{X: 24, Y: 15},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 15},
						To:   canvas.Point{X: 38, Y: 15},
					},

					c: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 3},
						To:   canvas.Point{X: 38, Y: 3},
					},
					d: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 3},
						To:   canvas.Point{X: 24, Y: 3},
					},
				},
			},
			otrasAreas: map[string]canvas.Rectangle{
				"muestra": canvas.Rectangle{
					From: canvas.Point{X: 25, Y: 8},
					To:   canvas.Point{X: 28, Y: 10},
				},
				"exteriorMesa": canvas.Rectangle{
					From: canvas.Point{X: 11, Y: 5},
					To:   canvas.Point{X: 42, Y: 13},
				},
				"interiorMesa": canvas.Rectangle{
					From: canvas.Point{X: 12, Y: 6},
					To:   canvas.Point{X: 41, Y: 12},
				},
				"estadisticas": canvas.Rectangle{
					From: canvas.Point{X: 57, Y: 2},
					To:   canvas.Point{X: 74, Y: 15},
				},
				"#Mano": canvas.Rectangle{
					From: canvas.Point{X: 66, Y: 3},
					To:   canvas.Point{X: 72, Y: 3},
				},
				"Mano": canvas.Rectangle{
					From: canvas.Point{X: 65, Y: 5},
					To:   canvas.Point{X: 72, Y: 5},
				},
				"Turno": canvas.Rectangle{
					From: canvas.Point{X: 66, Y: 7},
					To:   canvas.Point{X: 72, Y: 7},
				},
				"Puntuacion": canvas.Rectangle{
					From: canvas.Point{X: 71, Y: 9},
					To:   canvas.Point{X: 72, Y: 9},
				},
				"puntajeRojo": canvas.Rectangle{
					From: canvas.Point{X: 61, Y: 14},
					To:   canvas.Point{X: 62, Y: 14},
				},
				"puntajeAzul": canvas.Rectangle{
					From: canvas.Point{X: 68, Y: 14},
					To:   canvas.Point{X: 69, Y: 14},
				},
			},
		},
	}
}

func nuevaImpresora6() impresora6 {
	return impresora6{
		impresora{
			canvas: canvas.NewCanvas(75, 19),
			areasJugadores: map[string](map[posicion]canvas.Rectangle){
				"nombres": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 14},
						To:   canvas.Point{X: 24, Y: 14},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 14},
						To:   canvas.Point{X: 38, Y: 14},
					},
					c: canvas.Rectangle{
						From: canvas.Point{X: 44, Y: 9},
						To:   canvas.Point{X: 53, Y: 9},
					},
					d: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 4},
						To:   canvas.Point{X: 38, Y: 4},
					},
					e: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 4},
						To:   canvas.Point{X: 24, Y: 4},
					},
					f: canvas.Rectangle{
						From: canvas.Point{X: 0, Y: 9},
						To:   canvas.Point{X: 9, Y: 9},
					},
				},
				"tiradas": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 19, Y: 10},
						To:   canvas.Point{X: 24, Y: 12},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 10},
						To:   canvas.Point{X: 34, Y: 12},
					},
					c: canvas.Rectangle{
						From: canvas.Point{X: 35, Y: 8},
						To:   canvas.Point{X: 40, Y: 13},
					},
					d: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 6},
						To:   canvas.Point{X: 34, Y: 8},
					},
					e: canvas.Rectangle{
						From: canvas.Point{X: 19, Y: 6},
						To:   canvas.Point{X: 24, Y: 8},
					},
					f: canvas.Rectangle{
						From: canvas.Point{X: 13, Y: 8},
						To:   canvas.Point{X: 18, Y: 10},
					},
				},
				"posesiones": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 16},
						To:   canvas.Point{X: 24, Y: 18},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 16},
						To:   canvas.Point{X: 38, Y: 18},
					},
					c: canvas.Rectangle{
						From: canvas.Point{X: 44, Y: 10},
						To:   canvas.Point{X: 53, Y: 12},
					},
					d: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 0},
						To:   canvas.Point{X: 38, Y: 2},
					},
					e: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 0},
						To:   canvas.Point{X: 24, Y: 2},
					},
					f: canvas.Rectangle{
						From: canvas.Point{X: 0, Y: 10},
						To:   canvas.Point{X: 9, Y: 12},
					},
				},
				"tooltips": map[posicion]canvas.Rectangle{
					a: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 15},
						To:   canvas.Point{X: 24, Y: 15},
					},
					b: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 15},
						To:   canvas.Point{X: 38, Y: 15},
					},
					c: canvas.Rectangle{
						From: canvas.Point{X: 44, Y: 8},
						To:   canvas.Point{X: 53, Y: 8},
					},
					d: canvas.Rectangle{
						From: canvas.Point{X: 29, Y: 3},
						To:   canvas.Point{X: 38, Y: 3},
					},
					e: canvas.Rectangle{
						From: canvas.Point{X: 15, Y: 3},
						To:   canvas.Point{X: 24, Y: 3},
					},
					f: canvas.Rectangle{
						From: canvas.Point{X: 0, Y: 8},
						To:   canvas.Point{X: 5, Y: 12},
					},
				},
			},
			otrasAreas: map[string]canvas.Rectangle{
				"muestra": canvas.Rectangle{
					From: canvas.Point{X: 25, Y: 8},
					To:   canvas.Point{X: 28, Y: 10},
				},
				"exteriorMesa": canvas.Rectangle{
					From: canvas.Point{X: 11, Y: 5},
					To:   canvas.Point{X: 42, Y: 13},
				},
				"interiorMesa": canvas.Rectangle{
					From: canvas.Point{X: 12, Y: 6},
					To:   canvas.Point{X: 41, Y: 12},
				},
				"estadisticas": canvas.Rectangle{
					From: canvas.Point{X: 57, Y: 2},
					To:   canvas.Point{X: 74, Y: 15},
				},
				"#Mano": canvas.Rectangle{
					From: canvas.Point{X: 66, Y: 3},
					To:   canvas.Point{X: 72, Y: 3},
				},
				"Mano": canvas.Rectangle{
					From: canvas.Point{X: 65, Y: 5},
					To:   canvas.Point{X: 72, Y: 5},
				},
				"Turno": canvas.Rectangle{
					From: canvas.Point{X: 66, Y: 7},
					To:   canvas.Point{X: 72, Y: 7},
				},
				"Puntuacion": canvas.Rectangle{
					From: canvas.Point{X: 71, Y: 9},
					To:   canvas.Point{X: 72, Y: 9},
				},
				"puntajeRojo": canvas.Rectangle{
					From: canvas.Point{X: 61, Y: 14},
					To:   canvas.Point{X: 62, Y: 14},
				},
				"puntajeAzul": canvas.Rectangle{
					From: canvas.Point{X: 68, Y: 14},
					To:   canvas.Point{X: 69, Y: 14},
				},
			},
		},
	}
}
