package truco

import (
	"testing"
)

func TestNoDeberianTenerFlor(t *testing.T) {

	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Copa, Valor: 5})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Copa, Valor: 10},
					Carta{Palo: Copa, Valor: 7},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 3},
				},
			},
		},
	)

	tieneFlor, _ := p.Ronda.Manojos[0].tieneFlor(p.Ronda.Muestra)
	oops = tieneFlor == true
	if oops {
		t.Error(`Alvaro' NO deberia de tener 'flor'`)
		return
	}

	tieneFlor, _ = p.Ronda.Manojos[1].tieneFlor(p.Ronda.Muestra)
	oops = tieneFlor == true
	if oops {
		t.Error(`Roro' NO deberia de tener 'flor'`)
		return
	}

}

func TestDeberiaTenerFlor(t *testing.T) {

	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Copa, Valor: 5})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 4},
					Carta{Palo: Espada, Valor: 10},
					Carta{Palo: Espada, Valor: 7},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Oro, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Oro, Valor: 3},
				},
			},
		},
	)

	tieneFlor, _ := p.Ronda.Manojos[0].tieneFlor(p.Ronda.Muestra)
	oops = !(tieneFlor == true)
	if oops {
		t.Error(`Alvaro' deberia tener 'flor'`)
		return
	}

	tieneFlor, _ = p.Ronda.Manojos[1].tieneFlor(p.Ronda.Muestra)
	oops = !(tieneFlor == true)
	if oops {
		t.Error(`Roro' deberia tener 'flor'`)
		return
	}
}

func TestFlorFlorContraFlorQuiero(t *testing.T) {

	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	p.Ronda.setMuestra(Carta{Palo: Oro, Valor: 3})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // Alvaro tiene flor
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 6},
					Carta{Palo: Basto, Valor: 7},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // Roro
					Carta{Palo: Oro, Valor: 5},
					Carta{Palo: Espada, Valor: 5},
					Carta{Palo: Basto, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // Adolfo tiene flor
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Copa, Valor: 2},
					Carta{Palo: Copa, Valor: 3},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // Renzo tiene flor
					Carta{Palo: Oro, Valor: 4},
					Carta{Palo: Espada, Valor: 4},
					Carta{Palo: Espada, Valor: 1},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // Andres
					Carta{Palo: Copa, Valor: 10},
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Basto, Valor: 11},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // Richard tiene flor
					Carta{Palo: Oro, Valor: 10},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 1},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Flor")
	p.SetSigJugada("Roro Mazo")
	p.SetSigJugada("Renzo Flor")
	p.SetSigJugada("Adolfo Contra-flor-al-resto")
	p.SetSigJugada("Richard Quiero")
	p.Esperar()

	oops = p.Ronda.Envite.Estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia ser 'deshabilitado'`)
		return
	}

	oops = p.Ronda.Envite.Estado != DESHABILITADO
	if oops {
		t.Error(`El estado de la flor deberia ser 'deshabilitado'`)
		return
	}

	// duda: se suman solo las flores ganadoras
	// si contraflor AL RESTO -> no acumulativo
	// duda: deberia sumar tambien los puntos de las flores
	// oops = !(p.Puntajes[Azul] == 4*3+10 && p.Puntajes[Rojo] == 0)
	// puntos para ganar chico + todas las flores NO ACHICADAS
	oops = !(p.Puntajes[Azul] == 10 && p.Puntajes[Rojo] == 0)
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
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"jugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"jugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":12},{"palo":"Oro","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Basto","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Espada","valor":10},{"palo":"Basto","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Copa","valor":3},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":1},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro envido")
	// pero Richard tiene flor
	// y no le esta sumando esos puntos
	p.Esperar()

	if !(p.Ronda.Envite.Estado == DESHABILITADO) {
		t.Error(`El estado de la flor deberia ser 'deshabilitado'`)
	} else if !(p.Puntajes[Rojo] == 3) {
		t.Error(`El puntaje del equipo rojo deberia ser 3 por la flor de richard`)
	}

	p.SetSigJugada("alvaro 6 espada")
	p.SetSigJugada("alvaro 6 espada")
	p.SetSigJugada("roro 5 espada")
	p.SetSigJugada("adolfo 10 oro")
	p.SetSigJugada("renzo 6 basto")
	p.SetSigJugada("andres 6 copa")
	p.SetSigJugada("richard 3 espada")
	p.SetSigJugada("adolfo 10 copa")
	p.SetSigJugada("renzo 10 espada")
	p.SetSigJugada("andres 3 copa")
	p.SetSigJugada("richard 11 espada")
	p.SetSigJugada("alvaro 12 basto")
	p.SetSigJugada("roro 10 basto")
	p.Esperar()

	if !(p.Puntajes[Rojo] == 3) {
		t.Error(`El puntaje del equipo rojo deberia ser 3 por la flor de richard`)
	} else if !(p.Puntajes[Azul] == 1) {
		t.Error(`El puntaje del equipo azul deberia ser 1 por la ronda ganada`)
	}

	p.Terminar()
}

// bug a arreglar:
// hay 2 flores; se cantan ambas -> no pasa nada
func TestFixFlorBucle(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"jugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":12},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"jugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":12},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":11},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Oro","valor":5},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Basto","valor":1},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":2},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":12},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":10},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro mazo")
	p.SetSigJugada("roro flor")
	p.SetSigJugada("richard flor")
	p.Esperar()

	if !(p.Puntajes[Rojo] == 6) {
		t.Error(`El puntaje del equipo rojo deberia ser 6 por las 2 flores`)
	}

	p.Terminar()
}

// bug a arreglar:
// no se puede cantar contra flor
func TestFixContraFlor(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"jugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"jugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":1},{"palo":"Espada","valor":1},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Copa","valor":7},{"palo":"Oro","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Oro","valor":1},{"palo":"Oro","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Copa","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro 1 basto")
	p.SetSigJugada("roro 4 oro")
	p.SetSigJugada("adolfo flor")
	p.SetSigJugada("adolfo 11 basto")
	p.SetSigJugada("renzo 12 basto")
	p.SetSigJugada("renzo quiero")
	p.SetSigJugada("renzo contra-flor")
	p.SetSigJugada("adolfo quiero")

	p.Esperar()

	// if !(p.Puntajes[Rojo] == 3) {
	// 	t.Error(`El puntaje del equipo rojo deberia ser 3 por la flor de richard`)
	// } else if !(p.Puntajes[Azul] == 1) {
	// 	t.Error(`El puntaje del equipo azul deberia ser 1 por la ronda ganada`)
	// }

	p.Terminar()
}
