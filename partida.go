package truco

import (
	"fmt"
	"regexp"
	"strings"
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
	Ronda         Ronda

	sigJugada chan string
}

func (p *Partida) readLnJugada() error {
	return nil
}

// SetSigJugada nexo capa presentacion con capa logica
func (p *Partida) SetSigJugada(cmd string) error {
	// checkeo de sintaxis
	ok := regexp.MustCompile(`^(\w|-)+\s(\w|-)+$`).MatchString(cmd)
	if !ok {
		return fmt.Errorf("Comando incorrecto")
	}

	p.sigJugada <- cmd
	return nil
}

// devuelve solo la siguiente jugada VALIDA
// si no es valida es como si no hubiese pasado nada
func (p *Partida) getSigJugada() (IJugada, *Manojo) {
	var (
		manojo *Manojo
		jugada IJugada
		err    error
	)
	for {
		cmd := <-p.sigJugada
		params := strings.Fields(cmd)
		jugadaStr, jugadorStr := params[1], params[0]

		manojo, err = p.Ronda.getManojo(jugadorStr)
		if err == nil {
			jugada, err = p.parseJugada(jugadaStr)
			if err == nil {
				return jugada, manojo
			}
		}
		fmt.Println(err.Error())
	}
}

func (p *Partida) parseJugada(jugadaStr string) (IJugada, error) {
	var jugada IJugada

	jugadaStr = strings.ToLower(jugadaStr)

	switch jugadaStr {
	// toques
	case "envido":
		jugada = tocarEnvido{}
	case "real-envido":
		jugada = tocarRealEnvido{}
	case "falta-envido":
		jugada = tocarFaltaEnvido{}

	// cantos
	case "flor":
		jugada = cantarFlor{}
	case "contra-flor":
		jugada = cantarContraFlor{}
	case "contra-flor-al-resto":
		jugada = cantarContraFlorAlResto{}

	// gritos
	case "truco":
		jugada = gritarTruco{}
	case "re-truco":
		jugada = gritarReTruco{}
	case "vale-4":
		jugada = gritarVale4{}

	// respuestas
	case "quiero":
		jugada = responderQuiero{}
	case "no-Quiero":
		jugada = responderNoQuiero{}
	case "tiene":
		jugada = responderNoQuiero{}

	// acciones
	case "mazo":
		jugada = irseAlMazo{}
	default:
		return nil, fmt.Errorf("No esxiste esa jugada")
	}

	return jugada, nil
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

// NoAcabada retorna true si la partida acabo
func (p *Partida) NoAcabada() bool {
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
	return p.calcPtsFaltaEnvido(ganadorDelEnvite)
}

// retorna la cantidad de puntos que corresponden al Falta-Envido
func (p *Partida) calcPtsFaltaEnvido(ganadorDelEnvite Equipo) int {
	// si el que va ganando:
	// 		esta en Malas -> el ganador del envite (`ganadorDelEnvite`) gana el chico
	// 		esta en Buenas -> el ganador del envite (`ganadorDelEnvite`) gana lo que le falta al maximo para ganar la ronda

	if p.estaEnMalas(p.elQueVaGanando()) {
		loQueLeFaltaAlGANADORparaGanarElChico := p.elChico() - p.puntajes[ganadorDelEnvite]
		return loQueLeFaltaAlGANADORparaGanarElChico
	}
	//else {
	loQueLeFaltaAlQUEvaGANANDOparaGanarElChico := p.puntuacion.toInt() - p.puntajes[p.elQueVaGanando()]
	return loQueLeFaltaAlQUEvaGANANDOparaGanarElChico
	//}

}

// NuevaPartida retorna nueva partida; error si hubo
func NuevaPartida(puntuacion Puntuacion, equipoAzul, equipoRojo []string) (*Partida, error) {

	mismaCantidadDeJugadores := len(equipoRojo) == len(equipoAzul)
	cantJugadores := len(equipoRojo) + len(equipoAzul)
	cantidadCorrecta := contains([]int{2, 4, 6}, cantJugadores) // puede ser 2, 4 o 6
	ok := mismaCantidadDeJugadores && cantidadCorrecta
	if !ok {
		return nil, fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
	}
	// paso a crear los jugadores; intercalados
	var jugadores []Jugador
	// para cada rjo que agrego; le agrego tambien su mano
	for i := range equipoRojo {
		nuevoJugadorRojo := Jugador{equipoRojo[i], Rojo}
		nuevoJugadorAzul := Jugador{equipoAzul[i], Azul}
		jugadores = append(jugadores, nuevoJugadorAzul, nuevoJugadorRojo)
	}

	partida := Partida{
		puntuacion:    puntuacion,
		puntaje:       0,
		cantJugadores: cantJugadores,
		jugadores:     jugadores,
	}

	partida.puntajes[Rojo] = 0
	partida.puntajes[Azul] = 0

	//partida.dobleLinking()

	partida.Ronda = nuevaRonda(partida.jugadores)

	partida.sigJugada = make(chan string, 1)

	go func() {
		for {
			sjugada, sjugador := partida.getSigJugada()
			sjugada.hacer(&partida, sjugador)
		}
	}()

	return &partida, nil
}
