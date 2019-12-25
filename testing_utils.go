package truco

var (
	oops = false

	/* JUGADORES (en sentido antihorario) */

	// jugadores clasicos
	jugadores = []Jugador{
		// 2 jugadores
		Jugador{"Juan", Rojo, nil},
		Jugador{"Pedro", Azul, nil},
		// 4 jugadores
		Jugador{"Jacinto", Rojo, nil},
		Jugador{"Patricio", Azul, nil},
		// 6 jugadores
		Jugador{"Jaime", Rojo, nil},
		Jugador{"Paco", Azul, nil},
	}

	juan     = &jugadores[0] // rojo
	pedro    = &jugadores[1] // azul
	jacinto  = &jugadores[2] // rojo
	patricio = &jugadores[3] // azul

	// jugadores YT
	jugadoresYT = []Jugador{
		Jugador{"A", Rojo, nil}, // Dr. Favaloro
		Jugador{"B", Azul, nil}, // Fangio
		Jugador{"C", Rojo, nil}, // inodoro
		Jugador{"D", Azul, nil}, // San Martin
		Jugador{"E", Rojo, nil}, // Yupanqui
		Jugador{"F", Azul, nil}, // Gardel
	}

	// A Dr. Favaloro
	A = &jugadoresYT[0] // rojo
	// B Fangio
	B = &jugadoresYT[1] // azul
	// C inodoro
	C = &jugadoresYT[2] // rojo
	// D San Martin
	D = &jugadoresYT[3] // azul
	// E Yupanqui
	E = &jugadoresYT[4] // azul
	// F Gardel
	F = &jugadoresYT[5] // azul

	muestra = Carta{
		Palo:  Espada,
		Valor: 1,
	}

	manojos = []Manojo{
		Manojo{
			Cartas: [3]Carta{ // envido: 13
				Carta{Palo: Oro, Valor: 7},
				Carta{Palo: Oro, Valor: 6},
				Carta{Palo: Copa, Valor: 5},
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
		Manojo{
			Cartas: [3]Carta{
				Carta{Palo: Copa, Valor: 4},
				Carta{Palo: Oro, Valor: 5},
				Carta{Palo: Basto, Valor: 2},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{
				Carta{Palo: Copa, Valor: 10},
				Carta{Palo: Oro, Valor: 3},
				Carta{Palo: Basto, Valor: 11},
			},
			jugador: nil,
		},
	}

	manojosYT1 = []Manojo{
		Manojo{
			Cartas: [3]Carta{ // envido: 26
				Carta{Palo: Oro, Valor: 6},
				Carta{Palo: Oro, Valor: 12},
				Carta{Palo: Copa, Valor: 5},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{ // envido: 20
				Carta{Palo: Copa, Valor: 12},
				Carta{Palo: Copa, Valor: 11},
				Carta{Palo: Basto, Valor: 3},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{ // envido: 28
				Carta{Palo: Copa, Valor: 2},
				Carta{Palo: Copa, Valor: 6},
				Carta{Palo: Basto, Valor: 1},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{ // envido: 25
				Carta{Palo: Oro, Valor: 2},
				Carta{Palo: Oro, Valor: 3},
				Carta{Palo: Basto, Valor: 2},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{ // envido: 33
				Carta{Palo: Basto, Valor: 6},
				Carta{Palo: Basto, Valor: 7},
				Carta{Palo: Oro, Valor: 5},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{ // envido: 27
				Carta{Palo: Copa, Valor: 3},
				Carta{Palo: Copa, Valor: 4},
				Carta{Palo: Oro, Valor: 4},
			},
			jugador: nil,
		},
	}

	manojosYT2 = []Manojo{
		Manojo{
			Cartas: [3]Carta{ // envido: 21
				Carta{Palo: Basto, Valor: 1},
				Carta{Palo: Basto, Valor: 12},
				Carta{Palo: Copa, Valor: 5},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{ // envido: 23
				Carta{Palo: Oro, Valor: 12},
				Carta{Palo: Oro, Valor: 3},
				Carta{Palo: Basto, Valor: 4},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{ // envido: 23
				Carta{Palo: Basto, Valor: 10},
				Carta{Palo: Copa, Valor: 6},
				Carta{Palo: Basto, Valor: 3},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{ // envido: 30
				Carta{Palo: Oro, Valor: 6},
				Carta{Palo: Oro, Valor: 4},
				Carta{Palo: Copa, Valor: 1},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{ // envido: 30
				Carta{Palo: Basto, Valor: 6},
				Carta{Palo: Basto, Valor: 4},
				Carta{Palo: Oro, Valor: 1},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{ // envido: 31
				Carta{Palo: Espada, Valor: 5},
				Carta{Palo: Copa, Valor: 4},
				Carta{Palo: Espada, Valor: 3},
			},
			jugador: nil,
		},
	}

	partida4JugadoresEnvidoTesting = Partida{
		puntuacion:    a20,
		cantJugadores: 4,
		jugadores:     jugadores[:4],
		puntajes:      [2]int{4, 3},
		Ronda: Ronda{
			manoEnJuego: primera,
			elMano:      2,
			turno:       2,
			envido:      Envido{puntaje: 0, estado: NOCANTADOAUN},
			truco:       NOCANTADO,
			manojos:     manojos[:4],
			manos:       make([]Mano, 3),
			muestra:     muestra,
		},
	}

	partidaYT1 = Partida{
		puntuacion:    a20,
		cantJugadores: 6,
		jugadores:     jugadoresYT[:6],
		puntajes:      [2]int{4, 3},
		Ronda: Ronda{
			manoEnJuego: primera,
			elMano:      0,
			turno:       0,
			envido:      Envido{puntaje: 0, estado: NOCANTADOAUN},
			truco:       NOCANTADO,
			manojos:     manojosYT1,
			manos:       make([]Mano, 3),
			muestra:     muestra,
		},
	}

	partidaYT2 = Partida{
		puntuacion:    a20,
		cantJugadores: 6,
		jugadores:     jugadoresYT[:6],
		puntajes:      [2]int{4, 3},
		Ronda: Ronda{
			manoEnJuego: primera,
			elMano:      0,
			turno:       0,
			envido:      Envido{puntaje: 0, estado: NOCANTADOAUN},
			truco:       NOCANTADO,
			manojos:     manojosYT2,
			manos:       make([]Mano, 3),
			muestra:     muestra,
		},
	}
)

func getPartidaCustom1() Partida {
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
			manojos:     manojos[:2],
			manos:       make([]Mano, 3),
			muestra:     muestra,
		},
	}
	p.dobleLinking()
	p.Ronda.getManoActual().repartidor = p.Ronda.elMano
	return p
}
