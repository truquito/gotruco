package truco

import (
	"fmt"
)

// IJugada Interface para las jugadas
type IJugada interface {
	hacer(p *Partida)
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
func (jugada tirarCarta) hacer(p *Partida) {

	// checkeo si se fue al mazo
	noSeFueAlMazo := jugada.autor.SeFueAlMazo == false
	ok := noSeFueAlMazo
	if !ok {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible tirar una carta porque ya te fuiste al mazo"),
			},
		})

		return

	}

	// esto es un tanto redundante porque es imposible que no sea su turno
	// (checkeado mas adelante) y que al mismo tiempo tenga algo para tirar
	// luego de haber jugado sus 3 cartas; aun asi lo dejo
	yaTiroTodasSusCartas := jugada.autor.getCantCartasTiradas() == 3
	if yaTiroTodasSusCartas {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible tirar una carta porque ya las tiraste todas"),
			},
		})

		return

	}

	// checkeo flor en juego
	enviteEnJuego := p.Ronda.Envite.Estado >= ENVIDO
	if enviteEnJuego {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible tirar una carta ahora porque el envite esta en juego"),
			},
		})

		return

	}

	// primero que nada: tiene esa carta?
	idx, err := jugada.autor.getCartaIdx(jugada.Carta)
	if err != nil {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: err.Error(),
			},
		})

		return
	}

	// luego, era su turno?
	eraSuTurno := p.Ronda.getElTurno() == jugada.autor
	if !eraSuTurno {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No era su turno, no puede tirar la carta"),
			},
		})

		return

	}

	// checkeo si tiene flor
	florHabilitada := (p.Ronda.Envite.Estado >= NOCANTADOAUN && p.Ronda.Envite.Estado <= FLOR) && p.Ronda.ManoEnJuego == primera
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := contains(p.Ronda.Envite.jugadoresConFlorQueNoCantaron, jugada.autor)
	noPuedeTirar := florHabilitada && tieneFlor && noCantoFlorAun
	if noPuedeTirar {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible tirar una carta sin antes cantar la flor"),
			},
		})

		return

	}

	trucoGritado := contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado)
	unoDelEquipoContrarioGritoTruco := trucoGritado && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	elTrucoEsRespondible := trucoGritado && unoDelEquipoContrarioGritoTruco
	if elTrucoEsRespondible {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible tirar una carta porque tu equipo debe responder la propuesta del truco"),
			},
		})

		return

	}

	// ok la tiene y era su turno -> la juega
	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Tirar-Carta",
			Nota: fmt.Sprintf("%s tira la carta %s",
				jugada.autor.Jugador.Nombre, jugada.Carta.toString()),
		},
	})

	jugada.autor.CartasNoTiradas[idx] = false
	jugada.autor.UltimaTirada = idx
	p.Ronda.getManoActual().agregarTirada(jugada)

	// era el ultimo en tirar de esta mano?
	eraElUltimoEnTirar := p.Ronda.getSigHabilitado(*jugada.autor) == nil
	if eraElUltimoEnTirar {
		// de ser asi tengo que checkear el resultado de la mano
		empiezaNuevaRonda, pkts := p.evaluarMano()

		if pkts != nil {
			for _, pkt := range pkts {
				if pkt != nil {

					// antes:
					//write(p.Stdout, msg)
					// ahora:
					write(p.Stdout, pkt)

				}
			}
		}

		if !empiezaNuevaRonda {

			p.Ronda.ManoEnJuego++
			p.Ronda.setNextTurnoPosMano()

		} else {

			if !p.Terminada() {
				// ahora se deberia de incrementar el mano
				// y ser el turno de este
				sigMano := p.Ronda.getSigElMano()
				p.nuevaRonda(sigMano) // todo: el tema es que cuando llama aca
				// no manda mensaje de que arranco nueva ronda
				// falso: el padre que llama a .evaluarRonda tiene que fijarse si
				// retorno true
				// entonces debe crearla el
				// no es responsabilidad de evaluarRonda arrancar una ronda nueva!!
				// de hecho, si una ronda es terminable y se llama 2 veces consecutivas
				// al mismo metodo booleano, en ambas oportunidades retorna diferente
				// ridiculo
				write(p.Stdout, &Pkt{
					Dest: []string{"ALL"},
					Msg: Msg{
						Tipo: "Nueva-Ronda",
						Cont: nil,
					},
				})

			} else {
				p.byeBye()
			}

		}

		// el turno del siguiente queda dado por el ganador de esta
	} else {
		p.Ronda.setNextTurno()
	}

	return
}

// PRE: supongo que el jugador que toca este envido
// no tiene flor (es checkeada cuando es su turno)
type tocarEnvido struct {
	Jugada
}

func (jugada tocarEnvido) hacer(p *Partida) {
	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	if florEnJuego {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible tocar el envido ahora porque la flor esta en juego"),
			},
		})

		return

	}
	esPrimeraMano := p.Ronda.ManoEnJuego == primera
	esSuTurno := p.Ronda.getElTurno() == jugada.autor
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	esDelEquipoContrario := p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	envidoHabilitado := (p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Envite.Estado == ENVIDO)
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO
	apuestaSaturada := p.Ronda.Envite.Puntaje >= p.calcPtsFalta()
	ok := (envidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario) && (esSuTurno || yaEstabamosEnEnvido) && !apuestaSaturada

	if !ok {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf(`No es posible cantar 'Envido'`),
			},
		})

		return

	}

	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Toca-Envido",
			Cont: []byte(jugada.autor.Jugador.Nombre),
		},
	})

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.jugadoresConFlorQueNoCantaron) > 0
	if hayFlor {
		// todo: deberia ir al estado magico en el que espera
		// solo por jugadas de tipo flor-related
		// lo mismo para el real-envido; falta-envido
		manojosConFlor := p.Ronda.Envite.jugadoresConFlorQueNoCantaron
		siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
		siguienteJugada.hacer(p)

	} else {
		// 2 opciones: o bien no se jugo aun
		// o bien ya estabamos en envido
		yaSeHabiaCantadoElEnvido := p.Ronda.Envite.Estado == ENVIDO
		if yaSeHabiaCantadoElEnvido {
			// se aumenta el puntaje del envido en +2
			p.Ronda.Envite.Puntaje += 2
			p.Ronda.Envite.CantadoPor = jugada.autor

		} else { // no se habia jugado aun
			p.Ronda.Envite.CantadoPor = jugada.autor
			p.Ronda.Envite.Estado = ENVIDO
			p.Ronda.Envite.Puntaje = 2
		}
	}

	return
}

// donde 'j' el jugador que dijo 'quiero' al 'envido'/'real envido'
func (jugada tocarEnvido) eval(p *Partida) {
	p.Ronda.Envite.Estado = DESHABILITADO
	jIdx, max, out := p.Ronda.execElEnvido()
	print(out)

	jug := &p.jugadores[jIdx]

	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Sumar-Puntos",
			Nota: fmt.Sprintf(`El envido lo gano %s con %v, +%v puntos para el equipo %s`,
				jug.Nombre, max, p.Ronda.Envite.Puntaje, jug.Equipo),
			Cont: ContSumPts{
				Pts:    p.Ronda.Envite.Puntaje,
				Equipo: jug.Equipo.String(),
			}.ToBytes(),
		},
	})

	p.sumarPuntos(jug.Equipo, p.Ronda.Envite.Puntaje)

}

type tocarRealEnvido struct {
	Jugada
}

func (jugada tocarRealEnvido) hacer(p *Partida) {
	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	if florEnJuego {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible tocar real envido ahora porque la flor esta en juego"),
			},
		})

		return

	}
	esPrimeraMano := p.Ronda.ManoEnJuego == primera
	esSuTurno := p.Ronda.getElTurno() == jugada.autor
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	realEnvidoHabilitado := (p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Envite.Estado == ENVIDO)
	esDelEquipoContrario := p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO
	ok := realEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario && (esSuTurno || yaEstabamosEnEnvido)

	if !ok {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf(`No es posible cantar 'Real Envido'`),
			},
		})

		return

	}

	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Toca-RealEnvido",
			Cont: []byte(jugada.autor.Jugador.Nombre),
		},
	})

	p.Ronda.Envite.Estado = REALENVIDO
	p.Ronda.Envite.CantadoPor = jugada.autor

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.jugadoresConFlorQueNoCantaron) > 0

	if hayFlor {
		manojosConFlor := p.Ronda.Envite.jugadoresConFlorQueNoCantaron
		siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
		siguienteJugada.hacer(p)

	} else {
		// 2 opciones:
		// o bien el envido no se jugo aun,
		// o bien ya estabamos en envido
		if p.Ronda.Envite.Estado == NOCANTADOAUN { // no se habia jugado aun
			p.Ronda.Envite.Puntaje = 3
		} else { // ya se habia cantado ENVIDO x cantidad de veces
			p.Ronda.Envite.Puntaje += 3
		}
	}

	return
}

type tocarFaltaEnvido struct {
	Jugada
}

func (jugada tocarFaltaEnvido) hacer(p *Partida) {
	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	if florEnJuego {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible tocar falta envido ahora porque la flor esta en juego"),
			},
		})

		return

	}

	esSuTurno := p.Ronda.getElTurno() == jugada.autor
	esPrimeraMano := p.Ronda.ManoEnJuego == primera
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	faltaEnvidoHabilitado := p.Ronda.Envite.Estado >= NOCANTADOAUN && p.Ronda.Envite.Estado < FALTAENVIDO
	esDelEquipoContrario := p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO || p.Ronda.Envite.Estado == REALENVIDO
	ok := faltaEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario && (esSuTurno || yaEstabamosEnEnvido)

	if !ok {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf(`No es posible cantar 'Falta Envido'`),
			},
		})

		return

	}

	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Toca-FaltaEnvido",
			Cont: []byte(jugada.autor.Jugador.Nombre),
		},
	})

	p.Ronda.Envite.Estado = FALTAENVIDO
	p.Ronda.Envite.CantadoPor = jugada.autor

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.jugadoresConFlorQueNoCantaron) > 0
	if hayFlor {
		p.Ronda.Envite.Estado = DESHABILITADO
		manojosConFlor := p.Ronda.Envite.jugadoresConFlorQueNoCantaron
		siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
		siguienteJugada.hacer(p)
	}

	return
}

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

func (jugada tocarFaltaEnvido) eval(p *Partida) {
	p.Ronda.Envite.Estado = DESHABILITADO

	// computar envidos
	jIdx, max, out := p.Ronda.execElEnvido()

	print(out)

	// jug es el que gano el (falta) envido
	jug := &p.jugadores[jIdx]

	pts := p.calcPtsFaltaEnvido(jug.Equipo)

	p.Ronda.Envite.Puntaje += pts

	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Sumar-Puntos",
			Nota: fmt.Sprintf(`La falta envido la gano %s con %v, +%v puntos para el equipo %s`,
				jug.Nombre, max, p.Ronda.Envite.Puntaje, jug.Equipo),
			Cont: ContSumPts{
				Pts:    p.Ronda.Envite.Puntaje,
				Equipo: jug.Equipo.String(),
			}.ToBytes(),
		},
	})

	p.sumarPuntos(jug.Equipo, p.Ronda.Envite.Puntaje)

}

type cantarFlor struct {
	Jugada
}

// fix
// todas las jugadas tienen que checkear, al principio
// si la flor esta en juego
// si esta -> no es posible realizar dicha jugada
// las unicas jugadas que quedan extentas de esta regla son:
// mazo, flor, contra flor, quiero (si se esta jugando la contra flor),
// no quiero ~ con flor me achico, contra flor al resto

/*
todo:
actualmente no permite que si todos cantan flor
se pase a calcular el resultado solo de las flores acumuladas
se necesita timer

*/
func (jugada cantarFlor) hacer(p *Partida) {
	// manojo dice que puede cantar flor;
	// es esto verdad?
	florHabilitada := (p.Ronda.Envite.Estado >= NOCANTADOAUN && p.Ronda.Envite.Estado <= FLOR) && p.Ronda.ManoEnJuego == primera
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := contains(p.Ronda.Envite.jugadoresConFlorQueNoCantaron, jugada.autor)
	ok := florHabilitada && tieneFlor && noCantoFlorAun

	if !ok {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf(`No es posible cantar flor`),
			},
		})

		return

	}

	// yo canto
	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Canta-Flor",
			Cont: []byte(jugada.autor.Jugador.Nombre),
		},
	})

	// y me elimino de los que no-cantaron
	p.Ronda.Envite.jugadoresConFlorQueNoCantaron = eliminar(p.Ronda.Envite.jugadoresConFlorQueNoCantaron, jugada.autor)

	yaEstabamosEnFlor := p.Ronda.Envite.Estado == FLOR

	if yaEstabamosEnFlor {

		p.Ronda.Envite.Puntaje += 3
		p.Ronda.Envite.CantadoPor = jugada.autor

	} else {

		// se usa por si dicen "no quiero" -> se obtiene el equipo
		// al que pertenece el que la canto en un principio para
		// poder sumarle los puntos correspondientes
		p.Ronda.Envite.Puntaje = 3
		p.Ronda.Envite.CantadoPor = jugada.autor
		p.Ronda.Envite.Estado = FLOR

	}

	// es el ultimo en cantar flor que faltaba?
	// o simplemente es el unico que tiene flor (caso particular)

	todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.jugadoresConFlorQueNoCantaron) == 0
	if todosLosJugadoresConFlorCantaron {

		evalFlor(p)

	} else {

		// cachear esto
		// solos los de su equipo tienen flor?
		// si solos los de su equipo tienen flor (y los otros no) -> las canto todas
		soloLosDeSuEquipoTienenFlor := true
		for _, manojo := range p.Ronda.Envite.jugadoresConFlor {
			if manojo.Jugador.Equipo != jugada.autor.Jugador.Equipo {
				soloLosDeSuEquipoTienenFlor = false
				break
			}
		}

		if soloLosDeSuEquipoTienenFlor {
			// los quiero llamar a todos, pero no quiero hacer llamadas al pedo
			// entonces: llamo al primero sin cantar, y que este llame al proximo
			// y que el proximo llame al siguiente, y asi...
			primero := p.Ronda.Envite.jugadoresConFlorQueNoCantaron[0]
			siguienteJugada := cantarFlor{Jugada{autor: primero}}
			siguienteJugada.hacer(p)
		}

	}

	return
}

func evalFlor(p *Partida) {
	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.jugadoresConFlorQueNoCantaron) == 0
	ok := todosLosJugadoresConFlorCantaron && florEnJuego
	if !ok {
		return
	}

	// cual es la flor ganadora?
	manojoConLaFlorMasAlta, maxFlor := p.Ronda.getLaFlorMasAlta()
	equipoGanador := manojoConLaFlorMasAlta.Jugador.Equipo

	// que estaba en juego?
	switch p.Ronda.Envite.Estado {
	case FLOR:
		// ahora se quien es el ganador; necesito saber cuantos puntos
		// se le va a sumar a ese equipo:
		// los acumulados del envite hasta ahora
		puntosASumar := p.Ronda.Envite.Puntaje
		p.sumarPuntos(equipoGanador, puntosASumar)
		habiaSolo1JugadorConFlor := len(p.Ronda.Envite.jugadoresConFlor) == 1
		if habiaSolo1JugadorConFlor {

			write(p.Stdout, &Pkt{
				Dest: []string{"ALL"},
				Msg: Msg{
					Tipo: "Sumar-Puntos",
					Nota: fmt.Sprintf(`+%v puntos para el equipo %s (por ser la unica flor de esta ronda)`,
						puntosASumar, equipoGanador),
					Cont: ContSumPts{
						Pts:    puntosASumar,
						Equipo: equipoGanador.String(),
					}.ToBytes(),
				},
			})

		} else {

			write(p.Stdout, &Pkt{
				Dest: []string{"ALL"},
				Msg: Msg{
					Tipo: "Sumar-Puntos",
					Nota: fmt.Sprintf(`La flor mas alta es la de %s con %v, +%v puntos para el equipo %s`,
						manojoConLaFlorMasAlta.Jugador.Nombre, maxFlor, puntosASumar, equipoGanador),
					Cont: ContSumPts{
						Pts:    puntosASumar,
						Equipo: equipoGanador.String(),
					}.ToBytes(),
				},
			})

		}
	case CONTRAFLOR:
	case CONTRAFLORALRESTO:
	}

	p.Ronda.Envite.Estado = DESHABILITADO
}

type cantarContraFlor struct {
	Jugada
}

func (jugada cantarContraFlor) hacer(p *Partida) {
	// manojo dice que puede cantar flor;
	// es esto verdad?
	contraFlorHabilitada := p.Ronda.Envite.Estado == FLOR && p.Ronda.ManoEnJuego == primera
	esDelEquipoContrario := contraFlorHabilitada && p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := contains(p.Ronda.Envite.jugadoresConFlorQueNoCantaron, jugada.autor)
	ok := contraFlorHabilitada && tieneFlor && esDelEquipoContrario && noCantoFlorAun
	if !ok {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf(`No es posible cantar contra flor`),
			},
		})

		return

	}

	// la canta
	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Canta-ContraFlor",
			Cont: []byte(jugada.autor.Jugador.Nombre),
		},
	})

	p.Ronda.Envite.Estado = CONTRAFLOR
	p.Ronda.Envite.CantadoPor = jugada.autor
	// ahora la flor pasa a jugarse por 4 puntos
	p.Ronda.Envite.Puntaje = 4
	// y ahora tengo que esperar por la respuesta de la nueva
	// propuesta de todos menos de el que canto la contraflor
	// restauro la copia
	p.Ronda.Envite.jugadoresConFlorQueNoCantaron = eliminar(p.Ronda.Envite.jugadoresConFlor, jugada.autor)

	return
}

type cantarContraFlorAlResto struct {
	Jugada
}

func (jugada cantarContraFlorAlResto) hacer(p *Partida) {
	// manojo dice que puede cantar flor;
	// es esto verdad?
	contraFlorHabilitada := (p.Ronda.Envite.Estado == FLOR || p.Ronda.Envite.Estado == CONTRAFLOR) && p.Ronda.ManoEnJuego == primera
	esDelEquipoContrario := contraFlorHabilitada && p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := contains(p.Ronda.Envite.jugadoresConFlorQueNoCantaron, jugada.autor)
	ok := contraFlorHabilitada && tieneFlor && esDelEquipoContrario && noCantoFlorAun
	if !ok {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf(`No es posible cantar contra flor al resto`),
			},
		})

		return

	}

	// la canta
	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Canta-ContraFlorAlResto",
			Cont: []byte(jugada.autor.Jugador.Nombre),
		},
	})

	p.Ronda.Envite.Estado = CONTRAFLORALRESTO
	p.Ronda.Envite.CantadoPor = jugada.autor
	// ahora la flor pasa a jugarse por 4 puntos
	p.Ronda.Envite.Puntaje = 4
	// y ahora tengo que esperar por la respuesta de la nueva
	// propuesta de todos menos de el que canto la contraflor
	// restauro la copia
	p.Ronda.Envite.jugadoresConFlorQueNoCantaron = eliminar(p.Ronda.Envite.jugadoresConFlor, jugada.autor)

	return
}

type cantarConFlorMeAchico struct {
	Jugada
}

func (jugada cantarConFlorMeAchico) hacer(p *Partida) {
	return
}

type gritarTruco struct {
	Jugada
}

func (jugada gritarTruco) hacer(p *Partida) {
	// checkeos:
	noSeFueAlMazo := jugada.autor.SeFueAlMazo == false
	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= NOCANTADOAUN
	hayFlor := len(p.Ronda.Envite.jugadoresConFlorQueNoCantaron) > 0
	noSeCantoFlor := p.Ronda.Envite.Estado > DESHABILITADO && p.Ronda.Envite.Estado < FLOR
	laFlorEstaPrimero := hayFlor && noSeCantoFlor
	trucoNoSeJugoAun := p.Ronda.Truco.Estado == NOCANTADO
	// esSuTurno := p.Ronda.getElTurno() == jugada.autor
	trucoHabilitado := noSeFueAlMazo && trucoNoSeJugoAun && noSeEstaJugandoElEnvite && !laFlorEstaPrimero // && esSuTurno

	if !trucoHabilitado {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible cantar truco ahora"),
			},
		})

		if laFlorEstaPrimero {
			manojosConFlor := p.Ronda.Envite.jugadoresConFlorQueNoCantaron
			siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
			siguienteJugada.hacer(p)
		}

		return

	}

	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Grita-Truco",
			Cont: []byte(jugada.autor.Jugador.Nombre),
		},
	})

	p.Ronda.Truco.CantadoPor = jugada.autor
	p.Ronda.Truco.Estado = TRUCO
	p.Ronda.Envite.Estado = DESHABILITADO

	return
}

type gritarReTruco struct {
	Jugada
}

// checkeaos de este tipo:
// que pasa cuando gritan re-truco cuando el campo truco se encuentra nil
// ese fue el nil pointer exception
func (jugada gritarReTruco) hacer(p *Partida) {

	// checkeos generales:
	noSeFueAlMazo := jugada.autor.SeFueAlMazo == false
	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= NOCANTADOAUN
	hayFlor := len(p.Ronda.Envite.jugadoresConFlorQueNoCantaron) > 0
	noSeCantoFlor := p.Ronda.Envite.Estado > DESHABILITADO && p.Ronda.Envite.Estado < FLOR
	laFlorEstaPrimero := hayFlor && noSeCantoFlor

	/*
		Hay 2 casos para cantar rectruco:
		    - CASO I: Uno del equipo contrario grito el truco
			- CASO II: Uno de su equipo posee el quiero
	*/

	// CASO I:
	trucoGritado := p.Ronda.Truco.Estado == TRUCO
	unoDelEquipoContrarioGritoTruco := trucoGritado && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	casoI := trucoGritado && unoDelEquipoContrarioGritoTruco

	// CASO I:
	trucoYaQuerido := p.Ronda.Truco.Estado == TRUCOQUERIDO
	unoDeMiEquipoQuizo := trucoYaQuerido && p.Ronda.Truco.CantadoPor.Jugador.Equipo == jugada.autor.Jugador.Equipo
	// esTurnoDeMiEquipo := p.Ronda.getElTurno().Jugador.Equipo == jugada.autor.Jugador.Equipo
	casoII := trucoYaQuerido && unoDeMiEquipoQuizo // && esTurnoDeMiEquipo

	reTrucoHabilitado := noSeFueAlMazo && noSeEstaJugandoElEnvite && (casoI || casoII) && !laFlorEstaPrimero

	if !reTrucoHabilitado {

		if laFlorEstaPrimero {
			manojosConFlor := p.Ronda.Envite.jugadoresConFlorQueNoCantaron
			siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
			siguienteJugada.hacer(p)
		}

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible cantar re-truco ahora"),
			},
		})

		return

	}

	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Grita-ReTruco",
			Cont: []byte(jugada.autor.Jugador.Nombre),
		},
	})

	p.Ronda.Truco.CantadoPor = jugada.autor
	p.Ronda.Truco.Estado = RETRUCO

	return
}

type gritarVale4 struct {
	Jugada
}

func (jugada gritarVale4) hacer(p *Partida) {
	// checkeos:
	noSeFueAlMazo := jugada.autor.SeFueAlMazo == false

	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= NOCANTADOAUN
	hayFlor := len(p.Ronda.Envite.jugadoresConFlorQueNoCantaron) > 0
	noSeCantoFlor := p.Ronda.Envite.Estado > DESHABILITADO && p.Ronda.Envite.Estado < FLOR
	laFlorEstaPrimero := hayFlor && noSeCantoFlor

	/*
		Hay 2 casos para cantar rectruco:
		    - CASO I: Uno del equipo contrario grito el re-truco
			- CASO II: Uno de su equipo posee el quiero
	*/

	// CASO I:
	reTrucoGritado := p.Ronda.Truco.Estado == RETRUCO
	// para eviat el nil primero checkeo que haya sido gritado reTrucoGritado &&
	unoDelEquipoContrarioGritoReTruco := reTrucoGritado && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	casoI := reTrucoGritado && unoDelEquipoContrarioGritoReTruco

	// CASO I:
	retrucoYaQuerido := p.Ronda.Truco.Estado == RETRUCOQUERIDO
	// para eviat el nil primero checkeo que haya sido gritado reTrucoGritado &&
	suEquipotieneElQuiero := retrucoYaQuerido && p.Ronda.Truco.CantadoPor.Jugador.Equipo == jugada.autor.Jugador.Equipo
	casoII := retrucoYaQuerido && suEquipotieneElQuiero

	vale4Habilitado := noSeFueAlMazo && (casoI || casoII) && noSeEstaJugandoElEnvite && !laFlorEstaPrimero

	if !vale4Habilitado {

		if laFlorEstaPrimero {
			manojosConFlor := p.Ronda.Envite.jugadoresConFlorQueNoCantaron
			siguienteJugada := cantarFlor{Jugada{autor: manojosConFlor[0]}}
			siguienteJugada.hacer(p)
		}

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible cantar vale-4 ahora"),
			},
		})

		return

	}

	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Grita-Vale4",
			Cont: []byte(jugada.autor.Jugador.Nombre),
		},
	})

	p.Ronda.Truco.CantadoPor = jugada.autor
	p.Ronda.Truco.Estado = VALE4

	return
}

type responderQuiero struct {
	Jugada
}

func (jugada responderQuiero) hacer(p *Partida) {
	seFueAlMazo := jugada.autor.SeFueAlMazo
	if seFueAlMazo {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("Te fuiste al mazo; no podes hacer esta jugada"),
			},
		})

		return

	}

	// checkeo flor en juego
	// caso particular del checkeo: no se le puede decir quiero a la flor
	// pero si a la contra flor o contra flor al resto
	florEnJuego := p.Ronda.Envite.Estado == FLOR
	if florEnJuego {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible responder quiero ahora"),
			},
		})

		return

	}
	// se acepta una respuesta 'quiero' solo cuando:
	// - CASO I: se toco un envite+ (con autor del equipo contario)
	// - CASO II: se grito el truco+ (con autor del equipo contario)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO)
	// ojo: solo a la contraflor+ se le puede decir quiero; a la flor sola no
	laContraFlorEsRespondible := p.Ronda.Envite.Estado >= CONTRAFLOR && p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
	elTrucoEsRespondible := contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado) && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo

	ok := elEnvidoEsRespondible || laContraFlorEsRespondible || elTrucoEsRespondible
	if !ok {
		// si no, esta respondiendo al pedo

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf(`No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo`),
			},
		})

		return

	}

	if elEnvidoEsRespondible {

		esDelEquipoContrario := jugada.getAutor().Jugador.Equipo != p.Ronda.Envite.CantadoPor.Jugador.Equipo
		if !esDelEquipoContrario {

			write(p.Stdout, &Pkt{
				Dest: []string{jugada.autor.Jugador.Nombre},
				Msg: Msg{
					Tipo: "Error",
					Nota: fmt.Sprintf(`La jugada no es valida`),
				},
			})

			return

		}

		write(p.Stdout, &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Responde-Quiero",
				Cont: []byte(jugada.autor.Jugador.Nombre),
			},
		})

		if p.Ronda.Envite.Estado == FALTAENVIDO {
			tocarFaltaEnvido{Jugada{autor: jugada.autor}}.eval(p)
			return
		}
		// si no, era envido/real-envido o cualquier
		// combinacion valida de ellos
		tocarEnvido{Jugada{autor: jugada.autor}}.eval(p)
		return

	} else if laContraFlorEsRespondible {
		// tengo que verificar si efectivamente tiene flor
		tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
		esDelEquipoContrario := jugada.getAutor().Jugador.Equipo != p.Ronda.Envite.CantadoPor.Jugador.Equipo
		ok := tieneFlor && esDelEquipoContrario

		if !ok {

			write(p.Stdout, &Pkt{
				Dest: []string{jugada.autor.Jugador.Nombre},
				Msg: Msg{
					Tipo: "Error",
					Nota: fmt.Sprintf(`La jugada no es valida`),
				},
			})

			return

		}

		write(p.Stdout, &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Responde-Quiero",
				Cont: []byte(jugada.autor.Jugador.Nombre),
			},
		})

		// empieza cantando el autor del envite no el que "quizo"
		manojoConLaFlorMasAlta, maxFlor := p.Ronda.getLaFlorMasAlta()
		equipoGanador := manojoConLaFlorMasAlta.Jugador.Equipo

		if p.Ronda.Envite.Estado == CONTRAFLOR {
			puntosASumar := p.Ronda.Envite.Puntaje
			p.sumarPuntos(equipoGanador, puntosASumar)

			write(p.Stdout, &Pkt{
				Dest: []string{"ALL"},
				Msg: Msg{
					Tipo: "Sumar-Puntos",
					Nota: fmt.Sprintf(`La contra-flor la gano %s con %v, +%v puntos para el equipo %s`,
						manojoConLaFlorMasAlta.Jugador.Nombre, maxFlor, puntosASumar, equipoGanador),
					Cont: ContSumPts{
						Pts:    puntosASumar,
						Equipo: equipoGanador.String(),
					}.ToBytes(),
				},
			})

		} else {
			// el equipo del ganador de la contraflor al resto
			// gano la partida
			// duda se cuentan las flores?
			// puntosASumar := p.Ronda.Envite.Puntaje + p.calcPtsContraFlorAlResto(equipoGanador)
			puntosASumar := p.calcPtsContraFlorAlResto(equipoGanador)
			p.sumarPuntos(equipoGanador, puntosASumar)

			write(p.Stdout, &Pkt{
				Dest: []string{"ALL"},
				Msg: Msg{
					Tipo: "Sumar-Puntos",
					Nota: fmt.Sprintf(`La contra-flor-al-resto la gano %s con %v, +%v puntos para el equipo %s`,
						manojoConLaFlorMasAlta.Jugador.Nombre, maxFlor, puntosASumar, equipoGanador),
					Cont: ContSumPts{
						Pts:    puntosASumar,
						Equipo: equipoGanador.String(),
					}.ToBytes(),
				},
			})

		}

		p.Ronda.Envite.Estado = DESHABILITADO

	} else if elTrucoEsRespondible {

		write(p.Stdout, &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Responde-Quiero",
				Cont: []byte(jugada.autor.Jugador.Nombre),
			},
		})

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

	return

}

type responderNoQuiero struct {
	Jugada
}

func (jugada responderNoQuiero) hacer(p *Partida) {

	seFueAlMazo := jugada.autor.SeFueAlMazo
	if seFueAlMazo {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("Te fuiste al mazo; no podes hacer esta jugada"),
			},
		})

		return
	}

	// checkeo flor en juego
	// caso particular del checkeo: no se le puede decir quiero a la flor
	// pero si a la contra flor o contra flor al resto
	// FALSO porque el no quiero lo estoy contando como un "con flor me achico"
	// todo: agregar la jugada: "con flor me achico" y editar la variale:
	// AHORA:
	// laFlorEsRespondible := p.Ronda.Flor >= FLOR && p.Ronda.Envite.CantadoPor.Jugador.equipo != jugada.autor.Jugador.Equipo
	// LUEGO DE AGREGAR LA JUGADA "con flor me achico"
	// laFlorEsRespondible := p.Ronda.Flor > FLOR
	// FALSO ---> directamente se va la posibilidad de reponderle
	// "no quiero a la flor"

	// se acepta una respuesta 'no quiero' solo cuando:
	// - CASO I: se toco el envido (o similar)
	// - CASO II: se grito el truco (o similar)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO) && p.Ronda.Envite.CantadoPor != jugada.autor
	laFlorEsRespondible := p.Ronda.Envite.Estado >= FLOR && p.Ronda.Envite.CantadoPor != jugada.autor
	elTrucoEsRespondible := contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado) && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo

	ok := elEnvidoEsRespondible || laFlorEsRespondible || elTrucoEsRespondible

	if !ok {
		// si no, esta respondiendo al pedo

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf(`%s esta respondiendo al pedo; no hay nada respondible`, jugada.autor.Jugador.Nombre),
			},
		})

		return

	}

	if elEnvidoEsRespondible {

		esDelEquipoContrario := jugada.getAutor().Jugador.Equipo != p.Ronda.Envite.CantadoPor.Jugador.Equipo
		if !esDelEquipoContrario {

			write(p.Stdout, &Pkt{
				Dest: []string{jugada.autor.Jugador.Nombre},
				Msg: Msg{
					Tipo: "Error",
					Nota: fmt.Sprintf(`La jugada no es valida`),
				},
			})

			return

		}

		write(p.Stdout, &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Responde-NoQuiero",
				Cont: []byte(jugada.autor.Jugador.Nombre),
			},
		})

		//	no se toma en cuenta el puntaje total del ultimo toque

		var totalPts int

		switch p.Ronda.Envite.Estado {
		case ENVIDO:
			totalPts = p.Ronda.Envite.Puntaje - 1
		case REALENVIDO:
			totalPts = p.Ronda.Envite.Puntaje - 2
		case FALTAENVIDO:
			totalPts = p.Ronda.Envite.Puntaje + 1
		}

		p.Ronda.Envite.Estado = DESHABILITADO
		p.Ronda.Envite.Puntaje = totalPts

		write(p.Stdout, &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Sumar-Puntos",
				Nota: fmt.Sprintf(`+%v puntos para el equipo %s`,
					totalPts, p.Ronda.Envite.CantadoPor.Jugador.Equipo),
				Cont: ContSumPts{
					Pts:    totalPts,
					Equipo: p.Ronda.Envite.CantadoPor.Jugador.Equipo.String(),
				}.ToBytes(),
			},
		})

		p.sumarPuntos(p.Ronda.Envite.CantadoPor.Jugador.Equipo, totalPts)

	} else if laFlorEsRespondible {

		// tengo que verificar si efectivamente tiene flor
		tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
		esDelEquipoContrario := jugada.getAutor().Jugador.Equipo != p.Ronda.Envite.CantadoPor.Jugador.Equipo
		ok := tieneFlor && esDelEquipoContrario

		if !ok {

			write(p.Stdout, &Pkt{
				Dest: []string{jugada.autor.Jugador.Nombre},
				Msg: Msg{
					Tipo: "Error",
					Nota: fmt.Sprintf(`La jugada no es valida`),
				},
			})

			return

		}

		// todo ok: tiene flor; se pasa a jugar:
		write(p.Stdout, &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Responde-ConFlorMeAchico",
				Cont: []byte(jugada.autor.Jugador.Nombre),
			},
		})

		// cuenta como un "no quiero" (codigo copiado)
		// segun el estado de la apuesta actual:
		// los "me achico" no cuentan para la flor
		// Flor		xcg(+3) / xcg(+3)
		// Flor + Contra-Flor		xc(+3) / xCadaFlorDelQueHizoElDesafio(+3) + 1
		// Flor + [Contra-Flor] + ContraFlorAlResto		~Falta Envido + *TODAS* las flores no achicadas / xcg(+3) + 1

		// sumo todas las flores del equipo contrario
		totalPts := 0

		for _, m := range p.Ronda.Manojos {
			esDelEquipoContrario := p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
			tieneFlor, _ := m.tieneFlor(p.Ronda.Muestra)
			if tieneFlor && esDelEquipoContrario {
				totalPts += 3
			}
		}

		if p.Ronda.Envite.Estado == CONTRAFLOR || p.Ronda.Envite.Estado == CONTRAFLORALRESTO {
			// si es contraflor o al resto
			// se suma 1 por el `no quiero`
			totalPts++
		}

		p.Ronda.Envite.Estado = DESHABILITADO

		write(p.Stdout, &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Sumar-Puntos",
				Nota: fmt.Sprintf(`+%v puntos para el equipo %s por las flores`,
					totalPts, p.Ronda.Envite.CantadoPor.Jugador.Equipo),
				Cont: ContSumPts{
					Pts:    totalPts,
					Equipo: p.Ronda.Envite.CantadoPor.Jugador.Equipo.String(),
				}.ToBytes(),
			},
		})

		p.sumarPuntos(p.Ronda.Envite.CantadoPor.Jugador.Equipo, totalPts)

	} else if elTrucoEsRespondible {

		write(p.Stdout, &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Responde-NoQuiero",
				Cont: []byte(jugada.autor.Jugador.Nombre),
			},
		})

		// pongo al equipo que propuso el truco como ganador de la mano actual
		manoActual := p.Ronda.ManoEnJuego.toInt() - 1
		p.Ronda.Manos[manoActual].Ganador = p.Ronda.Truco.CantadoPor

		nuevaRonda, pkt := p.evaluarRonda()

		if pkt != nil {
			// antes:
			//write(p.Stdout, msg)
			// ahora:
			write(p.Stdout, pkt)
		}

		if nuevaRonda {

			if !p.Terminada() {
				// ahora se deberia de incrementar el mano
				// y ser el turno de este
				sigMano := p.Ronda.getSigElMano()
				p.nuevaRonda(sigMano) // todo: el tema es que cuando llama aca
				// no manda mensaje de que arranco nueva ronda
				// falso: el padre que llama a .evaluarRonda tiene que fijarse si
				// retorno true
				// entonces debe crearla el
				// no es responsabilidad de evaluarRonda arrancar una ronda nueva!!
				// de hecho, si una ronda es terminable y se llama 2 veces consecutivas
				// al mismo metodo booleano, en ambas oportunidades retorna diferente
				// ridiculo
				write(p.Stdout, &Pkt{
					Dest: []string{"ALL"},
					Msg: Msg{
						Tipo: "Nueva-Ronda",
						Cont: nil,
					},
				})

			} else {
				p.byeBye()
			}

		}

	}

	return
}

type irseAlMazo struct {
	Jugada
}

func (jugada irseAlMazo) hacer(p *Partida) {
	// checkeos:
	yaSeFueAlMazo := jugada.autor.SeFueAlMazo == true
	yaTiroTodasSusCartas := jugada.autor.getCantCartasTiradas() == 3
	if yaSeFueAlMazo || yaTiroTodasSusCartas {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible irse al mazo ahora"),
			},
		})

		return

	}

	seEstabaJugandoElEnvido := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO)
	seEstabaJugandoLaFlor := p.Ronda.Envite.Estado >= FLOR
	seEstabaJugandoElTruco := contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado)
	// no se puede ir al mazo sii:
	// 1. el fue el que canto el envido (y el envido esta en juego)
	// 2. tampoco se puede ir al mazo si el canto la flor o similar
	// 3. tampoco se puede ir al mazo si el grito el truco

	// envidoPropuesto := contains([]EstadoEnvite{ENVIDO, REALENVIDO, FALTAENVIDO}, p.Ronda.Envite.Estado)
	// envidoPropuestoPorSuEquipo := p.Ronda.Envite.CantadoPor.Jugador.Equipo == jugada.autor.Jugador.Equipo
	// trucoPropuesto := contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado)
	// trucoPropuestoPorSuEquipo := p.Ronda.Truco.CantadoPor.Jugador.Equipo == jugada.autor.Jugador.Equipo
	// condicionDelBobo := (envidoPropuesto && envidoPropuestoPorSuEquipo) || (trucoPropuesto && trucoPropuestoPorSuEquipo)

	// if condicionDelBobo {

	// write(p.Stdout, &Pkt{
	// 	Dest: []string{jugada.autor.Jugador.Nombre},
	// 	Msg: Msg{
	// 		Tipo: "Error",
	// 		Nota: fmt.Sprintf("No es posible irse al mazo ahora porque hay propuestas de tu equipo sin responder"),
	// 	},
	// })

	// return

	// }

	noSePuedeIrPorElEnvite := (seEstabaJugandoElEnvido || seEstabaJugandoLaFlor) && p.Ronda.Envite.CantadoPor == jugada.autor
	// la de la flor es igual al del envido; porque es un envite
	noSePuedeIrPorElTruco := seEstabaJugandoElTruco && p.Ronda.Truco.CantadoPor == jugada.autor
	if noSePuedeIrPorElEnvite || noSePuedeIrPorElTruco {

		write(p.Stdout, &Pkt{
			Dest: []string{jugada.autor.Jugador.Nombre},
			Msg: Msg{
				Tipo: "Error",
				Nota: fmt.Sprintf("No es posible irse al mazo ahora"),
			},
		})

		return

	}

	// ok -> se va al mazo:
	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Mazo",
			Cont: []byte(jugada.autor.Jugador.Nombre),
		},
	})

	jugada.autor.SeFueAlMazo = true
	equipoDelJugador := jugada.autor.Jugador.Equipo
	p.Ronda.CantJugadoresEnJuego[equipoDelJugador]--
	seFueronTodos := p.Ronda.CantJugadoresEnJuego[equipoDelJugador] == 0

	// si tenia flor -> ya no lo tomo en cuenta
	tieneFlor, _ := jugada.autor.tieneFlor(p.Ronda.Muestra)
	if tieneFlor {
		p.Ronda.Envite.jugadoresConFlor = eliminar(p.Ronda.Envite.jugadoresConFlor, jugada.autor)
		p.Ronda.Envite.jugadoresConFlorQueNoCantaron = eliminar(p.Ronda.Envite.jugadoresConFlorQueNoCantaron, jugada.autor)
		// que pasa si era el ultimo que se esperaba que cantara flor?
		// tengo que hacer el eval de la flor
		todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.jugadoresConFlorQueNoCantaron) == 0
		if todosLosJugadoresConFlorCantaron {
			evalFlor(p)
		}
	}

	// era el ultimo en tirar de esta mano?
	eraElUltimoEnTirar := p.Ronda.getSigHabilitado(*jugada.autor) == nil

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
			e := &p.Ronda.Envite
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

			write(p.Stdout, &Pkt{
				Dest: []string{"ALL"},
				Msg: Msg{
					Tipo: "Sumar-Puntos",
					Nota: fmt.Sprintf(`+%v puntos del envite para el equipo %s`,
						totalPts, e.CantadoPor.Jugador.Equipo),
					Cont: ContSumPts{
						Pts:    totalPts,
						Equipo: e.CantadoPor.Jugador.Equipo.String(),
					}.ToBytes(),
				},
			})

			p.sumarPuntos(p.Ronda.Envite.CantadoPor.Jugador.Equipo, totalPts)

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
				esDelEquipoContrario := p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.autor.Jugador.Equipo
				tieneFlor, _ := m.tieneFlor(p.Ronda.Muestra)
				if tieneFlor && esDelEquipoContrario {
					totalPts += 3
				}
			}

			if p.Ronda.Envite.Estado == CONTRAFLOR || p.Ronda.Envite.Estado == CONTRAFLORALRESTO {
				// si es contraflor o al resto
				// se suma 1 por el `no quiero`
				totalPts++
			}

			p.Ronda.Envite.Estado = DESHABILITADO

			write(p.Stdout, &Pkt{
				Dest: []string{"ALL"},
				Msg: Msg{
					Tipo: "Sumar-Puntos",
					Nota: fmt.Sprintf(`+%v puntos para el equipo %s por las flores`,
						totalPts, p.Ronda.Envite.CantadoPor.Jugador.Equipo),
					Cont: ContSumPts{
						Pts:    totalPts,
						Equipo: p.Ronda.Envite.CantadoPor.Jugador.Equipo.String(),
					}.ToBytes(),
				},
			})

		}
	}

	// evaluar ronda sii:
	// o bien se fueron todos
	// o bien este se fue al mazo, pero alguno de sus companeros no
	// (es decir que queda al menos 1 jugador en juego)
	hayQueEvaluarRonda := seFueronTodos || eraElUltimoEnTirar
	if hayQueEvaluarRonda {
		// de ser asi tengo que checkear el resultado de la mano
		// el turno del siguiente queda dado por el ganador de esta
		empiezaNuevaRonda, pkts := p.evaluarMano()

		if pkts != nil {
			for _, pkt := range pkts {
				if pkt != nil {

					// antes:
					//write(p.Stdout, msg)
					// ahora
					write(p.Stdout, pkt)

				}
			}
		}

		if !empiezaNuevaRonda {

			p.Ronda.ManoEnJuego++
			p.Ronda.setNextTurnoPosMano()

		} else {

			if !p.Terminada() {
				// ahora se deberia de incrementar el mano
				// y ser el turno de este
				sigMano := p.Ronda.getSigElMano()
				p.nuevaRonda(sigMano) // todo: el tema es que cuando llama aca
				// no manda mensaje de que arranco nueva ronda
				// falso: el padre que llama a .evaluarRonda tiene que fijarse si
				// retorno true
				// entonces debe crearla el
				// no es responsabilidad de evaluarRonda arrancar una ronda nueva!!
				// de hecho, si una ronda es terminable y se llama 2 veces consecutivas
				// al mismo metodo booleano, en ambas oportunidades retorna diferente
				// ridiculo
				write(p.Stdout, &Pkt{
					Dest: []string{"ALL"},
					Msg: Msg{
						Tipo: "Nueva-Ronda",
						Cont: nil,
					},
				})

			} else {
				p.byeBye()
			}

		}
	} else {
		// cambio de turno solo si era su turno
		eraSuTurno := p.Ronda.getElTurno() == jugada.autor
		if eraSuTurno {
			p.Ronda.setNextTurno()
		}
	}

	return
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
		fmt.Printf("%s: ", tipoJugada) //
		for _, jugada := range opciones {
			fmt.Printf("%s, ", jugada) //
		}
		fmt.Printf("\n") //
	}
	fmt.Println()
}
