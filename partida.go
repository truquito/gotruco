package truco

import "strings"

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

	sigJugada chan string
}

func (p *Partida) readLnJugada() error {
	return nil
}

// nexo capa presentacion con capa logica
func (p *Partida) setSigJugada(cmd string) {
	p.sigJugada <- cmd
}

func (p *Partida) getSigJugada() (IJugada, *Jugador) {
	cmd := <-p.sigJugada
	params := strings.Fields(cmd)
	jugadaStr, jugadorStr := params[1], params[0]
	return p.parseJugada(jugadorStr, jugadaStr)

}

func (p *Partida) parseJugada(jugadorStr string, jugadaStr string) (IJugada, *Jugador) {
	var (
		jugador, _ = parseJugador(jugadorStr, p.jugadores)
		jugada     IJugada
	)

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
		jugada = cantarFlor{}
	case "Contra-flor":
		jugada = cantarContraFlor{}
	case "Contra-flor-al-resto":
		jugada = cantarContraFlorAlResto{}

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

	return jugada, jugador
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

// retorna true si la partida acabo
func (p *Partida) noAcabada() bool {
	return p.getMaxPuntaje() < p.puntuacion.toInt()
}

func (p *Partida) elChico() int {
	return p.puntuacion.toInt() / 2
}

// retorna true si `e` esta en malas
func (p *Partida) estaEnMalas(e Equipo) bool {
	return p.puntajes[e] < p.elChico()
}

// retorna el equipo que va ganando
func (p *Partida) elQueVaGanando() Equipo {
	vaGanandoRojo := p.puntajes[Rojo] > p.puntajes[Azul]
	if vaGanandoRojo {
		return Rojo
	}
	return Azul
}

// retorna la cantidad de puntos que le corresponderian
// a `ganadorDelEnvite` si hubiese ganado un "Contra flor al resto"
// sin tener en cuenta los puntos acumulados de envites anteriores
func (p *Partida) calcPtsContraFlorAlResto(ganadorDelEnvite Equipo) int {
	// si el que va ganando:
	// 		esta en Malas -> el ganador del envite (`ganadorDelEnvite`) gana el chico
	// 		esta en Buenas -> el ganador del envite (`ganadorDelEnvite`) gana lo que le falta al maximo para ganar la ronda

	if p.estaEnMalas(p.elQueVaGanando()) {
		loQueLeFaltaAlGANADORparaGanarElChico := p.elChico() - p.puntajes[ganadorDelEnvite]
		return loQueLeFaltaAlGANADORparaGanarElChico
	} else {
		loQueLeFaltaAlQUEvaGANANDOparaGanarElChico := p.puntuacion.toInt() - p.puntajes[p.elQueVaGanando()]
		return loQueLeFaltaAlQUEvaGANANDOparaGanarElChico
	}

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
			flor:        NOCANTADA,
			truco:       NOCANTADO,
			manojos:     manojos[:2],
			manos:       make([]Mano, 3),
			muestra:     muestra,
		},
	}
	partida.puntajes[Rojo] = 0
	partida.puntajes[Azul] = 0

	partida.dobleLinking()

	partida.sigJugada = make(chan string, 1)

	return &partida
}
