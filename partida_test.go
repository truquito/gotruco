package truco

import (
	"fmt"
	"testing"
)

func TestParseJugada(t *testing.T) {
	p, _ := NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{ // Alvaro tiene flor
				Cartas: [3]Carta{
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 6},
					Carta{Palo: Basto, Valor: 7},
				},
			},
			Manojo{ // Roro no tiene flor
				Cartas: [3]Carta{
					Carta{Palo: Oro, Valor: 5},
					Carta{Palo: Espada, Valor: 5},
					Carta{Palo: Basto, Valor: 5},
				},
			},
		},
	)

	shouldBeOK := []string{
		"alvaro envido",
		"Alvaro real-envido",
		"Alvaro falta-envido",
		"Alvaro flor",
		"Alvaro contra-flor",
		"Alvaro contra-flor-al-resto",
		"Alvaro truco",
		"Alvaro re-truco",
		"Alvaro vale-4",
		"Alvaro quiero",
		"Alvaro no-quiero",
		"Alvaro mazo",
		// tiradas
		"Alvaro 2 oro",
		"Alvaro 2 ORO",
		"Alvaro 2 oRo",
		"Alvaro 6 basto",
		"Alvaro 7 basto",
		"Roro 5 Oro",
		"Roro 5 Espada",
		"Roro 5 Basto",
	}

	shouldNotBeOK := []string{
		"Juancito envido",
		"Juancito envido asd",
		"Juancito envido 33",
		"Juancito envid0",
		// tiradas
		"Alvaro 2 oroo",
		"Alvaro 2 oRo ",
		"Alvaro 6 espada*",
		"Alvaro 7 asd",
		"Alvaro 2  copa",
		"Alvaro 54 Oro ",
		"Alvaro 0 oro",
		"Alvaro 9 oro",
		"Alvaro 8 oro",
		"Alvaro 111 oro",
		// roro trata de usar las de alvaro
		// esto se debe testear en jugadas
		// "roro 2 oRo",
		// "roro 6 basto",
		// "roro 7 basto",
	}

	for _, cmd := range shouldBeOK {
		_, err := p.parseJugada(cmd)
		ok := err == nil
		if !ok {
			t.Error(err.Error())
		}
	}

	for _, cmd := range shouldNotBeOK {
		_, err := p.parseJugada(cmd)
		if err == nil {
			t.Error(`Deberia dar error`)
		}
	}
}

func TestPartida1(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	p.Ronda.setMuestra(Carta{Palo: Oro, Valor: 3})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{ // Alvaro tiene flor
				Cartas: [3]Carta{
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 6},
					Carta{Palo: Basto, Valor: 7},
				},
			},
			Manojo{ // Roro no tiene flor
				Cartas: [3]Carta{
					Carta{Palo: Oro, Valor: 5},
					Carta{Palo: Espada, Valor: 5},
					Carta{Palo: Basto, Valor: 5},
				},
			},
			Manojo{ // Adolfo tiene flor
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Copa, Valor: 2},
					Carta{Palo: Copa, Valor: 3},
				},
			},
			Manojo{ // Renzo tiene flor
				Cartas: [3]Carta{
					Carta{Palo: Oro, Valor: 4},
					Carta{Palo: Espada, Valor: 4},
					Carta{Palo: Espada, Valor: 1},
				},
			},
			Manojo{ // Andres no tiene  flor
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 10},
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Basto, Valor: 11},
				},
			},
			Manojo{ // Richard tiene flor
				Cartas: [3]Carta{
					Carta{Palo: Oro, Valor: 10},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 1},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Envido") // no estoy recibiendo output
	p.SetSigJugada("Alvaro Flor")
	p.SetSigJugada("Roro Mazo") // no estoy recibiendo output
	p.SetSigJugada("Adolfo Flor")
	p.SetSigJugada("Renzo Contra-flor")
	p.SetSigJugada("Alvaro Quiero")
	p.Terminar()
}

func TestPartidaComandosInvalidos(t *testing.T) {

	p, _ := NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})

	p.SetSigJugada("Alvaro Envido")
	p.SetSigJugada("Quiero")
	p.SetSigJugada("Schumacher Flor")
	p.SetSigJugada("Adolfo Flor")

	p.Terminar()

}

func TestPartidaJSON(t *testing.T) {
	p, _ := NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	fmt.Printf(p.ToJSON())

}

// - 11 le gana a 10 (de la muestra) no de sparda
// - si es parda pero el turno deberia de ser de el mano (alvaro)
// - adolfo deberia de poder cantar retruco
func TestFixNacho(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"jugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"jugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro 6 basto")
	p.SetSigJugada("roro 2 basto")
	p.SetSigJugada("Adolfo 4 basto")
	p.SetSigJugada("renzo 7 basto")
	p.SetSigJugada("andres 10 espada")
	p.SetSigJugada("richard flor")
	p.SetSigJugada("richard 11 espada")

	p.SetSigJugada("richard truco") // el envido deberia pasar a inhabilitado

	p.SetSigJugada("roro quiero") // no deberia poder ya que es de su mismo equipo
	p.SetSigJugada("adolfo quiero")
	p.SetSigJugada("richard 5 espada")
	p.SetSigJugada("alvaro mazo")
	p.SetSigJugada("roro quiero")     // no hay nada que querer
	p.SetSigJugada("roro retruco")    // syntaxis invalida
	p.SetSigJugada("roro re-truco")   // no debe permitir
	p.SetSigJugada("alvaro re-truco") // no deberia dejarlo porque se fue al mazo
	p.SetSigJugada("Adolfo re-truco") // ojo que nadie le acepto el re-truco
	p.SetSigJugada("roro 6 copa")     // no deberia dejarlo porque ya paso su turno
	p.SetSigJugada("adolfo re-truco")
	p.SetSigJugada("adolfo 1 espada")
	p.SetSigJugada("renzo retruco")
	p.SetSigJugada("renzo re-truco")
	p.SetSigJugada("renzo mazo")
	p.SetSigJugada("andres mazo")

	p.Esperar()

	p.Print()

	p.Terminar()
}

func TestFixNoFlor(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"jugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Basto","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":2},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"jugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Basto","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":2},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":4},{"palo":"Espada","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":7},{"palo":"Basto","valor":11},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":12},{"palo":"Basto","valor":1},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Espada","valor":7},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Basto","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":2},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro 4 basto")
	p.SetSigJugada("roro truco")

	p.Esperar()

	oops = !(p.Ronda.Envite.Estado == FLOR)
	if oops {
		t.Error(`El envido esta primero!`)
	}

	p.SetSigJugada("richard no-quiero")

	p.SetSigJugada("roro truco")

	p.SetSigJugada("adolfo 12 oro")
	// No era su turno, no puede tirar la carta

	p.SetSigJugada("roro 7 copa")
	p.SetSigJugada("andres quiero")
	p.SetSigJugada("adolfo 12 oro")
	p.SetSigJugada("renzo 5 oro")

	p.SetSigJugada("andres flor")
	// No es posible cantar flor

	p.SetSigJugada("andres 6 basto")

	p.SetSigJugada("richard flor")
	// No es posible cantar flor

	p.SetSigJugada("richard 11 copa")
	// termina la primera mano, la gana adolfo
	// entonces es el turno de adolfo

	p.SetSigJugada("adolfo re-truco")
	p.SetSigJugada("richard quiero")

	p.SetSigJugada("richard vale-4")

	p.Esperar()

	oops = !(p.Ronda.Truco.Estado == VALE4)
	if oops {
		t.Error(`Richard deberia poder gritar vale4`)
	}

	p.SetSigJugada("adolfo quiero")

	p.Esperar()

	oops = !(p.Ronda.Truco.Estado == VALE4QUERIDO)
	if oops {
		t.Error(`El estado del truco deberia ser VALE4QUERIDO`)
	}

	p.SetSigJugada("adolfo 1 basto")
	p.SetSigJugada("renzo 7 espada")
	p.SetSigJugada("andres 4 espada")
	p.SetSigJugada("richard 10 espada")
	p.SetSigJugada("alvaro 6 espada")
	p.SetSigJugada("roro re-truco")
	p.SetSigJugada("roro mazo")

	p.Esperar()

	oops = !(p.getMaxPuntaje() == 3+4)
	if oops {
		t.Error(`suma mal los puntos cuando roro se fue al mazo`)
	}

	p.Print()
	p.Terminar()
}
