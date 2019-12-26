package truco

import (
	"fmt"
)

// IJugada Interface para las jugadas
type IJugada interface {
	hacer(p *Partida, j *Jugador) error
}

// PRE: supongo que el jugador que toca este envido
// no tiene flor (es checkeada cuando es su turno)
type tocarEnvido struct{}

func (jugada tocarEnvido) hacer(p *Partida, j *Jugador) error {
	e := &p.Ronda.envido
	envidoHabilitado := e.estado == NOCANTADOAUN || e.estado == ENVIDO
	// no es necesario porque (todo:) cuando termina
	// la primera mano deberia de invalidar el envido
	// automaticamente
	// esPrimeraMano 		:= p.Ronda.manoEnJuego == primera
	ok := envidoHabilitado // && esPrimeraMano
	if !ok {
		return fmt.Errorf(`No es posible cantar 'Envido'`)
	}

	e.cantadoPor = j
	fmt.Printf(">> %s toca envido\n", j.nombre)
	// ahora checkeo si alguien tiene flor
	hayFlor, jFlor := p.Ronda.checkFlores()
	if hayFlor {
		p.Ronda.envido.estado = DESHABILITADO
		p.Ronda.flor = FLOR
		// Se cachea turno actual (del envido).
		// Cuando se termine de jugar la flor,
		// se reestablece a este.
		cacheTurnoEnvido := p.Ronda.turno
		nuevoTurnoFlor, _ := obtenerIdx(jFlor[0], p.jugadores)
		p.Ronda.turno = nuevoTurnoFlor
		siguienteJugada := cantarFlor{}
		siguienteJugada.hacer(p, jFlor[0])
		// una vez terminada, vuelve el turno al del envido
		p.Ronda.turno = cacheTurnoEnvido
	} else {
		// 2 opciones: o bien no se jugo aun
		// o bien ya estabamos en envido
		if e.estado == ENVIDO {
			// se aumenta el puntaje del envido en +2
			e.puntaje += 2
		} else if e.estado == NOCANTADOAUN { // no se habia jugado aun
			e.estado = ENVIDO
			e.puntaje = 2
		}
		// esperando respuestas
		cacheTurnoEnvido := p.Ronda.turno
		p.readLnJugada() // se juega la respuesta
		p.Ronda.turno = cacheTurnoEnvido
	}

	return nil
}

// donde 'j' el jugador que dijo 'quiero' al 'envido'/'real envido'
func (jugada tocarEnvido) eval(p *Partida, j *Jugador) error {
	p.Ronda.envido.estado = DESHABILITADO
	jIdx, max, out := p.Ronda.getElEnvido()
	print(out)
	jug := &p.jugadores[jIdx]
	p.puntajes[jug.equipo] += p.Ronda.envido.puntaje
	fmt.Printf(`>> El envido lo gano %s con %v, +%v puntos
	para el equipo %s`+"\n",
		jug.nombre, max, p.Ronda.envido.puntaje, jug.equipo)

	return nil
}

type tocarRealEnvido struct{}

func (jugada tocarRealEnvido) hacer(p *Partida, j *Jugador) error {
	e := &p.Ronda.envido
	realEnvidoHabilitado := e.estado == NOCANTADOAUN || e.estado == ENVIDO
	ok := realEnvidoHabilitado // && esPrimeraMano
	if ok {
		e.estado = REALENVIDO
		e.puntaje += 3
		e.cantadoPor = j
		fmt.Printf(">> %s toca real envido\n", j.nombre)
		// ahora checkeo si alguien tiene flor
		hayFlor, jFlor := p.Ronda.checkFlores()
		if hayFlor {
			p.Ronda.envido.estado = DESHABILITADO
			p.Ronda.flor = FLOR
			// Se cachea turno actual (del envido).
			// Cuando se termine de jugar la flor,
			// se reestablece a este.
			cacheTurnoEnvido := p.Ronda.turno
			nuevoTurnoFlor, _ := obtenerIdx(jFlor[0], p.jugadores)
			p.Ronda.turno = nuevoTurnoFlor
			siguienteJugada := cantarFlor{}
			siguienteJugada.hacer(p, jFlor[0])
			// una vez terminada, vuelve el turno al del envido
			p.Ronda.turno = cacheTurnoEnvido
		} else {
			// 2 opciones:
			// o bien el envido no se jugo aun,
			// o bien ya estabamos en envido
			if e.estado == ENVIDO {
				// se aumenta el puntaje del envido en +2
				e.puntaje += 3
			} else if e.estado == NOCANTADOAUN { // no se habia jugado aun
				e.estado = REALENVIDO
				e.puntaje = 3
			}
			// esperando respuestas
			cacheTurnoEnvido := p.Ronda.turno
			p.readLnJugada() // se juega la respuesta
			p.Ronda.turno = cacheTurnoEnvido
		}
	}
	return nil
}

type tocarFaltaEnvido struct{}

func (jugada tocarFaltaEnvido) hacer(p *Partida, j *Jugador) error {
	// si ambos jugadores estan en malas:
	//  x = lo que le falta al ganador para ganar
	// (es decir, el ganador gana la partida)
	// si (al menos) uno de los dos paso:
	//  x = lo que le falta a el que va ganando para ganar la partida (completar las buenas)
	e := &p.Ronda.envido
	faltaEnvidoHabilitado := e.estado >= NOCANTADOAUN && e.estado < FALTAENVIDO
	ok := faltaEnvidoHabilitado // && esPrimeraMano
	if ok {
		e.estado = FALTAENVIDO
		e.cantadoPor = j
		fmt.Printf(">> %s toca falta envido\n", j.nombre)
		// ahora checkeo si alguien tiene flor
		hayFlor, jFlor := p.Ronda.checkFlores()
		if hayFlor {
			p.Ronda.envido.estado = DESHABILITADO
			p.Ronda.flor = FLOR
			// Se cachea turno actual (del envido).
			// Cuando se termine de jugar la flor,
			// se reestablece a este.
			cacheTurnoEnvido := p.Ronda.turno
			nuevoTurnoFlor, _ := obtenerIdx(jFlor[0], p.jugadores)
			p.Ronda.turno = nuevoTurnoFlor
			siguienteJugada := cantarFlor{}
			siguienteJugada.hacer(p, jFlor[0])
			// una vez terminada, vuelve el turno al del envido
			p.Ronda.turno = cacheTurnoEnvido
		} else {
			// esperando respuestas
			cacheTurnoEnvido := p.Ronda.turno
			p.readLnJugada() // se juega la respuesta
			p.Ronda.turno = cacheTurnoEnvido
		}
	}
	// si no, esta tocando `Falta Envido` al pedo
	return fmt.Errorf(`No es posible cantar 'Falta Envido' 
	porque el 'envido' esta deshabilitado`)
}

// siendo j el jugador que dijo 'quiero' a la 'falta envido'

/**
 * forma actual de jugar:
 *		si estan en malas: el que gana el envido gana la partida.
				terminando asi la partida.
 *		si no: se juega por el resto del maximo puntaje
				no necesariamente terminando asi la partida.
 * forma alternativa:
 *		si estan en malas: se juega por completar las malas
 *		si no: se juega por el resto del maximo puntaje
*/

func (jugada tocarFaltaEnvido) eval(p *Partida, j *Jugador) error {
	p.Ronda.envido.estado = DESHABILITADO

	// computar envidos
	jIdx, max, out := p.Ronda.getElEnvido()

	print(out)

	// jug es el que gano el (falta) envido
	jug := &p.jugadores[jIdx]

	puntajeEnJuego := 0
	puntuacion := p.puntuacion.toInt()
	puntuacionMalas := p.getPuntuacionMalas()
	maxPuntaje := p.getMaxPuntaje()
	estanEnMalas := maxPuntaje < puntuacionMalas
	if estanEnMalas {
		// el que tiene el maximo envido ('jug')
		// gano la partida
		puntajeEnJuego = puntuacion - p.puntajes[jug.equipo]
	} else {
		// en caso contrario, la falta envido se juega
		// por los puntos que le falta al equipo que
		// va ganando
		puntajeEnJuego = puntuacion - maxPuntaje
	}

	p.puntajes[jug.equipo] += puntajeEnJuego
	fmt.Printf(`>> La falta envido la gano %s con %v, +%v puntos
	para el equipo %s`+"\n",
		jug.nombre, max, p.Ronda.envido.puntaje, jug.equipo)

	return nil
}

type cantarFlor struct{}

func (jugada cantarFlor) hacer(p *Partida, j *Jugador) error {
	// j dice que puede cantar flor;
	// es esto verdad?
	florHabilitada := (p.Ronda.flor == NOCANTADA || p.Ronda.flor == FLOR) && p.Ronda.manoEnJuego == primera
	tieneFlor, _ := j.manojo.tieneFlor(p.Ronda.muestra)
	ok := florHabilitada && tieneFlor
	if !ok {
		return fmt.Errorf(`No es posible cantar flor`)
	}
	// todo: en ningun momento se usa la variable cantadoPor
	// solo se setea
	// e.cantadoPor = j
	fmt.Printf(">> %s canta flor\n", j.nombre)
	p.Ronda.envido.estado = DESHABILITADO
	p.Ronda.flor = FLOR
	// ahora checkeo si alguien tiene flor
	// retorna TODOS los jugadores que tengan flor (si es que existen)
	// aPartirDe, _ := obtenerIdx(j, p.jugadores)
	hayFlor, jugadoresConFlor := p.Ronda.checkFlores()
	// creo una copia
	jugadoresConFlorCACHE := make([]*Jugador, len(jugadoresConFlor))
	copy(jugadoresConFlorCACHE, jugadoresConFlor)

	if hayFlor {
		// entonces tengo que esperar respuesta SOLO de alguno de ellos;
		// a menos de un "Me voy al mazo; esa tambien es aceptada"
		// las otras las descarto
		// si no recibo respuesta en menos de x tiempo la canto yo
		// por ellos

		// Se cachea turno actual (del que canto flor).
		// Cuando se termine de jugar la flor,
		// se reestablece a este.
		cacheTurnoFlor := p.Ronda.turno

		// esto de los turnos es al pedo
		// ahora deberia de hacerce con los checkers
		nuevoTurnoFlor, _ := obtenerIdx(jugadoresConFlor[0], p.jugadores)
		p.Ronda.turno = nuevoTurnoFlor

		todosLosJugadoresConFlorCantaron := false
		for !todosLosJugadoresConFlorCantaron {

			jugada, jugador := p.getSigJugada()
			esAlguienDelQueEspero := contains(jugadoresConFlor, jugador)

			_, esMeVoyAlMazo := jugada.(irseAlMazo)
			_, esCantoFlor := jugada.(cantarFlor)
			_, esCantoContraFlor := jugada.(cantarContraFlor)
			_, esCantoContraFlorAlResto := jugada.(cantarContraFlorAlResto)
			_, esCantoConFlorMeAchico := jugada.(cantarConFlorMeAchico)
			esTipoFlor := esCantoFlor || esCantoContraFlor || esCantoContraFlorAlResto || esCantoConFlorMeAchico
			_, esQuiero := jugada.(responderQuiero)
			_, esNoQuiero := jugada.(responderNoQuiero)
			esRespuesta := esQuiero || esNoQuiero
			seEstaJugandoLaContraFlorAlResto := p.Ronda.flor == CONTRAFLORALRESTO

			esDeAlguienQueNoEsperoYNoEsIrseAlMazo := !esAlguienDelQueEspero && !esMeVoyAlMazo
			esDeAlguienQueEsperoPeroNoEsNiFlorNiIrseAlMazo := esAlguienDelQueEspero && (!seEstaJugandoLaContraFlorAlResto && !(esTipoFlor || esMeVoyAlMazo)) || (seEstaJugandoLaContraFlorAlResto && !esRespuesta)

			noEsValida := esDeAlguienQueNoEsperoYNoEsIrseAlMazo || esDeAlguienQueEsperoPeroNoEsNiFlorNiIrseAlMazo

			// todo: que pasa si solo falta 1 por responder y se va al mazo?

			if noEsValida {
				// no deberia de salir de este loop
				// pero solo responderle al loco (?)
				return fmt.Errorf(`No es el momento de realizar
					esta jugada; ahora estoy esperando por cantos de flor (de
					aquellos que la poseen) o bien "Irse al mazo" (de cualquier jugador)`)
			}

			// solo queda 3 casos posibles:
			// CASO I: 	esEsperado & esFlor
			// CASO II:	esEsperado & esMazo
			// CASO III: 	!esEsperado & esMazo

			if esAlguienDelQueEspero {
				// lo descuento de los esperados
				jugadoresConFlor = eliminar(jugadoresConFlor, jugador)
				// era el ultimo que del que me faltaba escuchar?
				// y por ende -> fin del bucle ?
			}

			// la ejecuto porque por descarte ya se que es valida
			if esMeVoyAlMazo {
				jugada.hacer(p, jugador)

			} else if esCantoFlor {
				// ya se que estaba habilitado para cantar flor
				// porque estaba en `jugadoresConFlor`
				// ahora: se canto contraflor o mayor -> inhabilitado
				florHabilitada := p.Ronda.flor == FLOR
				if !florHabilitada {
					// entonces lo vuelvo a agregar a la lista de esperados; se equivoco
					jugadoresConFlor = append(jugadoresConFlor, jugador)
					return fmt.Errorf(`Ya no es posible cantar flor;`)
				}
				// en caso contrario; esta todo bien;
				// la canta
				fmt.Printf(">> %s canta flor\n", jugador.nombre)
				// ahora la flor pasa a jugarse por +3 puntos
				p.Ronda.envido.puntaje += 3

			} else if esCantoContraFlor {
				// ya se que estaba habilitado para cantar flor
				// porque estaba en `jugadoresConFlor`
				// ahora: se canto contraflor o algo asi? si si -> inhabilitado
				contraFlorHabilitada := p.Ronda.flor == FLOR
				if !contraFlorHabilitada {
					// entonces lo vuelvo a agregar a la lista de esperados; se equivoco
					jugadoresConFlor = append(jugadoresConFlor, jugador)
					return fmt.Errorf(`Ya no es posible cantar contra flor;`)
				}
				// en caso contrario; esta todo bien;
				// la canta
				fmt.Printf(">> %s canta contra-flor\n", jugador.nombre)
				p.Ronda.flor = CONTRAFLOR
				// ahora la flor pasa a jugarse por 4 puntos
				p.Ronda.envido.puntaje = 4
				// y ahora tengo que esperar por la respuesta de la nueva
				// propuesta de todos menos de el que canto la contraflor
				// restauro la copia
				jugadoresConFlor = make([]*Jugador, len(jugadoresConFlorCACHE))
				copy(jugadoresConFlor, jugadoresConFlorCACHE)
				// lo elimino de los que espero
				jugadoresConFlor = eliminar(jugadoresConFlor, jugador)

			} else if esCantoContraFlorAlResto {
				// ya se que estaba habilitado para cantar flor
				// porque estaba en `jugadoresConFlor`
				// ahora: puede cantarContraFlorAlResto?
				contraFlorAlRestoHabilitada := p.Ronda.flor == FLOR || p.Ronda.flor == CONTRAFLOR
				if !contraFlorAlRestoHabilitada {
					// entonces lo vuelvo a agregar a la lista de esperados; se equivoco
					// por ejemplo; ya otro jugador habia cantado contraFlorAlResto
					jugadoresConFlor = append(jugadoresConFlor, jugador)
					return fmt.Errorf(`Ya no es posible cantar contra flor al resto;`)
				}
				// en caso contrario; esta todo bien;
				// la canta
				fmt.Printf(">> %s canta contra-flor-al-resto\n", jugador.nombre)
				p.Ronda.flor = CONTRAFLORALRESTO
				p.Ronda.envido.cantadoPor = j
				// los puntos de la flor quedan acumulados
				// y ahora tengo que esperar por la respuesta de la nueva
				// propuesta de todos menos de el que canto la contraflor
				// restauro la copia
				jugadoresConFlor = make([]*Jugador, len(jugadoresConFlorCACHE))
				copy(jugadoresConFlor, jugadoresConFlorCACHE)
				// lo elimino de los que espero
				jugadoresConFlor = eliminar(jugadoresConFlor, jugador)

			} else if esQuiero && seEstaJugandoLaContraFlorAlResto {
				// solo con que uno *DEL EQUIPO CONTRARIO*
				// al que canto la contra-flor-al-resto diga quiero
				// es del equipo contrario?
				esDelEquipoContrario := jugador.equipo != p.Ronda.envido.cantadoPor.equipo
				if !esDelEquipoContrario {
					return fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
				}
				fmt.Printf(">> %s dice quiero \n", jugador.nombre)
				// ok; se cierra el envite; hora de calcular el ganador
				p.Ronda.flor = DESHABILITADA
				manojoConLaFlorMasAlta, maxFlor := p.Ronda.getLaFlorMasAlta()
				equipoGanador := manojoConLaFlorMasAlta.jugador.equipo
				// ahora se quien es el ganador; necesito saber cuantos puntos
				// se le va a sumar a ese equipo:
				// los acumulados del envite hasta ahora + la contrafloralresto
				puntosASumar := p.Ronda.envido.puntaje + p.calcPtsContraFlorAlResto(equipoGanador)
				p.puntajes[equipoGanador] += puntosASumar
				fmt.Printf(`>> La contra-flor-al-resto la gano %s con %v, +%v puntos
				para el equipo %s`+"\n",
					manojoConLaFlorMasAlta.jugador.nombre, maxFlor, puntosASumar, equipoGanador)
				// se corta el bucle de la flor:
				break

			} else if esNoQuiero && seEstaJugandoLaContraFlorAlResto {
				// solo con que uno *DEL EQUIPO CONTRARIO*
				// al que canto la contra-flor-al-resto diga quiero
				// es del equipo contrario?
				esDelEquipoContrario := jugador.equipo != p.Ronda.envido.cantadoPor.equipo
				if !esDelEquipoContrario {
					return fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
				}
				// ok; se cierra el envite; los puntos van para el que propuso el envite
				p.Ronda.flor = DESHABILITADA
				equipoGanador := p.Ronda.envido.cantadoPor.equipo
				// ahora se quien es el ganador; necesito saber cuantos puntos
				// se le va a sumar a ese equipo:
				// los acumulados del envite hasta ahora + la contrafloralresto
				puntosASumar := p.Ronda.envido.puntaje + p.calcPtsContraFlorAlResto(equipoGanador)
				p.puntajes[equipoGanador] += puntosASumar
				fmt.Printf(`>> La contra-flor-al-resto la gano %s, +%v puntos
				para el equipo %s`+"\n",
					p.Ronda.envido.cantadoPor.nombre, puntosASumar, equipoGanador)
				// se corta el bucle de la flor:
				break
			}

			todosLosJugadoresConFlorCantaron = len(jugadoresConFlor) == 0

		}

		// una vez terminada, vuelve el turno al del envido
		p.Ronda.turno = cacheTurnoFlor

		return nil

	}

	//else {
	// Nadie mas tiene flor; entonces j se lleva todos
	// los puntos en juego (+3)
	p.puntajes[j.equipo] += 3
	fmt.Printf(`>> +%v puntos para el equipo %s`+"\n",
		3, j.equipo)
	return nil
	//}

}

type cantarContraFlor struct{}

func (jugada cantarContraFlor) hacer(p *Partida, j *Jugador) error {
	return nil
}

type cantarContraFlorAlResto struct{}

func (jugada cantarContraFlorAlResto) hacer(p *Partida, j *Jugador) error {
	return nil
}

type cantarConFlorMeAchico struct{}

func (jugada cantarConFlorMeAchico) hacer(p *Partida, j *Jugador) error {
	return nil
}

func esCantoFlor(jugada IJugada) bool {
	var esCantoFlor bool = false
	switch jugada.(type) {
	case cantarFlor, cantarContraFlor, cantarContraFlorAlResto, cantarConFlorMeAchico:
		esCantoFlor = true
	default:
		esCantoFlor = false
	}
	return esCantoFlor
}

type gritarTruco struct{}

func (jugada gritarTruco) hacer(p *Partida, j *Jugador) error {
	return nil
}

type gritarReTruco struct{}

func (jugada gritarReTruco) hacer(p *Partida, j *Jugador) error {
	return nil
}

type gritarVale4 struct{}

func (jugada gritarVale4) hacer(p *Partida, j *Jugador) error {
	return nil
}

type responderQuiero struct{}

func (jugada responderQuiero) hacer(p *Partida, j *Jugador) error {
	// se acepta una respuesta 'quiero' solo cuando:
	// - CASO I: se toco el envido (o similar)
	// - CASO II: se grito el truco (o similar)
	// en caso contrario, es incorrecto -> error

	// CASO I: se toco el envido (o similar)
	elEnvidoEsRespondible := p.Ronda.envido.estado >= ENVIDO
	if elEnvidoEsRespondible {
		fmt.Printf(">> %s responde quiero\n", j.nombre)
		if p.Ronda.envido.estado == FALTAENVIDO {
			return tocarFaltaEnvido{}.eval(p, j)
		}
		// si no, era envido/real-envido o cualquier
		// combinacion valida de ellos
		return tocarEnvido{}.eval(p, j)
	}

	// CASO II: se grito truco
	elTrucoEsRespondible := p.Ronda.truco >= TRUCO
	if elTrucoEsRespondible {

	}

	// si no, esta respondiendo al pedo
	return fmt.Errorf(`No hay nada \"que querer\"; ya que: el 
	estado del envido no es "envido" (o mayor) y el estado del 
	truco no es "truco" (o mayor)`)
}

type responderNoQuiero struct{}

func (jugada responderNoQuiero) hacer(p *Partida, j *Jugador) error {
	// se acepta una respuesta 'no quiero' solo cuando:
	// - CASO I: se toco el envido (o similar)
	// - CASO II: se grito el truco (o similar)
	// en caso contrario, es incorrecto -> error

	// CASO I: se toco el envido (o similar)
	e := &p.Ronda.envido
	elEnvidoEsRespondible := e.estado >= ENVIDO
	if elEnvidoEsRespondible {
		fmt.Printf(">> %s responde no quiero\n", j.nombre)

		// se pasa a calcular el puntaje correspondiente:
		// si se canto envido solo 1 vez: total = 1 pts
		// si se canto real envido solo 1 vez: total = 1 pts
		// si se canto falta envido solo 1 vez: total = 1 pts
		// conclusion: si Envido.puntaje <= 3 -> total = 1 pts
		// si no:
		//	no se toma en cuenta el puntaje total del ultimo toque
		// 	~ se resta de 'Envido.puntaje' el puntaje correspondiente
		//		del ultmo toque:
		//			-2pts si el ultimo toque fue envido
		//			-3pts si el ultmo toque fue real envido
		// 			-0pts si el ultimo toque fue falta envido

		var totalPts int

		if e.puntaje <= 3 {
			totalPts = 1
			// fix caso especial
			fix := e.estado == FALTAENVIDO && e.puntaje > 2
			if fix {
				totalPts = e.puntaje
			}
		} else {
			switch e.estado {
			case ENVIDO:
				totalPts = e.puntaje - 2
			case REALENVIDO:
				totalPts = e.puntaje - 3
			case FALTAENVIDO:
				totalPts = e.puntaje
			}
		}

		cantadoPor := e.cantadoPor
		e.estado = DESHABILITADO
		e.puntaje = totalPts
		p.puntajes[cantadoPor.equipo] += totalPts
		fmt.Printf(`>> +%v puntos para el equipo %s`+"\n",
			totalPts, cantadoPor.equipo)

	}

	// CASO II: se grito truco
	elTrucoEsRespondible := p.Ronda.truco >= TRUCO
	if elTrucoEsRespondible {

	}

	// si no, esta respondiendo al pedo
	return fmt.Errorf(`%s esta respondiendo al pedo; no hay 
	nada respondible`, j.nombre)
}

type irseAlMazo struct{}

func (jugada irseAlMazo) hacer(p *Partida, j *Jugador) error {
	return nil
}

var jugadas = map[string]([]string){
	"Gritos": []string{
		"Truco",    // 2pts // el-segundo
		"Re-Truco", // 3 pts
		"Vale 4",   // 4 pts
	},
	"Toques": []string{
		"Envido",
		"Real envido",
		"Falta envido",
	},
	"Cantos": []string{
		"Flor",                 // 2pts (tanto o el-primero)
		"Contra flor",          // 3 pts
		"Contra flor al resto", // 4 pts

		//"Con flor me achico",
		//"Con flor quiero",
	},
	"Respuestas": []string{
		"Quiero",
		"No-Quiero",
	},
	"Acciones": []string{
		"Irse al mazo",
	},
}

// ImprimirJugadas imprime las jugadas posibles
func ImprimirJugadas() {
	for tipoJugada, opciones := range jugadas {
		fmt.Printf("%s: ", tipoJugada)
		for _, jugada := range opciones {
			fmt.Printf("%s, ", jugada)
		}
		fmt.Printf("\n")
	}
	fmt.Println()
}
