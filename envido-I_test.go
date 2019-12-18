package main

import (
	"testing"
	// "fmt"
	// "bufio"
	// "os"
)

// sinopsis:
// juan toca envido
// pedro responde quiero
// juan gana con 12 de envido vs 5 de pedro
func TestEnvidoAceptado(t *testing.T) {
	p := partidaDefault2Jugadores
	dobleLinking(&p)
	p.ronda.getManoActual().repartidor = p.ronda.elMano

	p.ronda.Print()

	// empieza primera ronda
	// empieza primera mano

	// Juan toca envido
	jugada := tocarEnvido{}
	jugada.hacer(&p, &p.jugadores[0])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error("El estado del envido deberia de ser `envido`")
		return
	}

	oops = p.ronda.envido.puntaje != 2
	if oops {
		t.Error("El `puntaje` del envido deberia de ser 2")
		return
	}

	// Pedro responde 'quiero'
	responderQuiero{}.hacer(&p, &p.jugadores[1])

	oops = p.ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado', 
		ya que fue aceptado por Pedro`)
		return
	}

	oops = p.ronda.envido.puntaje != 2
	if oops {
		t.Error(`El puntaje del envido deberia de ser 2`)
		return
	}
}

// sinopsis:
// juan: envido
// pedro: envido
// todo: ??
func TestDobleEnvido(t *testing.T) {
	p := partidaDefault2Jugadores

	dobleLinking(&p)
	p.ronda.getManoActual().repartidor = p.ronda.elMano

	p.ronda.Print()

	// empieza primera ronda
	// empieza primera mano

	// Juan toca envido
	jugada := tocarEnvido{}
	jugada.hacer(&p, &p.jugadores[0])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error("El estado del envido deberia de ser `envido`")
		return
	}

	oops = p.ronda.envido.puntaje != 2
	if oops {
		t.Error("El `puntaje` del envido deberia de ser 2")
		return
	}

	// Pedro redobla el envido
	tocarEnvido{}.hacer(&p, &p.jugadores[1])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error(`El estado del envido deberia de ser 'envido', incluso luego de que
		ambos Juan y Pedro lo hayan tocando`)
		return
	}

	oops = p.ronda.envido.puntaje != 4
	if oops {
		t.Error(`El puntaje del envido deberia ahora de ser '2 + 2 = 4'`)
		return
	}

	responderQuiero{}.hacer(&p, &p.jugadores[0])
}

// sinopsis:
// Juan: envido
// Pedro: envido
// Juan: no quiero
func TestDobleEnvidoNoAceptado(t *testing.T) {
	p := partidaDefault2Jugadores

	dobleLinking(&p)
	p.ronda.getManoActual().repartidor = p.ronda.elMano

	p.ronda.Print()

	// empieza primera ronda
	// empieza primera mano

	// Juan toca envido
	jugada := tocarEnvido{}
	jugada.hacer(&p, &p.jugadores[0])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error("El estado del envido deberia de ser `envido`")
		return
	}

	oops = p.ronda.envido.puntaje != 2
	if oops {
		t.Error("El `puntaje` del envido deberia de ser 2")
		return
	}

	// Pedro redobla el envido
	tocarEnvido{}.hacer(&p, &p.jugadores[1])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error(`El estado del envido deberia de ser 'envido', incluso luego de que
		ambos Juan y Pedro lo hayan tocando`)
		return
	}

	oops = p.ronda.envido.puntaje != 4
	if oops {
		t.Error(`El puntaje del envido deberia ahora de ser '2 + 2 = 4'`)
		return
	}

	// Juan responde 'no quiero'
	responderNoQuiero{}.hacer(&p, &p.jugadores[0])
}
