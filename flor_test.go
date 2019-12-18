package main

import (
	"testing"
	// "fmt"
	// "bufio"
	// "os"
)

var (
	muestra2 = Carta{
		Palo:  Copa,
		Valor: 5,
	}

	manojos2 = []Manojo{
		Manojo{
			Cartas: [3]Carta{
				Carta{Palo: Oro, Valor: 6},
				Carta{Palo: Copa, Valor: 10},
				Carta{Palo: Copa, Valor: 7},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{
				Carta{Palo: Copa, Valor: 1},
				Carta{Palo: Oro, Valor: 2},
				Carta{Palo: Basto, Valor: 3},
			},
			jugador: nil,
		},
	}
)

func TestFlorI(t *testing.T) {
	p := Partida{
		puntuacion:    a20,
		puntaje:       0,
		cantJugadores: 2,
		jugadores:     jugadores[:2],
		ronda: Ronda{
			manoEnJuego: primera,
			elMano:      0,
			turno:       0,
			envido:      Envido{puntaje: 0, estado: NOCANTADOAUN},
			truco:       NOCANTADO,
			manojos:     manojos2,
			manos:       make([]Mano, 3),
			muestra:     muestra2,
		},
	}

	dobleLinking(&p)
	p.ronda.getManoActual().repartidor = p.ronda.elMano

	p.ronda.Print()

	// juan
	juan := p.jugadores[0]
	tieneFlor, _ := juan.manojo.tieneFlor(p.ronda.muestra)
	oops = tieneFlor == true
	if oops {
		t.Error(`Juan' NO deberia de tener 'flor'`)
		return
	}

}
