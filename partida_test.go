package truco

import (
	"testing"
	"time"
)

func TestPartida1(t *testing.T) {
	cantJugadores := 6
	p := Partida{
		puntuacion:    a20,
		puntaje:       0,
		cantJugadores: cantJugadores, // puede ser 2, 4 o 6
		jugadores:     getDefaultJugadores(6),
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
					Cartas: [3]Carta{ // Alvaro tiene flor
						Carta{Palo: Oro, Valor: 2},
						Carta{Palo: Basto, Valor: 6},
						Carta{Palo: Basto, Valor: 7},
					},
					jugador: nil,
				},
				Manojo{
					Cartas: [3]Carta{ // Roro no tiene flor
						Carta{Palo: Oro, Valor: 5},
						Carta{Palo: Espada, Valor: 5},
						Carta{Palo: Basto, Valor: 5},
					},
					jugador: nil,
				},
				Manojo{
					Cartas: [3]Carta{ // Adolfo tiene flor
						Carta{Palo: Copa, Valor: 1},
						Carta{Palo: Copa, Valor: 2},
						Carta{Palo: Copa, Valor: 3},
					},
					jugador: nil,
				},
				Manojo{
					Cartas: [3]Carta{ // Renzo tiene flor
						Carta{Palo: Oro, Valor: 4},
						Carta{Palo: Espada, Valor: 4},
						Carta{Palo: Espada, Valor: 1},
					},
					jugador: nil,
				},
				Manojo{
					Cartas: [3]Carta{ // Andrés no tiene  flor
						Carta{Palo: Copa, Valor: 10},
						Carta{Palo: Oro, Valor: 7},
						Carta{Palo: Basto, Valor: 11},
					},
					jugador: nil,
				},
				Manojo{
					Cartas: [3]Carta{ // Richard tiene flor
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

	p.sigJugada = make(chan string, 1)

	ImprimirJugadas()
	p.Ronda.Print()

	go func() {
		for {
			sjugada, sjugador := p.getSigJugada()
			sjugada.hacer(&p, sjugador)
		}
	}()

	// fin capa logica/privada -----------------------

	p.SetSigJugada("Alvaro Envido")
	p.SetSigJugada("Alvaro Flor")
	p.SetSigJugada("Roro Mazo")
	p.SetSigJugada("Adolfo Flor")
	p.SetSigJugada("Renzo Contra-flor")
	p.SetSigJugada("Alvaro Quiero")
	// todo: si se envia un setSigJugada sin jugador
	// se muere el proc
	// ejemplo p.SetSigJugada("Quiero")

	time.Sleep(60 * time.Minute)

}
