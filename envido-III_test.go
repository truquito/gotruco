package main

import (
	"testing"
)

/* Tests de Youtube */
func TestYTEnvidoCalc(t *testing.T) {
	p := partidaYT1
	dobleLinking(&p)
	p.ronda.getManoActual().repartidor = 5

	expected := []int{ 26, 20, 28, 25, 33, 27 }
	for i, jugador := range p.jugadores {
		got := jugador.manojo.calcularEnvido(p.ronda.muestra)
		oops = expected[i] != got
		if oops {
			t.Errorf(
				`El resultado del envido del jugador %s es incorrecto.
				\nEXPECTED: %v
				\nGOT: %v`, 
				jugador.nombre, expected[i], got)
			return
		}
	}
	
}

func TestYTEnvidoI(t *testing.T) {
	p := partidaYT1
	dobleLinking(&p)
	p.ronda.getManoActual().repartidor = 5

	tocarEnvido{}.hacer(&p, D)
	responderQuiero{}.hacer(&p, C)
	oops = p.puntajes[Rojo] != 4+2
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

func TestYTEnvidoCalcII(t *testing.T) {
	p := partidaYT2
	dobleLinking(&p)
	p.ronda.getManoActual().repartidor = 5

	expected := []int{ 21, 23, 23, 30, 30, 31 }
	for i, jugador := range p.jugadores {
		got := jugador.manojo.calcularEnvido(p.ronda.muestra)
		oops = expected[i] != got
		if oops {
			t.Errorf(
				`El resultado del envido del jugador %s es incorrecto.
				\nEXPECTED: %v
				\nGOT: %v`, 
				jugador.nombre, expected[i], got)
			return
		}
	}	
}

func TestYTEnvidoII(t *testing.T) {
	p := partidaYT2
	dobleLinking(&p)
	p.ronda.getManoActual().repartidor = 5

	tocarEnvido{}.hacer(&p, D)
	responderQuiero{}.hacer(&p, C)
	oops = p.puntajes[Azul] != 3+2
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}