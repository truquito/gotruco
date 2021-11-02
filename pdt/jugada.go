package pdt

import (
	"fmt"

	"github.com/filevich/truco/enco"
)

// IJugada Interface para las jugadas
type IJugada interface {
	Ok(p *Partida) ([]*enco.Packet, bool)
	Hacer(p *Partida) []*enco.Packet
}

type TirarCarta struct {
	Manojo *Manojo
	Carta  Carta
}

// Retorna true si la jugada es valida
func (jugada TirarCarta) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// checkeo si se fue al mazo
	noSeFueAlMazo := !jugada.Manojo.SeFueAlMazo
	ok := noSeFueAlMazo
	if !ok {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible tirar una carta porque ya te fuiste al mazo"),
		))

		return pkts, false

	}

	// esto es un tanto redundante porque es imposible que no sea su turno
	// (checkeado mas adelante) y que al mismo tiempo tenga algo para tirar
	// luego de haber jugado sus 3 cartas; aun asi lo dejo
	yaTiroTodasSusCartas := jugada.Manojo.GetCantCartasTiradas() == 3
	if yaTiroTodasSusCartas {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible tirar una carta porque ya las tiraste todas"),
		))

		return pkts, false

	}

	// checkeo flor en juego
	enviteEnJuego := p.Ronda.Envite.Estado >= ENVIDO
	if enviteEnJuego {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible tirar una carta ahora porque el envite esta en juego"),
		))

		return pkts, false

	}

	// primero que nada: tiene esa carta?
	idx, err := jugada.Manojo.GetCartaIdx(jugada.Carta)
	if err != nil {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, err.Error()),
		))

		return pkts, false
	}

	// ya jugo esa carta?
	todaviaNoLaTiro := jugada.Manojo.CartasNoTiradas[idx]
	if !todaviaNoLaTiro {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "Ya tiraste esa carta"),
		))

		return pkts, false
	}

	// luego, era su turno?
	eraSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
	if !eraSuTurno {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No era su turno, no puede tirar la carta"),
		))

		return pkts, false

	}

	// checkeo si tiene flor
	florHabilitada := (p.Ronda.Envite.Estado >= NOCANTADOAUN && p.Ronda.Envite.Estado <= FLOR) && p.Ronda.ManoEnJuego == Primera
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := p.Ronda.Envite.noCantoFlorAun(jugada.Manojo.Jugador.ID)
	noPuedeTirar := florHabilitada && tieneFlor && noCantoFlorAun
	if noPuedeTirar {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible tirar una carta sin antes cantar la flor"),
		))

		return pkts, false

	}

	// cambio: ahora no puede tirar carta si el grito truco
	trucoGritado := Contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado)
	unoDelEquipoContrarioGritoTruco := trucoGritado && p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	yoGiteElTruco := trucoGritado && jugada.Manojo.Jugador.ID == p.Ronda.Truco.CantadoPor
	elTrucoEsRespondible := trucoGritado && unoDelEquipoContrarioGritoTruco && !yoGiteElTruco
	if elTrucoEsRespondible {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible tirar una carta porque tu equipo debe responder la propuesta del truco"),
		))

		return pkts, false

	}

	return pkts, true

}

func (jugada TirarCarta) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	// ok la tiene y era su turno -> la juega
	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TirarCarta,
			jugada.Manojo.Jugador.ID, int(jugada.Carta.Palo), jugada.Carta.Valor),
	))

	idx, _ := jugada.Manojo.GetCartaIdx(jugada.Carta)

	p.TirarCarta(jugada.Manojo, idx)

	// era el ultimo en tirar de esta mano?
	eraElUltimoEnTirar := p.Ronda.GetSigHabilitado(*jugada.Manojo) == nil
	if eraElUltimoEnTirar {
		// de ser asi tengo que checkear el resultado de la mano
		empiezaNuevaRonda, res := p.EvaluarMano()

		pkts = append(pkts, res...)

		if !empiezaNuevaRonda {

			// actualizo el mano
			p.Ronda.ManoEnJuego++
			p.Ronda.SetNextTurnoPosMano()
			// lo envio
			pkts = append(pkts, enco.Pkt(
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
					pkts = append(pkts, enco.Pkt(
						enco.Dest(m.Jugador.ID),
						enco.Msg(enco.NuevaRonda, p.PerspectivaCacheFlor(&m)),
					))

				}

			} // else {
			// p.byeBye()
			// }

		}

		// el turno del siguiente queda dado por el ganador de esta
	} else {
		p.Ronda.SetNextTurno()

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.SigTurno, int(p.Ronda.Turno)),
		))

	}

	return pkts
}

// PRE: supongo que el jugador que toca este envido
// no tiene flor (es checkeada cuando es su turno)
type TocarEnvido struct {
	Manojo *Manojo
}

func (jugada TocarEnvido) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	if florEnJuego {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible tocar el envido ahora porque la flor esta en juego"),
		))

		return pkts, false

	}
	seFueAlMazo := jugada.Manojo.SeFueAlMazo
	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	esSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	esDelEquipoContrario := p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	envidoHabilitado := (p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Envite.Estado == ENVIDO)
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO
	// apuestaSaturada := p.Ronda.Envite.Puntaje >= p.CalcPtsFalta()
	apuestaSaturada := p.Ronda.Envite.Puntaje >= 4
	trucoNoCantado := p.Ronda.Truco.Estado == NOCANTADO

	estaIniciandoPorPrimeraVezElEnvido := esSuTurno && p.Ronda.Envite.Estado == NOCANTADOAUN && trucoNoCantado
	estaRedoblandoLaApuesta := p.Ronda.Envite.Estado == ENVIDO && esDelEquipoContrario // cuando redobla una apuesta puede o no ser su turno
	elEnvidoEstaPrimero := !esSuTurno && p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	puedeTocarEnvido := estaIniciandoPorPrimeraVezElEnvido || estaRedoblandoLaApuesta || elEnvidoEstaPrimero

	ok := !seFueAlMazo && (envidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario) && puedeTocarEnvido && !apuestaSaturada

	if !ok {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, `No es posible cantar 'Envido'`),
		))

		return pkts, false

	}

	return pkts, true
}

func (jugada TocarEnvido) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO
	elEnvidoEstaPrimero := p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	if elEnvidoEstaPrimero {
		// deshabilito el truco
		p.Ronda.Truco.Estado = NOCANTADO
		p.Ronda.Truco.CantadoPor = ""

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.ElEnvidoEstaPrimero, jugada.Manojo.Jugador.ID),
		))

	}

	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TocarEnvido, jugada.Manojo.Jugador.ID),
	))

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.SinCantar) > 0
	if hayFlor {
		// todo: deberia ir al estado magico en el que espera
		// solo por jugadas de tipo flor-related
		// lo mismo para el real-envido; falta-envido
		jid := p.Ronda.Envite.SinCantar[0]
		j := p.Ronda.Manojo[jid]
		siguienteJugada := CantarFlor{j}
		res := siguienteJugada.Hacer(p)
		pkts = append(pkts, res...)

	} else {
		p.TocarEnvido(jugada.Manojo)
	}

	return pkts
}

/* el problema de esta funcion es que esta mas relacionada con el `quiero`
que con el envido. Deberia formar parte del eval del quiero */

// donde 'j' el jugador que dijo 'quiero' al 'envido'/'real envido'
func (jugada TocarEnvido) Eval(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)

	p.Ronda.Envite.Estado = DESHABILITADO
	jIdx, _, res := p.Ronda.ExecElEnvido()

	pkts = append(pkts, res...)

	jug := p.Ronda.Manojos[jIdx].Jugador

	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.SumaPts, jug.ID, enco.EnvidoGanado, p.Ronda.Envite.Puntaje),
	))

	p.SumarPuntos(jug.Equipo, p.Ronda.Envite.Puntaje)

	return pkts
}

type TocarRealEnvido struct {
	Manojo *Manojo
}

func (jugada TocarRealEnvido) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	if florEnJuego {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible tocar real envido ahora porque la flor esta en juego"),
		))

		return pkts, false

	}
	seFueAlMazo := jugada.Manojo.SeFueAlMazo
	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	esSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	realEnvidoHabilitado := (p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Envite.Estado == ENVIDO)
	esDelEquipoContrario := p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO
	trucoNoCantado := p.Ronda.Truco.Estado == NOCANTADO

	estaIniciandoPorPrimeraVezElEnvido := esSuTurno && p.Ronda.Envite.Estado == NOCANTADOAUN && trucoNoCantado
	estaRedoblandoLaApuesta := p.Ronda.Envite.Estado == ENVIDO && esDelEquipoContrario // cuando redobla una apuesta puede o no ser su turno
	elEnvidoEstaPrimero := !esSuTurno && p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	puedeTocarRealEnvido := estaIniciandoPorPrimeraVezElEnvido || estaRedoblandoLaApuesta || elEnvidoEstaPrimero
	ok := !seFueAlMazo && (realEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario) && puedeTocarRealEnvido

	if !ok {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, `No es posible cantar 'Real Envido'`),
		))

		return pkts, false

	}

	return pkts, true
}

func (jugada TocarRealEnvido) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO
	elEnvidoEstaPrimero := p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	if elEnvidoEstaPrimero {
		// deshabilito el truco
		p.Ronda.Truco.Estado = NOCANTADO
		p.Ronda.Truco.CantadoPor = ""

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.ElEnvidoEstaPrimero, jugada.Manojo.Jugador.ID),
		))

	}

	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TocarRealEnvido, jugada.Manojo.Jugador.ID),
	))

	p.TocarRealEnvido(jugada.Manojo)

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.SinCantar) > 0

	if hayFlor {
		jid := p.Ronda.Envite.SinCantar[0]
		j := p.Ronda.Manojo[jid]
		siguienteJugada := CantarFlor{j}
		res := siguienteJugada.Hacer(p)
		pkts = append(pkts, res...)
	}

	return pkts
}

type TocarFaltaEnvido struct {
	Manojo *Manojo
}

func (jugada TocarFaltaEnvido) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	if florEnJuego {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible tocar falta envido ahora porque la flor esta en juego"),
		))

		return pkts, false

	}
	seFueAlMazo := jugada.Manojo.SeFueAlMazo
	esSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	faltaEnvidoHabilitado := p.Ronda.Envite.Estado >= NOCANTADOAUN && p.Ronda.Envite.Estado < FALTAENVIDO
	esDelEquipoContrario := p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado >= ENVIDO
	trucoNoCantado := p.Ronda.Truco.Estado == NOCANTADO

	estaIniciandoPorPrimeraVezElEnvido := esSuTurno && p.Ronda.Envite.Estado == NOCANTADOAUN && trucoNoCantado
	estaRedoblandoLaApuesta := p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado < FALTAENVIDO && esDelEquipoContrario // cuando redobla una apuesta puede o no ser su turno
	elEnvidoEstaPrimero := !esSuTurno && p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	puedeTocarFaltaEnvido := estaIniciandoPorPrimeraVezElEnvido || estaRedoblandoLaApuesta || elEnvidoEstaPrimero
	ok := !seFueAlMazo && (faltaEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario) && puedeTocarFaltaEnvido

	if !ok {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, `No es posible cantar 'Falta Envido'`),
		))

		return pkts, false

	}

	return pkts, true
}

func (jugada TocarFaltaEnvido) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO || p.Ronda.Envite.Estado == REALENVIDO
	elEnvidoEstaPrimero := p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	if elEnvidoEstaPrimero {
		// deshabilito el truco
		p.Ronda.Truco.Estado = NOCANTADO
		p.Ronda.Truco.CantadoPor = ""

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.ElEnvidoEstaPrimero, jugada.Manojo.Jugador.ID),
		))

	}

	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TocarFaltaEnvido, jugada.Manojo.Jugador.ID),
	))

	p.TocarFaltaEnvido(jugada.Manojo)

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.SinCantar) > 0
	if hayFlor {
		jid := p.Ronda.Envite.SinCantar[0]
		j := p.Ronda.Manojo[jid]
		siguienteJugada := CantarFlor{j}
		res := siguienteJugada.Hacer(p)
		pkts = append(pkts, res...)
	}

	return pkts
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

func (jugada TocarFaltaEnvido) Eval(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)

	p.Ronda.Envite.Estado = DESHABILITADO

	// computar envidos
	jIdx, _, res := p.Ronda.ExecElEnvido()

	pkts = append(pkts, res...)

	// jug es el que gano el (falta) envido
	jug := p.Ronda.Manojos[jIdx].Jugador

	pts := p.CalcPtsFaltaEnvido(jug.Equipo)

	p.Ronda.Envite.Puntaje += pts

	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.SumaPts, jug.ID, enco.FaltaEnvidoGanado, p.Ronda.Envite.Puntaje),
	))

	p.SumarPuntos(jug.Equipo, p.Ronda.Envite.Puntaje)

	return pkts
}

type CantarFlor struct {
	Manojo *Manojo
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
func (jugada CantarFlor) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// manojo dice que puede cantar flor;
	// es esto verdad?
	seFueAlMazo := jugada.Manojo.SeFueAlMazo
	florHabilitada := (p.Ronda.Envite.Estado >= NOCANTADOAUN && p.Ronda.Envite.Estado <= FLOR) && p.Ronda.ManoEnJuego == Primera
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := p.Ronda.Envite.noCantoFlorAun(jugada.Manojo.Jugador.ID)
	ok := !seFueAlMazo && florHabilitada && tieneFlor && noCantoFlorAun

	if !ok {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, `No es posible cantar flor`),
		))

		return pkts, false

	}

	return pkts, true
}

func (jugada CantarFlor) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	// yo canto
	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.CantarFlor, jugada.Manojo.Jugador.ID),
	))

	// corresponde que desactive el truco?
	// si lo desactivo: es medio tedioso para el usuario tener q volver a gritar
	// si no lo desacivo: medio como que se olvida
	// QUEDA CONSISTENTE CON "EL ENVIDO ESTA PRIMERO"!
	p.Ronda.Truco.CantadoPor = ""
	p.Ronda.Truco.Estado = NOCANTADO

	// y me elimino de los que no-cantaron
	p.Ronda.Envite.cantoFlor(jugada.Manojo.Jugador.ID)

	p.CantarFlor(jugada.Manojo)

	// es el ultimo en cantar flor que faltaba?
	// o simplemente es el unico que tiene flor (caso particular)

	todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.SinCantar) == 0
	if todosLosJugadoresConFlorCantaron {

		pkts = append(pkts, evalFlor(p)...)

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
			// los quiero llamar a todos, pero no quiero Hacer llamadas al pedo
			// entonces: llamo al primero sin cantar, y que este llame al proximo
			// y que el proximo llame al siguiente, y asi...
			jid := p.Ronda.Envite.SinCantar[0]
			j := p.Ronda.Manojo[jid]
			siguienteJugada := CantarFlor{j}
			res := siguienteJugada.Hacer(p)
			pkts = append(pkts, res...)
		}

	}

	return pkts
}

func evalFlor(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)

	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.SinCantar) == 0
	ok := todosLosJugadoresConFlorCantaron && florEnJuego
	if !ok {
		return pkts
	}

	// cual es la flor ganadora?
	// empieza cantando el autor del envite no el que "quizo"
	autorIdx := p.Ronda.GetIdx(*p.Ronda.Manojo[p.Ronda.Envite.CantadoPor])
	manojoConLaFlorMasAlta, _, res := p.Ronda.ExecLaFlores(JugadorIdx(autorIdx))

	pkts = append(pkts, res...)

	equipoGanador := manojoConLaFlorMasAlta.Jugador.Equipo

	// que estaba en juego?
	switch p.Ronda.Envite.Estado {
	case FLOR:
		// ahora se quien es el ganador; necesito saber cuantos puntos
		// se le va a sumar a ese equipo:
		// los acumulados del envite hasta ahora
		puntosASumar := p.Ronda.Envite.Puntaje
		p.SumarPuntos(equipoGanador, puntosASumar)
		habiaSolo1JugadorConFlor := len(p.Ronda.Envite.JugadoresConFlor) == 1
		if habiaSolo1JugadorConFlor {

			pkts = append(pkts, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					manojoConLaFlorMasAlta.Jugador.ID,
					enco.LaUnicaFlor, puntosASumar),
			))

		} else {

			pkts = append(pkts, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					manojoConLaFlorMasAlta.Jugador.ID,
					enco.LaFlorMasAlta,
					puntosASumar),
			))

		}
	case CONTRAFLOR:
	case CONTRAFLORALRESTO:
	}

	p.Ronda.Envite.Estado = DESHABILITADO

	return pkts
}

type CantarContraFlor struct {
	Manojo *Manojo
}

func (jugada CantarContraFlor) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// manojo dice que puede cantar flor;
	// es esto verdad?
	seFueAlMazo := jugada.Manojo.SeFueAlMazo
	contraFlorHabilitada := p.Ronda.Envite.Estado == FLOR && p.Ronda.ManoEnJuego == Primera
	esDelEquipoContrario := contraFlorHabilitada && p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := p.Ronda.Envite.noCantoFlorAun(jugada.Manojo.Jugador.ID)
	ok := !seFueAlMazo && contraFlorHabilitada && tieneFlor && esDelEquipoContrario && noCantoFlorAun
	if !ok {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, `No es posible cantar contra flor`),
		))

		return pkts, false

	}

	return pkts, true
}

func (jugada CantarContraFlor) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	// la canta
	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.CantarContraFlor, jugada.Manojo.Jugador.ID),
	))

	p.CantarContraFlor(jugada.Manojo)
	// y ahora tengo que esperar por la respuesta de la nueva
	// propuesta de todos menos de el que canto la contraflor
	// restauro la copia
	p.Ronda.Envite.cantoFlor(jugada.Manojo.Jugador.ID)

	return pkts
}

type CantarContraFlorAlResto struct {
	Manojo *Manojo
}

func (jugada CantarContraFlorAlResto) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// manojo dice que puede cantar flor;
	// es esto verdad?
	seFueAlMazo := jugada.Manojo.SeFueAlMazo
	contraFlorHabilitada := (p.Ronda.Envite.Estado == FLOR || p.Ronda.Envite.Estado == CONTRAFLOR) && p.Ronda.ManoEnJuego == Primera
	esDelEquipoContrario := contraFlorHabilitada && p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := p.Ronda.Envite.noCantoFlorAun(jugada.Manojo.Jugador.ID)
	ok := !seFueAlMazo && contraFlorHabilitada && tieneFlor && esDelEquipoContrario && noCantoFlorAun
	if !ok {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, `No es posible cantar contra flor al resto`),
		))

		return pkts, false

	}

	return pkts, true
}

func (jugada CantarContraFlorAlResto) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	// la canta
	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.CantarContraFlorAlResto, jugada.Manojo.Jugador.ID),
	))

	p.CantarContraFlorAlResto(jugada.Manojo)
	// y ahora tengo que esperar por la respuesta de la nueva
	// propuesta de todos menos de el que canto la contraflor
	// restauro la copia
	p.Ronda.Envite.cantoFlor(jugada.Manojo.Jugador.ID)

	return pkts
}

type CantarConFlorMeAchico struct {
	Manojo *Manojo
}

// no implementada (porque no es necesaria?)
func (jugada CantarConFlorMeAchico) Hacer(p *Partida) []*enco.Packet {
	return nil
}

type GritarTruco struct {
	Manojo *Manojo
}

// se resuelve lo siguiente:
// ------------------------
// si yo tengo flor, o uno de mis companeros tienen flor
// (y no se ha cantado aun)
// entonces no puedo cantar truco

// si yo puedo gritar truco ->
// cambia el estado del Truco a TRUCO
// luego si alguien dice:
// 		quiero -> no debe poder si uno de su equipo tiene flor
// 		si dice flor -> debe resetear el Truco

func (jugada GritarTruco) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// checkeos:
	noSeFueAlMazo := !jugada.Manojo.SeFueAlMazo
	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= NOCANTADOAUN

	yoOUnoDeMisCompasTieneFlorYAunNoCanto := p.Ronda.hayEquipoSinCantar(jugada.Manojo.Jugador.Equipo)

	laFlorEstaPrimero := yoOUnoDeMisCompasTieneFlorYAunNoCanto
	trucoNoSeJugoAun := p.Ronda.Truco.Estado == NOCANTADO
	esSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
	trucoHabilitado := noSeFueAlMazo && trucoNoSeJugoAun && noSeEstaJugandoElEnvite && !laFlorEstaPrimero && esSuTurno

	if !trucoHabilitado {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible cantar truco ahora"),
		))

		return pkts, false

	}

	return pkts, true
}

func (jugada GritarTruco) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.GritarTruco, jugada.Manojo.Jugador.ID),
	))

	p.GritarTruco(jugada.Manojo)

	return pkts
}

type GritarReTruco struct {
	Manojo *Manojo
}

func (jugada GritarReTruco) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// checkeos generales:
	noSeFueAlMazo := !jugada.Manojo.SeFueAlMazo
	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= NOCANTADOAUN

	yoOUnoDeMisCompasTieneFlorYAunNoCanto := p.Ronda.hayEquipoSinCantar(jugada.Manojo.Jugador.Equipo)

	laFlorEstaPrimero := yoOUnoDeMisCompasTieneFlorYAunNoCanto

	/*
		Hay 2 casos para cantar rectruco:
		    - CASO I: Uno del equipo contrario grito el truco
			- CASO II: Uno de su equipo posee el quiero
	*/

	// CASO I:
	trucoGritado := p.Ronda.Truco.Estado == TRUCO
	unoDelEquipoContrarioGritoTruco := trucoGritado && p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	casoI := trucoGritado && unoDelEquipoContrarioGritoTruco

	// CASO II:
	trucoYaQuerido := p.Ronda.Truco.Estado == TRUCOQUERIDO
	unoDeMiEquipoQuizo := trucoYaQuerido && p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo == jugada.Manojo.Jugador.Equipo
	// esTurnoDeMiEquipo := p.Ronda.GetElTurno().Jugador.Equipo == jugada.Manojo.Jugador.Equipo
	casoII := trucoYaQuerido && unoDeMiEquipoQuizo // && esTurnoDeMiEquipo

	reTrucoHabilitado := noSeFueAlMazo && noSeEstaJugandoElEnvite && (casoI || casoII) && !laFlorEstaPrimero

	if !reTrucoHabilitado {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible cantar re-truco ahora"),
		))

		return pkts, false

	}

	return pkts, true
}

// checkeaos de este tipo:
// que pasa cuando gritan re-truco cuando el campo truco se encuentra nil
// ese fue el nil pointer exception
func (jugada GritarReTruco) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.GritarReTruco, jugada.Manojo.Jugador.ID),
	))

	p.GritarReTruco(jugada.Manojo)

	return pkts
}

type GritarVale4 struct {
	Manojo *Manojo
}

func (jugada GritarVale4) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// checkeos:
	noSeFueAlMazo := !jugada.Manojo.SeFueAlMazo

	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= NOCANTADOAUN

	yoOUnoDeMisCompasTieneFlorYAunNoCanto := p.Ronda.hayEquipoSinCantar(jugada.Manojo.Jugador.Equipo)

	laFlorEstaPrimero := yoOUnoDeMisCompasTieneFlorYAunNoCanto

	/*
		Hay 2 casos para cantar rectruco:
		    - CASO I: Uno del equipo contrario grito el re-truco
			- CASO II: Uno de su equipo posee el quiero
	*/

	// CASO I:
	reTrucoGritado := p.Ronda.Truco.Estado == RETRUCO
	// para eviat el nil primero checkeo que haya sido gritado reTrucoGritado &&
	unoDelEquipoContrarioGritoReTruco := reTrucoGritado && p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	casoI := reTrucoGritado && unoDelEquipoContrarioGritoReTruco

	// CASO I:
	retrucoYaQuerido := p.Ronda.Truco.Estado == RETRUCOQUERIDO
	// para eviat el nil primero checkeo que haya sido gritado reTrucoGritado &&
	suEquipotieneElQuiero := retrucoYaQuerido && p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo == jugada.Manojo.Jugador.Equipo
	casoII := retrucoYaQuerido && suEquipotieneElQuiero

	vale4Habilitado := noSeFueAlMazo && (casoI || casoII) && noSeEstaJugandoElEnvite && !laFlorEstaPrimero

	if !vale4Habilitado {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible cantar vale-4 ahora"),
		))

		return pkts, false

	}

	return pkts, true
}

func (jugada GritarVale4) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.GritarVale4, jugada.Manojo.Jugador.ID),
	))

	p.GritarVale4(jugada.Manojo)

	return pkts
}

type ResponderQuiero struct {
	Manojo *Manojo
}

func (jugada ResponderQuiero) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	seFueAlMazo := jugada.Manojo.SeFueAlMazo
	if seFueAlMazo {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "Te fuiste al mazo; no podes Hacer esta jugada"),
		))

		return pkts, false

	}

	// checkeo flor en juego
	// caso particular del checkeo:
	// no se le puede decir quiero ni al envido* ni al truco si se esta jugando la flor
	// no se le puede decir quiero a la flor -> si la flor esta en juego -> error
	// pero si a la contra flor o contra flor al resto

	// casos posibles:
	// alguien dijo envido/truco, otro responde quiero, pero hay uno que tiene flor que todavia no la jugo -> deberia saltar error: "alguien tiene flor y no la jugo aun"
	// alguien tiene flor, uno dice quiero -> no deberia dejarlo porque la flor no se responde con quiero
	// se esta jugando la contra-flor/CFAR -> ok

	florEnJuego := p.Ronda.Envite.Estado == FLOR
	if florEnJuego {
		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible responder quiero ahora"),
		))

		return pkts, false
	}

	noHanCantadoLaFlorAun := p.Ronda.Envite.Estado < FLOR
	yoOUnoDeMisCompasTieneFlorYAunNoCanto := p.Ronda.hayEquipoSinCantar(jugada.Manojo.Jugador.Equipo)
	if noHanCantadoLaFlorAun && yoOUnoDeMisCompasTieneFlorYAunNoCanto {
		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible responder 'quiero' porque alguien con flor no ha cantado aun"),
		))

		return pkts, false
	}
	// se acepta una respuesta 'quiero' solo cuando:
	// - CASO I: se toco un envite+ (con autor del equipo contario)
	// - CASO II: se grito el truco+ (con autor del equipo contario)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO)
	// ojo: solo a la contraflor+ se le puede decir quiero; a la flor sola no
	laContraFlorEsRespondible := p.Ronda.Envite.Estado >= CONTRAFLOR && p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	elTrucoEsRespondible := Contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado) && p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo

	ok := elEnvidoEsRespondible || laContraFlorEsRespondible || elTrucoEsRespondible
	if !ok {
		// si no, esta respondiendo al pedo

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, `No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo`),
		))

		return pkts, false

	}

	if elEnvidoEsRespondible {

		esDelEquipoContrario := jugada.Manojo.Jugador.Equipo != p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo
		if !esDelEquipoContrario {

			pkts = append(pkts, enco.Pkt(
				enco.Dest(jugada.Manojo.Jugador.ID),
				enco.Msg(enco.Error, `La jugada no es valida`),
			))

			return pkts, false

		}

	} else if laContraFlorEsRespondible {
		// tengo que verificar si efectivamente tiene flor
		tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
		esDelEquipoContrario := jugada.Manojo.Jugador.Equipo != p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo
		ok := tieneFlor && esDelEquipoContrario

		if !ok {

			pkts = append(pkts, enco.Pkt(
				enco.Dest(jugada.Manojo.Jugador.ID),
				enco.Msg(enco.Error, `La jugada no es valida`),
			))

			return pkts, false

		}

	}

	return pkts, true
}

func (jugada ResponderQuiero) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	// se acepta una respuesta 'quiero' solo cuando:
	// - CASO I: se toco un envite+ (con autor del equipo contario)
	// - CASO II: se grito el truco+ (con autor del equipo contario)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO)
	// ojo: solo a la contraflor+ se le puede decir quiero; a la flor sola no
	laContraFlorEsRespondible := p.Ronda.Envite.Estado >= CONTRAFLOR && p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
	elTrucoEsRespondible := Contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado) && p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo

	if elEnvidoEsRespondible {

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.QuieroEnvite, jugada.Manojo.Jugador.ID),
		))

		if p.Ronda.Envite.Estado == FALTAENVIDO {
			res := TocarFaltaEnvido(jugada).Eval(p)
			return append(pkts, res...)
		}
		// si no, era envido/real-envido o cualquier
		// combinacion valida de ellos

		res := TocarEnvido(jugada).Eval(p)
		return append(pkts, res...)

	} else if laContraFlorEsRespondible {

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.QuieroEnvite, jugada.Manojo.Jugador.ID),
		))

		// empieza cantando el autor del envite no el que "quizo"
		autorIdx := p.Ronda.GetIdx(*p.Ronda.Manojo[p.Ronda.Envite.CantadoPor])
		manojoConLaFlorMasAlta, _, res := p.Ronda.ExecLaFlores(JugadorIdx(autorIdx))

		pkts = append(pkts, res...)

		// manojoConLaFlorMasAlta, _ := p.Ronda.GetLaFlorMasAlta()
		equipoGanador := manojoConLaFlorMasAlta.Jugador.Equipo

		if p.Ronda.Envite.Estado == CONTRAFLOR {
			puntosASumar := p.Ronda.Envite.Puntaje
			p.SumarPuntos(equipoGanador, puntosASumar)

			pkts = append(pkts, enco.Pkt(
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

			pkts = append(pkts, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					manojoConLaFlorMasAlta.Jugador.ID,
					enco.ContraFlorAlRestoGanada,
					puntosASumar),
			))

		}

		p.Ronda.Envite.Estado = DESHABILITADO

	} else if elTrucoEsRespondible {

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.QuieroTruco, jugada.Manojo.Jugador.ID),
		))

		p.QuererTruco(jugada.Manojo)
	}

	return pkts

}

type ResponderNoQuiero struct {
	Manojo *Manojo
}

func (jugada ResponderNoQuiero) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	seFueAlMazo := jugada.Manojo.SeFueAlMazo
	if seFueAlMazo {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "Te fuiste al mazo; no podes Hacer esta jugada"),
		))

		return pkts, false
	}

	// checkeo flor en juego
	// caso particular del checkeo: no se le puede decir quiero a la flor
	// pero si a la contra flor o contra flor al resto
	// FALSO porque el no quiero lo estoy contando como un "con flor me achico"
	// todo: agregar la jugada: "con flor me achico" y editar la variale:
	// AHORA:
	// laFlorEsRespondible := p.Ronda.Flor >= FLOR && p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.equipo != jugada.Manojo.Jugador.Equipo
	// LUEGO DE AGREGAR LA JUGADA "con flor me achico"
	// laFlorEsRespondible := p.Ronda.Flor > FLOR
	// FALSO ---> directamente se va la posibilidad de reponderle
	// "no quiero a la flor"

	// se acepta una respuesta 'no quiero' solo cuando:
	// - CASO I: se toco el envido (o similar)
	// - CASO II: se grito el truco (o similar)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO) && p.Ronda.Envite.CantadoPor != jugada.Manojo.Jugador.ID
	laFlorEsRespondible := p.Ronda.Envite.Estado >= FLOR && p.Ronda.Envite.CantadoPor != jugada.Manojo.Jugador.ID
	elTrucoEsRespondible := Contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado) && p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo

	ok := elEnvidoEsRespondible || laFlorEsRespondible || elTrucoEsRespondible

	if !ok {
		// si no, esta respondiendo al pedo

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, fmt.Sprintf(`%s esta respondiendo al pedo; no hay nada respondible`, jugada.Manojo.Jugador.ID)),
		))

		return pkts, false
	}

	if elEnvidoEsRespondible {

		esDelEquipoContrario := jugada.Manojo.Jugador.Equipo != p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo
		if !esDelEquipoContrario {

			pkts = append(pkts, enco.Pkt(
				enco.Dest(jugada.Manojo.Jugador.ID),
				enco.Msg(enco.Error, `La jugada no es valida`),
			))

			return pkts, false

		}

	} else if laFlorEsRespondible {

		// tengo que verificar si efectivamente tiene flor
		tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
		esDelEquipoContrario := jugada.Manojo.Jugador.Equipo != p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo
		ok := tieneFlor && esDelEquipoContrario

		if !ok {

			pkts = append(pkts, enco.Pkt(
				enco.Dest(jugada.Manojo.Jugador.ID),
				enco.Msg(enco.Error, `La jugada no es valida`),
			))

			return pkts, false

		}

	}

	return pkts, true
}

func (jugada ResponderNoQuiero) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO) && p.Ronda.Envite.CantadoPor != jugada.Manojo.Jugador.ID
	laFlorEsRespondible := p.Ronda.Envite.Estado >= FLOR && p.Ronda.Envite.CantadoPor != jugada.Manojo.Jugador.ID
	elTrucoEsRespondible := Contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado) && p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo

	if elEnvidoEsRespondible {

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.NoQuiero, jugada.Manojo.Jugador.ID),
		))

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

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.SumaPts,
				p.Ronda.Envite.CantadoPor,
				enco.EnviteNoQuerido,
				totalPts),
		))

		p.SumarPuntos(p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo, totalPts)

	} else if laFlorEsRespondible {

		// todo ok: tiene flor; se pasa a jugar:
		pkts = append(pkts, enco.Pkt(
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
			esDelEquipoContrario := p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
			tieneFlor, _ := m.TieneFlor(p.Ronda.Muestra)
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

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.SumaPts,
				p.Ronda.Envite.CantadoPor,
				enco.FlorAchicada,
				totalPts),
		))

		p.SumarPuntos(p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo, totalPts)

	} else if elTrucoEsRespondible {

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.NoQuiero, jugada.Manojo.Jugador.ID),
		))

		// pongo al equipo que propuso el truco como ganador de la mano actual
		manoActual := p.Ronda.ManoEnJuego.ToInt() - 1
		p.Ronda.Manos[manoActual].Ganador = p.Ronda.Truco.CantadoPor
		equipoGanador := GanoAzul
		if p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo == Rojo {
			equipoGanador = GanoRojo
		}
		p.Ronda.Manos[manoActual].Resultado = equipoGanador

		NuevaRonda, res := p.EvaluarRonda()

		pkts = append(pkts, res...)

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

					pkts = append(pkts, enco.Pkt(
						enco.Dest(m.Jugador.ID),
						enco.Msg(enco.NuevaRonda, p.PerspectivaCacheFlor(&m)),
					))
				}

			} // else {
			// p.byeBye()
			// }

		}

	}

	return pkts
}

type IrseAlMazo struct {
	Manojo *Manojo
}

func (jugada IrseAlMazo) Ok(p *Partida) ([]*enco.Packet, bool) {
	pkts := make([]*enco.Packet, 0)

	// checkeos:
	yaSeFueAlMazo := jugada.Manojo.SeFueAlMazo
	yaTiroTodasSusCartas := jugada.Manojo.GetCantCartasTiradas() == 3
	if yaSeFueAlMazo || yaTiroTodasSusCartas {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible irse al mazo ahora"),
		))

		return pkts, false

	}

	seEstabaJugandoElEnvido := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO)
	seEstabaJugandoLaFlor := p.Ronda.Envite.Estado >= FLOR
	seEstabaJugandoElTruco := Contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado)
	// no se puede ir al mazo sii:
	// 1. el fue el que canto el envido (y el envido esta en juego)
	// 2. tampoco se puede ir al mazo si el canto la flor o similar
	// 3. tampoco se puede ir al mazo si el grito el truco

	// envidoPropuesto := Contains([]EstadoEnvite{ENVIDO, REALENVIDO, FALTAENVIDO}, p.Ronda.Envite.Estado)
	// envidoPropuestoPorSuEquipo := p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo == jugada.Manojo.Jugador.Equipo
	// trucoPropuesto := Contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado)
	// trucoPropuestoPorSuEquipo := p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo == jugada.Manojo.Jugador.Equipo
	// condicionDelBobo := (envidoPropuesto && envidoPropuestoPorSuEquipo) || (trucoPropuesto && trucoPropuestoPorSuEquipo)

	// if condicionDelBobo {

	// enco.Write(p.Stdout, enco.Pkt(
	// 	enco.Dest(jugada.Manojo.Jugador.ID),
	// 	enco.Msg(enco.Error,  fmt.Sprintf("No es posible irse al mazo ahora porque hay propuestas de tu equipo sin responder")),
	// ))

	// return

	// }

	noSePuedeIrPorElEnvite := (seEstabaJugandoElEnvido || seEstabaJugandoLaFlor) && p.Ronda.Envite.CantadoPor == jugada.Manojo.Jugador.ID
	// la de la flor es igual al del envido; porque es un envite
	noSePuedeIrPorElTruco := seEstabaJugandoElTruco && p.Ronda.Truco.CantadoPor == jugada.Manojo.Jugador.ID
	if noSePuedeIrPorElEnvite || noSePuedeIrPorElTruco {

		pkts = append(pkts, enco.Pkt(
			enco.Dest(jugada.Manojo.Jugador.ID),
			enco.Msg(enco.Error, "No es posible irse al mazo ahora"),
		))

		return pkts, false

	}

	return pkts, true
}

func (jugada IrseAlMazo) Hacer(p *Partida) []*enco.Packet {

	pkts := make([]*enco.Packet, 0)
	pre, ok := jugada.Ok(p)
	pkts = append(pkts, pre...)

	if !ok {
		return pkts
	}

	// ok -> se va al mazo:
	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.Mazo, jugada.Manojo.Jugador.ID),
	))

	p.IrAlMazo(jugada.Manojo)

	equipoDelJugador := jugada.Manojo.Jugador.Equipo

	seFueronTodos := p.Ronda.CantJugadoresEnJuego[equipoDelJugador] == 0

	// si tenia flor -> ya no lo tomo en cuenta
	tieneFlor, _ := jugada.Manojo.TieneFlor(p.Ronda.Muestra)
	if tieneFlor {
		p.Ronda.Envite.JugadoresConFlor = Eliminar(p.Ronda.Envite.JugadoresConFlor, jugada.Manojo)
		p.Ronda.Envite.cantoFlor(jugada.Manojo.Jugador.ID)
		// que pasa si era el ultimo que se esperaba que cantara flor?
		// tengo que Hacer el Eval de la flor
		todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.SinCantar) == 0
		if todosLosJugadoresConFlorCantaron {
			pkts = append(pkts, evalFlor(p)...)
		}
	}

	// era el ultimo en tirar de esta mano?
	eraElUltimoEnTirar := p.Ronda.GetSigHabilitado(*jugada.Manojo) == nil

	if seFueronTodos {

		seEstabaJugandoElEnvido := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO)
		seEstabaJugandoLaFlor := p.Ronda.Envite.Estado >= FLOR

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

			pkts = append(pkts, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					e.CantadoPor,
					enco.EnviteNoQuerido,
					totalPts),
			))

			p.SumarPuntos(p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo, totalPts)

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
				esDelEquipoContrario := p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo != jugada.Manojo.Jugador.Equipo
				tieneFlor, _ := m.TieneFlor(p.Ronda.Muestra)
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

			pkts = append(pkts, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SumaPts,
					p.Ronda.Envite.CantadoPor,
					enco.FlorAchicada,
					totalPts),
			))

			p.SumarPuntos(p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo, totalPts)

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
		empiezaNuevaRonda, res := p.EvaluarMano()

		pkts = append(pkts, res...)

		if !empiezaNuevaRonda {

			// actualizo el mano
			p.Ronda.ManoEnJuego++
			p.Ronda.SetNextTurnoPosMano()
			// lo envio
			pkts = append(pkts, enco.Pkt(
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

					pkts = append(pkts, enco.Pkt(
						enco.Dest(m.Jugador.ID),
						enco.Msg(enco.NuevaRonda, p.PerspectivaCacheFlor(&m)),
					))
				}

			} // else {
			// p.byeBye()
			// }

		}
	} else {
		// cambio de turno solo si era su turno
		eraSuTurno := p.Ronda.GetElTurno() == jugada.Manojo
		if eraSuTurno {
			p.Ronda.SetNextTurno()

			pkts = append(pkts, enco.Pkt(
				enco.Dest("ALL"),
				enco.Msg(enco.SigTurno, int(p.Ronda.Turno)),
			))

		}
	}

	return pkts
}
