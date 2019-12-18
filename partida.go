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
	case a20: return 20
	case a30: return 30
	default:	return 40		
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

// nuevaRonda : crea una nueva ronda al azar
func (p *Partida) nuevaRonda() {
	p.ronda = Ronda{
		manoEnJuego: primera,
		elMano:      0,
		turno:       0,
		envido:      Envido{puntaje: 0, estado: NOCANTADOAUN},
		truco:       NOCANTADO,
		manojos:     make([]Manojo, p.cantJugadores),
		manos:       make([]Mano, 3),
	}

	// reparto 3 cartas al azar a cada jugador
	// y ademas una muestra, tambien al azar.
	dealCards(&p.ronda.manojos, &p.ronda.muestra)

	// hago el doble-linking "jugadores <-> manojos"
	for i := 0; i < p.cantJugadores; i++ {
		p.jugadores[i].manojo = &p.ronda.manojos[i]
		p.ronda.manojos[i].jugador = &p.jugadores[i]
	}

	// seteo el repartidor de la primera mano como
	// el mano de la ronda (segun las reglas)
	p.ronda.getManoActual().repartidor = p.ronda.elMano
	// p.ronda.setTurno()
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

func dobleLinking(p *Partida) {
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
	esElUltimo := int(j) == p.cantJugadores - 1
	if esElUltimo { 
		return 0
	}
	return j + 1
}

/**
* getElEnvido computa el envido de la ronda
* @return `jIdx JugadorIdx` Es el indice del jugador con
* el envido mas alto (i.e., ganador)
* @return `max int` Es el valor numerico del maximo envido
* @return `stdOut []string` Es conjunto ordenado de todos
* los mensajes que normalmente se escucharian en una ronda
* de envido en la vida real.
* e.g.:
* 	`stdOut[0] = Jacinto dice: "tengo 9"`
*   `stdOut[1] = Patricio dice: "son buenas" (tenia 3)`
*   `stdOut[2] = Pedro dice: "30 son mejores!"`
*		`stdOut[3] = Juan dice: "33 son mejores!"`
*
* NOTA: todo: Eventualmente se cambiaria []string por algo
* "mas serializable" para usar con el front-end 
* e.g., []{JugadorIdx, string} donde `string` no deberia de
* contener cosas como "tenia 3". O tal vez un hibrido de 
* ambos con un parametro-flag que decida:
* si el juego esta en "modo json" o "modo consola"
*/
func (p *Partida) getElEnvido() (jIdx JugadorIdx, 
	max int, stdOut []string) {
	
	// decir envidos en orden segun las reglas:
	// empieza la mano
	// canta el siguiente en sentido anti-horario sii
	// tiene MAS pts que el maximo actual y es de equipo
	// contrario. de no ser asi: o bien "pasa" o bien dice
	// "son buenas" dependiendo del caso
	// asi hasta terminar una ronda completa sin decir nada
	
	// calculo y cacheo todos los envidos
	envidos := make([]int, p.cantJugadores)
	for i := range envidos {
		envidos[i] = p.jugadores[i].manojo.calcularEnvido(p.ronda.muestra)
	}
	
	// `yaDijeron` indica que jugador ya "dijo"
	// si tenia mejor, o peor envido. Por lo tanto,
	// ya no es tenido en cuenta.
	yaDijeron := make([]bool, p.cantJugadores)
	// `jIdx` indica el jugador con el envido mas alto
	// var jIdx JugadorIdx
	
	// empieza la mano
	jIdx = p.ronda.elMano
	yaDijeron[jIdx] = true
	out := fmt.Sprintf("   %s dice: \"tengo %v\"\n", p.jugadores[jIdx].nombre,
	 envidos[jIdx])
	stdOut = append(stdOut, out)

	// `todaviaNoDijeronSonMejores` se usa para
	// no andar repitiendo "son bueanas" "son buenas"
	// por cada jugador que haya jugado "de callado"
	// y ahora resulte tener peor envido.
	// agiliza el juego, de forma tal que solo se
	// escucha decir "xx son mejores", "yy son mejores"
	// "zz son mejores"
	todaviaNoDijeronSonMejores := true

	// iterador
	i := p.ronda.elMano + 1

	// termina el bucle cuando se haya dado
	// "una vuelta completa" de:mano+1 hasta:mano
	// ergo, cuando se "resetea" el iterador,
	// se setea a `p.ronda.elMano + 1`
	for i != p.ronda.elMano {
			todaviaEsTenidoEnCuenta := !yaDijeron[i]
			if todaviaEsTenidoEnCuenta {
				
				esDeEquipoContrario := p.jugadores[i].equipo != p.jugadores[jIdx].equipo
				tieneEnvidoMasAlto 	:= envidos[i] > envidos[jIdx]
				tieneEnvidoIgual 		:= envidos[i] == envidos[jIdx]
				leGanaDeMano 				:= leGanaDeMano(i, jIdx, p.ronda.elMano, p.cantJugadores)
				sonMejores 					:= tieneEnvidoMasAlto || (tieneEnvidoIgual && leGanaDeMano)
				
				if sonMejores {
					if esDeEquipoContrario {
						out := fmt.Sprintf("   %s dice: \"%v son mejores!\"\n", 
						p.jugadores[i].nombre, envidos[i])
						stdOut = append(stdOut, out)
						jIdx = i
						yaDijeron[i] = true
						todaviaNoDijeronSonMejores = false
						// se "resetea" el bucle
						i = p.sig(p.ronda.elMano)

					} else /* es del mismo equipo */ {
						// no dice nada si es del mismo equipo
						// juega de callado & sigue siendo tenido
						// en cuenta
						i = p.sig(i)
					}

				}	else /* tiene el envido mas chico */ {
					if esDeEquipoContrario {
						if todaviaNoDijeronSonMejores {
							out := fmt.Sprintf("   %s dice: \"son buenas\" (tenia %v)\n", 
							p.jugadores[i].nombre, envidos[i])
							stdOut = append(stdOut, out)
							yaDijeron[i] = true
							// pasa al siguiente
						}
						i = p.sig(i)
					} else {
						// es del mismo equipo pero tiene un envido
						// mas bajo del que ya canto su equipo.
						// ya no lo tengo en cuenta, pero no dice nada.
						yaDijeron[i] = true
						i = p.sig(i)
					}
				}

			} else {
			// si no es tenido en cuenta,
			// simplemente pasar al siguiente
			i = p.sig(i)
			}
	} // fin bucle while

	max = envidos[jIdx]

	return jIdx, max, stdOut
}

func nuevaPartida() *Partida {
	cantJugadores := 2
	partida := Partida{
		puntuacion:    a20,
		puntaje:       0,
		cantJugadores: cantJugadores, // puede ser 2, 4 o 6
		jugadores: []Jugador{
			*NuevoJugador("Juan"),
			*NuevoJugador("Pedro"),
		},
	}
	partida.puntajes[Rojo] = 0
	partida.puntajes[Azul] = 0
	return &partida
}
