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
	// << Alvaro tira la carta 4 de Basto

	p.SetSigJugada("roro truco")
	// No es posible cantar truco ahora
	// "la flor esta primero":
	// << Andres canta flor

	p.Esperar()

	oops = !(p.Ronda.Envite.Estado == FLOR)
	if oops {
		t.Error(`El envido esta primero!`)
	}

	// el otro que tiene flor, pero se arruga
	p.SetSigJugada("richard no-quiero")
	// FIX: NO ESTA OUTPUTEANDO EL NO QUIERO
	// << +6 puntos para el equipo Azul por las flores

	p.SetSigJugada("roro truco")
	// << Roro grita truco

	p.SetSigJugada("adolfo 12 oro")
	// No era su turno, no puede tirar la carta

	p.SetSigJugada("roro 7 copa")
	// << Roro tira la carta 7 de Copa

	p.SetSigJugada("andres quiero")
	// << Andres responde quiero

	p.SetSigJugada("adolfo 12 oro")
	// << Adolfo tira la carta 12 de Oro

	p.SetSigJugada("renzo 5 oro")
	// << Renzo tira la carta 5 de Oro

	p.SetSigJugada("andres flor")
	// No es posible cantar flor

	p.SetSigJugada("andres 6 basto")
	// << Andres tira la carta 6 de Basto

	p.SetSigJugada("richard flor")
	// No es posible cantar flor

	p.SetSigJugada("richard 11 copa")
	// << Richard tira la carta 11 de Copa

	/* *********************************** */
	// << La primera mano la gano Adolfo (equipo Azul)
	// << Es el turno de Adolfo
	/* *********************************** */

	p.SetSigJugada("adolfo re-truco")
	// << Adolfo grita re-truco

	p.SetSigJugada("richard quiero")
	// << Richard responde quiero

	p.SetSigJugada("richard vale-4")
	// << Richard grita vale 4

	p.Esperar()

	oops = !(p.Ronda.Truco.Estado == VALE4)
	if oops {
		t.Error(`Richard deberia poder gritar vale4`)
	}

	p.SetSigJugada("adolfo quiero")
	// << Adolfo responde quiero

	p.Esperar()

	oops = !(p.Ronda.Truco.Estado == VALE4QUERIDO)
	if oops {
		t.Error(`El estado del truco deberia ser VALE4QUERIDO`)
	}

	/* *********************************** */
	// ACA EMPIEZAN A TIRAR CARTAS PARA LA SEGUNDA MANO
	// muesta: 3 espada
	/* *********************************** */

	p.SetSigJugada("adolfo 1 basto")
	// << Adolfo tira la carta 1 de Basto

	p.SetSigJugada("renzo 7 espada")
	// << Renzo tira la carta 7 de Espada

	p.SetSigJugada("andres 4 espada")
	// << Andres tira la carta 4 de Espada

	p.SetSigJugada("richard 10 espada")
	// << Richard tira la carta 10 de Espada

	p.SetSigJugada("alvaro 6 espada")
	// << Alvaro tira la carta 6 de Espada

	p.SetSigJugada("roro re-truco")
	// << Alvaro tira la carta 6 de Espada

	p.SetSigJugada("roro mazo")
	// << Roro se va al mazo

	// era el ultimo que quedaba por tirar en esta mano
	// -> que evalue la mano

	// << +4 puntos para el equipo Azul por el vale4Querido no querido por Roro
	// << Empieza una nueva ronda
	// << Empieza una nueva ronda

	p.Esperar()

	oops = !(p.getMaxPuntaje() == 6+4) // 6 de las 2 flores
	if oops {
		t.Error(`suma mal los puntos cuando roro se fue al mazo`)
	}

	p.Print()
	p.Terminar()
}

func TestFixPanic(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"jugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"jugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro 6 basto")
	// << Alvaro tira la carta 6 de Basto

	p.SetSigJugada("roro 2 basto")
	// << Roro tira la carta 2 de Basto

	p.SetSigJugada("Adolfo 4 basto")
	// << Adolfo tira la carta 4 de Basto

	p.SetSigJugada("renzo 7 basto")
	// << Renzo tira la carta 7 de Basto

	p.SetSigJugada("andres 10 espada")
	// << Andres tira la carta 10 de Espada

	p.SetSigJugada("richard flor")
	// << Richard canta flor
	// << +3 puntos para el equipo Rojo (por ser la unica flor de esta ronda)

	p.SetSigJugada("richard 11 espada")
	// << Richard tira la carta 11 de Espada

	/*
		// << La Mano resulta parda
		// << Es el turno de Richard
	*/

	// ERROR: no la deberia ganar andres porque es mano (DUDA)

	p.SetSigJugada("richard truco")
	// << Richard grita truco

	p.SetSigJugada("roro quiero")
	// (Para Roro) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.SetSigJugada("adolfo quiero")
	// << Adolfo responde quiero

	p.SetSigJugada("richard 5 espada")
	// << Richard tira la carta 5 de Espada

	p.SetSigJugada("alvaro mazo")
	// << Alvaro se va al mazo

	p.SetSigJugada("roro quiero")
	// (Para Roro) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.SetSigJugada("roro retruco")
	// << No esxiste esa jugada

	p.SetSigJugada("roro re-truco")
	// No es posible cantar re-truco ahora

	p.SetSigJugada("alvaro re-truco") // ya que se fue al mazo
	// No es posible cantar re-truco ahora

	p.SetSigJugada("Adolfo re-truco") // no es su turno ni el de su equipo
	// No es posible cantar re-truco ahora

	p.SetSigJugada("roro 6 copa")
	// << Roro tira la carta 6 de Copa

	p.SetSigJugada("adolfo re-truco")
	// << Adolfo grita re-truco

	p.SetSigJugada("adolfo 1 espada")
	// << Adolfo tira la carta 1 de Espada

	p.SetSigJugada("renzo retruco")
	// << No esxiste esa jugada

	p.SetSigJugada("renzo re-truco") // ya que ya lo canto adolfo
	// No es posible cantar re-truco ahora

	p.SetSigJugada("renzo mazo")
	// << Renzo se va al mazo

	p.Esperar()
	p.Print()

	p.SetSigJugada("andres mazo")
	// << Andres se va al mazo

	// andres se va al mazo y era el ultimo que quedaba por jugar
	// ademas era la segunda mano -> ya se decide
	// aunque hay un retruco propuesto <--------------
	// si hay algo propuesto por su equipo no se puede ir <-------

	// << La segunda mano la gano el equipo Rojo gracia a Richard
	// << La ronda ha sido ganada por el equipo Rojo
	// << +0 puntos para el equipo Rojo por el reTruco no querido
	// << Empieza una nueva ronda

	p.Esperar()
	p.Terminar()
}

func TestFixBocha(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro mazo")
	// << Alvaro se va al mazo

	p.SetSigJugada("adolfo mazo")
	// << Adolfo se va al mazo

	p.SetSigJugada("andres mazo")
	// << Andres se va al mazo

	/*
		- todos los de azul se fueron al mazo ->
		la deberia de haber ganado los rojos
		- deberia ser el turno de Roro (ponele ???)
	*/

	// << La ronda ha sido ganada por el equipo Rojo
	// << +1 puntos para el equipo Rojo por el noCantado ganado
	// << Empieza una nueva ronda

	p.Esperar()
	p.Print()
	p.Terminar()
}

func TestFixBochaParte2(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("roro envido")
	// No es posible cantar 'Envido'

	p.SetSigJugada("andres quiero")
	// (Para Andres) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.SetSigJugada("andres quiero")
	// (Para Andres) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.SetSigJugada("alvaro mazo")
	// << Alvaro se va al mazo

	p.SetSigJugada("adolfo 1 copa")
	// Esa carta no se encuentra en este manojo

	p.SetSigJugada("adolfo mazo")
	// << Adolfo se va al mazo

	p.SetSigJugada("andres mazo")
	// << Andres se va al mazo

	// << La ronda ha sido ganada por el equipo Rojo
	// << +1 puntos para el equipo Rojo por el noCantado ganado
	// << Empieza una nueva ronda

	p.Esperar()
	p.Print()
	p.Terminar()
}

func TestFixBochaParte3(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("richard flor")
	// No es posible cantar flor

	p.SetSigJugada("andres quiero")
	// (Para Andres) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.SetSigJugada("andres contra-flor")
	// No es posible cantar contra flor

	p.SetSigJugada("richard contra-flor")
	// No es posible cantar contra flor

	p.SetSigJugada("richard quiero")
	// (Para Richard) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Esperar()
	p.Print()
	p.Terminar()
}

func TestFixAutoQuerer(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro envido")
	p.SetSigJugada("alvaro quiero")
	p.SetSigJugada("adolfo quiero")

	p.Print()
	p.Terminar()
}

func TestFixNilPointer(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo": "Oro", "valor":11}, {"palo": "Espada", "valor":10}, {"palo": "Basto", "valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo": "Oro", "valor":12}, {"palo": "Copa", "valor":5}, {"palo": "Copa", "valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo": "Espada", "valor":3}, {"palo": "Copa", "valor":7}, {"palo": "Basto", "valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo": "Basto", "valor":6}, {"palo": "Basto", "valor":1}, {"palo": "Copa", "valor":4 }],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo": "Oro", "valor":3}, {"palo": "Copa", "valor":6}, {"palo": "Copa", "valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo": "Espada", "valor":4}, {"palo": "Basto", "valor":10}, {"palo": "Copa", "valor":10 }],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Ronda.Turno = 3 // es el turno de renzo
	p.Print()

	p.SetSigJugada("renzo6 basto")
	p.SetSigJugada("renzo 6 basto")
	p.SetSigJugada("andres truco")
	p.SetSigJugada("renzo quiero")
	p.SetSigJugada("andres re-truco")
	p.SetSigJugada("andres 3 oro")
	p.SetSigJugada("richard vale-4")
	p.SetSigJugada("richard re-truco")
	p.SetSigJugada("andres quiero")
	p.SetSigJugada("richard mazo")
	p.SetSigJugada("alvaro vale-4")
	p.SetSigJugada("andres quiero")
	p.SetSigJugada("roro quiero")
	p.SetSigJugada("alvaro mazo")
	p.SetSigJugada("roro mazo")
	p.SetSigJugada("roro 12 oro")
	p.SetSigJugada("adolfo mazo")
	p.SetSigJugada("Renzo flor")

	p.Esperar()
	// output := p.Dispatch()
	// for _, msg := range output {
	// 	fmt.Println(msg)
	// }

	// p.Print()
	p.Terminar()
}

func TestFixNoDejaIrseAlMazo(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"jugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Espada","valor":10},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"jugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Espada","valor":10},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Oro","valor":7},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Copa","valor":2},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Oro","valor":4},{"palo":"Oro","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Espada","valor":10},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Copa","valor":7},{"palo":"Basto","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":2},{"palo":"Copa","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro 11 copa")
	p.SetSigJugada("roro 6 basto")
	p.SetSigJugada("adolfo 12 espada")
	p.SetSigJugada("renzo 2 espada")
	p.SetSigJugada("renzo flor")
	p.SetSigJugada("andres 6 oro")
	p.SetSigJugada("richard truco")
	p.SetSigJugada("andres quiero")
	p.SetSigJugada("alvaro 7 oro")
	p.SetSigJugada("roro 2 copa")
	p.SetSigJugada("richard 2 oro")
	p.SetSigJugada("renzo mazo")
	p.SetSigJugada("andres mazo")
	p.SetSigJugada("andres mazo")

	p.Esperar()
	// output := p.Dispatch()
	// for _, msg := range output {
	// 	fmt.Println(msg)
	// }

	p.Print()
	p.Terminar()
}

func TestFixFlorObligatoria(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"jugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}],"jugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":7},{"palo":"Basto","valor":5},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":1},{"palo":"Copa","valor":11},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Basto","valor":3},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro 2 basto") // alvaro deberia primero cantar flor
	// p.SetSigJugada("alvaro flor")
	p.SetSigJugada("roro 5 copa")
	p.SetSigJugada("adolfo 7 espada")
	p.SetSigJugada("renzo 1 espada")
	p.SetSigJugada("andres 10 espada")
	p.SetSigJugada("richard 3 basto")
	p.SetSigJugada("alvaro envido")
	p.SetSigJugada("alvaro 1 oro")
	p.SetSigJugada("roro 2 espada")
	p.SetSigJugada("adolfo truco")
	p.SetSigJugada("roro quiero")
	p.SetSigJugada("renzo quiero")
	p.SetSigJugada("adolfo 5 basto")
	p.SetSigJugada("renzo quiero")
	p.SetSigJugada("renzo 11 copa")
	p.SetSigJugada("andres 2 oro")
	p.SetSigJugada("richard 10 oro")
	p.SetSigJugada("roro 1 oro")

	p.Esperar()
	// output := p.Dispatch()
	// for _, msg := range output {
	// 	fmt.Println(msg)
	// }

	p.Print()
	p.Terminar()
}

func TestFixNoPermiteContraFlor(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":1},{"palo":"Espada","valor":1},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Copa","valor":7},{"palo":"Oro","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Oro","valor":1},{"palo":"Oro","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Copa","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	p.Print()

	p.SetSigJugada("alvaro 1 basto")
	p.SetSigJugada("roro 4 oro")
	p.SetSigJugada("adolfo flor")
	p.SetSigJugada("adolfo 11 basto")
	p.SetSigJugada("renzo 12 basto")
	p.SetSigJugada("renzo quiero")
	p.SetSigJugada("renzo contra-flor")

	p.Esperar()
	// output := p.Dispatch()
	// for _, msg := range output {
	// 	fmt.Println(msg)
	// }

	p.Print()
	p.Terminar()
}
