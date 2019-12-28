package truco

import (
	"sync"
	"testing"
	// "fmt"
	// "bufio"
	// "os"
)

var (
	CincoDeCopa = Carta{
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

func TestNoDeberiaTenerFlor(t *testing.T) {

	// juan := jugadores[0]
	// juan.manojo = &manojos2[0]
	// muestra := Carta{Palo: Copa, Valor: 5}
	// tieneFlor, _ := juan.manojo.tieneFlor(muestra)
	// oops = tieneFlor == true
	// if oops {
	// 	t.Error(`Juan' NO deberia de tener 'flor'`)
	// 	return
	// }

}

func TestDeberiaTenerFlor(t *testing.T) {
	// juan := jugadores[0]
	// juan.manojo = &Manojo{
	// 	Cartas: [3]Carta{
	// 		Carta{Palo: Copa, Valor: 4},
	// 		Carta{Palo: Espada, Valor: 10},
	// 		Carta{Palo: Espada, Valor: 7},
	// 	},
	// 	jugador: nil,
	// }
	// muestra := Carta{Palo: Copa, Valor: 5}
	// tieneFlor, _ := juan.manojo.tieneFlor(muestra)
	// oops = tieneFlor == false
	// if oops {
	// 	t.Error(`Juan' Deberia de tener 'flor'`)
	// 	return
	// }
}

func TestJuanNoDeberiaTenerFlor(t *testing.T) {
	p := Partida{
		puntuacion:    a20,
		puntaje:       0,
		cantJugadores: 2,
		jugadores:     jugadores[:2],
		Ronda: Ronda{
			manoEnJuego: primera,
			elMano:      0,
			turno:       0,
			envido:      Envido{puntaje: 0, estado: NOCANTADOAUN},
			truco:       NOCANTADO,
			manojos:     manojos2,
			manos:       make([]Mano, 3),
			muestra:     CincoDeCopa,
		},
	}

	p.Ronda.singleLinking(p.jugadores)
	p.Ronda.getManoActual().repartidor = p.Ronda.elMano

	p.Ronda.Print()

	// juan
	manojoJuan := p.Ronda.manojos[0]
	tieneFlor, _ := manojoJuan.tieneFlor(p.Ronda.muestra)
	oops = tieneFlor == true
	if oops {
		t.Error(`Juan' NO deberia de tener 'flor'`)
		return
	}

}

func TestFlor(t *testing.T) {

	p := Partida{
		puntuacion:    a20,
		puntaje:       0,
		cantJugadores: 6, // puede ser 2, 4 o 6
		jugadores:     jugadores,
		Ronda: Ronda{
			manoEnJuego: primera,
			elMano:      0,
			turno:       0,
			envido:      Envido{puntaje: 0, estado: NOCANTADOAUN},
			flor:        NOCANTADA,
			truco:       NOCANTADO,
			muestra:     Carta{Palo: Oro, Valor: 3},
			manojos: []Manojo{
				Manojo{
					Cartas: [3]Carta{ // Juan tiene flor
						Carta{Palo: Oro, Valor: 2},
						Carta{Palo: Basto, Valor: 6},
						Carta{Palo: Basto, Valor: 7},
					},
					jugador: nil,
				},
				Manojo{
					Cartas: [3]Carta{ // Pedro no tiene flor
						Carta{Palo: Oro, Valor: 5},
						Carta{Palo: Espada, Valor: 5},
						Carta{Palo: Basto, Valor: 5},
					},
					jugador: nil,
				},
				Manojo{
					Cartas: [3]Carta{ // Jacinto tiene flor
						Carta{Palo: Copa, Valor: 1},
						Carta{Palo: Copa, Valor: 2},
						Carta{Palo: Copa, Valor: 3},
					},
					jugador: nil,
				},
				Manojo{
					Cartas: [3]Carta{ // Patricio tiene flor
						Carta{Palo: Oro, Valor: 4},
						Carta{Palo: Espada, Valor: 4},
						Carta{Palo: Espada, Valor: 1},
					},
					jugador: nil,
				},
				Manojo{
					Cartas: [3]Carta{ // Jaime no tiene  flor
						Carta{Palo: Copa, Valor: 10},
						Carta{Palo: Oro, Valor: 7},
						Carta{Palo: Basto, Valor: 11},
					},
					jugador: nil,
				},
				Manojo{
					Cartas: [3]Carta{ // Paco tiene flor
						Carta{Palo: Oro, Valor: 10},
						Carta{Palo: Oro, Valor: 2},
						Carta{Palo: Basto, Valor: 1},
					},
					jugador: nil,
				},
			},
			manos: make([]Mano, 3),
		},
	}
	p.puntajes[Rojo] = 0
	p.puntajes[Azul] = 0

	p.Ronda.singleLinking(p.jugadores)

	ImprimirJugadas()

	go func() {
		for {
			sigJugada := p.getSigJugada()
			sigJugada.hacer(&p)
		}
	}()

	// fin capa logica/privada -----------------------

	p.SetSigJugada("Juan Flor")
	p.SetSigJugada("Pedro Mazo")
	// p.SetSigJugada("Pedro Mazo") // test 2 vecees se va al mazo
	// p.SetSigJugada("Pedro Mazo") // test despues de irse al mazo manda flor o algo asi
	p.SetSigJugada("Patricio Flor")
	p.SetSigJugada("Jacinto Contra-flor-al-resto")
	p.SetSigJugada("Paco Quiero")

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

}
