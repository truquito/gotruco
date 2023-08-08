package pdt

import (
	"strconv"
	"strings"

	"github.com/filevich/canvas"
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
	dibujarEstadisticas(p *Partida)
	dibujarMuestra(muestra Carta)
	dibujarNombres(manojos []Manojo, muestra Carta)
	dibujarTiradas(manojos []Manojo)
	dibujarPosesiones(manojos []Manojo)
	dibujarTooltips(r Ronda)
	dibujarDialogos(r Ronda, dialogos ...Dialogo)
	// Print(p *Partida)
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

func (pr impresora) dibujarEstadisticas(p *Partida) {
	template := pr.templates.estadisticas()
	pr.canvas.DrawAt(pr.otrasAreas["estadisticas"].From, template)

	// NUMERO DE Mano en juego
	numMano := p.Ronda.ManoEnJuego.String()
	pr.canvas.DrawAt(pr.otrasAreas["#Mano"].From, numMano)

	// Mano
	mano := p.Ronda.GetElMano().Jugador.ID
	mano = chop(mano, 8)
	pr.canvas.DrawAt(pr.otrasAreas["Mano"].From, mano)

	// Turno
	turno := p.Ronda.GetElTurno().Jugador.ID
	turno = chop(turno, 8)
	pr.canvas.DrawAt(pr.otrasAreas["Turno"].From, turno)

	// puntuacion
	puntuacion := strconv.Itoa((int(p.Puntuacion)))
	pr.canvas.DrawAt(pr.otrasAreas["Puntuacion"].From, puntuacion)

	// ------------------

	// Envite
	envite := p.Ronda.Envite.Estado.String()
	envite = chop(envite, 6)
	pr.canvas.DrawAt(pr.otrasAreas["Envite"].From, envite)

	// Envite Autor
	envite_por := p.Ronda.Envite.CantadoPor
	envite_por = chop(envite_por, 6)
	pr.canvas.DrawAt(pr.otrasAreas["EnvitePor"].From, envite_por)

	// Truco
	truco := p.Ronda.Truco.Estado.String()
	truco = chop(truco, 8)
	pr.canvas.DrawAt(pr.otrasAreas["Truco"].From, truco)

	// Truco Autor
	truco_por := p.Ronda.Truco.CantadoPor
	truco_por = chop(truco_por, 9)
	pr.canvas.DrawAt(pr.otrasAreas["TrucoPor"].From, truco_por)

	// ------------------

	// nombres
	var nombreR, nombreA string
	if p.Ronda.Manojos[0].Jugador.Equipo == Rojo {
		nombreR = p.Ronda.Manojos[0].Jugador.ID
		nombreA = p.Ronda.Manojos[1].Jugador.ID
	} else {
		nombreR = p.Ronda.Manojos[1].Jugador.ID
		nombreA = p.Ronda.Manojos[0].Jugador.ID
	}
	nombreR, nombreA = chop(nombreR, 5), chop(nombreA, 5)
	pr.canvas.DrawAt(pr.otrasAreas["nombreRojo"].From, nombreR)
	pr.canvas.DrawAt(pr.otrasAreas["nombreAzul"].From, nombreA)

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
		nombre := manojo.Jugador.ID
		// tieneFlor, _ := manojo.tieneFlor(muestra)
		// if tieneFlor {
		// 	nombre = "❀ " + nombre
		// }
		nombre = chop(nombre, 10)
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
		var carta *Carta = nil
		if manojo.UltimaTirada >= 0 {
			carta = manojo.Cartas[manojo.UltimaTirada]
		}
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
		if manojos[i].SeFueAlMazo {
			continue
		}
		area = pr.areasJugadores["posesiones"][posicion(i)]
		manojo := manojos[i]

		var cartasEnPosesion []*Carta
		for j, c := range manojo.Cartas {
			if !manojo.Tiradas[j] {
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
				template = pr.templates.cartaOculta()
			case 2:
				template = pr.templates.cartaDobleOculta()
			case 3:
				template = pr.templates.cartaTripleOculta()
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

		if manojo.SeFueAlMazo {
			tooltip += "✗ "
		}

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
				tooltip += " ↑"
			default:
				tooltip += " ↓"
			}

		}
		tooltip = strings.Trim(tooltip, " ")
		area := pr.areasJugadores["tooltips"][posicion(i)]
		pr.canvas.DrawAt(area.From, area.Center(tooltip))
	}

}

func (pr impresora) dibujarDialogos(r Ronda, dialogos ...Dialogo) {
	for _, d := range dialogos {
		var pos int
		for i, m := range r.Manojos {
			if m.Jugador.ID == d.ID {
				pos = i
				break
			}
		}
		// d.Msg
		up := pos == 0 || pos == 1
		m := pr.templates.dialogo(up)

		area := pr.areasJugadores["dialogos"][posicion(pos)]
		pr.canvas.DrawAt(area.From, m)

		m = chop(d.Msg, max_len_msg)
		area = pr.areasJugadores["msgs"][posicion(pos)]
		pr.canvas.DrawAt(area.From, area.Center(m))
	}
}

const max_len_msg = 13

type Dialogo struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

// Renderizar .
func Renderizar(p *Partida, dialogos ...Dialogo) string {

	// como tiene el parametro en Print
	// basta con tener una sola instancia de impresora
	// para imprimir varias instancias de partidas diferentes
	var pr iPrinter

	switch len(p.Ronda.Manojos) {
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
	pr.dibujarDialogos(p.Ronda, dialogos...)

	return pr.render()
}

/* overrides */
func (pr impresora2) dibujarTooltips(r Ronda) {
	turno := int(r.Turno)

	for i, manojo := range r.Manojos {

		tooltip := ""

		if lasConoce(manojo.Cartas[:]) {
			tieneFlor, _ := manojo.TieneFlor(r.Muestra)
			if tieneFlor {
				tooltip += "❀"
			}
		}

		esSuTurno := turno == i
		if esSuTurno {
			tooltip += " ↓"
		}
		tooltip = strings.Trim(tooltip, " ")
		area := pr.areasJugadores["tooltips"][posicion(i)]
		pr.canvas.DrawAt(area.From, area.Center(tooltip))
	}

}

func (pr impresora2) dibujarDialogos(r Ronda, dialogos ...Dialogo) {
	dialogo := pr.templates.dialogo(false)

	for _, d := range dialogos {
		// que usuario es? ~ que pos le corresponde?
		var pos int
		for i, m := range r.Manojos {
			if m.Jugador.ID == d.ID {
				pos = i
				break
			}
		}
		area := pr.areasJugadores["dialogos"][posicion(pos)]
		pr.canvas.DrawAt(area.From, dialogo)

		area = pr.areasJugadores["msgs"][posicion(pos)]
		m := chop(d.Msg, max_len_msg)
		pr.canvas.DrawAt(area.From, area.Center(m))
	}
}

func nuevaImpresora2() impresora2 {
	im6 := nuevaImpresora6()

	mask6 := func(attr string) map[posicion]canvas.Rectangle {
		return map[posicion]canvas.Rectangle{
			a: im6.areasJugadores[attr][c],
			b: im6.areasJugadores[attr][f],
		}
	}

	return impresora2{
		impresora{
			canvas: canvas.NewCanvas(80, 25),
			areasJugadores: map[string](map[posicion]canvas.Rectangle){
				"nombres":    mask6("nombres"),
				"tiradas":    mask6("tiradas"),
				"posesiones": mask6("posesiones"),
				"tooltips":   mask6("tooltips"),
				"dialogos":   mask6("dialogos"),
				"msgs":       mask6("msgs"),
			},
			otrasAreas: im6.otrasAreas,
		},
	}
}

func nuevaImpresora4() impresora4 {
	im6 := nuevaImpresora6()

	mask6 := func(attr string) map[posicion]canvas.Rectangle {
		return map[posicion]canvas.Rectangle{
			a: im6.areasJugadores[attr][a],
			b: im6.areasJugadores[attr][b],
			c: im6.areasJugadores[attr][d],
			d: im6.areasJugadores[attr][e],
		}
	}

	return impresora4{
		impresora{
			canvas: canvas.NewCanvas(80, 25),
			areasJugadores: map[string](map[posicion]canvas.Rectangle){
				"nombres":    mask6("nombres"),
				"tiradas":    mask6("tiradas"),
				"posesiones": mask6("posesiones"),
				"tooltips":   mask6("tooltips"),
				"dialogos":   mask6("dialogos"),
				"msgs":       mask6("msgs"),
			},
			otrasAreas: im6.otrasAreas,
		},
	}
}

func nuevaImpresora6() impresora6 {
	return impresora6{
		impresora{
			canvas: canvas.NewCanvas(80, 25),
			areasJugadores: map[string](map[posicion]canvas.Rectangle){
				"nombres": map[posicion]canvas.Rectangle{
					a: {
						From: canvas.Point{X: 19, Y: 17},
						To:   canvas.Point{X: 28, Y: 17},
					},
					b: {
						From: canvas.Point{X: 33, Y: 17},
						To:   canvas.Point{X: 42, Y: 17},
					},
					c: {
						From: canvas.Point{X: 50, Y: 12},
						To:   canvas.Point{X: 59, Y: 12},
					},
					d: {
						From: canvas.Point{X: 33, Y: 7},
						To:   canvas.Point{X: 42, Y: 7},
					},
					e: {
						From: canvas.Point{X: 19, Y: 7},
						To:   canvas.Point{X: 28, Y: 7},
					},
					f: {
						From: canvas.Point{X: 2, Y: 12},
						To:   canvas.Point{X: 11, Y: 12},
					},
				},
				"tiradas": map[posicion]canvas.Rectangle{
					a: {
						From: canvas.Point{X: 23, Y: 13},
						To:   canvas.Point{X: 28, Y: 15},
					},
					b: {
						From: canvas.Point{X: 33, Y: 13},
						To:   canvas.Point{X: 38, Y: 15},
					},
					c: {
						From: canvas.Point{X: 39, Y: 11},
						To:   canvas.Point{X: 44, Y: 16},
					},
					d: {
						From: canvas.Point{X: 33, Y: 9},
						To:   canvas.Point{X: 38, Y: 11},
					},
					e: {
						From: canvas.Point{X: 23, Y: 9},
						To:   canvas.Point{X: 28, Y: 11},
					},
					f: {
						From: canvas.Point{X: 17, Y: 11},
						To:   canvas.Point{X: 22, Y: 13},
					},
				},
				"posesiones": map[posicion]canvas.Rectangle{
					a: {
						From: canvas.Point{X: 19, Y: 19},
						To:   canvas.Point{X: 28, Y: 21},
					},
					b: {
						From: canvas.Point{X: 33, Y: 19},
						To:   canvas.Point{X: 42, Y: 21},
					},
					c: {
						From: canvas.Point{X: 50, Y: 13},
						To:   canvas.Point{X: 59, Y: 15},
					},
					d: {
						From: canvas.Point{X: 33, Y: 3},
						To:   canvas.Point{X: 42, Y: 5},
					},
					e: {
						From: canvas.Point{X: 19, Y: 3},
						To:   canvas.Point{X: 28, Y: 5},
					},
					f: {
						From: canvas.Point{X: 2, Y: 13},
						To:   canvas.Point{X: 11, Y: 15},
					},
				},
				"tooltips": map[posicion]canvas.Rectangle{
					a: {
						From: canvas.Point{X: 19, Y: 18},
						To:   canvas.Point{X: 28, Y: 18},
					},
					b: {
						From: canvas.Point{X: 33, Y: 18},
						To:   canvas.Point{X: 42, Y: 18},
					},
					c: {
						From: canvas.Point{X: 48, Y: 11},
						To:   canvas.Point{X: 57, Y: 11},
					},
					d: {
						From: canvas.Point{X: 33, Y: 6},
						To:   canvas.Point{X: 42, Y: 6},
					},
					e: {
						From: canvas.Point{X: 19, Y: 6},
						To:   canvas.Point{X: 28, Y: 6},
					},
					f: {
						From: canvas.Point{X: 4, Y: 11},
						To:   canvas.Point{X: 9, Y: 15},
					},
				},
				"dialogos": map[posicion]canvas.Rectangle{
					a: {
						From: canvas.Point{X: 16, Y: 22},
						To:   canvas.Point{X: 30, Y: 25},
					},
					b: {
						From: canvas.Point{X: 31, Y: 22},
						To:   canvas.Point{X: 45, Y: 25},
					},
					c: {
						From: canvas.Point{X: 47, Y: 8},
						To:   canvas.Point{X: 56, Y: 10},
					},
					d: {
						From: canvas.Point{X: 31, Y: 0},
						To:   canvas.Point{X: 45, Y: 3},
					},
					e: {
						From: canvas.Point{X: 16, Y: 0},
						To:   canvas.Point{X: 30, Y: 3},
					},
					f: {
						From: canvas.Point{X: 0, Y: 8},
						To:   canvas.Point{X: 14, Y: 10},
					},
				},
				"msgs": map[posicion]canvas.Rectangle{
					a: {
						From: canvas.Point{X: 17, Y: 23},
						To:   canvas.Point{X: 29, Y: 23},
					},
					b: {
						From: canvas.Point{X: 32, Y: 23},
						To:   canvas.Point{X: 44, Y: 23},
					},
					c: {
						From: canvas.Point{X: 48, Y: 9},
						To:   canvas.Point{X: 60, Y: 9},
					},
					d: {
						From: canvas.Point{X: 32, Y: 1},
						To:   canvas.Point{X: 44, Y: 1},
					},
					e: {
						From: canvas.Point{X: 17, Y: 1},
						To:   canvas.Point{X: 29, Y: 1},
					},
					f: {
						From: canvas.Point{X: 1, Y: 9},
						To:   canvas.Point{X: 13, Y: 9},
					},
				},
			},
			otrasAreas: map[string]canvas.Rectangle{
				"muestra": {
					From: canvas.Point{X: 29, Y: 11},
					To:   canvas.Point{X: 32, Y: 13},
				},
				"exteriorMesa": {
					From: canvas.Point{X: 15, Y: 8},
					To:   canvas.Point{X: 46, Y: 16},
				},
				"interiorMesa": {
					From: canvas.Point{X: 16, Y: 9},
					To:   canvas.Point{X: 45, Y: 15},
				},
				"estadisticas": {
					From: canvas.Point{X: 62, Y: 2},
					To:   canvas.Point{X: 79, Y: 20},
				},
				// -----------------------
				"#Mano": {
					From: canvas.Point{X: 71, Y: 3},
					To:   canvas.Point{X: 77, Y: 3},
				},
				"Mano": {
					From: canvas.Point{X: 70, Y: 5},
					To:   canvas.Point{X: 77, Y: 5},
				},
				"Turno": {
					From: canvas.Point{X: 71, Y: 7},
					To:   canvas.Point{X: 77, Y: 7},
				},
				"Puntuacion": {
					From: canvas.Point{X: 76, Y: 9},
					To:   canvas.Point{X: 77, Y: 9},
				},
				// -----------------------
				"Envite": {
					From: canvas.Point{X: 72, Y: 12},
					To:   canvas.Point{X: 77, Y: 12},
				},
				"EnvitePor": {
					From: canvas.Point{X: 69, Y: 13},
					To:   canvas.Point{X: 77, Y: 13},
				},
				"Truco": {
					From: canvas.Point{X: 71, Y: 15},
					To:   canvas.Point{X: 77, Y: 15},
				},
				"TrucoPor": {
					From: canvas.Point{X: 69, Y: 16},
					To:   canvas.Point{X: 77, Y: 16},
				},
				// -----------------------
				"nombreAzul": {
					From: canvas.Point{X: 64, Y: 19},
					To:   canvas.Point{X: 77, Y: 19}, // tam 5
				},
				"nombreRojo": {
					From: canvas.Point{X: 73, Y: 19},
					To:   canvas.Point{X: 77, Y: 19},
				},
				"puntajeAzul": {
					From: canvas.Point{X: 66, Y: 21},
					To:   canvas.Point{X: 67, Y: 21},
				},
				"puntajeRojo": {
					From: canvas.Point{X: 75, Y: 21},
					To:   canvas.Point{X: 76, Y: 21},
				},
			},
		},
	}
}
