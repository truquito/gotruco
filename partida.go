package truco

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	quit       chan bool    = make(chan bool, 1)
	wait       chan bool    = make(chan bool, 1)
	sigJugada  chan IJugada = make(chan IJugada, 1)
	sigComando chan string  = make(chan string, 1)
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
	Azul Equipo = 0
	Rojo Equipo = 1
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
}

func (p *Partida) readLnJugada() error {
	return nil
}

// SetSigJugada nexo capa presentacion con capa logica
func (p *Partida) SetSigJugada(cmd string) error {
	// checkeo de sintaxis
	ok := regexp.MustCompile(`^(\w|-)+\s(\w|-)+\n?$`).MatchString(cmd)
	if !ok {
		return fmt.Errorf("Comando incorrecto")
	}

	sigComando <- cmd
	return nil
}

// devuelve solo la siguiente jugada VALIDA
// si no es valida es como si no hubiese pasado nada
func (p *Partida) getSigJugada() IJugada {
	var (
		iJugada IJugada
		valid   bool
	)
	for {
		iJugada, valid = <-sigJugada
		if !valid {
			quit <- true
		} else if iJugada == nil {
			wait <- true
		} else {
			break
		}
	}
	return iJugada
}

func (p *Partida) parseJugada(jugadaStr, jugadorStr string) (IJugada, error) {
	var jugada IJugada

	manojo, err := p.Ronda.getManojo(jugadorStr)
	if err != nil {
		return nil, fmt.Errorf("Usuario %s no encontrado", jugadorStr)
	}

	jugadaStr = strings.ToLower(jugadaStr)

	switch jugadaStr {
	// toques
	case "envido":
		jugada = tocarEnvido{Jugada{autor: manojo}}
	case "real-envido":
		jugada = tocarRealEnvido{Jugada{autor: manojo}}
	case "falta-envido":
		jugada = tocarFaltaEnvido{Jugada{autor: manojo}}

	// cantos
	case "flor":
		jugada = cantarFlor{Jugada{autor: manojo}}
	case "contra-flor":
		jugada = cantarContraFlor{Jugada{autor: manojo}}
	case "contra-flor-al-resto":
		jugada = cantarContraFlorAlResto{Jugada{autor: manojo}}

	// gritos
	case "truco":
		jugada = gritarTruco{Jugada{autor: manojo}}
	case "re-truco":
		jugada = gritarReTruco{Jugada{autor: manojo}}
	case "vale-4":
		jugada = gritarVale4{Jugada{autor: manojo}}

	// respuestas
	case "quiero":
		jugada = responderQuiero{Jugada{autor: manojo}}
	case "no-quiero":
		jugada = responderNoQuiero{Jugada{autor: manojo}}
	case "tiene":
		jugada = responderNoQuiero{Jugada{autor: manojo}}

	// acciones
	case "mazo":
		jugada = irseAlMazo{Jugada{autor: manojo}}
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

	p := Partida{
		puntuacion:    puntuacion,
		puntaje:       0,
		cantJugadores: cantJugadores,
		jugadores:     jugadores,
	}

	p.puntajes[Rojo] = 0
	p.puntajes[Azul] = 0

	p.Ronda = nuevaRonda(p.jugadores)

	go func() {
		for {
			sjugada := p.getSigJugada()
			sjugada.hacer(&p)
		}
	}()

	go func() {
		for {
			select {
			// este canal agarra solo los comandos en forma de string
			// luego lo pasa al otro canal de jugadas ya aceptadas
			// en la que espera la parte interna del codigo
			case cmd := <-sigComando:
				switch cmd {
				case "__TERMINAR__":
					close(sigJugada)
				case "__WAIT__":
					sigJugada <- nil
				default:
					params := strings.Fields(cmd)
					jugadaStr, jugadorStr := params[1], params[0]
					jugada, err := p.parseJugada(jugadaStr, jugadorStr)
					if err != nil {
						fmt.Println(err.Error())
					} else {
						sigJugada <- jugada
					}
				}

				// case <-p.quit:
				// case <-time.After(1 * time.Second):
				// default:
			}
		}
	}()

	return &p, nil
}

// Terminar espera a que se consuma toda la fila de jugadas
// si se quisiera terminar abruptamente se deberia
// usar otro canal tipo `p.quit<-true` y agregarle
// el caso que corresponda al `select{...}`
func (p *Partida) Terminar() {
	sigComando <- "__TERMINAR__"
	<-quit
}

// Esperar espera a que se consuma toda la fila de jugadas
// para continuar; pero sin cerrar ningun canal
func (p *Partida) Esperar() {
	sigComando <- "__WAIT__"
	<-wait
}
