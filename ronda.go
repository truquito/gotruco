package truco

import (
	"fmt"
)

// EstadoTruco : enum
type EstadoTruco int

// enums del truco
const (
	NOCANTADO EstadoTruco = 0
	TRUCO     EstadoTruco = 1
	RETRUCO   EstadoTruco = 2
	VALE4     EstadoTruco = 3
)

type truco struct {
	cantadoPor *Manojo
	estado     EstadoTruco
}

// Ronda :
type Ronda struct {
	manoEnJuego          NumMano // Numero de mano: 1era 2da 3era
	cantJugadoresEnJuego [2]int  // Numero de jugadores que no se fueron al mazo

	/* Indices */
	elMano JugadorIdx
	turno  JugadorIdx
	pies   [2]JugadorIdx

	/* toques, gritos y cantos */
	envido Envido     // Estado del envido, la primera o la mentira
	flor   EstadoFlor // Estado del envido, la primera o la mentira
	truco  truco      // Estado del truco, la segunda o el rabón

	/* cartas */
	manojos []Manojo
	muestra Carta

	manos []Mano
}

func (r *Ronda) checkFlorDelMano() {
	tieneFlor, tipoFlor := r.getElMano().tieneFlor(r.muestra)
	if tieneFlor {
		r.getElMano().cantarFlor(tipoFlor, r.muestra)
	}
}

// devuelve `false` si la ronda se acabo
func (r *Ronda) enJuego() bool {
	return true
}

// los anteriores a `aPartirDe` (incluido este) no
// son necesarios de checkear porque ya han sido
// checkeados si tenian flor
func (r Ronda) cantarFloresSiLasHay(aPartirDe JugadorIdx) {
	for _, jugador := range r.manojos[aPartirDe+1:] {
		tieneFlor, tipoFlor := jugador.tieneFlor(r.muestra)
		if tieneFlor {
			// todo:
			tieneFlor = false
			tipoFlor = tipoFlor + 1
			// var jugada IJugada = responderNoQuiero{}
			// jugador.cantarFlor(tipoFlor, r.muestra)
			r.envido.estado = DESHABILITADO
			break
		}
	}
}

// retorna todos los manojos que tienen flor
func (r Ronda) getFlores() (hayFlor bool,
	manojosConFlor []*Manojo) {
	for i, manojo := range r.manojos {
		tieneFlor, _ := manojo.tieneFlor(r.muestra)
		if tieneFlor {
			manojosConFlor = append(manojosConFlor, &r.manojos[i])
		}
	}
	hayFlor = len(manojosConFlor) > 0
	return hayFlor, manojosConFlor
}

func (r Ronda) getElMano() *Manojo {
	return &r.manojos[r.elMano]
}

func (r Ronda) getElTurno() *Manojo {
	return &r.manojos[r.turno]
}

func (r Ronda) getManoAnterior() *Mano {
	return &r.manos[r.manoEnJuego-1]
}

func (r Ronda) getManoActual() *Mano {
	return &r.manos[r.manoEnJuego]
}

func (r Ronda) setTurno() {
	// si es la primera mano que se juega
	// entonces es el turno del mano
	if r.manoEnJuego == primera {
		r.turno = r.elMano
		// si no, es turno del ganador de
		// la mano anterior
	} else {
		r.turno = r.getManoAnterior().ganador
	}
}

// Print Imprime la informacion de la ronda
func (r Ronda) Print() {
	for i := range r.manojos {
		fmt.Printf("%s:\n", r.manojos[i].jugador.nombre)
		r.manojos[i].Print()
	}

	fmt.Printf("\nY la muestra es\n    - %s\n", r.muestra.toString())
	fmt.Printf("\nEl mano actual es: %s\nEs el turno de %s\n\n",
		r.getElMano().jugador.nombre, r.getElTurno().jugador.nombre)
}

// sig devuelve el `JugadorIdx` del
// jugador siguiente a j
func (r *Ronda) sig(j JugadorIdx) JugadorIdx {
	cantJugadores := len(r.manojos)
	esElUltimo := int(j) == cantJugadores-1
	if esElUltimo {
		return 0
	}
	return j + 1
}

// retorna el manojo con la flor mas alta en la ronda
// y su valor
// pre-requisito: hay flor en la ronda
func (r *Ronda) getLaFlorMasAlta() (*Manojo, int) {
	var (
		maxFlor     = -1
		maxIdx  int = -1
	)
	for i := range r.manojos {
		valorFlor, _ := r.manojos[i].calcFlor(r.muestra)
		if valorFlor > maxFlor {
			maxFlor = valorFlor
			maxIdx = i
		}
	}
	return &r.manojos[maxIdx], maxFlor
}

// todo: esto anda bien; es legacy; pero hacer que devuelva punteros
// no indices
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
func (r *Ronda) getElEnvido() (jIdx JugadorIdx,
	max int, stdOut []string) {

	cantJugadores := len(r.manojos)

	// decir envidos en orden segun las reglas:
	// empieza la mano
	// canta el siguiente en sentido anti-horario sii
	// tiene MAS pts que el maximo actual y es de equipo
	// contrario. de no ser asi: o bien "pasa" o bien dice
	// "son buenas" dependiendo del caso
	// asi hasta terminar una ronda completa sin decir nada

	// calculo y cacheo todos los envidos
	envidos := make([]int, cantJugadores)
	for i := range envidos {
		envidos[i] = r.manojos[i].calcularEnvido(r.muestra)
	}

	// `yaDijeron` indica que jugador ya "dijo"
	// si tenia mejor, o peor envido. Por lo tanto,
	// ya no es tenido en cuenta.
	yaDijeron := make([]bool, cantJugadores)
	// `jIdx` indica el jugador con el envido mas alto
	// var jIdx JugadorIdx

	// empieza la mano
	jIdx = r.elMano
	yaDijeron[jIdx] = true
	out := fmt.Sprintf("   %s dice: \"tengo %v\"\n", r.manojos[jIdx].jugador.nombre,
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
	i := r.elMano + 1

	// termina el bucle cuando se haya dado
	// "una vuelta completa" de:mano+1 hasta:mano
	// ergo, cuando se "resetea" el iterador,
	// se setea a `p.Ronda.elMano + 1`
	for i != r.elMano {
		todaviaEsTenidoEnCuenta := !yaDijeron[i]
		if todaviaEsTenidoEnCuenta {

			esDeEquipoContrario := r.manojos[i].jugador.equipo != r.manojos[jIdx].jugador.equipo
			tieneEnvidoMasAlto := envidos[i] > envidos[jIdx]
			tieneEnvidoIgual := envidos[i] == envidos[jIdx]
			leGanaDeMano := leGanaDeMano(i, jIdx, r.elMano, cantJugadores)
			sonMejores := tieneEnvidoMasAlto || (tieneEnvidoIgual && leGanaDeMano)

			if sonMejores {
				if esDeEquipoContrario {
					out := fmt.Sprintf("   %s dice: \"%v son mejores!\"\n",
						r.manojos[i].jugador.nombre, envidos[i])
					stdOut = append(stdOut, out)
					jIdx = i
					yaDijeron[i] = true
					todaviaNoDijeronSonMejores = false
					// se "resetea" el bucle
					i = r.sig(r.elMano)

				} else /* es del mismo equipo */ {
					// no dice nada si es del mismo equipo
					// juega de callado & sigue siendo tenido
					// en cuenta
					i = r.sig(i)
				}

			} else /* tiene el envido mas chico */ {
				if esDeEquipoContrario {
					if todaviaNoDijeronSonMejores {
						out := fmt.Sprintf("   %s dice: \"son buenas\" (tenia %v)\n",
							r.manojos[i].jugador.nombre, envidos[i])
						stdOut = append(stdOut, out)
						yaDijeron[i] = true
						// pasa al siguiente
					}
					i = r.sig(i)
				} else {
					// es del mismo equipo pero tiene un envido
					// mas bajo del que ya canto su equipo.
					// ya no lo tengo en cuenta, pero no dice nada.
					yaDijeron[i] = true
					i = r.sig(i)
				}
			}

		} else {
			// si no es tenido en cuenta,
			// simplemente pasar al siguiente
			i = r.sig(i)
		}
	} // fin bucle while

	max = envidos[jIdx]

	return jIdx, max, stdOut
}

/**
* cantarFlores computa la flor
* @return `j *Manojo` Es el ptr al manojo con
* la flor mas alta (i.e., ganador)
* @return `max int` Es el valor numerico de la flor mas alta
* @return `stdOut []string` Es conjunto ordenado de todos
* los mensajes que normalmente se escucharian en una ronda
* de flor en la vida real empezando desde jIdx
* e.g.:
* 	`stdOut[0] = Jacinto dice: "tengo 9"`
*   `stdOut[1] = Patricio dice: "son buenas" (tenia 3)`
*   `stdOut[2] = Pedro dice: "30 son mejores!"`
*		`stdOut[3] = Juan dice: "33 son mejores!"`
*
* NOTA: todo: Eventualmente se cambiaria []string por algo
* "mas serializable" para usar con el front-end
 */
func (r *Ronda) cantarFlores(aPartirDe JugadorIdx) (j *Manojo,
	max int, stdOut []string) {

	cantJugadores := len(r.manojos)

	// decir flores en orden segun las reglas:
	// empieza el autor del envite
	// canta el siguiente en sentido anti-horario sii
	// tiene MAS pts que el maximo actual y es de equipo
	// contrario. de no ser asi: o bien "pasa" o bien dice
	// "son buenas" dependiendo del caso
	// asi hasta terminar una ronda completa sin decir nada

	// calculo y cacheo todas las flores
	flores := make([]int, cantJugadores)

	// `yaDijeron` indica que jugador ya "dijo"
	// si tenia mejor, o peor envido. Por lo tanto,
	// ya no es tenido en cuenta.
	yaDijeron := make([]bool, cantJugadores)

	for i := range r.manojos {
		flores[i], _ = r.manojos[i].calcFlor(r.muestra)
		tieneFlor := flores[i] > 0
		seFueAlMazo := r.manojos[i].seFueAlMazo == false
		if tieneFlor && !seFueAlMazo {
			yaDijeron[i] = false
		} else {
			yaDijeron[i] = true
		}
	}

	// `jIdx` indica el jugador con la flor mas alta

	// empieza el del parametro
	if flores[aPartirDe] > 0 {
		yaDijeron[aPartirDe] = true
		out := fmt.Sprintf("   %s dice: \"tengo %v\"\n", r.manojos[aPartirDe].jugador.nombre,
			flores[aPartirDe])
		stdOut = append(stdOut, out)
	}

	// `todaviaNoDijeronSonMejores` se usa para
	// no andar repitiendo "son bueanas" "son buenas"
	// por cada jugador que haya jugado "de callado"
	// y ahora resulte tener peor envido.
	// agiliza el juego, de forma tal que solo se
	// escucha decir "xx son mejores", "yy son mejores"
	// "zz son mejores"
	todaviaNoDijeronSonMejores := true
	jIdx := aPartirDe
	i := r.sig(aPartirDe)

	// termina el bucle cuando se haya dado
	// "una vuelta completa" de:aPartirDe hasta:aPartirDe
	// ergo, cuando se "resetea" el iterador,
	for i != aPartirDe {
		todaviaEsTenidoEnCuenta := !yaDijeron[i]
		if todaviaEsTenidoEnCuenta {

			esDeEquipoContrario := r.manojos[i].jugador.equipo != r.manojos[jIdx].jugador.equipo
			tieneEnvidoMasAlto := flores[i] > flores[jIdx]
			tieneEnvidoIgual := flores[i] == flores[jIdx]
			leGanaDeMano := leGanaDeMano(i, jIdx, r.elMano, cantJugadores)
			sonMejores := tieneEnvidoMasAlto || (tieneEnvidoIgual && leGanaDeMano)

			if sonMejores {
				if esDeEquipoContrario {
					out := fmt.Sprintf("   %s dice: \"%v son mejores!\"\n",
						r.manojos[i].jugador.nombre, flores[i])
					stdOut = append(stdOut, out)
					jIdx = i
					yaDijeron[i] = true
					todaviaNoDijeronSonMejores = false
					// se "resetea" el bucle
					i = r.sig(aPartirDe)

				} else /* es del mismo equipo */ {
					// no dice nada si es del mismo equipo
					// juega de callado & sigue siendo tenido
					// en cuenta
					i = r.sig(i)
				}

			} else /* tiene el envido mas chico */ {
				if esDeEquipoContrario {
					if todaviaNoDijeronSonMejores {
						out := fmt.Sprintf("   %s dice: \"son buenas\" (tenia %v)\n",
							r.manojos[i].jugador.nombre, flores[i])
						stdOut = append(stdOut, out)
						yaDijeron[i] = true
						// pasa al siguiente
					}
					i = r.sig(i)
				} else {
					// es del mismo equipo pero tiene un envido
					// mas bajo del que ya canto su equipo.
					// ya no lo tengo en cuenta, pero no dice nada.
					yaDijeron[i] = true
					i = r.sig(i)
				}
			}

		} else {
			// si no es tenido en cuenta,
			// simplemente pasar al siguiente
			i = r.sig(i)
		}
	}

	max = flores[jIdx]

	return r.getManojo(jIdx), max, stdOut
}

func (r *Ronda) getManojo(jIdx JugadorIdx) *Manojo {
	return &r.manojos[jIdx]
}

func (r *Ronda) singleLinking(jugadores []Jugador) {
	cantJugadores := len(jugadores)
	for i := 0; i < cantJugadores; i++ {
		r.manojos[i].jugador = &jugadores[i]
	}
}

// todo: esto es ineficiente
// getManojo devuelve el puntero al manojo,
// dado un string que identifique al jugador duenio de ese manojo
func (r *Ronda) getManojoByStr(idJugador string) (*Manojo, error) {
	for i := range r.manojos {
		if r.manojos[i].jugador.nombre == idJugador {
			return &r.manojos[i], nil
		}
	}
	return nil, fmt.Errorf("Jugador `%s` no encontrado", idJugador)
}

func (r *Ronda) setManojos(manojos []Manojo) {
	for m, manojo := range manojos {
		for c, carta := range manojo.Cartas {
			r.manojos[m].Cartas[c] = carta
		}
	}
}

func (r *Ronda) setMuestra(muestra Carta) {
	r.muestra = muestra
}

/*
 * Reparte 3 cartas al azar a cada manojo de c/jugador
 * y 1 a la `muestra` (se las actualiza)
 */
func (r *Ronda) dealCards() {
	cantJugadores := cap(r.manojos)
	// genero `3*cantJugadores + 1` cartas al azar
	randomCards := getCartasRandom(3*cantJugadores + 1)

	for numJugador := 0; numJugador < cantJugadores; numJugador++ {
		for numCarta := 0; numCarta < 3; numCarta++ {
			cartaID := CartaID(randomCards[3*numJugador+numCarta])
			carta := nuevaCarta(cartaID)
			r.manojos[numJugador].Cartas[numCarta] = carta
		}
	}

	// la ultima es la muestra
	n := cap(randomCards)
	r.muestra = nuevaCarta(CartaID(randomCards[n-1]))
}

// nuevaRonda : crea una nueva ronda al azar
func nuevaRonda(jugadores []Jugador) Ronda {
	cantJugadores := len(jugadores)
	cantJugadoresPorEquipo := cantJugadores / 2
	ronda := Ronda{
		manoEnJuego:          primera,
		cantJugadoresEnJuego: [2]int{cantJugadoresPorEquipo, cantJugadoresPorEquipo},
		elMano:               0,
		turno:                0,
		envido:               Envido{puntaje: 0, estado: NOCANTADOAUN},
		flor:                 NOCANTADA,
		truco:                truco{cantadoPor: nil, estado: NOCANTADO},
		manojos:              make([]Manojo, cantJugadores),
		manos:                make([]Mano, 3),
	}

	// reparto 3 cartas al azar a cada jugador
	// y ademas una muestra, tambien al azar.
	ronda.dealCards()

	// // hago el SINGLE-linking "jugadores <- manojos"
	ronda.singleLinking(jugadores)

	// seteo el repartidor de la primera mano como
	// el mano de la ronda (segun las reglas)
	ronda.getManoActual().repartidor = ronda.elMano
	// p.Ronda.setTurno()

	return ronda
}
