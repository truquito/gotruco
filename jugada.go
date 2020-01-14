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
	fmt.Printf("<< %s tira la carta %s\n",
		jugada.autor.Jugador.Nombre,
		jugada.Carta.toString())
	jugada.autor.CartasNoTiradas[idx] = false
	jugada.autor.UltimaTirada = idx
	p.Ronda.getManoActual().agregarTirada(jugada)

	// era el ultimo en tirar de esta mano?
	eraElUltimoEnTirar := p.Ronda.sigHabilitado(*jugada.autor) == nil
	if eraElUltimoEnTirar {
		// de ser asi tengo que checkear el resultado de la mano
		p.evaluarMano()
		// el turno del siguiente queda dado por el ganador de esta
	} else {
		p.Ronda.nextTurno()
	}

	return nil
}

// PRE: supongo que el jugador que toca este envido
// no tiene flor (es checkeada cuando es su turno)
type tocarEnvido struct {
	Jugada
}

func (jugada tocarEnvido) hacer(p *Partida) error {
	esPrimeraMano := p.Ronda.ManoEnJuego == primera
	esSuTurno := p.getJugador(p.Ronda.Turno) == jugada.autor.Jugador
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	esDelEquipoContrario := p.Ronda.Envido.Estado == NOCANTADOAUN || p.Ronda.Envido.CantadoPor.Equipo != jugada.autor.Jugador.Equipo
	envidoHabilitado := (p.Ronda.Envido.Estado == NOCANTADOAUN || p.Ronda.Envido.Estado == ENVIDO) && p.Ronda.Flor == NOCANTADA
	ok := envidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario && esSuTurno

	if !ok {
		return fmt.Errorf(`No es posible cantar 'Envido'`)
	}

	fmt.Printf("<< %s toca envido\n", jugada.autor.Jugador.Nombre)

	// ahora checkeo si alguien tiene flor
	hayFlor, manojosConFlor := p.Ronda.getFlores()
	if hayFlor {
		// todo: deberia ir al estado magico en el que espera
		// solo por jugadas de tipo flor-related
		// lo mismo para el real-envido; falta-envido
		p.Ronda.Envido.Estado = DESHABILITADO
		siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
		siguienteJugada.hacer(p)

	} else {
		// 2 opciones: o bien no se jugo aun
		// o bien ya estabamos en envido
		yaSeHabiaCantadoElEnvido := p.Ronda.Envido.Estado == ENVIDO
		if yaSeHabiaCantadoElEnvido {
			// se aumenta el puntaje del envido en +2
			p.Ronda.Envido.Puntaje += 2
			p.Ronda.Envido.CantadoPor = jugada.autor.Jugador

		} else { // no se habia jugado aun
			p.Ronda.Envido.CantadoPor = jugada.autor.Jugador
			p.Ronda.Envido.Estado = ENVIDO
			p.Ronda.Envido.Puntaje = 2
		}
	}

	return nil
}

// donde 'j' el jugador que dijo 'quiero' al 'envido'/'real envido'
func (jugada tocarEnvido) eval(p *Partida) error {
	p.Ronda.Envido.Estado = DESHABILITADO
	jIdx, max, out := p.Ronda.getElEnvido()
	print(out)

	jug := &p.jugadores[jIdx]
	p.Puntajes[jug.Equipo] += p.Ronda.Envido.Puntaje
	fmt.Printf(`<< El envido lo gano %s con %v, +%v puntos para el equipo %s`+"\n",
		jug.Nombre, max, p.Ronda.Envido.Puntaje, jug.Equipo)

	return nil
}

type tocarRealEnvido struct {
	Jugada
}

func (jugada tocarRealEnvido) hacer(p *Partida) error {
	esPrimeraMano := p.Ronda.ManoEnJuego == primera
	esSuTurno := p.getJugador(p.Ronda.Turno) == jugada.autor.Jugador
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	realEnvidoHabilitado := (p.Ronda.Envido.Estado == NOCANTADOAUN || p.Ronda.Envido.Estado == ENVIDO) && p.Ronda.Flor == NOCANTADA
	esDelEquipoContrario := p.Ronda.Envido.Estado == NOCANTADOAUN || p.Ronda.Envido.CantadoPor.Equipo != jugada.autor.Jugador.Equipo
	ok := realEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario && esSuTurno

	if !ok {
		return fmt.Errorf(`No es posible cantar 'Real Envido'`)
	}

	fmt.Printf("<< %s toca real envido\n", jugada.autor.Jugador.Nombre)
	p.Ronda.Envido.Estado = REALENVIDO
	p.Ronda.Envido.CantadoPor = jugada.autor.Jugador

	// ahora checkeo si alguien tiene flor
	hayFlor, manojosConFlor := p.Ronda.getFlores()

	if hayFlor {
		p.Ronda.Envido.Estado = DESHABILITADO
		siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
		siguienteJugada.hacer(p)

	} else {
		// 2 opciones:
		// o bien el envido no se jugo aun,
		// o bien ya estabamos en envido
		if p.Ronda.Envido.Estado == NOCANTADOAUN { // no se habia jugado aun
			p.Ronda.Envido.Puntaje = 3
		} else { // ya se habia cantado ENVIDO x cantidad de veces
			p.Ronda.Envido.Puntaje += 3
		}
	}

	return nil
}

type tocarFaltaEnvido struct {
	Jugada
}

func (jugada tocarFaltaEnvido) hacer(p *Partida) error {
	esSuTurno := p.getJugador(p.Ronda.Turno) == jugada.autor.Jugador
	esPrimeraMano := p.Ronda.ManoEnJuego == primera
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	faltaEnvidoHabilitado := p.Ronda.Envido.Estado >= NOCANTADOAUN && p.Ronda.Envido.Estado < FALTAENVIDO
	esDelEquipoContrario := p.Ronda.Envido.Estado == NOCANTADOAUN || p.Ronda.Envido.CantadoPor.Equipo != jugada.autor.Jugador.Equipo
	ok := faltaEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario && esSuTurno

	if !ok {
		return fmt.Errorf(`No es posible cantar 'Falta Envido'`)
	}

	fmt.Printf("<< %s toca falta envido\n", jugada.autor.Jugador.Nombre)
	p.Ronda.Envido.Estado = FALTAENVIDO
	p.Ronda.Envido.CantadoPor = jugada.autor.Jugador

	// ahora checkeo si alguien tiene flor
	hayFlor, manojosConFlor := p.Ronda.getFlores()
	if hayFlor {
		p.Ronda.Envido.Estado = DESHABILITADO
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
	p.Ronda.Envido.Estado = DESHABILITADO

	// computar envidos
	jIdx, max, out := p.Ronda.getElEnvido()

	print(out)

	// jug es el que gano el (falta) envido
	jug := &p.jugadores[jIdx]

	pts := p.calcPtsFaltaEnvido(jug.Equipo)

	p.Ronda.Envido.Puntaje += pts
	p.Puntajes[jug.Equipo] += p.Ronda.Envido.Puntaje
	fmt.Printf(`<< La falta envido la gano %s con %v, +%v puntos para el equipo %s`+"\n",
		jug.Nombre, max, p.Ronda.Envido.Puntaje, jug.Equipo)

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
	florHabilitada := (p.Ronda.Flor == NOCANTADA || p.Ronda.Flor == FLOR) && p.Ronda.ManoEnJuego == primera
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	ok := florHabilitada && tieneFlor
	if !ok {
		return fmt.Errorf(`No es posible cantar flor`)
	}

	// se usa por si dicen "no quiero" -> se obtiene el equipo
	// al que pertenece el que la canto en un principio para
	// poder sumarle los puntos correspondientes
	fmt.Printf("<< %s canta flor\n", jugada.autor.Jugador.Nombre)
	p.Ronda.Envido.Estado = DESHABILITADO
	p.Ronda.Envido.Puntaje = 3
	p.Ronda.Envido.CantadoPor = jugada.autor.Jugador
	p.Ronda.Flor = FLOR

	// ahora checkeo si alguien tiene flor
	// retorna TODOS los jugadores que tengan flor (si es que existen)
	// aPartirDe, _ := obtenerIdx(j, p.jugadores)
	_, jugadoresConFlor := p.Ronda.getFlores()
	// creo una copia
	jugadoresConFlorCACHE := make([]*Manojo, len(jugadoresConFlor))
	copy(jugadoresConFlorCACHE, jugadoresConFlor)

	hayFlor := len(eliminar(jugadoresConFlor, jugada.getAutor())) > 0

	if !hayFlor {
		// Nadie mas tiene flor; entonces manojo se lleva todos
		// los puntos en juego (+3)
		p.Puntajes[jugada.autor.Jugador.Equipo] += p.Ronda.Envido.Puntaje // +3
		fmt.Printf(`<< +%v puntos para el equipo %s`+"\n",
			3, jugada.autor.Jugador.Equipo)
		p.Ronda.Envido.Estado = DESHABILITADO
		p.Ronda.Flor = DESHABILITADA
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
		seEstaJugandoLaFlor := p.Ronda.Flor == FLOR
		seEstaJugandoLaContraFlor := p.Ronda.Flor == CONTRAFLOR
		seEstaJugandoLaContraFlorAlResto := p.Ronda.Flor == CONTRAFLORALRESTO

		esDeAlguienQueNoEsperoYNoEsIrseAlMazo := !esAlguienDelQueEspero && !esMeVoyAlMazo
		esDeAlguienQueEsperoPeroNoEsNiFlorNiIrseAlMazo := esAlguienDelQueEspero && ((seEstaJugandoLaFlor && !(esTipoFlor || esMeVoyAlMazo)) || (seEstaJugandoLaContraFlor && !(esTipoFlor || esMeVoyAlMazo || esRespuesta)) || (seEstaJugandoLaContraFlorAlResto && !esRespuesta))

		noEsValida := esDeAlguienQueNoEsperoYNoEsIrseAlMazo || esDeAlguienQueEsperoPeroNoEsNiFlorNiIrseAlMazo

		// todo: que pasa si solo falta 1 por responder y se va al mazo?

		if noEsValida {
			// no deberia de salir de este loop
			// pero solo responderle al loco (?)
			return fmt.Errorf(`No es el momento de realizar esta jugada; ahora estoy esperando por cantos de flor (de aquellos que la poseen) o bien "Irse al mazo" (de cualquier jugador)`)
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
			florHabilitada := p.Ronda.Flor == FLOR
			if !florHabilitada {
				// entonces lo vuelvo a agregar a la lista de esperados; se equivoco
				jugadoresConFlor = append(jugadoresConFlor, sigJugada.getAutor())
				return fmt.Errorf(`Ya no es posible cantar flor;`)
			}
			// en caso contrario; esta todo bien;
			// la canta
			fmt.Printf("<< %s canta flor\n", sigJugada.getAutor().Jugador.Nombre)
			p.Ronda.Envido.CantadoPor = sigJugada.getAutor().Jugador
			// ahora la flor pasa a jugarse por +3 puntos
			p.Ronda.Envido.Puntaje += 3

		} else if esCantoContraFlor {
			// ya se que estaba habilitado para cantar flor
			// porque estaba en `jugadoresConFlor`
			// ahora: se canto contraflor o algo asi? si si -> inhabilitado
			contraFlorHabilitada := p.Ronda.Flor == FLOR
			if !contraFlorHabilitada {
				// entonces lo vuelvo a agregar a la lista de esperados; se equivoco
				jugadoresConFlor = append(jugadoresConFlor, sigJugada.getAutor())
				return fmt.Errorf(`Ya no es posible cantar contra flor;`)
			}
			// en caso contrario; esta todo bien;
			// la canta
			fmt.Printf("<< %s canta contra-flor\n", sigJugada.getAutor().Jugador.Nombre)
			p.Ronda.Flor = CONTRAFLOR
			p.Ronda.Envido.CantadoPor = sigJugada.getAutor().Jugador
			// ahora la flor pasa a jugarse por 4 puntos
			p.Ronda.Envido.Puntaje = 4
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
			contraFlorAlRestoHabilitada := p.Ronda.Flor == FLOR || p.Ronda.Flor == CONTRAFLOR
			if !contraFlorAlRestoHabilitada {
				// entonces lo vuelvo a agregar a la lista de esperados; se equivoco
				// por ejemplo; ya otro jugador habia cantado contraFlorAlResto
				// ya que solo espero quiero|noQuiero|alMazo del el
				jugadoresConFlor = append(jugadoresConFlor, sigJugada.getAutor())
				return fmt.Errorf(`Ya no es posible cantar contra flor al resto;`)
			}
			// en caso contrario; esta todo bien;
			// la canta
			fmt.Printf("<< %s canta contra-flor-al-resto\n", sigJugada.getAutor().Jugador.Nombre)
			p.Ronda.Flor = CONTRAFLORALRESTO
			p.Ronda.Envido.CantadoPor = sigJugada.getAutor().Jugador
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
			esDelEquipoContrario := sigJugada.getAutor().Jugador.Equipo != p.Ronda.Envido.CantadoPor.Equipo
			if !esDelEquipoContrario {
				return fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
			}
			fmt.Printf("<< %s dice quiero \n", sigJugada.getAutor().Jugador.Nombre)
			// ok; se cierra el envite; hora de calcular el ganador
			p.Ronda.Flor = DESHABILITADA
			manojoConLaFlorMasAlta, maxFlor := p.Ronda.getLaFlorMasAlta()
			equipoGanador := manojoConLaFlorMasAlta.Jugador.Equipo
			// ahora se quien es el ganador; necesito saber cuantos puntos
			// se le va a sumar a ese equipo:
			// los acumulados del envite hasta ahora
			puntosASumar := p.Ronda.Envido.Puntaje
			p.Puntajes[equipoGanador] += puntosASumar
			fmt.Printf(`<< La contra-flor-al-resto la gano %s con %v, +%v puntos para el equipo %s`+"\n",
				manojoConLaFlorMasAlta.Jugador.Nombre, maxFlor, puntosASumar, equipoGanador)
			// se corta el bucle de la flor:
			break

		} else if esNoQuiero && seEstaJugandoLaContraFlorAlResto {
			// solo con que uno *DEL EQUIPO CONTRARIO*
			// al que canto la contra-flor-al-resto diga quiero
			// es del equipo contrario?
			esDelEquipoContrario := sigJugada.getAutor().Jugador.Equipo != p.Ronda.Envido.CantadoPor.Equipo
			if !esDelEquipoContrario {
				return fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
			}
			// ok; se cierra el envite; los puntos van para el que propuso el envite
			p.Ronda.Flor = DESHABILITADA
			equipoGanador := p.Ronda.Envido.CantadoPor.Equipo
			// ahora se quien es el ganador; necesito saber cuantos puntos
			// se le va a sumar a ese equipo:
			// los acumulados del envite hasta ahora + la contrafloralresto
			puntosASumar := p.Ronda.Envido.Puntaje
			p.Puntajes[equipoGanador] += puntosASumar
			fmt.Printf(`<< La contra-flor la gano %s, +%v puntos para el equipo %s`+"\n",
				p.Ronda.Envido.CantadoPor.Nombre, puntosASumar, equipoGanador)
			// se corta el bucle de la flor:
			break

		} else if esQuiero && seEstaJugandoLaContraFlorAlResto {
			// solo con que uno *DEL EQUIPO CONTRARIO*
			// al que canto la contra-flor-al-resto diga quiero
			// es del equipo contrario?
			esDelEquipoContrario := sigJugada.getAutor().Jugador.Equipo != p.Ronda.Envido.CantadoPor.Equipo
			if !esDelEquipoContrario {
				return fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
			}
			fmt.Printf("<< %s dice quiero \n", sigJugada.getAutor().Jugador.Nombre)
			// ok; se cierra el envite; hora de calcular el ganador
			p.Ronda.Flor = DESHABILITADA
			manojoConLaFlorMasAlta, maxFlor := p.Ronda.getLaFlorMasAlta()
			equipoGanador := manojoConLaFlorMasAlta.Jugador.Equipo
			// ahora se quien es el ganador; necesito saber cuantos puntos
			// se le va a sumar a ese equipo:
			// los acumulados del envite hasta ahora + la contrafloralresto
			puntosASumar := p.Ronda.Envido.Puntaje + p.calcPtsContraFlorAlResto(equipoGanador)
			p.Puntajes[equipoGanador] += puntosASumar
			fmt.Printf(`<< La contra-flor-al-resto la gano %s con %v, +%v puntos para el equipo %s`+"\n",
				manojoConLaFlorMasAlta.Jugador.Nombre, maxFlor, puntosASumar, equipoGanador)
			// se corta el bucle de la flor:
			break

		} else if esNoQuiero && seEstaJugandoLaContraFlorAlResto {
			// solo con que uno *DEL EQUIPO CONTRARIO*
			// al que canto la contra-flor-al-resto diga quiero
			// es del equipo contrario?
			esDelEquipoContrario := sigJugada.getAutor().Jugador.Equipo != p.Ronda.Envido.CantadoPor.Equipo
			if !esDelEquipoContrario {
				return fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
			}
			// ok; se cierra el envite; los puntos van para el que propuso el envite
			p.Ronda.Flor = DESHABILITADA
			equipoGanador := p.Ronda.Envido.CantadoPor.Equipo
			// ahora se quien es el ganador; necesito saber cuantos puntos
			// se le va a sumar a ese equipo:
			// los acumulados del envite hasta ahora + la contrafloralresto
			puntosASumar := p.Ronda.Envido.Puntaje + p.calcPtsContraFlorAlResto(equipoGanador)
			p.Puntajes[equipoGanador] += puntosASumar
			fmt.Printf(`<< La contra-flor-al-resto la gano %s, +%v puntos para el equipo %s`+"\n",
				p.Ronda.Envido.CantadoPor.Nombre, puntosASumar, equipoGanador)
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
	noSeFueAlMazo := jugada.autor.SeFueAlMazo == false
	trucoNoSeJugoAun := p.Ronda.Truco.Estado == NOCANTADO
	esSuTurno := p.getJugador(p.Ronda.Turno) == jugada.autor.Jugador
	noSeEstaJugandoElEnvido := p.Ronda.Envido.Estado == NOCANTADOAUN || p.Ronda.Envido.Estado == DESHABILITADO
	noSeEstaJugandoLaFlor := p.Ronda.Flor == NOCANTADA || p.Ronda.Flor == DESHABILITADA
	trucoHabilitado := noSeFueAlMazo && trucoNoSeJugoAun && esSuTurno && noSeEstaJugandoElEnvido && noSeEstaJugandoLaFlor

	if !trucoHabilitado {
		return fmt.Errorf("No es posible cantar truco ahora")
	}

	fmt.Printf("<< %s grita truco\n", jugada.autor.Jugador.Nombre)
	p.Ronda.Truco.CantadoPor = jugada.autor
	p.Ronda.Truco.Estado = TRUCO

	return nil
}

type gritarReTruco struct {
	Jugada
}

func (jugada gritarReTruco) hacer(p *Partida) error {
	// checkeos:
	noSeFueAlMazo := jugada.autor.SeFueAlMazo == false
	trucoYaQuerido := p.Ronda.Truco.Estado == TRUCOQUERIDO
	tieneElQuiero := p.Ronda.Truco.CantadoPor == jugada.autor
	esSuTurno := p.getJugador(p.Ronda.Turno) == jugada.autor.Jugador
	noSeEstaJugandoElEnvido := p.Ronda.Envido.Estado == NOCANTADOAUN || p.Ronda.Envido.Estado == DESHABILITADO
	noSeEstaJugandoLaFlor := p.Ronda.Flor == NOCANTADA || p.Ronda.Flor == DESHABILITADA
	esDelEquipoContrario := p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	reTrucoHabilitado := noSeFueAlMazo && trucoYaQuerido && tieneElQuiero && esSuTurno && noSeEstaJugandoElEnvido && noSeEstaJugandoLaFlor && esDelEquipoContrario

	if !reTrucoHabilitado {
		return fmt.Errorf("No es posible cantar re-truco ahora")
	}

	fmt.Printf("<< %s grita re-truco\n", jugada.autor.Jugador.Nombre)
	p.Ronda.Truco.CantadoPor = jugada.autor
	p.Ronda.Truco.Estado = RETRUCO

	return nil
}

type gritarVale4 struct {
	Jugada
}

func (jugada gritarVale4) hacer(p *Partida) error {
	// checkeos:
	noSeFueAlMazo := jugada.autor.SeFueAlMazo == false
	retrucoYaQuerido := p.Ronda.Truco.Estado == RETRUCOQUERIDO
	tieneElQuiero := p.Ronda.Truco.CantadoPor == jugada.autor
	esSuTurno := p.getJugador(p.Ronda.Turno) == jugada.autor.Jugador
	noSeEstaJugandoElEnvido := p.Ronda.Envido.Estado == NOCANTADOAUN || p.Ronda.Envido.Estado == DESHABILITADO
	noSeEstaJugandoLaFlor := p.Ronda.Flor == NOCANTADA || p.Ronda.Flor == DESHABILITADA
	esDelEquipoContrario := p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	vale4Habilitado := noSeFueAlMazo && retrucoYaQuerido && tieneElQuiero && esSuTurno && noSeEstaJugandoElEnvido && noSeEstaJugandoLaFlor && esDelEquipoContrario

	if !vale4Habilitado {
		return fmt.Errorf("No es posible cantar re-truco ahora")
	}

	fmt.Printf("<< %s grita re-truco\n", jugada.autor.Jugador.Nombre)
	p.Ronda.Truco.CantadoPor = jugada.autor
	p.Ronda.Truco.Estado = VALE4

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

	elEnvidoEsRespondible := p.Ronda.Envido.Estado >= ENVIDO && p.Ronda.Envido.CantadoPor != jugada.autor.Jugador
	// ojo: solo a la contraflor+ se le puede decir quiero; a la flor sola no
	laContraFlorEsRespondible := p.Ronda.Flor >= CONTRAFLOR && p.Ronda.Envido.CantadoPor != jugada.autor.Jugador
	elTrucoEsRespondible := contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado) && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo

	ok := elEnvidoEsRespondible || laContraFlorEsRespondible || elTrucoEsRespondible
	if !ok {
		// si no, esta respondiendo al pedo
		return fmt.Errorf(`No hay nada \"que querer\"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor)`)
	}

	if elEnvidoEsRespondible {
		fmt.Printf("<< %s responde quiero\n", jugada.autor.Jugador.Nombre)
		if p.Ronda.Envido.Estado == FALTAENVIDO {
			return tocarFaltaEnvido{Jugada{autor: jugada.autor}}.eval(p)
		}
		// si no, era envido/real-envido o cualquier
		// combinacion valida de ellos
		return tocarEnvido{Jugada{autor: jugada.autor}}.eval(p)

	} else if laContraFlorEsRespondible {
		// tengo que verificar si efectivamente tiene flor
		tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)

		if !tieneFlor {
			return fmt.Errorf(`No tiene flor; la jugada es incompatible`)
		}

		// todo ok: tiene flor; se pasa a jugar:
		// empieza cantando el autor del envite no el que "quizo"
		aPartirDe, _ := obtenerIdx(p.Ronda.Envido.CantadoPor, p.jugadores)
		manojoConLaFlorGanadora, _, _ := p.Ronda.cantarFlores(aPartirDe)
		if p.Ronda.Flor == CONTRAFLOR {
			// sumo +3 por cada flor (todas las flores de la ronda, que no se haya ido)
			// al equipo del ganador
			_, flores := p.Ronda.getFlores()
			totalPts := 0
			for _, m := range flores {
				if !m.SeFueAlMazo {
					totalPts += 3
				}
			}
			fmt.Printf("<< %s gano la contra flor. +%v puntos para el equipo %s\n",
				manojoConLaFlorGanadora.Jugador.Nombre, totalPts, manojoConLaFlorGanadora.Jugador.Equipo)

		} else {
			// el equipo del ganador de la contraflor al resto
			// gano la partida
			equipoDelGanador := manojoConLaFlorGanadora.Jugador.Equipo
			ptsFaltantes := p.Puntuacion.toInt() - p.Puntajes[equipoDelGanador]

			fmt.Printf("<< %s gano la contra flor al resto: +%v puntos para el equipo %s\n",
				manojoConLaFlorGanadora.Jugador.Nombre, ptsFaltantes, manojoConLaFlorGanadora.Jugador.Equipo)
			fmt.Printf("<< el equipo %s gano la partida\n", manojoConLaFlorGanadora.Jugador.Equipo)

		}

	} else if elTrucoEsRespondible {
		p.Ronda.Truco.CantadoPor = jugada.autor
		switch p.Ronda.Truco.Estado {
		case TRUCO:
			p.Ronda.Truco.Estado = TRUCOQUERIDO
		case RETRUCO:
			p.Ronda.Truco.Estado = RETRUCOQUERIDO
		case VALE4:
			p.Ronda.Truco.Estado = VALE4QUERIDO
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

	elEnvidoEsRespondible := p.Ronda.Envido.Estado >= ENVIDO && p.Ronda.Envido.CantadoPor.Equipo != jugada.autor.Jugador.Equipo
	laFlorEsRespondible := p.Ronda.Flor >= FLOR && p.Ronda.Envido.CantadoPor.Equipo != jugada.autor.Jugador.Equipo
	elTrucoEsRespondible := contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado) && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo

	ok := elEnvidoEsRespondible || laFlorEsRespondible || elTrucoEsRespondible

	if !ok {
		// si no, esta respondiendo al pedo
		return fmt.Errorf(`%s esta respondiendo al pedo; no hay nada respondible`, jugada.autor.Jugador.Nombre)
	}

	if elEnvidoEsRespondible {
		fmt.Printf("<< %s responde no quiero\n", jugada.autor.Jugador.Nombre)

		//	no se toma en cuenta el puntaje total del ultimo toque

		var totalPts int

		switch p.Ronda.Envido.Estado {
		case ENVIDO:
			totalPts = p.Ronda.Envido.Puntaje - 1
		case REALENVIDO:
			totalPts = p.Ronda.Envido.Puntaje - 2
		case FALTAENVIDO:
			totalPts = p.Ronda.Envido.Puntaje + 1
		}

		p.Ronda.Envido.Estado = DESHABILITADO
		p.Ronda.Envido.Puntaje = totalPts
		p.Puntajes[p.Ronda.Envido.CantadoPor.Equipo] += totalPts
		fmt.Printf(`<< +%v puntos para el equipo %s`+"\n",
			totalPts, p.Ronda.Envido.CantadoPor.Equipo)

	} else if laFlorEsRespondible {

		// cuenta como un "no quiero" (codigo copiado)
		// segun el estado de la apuesta actual:
		// los "me achico" no cuentan para la flor
		// Flor		xcg(+3) / xcg(+3)
		// Flor + Contra-Flor		xc(+3) / xCadaFlorDelQueHizoElDesafio(+3) + 1
		// Flor + [Contra-Flor] + ContraFlorAlResto		~Falta Envido + *TODAS* las flores no achicadas / xcg(+3) + 1

		// sumo todas las flores del equipo contrario
		totalPts := 0

		for _, m := range p.Ronda.Manojos {
			esDelEquipoContrario := p.Ronda.Envido.CantadoPor.Equipo != jugada.autor.Jugador.Equipo
			tieneFlor, _ := m.tieneFlor(p.Ronda.Muestra)
			if tieneFlor && esDelEquipoContrario {
				totalPts += 3
			}
		}

		if p.Ronda.Flor == CONTRAFLOR || p.Ronda.Flor == CONTRAFLORALRESTO {
			// si es contraflor o al resto
			// se suma 1 por el `no quiero`
			totalPts++
		}

		fmt.Printf(`<< +%v puntos para el equipo %s por las flores`+"\n",
			totalPts, p.Ronda.Envido.CantadoPor.Equipo)

	} else if elTrucoEsRespondible {

		// si dice no quiero:
		// todos los puntos en juego van para el equipo que hizo la apuesta de truco
		// y se termina la ronda
		var totalPts int
		switch p.Ronda.Truco.Estado {
		case TRUCO:
			totalPts = 1
		case RETRUCO:
			totalPts = 2
		case VALE4:
			totalPts = 3
		}
		fmt.Printf(`<< +%v puntos para el equipo %s por el %s no querido por %s`+"\n",
			totalPts,
			p.Ronda.Truco.CantadoPor.Jugador.Equipo,
			p.Ronda.Truco.Estado.String(),
			jugada.autor.Jugador.Nombre)
		termino := p.sumarPuntos(p.Ronda.Truco.CantadoPor.Jugador.Equipo, totalPts)
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
	yaSeFueAlMazo := jugada.autor.SeFueAlMazo == true
	seEstabaJugandoElEnvido := p.Ronda.Envido.Estado >= ENVIDO
	seEstabaJugandoLaFlor := p.Ronda.Flor >= FLOR
	seEstabaJugandoElTruco := p.Ronda.Truco.Estado >= TRUCO

	if yaSeFueAlMazo {
		return fmt.Errorf("No es posible irse al mazo ahora")
	}

	// no se puede ir al mazo sii:
	// 1. el fue el que canto el envido (y el envido esta en juego)
	// 2. tampoco se puede ir al mazo si el canto la flor o similar
	// 3. tampoco se puede ir al mazo si el grito el truco

	noSePuedeIrPorElEnvite := (seEstabaJugandoElEnvido || seEstabaJugandoLaFlor) && p.Ronda.Envido.CantadoPor == jugada.autor.Jugador
	// la de la flor es igual al del envido; porque es un envite
	noSePuedeIrPorElTruco := seEstabaJugandoElTruco && p.Ronda.Truco.CantadoPor == jugada.autor
	if noSePuedeIrPorElEnvite || noSePuedeIrPorElTruco {
		return fmt.Errorf("No es posible irse al mazo ahora")
	}

	// ok -> se va al mazo:
	fmt.Printf("<< %s se va al mazo\n", jugada.autor.Jugador.Nombre)
	jugada.autor.SeFueAlMazo = true
	equipoDelJugador := jugada.autor.Jugador.Equipo
	p.Ronda.CantJugadoresEnJuego[equipoDelJugador]--
	seFueronTodos := p.Ronda.CantJugadoresEnJuego[equipoDelJugador] == 0

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
			e := &p.Ronda.Envido
			switch e.Estado {
			case ENVIDO:
				totalPts = e.Puntaje - 1
			case REALENVIDO:
				totalPts = e.Puntaje - 2
			case FALTAENVIDO:
				totalPts = e.Puntaje + 1
			}
			e.Estado = DESHABILITADO
			e.Puntaje = totalPts
			p.Puntajes[e.CantadoPor.Equipo] += totalPts
			fmt.Printf(`<< +%v puntos del envite para el equipo %s`+"\n",
				totalPts, e.CantadoPor.Equipo)

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

			for _, m := range p.Ronda.Manojos {
				esDelEquipoContrario := p.Ronda.Envido.CantadoPor.Equipo != jugada.autor.Jugador.Equipo
				tieneFlor, _ := m.tieneFlor(p.Ronda.Muestra)
				if tieneFlor && esDelEquipoContrario {
					totalPts += 3
				}
			}

			if p.Ronda.Flor == CONTRAFLOR || p.Ronda.Flor == CONTRAFLORALRESTO {
				// si es contraflor o al resto
				// se suma 1 por el `no quiero`
				totalPts++
			}

			fmt.Printf(`<< +%v puntos para el equipo %s por las flores`+"\n",
				totalPts, p.Ronda.Envido.CantadoPor.Equipo)

		}

		if seEstabaJugandoElTruco {
			// parecido a un "no quiero"
			// todos los puntos en juego van para el equipo que hizo la apuesta de truco
			// y se termina la ronda
			var totalPts int
			switch p.Ronda.Truco.Estado {
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
			fmt.Printf(`<< +%v puntos para el equipo %s por el %s no querido por %s`+"\n",
				totalPts,
				p.Ronda.Truco.CantadoPor.Jugador.Equipo,
				p.Ronda.Truco.Estado.String(),
				jugada.autor.Jugador.Nombre)
			termino := p.sumarPuntos(p.Ronda.Truco.CantadoPor.Jugador.Equipo, totalPts)
			if !termino {
				sigMano := p.Ronda.getSigMano()
				p.nuevaRonda(sigMano)
			}
		}

		noHabiaNadaEnJuego := !(seEstabaJugandoElEnvido || seEstabaJugandoLaFlor || seEstabaJugandoElTruco)
		if noHabiaNadaEnJuego {
			equipoContrario := jugada.autor.Jugador.getEquipoContrario()
			p.Puntajes[equipoContrario]++
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
		"Re-truco", // 2/3
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
		"No quiero",
	},
	"Acciones": []string{
		"Irse al mazo",
		"Tirar carta",
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
