package truco

import (
	"fmt"
)

// IJugada Interface para las jugadas
type IJugada interface {
	hacer(p *Partida) error
	getAutor() *Manojo
}

// Jugada ...
type Jugada struct {
	autor *Manojo
}

func (j Jugada) getAutor() *Manojo {
	return j.autor
}

type tirarCarta struct {
	Jugada
	Carta
}

// el jugador tira una carta;
// el parametro se encuentra en la struct como atributo
func (jugada tirarCarta) hacer(p *Partida) error {
	// primero que nada: tiene esa carta?
	idx, err := jugada.autor.getCartaIdx(jugada.Carta)
	if err != nil {
		return err
	}

	// luego, era su turno?
	eraSuTurno := p.Ronda.getElTurno() == jugada.autor
	if !eraSuTurno {
		return fmt.Errorf("No era su turno, no puede tirar la carta")
	}

	// ok la tiene y era su turno -> la juega
	fmt.Printf("%s tira la carta %s",
		jugada.autor.jugador.nombre,
		jugada.Carta.toString())
	jugada.autor.cartasNoJugadas[idx] = false
	p.Ronda.getManoActual().agregarTirada(jugada)

	// era el ultimo en tirar de esta mano?
	eraElUltimoEnTirar := p.Ronda.sigHabilitado(*jugada.autor) == nil
	if eraElUltimoEnTirar {
		// de ser asi tengo que checkear el resultado de la mano
		p.evaluarMano()
		// el turno del siguiente queda dado por el ganador de esta
	}

	// TODO:::::::::
	// paso de turno
	// p.Ronda.sigTurno()

	return nil
}

// PRE: supongo que el jugador que toca este envido
// no tiene flor (es checkeada cuando es su turno)
type tocarEnvido struct {
	Jugada
}

func (jugada tocarEnvido) hacer(p *Partida) error {
	esPrimeraMano := p.Ronda.manoEnJuego == primera
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.muestra)
	esDelEquipoContrario := p.Ronda.envido.estado == NOCANTADOAUN || p.Ronda.envido.cantadoPor.equipo != jugada.autor.jugador.equipo
	envidoHabilitado := (p.Ronda.envido.estado == NOCANTADOAUN || p.Ronda.envido.estado == ENVIDO) && p.Ronda.flor == NOCANTADA
	ok := envidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario

	if !ok {
		return fmt.Errorf(`No es posible cantar 'Envido'`)
	}

	fmt.Printf(">> %s toca envido\n", jugada.autor.jugador.nombre)

	// ahora checkeo si alguien tiene flor
	hayFlor, manojosConFlor := p.Ronda.getFlores()
	if hayFlor {
		// todo: deberia ir al estado magico en el que espera
		// solo por jugadas de tipo flor-related
		// lo mismo para el real-envido; falta-envido
		p.Ronda.envido.estado = DESHABILITADO
		siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
		siguienteJugada.hacer(p)

	} else {
		// 2 opciones: o bien no se jugo aun
		// o bien ya estabamos en envido
		yaSeHabiaCantadoElEnvido := p.Ronda.envido.estado == ENVIDO
		if yaSeHabiaCantadoElEnvido {
			// se aumenta el puntaje del envido en +2
			p.Ronda.envido.puntaje += 2
			p.Ronda.envido.cantadoPor = jugada.autor.jugador

		} else { // no se habia jugado aun
			p.Ronda.envido.cantadoPor = jugada.autor.jugador
			p.Ronda.envido.estado = ENVIDO
			p.Ronda.envido.puntaje = 2
		}
	}

	return nil
}

// donde 'j' el jugador que dijo 'quiero' al 'envido'/'real envido'
func (jugada tocarEnvido) eval(p *Partida) error {
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

type tocarRealEnvido struct {
	Jugada
}

func (jugada tocarRealEnvido) hacer(p *Partida) error {
	esPrimeraMano := p.Ronda.manoEnJuego == primera
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.muestra)
	realEnvidoHabilitado := (p.Ronda.envido.estado == NOCANTADOAUN || p.Ronda.envido.estado == ENVIDO) && p.Ronda.flor == NOCANTADA
	esDelEquipoContrario := p.Ronda.envido.estado == NOCANTADOAUN || p.Ronda.envido.cantadoPor.equipo != jugada.autor.jugador.equipo
	ok := realEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario

	if !ok {
		return fmt.Errorf(`No es posible cantar 'Real Envido'`)
	}

	fmt.Printf(">> %s toca real envido\n", jugada.autor.jugador.nombre)
	p.Ronda.envido.estado = REALENVIDO
	p.Ronda.envido.cantadoPor = jugada.autor.jugador

	// ahora checkeo si alguien tiene flor
	hayFlor, manojosConFlor := p.Ronda.getFlores()

	if hayFlor {
		p.Ronda.envido.estado = DESHABILITADO
		siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
		siguienteJugada.hacer(p)

	} else {
		// 2 opciones:
		// o bien el envido no se jugo aun,
		// o bien ya estabamos en envido
		if p.Ronda.envido.estado == NOCANTADOAUN { // no se habia jugado aun
			p.Ronda.envido.puntaje = 3
		} else { // ya se habia cantado ENVIDO x cantidad de veces
			p.Ronda.envido.puntaje += 3
		}
	}

	return nil
}

type tocarFaltaEnvido struct {
	Jugada
}

func (jugada tocarFaltaEnvido) hacer(p *Partida) error {

	esPrimeraMano := p.Ronda.manoEnJuego == primera
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.muestra)
	faltaEnvidoHabilitado := p.Ronda.envido.estado >= NOCANTADOAUN && p.Ronda.envido.estado < FALTAENVIDO
	esDelEquipoContrario := p.Ronda.envido.estado == NOCANTADOAUN || p.Ronda.envido.cantadoPor.equipo != jugada.autor.jugador.equipo
	ok := faltaEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario

	if !ok {
		return fmt.Errorf(`No es posible cantar 'Falta Envido'`)
	}

	fmt.Printf(">> %s toca falta envido\n", jugada.autor.jugador.nombre)
	p.Ronda.envido.estado = FALTAENVIDO
	p.Ronda.envido.cantadoPor = jugada.autor.jugador

	// ahora checkeo si alguien tiene flor
	hayFlor, manojosConFlor := p.Ronda.getFlores()
	if hayFlor {
		p.Ronda.envido.estado = DESHABILITADO
		siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
		siguienteJugada.hacer(p)
	}

	return nil
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

func (jugada tocarFaltaEnvido) eval(p *Partida) error {
	p.Ronda.envido.estado = DESHABILITADO

	// computar envidos
	jIdx, max, out := p.Ronda.getElEnvido()

	print(out)

	// jug es el que gano el (falta) envido
	jug := &p.jugadores[jIdx]

	pts := p.calcPtsFaltaEnvido(jug.equipo)

	p.Ronda.envido.puntaje += pts
	p.puntajes[jug.equipo] += p.Ronda.envido.puntaje
	fmt.Printf(`>> La falta envido la gano %s con %v, +%v puntos
	para el equipo %s`+"\n",
		jug.nombre, max, p.Ronda.envido.puntaje, jug.equipo)

	return nil
}

type cantarFlor struct {
	Jugada
}

/*
todo:
actualmente no permite que si todos cantan flor
se pase a calcular el resultado solo de las flores acumuladas
se necesita timer

*/
func (jugada cantarFlor) hacer(p *Partida) error {
	// manojo dice que puede cantar flor;
	// es esto verdad?
	florHabilitada := (p.Ronda.flor == NOCANTADA || p.Ronda.flor == FLOR) && p.Ronda.manoEnJuego == primera
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.muestra)
	ok := florHabilitada && tieneFlor
	if !ok {
		return fmt.Errorf(`No es posible cantar flor`)
	}

	// se usa por si dicen "no quiero" -> se obtiene el equipo
	// al que pertenece el que la canto en un principio para
	// poder sumarle los puntos correspondientes
	fmt.Printf(">> %s canta flor\n", jugada.autor.jugador.nombre)
	p.Ronda.envido.estado = DESHABILITADO
	p.Ronda.envido.puntaje = 3
	p.Ronda.envido.cantadoPor = jugada.autor.jugador
	p.Ronda.flor = FLOR

	// ahora checkeo si alguien tiene flor
	// retorna TODOS los jugadores que tengan flor (si es que existen)
	// aPartirDe, _ := obtenerIdx(j, p.jugadores)
	hayFlor, jugadoresConFlor := p.Ronda.getFlores()
	// creo una copia
	jugadoresConFlorCACHE := make([]*Manojo, len(jugadoresConFlor))
	copy(jugadoresConFlorCACHE, jugadoresConFlor)

	if !hayFlor {
		// Nadie mas tiene flor; entonces manojo se lleva todos
		// los puntos en juego (+3)
		p.puntajes[jugada.autor.jugador.equipo] += p.Ronda.envido.puntaje // +3
		fmt.Printf(`>> +%v puntos para el equipo %s`+"\n",
			3, jugada.autor.jugador.equipo)
		p.Ronda.envido.estado = DESHABILITADO
		p.Ronda.flor = DESHABILITADA
		return nil
	}

	// si hayFlor:
	// entonces tengo que esperar respuesta SOLO de alguno de ellos;
	// a menos de un "Me voy al mazo; esa tambien es aceptada"
	// las otras las descarto
	// si no recibo respuesta en menos de x tiempo la canto yo
	// por ellos

	// Se cachea turno actual (del que canto flor).
	// Cuando se termine de jugar la flor,
	// se reestablece a este.

	todosLosJugadoresConFlorCantaron := false
	for !todosLosJugadoresConFlorCantaron {

		sigJugada := p.getSigJugada()
		esAlguienDelQueEspero := contains(jugadoresConFlor, sigJugada.getAutor())

		_, esMeVoyAlMazo := sigJugada.(irseAlMazo)
		_, esCantoFlor := sigJugada.(cantarFlor)
		_, esCantoContraFlor := sigJugada.(cantarContraFlor)
		_, esCantoContraFlorAlResto := sigJugada.(cantarContraFlorAlResto)
		_, esCantoConFlorMeAchico := sigJugada.(cantarConFlorMeAchico)
		esTipoFlor := esCantoFlor || esCantoContraFlor || esCantoContraFlorAlResto || esCantoConFlorMeAchico
		_, esQuiero := sigJugada.(responderQuiero)
		_, esNoQuiero := sigJugada.(responderNoQuiero)
		esRespuesta := esQuiero || esNoQuiero
		seEstaJugandoLaFlor := p.Ronda.flor == FLOR
		seEstaJugandoLaContraFlor := p.Ronda.flor == CONTRAFLOR
		seEstaJugandoLaContraFlorAlResto := p.Ronda.flor == CONTRAFLORALRESTO

		esDeAlguienQueNoEsperoYNoEsIrseAlMazo := !esAlguienDelQueEspero && !esMeVoyAlMazo
		esDeAlguienQueEsperoPeroNoEsNiFlorNiIrseAlMazo := esAlguienDelQueEspero && ((seEstaJugandoLaFlor && !(esTipoFlor || esMeVoyAlMazo)) || (seEstaJugandoLaContraFlor && !(esTipoFlor || esMeVoyAlMazo || esRespuesta)) || (seEstaJugandoLaContraFlorAlResto && !esRespuesta))

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
			jugadoresConFlor = eliminar(jugadoresConFlor, sigJugada.getAutor())
			// era el ultimo que del que me faltaba escuchar?
			// y por ende -> fin del bucle ?
		}

		// la ejecuto porque por descarte ya se que es valida
		if esMeVoyAlMazo {
			sigJugada.hacer(p)

		} else if esCantoFlor {
			// ya se que estaba habilitado para cantar flor
			// porque estaba en `jugadoresConFlor`
			// ahora: se canto contraflor o mayor -> inhabilitado
			florHabilitada := p.Ronda.flor == FLOR
			if !florHabilitada {
				// entonces lo vuelvo a agregar a la lista de esperados; se equivoco
				jugadoresConFlor = append(jugadoresConFlor, sigJugada.getAutor())
				return fmt.Errorf(`Ya no es posible cantar flor;`)
			}
			// en caso contrario; esta todo bien;
			// la canta
			fmt.Printf(">> %s canta flor\n", sigJugada.getAutor().jugador.nombre)
			p.Ronda.envido.cantadoPor = sigJugada.getAutor().jugador
			// ahora la flor pasa a jugarse por +3 puntos
			p.Ronda.envido.puntaje += 3

		} else if esCantoContraFlor {
			// ya se que estaba habilitado para cantar flor
			// porque estaba en `jugadoresConFlor`
			// ahora: se canto contraflor o algo asi? si si -> inhabilitado
			contraFlorHabilitada := p.Ronda.flor == FLOR
			if !contraFlorHabilitada {
				// entonces lo vuelvo a agregar a la lista de esperados; se equivoco
				jugadoresConFlor = append(jugadoresConFlor, sigJugada.getAutor())
				return fmt.Errorf(`Ya no es posible cantar contra flor;`)
			}
			// en caso contrario; esta todo bien;
			// la canta
			fmt.Printf(">> %s canta contra-flor\n", sigJugada.getAutor().jugador.nombre)
			p.Ronda.flor = CONTRAFLOR
			p.Ronda.envido.cantadoPor = sigJugada.getAutor().jugador
			// ahora la flor pasa a jugarse por 4 puntos
			p.Ronda.envido.puntaje = 4
			// y ahora tengo que esperar por la respuesta de la nueva
			// propuesta de todos menos de el que canto la contraflor
			// restauro la copia
			jugadoresConFlor = make([]*Manojo, len(jugadoresConFlorCACHE))
			copy(jugadoresConFlor, jugadoresConFlorCACHE)
			// lo elimino de los que espero
			jugadoresConFlor = eliminar(jugadoresConFlor, sigJugada.getAutor())

		} else if esCantoContraFlorAlResto {
			// ya se que estaba habilitado para cantar flor
			// porque estaba en `jugadoresConFlor`
			// ahora: puede cantarContraFlorAlResto?
			contraFlorAlRestoHabilitada := p.Ronda.flor == FLOR || p.Ronda.flor == CONTRAFLOR
			if !contraFlorAlRestoHabilitada {
				// entonces lo vuelvo a agregar a la lista de esperados; se equivoco
				// por ejemplo; ya otro jugador habia cantado contraFlorAlResto
				// ya que solo espero quiero|noQuiero|alMazo del el
				jugadoresConFlor = append(jugadoresConFlor, sigJugada.getAutor())
				return fmt.Errorf(`Ya no es posible cantar contra flor al resto;`)
			}
			// en caso contrario; esta todo bien;
			// la canta
			fmt.Printf(">> %s canta contra-flor-al-resto\n", sigJugada.getAutor().jugador.nombre)
			p.Ronda.flor = CONTRAFLORALRESTO
			p.Ronda.envido.cantadoPor = sigJugada.getAutor().jugador
			// los puntos de la flor quedan acumulados
			// y ahora tengo que esperar por la respuesta de la nueva
			// propuesta de todos menos de el que canto la contraflor
			// restauro la copia
			jugadoresConFlor = make([]*Manojo, len(jugadoresConFlorCACHE))
			copy(jugadoresConFlor, jugadoresConFlorCACHE)
			// lo elimino de los que espero
			jugadoresConFlor = eliminar(jugadoresConFlor, sigJugada.getAutor())

		} else if esQuiero && seEstaJugandoLaContraFlor {
			// solo con que uno *DEL EQUIPO CONTRARIO*
			// al que canto la contra-flor diga quiero
			// es del equipo contrario?
			esDelEquipoContrario := sigJugada.getAutor().jugador.equipo != p.Ronda.envido.cantadoPor.equipo
			if !esDelEquipoContrario {
				return fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
			}
			fmt.Printf(">> %s dice quiero \n", sigJugada.getAutor().jugador.nombre)
			// ok; se cierra el envite; hora de calcular el ganador
			p.Ronda.flor = DESHABILITADA
			manojoConLaFlorMasAlta, maxFlor := p.Ronda.getLaFlorMasAlta()
			equipoGanador := manojoConLaFlorMasAlta.jugador.equipo
			// ahora se quien es el ganador; necesito saber cuantos puntos
			// se le va a sumar a ese equipo:
			// los acumulados del envite hasta ahora
			puntosASumar := p.Ronda.envido.puntaje
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
			esDelEquipoContrario := sigJugada.getAutor().jugador.equipo != p.Ronda.envido.cantadoPor.equipo
			if !esDelEquipoContrario {
				return fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
			}
			// ok; se cierra el envite; los puntos van para el que propuso el envite
			p.Ronda.flor = DESHABILITADA
			equipoGanador := p.Ronda.envido.cantadoPor.equipo
			// ahora se quien es el ganador; necesito saber cuantos puntos
			// se le va a sumar a ese equipo:
			// los acumulados del envite hasta ahora + la contrafloralresto
			puntosASumar := p.Ronda.envido.puntaje
			p.puntajes[equipoGanador] += puntosASumar
			fmt.Printf(`>> La contra-flor la gano %s, +%v puntos
					para el equipo %s`+"\n",
				p.Ronda.envido.cantadoPor.nombre, puntosASumar, equipoGanador)
			// se corta el bucle de la flor:
			break

		} else if esQuiero && seEstaJugandoLaContraFlorAlResto {
			// solo con que uno *DEL EQUIPO CONTRARIO*
			// al que canto la contra-flor-al-resto diga quiero
			// es del equipo contrario?
			esDelEquipoContrario := sigJugada.getAutor().jugador.equipo != p.Ronda.envido.cantadoPor.equipo
			if !esDelEquipoContrario {
				return fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
			}
			fmt.Printf(">> %s dice quiero \n", sigJugada.getAutor().jugador.nombre)
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
			esDelEquipoContrario := sigJugada.getAutor().jugador.equipo != p.Ronda.envido.cantadoPor.equipo
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

	return nil

}

type cantarContraFlor struct {
	Jugada
}

func (jugada cantarContraFlor) hacer(p *Partida) error {
	// si llego aca es porque canto contra-flor ANTES que flor;
	// lo cual no dberia pasar;
	// ya que lo deberia de tomar el listener de la flor
	return fmt.Errorf(`No es posible cantar contra-flor ahora`)
}

type cantarContraFlorAlResto struct {
	Jugada
}

func (jugada cantarContraFlorAlResto) hacer(p *Partida) error {
	// si llego aca es porque canto contra-flor ANTES que flor;
	// lo cual no dberia pasar;
	// ya que lo deberia de tomar el listener de la flor
	return fmt.Errorf(`No es posible cantar contra-flor ahora`)
}

type cantarConFlorMeAchico struct {
	Jugada
}

func (jugada cantarConFlorMeAchico) hacer(p *Partida) error {
	return nil
}

type gritarTruco struct {
	Jugada
}

func (jugada gritarTruco) hacer(p *Partida) error {
	// checkeos:
	noSeFueAlMazo := jugada.autor.seFueAlMazo == false
	trucoNoSeJugoAun := p.Ronda.truco.estado == NOCANTADO
	esSuTurno := p.getJugador(p.Ronda.turno) == jugada.autor.jugador
	noSeEstaJugandoElEnvido := p.Ronda.envido.estado == NOCANTADOAUN || p.Ronda.envido.estado == DESHABILITADO
	noSeEstaJugandoLaFlor := p.Ronda.flor == NOCANTADA || p.Ronda.flor == DESHABILITADA
	trucoHabilitado := noSeFueAlMazo && trucoNoSeJugoAun && esSuTurno && noSeEstaJugandoElEnvido && noSeEstaJugandoLaFlor

	if !trucoHabilitado {
		return fmt.Errorf("No es posible cantar truco ahora")
	}

	fmt.Printf(">> %s grita truco\n", jugada.autor.jugador.nombre)
	p.Ronda.truco.cantadoPor = jugada.autor
	p.Ronda.truco.estado = TRUCO

	return nil
}

type gritarReTruco struct {
	Jugada
}

func (jugada gritarReTruco) hacer(p *Partida) error {
	// checkeos:
	noSeFueAlMazo := jugada.autor.seFueAlMazo == false
	trucoYaQuerido := p.Ronda.truco.estado == TRUCOQUERIDO
	tieneElQuiero := p.Ronda.truco.cantadoPor == jugada.autor
	esSuTurno := p.getJugador(p.Ronda.turno) == jugada.autor.jugador
	noSeEstaJugandoElEnvido := p.Ronda.envido.estado == NOCANTADOAUN || p.Ronda.envido.estado == DESHABILITADO
	noSeEstaJugandoLaFlor := p.Ronda.flor == NOCANTADA || p.Ronda.flor == DESHABILITADA
	esDelEquipoContrario := p.Ronda.truco.cantadoPor.jugador.equipo != jugada.autor.jugador.equipo
	reTrucoHabilitado := noSeFueAlMazo && trucoYaQuerido && tieneElQuiero && esSuTurno && noSeEstaJugandoElEnvido && noSeEstaJugandoLaFlor && esDelEquipoContrario

	if !reTrucoHabilitado {
		return fmt.Errorf("No es posible cantar re-truco ahora")
	}

	fmt.Printf(">> %s grita re-truco\n", jugada.autor.jugador.nombre)
	p.Ronda.truco.cantadoPor = jugada.autor
	p.Ronda.truco.estado = RETRUCO

	return nil
}

type gritarVale4 struct {
	Jugada
}

func (jugada gritarVale4) hacer(p *Partida) error {
	// checkeos:
	noSeFueAlMazo := jugada.autor.seFueAlMazo == false
	retrucoYaQuerido := p.Ronda.truco.estado == RETRUCOQUERIDO
	tieneElQuiero := p.Ronda.truco.cantadoPor == jugada.autor
	esSuTurno := p.getJugador(p.Ronda.turno) == jugada.autor.jugador
	noSeEstaJugandoElEnvido := p.Ronda.envido.estado == NOCANTADOAUN || p.Ronda.envido.estado == DESHABILITADO
	noSeEstaJugandoLaFlor := p.Ronda.flor == NOCANTADA || p.Ronda.flor == DESHABILITADA
	esDelEquipoContrario := p.Ronda.truco.cantadoPor.jugador.equipo != jugada.autor.jugador.equipo
	vale4Habilitado := noSeFueAlMazo && retrucoYaQuerido && tieneElQuiero && esSuTurno && noSeEstaJugandoElEnvido && noSeEstaJugandoLaFlor && esDelEquipoContrario

	if !vale4Habilitado {
		return fmt.Errorf("No es posible cantar re-truco ahora")
	}

	fmt.Printf(">> %s grita re-truco\n", jugada.autor.jugador.nombre)
	p.Ronda.truco.cantadoPor = jugada.autor
	p.Ronda.truco.estado = VALE4

	return nil
}

type responderQuiero struct {
	Jugada
}

func (jugada responderQuiero) hacer(p *Partida) error {
	// se acepta una respuesta 'quiero' solo cuando:
	// - CASO I: se toco un envite+ (con autor del equipo contario)
	// - CASO II: se grito el truco+ (con autor del equipo contario)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := p.Ronda.envido.estado >= ENVIDO && p.Ronda.envido.cantadoPor != jugada.autor.jugador
	// ojo: solo a la contraflor+ se le puede decir quiero; a la flor sola no
	laContraFlorEsRespondible := p.Ronda.flor >= CONTRAFLOR && p.Ronda.envido.cantadoPor != jugada.autor.jugador
	elTrucoEsRespondible := contains([3]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.truco.estado) && p.Ronda.truco.cantadoPor.jugador.equipo != jugada.autor.jugador.equipo

	ok := elEnvidoEsRespondible || laContraFlorEsRespondible || elTrucoEsRespondible
	if !ok {
		// si no, esta respondiendo al pedo
		return fmt.Errorf(`No hay nada \"que querer\"; ya que: el 
		estado del envido no es "envido" (o mayor) y el estado del 
		truco no es "truco" (o mayor)`)
	}

	if elEnvidoEsRespondible {
		fmt.Printf(">> %s responde quiero\n", jugada.autor.jugador.nombre)
		if p.Ronda.envido.estado == FALTAENVIDO {
			return tocarFaltaEnvido{Jugada{autor: jugada.autor}}.eval(p)
		}
		// si no, era envido/real-envido o cualquier
		// combinacion valida de ellos
		return tocarEnvido{Jugada{autor: jugada.autor}}.eval(p)

	} else if laContraFlorEsRespondible {
		// tengo que verificar si efectivamente tiene flor
		tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.muestra)

		if !tieneFlor {
			return fmt.Errorf(`No tiene flor; la jugada es incompatible`)
		}

		// todo ok: tiene flor; se pasa a jugar:
		// empieza cantando el autor del envite no el que "quizo"
		aPartirDe, _ := obtenerIdx(p.Ronda.envido.cantadoPor, p.jugadores)
		manojoConLaFlorGanadora, _, _ := p.Ronda.cantarFlores(aPartirDe)
		if p.Ronda.flor == CONTRAFLOR {
			// sumo +3 por cada flor (todas las flores de la ronda, que no se haya ido)
			// al equipo del ganador
			_, flores := p.Ronda.getFlores()
			totalPts := 0
			for _, m := range flores {
				if !m.seFueAlMazo {
					totalPts += 3
				}
			}
			fmt.Printf(">> %s gano la contra flor. +%v puntos para el equipo %s\n",
				manojoConLaFlorGanadora.jugador.nombre, totalPts, manojoConLaFlorGanadora.jugador.equipo)

		} else {
			// el equipo del ganador de la contraflor al resto
			// gano la partida
			equipoDelGanador := manojoConLaFlorGanadora.jugador.equipo
			ptsFaltantes := p.puntuacion.toInt() - p.puntajes[equipoDelGanador]

			fmt.Printf(">> %s gano la contra flor al resto: +%v puntos para el equipo %s\n",
				manojoConLaFlorGanadora.jugador.nombre, ptsFaltantes, manojoConLaFlorGanadora.jugador.equipo)
			fmt.Printf(">> el equipo %s gano la partida\n", manojoConLaFlorGanadora.jugador.equipo)

		}

	} else if elTrucoEsRespondible {
		p.Ronda.truco.cantadoPor = jugada.autor
		switch p.Ronda.truco.estado {
		case TRUCO:
			p.Ronda.truco.estado = TRUCOQUERIDO
		case RETRUCO:
			p.Ronda.truco.estado = RETRUCOQUERIDO
		case VALE4:
			p.Ronda.truco.estado = VALE4QUERIDO
		}
	}

	return nil

}

type responderNoQuiero struct {
	Jugada
}

func (jugada responderNoQuiero) hacer(p *Partida) error {
	// se acepta una respuesta 'no quiero' solo cuando:
	// - CASO I: se toco el envido (o similar)
	// - CASO II: se grito el truco (o similar)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := p.Ronda.envido.estado >= ENVIDO && p.Ronda.envido.cantadoPor.equipo != jugada.autor.jugador.equipo
	laFlorEsRespondible := p.Ronda.flor >= FLOR && p.Ronda.envido.cantadoPor.equipo != jugada.autor.jugador.equipo
	elTrucoEsRespondible := contains([3]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.truco.estado) && p.Ronda.truco.cantadoPor.jugador.equipo != jugada.autor.jugador.equipo

	ok := elEnvidoEsRespondible || laFlorEsRespondible || elTrucoEsRespondible

	if !ok {
		// si no, esta respondiendo al pedo
		return fmt.Errorf(`%s esta respondiendo al pedo; no hay 
		nada respondible`, jugada.autor.jugador.nombre)
	}

	if elEnvidoEsRespondible {
		fmt.Printf(">> %s responde no quiero\n", jugada.autor.jugador.nombre)

		//	no se toma en cuenta el puntaje total del ultimo toque

		var totalPts int

		switch p.Ronda.envido.estado {
		case ENVIDO:
			totalPts = p.Ronda.envido.puntaje - 1
		case REALENVIDO:
			totalPts = p.Ronda.envido.puntaje - 2
		case FALTAENVIDO:
			totalPts = p.Ronda.envido.puntaje + 1
		}

		p.Ronda.envido.estado = DESHABILITADO
		p.Ronda.envido.puntaje = totalPts
		p.puntajes[p.Ronda.envido.cantadoPor.equipo] += totalPts
		fmt.Printf(`>> +%v puntos para el equipo %s`+"\n",
			totalPts, p.Ronda.envido.cantadoPor.equipo)

	} else if laFlorEsRespondible {

		// cuenta como un "no quiero" (codigo copiado)
		// segun el estado de la apuesta actual:
		// los "me achico" no cuentan para la flor
		// Flor		xcg(+3) / xcg(+3)
		// Flor + Contra-Flor		xc(+3) / xCadaFlorDelQueHizoElDesafio(+3) + 1
		// Flor + [Contra-Flor] + ContraFlorAlResto		~Falta Envido + *TODAS* las flores no achicadas / xcg(+3) + 1

		// sumo todas las flores del equipo contrario
		totalPts := 0

		for _, m := range p.Ronda.manojos {
			esDelEquipoContrario := p.Ronda.envido.cantadoPor.equipo != jugada.autor.jugador.equipo
			tieneFlor, _ := m.tieneFlor(p.Ronda.muestra)
			if tieneFlor && esDelEquipoContrario {
				totalPts += 3
			}
		}

		if p.Ronda.flor == CONTRAFLOR || p.Ronda.flor == CONTRAFLORALRESTO {
			// si es contraflor o al resto
			// se suma 1 por el `no quiero`
			totalPts++
		}

		fmt.Printf(`>> +%v puntos para el equipo %s por las flores`+"\n",
			totalPts, p.Ronda.envido.cantadoPor.equipo)

	} else if elTrucoEsRespondible {

		// si dice no quiero:
		// todos los puntos en juego van para el equipo que hizo la apuesta de truco
		// y se termina la ronda
		var totalPts int
		switch p.Ronda.truco.estado {
		case TRUCO:
			totalPts = 1
		case RETRUCO:
			totalPts = 2
		case VALE4:
			totalPts = 3
		}
		fmt.Printf(`>> +%v puntos para el equipo %s por el %s no querido por %s`+"\n",
			totalPts,
			p.Ronda.truco.cantadoPor.jugador.equipo,
			p.Ronda.truco.estado.toString(),
			jugada.autor.jugador.nombre)
		termino := p.sumarPuntos(p.Ronda.truco.cantadoPor.jugador.equipo, totalPts)
		if !termino {
			sigMano := p.Ronda.getSigMano()
			p.nuevaRonda(sigMano)
		}
	}

	return nil
}

type irseAlMazo struct {
	Jugada
}

func (jugada irseAlMazo) hacer(p *Partida) error {
	// checkeos:
	yaSeFueAlMazo := jugada.autor.seFueAlMazo == false
	seEstabaJugandoElEnvido := p.Ronda.envido.estado >= ENVIDO
	seEstabaJugandoLaFlor := p.Ronda.flor >= FLOR
	seEstabaJugandoElTruco := p.Ronda.truco.estado >= TRUCO

	if yaSeFueAlMazo {
		return fmt.Errorf("No es posible irse al mazo ahora")
	}

	// no se puede ir al mazo sii:
	// 1. el fue el que canto el envido (y el envido esta en juego)
	// 2. tampoco se puede ir al mazo si el canto la flor o similar
	// 3. tampoco se puede ir al mazo si el grito el truco

	noSePuedeIrPorElEnvite := (seEstabaJugandoElEnvido || seEstabaJugandoLaFlor) && p.Ronda.envido.cantadoPor == jugada.autor.jugador
	// la de la flor es igual al del envido; porque es un envite
	noSePuedeIrPorElTruco := seEstabaJugandoElTruco && p.Ronda.truco.cantadoPor == jugada.autor
	if noSePuedeIrPorElEnvite || noSePuedeIrPorElTruco {
		return fmt.Errorf("No es posible irse al mazo ahora")
	}

	// ok -> se va al mazo:
	fmt.Printf(">> %s se va al mazo\n", jugada.autor.jugador.nombre)
	jugada.autor.seFueAlMazo = true
	equipoDelJugador := jugada.autor.jugador.equipo
	p.Ronda.cantJugadoresEnJuego[equipoDelJugador]--
	seFueronTodos := p.Ronda.cantJugadoresEnJuego[equipoDelJugador] == 0

	if seFueronTodos {
		// el equipo contrario gana la ronda
		// y todo lo que estaba en juego hasta ahora
		// envido; flor; truco;
		// si no habia nada en juego -> suma 1 punto
		if seEstabaJugandoElEnvido {
			// cuenta como un "no quiero"

			// codigo copiado de "no quiero"
			//	no se toma en cuenta el puntaje total del ultimo toque
			var totalPts int
			e := &p.Ronda.envido
			switch e.estado {
			case ENVIDO:
				totalPts = e.puntaje - 1
			case REALENVIDO:
				totalPts = e.puntaje - 2
			case FALTAENVIDO:
				totalPts = e.puntaje + 1
			}
			e.estado = DESHABILITADO
			e.puntaje = totalPts
			p.puntajes[e.cantadoPor.equipo] += totalPts
			fmt.Printf(`>> +%v puntos del envite para el equipo %s`+"\n",
				totalPts, e.cantadoPor.equipo)

		}

		if seEstabaJugandoLaFlor {
			// cuenta como un "no quiero"
			// segun el estado de la apuesta actual:
			// los "me achico" no cuentan para la flor
			// Flor		xcg(+3) / xcg(+3)
			// Flor + Contra-Flor		xc(+3) / xCadaFlorDelQueHizoElDesafio(+3) + 1
			// Flor + [Contra-Flor] + ContraFlorAlResto		~Falta Envido + *TODAS* las flores no achicadas / xcg(+3) + 1

			// sumo todas las flores del equipo contrario
			totalPts := 0

			for _, m := range p.Ronda.manojos {
				esDelEquipoContrario := p.Ronda.envido.cantadoPor.equipo != jugada.autor.jugador.equipo
				tieneFlor, _ := m.tieneFlor(p.Ronda.muestra)
				if tieneFlor && esDelEquipoContrario {
					totalPts += 3
				}
			}

			if p.Ronda.flor == CONTRAFLOR || p.Ronda.flor == CONTRAFLORALRESTO {
				// si es contraflor o al resto
				// se suma 1 por el `no quiero`
				totalPts++
			}

			fmt.Printf(`>> +%v puntos para el equipo %s por las flores`+"\n",
				totalPts, p.Ronda.envido.cantadoPor.equipo)

		}

		if seEstabaJugandoElTruco {
			// parecido a un "no quiero"
			// todos los puntos en juego van para el equipo que hizo la apuesta de truco
			// y se termina la ronda
			var totalPts int
			switch p.Ronda.truco.estado {
			case TRUCO:
				totalPts = 1
			case TRUCOQUERIDO:
				totalPts = 2
			case RETRUCO:
				totalPts = 2
			case RETRUCOQUERIDO:
				totalPts = 3
			case VALE4:
				totalPts = 3
			case VALE4QUERIDO:
				totalPts = 4
			}
			fmt.Printf(`>> +%v puntos para el equipo %s por el %s no querido por %s`+"\n",
				totalPts,
				p.Ronda.truco.cantadoPor.jugador.equipo,
				p.Ronda.truco.estado.toString(),
				jugada.autor.jugador.nombre)
			termino := p.sumarPuntos(p.Ronda.truco.cantadoPor.jugador.equipo, totalPts)
			if !termino {
				sigMano := p.Ronda.getSigMano()
				p.nuevaRonda(sigMano)
			}
		}

		noHabiaNadaEnJuego := !(seEstabaJugandoElEnvido || seEstabaJugandoLaFlor || seEstabaJugandoElTruco)
		if noHabiaNadaEnJuego {
			equipoContrario := jugada.autor.jugador.getEquipoContrario()
			p.puntajes[equipoContrario]++
		}

		// como se fueron todos:
		sigMano := p.Ronda.getSigMano()
		p.nuevaRonda(sigMano)
	}
	return nil
}

var jugadas = map[string]([]string){
	"Gritos": []string{
		"Truco",    // 1/2
		"Re-Truco", // 2/3
		"Vale 4",   // 3/4
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
