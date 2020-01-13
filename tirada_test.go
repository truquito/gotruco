package truco

import (
	"testing"
)

func TestTirada1(t *testing.T) {
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

	p.SetSigJugada("Alvaro 2 Oro")
	p.SetSigJugada("Roro 5 Oro")
	p.SetSigJugada("Adolfo 1 Copa")
	p.SetSigJugada("Renzo 4 Oro")
	p.SetSigJugada("Andres 10 Copa")
	p.SetSigJugada("Richard 10 Oro")
	p.Esperar()

	// como la muestra es Palo: Oro, Valor: 3 -> gana alvaro
	if !(len(p.Ronda.Manos[primera].CartasTiradas) == 6) {
		t.Error("La cantidad de cartas tiradas deberia ser 6")
		return

	} else if !(p.Ronda.Manos[primera].Ganador.Jugador.Nombre == "Alvaro") {
		t.Error("El ganador de la priemra mano deberia ser Alvaro")
		return

	} else if !(p.Ronda.Manos[primera].Resultado == GanoAzul) {
		t.Error("El equipo ganador de la priemra mano deberia ser Azul")
		return
	}

	// como alvaro gano la mano anterior -> empieza tirando el
	p.SetSigJugada("Alvaro 6 Basto")
	p.SetSigJugada("Roro 5 Espada")
	p.SetSigJugada("Adolfo 2 Copa")
	p.SetSigJugada("Renzo 4 Espada")
	p.SetSigJugada("Andres 7 Oro")
	p.SetSigJugada("Richard 2 Oro")
	p.Esperar()

	// como la muestra es Palo: Oro, Valor: 3 -> gana richard
	if !(len(p.Ronda.Manos[segunda].CartasTiradas) == 6) {
		t.Error("La cantidad de cartas tiradas deberia ser 6")
		return

	} else if !(p.Ronda.Manos[segunda].Ganador.Jugador.Nombre == "Richard") {
		t.Error("El ganador de la priemra mano deberia ser Richard")
		return

	} else if !(p.Ronda.Manos[segunda].Resultado == GanoRojo) {
		t.Error("El equipo ganador de la priemra mano deberia ser Rojo")
		return
	}

	// vuelvo a checkear que el estado de la primera nos se haya editado
	if !(len(p.Ronda.Manos[primera].CartasTiradas) == 6) {
		t.Error("La cantidad de cartas tiradas deberia ser 6")
		return

	} else if !(p.Ronda.Manos[primera].Ganador.Jugador.Nombre == "Alvaro") {
		t.Error("El ganador de la priemra mano deberia ser Alvaro")
		return

	} else if !(p.Ronda.Manos[primera].Resultado == GanoAzul) {
		t.Error("El equipo ganador de la priemra mano deberia ser Azul")
		return
	}

	// como richard gano la mano anterior -> empieza tirando el
	p.SetSigJugada("Richard 1 Basto")
	p.SetSigJugada("Alvaro 7 Basto")
	p.SetSigJugada("Roro 5 Basto")
	p.SetSigJugada("Adolfo 3 Copa")
	p.SetSigJugada("Renzo 1 Espada")
	p.SetSigJugada("Andres 11 Basto")
	p.Esperar()

	// para este momento ya cambio a una nueva ronda
	// como la muestra es Palo: Oro, Valor: 3 -> gana Renzo con el 1 de espada
	// 1 mano ganada por azul; 2 por rojo -> ronda ganada por rojo
	if !(p.Puntajes[Rojo] == 1) {
		t.Error("El puntaje del equipo Rojo deberia ser 1 porque gano la ronda")
		return

	}

	p.Esperar()
	p.Terminar()
}
