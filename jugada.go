package truco

import (
	"fmt"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/pdt"
)

// IJugada Interface para las jugadas
type IJugada interface {
	hacer(p *Partida)
}

type tirarCarta struct {
	*pdt.Manojo
	pdt.Carta
}

// el jugador tira una carta;
// el parametro se encuentra en la struct como atributo
func (jugada tirarCarta) hacer(p *Partida) {

	// checkeo si se fue al mazo
	noSeFueAlMazo := jugada.Manojo.SeFueAlMazo == false
	ok := noSeFueAlMazo
	if !ok {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible tirar una carta porque ya te fuiste al mazo"),
		))

		return

	}

	// esto es un tanto redundante porque es imposible que no sea su turno
	// (checkeado mas adelante) y que al mismo tiempo tenga algo para tirar
	// luego de haber jugado sus 3 cartas; aun asi lo dejo
	yaTiroTodasSusCartas := jugada.Manojo.GetCantCartasTiradas() == 3
	if yaTiroTodasSusCartas {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible tirar una carta porque ya las tiraste todas"),
		))

		return

	}

	// checkeo flor en juego
	enviteEnJuego := p.Ronda.Envite.Estado >= pdt.ENVIDO
	if enviteEnJuego {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible tirar una carta ahora porque el envite esta en juego"),
		))

		return

	}

	// primero que nada: tiene esa carta?
	idx, err := jugada.Manojo.GetCartaIdx(jugada.Carta)
	if err != nil {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, err.Error()),
		))

		return
	}

	// ya jugo esa carta?
	todaviaNoLaTiro := jugada.Manojo.CartasNoTiradas[idx]
	if !todaviaNoLaTiro {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "Ya tiraste esa carta"),
		))

		return
	}

	// luego, era su turno?
	eraSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
	if !eraSuTurno {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No era su turno, no puede tirar la carta"),
		))

		return

	}

	// checkeo si tiene flor
	florHabilitada := (p.Ronda.Envite.Estado >= pdt.NOCANTADOAUN && p.Ronda.Envite.Estado <= pdt.FLOR) && p.Ronda.ManoEnJuego == pdt.Primera
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := pdt.Contains(p.Ronda.Envite.JugadoresConFlorQueNoCantaron, jugada.Manojo)
	noPuedeTirar := florHabilitada && tieneFlor && noCantoFlorAun
	if noPuedeTirar {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible tirar una carta sin antes cantar la flor"),
		))

		return

	}

	trucoGritado := pdt.Contains([]pdt.EstadoTruco{pdt.TRUCO, pdt.RETRUCO, pdt.VALE4}, p.Ronda.Truco.Estado)
	unoDelEquipoContrarioGritoTruco := trucoGritado && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	elTrucoEsRespondible := trucoGritado && unoDelEquipoContrarioGritoTruco
	if elTrucoEsRespondible {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible tirar una carta porque tu equipo debe responder la propuesta del truco"),
		))

		return

	}

	// ok la tiene y era su turno -> la juega
	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TirarCarta,
			jugada.Manojo.Jugador.Nombre, int(jugada.Carta.Palo), jugada.Carta.Valor),
	))

	p.PartidaDT.TirarCarta(jugada.Manojo, idx)

	// era el ultimo en tirar de esta mano?
	eraElUltimoEnTirar := p.Ronda.GetSigHabilitado(*jugada.Manojo) == nil
	if eraElUltimoEnTirar {
		// de ser asi tengo que checkear el resultado de la mano
		empiezaNuevaRonda, pkts := p.EvaluarMano()

		if pkts != nil {
			for _, pkt := range pkts {
				if pkt != nil {

					// antes:
					//write(p.Stdout, msg)
					// ahora:
					enco.Write(p.out, pkt)

				}
			}
		}

		if !empiezaNuevaRonda {

			// actualizo el mano
			p.Ronda.ManoEnJuego++
			p.Ronda.SetNextTurnoPosMano()
			// lo envio
			enco.Write(p.out, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SigTurnoPosMano, int(p.Ronda.Turno)),
			))

		} else {

			if !p.Terminada() {
				// ahora se deberia de incrementar el mano
				// y ser el turno de este
				sigMano := p.Ronda.GetSigElMano()
				p.NuevaRonda(sigMano) // todo: el tema es que cuando llama aca
				// no manda mensaje de que arranco nueva ronda
				// falso: el padre que llama a .EvaluarRonda tiene que fijarse si
				// retorno true
				// entonces debe crearla el
				// no es responsabilidad de EvaluarRonda arrancar una ronda nueva!!
				// de hecho, si una ronda es terminable y se llama 2 veces consecutivas
				// al mismo metodo booleano, en ambas oportunidades retorna diferente
				// ridiculo

				for _, m := range p.Ronda.Manojos {
					enco.Write(p.out, enco.Pkt(
						enco.Dest(m.Jugador.ID),
						enco.Msg(enco.NuevaRonda, p.PartidaDT.PerspectivaCacheFlor(&m)),
					))

				}

			} else {
				p.byeBye()
			}

		}

		// el turno del siguiente queda dado por el ganador de esta
	} else {
		p.Ronda.SetNextTurno()

		enco.Write(p.out, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.SigTurno, int(p.Ronda.Turno)),
		))

	}

	return
}

// PRE: supongo que el jugador que toca este envido
// no tiene flor (es checkeada cuando es su turno)
type tocarEnvido struct {
	*pdt.Manojo
}

func (jugada tocarEnvido) hacer(p *Partida) {
	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= pdt.FLOR
	if florEnJuego {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible tocar el envido ahora porque la flor esta en juego"),
		))

		return

	}
	esPrimeraMano := p.Ronda.ManoEnJuego == pdt.Primera
	esSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	esDelEquipoContrario := p.Ronda.Envite.Estado == pdt.NOCANTADOAUN || p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	envidoHabilitado := (p.Ronda.Envite.Estado == pdt.NOCANTADOAUN || p.Ronda.Envite.Estado == pdt.ENVIDO)
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == pdt.ENVIDO
	apuestaSaturada := p.Ronda.Envite.Puntaje >= p.CalcPtsFalta()
	ok := (envidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario) && (esSuTurno || yaEstabamosEnEnvido) && !apuestaSaturada

	if !ok {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, `No es posible cantar 'Envido'`),
		))

		return

	}

	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TocarEnvido, jugada.Manojo.Jugador.Nombre),
	))

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.JugadoresConFlorQueNoCantaron) > 0
	if hayFlor {
		// todo: deberia ir al estado magico en el que espera
		// solo por jugadas de tipo flor-related
		// lo mismo para el real-envido; falta-envido
		manojosConFlor := p.Ronda.Envite.JugadoresConFlorQueNoCantaron
		siguienteJugada := cantarFlor{manojosConFlor[0]}
		siguienteJugada.hacer(p)

	} else {
		p.PartidaDT.TocarEnvido(jugada.Manojo)
	}

	return
}

// donde 'j' el jugador que dijo 'quiero' al 'envido'/'real envido'
func (jugada tocarEnvido) eval(p *Partida) {
	p.Ronda.Envite.Estado = pdt.DESHABILITADO
	jIdx, _, pkts := p.Ronda.ExecElEnvido()

	if pkts != nil {
		for _, pkt := range pkts {
			if pkt != nil {
				enco.Write(p.out, pkt)
			}
		}
	}

	jug := &p.Jugadores[jIdx]

	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.SumaPts, jug.Nombre, enco.EnvidoGanado, p.Ronda.Envite.Puntaje),
	))

	p.SumarPuntos(jug.Equipo, p.Ronda.Envite.Puntaje)

}

type tocarRealEnvido struct {
	*pdt.Manojo
}

func (jugada tocarRealEnvido) hacer(p *Partida) {
	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= pdt.FLOR
	if florEnJuego {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible tocar real envido ahora porque la flor esta en juego"),
		))

		return

	}
	esPrimeraMano := p.Ronda.ManoEnJuego == pdt.Primera
	esSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	realEnvidoHabilitado := (p.Ronda.Envite.Estado == pdt.NOCANTADOAUN || p.Ronda.Envite.Estado == pdt.ENVIDO)
	esDelEquipoContrario := p.Ronda.Envite.Estado == pdt.NOCANTADOAUN || p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == pdt.ENVIDO
	ok := realEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario && (esSuTurno || yaEstabamosEnEnvido)

	if !ok {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, `No es posible cantar 'Real Envido'`),
		))

		return

	}

	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TocarRealEnvido, jugada.Manojo.Jugador.Nombre),
	))

	p.PartidaDT.TocarRealEnvido(jugada.Manojo)

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.JugadoresConFlorQueNoCantaron) > 0

	if hayFlor {
		manojosConFlor := p.Ronda.Envite.JugadoresConFlorQueNoCantaron
		siguienteJugada := cantarFlor{manojosConFlor[0]}
		siguienteJugada.hacer(p)

	}

	return
}

type tocarFaltaEnvido struct {
	*pdt.Manojo
}

func (jugada tocarFaltaEnvido) hacer(p *Partida) {
	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= pdt.FLOR
	if florEnJuego {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible tocar falta envido ahora porque la flor esta en juego"),
		))

		return

	}

	esSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
	esPrimeraMano := p.Ronda.ManoEnJuego == pdt.Primera
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	faltaEnvidoHabilitado := p.Ronda.Envite.Estado >= pdt.NOCANTADOAUN && p.Ronda.Envite.Estado < pdt.FALTAENVIDO
	esDelEquipoContrario := p.Ronda.Envite.Estado == pdt.NOCANTADOAUN || p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == pdt.ENVIDO || p.Ronda.Envite.Estado == pdt.REALENVIDO
	ok := faltaEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario && (esSuTurno || yaEstabamosEnEnvido)

	if !ok {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, `No es posible cantar 'Falta Envido'`),
		))

		return

	}

	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TocarFaltaEnvido, jugada.Manojo.Jugador.Nombre),
	))

	p.PartidaDT.TocarFaltaEnvido(jugada.Manojo)

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.JugadoresConFlorQueNoCantaron) > 0
	if hayFlor {
		manojosConFlor := p.Ronda.Envite.JugadoresConFlorQueNoCantaron
		siguienteJugada := cantarFlor{manojosConFlor[0]}
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
	p.Ronda.Envite.Estado = pdt.DESHABILITADO

	// computar envidos
	jIdx, _, pkts := p.Ronda.ExecElEnvido()

	if pkts != nil {
		for _, pkt := range pkts {
			if pkt != nil {
				enco.Write(p.out, pkt)
			}
		}
	}

	// jug es el que gano el (falta) envido
	jug := &p.Jugadores[jIdx]

	pts := p.CalcPtsFaltaEnvido(jug.Equipo)

	p.Ronda.Envite.Puntaje += pts

	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.SumaPts, jug.Nombre, enco.FaltaEnvidoGanado, p.Ronda.Envite.Puntaje),
	))

	p.SumarPuntos(jug.Equipo, p.Ronda.Envite.Puntaje)

}

type cantarFlor struct {
	*pdt.Manojo
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
	florHabilitada := (p.Ronda.Envite.Estado >= pdt.NOCANTADOAUN && p.Ronda.Envite.Estado <= pdt.FLOR) && p.Ronda.ManoEnJuego == pdt.Primera
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := pdt.Contains(p.Ronda.Envite.JugadoresConFlorQueNoCantaron, jugada.Manojo)
	ok := florHabilitada && tieneFlor && noCantoFlorAun

	if !ok {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, `No es posible cantar flor`),
		))

		return

	}

	// yo canto
	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.CantarFlor, jugada.Manojo.Jugador.Nombre),
	))

	// y me elimino de los que no-cantaron
	p.Ronda.Envite.JugadoresConFlorQueNoCantaron = pdt.Eliminar(p.Ronda.Envite.JugadoresConFlorQueNoCantaron, jugada.Manojo)

	p.PartidaDT.CantarFlor(jugada.Manojo)

	// es el ultimo en cantar flor que faltaba?
	// o simplemente es el unico que tiene flor (caso particular)

	todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.JugadoresConFlorQueNoCantaron) == 0
	if todosLosJugadoresConFlorCantaron {

		evalFlor(p)

	} else {

		// cachear esto
		// solos los de su equipo tienen flor?
		// si solos los de su equipo tienen flor (y los otros no) -> las canto todas
		soloLosDeSuEquipoTienenFlor := true
		for _, manojo := range p.Ronda.Envite.JugadoresConFlor {
			if manojo.Jugador.Equipo != jugada.Manojo.Jugador.Equipo {
				soloLosDeSuEquipoTienenFlor = false
				break
			}
		}

		if soloLosDeSuEquipoTienenFlor {
			// los quiero llamar a todos, pero no quiero hacer llamadas al pedo
			// entonces: llamo al primero sin cantar, y que este llame al proximo
			// y que el proximo llame al siguiente, y asi...
			primero := p.Ronda.Envite.JugadoresConFlorQueNoCantaron[0]
			siguienteJugada := cantarFlor{primero}
			siguienteJugada.hacer(p)
		}

	}

	return
}

func evalFlor(p *Partida) {
	florEnJuego := p.Ronda.Envite.Estado >= pdt.FLOR
	todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.JugadoresConFlorQueNoCantaron) == 0
	ok := todosLosJugadoresConFlorCantaron && florEnJuego
	if !ok {
		return
	}

	// cual es la flor ganadora?
	manojoConLaFlorMasAlta, _ := p.Ronda.GetLaFlorMasAlta()
	equipoGanador := manojoConLaFlorMasAlta.Jugador.Equipo

	// que estaba en juego?
	switch p.Ronda.Envite.Estado {
	case pdt.FLOR:
		// ahora se quien es el ganador; necesito saber cuantos puntos
		// se le va a sumar a ese equipo:
		// los acumulados del envite hasta ahora
		puntosASumar := p.Ronda.Envite.Puntaje
		p.SumarPuntos(equipoGanador, puntosASumar)
		habiaSolo1JugadorConFlor := len(p.Ronda.Envite.JugadoresConFlor) == 1
		if habiaSolo1JugadorConFlor {

			enco.Write(p.out, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					manojoConLaFlorMasAlta.Jugador.ID,
					enco.LaUnicaFlor, puntosASumar),
			))

		} else {

			enco.Write(p.out, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					manojoConLaFlorMasAlta.Jugador.Nombre,
					enco.LaFlorMasAlta,
					puntosASumar),
			))

		}
	case pdt.CONTRAFLOR:
	case pdt.CONTRAFLORALRESTO:
	}

	p.Ronda.Envite.Estado = pdt.DESHABILITADO
}

type cantarContraFlor struct {
	*pdt.Manojo
}

func (jugada cantarContraFlor) hacer(p *Partida) {
	// manojo dice que puede cantar flor;
	// es esto verdad?
	contraFlorHabilitada := p.Ronda.Envite.Estado == pdt.FLOR && p.Ronda.ManoEnJuego == pdt.Primera
	esDelEquipoContrario := contraFlorHabilitada && p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := pdt.Contains(p.Ronda.Envite.JugadoresConFlorQueNoCantaron, jugada.Manojo)
	ok := contraFlorHabilitada && tieneFlor && esDelEquipoContrario && noCantoFlorAun
	if !ok {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, `No es posible cantar contra flor`),
		))

		return

	}

	// la canta
	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.CantarContraFlor, jugada.Manojo.Jugador.ID),
	))

	p.PartidaDT.CantarContraFlor(jugada.Manojo)
	// y ahora tengo que esperar por la respuesta de la nueva
	// propuesta de todos menos de el que canto la contraflor
	// restauro la copia
	p.Ronda.Envite.JugadoresConFlorQueNoCantaron = pdt.Eliminar(p.Ronda.Envite.JugadoresConFlor, jugada.Manojo)

	return
}

type cantarContraFlorAlResto struct {
	*pdt.Manojo
}

func (jugada cantarContraFlorAlResto) hacer(p *Partida) {
	// manojo dice que puede cantar flor;
	// es esto verdad?
	contraFlorHabilitada := (p.Ronda.Envite.Estado == pdt.FLOR || p.Ronda.Envite.Estado == pdt.CONTRAFLOR) && p.Ronda.ManoEnJuego == pdt.Primera
	esDelEquipoContrario := contraFlorHabilitada && p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := pdt.Contains(p.Ronda.Envite.JugadoresConFlorQueNoCantaron, jugada.Manojo)
	ok := contraFlorHabilitada && tieneFlor && esDelEquipoContrario && noCantoFlorAun
	if !ok {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, `No es posible cantar contra flor al resto`),
		))

		return

	}

	// la canta
	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.CantarContraFlorAlResto, jugada.Manojo.Jugador.ID),
	))

	p.PartidaDT.CantarContraFlorAlResto(jugada.Manojo)
	// y ahora tengo que esperar por la respuesta de la nueva
	// propuesta de todos menos de el que canto la contraflor
	// restauro la copia
	p.Ronda.Envite.JugadoresConFlorQueNoCantaron = pdt.Eliminar(p.Ronda.Envite.JugadoresConFlor, jugada.Manojo)

	return
}

type cantarConFlorMeAchico struct {
	*pdt.Manojo
}

func (jugada cantarConFlorMeAchico) hacer(p *Partida) {
	return
}

type gritarTruco struct {
	*pdt.Manojo
}

func (jugada gritarTruco) hacer(p *Partida) {
	// checkeos:
	noSeFueAlMazo := jugada.Manojo.SeFueAlMazo == false
	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= pdt.NOCANTADOAUN
	hayFlor := len(p.Ronda.Envite.JugadoresConFlorQueNoCantaron) > 0
	noSeCantoFlor := p.Ronda.Envite.Estado > pdt.DESHABILITADO && p.Ronda.Envite.Estado < pdt.FLOR
	laFlorEstaPrimero := hayFlor && noSeCantoFlor
	trucoNoSeJugoAun := p.Ronda.Truco.Estado == pdt.NOCANTADO
	// esSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
	trucoHabilitado := noSeFueAlMazo && trucoNoSeJugoAun && noSeEstaJugandoElEnvite && !laFlorEstaPrimero // && esSuTurno

	if !trucoHabilitado {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible cantar truco ahora"),
		))

		if laFlorEstaPrimero {
			manojosConFlor := p.Ronda.Envite.JugadoresConFlorQueNoCantaron
			siguienteJugada := cantarFlor{manojosConFlor[0]}
			siguienteJugada.hacer(p)
		}

		return

	}

	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.GritarTruco, jugada.Manojo.Jugador.ID),
	))

	p.PartidaDT.GritarTruco(jugada.Manojo)

	return
}

type gritarReTruco struct {
	*pdt.Manojo
}

// checkeaos de este tipo:
// que pasa cuando gritan re-truco cuando el campo truco se encuentra nil
// ese fue el nil pointer exception
func (jugada gritarReTruco) hacer(p *Partida) {

	// checkeos generales:
	noSeFueAlMazo := jugada.Manojo.SeFueAlMazo == false
	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= pdt.NOCANTADOAUN
	hayFlor := len(p.Ronda.Envite.JugadoresConFlorQueNoCantaron) > 0
	noSeCantoFlor := p.Ronda.Envite.Estado > pdt.DESHABILITADO && p.Ronda.Envite.Estado < pdt.FLOR
	laFlorEstaPrimero := hayFlor && noSeCantoFlor

	/*
		Hay 2 casos para cantar rectruco:
		    - CASO I: Uno del equipo contrario grito el truco
			- CASO II: Uno de su equipo posee el quiero
	*/

	// CASO I:
	trucoGritado := p.Ronda.Truco.Estado == pdt.TRUCO
	unoDelEquipoContrarioGritoTruco := trucoGritado && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	casoI := trucoGritado && unoDelEquipoContrarioGritoTruco

	// CASO I:
	trucoYaQuerido := p.Ronda.Truco.Estado == pdt.TRUCOQUERIDO
	unoDeMiEquipoQuizo := trucoYaQuerido && p.Ronda.Truco.CantadoPor.Jugador.Equipo == jugada.Manojo.Jugador.Equipo
	// esTurnoDeMiEquipo := p.Ronda.GetElTurno().Jugador.Equipo == jugada.Manojo.Jugador.Equipo
	casoII := trucoYaQuerido && unoDeMiEquipoQuizo // && esTurnoDeMiEquipo

	reTrucoHabilitado := noSeFueAlMazo && noSeEstaJugandoElEnvite && (casoI || casoII) && !laFlorEstaPrimero

	if !reTrucoHabilitado {

		if laFlorEstaPrimero {
			manojosConFlor := p.Ronda.Envite.JugadoresConFlorQueNoCantaron
			siguienteJugada := cantarFlor{manojosConFlor[0]}
			siguienteJugada.hacer(p)
		}

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible cantar re-truco ahora"),
		))

		return

	}

	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.GritarReTruco, jugada.Manojo.Jugador.ID),
	))

	p.PartidaDT.GritarReTruco(jugada.Manojo)

	return
}

type gritarVale4 struct {
	*pdt.Manojo
}

func (jugada gritarVale4) hacer(p *Partida) {
	// checkeos:
	noSeFueAlMazo := jugada.Manojo.SeFueAlMazo == false

	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= pdt.NOCANTADOAUN
	hayFlor := len(p.Ronda.Envite.JugadoresConFlorQueNoCantaron) > 0
	noSeCantoFlor := p.Ronda.Envite.Estado > pdt.DESHABILITADO && p.Ronda.Envite.Estado < pdt.FLOR
	laFlorEstaPrimero := hayFlor && noSeCantoFlor

	/*
		Hay 2 casos para cantar rectruco:
		    - CASO I: Uno del equipo contrario grito el re-truco
			- CASO II: Uno de su equipo posee el quiero
	*/

	// CASO I:
	reTrucoGritado := p.Ronda.Truco.Estado == pdt.RETRUCO
	// para eviat el nil primero checkeo que haya sido gritado reTrucoGritado &&
	unoDelEquipoContrarioGritoReTruco := reTrucoGritado && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	casoI := reTrucoGritado && unoDelEquipoContrarioGritoReTruco

	// CASO I:
	retrucoYaQuerido := p.Ronda.Truco.Estado == pdt.RETRUCOQUERIDO
	// para eviat el nil primero checkeo que haya sido gritado reTrucoGritado &&
	suEquipotieneElQuiero := retrucoYaQuerido && p.Ronda.Truco.CantadoPor.Jugador.Equipo == jugada.Manojo.Jugador.Equipo
	casoII := retrucoYaQuerido && suEquipotieneElQuiero

	vale4Habilitado := noSeFueAlMazo && (casoI || casoII) && noSeEstaJugandoElEnvite && !laFlorEstaPrimero

	if !vale4Habilitado {

		if laFlorEstaPrimero {
			manojosConFlor := p.Ronda.Envite.JugadoresConFlorQueNoCantaron
			siguienteJugada := cantarFlor{manojosConFlor[0]}
			siguienteJugada.hacer(p)
		}

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible cantar vale-4 ahora"),
		))

		return

	}

	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.GritarVale4, jugada.Manojo.Jugador.ID),
	))

	p.PartidaDT.GritarVale4(jugada.Manojo)

	return
}

type responderQuiero struct {
	*pdt.Manojo
}

func (jugada responderQuiero) hacer(p *Partida) {
	seFueAlMazo := jugada.Manojo.SeFueAlMazo
	if seFueAlMazo {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "Te fuiste al mazo; no podes hacer esta jugada"),
		))

		return

	}

	// checkeo flor en juego
	// caso particular del checkeo: no se le puede decir quiero a la flor
	// pero si a la contra flor o contra flor al resto
	florEnJuego := p.Ronda.Envite.Estado == pdt.FLOR
	if florEnJuego {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible responder quiero ahora"),
		))

		return

	}
	// se acepta una respuesta 'quiero' solo cuando:
	// - CASO I: se toco un envite+ (con autor del equipo contario)
	// - CASO II: se grito el truco+ (con autor del equipo contario)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= pdt.ENVIDO && p.Ronda.Envite.Estado <= pdt.FALTAENVIDO)
	// ojo: solo a la contraflor+ se le puede decir quiero; a la flor sola no
	laContraFlorEsRespondible := p.Ronda.Envite.Estado >= pdt.CONTRAFLOR && p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	elTrucoEsRespondible := pdt.Contains([]pdt.EstadoTruco{pdt.TRUCO, pdt.RETRUCO, pdt.VALE4}, p.Ronda.Truco.Estado) && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo

	ok := elEnvidoEsRespondible || laContraFlorEsRespondible || elTrucoEsRespondible
	if !ok {
		// si no, esta respondiendo al pedo

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, `No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo`),
		))

		return

	}

	if elEnvidoEsRespondible {

		esDelEquipoContrario := jugada.Manojo.Jugador.Equipo != p.Ronda.Envite.CantadoPor.Jugador.Equipo
		if !esDelEquipoContrario {

			enco.Write(p.out, enco.Pkt(
				enco.Dest(jugada.Manojo.Jugador.Nombre),
				enco.Msg(enco.Error, `La jugada no es valida`),
			))

			return

		}

		enco.Write(p.out, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.QuieroEnvite, jugada.Manojo.Jugador.ID),
		))

		if p.Ronda.Envite.Estado == pdt.FALTAENVIDO {
			tocarFaltaEnvido{jugada.Manojo}.eval(p)
			return
		}
		// si no, era envido/real-envido o cualquier
		// combinacion valida de ellos
		tocarEnvido{jugada.Manojo}.eval(p)
		return

	} else if laContraFlorEsRespondible {
		// tengo que verificar si efectivamente tiene flor
		tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
		esDelEquipoContrario := jugada.Manojo.Jugador.Equipo != p.Ronda.Envite.CantadoPor.Jugador.Equipo
		ok := tieneFlor && esDelEquipoContrario

		if !ok {

			enco.Write(p.out, enco.Pkt(
				enco.Dest(jugada.Manojo.Jugador.Nombre),
				enco.Msg(enco.Error, `La jugada no es valida`),
			))

			return

		}

		enco.Write(p.out, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.QuieroEnvite, jugada.Manojo.Jugador.Nombre),
		))

		// empieza cantando el autor del envite no el que "quizo"
		manojoConLaFlorMasAlta, _ := p.Ronda.GetLaFlorMasAlta()
		equipoGanador := manojoConLaFlorMasAlta.Jugador.Equipo

		if p.Ronda.Envite.Estado == pdt.CONTRAFLOR {
			puntosASumar := p.Ronda.Envite.Puntaje
			p.SumarPuntos(equipoGanador, puntosASumar)

			enco.Write(p.out, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					manojoConLaFlorMasAlta.Jugador.ID,
					enco.ContraFlorGanada,
					puntosASumar),
			))

		} else {
			// el equipo del ganador de la contraflor al resto
			// gano la partida
			// duda se cuentan las flores?
			// puntosASumar := p.Ronda.Envite.Puntaje + p.CalcPtsContraFlorAlResto(equipoGanador)
			puntosASumar := p.CalcPtsContraFlorAlResto(equipoGanador)
			p.SumarPuntos(equipoGanador, puntosASumar)

			enco.Write(p.out, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					manojoConLaFlorMasAlta.Jugador.ID,
					enco.ContraFlorAlRestoGanada,
					puntosASumar),
			))

		}

		p.Ronda.Envite.Estado = pdt.DESHABILITADO

	} else if elTrucoEsRespondible {

		enco.Write(p.out, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.QuieroTruco, jugada.Manojo.Jugador.ID),
		))

		p.PartidaDT.QuererTruco(jugada.Manojo)
	}

	return

}

type responderNoQuiero struct {
	*pdt.Manojo
}

func (jugada responderNoQuiero) hacer(p *Partida) {

	seFueAlMazo := jugada.Manojo.SeFueAlMazo
	if seFueAlMazo {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "Te fuiste al mazo; no podes hacer esta jugada"),
		))

		return
	}

	// checkeo flor en juego
	// caso particular del checkeo: no se le puede decir quiero a la flor
	// pero si a la contra flor o contra flor al resto
	// FALSO porque el no quiero lo estoy contando como un "con flor me achico"
	// todo: agregar la jugada: "con flor me achico" y editar la variale:
	// AHORA:
	// laFlorEsRespondible := p.Ronda.Flor >= pdt.FLOR && p.Ronda.Envite.CantadoPor.Jugador.equipo != jugada.Manojo.Jugador.Equipo
	// LUEGO DE AGREGAR LA JUGADA "con flor me achico"
	// laFlorEsRespondible := p.Ronda.Flor > pdt.FLOR
	// FALSO ---> directamente se va la posibilidad de reponderle
	// "no quiero a la flor"

	// se acepta una respuesta 'no quiero' solo cuando:
	// - CASO I: se toco el envido (o similar)
	// - CASO II: se grito el truco (o similar)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= pdt.ENVIDO && p.Ronda.Envite.Estado <= pdt.FALTAENVIDO) && p.Ronda.Envite.CantadoPor != jugada.Manojo
	laFlorEsRespondible := p.Ronda.Envite.Estado >= pdt.FLOR && p.Ronda.Envite.CantadoPor != jugada.Manojo
	elTrucoEsRespondible := pdt.Contains([]pdt.EstadoTruco{pdt.TRUCO, pdt.RETRUCO, pdt.VALE4}, p.Ronda.Truco.Estado) && p.Ronda.Truco.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo

	ok := elEnvidoEsRespondible || laFlorEsRespondible || elTrucoEsRespondible

	if !ok {
		// si no, esta respondiendo al pedo

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, fmt.Sprintf(`%s esta respondiendo al pedo; no hay nada respondible`, jugada.Manojo.Jugador.Nombre)),
		))

		return

	}

	if elEnvidoEsRespondible {

		esDelEquipoContrario := jugada.Manojo.Jugador.Equipo != p.Ronda.Envite.CantadoPor.Jugador.Equipo
		if !esDelEquipoContrario {

			enco.Write(p.out, enco.Pkt(
				enco.Dest(jugada.Manojo.Jugador.Nombre),
				enco.Msg(enco.Error, `La jugada no es valida`),
			))

			return

		}

		enco.Write(p.out, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.NoQuiero, jugada.Manojo.Jugador.Nombre),
		))

		//	no se toma en cuenta el puntaje total del ultimo toque

		var totalPts int

		switch p.Ronda.Envite.Estado {
		case pdt.ENVIDO:
			totalPts = p.Ronda.Envite.Puntaje - 1
		case pdt.REALENVIDO:
			totalPts = p.Ronda.Envite.Puntaje - 2
		case pdt.FALTAENVIDO:
			totalPts = p.Ronda.Envite.Puntaje + 1
		}

		p.Ronda.Envite.Estado = pdt.DESHABILITADO
		p.Ronda.Envite.Puntaje = totalPts

		enco.Write(p.out, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.SumaPts,
				p.Ronda.Envite.CantadoPor.Jugador.ID,
				enco.EnviteNoQuerido,
				totalPts),
		))

		p.SumarPuntos(p.Ronda.Envite.CantadoPor.Jugador.Equipo, totalPts)

	} else if laFlorEsRespondible {

		// tengo que verificar si efectivamente tiene flor
		tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
		esDelEquipoContrario := jugada.Manojo.Jugador.Equipo != p.Ronda.Envite.CantadoPor.Jugador.Equipo
		ok := tieneFlor && esDelEquipoContrario

		if !ok {

			enco.Write(p.out, enco.Pkt(
				enco.Dest(jugada.Manojo.Jugador.Nombre),
				enco.Msg(enco.Error, `La jugada no es valida`),
			))

			return

		}

		// todo ok: tiene flor; se pasa a jugar:
		enco.Write(p.out, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.ConFlorMeAchico, jugada.Manojo.Jugador.ID),
		))

		// cuenta como un "no quiero" (codigo copiado)
		// segun el estado de la apuesta actual:
		// los "me achico" no cuentan para la flor
		// Flor		xcg(+3) / xcg(+3)
		// Flor + Contra-Flor		xc(+3) / xCadaFlorDelQueHizoElDesafio(+3) + 1
		// Flor + [Contra-Flor] + ContraFlorAlResto		~Falta Envido + *TODAS* las flores no achicadas / xcg(+3) + 1

		// sumo todas las flores del equipo contrario
		totalPts := 0

		for _, m := range p.Ronda.Manojos {
			esDelEquipoContrario := p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
			tieneFlor, _ := m.TieneFlor(p.Ronda.Muestra)
			if tieneFlor && esDelEquipoContrario {
				totalPts += 3
			}
		}

		if p.Ronda.Envite.Estado == pdt.CONTRAFLOR || p.Ronda.Envite.Estado == pdt.CONTRAFLORALRESTO {
			// si es contraflor o al pdt.resto
			// se suma 1 por el `no quiero`
			totalPts++
		}

		p.Ronda.Envite.Estado = pdt.DESHABILITADO

		enco.Write(p.out, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.SumaPts,
				p.Ronda.Envite.CantadoPor.Jugador.ID,
				enco.FlorAchicada,
				totalPts),
		))

		p.SumarPuntos(p.Ronda.Envite.CantadoPor.Jugador.Equipo, totalPts)

	} else if elTrucoEsRespondible {

		enco.Write(p.out, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.NoQuiero, jugada.Manojo.Jugador.Nombre),
		))

		// pongo al equipo que propuso el truco como ganador de la mano actual
		manoActual := p.Ronda.ManoEnJuego.ToInt() - 1
		p.Ronda.Manos[manoActual].Ganador = p.Ronda.Truco.CantadoPor
		equipoGanador := pdt.GanoAzul
		if p.Ronda.Truco.CantadoPor.Jugador.Equipo == pdt.Rojo {
			equipoGanador = pdt.GanoRojo
		}
		p.Ronda.Manos[manoActual].Resultado = equipoGanador

		NuevaRonda, pkts := p.EvaluarRonda()

		if pkts != nil {
			for _, pkt := range pkts {
				if pkt != nil {
					enco.Write(p.out, pkt)
				}
			}
		}

		if NuevaRonda {

			if !p.Terminada() {
				// ahora se deberia de incrementar el mano
				// y ser el turno de este
				sigMano := p.Ronda.GetSigElMano()
				p.NuevaRonda(sigMano) // todo: el tema es que cuando llama aca
				// no manda mensaje de que arranco nueva ronda
				// falso: el padre que llama a .EvaluarRonda tiene que fijarse si
				// retorno true
				// entonces debe crearla el
				// no es responsabilidad de EvaluarRonda arrancar una ronda nueva!!
				// de hecho, si una ronda es terminable y se llama 2 veces consecutivas
				// al mismo metodo booleano, en ambas oportunidades retorna diferente
				// ridiculo
				for _, m := range p.Ronda.Manojos {

					enco.Write(p.out, enco.Pkt(
						enco.Dest(m.Jugador.ID),
						enco.Msg(enco.NuevaRonda, p.PartidaDT.PerspectivaCacheFlor(&m)),
					))
				}

			} else {
				p.byeBye()
			}

		}

	}

	return
}

type irseAlMazo struct {
	*pdt.Manojo
}

func (jugada irseAlMazo) hacer(p *Partida) {
	// checkeos:
	yaSeFueAlMazo := jugada.Manojo.SeFueAlMazo == true
	yaTiroTodasSusCartas := jugada.Manojo.GetCantCartasTiradas() == 3
	if yaSeFueAlMazo || yaTiroTodasSusCartas {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible irse al mazo ahora"),
		))

		return

	}

	seEstabaJugandoElEnvido := (p.Ronda.Envite.Estado >= pdt.ENVIDO && p.Ronda.Envite.Estado <= pdt.FALTAENVIDO)
	seEstabaJugandoLaFlor := p.Ronda.Envite.Estado >= pdt.FLOR
	seEstabaJugandoElTruco := pdt.Contains([]pdt.EstadoTruco{pdt.TRUCO, pdt.RETRUCO, pdt.VALE4}, p.Ronda.Truco.Estado)
	// no se puede ir al mazo sii:
	// 1. el fue el que canto el envido (y el envido esta en juego)
	// 2. tampoco se puede ir al mazo si el canto la flor o similar
	// 3. tampoco se puede ir al mazo si el grito el truco

	// envidoPropuesto := pdt.Contains([]EstadoEnvite{pdt.ENVIDO, pdt.REALENVIDO, pdt.FALTAENVIDO}, p.Ronda.Envite.Estado)
	// envidoPropuestoPorSuEquipo := p.Ronda.Envite.CantadoPor.Jugador.Equipo == jugada.Manojo.Jugador.Equipo
	// trucoPropuesto := pdt.Contains([]pdt.EstadoTruco{pdt.TRUCO, pdt.RETRUCO, pdt.VALE4}, p.Ronda.Truco.Estado)
	// trucoPropuestoPorSuEquipo := p.Ronda.Truco.CantadoPor.Jugador.Equipo == jugada.Manojo.Jugador.Equipo
	// condicionDelBobo := (envidoPropuesto && envidoPropuestoPorSuEquipo) || (trucoPropuesto && trucoPropuestoPorSuEquipo)

	// if condicionDelBobo {

	// enco.Write(p.Stdout, enco.Pkt(
	// 	enco.Dest(jugada.Manojo.Jugador.Nombre),
	// 	enco.Msg(enco.Error,  fmt.Sprintf("No es posible irse al mazo ahora porque hay propuestas de tu equipo sin responder")),
	// ))

	// return

	// }

	noSePuedeIrPorElEnvite := (seEstabaJugandoElEnvido || seEstabaJugandoLaFlor) && p.Ronda.Envite.CantadoPor == jugada.Manojo
	// la de la flor es igual al del envido; porque es un envite
	noSePuedeIrPorElTruco := seEstabaJugandoElTruco && p.Ronda.Truco.CantadoPor == jugada.Manojo
	if noSePuedeIrPorElEnvite || noSePuedeIrPorElTruco {

		enco.Write(p.out, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.Nombre),
			enco.Msg(enco.Error, "No es posible irse al mazo ahora"),
		))

		return

	}

	// ok -> se va al mazo:
	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.Mazo, jugada.Manojo.Jugador.ID),
	))

	p.PartidaDT.IrAlMazo(jugada.Manojo)

	equipoDelJugador := jugada.Manojo.Jugador.Equipo

	seFueronTodos := p.Ronda.CantJugadoresEnJuego[equipoDelJugador] == 0

	// si tenia flor -> ya no lo tomo en cuenta
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	if tieneFlor {
		p.Ronda.Envite.JugadoresConFlor = pdt.Eliminar(p.Ronda.Envite.JugadoresConFlor, jugada.Manojo)
		p.Ronda.Envite.JugadoresConFlorQueNoCantaron = pdt.Eliminar(p.Ronda.Envite.JugadoresConFlorQueNoCantaron, jugada.Manojo)
		// que pasa si era el ultimo que se esperaba que cantara flor?
		// tengo que hacer el eval de la flor
		todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.JugadoresConFlorQueNoCantaron) == 0
		if todosLosJugadoresConFlorCantaron {
			evalFlor(p)
		}
	}

	// era el ultimo en tirar de esta mano?
	eraElUltimoEnTirar := p.Ronda.GetSigHabilitado(*jugada.Manojo) == nil

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
			case pdt.ENVIDO:
				totalPts = e.Puntaje - 1
			case pdt.REALENVIDO:
				totalPts = e.Puntaje - 2
			case pdt.FALTAENVIDO:
				totalPts = e.Puntaje + 1
			}
			e.Estado = pdt.DESHABILITADO
			e.Puntaje = totalPts

			enco.Write(p.out, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					e.CantadoPor.Jugador.ID,
					enco.EnviteNoQuerido,
					totalPts),
			))

			p.SumarPuntos(p.Ronda.Envite.CantadoPor.Jugador.Equipo, totalPts)

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
				esDelEquipoContrario := p.Ronda.Envite.CantadoPor.Jugador.Equipo != jugada.Manojo.Jugador.Equipo
				tieneFlor, _ := m.TieneFlor(p.Ronda.Muestra)
				if tieneFlor && esDelEquipoContrario {
					totalPts += 3
				}
			}

			if p.Ronda.Envite.Estado == pdt.CONTRAFLOR || p.Ronda.Envite.Estado == pdt.CONTRAFLORALRESTO {
				// si es contraflor o al pdt.resto
				// se suma 1 por el `no quiero`
				totalPts++
			}

			p.Ronda.Envite.Estado = pdt.DESHABILITADO

			enco.Write(p.out, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					p.Ronda.Envite.CantadoPor.Jugador.ID,
					enco.FlorAchicada,
					totalPts),
			))

			p.SumarPuntos(p.Ronda.Envite.CantadoPor.Jugador.Equipo, totalPts)

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
		empiezaNuevaRonda, pkts := p.EvaluarMano()

		if pkts != nil {
			for _, pkt := range pkts {
				if pkt != nil {

					// antes:
					//write(p.Stdout, msg)
					// ahora
					enco.Write(p.out, pkt)

				}
			}
		}

		if !empiezaNuevaRonda {

			// actualizo el mano
			p.Ronda.ManoEnJuego++
			p.Ronda.SetNextTurnoPosMano()
			// lo envio
			enco.Write(p.out, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SigTurnoPosMano, int(p.Ronda.Turno)),
			))

		} else {

			if !p.Terminada() {
				// ahora se deberia de incrementar el mano
				// y ser el turno de este
				sigMano := p.Ronda.GetSigElMano()
				p.NuevaRonda(sigMano) // todo: el tema es que cuando llama aca
				// no manda mensaje de que arranco nueva ronda
				// falso: el padre que llama a .EvaluarRonda tiene que fijarse si
				// retorno true
				// entonces debe crearla el
				// no es responsabilidad de EvaluarRonda arrancar una ronda nueva!!
				// de hecho, si una ronda es terminable y se llama 2 veces consecutivas
				// al mismo metodo booleano, en ambas oportunidades retorna diferente
				// ridiculo

				for _, m := range p.Ronda.Manojos {

					enco.Write(p.out, enco.Pkt(
						enco.Dest(m.Jugador.ID),
						enco.Msg(enco.NuevaRonda, p.PartidaDT.PerspectivaCacheFlor(&m)),
					))
				}

			} else {
				p.byeBye()
			}

		}
	} else {
		// cambio de turno solo si era su turno
		eraSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
		if eraSuTurno {
			p.Ronda.SetNextTurno()

			enco.Write(p.out, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SigTurno, int(p.Ronda.Turno)),
			))

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
