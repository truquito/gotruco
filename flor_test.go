package truco

import (
	"testing"

	"github.com/jpfilevich/truco/pdt"
)

func TestNoDeberianTenerFlor(t *testing.T) {

	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Copa, Valor: 5})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 10},
					&pdt.Carta{Palo: pdt.Copa, Valor: 7},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	tieneFlor, _ := p.Ronda.Manojos[0].TieneFlor(p.Ronda.Muestra)
	oops = tieneFlor == true
	if oops {
		t.Error(`Alvaro' NO deberia de tener 'flor'`)
		return
	}

	tieneFlor, _ = p.Ronda.Manojos[1].TieneFlor(p.Ronda.Muestra)
	oops = tieneFlor == true
	if oops {
		t.Error(`Roro' NO deberia de tener 'flor'`)
		return
	}

}

func TestDeberiaTenerFlor(t *testing.T) {

	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Copa, Valor: 5})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 4},
					&pdt.Carta{Palo: pdt.Espada, Valor: 10},
					&pdt.Carta{Palo: pdt.Espada, Valor: 7},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Oro, Valor: 3},
				},
			},
		},
	)

	tieneFlor, _ := p.Ronda.Manojos[0].TieneFlor(p.Ronda.Muestra)
	oops = !(tieneFlor == true)
	if oops {
		t.Error(`Alvaro' deberia tener 'flor'`)
		return
	}

	tieneFlor, _ = p.Ronda.Manojos[1].TieneFlor(p.Ronda.Muestra)
	oops = !(tieneFlor == true)
	if oops {
		t.Error(`Roro' deberia tener 'flor'`)
		return
	}
}

func TestFlorFlorContraFlorQuiero(t *testing.T) {

	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // Alvaro tiene flor
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 6},
					&pdt.Carta{Palo: pdt.Basto, Valor: 7},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // Roro
					&pdt.Carta{Palo: pdt.Oro, Valor: 5},
					&pdt.Carta{Palo: pdt.Espada, Valor: 5},
					&pdt.Carta{Palo: pdt.Basto, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // Adolfo tiene flor
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Copa, Valor: 2},
					&pdt.Carta{Palo: pdt.Copa, Valor: 3},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // Renzo tiene flor
					&pdt.Carta{Palo: pdt.Oro, Valor: 4},
					&pdt.Carta{Palo: pdt.Espada, Valor: 4},
					&pdt.Carta{Palo: pdt.Espada, Valor: 1},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // Andres
					&pdt.Carta{Palo: pdt.Copa, Valor: 10},
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Basto, Valor: 11},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // Richard tiene flor
					&pdt.Carta{Palo: pdt.Oro, Valor: 10},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 1},
				},
			},
		},
	)

	p.Cmd("Alvaro Flor")
	p.Cmd("Roro Mazo")
	p.Cmd("Renzo Flor")
	p.Cmd("Adolfo Contra-flor-al-resto")
	p.Cmd("Richard Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia ser 'deshabilitado'`)
		return
	}

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado de la flor deberia ser 'deshabilitado'`)
		return
	}

	// duda: se suman solo las flores ganadoras
	// si contraflor AL RESTO -> no acumulativo
	// duda: deberia sumar tambien los puntos de las flores
	// oops = !(p.Puntajes[pdt.Azul] == 4*3+10 && p.Puntajes[pdt.Rojo] == 0)
	// puntos para ganar chico + todas las flores NO ACHICADAS
	oops = !(p.Puntajes[pdt.Azul] == 10 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia ser 2`)
		return
	}
}

// Tests:
// los "me achico" no cuentan para la flor
// Flor		xcg(+3) / xcg(+3)
// Flor + Contra-Flor		xc(+3) / xCadaFlorDelQueHizoElDesafio(+3) + 1
// Flor + [Contra-Flor] + ContraFlorAlResto		~Falta Envido + *TODAS* las flores no achicadas / xcg(+3) + 1

func TestFixFlor(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":12},{"palo":"Oro","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Basto","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Espada","valor":10},{"palo":"Basto","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Copa","valor":3},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":1},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro envido")
	// pero Richard tiene flor
	// y no le esta sumando esos puntos

	if !(p.Ronda.Envite.Estado == pdt.DESHABILITADO) {
		t.Error(`El estado de la flor deberia ser 'deshabilitado'`)
	} else if !(p.Puntajes[pdt.Rojo] == 3) {
		t.Error(`El puntaje del equipo rojo deberia ser 3 por la flor de richard`)
	}

	p.Cmd("alvaro 6 espada")
	p.Cmd("alvaro 6 espada")
	p.Cmd("roro 5 espada")
	p.Cmd("adolfo 10 oro")
	p.Cmd("renzo 6 basto")
	p.Cmd("andres 6 copa")
	p.Cmd("richard 3 espada")
	p.Cmd("adolfo 10 copa")
	p.Cmd("renzo 10 espada")
	p.Cmd("andres 3 copa")
	p.Cmd("richard 11 espada")
	p.Cmd("alvaro 12 basto")
	p.Cmd("roro 10 basto")

	if !(p.Puntajes[pdt.Rojo] == 3) {
		t.Error(`El puntaje del equipo rojo deberia ser 3 por la flor de richard`)
	} else if !(p.Puntajes[pdt.Azul] == 1) {
		t.Error(`El puntaje del equipo azul deberia ser 1 por la ronda ganada`)
	}

}

// bug a arreglar:
// hay 2 flores; se cantan ambas -> no pasa nada
func TestFixFlorBucle(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":12},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":12},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":11},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Oro","valor":5},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Basto","valor":1},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":2},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":12},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":10},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro mazo")
	p.Cmd("roro flor")
	p.Cmd("richard flor")

	if !(p.Puntajes[pdt.Rojo] == 6) {
		t.Error(`El puntaje del equipo rojo deberia ser 6 por las 2 flores`)
	}

}

// bug a arreglar:
// no se puede cantar contra flor
func TestFixContraFlor(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":1},{"palo":"Espada","valor":1},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Copa","valor":7},{"palo":"Oro","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Oro","valor":1},{"palo":"Oro","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Copa","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro 1 basto")
	p.Cmd("roro 4 oro")
	p.Cmd("adolfo flor")
	p.Cmd("adolfo 11 basto")
	p.Cmd("renzo 12 basto")
	p.Cmd("renzo quiero")
	p.Cmd("renzo contra-flor")
	p.Cmd("adolfo quiero")

	// if !(p.Puntajes[pdt.Rojo] == 3) {
	// 	t.Error(`El puntaje del equipo rojo deberia ser 3 por la flor de richard`)
	// } else if !(p.Puntajes[pdt.Azul] == 1) {
	// 	t.Error(`El puntaje del equipo azul deberia ser 1 por la ronda ganada`)
	// }

}
