package main

import (
	"bufio"
	"fmt"
	"os"
)

// Puntuacion : Enum para el puntaje maximo de la partida
type Puntuacion int

// hasta 15 pts, 20 pts, 30 pts o 40 pts
const (
	a20 Puntuacion = 0
	a30 Puntuacion = 1
	a40 Puntuacion = 2
)

func (pt Puntuacion) toInt() int {
	switch pt {
	case a20:
		return 20
	case a30:
		return 30
	default:
		return 40
	}
}

// Equipo : Enum para el puntaje maximo de la partida
type Equipo int

// rojo o azul
const (
	Rojo Equipo = 0
	Azul Equipo = 1
)

func (e Equipo) String() string {
	if e == Rojo {
		return "Rojo"
	}
	return "Azul"
}

// Partida :
type Partida struct {
	cantJugadores int
	jugadores     []Jugador
	puntuacion    Puntuacion
	puntaje       int
	puntajes      [2]int // Rojo o Azul
	ronda         Ronda
}

func (p *Partida) esperandoJugada() {
	if debuggingMode {
		return
	}

	imprimirJugadas()
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\nIngresar jugador:\n")
	jugadorStr, _ := reader.ReadString('\n')
	jugador, _ := decodeJugador(jugadorStr, p.jugadores)

	fmt.Printf("\nIngresar jugada:\n")
	jugadaStr, _ := reader.ReadString('\n')

	p.procesarJugada(jugador, jugadaStr)
}

func (p *Partida) procesarJugada(jugador *Jugador, jugadaStr string) {
	var jugada IJugada

	switch jugadaStr {
	// toques
	case "Envido":
		jugada = tocarEnvido{}
	case "Real-envido":
		jugada = tocarRealEnvido{}
	case "Falta-envido":
		jugada = tocarFaltaEnvido{}

	// cantos
	case "Flor":
		jugada = tocarEnvido{}
	case "Contra-flor":
		jugada = tocarRealEnvido{}
	case "Contra-flor-al-resto":
		jugada = tocarFaltaEnvido{}

	// gritos
	case "Truco":
		jugada = gritarTruco{}
	case "Re-truco":
		jugada = gritarReTruco{}
	case "Vale-4":
		jugada = gritarVale4{}

	// respuestas
	case "Quiero":
		jugada = responderQuiero{}
	case "No-Quiero":
		jugada = responderNoQuiero{}
	case "Tiene":
		jugada = responderNoQuiero{}

	// acciones
	case "Mazo":
		jugada = irseAlMazo{}
	default:
		panic("lols")
	}

	jugada.hacer(p, jugador)
}

func (p *Partida) dobleLinking() {
	// hago el doble-linking "jugadores <-> manojos"
	for i := 0; i < p.cantJugadores; i++ {
		p.jugadores[i].manojo = &p.ronda.manojos[i]
		p.ronda.manojos[i].jugador = &p.jugadores[i]
	}
}

func (p *Partida) getMaxPuntaje() int {
	if p.puntajes[Rojo] > p.puntajes[Azul] {
		return p.puntajes[Rojo]
	}
	return p.puntajes[Azul]
}

// getPuntuacionMalas devuelve la mitad de la puntuacion
// total jugable durante toda la partida
func (p *Partida) getPuntuacionMalas() int {
	return p.puntuacion.toInt() / 2
}

// getJugador dado un indice de jugador,
// devuelve su puntero correspondiente
func (p *Partida) getJugador(jIdx JugadorIdx) *Jugador {
	return &p.jugadores[jIdx]
}

// sig devuelve el `JugadorIdx` del
// jugador siguiente a j
func (p *Partida) sig(j JugadorIdx) JugadorIdx {
	esElUltimo := int(j) == p.cantJugadores-1
	if esElUltimo {
		return 0
	}
	return j + 1
}

func nuevaPartida(puntuacion Puntuacion, jugadores []Jugador) *Partida {
	partida := Partida{
		puntuacion:    puntuacion,
		puntaje:       0,
		cantJugadores: len(jugadores), // puede ser 2, 4 o 6
		jugadores:     jugadores,
		ronda: Ronda{
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
	partida.puntajes[Rojo] = 0
	partida.puntajes[Azul] = 0

	partida.dobleLinking()

	return &partida
}
