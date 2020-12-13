package truco

import (
	"fmt"
	"testing"

	"github.com/filevich/truco/deco"

	"github.com/filevich/truco/out"
	"github.com/filevich/truco/pdt"
)

func TestParseJugada(t *testing.T) {
	p, _ := NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 6},
					&pdt.Carta{Palo: pdt.Basto, Valor: 7},
				},
			},
			pdt.Manojo{ // Roro no tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 5},
					&pdt.Carta{Palo: pdt.Espada, Valor: 5},
					&pdt.Carta{Palo: pdt.Basto, Valor: 5},
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
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 6},
					&pdt.Carta{Palo: pdt.Basto, Valor: 7},
				},
			},
			pdt.Manojo{ // Roro no tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 5},
					&pdt.Carta{Palo: pdt.Espada, Valor: 5},
					&pdt.Carta{Palo: pdt.Basto, Valor: 5},
				},
			},
			pdt.Manojo{ // Adolfo tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Copa, Valor: 2},
					&pdt.Carta{Palo: pdt.Copa, Valor: 3},
				},
			},
			pdt.Manojo{ // Renzo tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 4},
					&pdt.Carta{Palo: pdt.Espada, Valor: 4},
					&pdt.Carta{Palo: pdt.Espada, Valor: 1},
				},
			},
			pdt.Manojo{ // Andres no tiene  flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 10},
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Basto, Valor: 11},
				},
			},
			pdt.Manojo{ // Richard tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 10},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 1},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido") // no estoy recibiendo output
	p.Cmd("Alvaro Flor")
	p.Cmd("Roro Mazo") // no estoy recibiendo output
	p.Cmd("Adolfo Flor")
	p.Cmd("Renzo Contra-flor")
	p.Cmd("Alvaro Quiero")

}

func TestPartidaComandosInvalidos(t *testing.T) {

	p, _ := NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})

	p.Cmd("Alvaro Envido")
	p.Cmd("Quiero")
	p.Cmd("Schumacher Flor")
	p.Cmd("Adolfo Flor")

}

func TestPartidaJSON(t *testing.T) {
	p, _ := NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	pJSON, err := p.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf(string(pJSON))
}

// - 11 le gana a 10 (de la muestra) no de sparda
// - si es parda pero el turno deberia de ser de el mano (alvaro)
// - adolfo deberia de poder cantar retruco
func TestFixNacho(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)

	p.Cmd("alvaro 6 basto")
	p.Cmd("roro 2 basto")
	p.Cmd("Adolfo 4 basto")
	p.Cmd("renzo 7 basto")
	p.Cmd("andres 10 espada")
	p.Cmd("richard flor")
	p.Cmd("richard 11 espada")
	p.Cmd("richard truco") // el envido deberia estar inhabilitado por la flor, no por esta accion
	p.Cmd("roro quiero")   // no deberia poder ya que es de su mismo equipo
	p.Cmd("adolfo quiero")
	p.Cmd("richard 5 espada")
	p.Cmd("alvaro mazo")
	p.Cmd("roro quiero")     // no hay nada que querer
	p.Cmd("roro retruco")    // syntaxis invalida
	p.Cmd("roro re-truco")   // no debe permitir
	p.Cmd("alvaro re-truco") // no deberia dejarlo porque se fue al mazo
	p.Cmd("Adolfo re-truco") // ojo que nadie le acepto el re-truco
	p.Cmd("renzo quiero")
	p.Cmd("roro 6 copa")     // no deberia dejarlo porque ya paso su turno
	p.Cmd("adolfo re-truco") // no deberia dejarlo
	p.Cmd("adolfo 1 espada")
	p.Cmd("renzo 3 oro") // no deberia de dejarlo porque el equipo contrario	// propuso un re-truco

	oops = !(p.Ronda.GetElTurno().Jugador.Nombre == "Andres")
	if oops {
		t.Error(`Deberia ser el turno de Andres`)
		return
	}

	p.Cmd("andres mazo")

	out.Consume(p.Stdout, out.Print)
	p.Print()

}

func TestFixNoFlor(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Basto","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":2},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Basto","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":2},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":4},{"palo":"Espada","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":7},{"palo":"Basto","valor":11},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":12},{"palo":"Basto","valor":1},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Espada","valor":7},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Basto","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":2},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro 4 basto")
	// << Alvaro tira la carta 4 de pdt.Basto

	p.Cmd("roro truco")
	// No es posible cantar truco ahora
	// "la flor esta primero":
	// << Andres canta flor

	oops = !(p.Ronda.Envite.Estado == pdt.FLOR)
	if oops {
		t.Error(`El envido esta primero!`)
	}

	// el otro que tiene flor, pero se arruga
	p.Cmd("richard no-quiero")
	// FIX: NO ESTA OUTPUTEANDO EL NO QUIERO
	// << +6 puntos para el equipo pdt.Azul por las flores

	p.Cmd("roro truco")
	// << Roro grita truco

	p.Cmd("adolfo 12 oro")
	// No era su turno, no puede tirar la carta

	p.Cmd("roro 7 copa")
	// << Roro tira la carta 7 de pdt.Copa

	p.Cmd("andres quiero")
	// << Andres responde quiero

	p.Cmd("adolfo 12 oro")
	// << Adolfo tira la carta 12 de pdt.Oro

	p.Cmd("renzo 5 oro")
	// << Renzo tira la carta 5 de pdt.Oro

	p.Cmd("andres flor")
	// No es posible cantar flor

	p.Cmd("andres 6 basto")
	// << Andres tira la carta 6 de pdt.Basto

	p.Cmd("richard flor")
	// No es posible cantar flor

	p.Cmd("richard 11 copa")
	// << Richard tira la carta 11 de pdt.Copa

	/* *********************************** */
	// << La Primera mano la gano Adolfo (equipo pdt.Azul)
	// << Es el turno de Adolfo
	/* *********************************** */

	p.Cmd("adolfo re-truco")
	// << Adolfo grita re-truco

	out.Consume(p.Stdout, deco.Print(&p.PartidaDT))
	p.Print()

	p.Cmd("richard quiero")
	// << Richard responde quiero

	p.Cmd("richard vale-4")
	// << Richard grita vale 4

	oops = !(p.Ronda.Truco.Estado == pdt.VALE4)
	if oops {
		t.Error(`Richard deberia poder gritar vale4`)
	}

	p.Cmd("adolfo quiero")
	// << Adolfo responde quiero

	oops = !(p.Ronda.Truco.Estado == pdt.VALE4QUERIDO)
	if oops {
		t.Error(`El estado del truco deberia ser VALE4QUERIDO`)
	}

	/* *********************************** */
	// ACA EMPIEZAN A TIRAR CARTAS PARA LA SEGUNDA MANO
	// muesta: 3 espada
	/* *********************************** */

	p.Cmd("adolfo 1 basto")
	// << Adolfo tira la carta 1 de pdt.Basto

	p.Cmd("renzo 7 espada")
	// << Renzo tira la carta 7 de pdt.Espada

	p.Cmd("andres 4 espada")
	// << Andres tira la carta 4 de pdt.Espada

	p.Cmd("richard 10 espada")
	// << Richard tira la carta 10 de pdt.Espada

	p.Cmd("alvaro 6 espada")
	// << Alvaro tira la carta 6 de pdt.Espada

	p.Cmd("roro re-truco")
	// << Alvaro tira la carta 6 de pdt.Espada

	p.Cmd("roro mazo")
	// << Roro se va al mazo

	// era el ultimo que quedaba por tirar en esta mano
	// -> que evalue la mano

	// << +4 puntos para el equipo pdt.Azul por el vale4Querido no querido por Roro
	// << Empieza una nueva ronda
	// << Empieza una nueva ronda

	oops = !(p.GetMaxPuntaje() == 6+4) // 6 de las 2 flores
	if oops {
		t.Error(`suma mal los puntos cuando roro se fue al mazo`)
	}

	p.Print()

}

func TestFixPanic(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro 6 basto")
	// << Alvaro tira la carta 6 de pdt.Basto

	p.Cmd("roro 2 basto")
	// << Roro tira la carta 2 de pdt.Basto

	p.Cmd("Adolfo 4 basto")
	// << Adolfo tira la carta 4 de pdt.Basto

	p.Cmd("renzo 7 basto")
	// << Renzo tira la carta 7 de pdt.Basto

	p.Cmd("andres 10 espada")
	// << Andres tira la carta 10 de pdt.Espada

	p.Cmd("richard flor")
	// << Richard canta flor
	// << +3 puntos para el equipo pdt.Rojo (por ser la unica flor de esta ronda)

	p.Cmd("richard 11 espada")
	// << Richard tira la carta 11 de pdt.Espada

	/*
		// << La Mano resulta parda
		// << Es el turno de Richard
	*/

	// ERROR: no la deberia ganar andres porque es mano (DUDA)

	p.Cmd("richard truco")
	// << Richard grita truco

	p.Cmd("roro quiero")
	// (Para Roro) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Cmd("adolfo quiero")
	// << Adolfo responde quiero

	p.Cmd("richard 5 espada")
	// << Richard tira la carta 5 de pdt.Espada

	p.Cmd("alvaro mazo")
	// << Alvaro se va al mazo

	p.Cmd("roro quiero")
	// (Para Roro) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Cmd("roro retruco")
	// << No esxiste esa jugada

	p.Cmd("roro re-truco")
	// No es posible cantar re-truco ahora

	p.Cmd("alvaro re-truco") // ya que se fue al mazo
	// No es posible cantar re-truco ahora

	p.Cmd("Adolfo re-truco") // no es su turno ni el de su equipo
	// No es posible cantar re-truco ahora

	p.Cmd("roro 6 copa")
	// << Roro tira la carta 6 de pdt.Copa

	p.Cmd("adolfo re-truco")
	// << Adolfo grita re-truco

	p.Cmd("adolfo 1 espada")
	// << Adolfo tira la carta 1 de pdt.Espada

	p.Cmd("renzo retruco")
	// << No esxiste esa jugada

	p.Cmd("renzo re-truco") // ya que ya lo canto adolfo
	// No es posible cantar re-truco ahora

	p.Cmd("renzo mazo")
	// << Renzo se va al mazo

	p.Print()

	p.Cmd("andres mazo")
	// << Andres se va al mazo

	// andres se va al mazo y era el ultimo que quedaba por jugar
	// ademas era la Segunda mano -> ya se decide
	// aunque hay un retruco propuesto <--------------
	// si hay algo propuesto por su equipo no se puede ir <-------

	// << La Segunda mano la gano el equipo pdt.Rojo gracia a Richard
	// << La ronda ha sido ganada por el equipo pdt.Rojo
	// << +0 puntos para el equipo pdt.Rojo por el reTruco no querido
	// << Empieza una nueva ronda

}

func TestFixBocha(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro mazo")
	// << Alvaro se va al mazo

	p.Cmd("adolfo mazo")
	// << Adolfo se va al mazo

	p.Cmd("andres mazo")
	// << Andres se va al mazo

	/*
		- todos los de azul se fueron al mazo ->
		la deberia de haber ganado los rojos
		- deberia ser el turno de Roro (ponele ???)
	*/

	// << La ronda ha sido ganada por el equipo pdt.Rojo
	// << +1 puntos para el equipo pdt.Rojo por el noCantado ganado
	// << Empieza una nueva ronda

	p.Print()

}

func TestFixBochaParte2(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("roro envido")
	// No es posible cantar 'Envido'

	p.Cmd("andres quiero")
	// (Para Andres) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Cmd("andres quiero")
	// (Para Andres) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Cmd("alvaro mazo")
	// << Alvaro se va al mazo

	p.Cmd("adolfo 1 copa")
	// Esa carta no se encuentra en este manojo

	p.Cmd("adolfo mazo")
	// << Adolfo se va al mazo

	p.Cmd("andres mazo")
	// << Andres se va al mazo

	// << La ronda ha sido ganada por el equipo pdt.Rojo
	// << +1 puntos para el equipo pdt.Rojo por el noCantado ganado
	// << Empieza una nueva ronda

	p.Print()

}

func TestFixBochaParte3(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("richard flor")
	// No es posible cantar flor

	p.Cmd("andres quiero")
	// (Para Andres) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Cmd("andres contra-flor")
	// No es posible cantar contra flor

	p.Cmd("richard contra-flor")
	// No es posible cantar contra flor

	p.Cmd("richard quiero")
	// (Para Richard) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Print()

}

func TestFixAutoQuerer(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro envido")
	p.Cmd("alvaro quiero")
	p.Cmd("adolfo quiero")

	p.Print()

}

func TestFixNilPointer(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo": "Oro", "valor":11}, {"palo": "Espada", "valor":10}, {"palo": "Basto", "valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo": "Oro", "valor":12}, {"palo": "Copa", "valor":5}, {"palo": "Copa", "valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo": "Espada", "valor":3}, {"palo": "Copa", "valor":7}, {"palo": "Basto", "valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo": "Basto", "valor":6}, {"palo": "Basto", "valor":1}, {"palo": "Copa", "valor":4 }],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo": "Oro", "valor":3}, {"palo": "Copa", "valor":6}, {"palo": "Copa", "valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo": "Espada", "valor":4}, {"palo": "Basto", "valor":10}, {"palo": "Copa", "valor":10 }],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Ronda.Turno = 3 // es el turno de renzo
	p.Print()

	p.Cmd("renzo6 basto")
	p.Cmd("renzo 6 basto")
	p.Cmd("andres truco")
	p.Cmd("renzo quiero")
	p.Cmd("andres re-truco")
	p.Cmd("andres 3 oro")
	p.Cmd("richard vale-4")
	p.Cmd("richard re-truco")
	p.Cmd("andres quiero")
	p.Cmd("richard mazo")
	p.Cmd("alvaro vale-4")
	p.Cmd("andres quiero")
	p.Cmd("roro quiero")
	p.Cmd("alvaro mazo")
	p.Cmd("roro mazo")
	p.Cmd("roro 12 oro")
	p.Cmd("adolfo mazo")
	p.Cmd("Renzo flor")

	// output := p.Dispatch()
	// for _, msg := range output {
	// 	fmt.Println(msg)
	// }

	// p.Print()

}

func TestFixNoDejaIrseAlMazo(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Espada","valor":10},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Espada","valor":10},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Oro","valor":7},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Copa","valor":2},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Oro","valor":4},{"palo":"Oro","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Espada","valor":10},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Copa","valor":7},{"palo":"Basto","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":2},{"palo":"Copa","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro 11 copa")
	p.Cmd("roro 6 basto")
	p.Cmd("adolfo 12 espada")
	p.Cmd("renzo 2 espada")
	p.Cmd("renzo flor")
	p.Cmd("andres 6 oro")
	p.Cmd("richard truco")
	p.Cmd("andres quiero")
	p.Cmd("alvaro 7 oro")
	p.Cmd("roro 2 copa")
	p.Cmd("richard 2 oro")
	p.Cmd("renzo mazo")
	p.Cmd("andres mazo")
	p.Cmd("andres mazo")

	// output := p.Dispatch()
	// for _, msg := range output {
	// 	fmt.Println(msg)
	// }

	p.Print()

}

func TestFixFlorObligatoria(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":7},{"palo":"Basto","valor":5},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":1},{"palo":"Copa","valor":11},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Basto","valor":3},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro 2 basto") // alvaro deberia primero cantar flor
	// p.Cmd("alvaro flor")
	p.Cmd("roro 5 copa")
	p.Cmd("adolfo 7 espada")
	p.Cmd("renzo 1 espada")
	p.Cmd("andres 10 espada")
	p.Cmd("richard 3 basto")
	p.Cmd("alvaro envido")
	p.Cmd("alvaro 1 oro")
	p.Cmd("roro 2 espada")
	p.Cmd("adolfo truco")
	p.Cmd("roro quiero")
	p.Cmd("renzo quiero")
	p.Cmd("adolfo 5 basto")
	p.Cmd("renzo quiero")
	p.Cmd("renzo 11 copa")
	p.Cmd("andres 2 oro")
	p.Cmd("richard 10 oro")
	p.Cmd("roro 1 oro")

	// output := p.Dispatch()
	// for _, msg := range output {
	// 	fmt.Println(msg)
	// }

	p.Print()

}

func TestFixNoPermiteContraFlor(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":1},{"palo":"Espada","valor":1},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Copa","valor":7},{"palo":"Oro","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Oro","valor":1},{"palo":"Oro","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Copa","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro 1 basto")
	p.Cmd("roro 4 oro")
	p.Cmd("adolfo flor")
	p.Cmd("adolfo 11 basto")
	p.Cmd("renzo 12 basto")
	p.Cmd("renzo quiero")
	p.Cmd("renzo contra-flor")

	// output := p.Dispatch()
	// for _, msg := range output {
	// 	fmt.Println(msg)
	// }

	p.Print()

}

func TestFixDeberiaGanarAzul(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":null,"JugadoresConFlorQueNoCantaron":[]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":10},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Oro","valor":5},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Espada","valor":3},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":4},{"palo":"Espada","valor":6},{"palo":"Copa","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":12},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro 10 oro")
	p.Cmd("roro 10 copa")
	p.Cmd("adolfo 2 copa")
	p.Cmd("renzo 4 basto")
	p.Cmd("renzo mazo")
	p.Cmd("alvaro mazo")
	p.Cmd("roro mazo")

	oops = !(p.Puntajes[pdt.Rojo] == 0 && p.Puntajes[pdt.Azul] > 0)
	if oops {
		t.Error(`La ronda deberia de haber sido ganado por pdt.Azul`)
		return
	}

}

func TestFixPierdeTurno(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":7},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Espada","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Oro","valor":2},{"palo":"Espada","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro 5 espada")
	p.Cmd("adolfo mazo")

	out.Consume(p.Stdout, out.Print)
	p.Print()

	oops = !(p.Ronda.GetElTurno().Jugador.Nombre == "Roro")
	if oops {
		t.Error(`Deberia ser el turno de Roro`)
		return
	}

}

func TestFixTieneFlor(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":7},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Espada","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Oro","valor":2},{"palo":"Espada","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro 5 copa")
	// no deberia dejarlo jugar porque tiene flor

	oops = !(p.Ronda.GetElTurno().Jugador.Nombre == "Alvaro")
	if oops {
		t.Error(`Deberia ser el turno de Alvaro`)
		return
	}

}

func Test2FloresSeVaAlMazo(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro 1 copa") // no deberia dejarlo porque tiene flor
	p.Cmd("alvaro envido") // no deberia dejarlo porque tiene flor
	p.Cmd("alvaro flor")
	p.Cmd("richard mazo") // lo deja que se vaya

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestTodoTienenFlor(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 6},
					&pdt.Carta{Palo: pdt.Basto, Valor: 7},
				},
			},
			pdt.Manojo{ // Roro no tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 5},
					&pdt.Carta{Palo: pdt.Espada, Valor: 5},
					&pdt.Carta{Palo: pdt.Basto, Valor: 5},
				},
			},
			pdt.Manojo{ // Adolfo tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Copa, Valor: 2},
					&pdt.Carta{Palo: pdt.Copa, Valor: 3},
				},
			},
			pdt.Manojo{ // Renzo no tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 10},
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Basto, Valor: 11},
				},
			},
			pdt.Manojo{ // Andres tiene  flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 4},
					&pdt.Carta{Palo: pdt.Espada, Valor: 4},
					&pdt.Carta{Palo: pdt.Espada, Valor: 1},
				},
			},
			pdt.Manojo{ // Richard no tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 11},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 1},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido") // no estoy recibiendo output
	p.Cmd("Alvaro Flor")
	p.Cmd("Roro Mazo") // no estoy recibiendo output
	p.Cmd("Adolfo Flor")

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestFixTopeEnvido(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":10,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":null,"JugadoresConFlorQueNoCantaron":[]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":1},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":11},{"palo":"Basto","valor":3},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":1},{"palo":"Espada","valor":10},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Copa","valor":12},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":7},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":1},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	// azul va 10 pts de 20,
	// asi que el maximo permitido de envite deberia ser 10
	// ~ 5 envidos
	// al 6to saltar error

	p.Cmd("alvaro envido")
	p.Cmd("Roro envido")

	p.Cmd("alvaro envido")
	p.Cmd("Roro envido")

	p.Cmd("alvaro envido")

	p.Cmd("Roro envido") // debe retornar error

	p.Cmd("Roro quiero")

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestAutoQuererse(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":10,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":null,"JugadoresConFlorQueNoCantaron":[]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":1},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":11},{"palo":"Basto","valor":3},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":1},{"palo":"Espada","valor":10},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Copa","valor":12},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":7},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":1},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	// no deberia poder auto quererse **ni auto no-quererse**

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Roro Quiero")     // no lo deberia dejar
	p.Cmd("Roro no-quiero")  // no lo deberia dejar
	p.Cmd("Renzo Quiero")    // no lo deberia dejar
	p.Cmd("Renzo no-quiero") // no lo deberia dejar

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestJsonSinFlores(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	// partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":6},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":1},{"palo":"Copa","valor":2},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Espada","valor":4},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Oro","valor":2},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":6},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":1},{"palo":"Copa","valor":2},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Espada","valor":4},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Oro","valor":2},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":6},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Espada","valor":5},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":1},{"palo":"Copa","valor":2},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Espada","valor":4},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Oro","valor":7},{"palo":"Basto","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Oro","valor":2},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":6},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Espada","valor":5},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":1},{"palo":"Copa","valor":2},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Espada","valor":4},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Oro","valor":7},{"palo":"Basto","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Oro","valor":2},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()
}

func TestFixEnvidoManoEsElUltimo(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":10,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":5,"turno":3,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":null,"JugadoresConFlorQueNoCantaron":[]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":1},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":11},{"palo":"Basto","valor":3},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":1},{"palo":"Espada","valor":10},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Copa","valor":12},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":7},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":1},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("renzo Envido")
	p.Cmd("andres Envido")
	p.Cmd("richard quiero")

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestEnvidoManoSeFue(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":10,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":5,"turno":3,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":null,"JugadoresConFlorQueNoCantaron":[]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":1},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":11},{"palo":"Basto","valor":3},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":1},{"palo":"Espada","valor":10},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":1},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":7},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Copa","valor":12},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("renzo Envido")
	p.Cmd("andres Envido")
	p.Cmd("richard mazo")
	p.Cmd("renzo quiero")

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestFlorBlucle(t *testing.T) {
	p, _ := NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 6},
					&pdt.Carta{Palo: pdt.Basto, Valor: 7},
				},
			},
			pdt.Manojo{ // Roro tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 5},
					&pdt.Carta{Palo: pdt.Espada, Valor: 5},
					&pdt.Carta{Palo: pdt.Espada, Valor: 5},
				},
			},
		},
	)

	p.Cmd("alvaro flor")
	p.Cmd("roro flor")

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestQuieroContraflorDesdeMazo(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 6},
					&pdt.Carta{Palo: pdt.Basto, Valor: 7},
				},
			},
			pdt.Manojo{ // Roro no tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 5},
					&pdt.Carta{Palo: pdt.Espada, Valor: 5},
					&pdt.Carta{Palo: pdt.Basto, Valor: 5},
				},
			},
			pdt.Manojo{ // Adolfo tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Copa, Valor: 2},
					&pdt.Carta{Palo: pdt.Copa, Valor: 3},
				},
			},
			pdt.Manojo{ // Renzo no tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 10},
					&pdt.Carta{Palo: pdt.Copa, Valor: 4},
					&pdt.Carta{Palo: pdt.Copa, Valor: 11},
				},
			},
			pdt.Manojo{ // Andres tiene  flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 4},
					&pdt.Carta{Palo: pdt.Espada, Valor: 4},
					&pdt.Carta{Palo: pdt.Espada, Valor: 1},
				},
			},
			pdt.Manojo{ // Richard no tiene flor
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 12},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 1},
				},
			},
		},
	)

	p.Cmd("alvaro flor")
	p.Cmd("andres mazo") // cambiarlo por andres depues
	p.Cmd("renzo contra-flor")
	p.Cmd("andres quiero")

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestFixSeVaAlMazoYTeniaFlor(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Espada","valor":10},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Espada","valor":10},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Espada","valor":10},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":7},{"palo":"Basto","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":12},{"palo":"Basto","valor":3},{"palo":"Espada","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Oro","valor":7},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":5},{"palo":"Basto","valor":11},{"palo":"Copa","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":5},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":5},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro mazo")
	p.Cmd("roro truco")

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestFixDesconcertante(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":2},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":4},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":2},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":4},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":2},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":3},{"palo":"Espada","valor":11},{"palo":"Copa","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":7},{"palo":"Copa","valor":10},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":4},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Oro","valor":7},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":3},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro flor")
	p.Cmd("alvaro mazo")
	p.Cmd("roro truco")

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestMalaAsignacionPts(t *testing.T) {
	p, _ := NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})

	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Basto, Valor: 12})
	p.Puntajes[pdt.Rojo] = 3
	p.Puntajes[pdt.Azul] = 2
	p.Ronda.Turno = 1
	p.Ronda.ElMano = 1
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{ // Alvaro
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 10},
					&pdt.Carta{Palo: pdt.Oro, Valor: 12},
					&pdt.Carta{Palo: pdt.Espada, Valor: 5},
				},
			},
			pdt.Manojo{ // Roro
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Oro, Valor: 1},
					&pdt.Carta{Palo: pdt.Espada, Valor: 1},
					&pdt.Carta{Palo: pdt.Espada, Valor: 11},
				},
			},
		},
	)

	p.Print()

	p.Cmd("alvaro vale-4")
	p.Cmd("alvaro truco") // vigente
	p.Cmd("roro truco")
	p.Cmd("alvaro re-truco")
	p.Cmd("roro vale-4")
	p.Cmd("alvaro quiero")

	p.Cmd("roro quiero")
	p.Cmd("roro 1 espada")
	p.Cmd("alvaro 12 oro")
	p.Cmd("roro 1 oro")
	p.Cmd("alvaro 5 espada")

	out.Consume(p.Stdout, out.Print)
	p.Print()

	oops = !(p.Puntajes[pdt.Rojo] == 5 && p.Puntajes[pdt.Azul] == 2)
	if oops {
		t.Error(`Asigno mal los puntos`)
		return
	}
}

func TestFixRondaNueva(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Oro","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Oro","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":3},{"palo":"Copa","valor":1},{"palo":"Espada","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Espada","valor":4},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Oro","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":4},{"palo":"Basto","valor":7},{"palo":"Espada","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":12},{"palo":"Espada","valor":10},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("renzo flor")
	p.Cmd("alvaro 1 copa")
	p.Cmd("roro truco")
	p.Cmd("adolfo re-truco")
	p.Cmd("renzo vale-4")
	p.Cmd("adolfo quiero")
	p.Cmd("roro 5 oro")
	p.Cmd("adolfo 6 basto")
	p.Cmd("renzo 19 oro")
	p.Cmd("andres 4 copa")
	p.Cmd("richard 10 espada")
	p.Cmd("renzo 10 oro")
	p.Cmd("andres 7 basto")
	p.Cmd("richard 10 espada")
	p.Cmd("richard 12 copa")
	p.Cmd("alvaro 3 espada")
	p.Cmd("roro 4 espada")
	p.Cmd("adolfo 2 copa")
	p.Cmd("renzo 4 oro")
	p.Cmd("andres 5 espada")

	out.Consume(p.Stdout, out.Print)
	p.Print()
}

func TestFixIrseAlMazo2(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Espada","valor":7},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":3},{"palo":"Copa","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":5},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":12},{"palo":"Oro","valor":2},{"palo":"Copa","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":7},{"palo":"Basto","valor":7},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Copa","valor":3},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("renzo flor")
	p.Cmd("alvaro envido")
	p.Cmd("alvaro 6 basto")
	p.Cmd("roro 11 oro")
	p.Cmd("adolfo 2 basto")
	p.Cmd("renzo 12 oro")
	p.Cmd("andres 7 copa")
	p.Cmd("richard 12 basto")
	p.Cmd("roro 3 espada")
	p.Cmd("adolfo 5 oro")
	p.Cmd("renzo 2 oro")
	p.Cmd("andres 7 basto")
	p.Cmd("richard truco")
	p.Cmd("renzo quiero")
	p.Cmd("andres re-truco")
	p.Cmd("richard vale-4")
	p.Cmd("alvaro quiero")
	p.Cmd("roro mazo")

	out.Consume(p.Stdout, out.Print)

	p.Print()

}

func TestFixDecirQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":10},{"palo":"Espada","valor":12},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":10},{"palo":"Espada","valor":12},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":1},{"palo":"Copa","valor":6},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":7},{"palo":"Copa","valor":12},{"palo":"Oro","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":2},{"palo":"Espada","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":10},{"palo":"Espada","valor":12},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":3},{"palo":"Basto","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":6},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":2},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("renzo flor")
	p.Cmd("alvaro truco")
	p.Cmd("renzo quiero")
	p.Cmd("alvaro re-truco")
	p.Cmd("renzo vale-4")
	p.Cmd("alvaro quiero")
	p.Cmd("renzo re-truco")

	out.Consume(p.Stdout, out.Print)
	p.Print()

}

func TestFixPanicNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":1},{"palo":"Basto","valor":10},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":4},{"palo":"Espada","valor":6},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":5},{"palo":"Basto","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Oro","valor":12},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":7},{"palo":"Oro","valor":11},{"palo":"Oro","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":2},{"palo":"Basto","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":1},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	p.Cmd("alvaro flor")
	p.Cmd("alvaro 1 basto")
	p.Cmd("renzo flor")

	p.Cmd("alvaro 1 basto")
	p.Cmd("roro 4 copa")
	p.Cmd("adolfo 2 espada")
	p.Cmd("renzo 4 oro") // la Primera mano la gana renzo
	p.Cmd("andres 7 espada")
	p.Cmd("richard 10 copa")

	p.Cmd("renzo 12 oro")
	p.Cmd("andres 11 oro") // la seguna mano la gana andres
	p.Cmd("richard 2 copa")
	p.Cmd("alvaro 10 basto")
	p.Cmd("roro 6 espada")
	p.Cmd("adolfo 5 copa")

	p.Cmd("andres 3 oro")
	p.Cmd("richard truco")
	p.Cmd("richard mazo")
	p.Cmd("alvaro quiero")
	p.Cmd("alvaro re-truco")
	p.Cmd("renzo quiero")
	p.Cmd("roro vale-4")
	p.Cmd("andres no-quiero")

	out.Consume(p.Stdout, out.Print)
	p.Print()

}

func TestFixOutput(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":4},{"palo":"Oro","valor":3},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":1},{"palo":"Oro","valor":12},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Espada","valor":4},{"palo":"Basto","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":1},{"palo":"Basto","valor":11},{"palo":"Espada","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":6},{"palo":"Basto","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":10},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.Force(partidaJSON)
	p.Print()

	// p.Cmd("alvaro 4 copa")
	p.Cmd("roro flor")
	// p.Cmd("roro 1 oro")
	// p.Cmd("adolfo flor")
	// p.Cmd("adolfo 2 copa")
	// p.Cmd("renzo mazo")
	// p.Cmd("andres 12 espada")
	// p.Cmd("richard truco")
	// p.Cmd("andres quiero")
	// p.Cmd("richard 11 oro")
	// p.Cmd("andres re-truco")
	// p.Cmd("renzo quiero")
	// p.Cmd("roro quiero")
	// p.Cmd("andres 5 basto")
	// p.Cmd("richard 6 espada")
	// p.Cmd("alvaro 7 basto")
	// p.Cmd("roro 12 oro")
	// p.Cmd("adolfo 4 espada")

	out.Consume(p.Stdout, out.Print)
	p.Print()

}
