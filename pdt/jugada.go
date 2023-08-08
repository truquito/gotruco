package pdt

import (
	"fmt"

	"github.com/filevich/truco/enco"
)

type IJUGADA_ID int

const (
	JID_TIRAR_CARTA IJUGADA_ID = iota
	JID_ENVIDO
	JID_REAL_ENVIDO
	JID_FALTA_ENVIDO
	JID_FLOR
	JID_CONTRA_FLOR
	JID_CONTRA_FLOR_AL_RESTO
	JID_TRUCO
	JID_RE_TRUCO
	JID_VALE_4
	JID_QUIERO
	JID_NO_QUIERO
	JID_MAZO
)

// IJugada Interface para las jugadas
type IJugada interface {
	Ok(p *Partida) ([]enco.Envelope, bool)
	Hacer(p *Partida) []enco.Envelope
	String() string
	ID() IJUGADA_ID
}

type TirarCarta struct {
	// Manojo *Manojo
	JID   string
	Carta Carta
}

func (jugada TirarCarta) ID() IJUGADA_ID {
	return JID_TIRAR_CARTA
}

func (jugada TirarCarta) String() string {
	return fmt.Sprintf("%s %d %s",
		jugada.JID,
		jugada.Carta.Valor,
		jugada.Carta.Palo.String(),
	)
}

// Retorna true si la jugada es valida
func (jugada TirarCarta) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// checkeo si se fue al mazo
	noSeFueAlMazo := !p.Manojo(jugada.JID).SeFueAlMazo
	ok := noSeFueAlMazo
	if !ok {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tirar una carta porque ya te fuiste al mazo"),
			))
		}

		return pkts2, false

	}

	// esto es un tanto redundante porque es imposible que no sea su turno
	// (checkeado mas adelante) y que al mismo tiempo tenga algo para tirar
	// luego de haber jugado sus 3 cartas; aun asi lo dejo
	yaTiroTodasSusCartas := p.Manojo(jugada.JID).GetCantCartasTiradas() == 3
	if yaTiroTodasSusCartas {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tirar una carta porque ya las tiraste todas"),
			))
		}

		return pkts2, false

	}

	// checkeo flor en juego
	enviteEnJuego := p.Ronda.Envite.Estado >= ENVIDO
	if enviteEnJuego {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tirar una carta ahora porque el envite esta en juego"),
			))
		}

		return pkts2, false

	}

	// primero que nada: tiene esa carta?
	// pide por un error, pero p.Manojo NO retorna error alguno !
	idx, err := p.Manojo(jugada.JID).GetCartaIdx(jugada.Carta)
	if err != nil {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error(err.Error()),
			))
		}

		return pkts2, false
	}

	// ya jugo esa carta?
	todaviaNoLaTiro := !p.Manojo(jugada.JID).Tiradas[idx]
	if !todaviaNoLaTiro {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("Ya tiraste esa carta"),
			))
		}

		return pkts2, false
	}

	// luego, era su turno?
	eraSuTurno := p.Ronda.GetElTurno().Jugador.ID == p.Manojo(jugada.JID).Jugador.ID
	if !eraSuTurno {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No era su turno, no puede tirar la carta"),
			))
		}

		return pkts2, false

	}

	// checkeo si tiene flor
	florHabilitada := (p.Ronda.Envite.Estado >= NOCANTADOAUN && p.Ronda.Envite.Estado <= FLOR) && p.Ronda.ManoEnJuego == Primera
	tieneFlor, _ := p.Manojo(jugada.JID).TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := p.Ronda.Envite.noCantoFlorAun(p.Manojo(jugada.JID).Jugador.ID)
	noPuedeTirar := florHabilitada && tieneFlor && noCantoFlorAun
	if noPuedeTirar {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tirar una carta sin antes cantar la flor"),
			))
		}

		return pkts2, false

	}

	// cambio: ahora no puede tirar carta si el grito truco
	trucoGritado := p.Ronda.Truco.Estado.esTrucoRespondible()
	unoDelEquipoContrarioGritoTruco := trucoGritado && p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
	yoGiteElTruco := trucoGritado && p.Manojo(jugada.JID).Jugador.ID == p.Ronda.Truco.CantadoPor
	elTrucoEsRespondible := trucoGritado && unoDelEquipoContrarioGritoTruco && !yoGiteElTruco
	if elTrucoEsRespondible {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tirar una carta porque tu equipo debe responder la propuesta del truco"),
			))
		}

		return pkts2, false

	}

	return pkts2, true

}

func (jugada TirarCarta) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.TirarCarta{
				Autor: jugada.JID,
				Palo:  jugada.Carta.Palo.String(),
				Valor: jugada.Carta.Valor,
			},
		))
	}

	idx, _ := p.Manojo(jugada.JID).GetCartaIdx(jugada.Carta)

	p.TirarCarta(p.Manojo(jugada.JID), idx)

	// era el ultimo en tirar de esta mano?
	eraElUltimoEnTirar := p.Ronda.GetSigHabilitado(*p.Manojo(jugada.JID)) == nil
	if eraElUltimoEnTirar {
		// de ser asi tengo que checkear el resultado de la mano
		empiezaNuevaRonda, res2 := p.EvaluarMano()

		if p.Verbose {
			pkts2 = append(pkts2, res2...)
		}

		if !empiezaNuevaRonda {

			seTerminoLaPrimeraMano := p.Ronda.ManoEnJuego == Primera
			nadieCantoEnvite := p.Ronda.Envite.Estado == NOCANTADOAUN
			if seTerminoLaPrimeraMano && nadieCantoEnvite {
				p.Ronda.Envite.Estado = DESHABILITADO
				p.Ronda.Envite.SinCantar = []string{}
			}

			// actualizo el mano
			p.Ronda.ManoEnJuego++
			p.Ronda.SetNextTurnoPosMano()
			// lo envio

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.ALL,
					enco.SigTurnoPosMano(p.Ronda.Turno),
				))
			}

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

				if p.Verbose {
					for _, m := range p.Ronda.Manojos {
						pkts2 = append(pkts2, enco.Pkt(
							enco.Dest(m.Jugador.ID),
							enco.NuevaRonda{
								Perspectiva: p.PerspectivaCacheFlor(&m),
							},
						))
					}
				}

			} // else {
			// p.byeBye()
			// }

		}

		// el turno del siguiente queda dado por el ganador de esta
	} else {
		p.Ronda.SetNextTurno()

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.SigTurno(p.Ronda.Turno),
			))
		}

	}

	return pkts2
}

// PRE: supongo que el jugador que toca este envido
// no tiene flor (es checkeada cuando es su turno)
type TocarEnvido struct {
	// Manojo *Manojo
	JID string
}

func (jugada TocarEnvido) ID() IJUGADA_ID {
	return JID_ENVIDO
}

func (jugada TocarEnvido) String() string {
	return jugada.JID + " envido"
}

func (jugada TocarEnvido) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	if florEnJuego {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tocar el envido ahora porque la flor esta en juego"),
			))
		}

		return pkts2, false
	}
	seFueAlMazo := p.Manojo(jugada.JID).SeFueAlMazo
	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	esSuTurno := p.Ronda.GetElTurno().Jugador.ID == p.Manojo(jugada.JID).Jugador.ID
	tieneFlor, _ := p.Manojo(jugada.JID).TieneFlor(p.Ronda.Muestra)
	envidoHabilitado := (p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Envite.Estado == ENVIDO)

	if !envidoHabilitado {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tocar envido ahora"),
			))
		}

		return pkts2, false
	}

	// supuestamente esto que sigue no es posible.
	// pero en el randomWalker me parecio haberlo visto.
	// lo dejo comentado por las dudas.
	// Tal vez fue error de codificacion del randomWalker.

	// puede cantar envite (desde 0; es decir, empezar un envido) solo si no tirno
	// niguna carta aun;
	// pero si puede responder a un envido incluso si ya tiro

	// yaTiroAlgunaCarta := p.Manojo(jugada.JID).yaTiroCarta(Primera)
	// estaIniciandoElEnvite := p.Ronda.Envite.Estado == NOCANTADOAUN
	// envidoHabilitado = !(yaTiroAlgunaCarta && estaIniciandoElEnvite)
	// if !envidoHabilitado {
	// 	pkts = append(pkts, enco.Pkt1(
	// 		enco.Dest(jugada.JID),
	// 		enco.Msg(enco.Error, "No es posible tocar envido ahora"),
	// 	))
	// 	return pkts2, false
	// }

	esDelEquipoContrario := p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO
	// antes:
	// apuestaSaturada := p.Ronda.Envite.Puntaje >= p.CalcPtsFalta()
	// apuestaSaturada := p.Ronda.Envite.Puntaje >= 4
	// ahora:
	lim := p.LimiteEnvido
	if lim == -1 {
		lim = p.CalcPtsFalta()
	}
	apuestaSaturada := p.Ronda.Envite.Puntaje >= lim
	trucoNoCantado := p.Ronda.Truco.Estado == NOGRITADOAUN

	estaIniciandoPorPrimeraVezElEnvido := esSuTurno && p.Ronda.Envite.Estado == NOCANTADOAUN && trucoNoCantado
	estaRedoblandoLaApuesta := p.Ronda.Envite.Estado == ENVIDO && esDelEquipoContrario // cuando redobla una apuesta puede o no ser su turno
	elEnvidoEstaPrimero := !esSuTurno && p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	puedeTocarEnvido := estaIniciandoPorPrimeraVezElEnvido || estaRedoblandoLaApuesta || elEnvidoEstaPrimero

	ok := !seFueAlMazo && (envidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario) && puedeTocarEnvido && !apuestaSaturada

	if !ok {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error(`No es posible cantar 'Envido'`),
			))
		}

		return pkts2, false

	}

	return pkts2, true
}

func (jugada TocarEnvido) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO
	elEnvidoEstaPrimero := p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	if elEnvidoEstaPrimero {
		// deshabilito el truco
		p.Ronda.Truco.Estado = NOGRITADOAUN
		p.Ronda.Truco.CantadoPor = ""

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.ElEnvidoEstaPrimero(jugada.JID),
			))
		}

	}

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.TocarEnvido(jugada.JID),
		))
	}

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.SinCantar) > 0
	if hayFlor {
		// todo: deberia ir al estado magico en el que espera
		// solo por jugadas de tipo flor-related
		// lo mismo para el real-envido; falta-envido
		jid := p.Ronda.Envite.SinCantar[0]
		// j := p.Ronda.Manojo(jid)
		siguienteJugada := CantarFlor{jid}
		res2 := siguienteJugada.Hacer(p)
		if p.Verbose {
			pkts2 = append(pkts2, res2...)
		}

	} else {
		p.TocarEnvido(p.Manojo(jugada.JID))
	}

	return pkts2
}

/* el problema de esta funcion es que esta mas relacionada con el `quiero`
que con el envido. Deberia formar parte del eval del quiero */

// donde 'j' el jugador que dijo 'quiero' al 'envido'/'real envido'
func (jugada TocarEnvido) Eval(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)

	p.Ronda.Envite.Estado = DESHABILITADO
	p.Ronda.Envite.SinCantar = []string{}
	jIdx, _, res2 := p.Ronda.ExecElEnvido(p.Verbose)

	if p.Verbose {
		pkts2 = append(pkts2, res2...)
	}

	jug := p.Ronda.Manojos[jIdx].Jugador

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.SumaPts{
				Autor:  jug.ID,
				Razon:  enco.EnvidoGanado,
				Puntos: p.Ronda.Envite.Puntaje,
			},
		))
	}

	p.SumarPuntos(jug.Equipo, p.Ronda.Envite.Puntaje)

	return pkts2
}

type TocarRealEnvido struct {
	// Manojo *Manojo
	JID string
}

func (jugada TocarRealEnvido) ID() IJUGADA_ID {
	return JID_REAL_ENVIDO
}

func (jugada TocarRealEnvido) String() string {
	return jugada.JID + " real-envido"
}

func (jugada TocarRealEnvido) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	if florEnJuego {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tocar real envido ahora porque la flor esta en juego"),
			))
		}

		return pkts2, false

	}
	seFueAlMazo := p.Manojo(jugada.JID).SeFueAlMazo
	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	esSuTurno := p.Ronda.GetElTurno().Jugador.ID == p.Manojo(jugada.JID).Jugador.ID
	tieneFlor, _ := p.Manojo(jugada.JID).TieneFlor(p.Ronda.Muestra)
	realEnvidoHabilitado := (p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Envite.Estado == ENVIDO)

	if !realEnvidoHabilitado {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tocar real-envido ahora"),
			))
		}

		return pkts2, false
	}

	esDelEquipoContrario := p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO
	trucoNoCantado := p.Ronda.Truco.Estado == NOGRITADOAUN

	estaIniciandoPorPrimeraVezElEnvido := esSuTurno && p.Ronda.Envite.Estado == NOCANTADOAUN && trucoNoCantado
	estaRedoblandoLaApuesta := p.Ronda.Envite.Estado == ENVIDO && esDelEquipoContrario // cuando redobla una apuesta puede o no ser su turno
	elEnvidoEstaPrimero := !esSuTurno && p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	puedeTocarRealEnvido := estaIniciandoPorPrimeraVezElEnvido || estaRedoblandoLaApuesta || elEnvidoEstaPrimero
	ok := !seFueAlMazo && (realEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario) && puedeTocarRealEnvido

	if !ok {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error(`No es posible cantar 'Real Envido'`),
			))
		}

		return pkts2, false

	}

	return pkts2, true
}

func (jugada TocarRealEnvido) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO
	elEnvidoEstaPrimero := p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	if elEnvidoEstaPrimero {
		// deshabilito el truco
		p.Ronda.Truco.Estado = NOGRITADOAUN
		p.Ronda.Truco.CantadoPor = ""

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.ElEnvidoEstaPrimero(jugada.JID),
			))
		}

	}

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.TocarRealEnvido(jugada.JID),
		))
	}

	p.TocarRealEnvido(p.Manojo(jugada.JID))

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.SinCantar) > 0

	if hayFlor {
		jid := p.Ronda.Envite.SinCantar[0]
		// j := p.Ronda.Manojo(jid)
		siguienteJugada := CantarFlor{jid}
		res2 := siguienteJugada.Hacer(p)
		if p.Verbose {
			pkts2 = append(pkts2, res2...)
		}
	}

	return pkts2
}

type TocarFaltaEnvido struct {
	// Manojo *Manojo
	JID string
}

func (jugada TocarFaltaEnvido) ID() IJUGADA_ID {
	return JID_FALTA_ENVIDO
}

func (jugada TocarFaltaEnvido) String() string {
	return jugada.JID + " falta-envido"
}

func (jugada TocarFaltaEnvido) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// checkeo flor en juego
	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	if florEnJuego {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tocar falta envido ahora porque la flor esta en juego"),
			))
		}

		return pkts2, false

	}
	seFueAlMazo := p.Manojo(jugada.JID).SeFueAlMazo
	esSuTurno := p.Ronda.GetElTurno().Jugador.ID == p.Manojo(jugada.JID).Jugador.ID
	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	tieneFlor, _ := p.Manojo(jugada.JID).TieneFlor(p.Ronda.Muestra)
	faltaEnvidoHabilitado := p.Ronda.Envite.Estado >= NOCANTADOAUN && p.Ronda.Envite.Estado < FALTAENVIDO

	if !faltaEnvidoHabilitado {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible tocar real-envido ahora"),
			))
		}

		return pkts2, false
	}

	esDelEquipoContrario := p.Ronda.Envite.Estado == NOCANTADOAUN || p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado >= ENVIDO
	trucoNoCantado := p.Ronda.Truco.Estado == NOGRITADOAUN

	estaIniciandoPorPrimeraVezElEnvido := esSuTurno && p.Ronda.Envite.Estado == NOCANTADOAUN && trucoNoCantado
	estaRedoblandoLaApuesta := p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado < FALTAENVIDO && esDelEquipoContrario // cuando redobla una apuesta puede o no ser su turno
	elEnvidoEstaPrimero := !esSuTurno && p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	puedeTocarFaltaEnvido := estaIniciandoPorPrimeraVezElEnvido || estaRedoblandoLaApuesta || elEnvidoEstaPrimero
	ok := !seFueAlMazo && (faltaEnvidoHabilitado && esPrimeraMano && !tieneFlor && esDelEquipoContrario) && puedeTocarFaltaEnvido

	if !ok {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error(`No es posible cantar 'Falta Envido'`),
			))
		}

		return pkts2, false

	}

	return pkts2, true
}

func (jugada TocarFaltaEnvido) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	yaEstabamosEnEnvido := p.Ronda.Envite.Estado == ENVIDO || p.Ronda.Envite.Estado == REALENVIDO
	elEnvidoEstaPrimero := p.Ronda.Truco.Estado == TRUCO && !yaEstabamosEnEnvido && esPrimeraMano

	if elEnvidoEstaPrimero {
		// deshabilito el truco
		p.Ronda.Truco.Estado = NOGRITADOAUN
		p.Ronda.Truco.CantadoPor = ""

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.ElEnvidoEstaPrimero(jugada.JID),
			))
		}

	}

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.TocarFaltaEnvido(jugada.JID),
		))
	}

	p.TocarFaltaEnvido(p.Manojo(jugada.JID))

	// ahora checkeo si alguien tiene flor
	hayFlor := len(p.Ronda.Envite.SinCantar) > 0
	if hayFlor {
		jid := p.Ronda.Envite.SinCantar[0]
		// j := p.Ronda.Manojo(jid)
		siguienteJugada := CantarFlor{jid}
		res2 := siguienteJugada.Hacer(p)
		if p.Verbose {
			pkts2 = append(pkts2, res2...)
		}
	}

	return pkts2
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

func (jugada TocarFaltaEnvido) Eval(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)

	p.Ronda.Envite.Estado = DESHABILITADO
	p.Ronda.Envite.SinCantar = []string{}

	// computar envidos
	jIdx, _, res2 := p.Ronda.ExecElEnvido(p.Verbose)

	if p.Verbose {
		pkts2 = append(pkts2, res2...)
	}

	// jug es el que gano el (falta) envido
	jug := p.Ronda.Manojos[jIdx].Jugador

	pts := p.CalcPtsFaltaEnvido(jug.Equipo)

	p.Ronda.Envite.Puntaje += pts

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.SumaPts{
				Autor:  jug.ID,
				Razon:  enco.FaltaEnvidoGanado,
				Puntos: p.Ronda.Envite.Puntaje,
			},
		))
	}

	p.SumarPuntos(jug.Equipo, p.Ronda.Envite.Puntaje)

	return pkts2
}

type CantarFlor struct {
	// Manojo *Manojo
	JID string
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
func (jugada CantarFlor) ID() IJUGADA_ID {
	return JID_FLOR
}

func (jugada CantarFlor) String() string {
	return jugada.JID + " flor"
}

func (jugada CantarFlor) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// manojo dice que puede cantar flor;
	// es esto verdad?
	seFueAlMazo := p.Manojo(jugada.JID).SeFueAlMazo
	florHabilitada := (p.Ronda.Envite.Estado >= NOCANTADOAUN) && p.Ronda.ManoEnJuego == Primera
	tieneFlor, _ := p.Manojo(jugada.JID).TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := p.Ronda.Envite.noCantoFlorAun(p.Manojo(jugada.JID).Jugador.ID)

	// caso especial:
	// tienen flor: alice bob ben.
	// alice:flor -> bob:contra-flor -> alice:mazo => ben ??? no canto su flor
	// entonces, puede cantar flor, (SIN disminuir su estado) si tiene flor Y NO CANTO AUN
	// por eso le elimino la clausura:
	// `p.Ronda.Envite.Estado <= FLOR` en la variable `florHabilitada`

	ok := !seFueAlMazo && florHabilitada && tieneFlor && noCantoFlorAun

	if !ok {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error(`No es posible cantar flor`),
			))
		}

		return pkts2, false

	}

	return pkts2, true
}

func (jugada CantarFlor) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	// yo canto

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.CantarFlor(jugada.JID),
		))
	}

	// corresponde que desactive el truco?
	// si lo desactivo: es medio tedioso para el usuario tener q volver a gritar
	// si no lo desacivo: medio como que se olvida
	// QUEDA CONSISTENTE CON "EL ENVIDO ESTA PRIMERO"!
	p.Ronda.Truco.CantadoPor = ""
	p.Ronda.Truco.Estado = NOGRITADOAUN

	// y me elimino de los que no-cantaron
	p.Ronda.Envite.cantoFlor(p.Manojo(jugada.JID).Jugador.ID)

	p.CantarFlor(p.Manojo(jugada.JID))

	// es el ultimo en cantar flor que faltaba?
	// o simplemente es el unico que tiene flor (caso particular)

	todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.SinCantar) == 0
	if todosLosJugadoresConFlorCantaron {

		florPkts2 := evalFlor(p)
		if p.Verbose {
			pkts2 = append(pkts2, florPkts2...)
		}

	} else {

		// cachear esto
		// solos los de su equipo tienen flor?
		// si solos los de su equipo tienen flor (y los otros no) -> las canto todas
		soloLosDeSuEquipoTienenFlor := true
		for _, manojo := range p.Ronda.Envite.JugadoresConFlor {
			if manojo.Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo {
				soloLosDeSuEquipoTienenFlor = false
				break
			}
		}

		if soloLosDeSuEquipoTienenFlor {
			// los quiero llamar a todos, pero no quiero Hacer llamadas al pedo
			// entonces: llamo al primero sin cantar, y que este llame al proximo
			// y que el proximo llame al siguiente, y asi...
			jid := p.Ronda.Envite.SinCantar[0]
			// j := p.Ronda.Manojo(jid)
			siguienteJugada := CantarFlor{jid}
			res2 := siguienteJugada.Hacer(p)
			if p.Verbose {
				pkts2 = append(pkts2, res2...)
			}
		}

	}

	return pkts2
}

func evalFlor(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)

	florEnJuego := p.Ronda.Envite.Estado >= FLOR
	todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.SinCantar) == 0
	ok := todosLosJugadoresConFlorCantaron && florEnJuego
	if !ok {
		return pkts2
	}

	// cual es la flor ganadora?
	// empieza cantando el autor del envite no el que "quizo"
	autorIdx := p.Ronda.JIX(p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.ID)
	manojoConLaFlorMasAlta, _, res2 := p.Ronda.ExecLaFlores(JIX(autorIdx), p.Verbose)

	if p.Verbose {
		pkts2 = append(pkts2, res2...)
	}

	equipoGanador := manojoConLaFlorMasAlta.Jugador.Equipo

	// que estaba en juego?
	// switch p.Ronda.Envite.Estado {
	// case FLOR:
	// ahora se quien es el ganador; necesito saber cuantos puntos
	// se le va a sumar a ese equipo:
	// los acumulados del envite hasta ahora
	puntosASumar := p.Ronda.Envite.Puntaje
	p.SumarPuntos(equipoGanador, puntosASumar)
	habiaSolo1JugadorConFlor := len(p.Ronda.Envite.JugadoresConFlor) == 1
	if habiaSolo1JugadorConFlor {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.SumaPts{
					Autor:  manojoConLaFlorMasAlta.Jugador.ID,
					Razon:  enco.LaUnicaFlor,
					Puntos: puntosASumar,
				},
			))
		}

	} else {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.SumaPts{
					Autor:  manojoConLaFlorMasAlta.Jugador.ID,
					Razon:  enco.LaFlorMasAlta,
					Puntos: puntosASumar,
				},
			))
		}

	}
	// case CONTRAFLOR:
	// case CONTRAFLORALRESTO:
	// }

	p.Ronda.Envite.Estado = DESHABILITADO
	p.Ronda.Envite.SinCantar = []string{}

	return pkts2
}

type CantarContraFlor struct {
	// Manojo *Manojo
	JID string
}

func (jugada CantarContraFlor) ID() IJUGADA_ID {
	return JID_CONTRA_FLOR
}

func (jugada CantarContraFlor) String() string {
	return jugada.JID + " contra-flor"
}

func (jugada CantarContraFlor) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// manojo dice que puede cantar flor;
	// es esto verdad?
	seFueAlMazo := p.Manojo(jugada.JID).SeFueAlMazo
	contraFlorHabilitada := p.Ronda.Envite.Estado == FLOR && p.Ronda.ManoEnJuego == Primera
	esDelEquipoContrario := contraFlorHabilitada && p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
	tieneFlor, _ := p.Manojo(jugada.JID).TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := p.Ronda.Envite.noCantoFlorAun(p.Manojo(jugada.JID).Jugador.ID)
	ok := !seFueAlMazo && contraFlorHabilitada && tieneFlor && esDelEquipoContrario && noCantoFlorAun
	if !ok {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error(`No es posible cantar contra flor`),
			))
		}

		return pkts2, false

	}

	return pkts2, true
}

func (jugada CantarContraFlor) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	// la canta

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.CantarContraFlor(jugada.JID),
		))
	}

	p.CantarContraFlor(p.Manojo(jugada.JID))
	// y ahora tengo que esperar por la respuesta de la nueva
	// propuesta de todos menos de el que canto la contraflor
	// restauro la copia
	p.Ronda.Envite.cantoFlor(p.Manojo(jugada.JID).Jugador.ID)

	return pkts2
}

type CantarContraFlorAlResto struct {
	// Manojo *Manojo
	JID string
}

func (jugada CantarContraFlorAlResto) ID() IJUGADA_ID {
	return JID_CONTRA_FLOR_AL_RESTO
}

func (jugada CantarContraFlorAlResto) String() string {
	return jugada.JID + " contra-flor-al-resto"
}

func (jugada CantarContraFlorAlResto) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// manojo dice que puede cantar flor;
	// es esto verdad?
	seFueAlMazo := p.Manojo(jugada.JID).SeFueAlMazo
	contraFlorHabilitada := (p.Ronda.Envite.Estado == FLOR || p.Ronda.Envite.Estado == CONTRAFLOR) && p.Ronda.ManoEnJuego == Primera
	esDelEquipoContrario := contraFlorHabilitada && p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
	tieneFlor, _ := p.Manojo(jugada.JID).TieneFlor(p.Ronda.Muestra)
	noCantoFlorAun := p.Ronda.Envite.noCantoFlorAun(p.Manojo(jugada.JID).Jugador.ID)
	ok := !seFueAlMazo && contraFlorHabilitada && tieneFlor && esDelEquipoContrario && noCantoFlorAun
	if !ok {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error(`No es posible cantar contra flor al resto`),
			))
		}

		return pkts2, false

	}

	return pkts2, true
}

func (jugada CantarContraFlorAlResto) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	// la canta

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.CantarContraFlorAlResto(jugada.JID),
		))
	}

	p.CantarContraFlorAlResto(p.Manojo(jugada.JID))
	// y ahora tengo que esperar por la respuesta de la nueva
	// propuesta de todos menos de el que canto la contraflor
	// restauro la copia
	p.Ronda.Envite.cantoFlor(p.Manojo(jugada.JID).Jugador.ID)

	return pkts2
}

type CantarConFlorMeAchico struct {
	// Manojo *Manojo
	JID string
}

// no implementada (porque no es necesaria?)
func (jugada CantarConFlorMeAchico) Hacer(p *Partida) []enco.Envelope {
	return nil
}

type GritarTruco struct {
	// Manojo *Manojo
	JID string
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

func (jugada GritarTruco) ID() IJUGADA_ID {
	return JID_TRUCO
}

func (jugada GritarTruco) String() string {
	return jugada.JID + " truco"
}

func (jugada GritarTruco) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// checkeos:
	noSeFueAlMazo := !p.Manojo(jugada.JID).SeFueAlMazo
	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= NOCANTADOAUN

	yoOUnoDeMisCompasTieneFlorYAunNoCanto := p.Ronda.hayEquipoSinCantar(p.Manojo(jugada.JID).Jugador.Equipo)

	laFlorEstaPrimero := yoOUnoDeMisCompasTieneFlorYAunNoCanto
	trucoNoSeJugoAun := p.Ronda.Truco.Estado == NOGRITADOAUN
	esSuTurno := p.Ronda.GetElTurno().Jugador.ID == p.Manojo(jugada.JID).Jugador.ID
	trucoHabilitado := noSeFueAlMazo && trucoNoSeJugoAun && noSeEstaJugandoElEnvite && !laFlorEstaPrimero && esSuTurno

	if !trucoHabilitado {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible cantar truco ahora"),
			))
		}

		return pkts2, false

	}

	return pkts2, true
}

func (jugada GritarTruco) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.GritarTruco(jugada.JID),
		))
	}

	p.GritarTruco(p.Manojo(jugada.JID))

	return pkts2
}

type GritarReTruco struct {
	// Manojo *Manojo
	JID string
}

func (jugada GritarReTruco) ID() IJUGADA_ID {
	return JID_RE_TRUCO
}

func (jugada GritarReTruco) String() string {
	return jugada.JID + " re-truco"
}

func (jugada GritarReTruco) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// checkeos generales:
	noSeFueAlMazo := !p.Manojo(jugada.JID).SeFueAlMazo
	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= NOCANTADOAUN

	yoOUnoDeMisCompasTieneFlorYAunNoCanto := p.Ronda.hayEquipoSinCantar(p.Manojo(jugada.JID).Jugador.Equipo)

	laFlorEstaPrimero := yoOUnoDeMisCompasTieneFlorYAunNoCanto

	/*
		Hay 2 casos para cantar rectruco:
		    - CASO I: Uno del equipo contrario grito el truco
			- CASO II: Uno de su equipo posee el quiero
	*/

	// CASO I:
	trucoGritado := p.Ronda.Truco.Estado == TRUCO
	unoDelEquipoContrarioGritoTruco := trucoGritado && p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
	casoI := trucoGritado && unoDelEquipoContrarioGritoTruco

	// CASO II:
	trucoYaQuerido := p.Ronda.Truco.Estado == TRUCOQUERIDO
	unoDeMiEquipoQuizo := trucoYaQuerido && p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == p.Manojo(jugada.JID).Jugador.Equipo
	// esTurnoDeMiEquipo := p.Ronda.GetElTurno().Jugador.Equipo == p.Manojo(jugada.JID).Jugador.Equipo
	casoII := trucoYaQuerido && unoDeMiEquipoQuizo // && esTurnoDeMiEquipo

	reTrucoHabilitado := noSeFueAlMazo && noSeEstaJugandoElEnvite && (casoI || casoII) && !laFlorEstaPrimero

	if !reTrucoHabilitado {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible cantar re-truco ahora"),
			))
		}

		return pkts2, false

	}

	return pkts2, true
}

// checkeaos de este tipo:
// que pasa cuando gritan re-truco cuando el campo truco se encuentra nil
// ese fue el nil pointer exception
func (jugada GritarReTruco) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.GritarReTruco(jugada.JID),
		))
	}

	p.GritarReTruco(p.Manojo(jugada.JID))

	return pkts2
}

type GritarVale4 struct {
	// Manojo *Manojo
	JID string
}

func (jugada GritarVale4) ID() IJUGADA_ID {
	return JID_VALE_4
}

func (jugada GritarVale4) String() string {
	return jugada.JID + " vale-4"
}

func (jugada GritarVale4) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// checkeos:
	noSeFueAlMazo := !p.Manojo(jugada.JID).SeFueAlMazo

	noSeEstaJugandoElEnvite := p.Ronda.Envite.Estado <= NOCANTADOAUN

	yoOUnoDeMisCompasTieneFlorYAunNoCanto := p.Ronda.hayEquipoSinCantar(p.Manojo(jugada.JID).Jugador.Equipo)

	laFlorEstaPrimero := yoOUnoDeMisCompasTieneFlorYAunNoCanto

	/*
		Hay 2 casos para cantar rectruco:
		    - CASO I: Uno del equipo contrario grito el re-truco
			- CASO II: Uno de su equipo posee el quiero
	*/

	// CASO I:
	reTrucoGritado := p.Ronda.Truco.Estado == RETRUCO
	// para eviat el nil primero checkeo que haya sido gritado reTrucoGritado &&
	unoDelEquipoContrarioGritoReTruco := reTrucoGritado && p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
	casoI := reTrucoGritado && unoDelEquipoContrarioGritoReTruco

	// CASO I:
	retrucoYaQuerido := p.Ronda.Truco.Estado == RETRUCOQUERIDO
	// para eviat el nil primero checkeo que haya sido gritado reTrucoGritado &&
	suEquipotieneElQuiero := retrucoYaQuerido && p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == p.Manojo(jugada.JID).Jugador.Equipo
	casoII := retrucoYaQuerido && suEquipotieneElQuiero

	vale4Habilitado := noSeFueAlMazo && (casoI || casoII) && noSeEstaJugandoElEnvite && !laFlorEstaPrimero

	if !vale4Habilitado {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible cantar vale-4 ahora"),
			))
		}

		return pkts2, false

	}

	return pkts2, true
}

func (jugada GritarVale4) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.GritarVale4(jugada.JID),
		))
	}

	p.GritarVale4(p.Manojo(jugada.JID))

	return pkts2
}

type ResponderQuiero struct {
	// Manojo *Manojo
	JID string
}

func (jugada ResponderQuiero) ID() IJUGADA_ID {
	return JID_QUIERO
}

func (jugada ResponderQuiero) String() string {
	return jugada.JID + " quiero"
}

func (jugada ResponderQuiero) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	seFueAlMazo := p.Manojo(jugada.JID).SeFueAlMazo
	if seFueAlMazo {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("Te fuiste al mazo; no podes Hacer esta jugada"),
			))
		}

		return pkts2, false

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

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible responder quiero ahora"),
			))
		}

		return pkts2, false
	}

	noHanCantadoLaFlorAun := p.Ronda.Envite.Estado < FLOR
	yoOUnoDeMisCompasTieneFlorYAunNoCanto := p.Ronda.hayEquipoSinCantar(p.Manojo(jugada.JID).Jugador.Equipo)
	if noHanCantadoLaFlorAun && yoOUnoDeMisCompasTieneFlorYAunNoCanto {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible responder 'quiero' porque alguien con flor no ha cantado aun"),
			))
		}

		return pkts2, false
	}
	// se acepta una respuesta 'quiero' solo cuando:
	// - CASO I: se toco un envite+ (con autor del equipo contario)
	// - CASO II: se grito el truco+ (con autor del equipo contario)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO)
	// ojo: solo a la contraflor+ se le puede decir quiero; a la flor sola no
	laContraFlorEsRespondible := p.Ronda.Envite.Estado >= CONTRAFLOR && p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
	elTrucoEsRespondible := p.Ronda.Truco.Estado.esTrucoRespondible() && p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo

	ok := elEnvidoEsRespondible || laContraFlorEsRespondible || elTrucoEsRespondible
	if !ok {
		// si no, esta respondiendo al pedo

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error(`No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo`),
			))
		}

		return pkts2, false

	}

	if elEnvidoEsRespondible {

		esDelEquipoContrario := p.Manojo(jugada.JID).Jugador.Equipo != p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo
		if !esDelEquipoContrario {

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.Dest(jugada.JID),
					enco.Error(`La jugada no es valida`),
				))
			}

			return pkts2, false

		}

	} else if laContraFlorEsRespondible {
		// tengo que verificar si efectivamente tiene flor
		tieneFlor, _ := p.Manojo(jugada.JID).TieneFlor(p.Ronda.Muestra)
		esDelEquipoContrario := p.Manojo(jugada.JID).Jugador.Equipo != p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo
		ok := tieneFlor && esDelEquipoContrario

		if !ok {

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.Dest(jugada.JID),
					enco.Error(`La jugada no es valida`),
				))
			}

			return pkts2, false

		}

	}

	return pkts2, true
}

func (jugada ResponderQuiero) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	// se acepta una respuesta 'quiero' solo cuando:
	// - CASO I: se toco un envite+ (con autor del equipo contario)
	// - CASO II: se grito el truco+ (con autor del equipo contario)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO)
	// ojo: solo a la contraflor+ se le puede decir quiero; a la flor sola no
	laContraFlorEsRespondible := p.Ronda.Envite.Estado >= CONTRAFLOR && p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
	elTrucoEsRespondible := p.Ronda.Truco.Estado.esTrucoRespondible() && p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo

	if elEnvidoEsRespondible {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.QuieroEnvite(jugada.JID),
			))
		}

		if p.Ronda.Envite.Estado == FALTAENVIDO {
			res2 := TocarFaltaEnvido(jugada).Eval(p)
			return append(pkts2, res2...)
		}
		// si no, era envido/real-envido o cualquier
		// combinacion valida de ellos

		res2 := TocarEnvido(jugada).Eval(p)
		return append(pkts2, res2...)

	} else if laContraFlorEsRespondible {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.QuieroEnvite(jugada.JID),
			))
		}

		// empieza cantando el autor del envite no el que "quizo"
		autorIdx := p.Ronda.JIX(p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.ID)
		manojoConLaFlorMasAlta, _, res2 := p.Ronda.ExecLaFlores(JIX(autorIdx), p.Verbose)

		if p.Verbose {
			pkts2 = append(pkts2, res2...)
		}

		// manojoConLaFlorMasAlta, _ := p.Ronda.GetLaFlorMasAlta()
		equipoGanador := manojoConLaFlorMasAlta.Jugador.Equipo

		if p.Ronda.Envite.Estado == CONTRAFLOR {
			puntosASumar := p.Ronda.Envite.Puntaje
			p.SumarPuntos(equipoGanador, puntosASumar)

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.ALL,
					enco.SumaPts{
						Autor:  manojoConLaFlorMasAlta.Jugador.ID,
						Razon:  enco.ContraFlorGanada,
						Puntos: puntosASumar,
					},
				))
			}

		} else {
			// el equipo del ganador de la contraflor al resto
			// gano la partida
			// duda se cuentan las flores?
			// puntosASumar := p.Ronda.Envite.Puntaje + p.CalcPtsContraFlorAlResto(equipoGanador)
			puntosASumar := p.CalcPtsContraFlorAlResto(equipoGanador)
			p.SumarPuntos(equipoGanador, puntosASumar)

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.ALL,
					enco.SumaPts{
						Autor:  manojoConLaFlorMasAlta.Jugador.ID,
						Razon:  enco.ContraFlorAlRestoGanada,
						Puntos: puntosASumar,
					},
				))
			}

		}

		p.Ronda.Envite.Estado = DESHABILITADO
		p.Ronda.Envite.SinCantar = []string{}

	} else if elTrucoEsRespondible {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.QuieroTruco(jugada.JID),
			))
		}

		p.QuererTruco(p.Manojo(jugada.JID))
	}

	return pkts2

}

type ResponderNoQuiero struct {
	// Manojo *Manojo
	JID string
}

func (jugada ResponderNoQuiero) ID() IJUGADA_ID {
	return JID_NO_QUIERO
}

func (jugada ResponderNoQuiero) String() string {
	return jugada.JID + " no-quiero"
}

func (jugada ResponderNoQuiero) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	seFueAlMazo := p.Manojo(jugada.JID).SeFueAlMazo
	if seFueAlMazo {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("Te fuiste al mazo; no podes Hacer esta jugada"),
			))
		}

		return pkts2, false
	}

	// checkeo flor en juego
	// caso particular del checkeo: no se le puede decir quiero a la flor
	// pero si a la contra flor o contra flor al resto
	// FALSO porque el no quiero lo estoy contando como un "con flor me achico"
	// todo: agregar la jugada: "con flor me achico" y editar la variale:
	// AHORA:
	// laFlorEsRespondible := p.Ronda.Flor >= FLOR && p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.equipo != p.Manojo(jugada.JID).Jugador.Equipo
	// LUEGO DE AGREGAR LA JUGADA "con flor me achico"
	// laFlorEsRespondible := p.Ronda.Flor > FLOR
	// FALSO ---> directamente se va la posibilidad de reponderle
	// "no quiero a la flor"

	// se acepta una respuesta 'no quiero' solo cuando:
	// - CASO I: se toco el envido (o similar)
	// - CASO II: se grito el truco (o similar)
	// en caso contrario, es incorrecto -> error

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO) && p.Ronda.Envite.CantadoPor != p.Manojo(jugada.JID).Jugador.ID
	laFlorEsRespondible := p.Ronda.Envite.Estado >= FLOR && p.Ronda.Envite.CantadoPor != p.Manojo(jugada.JID).Jugador.ID
	elTrucoEsRespondible := p.Ronda.Truco.Estado.esTrucoRespondible() && p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo

	ok := elEnvidoEsRespondible || laFlorEsRespondible || elTrucoEsRespondible

	if !ok {
		// si no, esta respondiendo al pedo
		err := p.Manojo(jugada.JID).Jugador.ID + ` esta respondiendo al pedo; no hay nada respondible`

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error(err),
			))
		}

		return pkts2, false
	}

	if elEnvidoEsRespondible {

		esDelEquipoContrario := p.Manojo(jugada.JID).Jugador.Equipo != p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo
		if !esDelEquipoContrario {

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.Dest(jugada.JID),
					enco.Error(`La jugada no es valida`),
				))
			}

			return pkts2, false

		}

	} else if laFlorEsRespondible {

		// tengo que verificar si efectivamente tiene flor
		tieneFlor, _ := p.Manojo(jugada.JID).TieneFlor(p.Ronda.Muestra)
		esDelEquipoContrario := p.Manojo(jugada.JID).Jugador.Equipo != p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo
		ok := tieneFlor && esDelEquipoContrario

		if !ok {

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.Dest(jugada.JID),
					enco.Error(`La jugada no es valida`),
				))
			}

			return pkts2, false

		}

	}

	return pkts2, true
}

func (jugada ResponderNoQuiero) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	elEnvidoEsRespondible := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO) && p.Ronda.Envite.CantadoPor != p.Manojo(jugada.JID).Jugador.ID
	laFlorEsRespondible := p.Ronda.Envite.Estado >= FLOR && p.Ronda.Envite.CantadoPor != p.Manojo(jugada.JID).Jugador.ID
	elTrucoEsRespondible := p.Ronda.Truco.Estado.esTrucoRespondible() && p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo

	if elEnvidoEsRespondible {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.NoQuiero(jugada.JID),
			))
		}

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
		p.Ronda.Envite.SinCantar = []string{}
		p.Ronda.Envite.Puntaje = totalPts

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.SumaPts{
					Autor:  p.Ronda.Envite.CantadoPor,
					Razon:  enco.EnviteNoQuerido,
					Puntos: totalPts,
				},
			))
		}

		p.SumarPuntos(p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo, totalPts)

	} else if laFlorEsRespondible {

		// todo ok: tiene flor; se pasa a jugar:

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.ConFlorMeAchico(jugada.JID),
			))
		}

		// cuenta como un "no quiero" (codigo copiado)
		// segun el estado de la apuesta actual:
		// los "me achico" no cuentan para la flor
		// Flor		xcg(+3) / xcg(+3)
		// Flor + Contra-Flor		xc(+3) / xCadaFlorDelQueHizoElDesafio(+3) + 1
		// Flor + [Contra-Flor] + ContraFlorAlResto		~Falta Envido + *TODAS* las flores no achicadas / xcg(+3) + 1

		// sumo todas las flores del equipo contrario
		totalPts := 0

		for _, m := range p.Ronda.Manojos {
			esDelEquipoContrario := p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
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
		p.Ronda.Envite.SinCantar = []string{}

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.SumaPts{
					Autor:  p.Ronda.Envite.CantadoPor,
					Razon:  enco.FlorAchicada,
					Puntos: totalPts,
				},
			))
		}

		p.SumarPuntos(p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo, totalPts)

	} else if elTrucoEsRespondible {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.ALL,
				enco.NoQuiero(jugada.JID),
			))
		}

		// pongo al equipo que propuso el truco como ganador de la mano actual
		manoActual := p.Ronda.ManoEnJuego.ToInt() - 1
		p.Ronda.Manos[manoActual].Ganador = p.Ronda.Truco.CantadoPor
		equipoGanador := GanoAzul
		if p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == Rojo {
			equipoGanador = GanoRojo
		}
		p.Ronda.Manos[manoActual].Resultado = equipoGanador

		NuevaRonda, res2 := p.EvaluarRonda()

		if p.Verbose {
			pkts2 = append(pkts2, res2...)
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
				if p.Verbose {
					for _, m := range p.Ronda.Manojos {
						pkts2 = append(pkts2, enco.Pkt(
							enco.ALL,
							enco.NuevaRonda{
								Perspectiva: p.PerspectivaCacheFlor(&m),
							},
						))
					}
				}

			} // else {
			// p.byeBye()
			// }

		}

	}

	return pkts2
}

type IrseAlMazo struct {
	// Manojo *Manojo
	JID string
}

func (jugada IrseAlMazo) ID() IJUGADA_ID {
	return JID_MAZO
}

func (jugada IrseAlMazo) String() string {
	return jugada.JID + " mazo"
}

func (jugada IrseAlMazo) Ok(p *Partida) ([]enco.Envelope, bool) {

	pkts2 := make([]enco.Envelope, 0)

	// checkeos:
	yaSeFueAlMazo := p.Manojo(jugada.JID).SeFueAlMazo
	yaTiroTodasSusCartas := p.Manojo(jugada.JID).GetCantCartasTiradas() == 3
	if yaSeFueAlMazo || yaTiroTodasSusCartas {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible irse al mazo ahora"),
			))
		}

		return pkts2, false

	}

	seEstabaJugandoElEnvido := (p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO)
	seEstabaJugandoLaFlor := p.Ronda.Envite.Estado >= FLOR
	seEstabaJugandoElTruco := p.Ronda.Truco.Estado.esTrucoRespondible()
	// no se puede ir al mazo sii:
	// 1. el fue el que canto el envido (y el envido esta en juego)
	// 2. tampoco se puede ir al mazo si el canto la flor o similar
	// 3. tampoco se puede ir al mazo si el grito el truco

	// envidoPropuesto := Contains([]EstadoEnvite{ENVIDO, REALENVIDO, FALTAENVIDO}, p.Ronda.Envite.Estado)
	// envidoPropuestoPorSuEquipo := p.Ronda.Manojo[p.Ronda.Envite.CantadoPor].Jugador.Equipo == p.Manojo(jugada.JID).Jugador.Equipo
	// trucoPropuesto := p.Ronda.Truco.Estado.esTrucoRespondible()
	// trucoPropuestoPorSuEquipo := p.Ronda.Manojo[p.Ronda.Truco.CantadoPor].Jugador.Equipo == p.Manojo(jugada.JID).Jugador.Equipo
	// condicionDelBobo := (envidoPropuesto && envidoPropuestoPorSuEquipo) || (trucoPropuesto && trucoPropuestoPorSuEquipo)

	// if condicionDelBobo {

	// enco.Write(p.Stdout, enco.Pkt1(
	// 	enco.Dest(jugada.JID),
	// 	enco.Msg(enco.Error,  fmt.Sprintf("No es posible irse al mazo ahora porque hay propuestas de tu equipo sin responder")),
	// ))

	// return

	// }

	noSePuedeIrPorElEnvite := (seEstabaJugandoElEnvido || seEstabaJugandoLaFlor) && p.Ronda.Envite.CantadoPor == p.Manojo(jugada.JID).Jugador.ID
	// la de la flor es igual al del envido; porque es un envite
	noSePuedeIrPorElTruco := seEstabaJugandoElTruco && p.Ronda.Truco.CantadoPor == p.Manojo(jugada.JID).Jugador.ID
	if noSePuedeIrPorElEnvite || noSePuedeIrPorElTruco {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible irse al mazo ahora"),
			))
		}

		return pkts2, false

	}

	// por como esta hecho el algoritmo EvaluarMano:

	esPrimeraMano := p.Ronda.ManoEnJuego == Primera
	tiradas := p.Ronda.GetManoActual().CartasTiradas
	n := len(tiradas)

	soloMiEquipoTiro := n == 1 && p.Ronda.Manojo(tiradas[n-1].Jugador).Jugador.Equipo == p.Manojo(jugada.JID).Jugador.Equipo

	equipoDelJugador := p.Manojo(jugada.JID).Jugador.Equipo
	soyElUnicoDeMiEquipo := p.Ronda.CantJugadoresEnJuego[equipoDelJugador] == 1
	noSePuedeIr := esPrimeraMano && soloMiEquipoTiro && soyElUnicoDeMiEquipo

	// que pasa si alguien dice truco y se va al mazo?

	if noSePuedeIr {

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest(jugada.JID),
				enco.Error("No es posible irse al mazo ahora"),
			))
		}

		return pkts2, false
	}

	return pkts2, true
}

// Eliminar .
func Eliminar(manojos []*Manojo, manojo *Manojo) []*Manojo {
	var i int
	// primero encuentro el elemento
	for i = 0; i <= len(manojos); i++ {
		noLoContiene := i == len(manojos)
		if noLoContiene {
			return manojos
		}
		if manojos[i].Jugador.ID == manojo.Jugador.ID {
			break
		}
	}
	manojos[i] = manojos[len(manojos)-1] // Copy last element to index i.
	return manojos[:len(manojos)-1]      // Truncate slice.
}

func (jugada IrseAlMazo) Hacer(p *Partida) []enco.Envelope {

	pkts2 := make([]enco.Envelope, 0)
	pre2, ok := jugada.Ok(p)
	if p.Verbose {
		pkts2 = append(pkts2, pre2...)
	}

	if !ok {
		return pkts2
	}

	// ok -> se va al mazo:

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.ALL,
			enco.Mazo(jugada.JID),
		))
	}

	p.IrAlMazo(p.Manojo(jugada.JID))

	equipoDelJugador := p.Manojo(jugada.JID).Jugador.Equipo

	seFueronTodos := p.Ronda.CantJugadoresEnJuego[equipoDelJugador] == 0

	// si tenia flor -> ya no lo tomo en cuenta
	tieneFlor, _ := p.Manojo(jugada.JID).TieneFlor(p.Ronda.Muestra)
	if tieneFlor {
		p.Ronda.Envite.JugadoresConFlor = Eliminar(p.Ronda.Envite.JugadoresConFlor, p.Manojo(jugada.JID))
		p.Ronda.Envite.cantoFlor(p.Manojo(jugada.JID).Jugador.ID)
		// que pasa si era el ultimo que se esperaba que cantara flor?
		// tengo que Hacer el Eval de la flor
		todosLosJugadoresConFlorCantaron := len(p.Ronda.Envite.SinCantar) == 0
		if todosLosJugadoresConFlorCantaron {
			florPkts2 := evalFlor(p)
			if p.Verbose {
				pkts2 = append(pkts2, florPkts2...)
			}
		}
	}

	// era el ultimo en tirar de esta mano?
	eraElUltimoEnTirar := p.Ronda.GetSigHabilitado(*p.Manojo(jugada.JID)) == nil

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
			p.Ronda.Envite.SinCantar = []string{}
			e.Puntaje = totalPts

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.ALL,
					enco.SumaPts{
						Autor:  e.CantadoPor,
						Razon:  enco.EnviteNoQuerido,
						Puntos: totalPts,
					},
				))
			}

			p.SumarPuntos(p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo, totalPts)

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
				esDelEquipoContrario := p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo != p.Manojo(jugada.JID).Jugador.Equipo
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
			p.Ronda.Envite.SinCantar = []string{}

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.ALL,
					enco.SumaPts{
						Autor:  p.Ronda.Envite.CantadoPor,
						Razon:  enco.FlorAchicada,
						Puntos: totalPts,
					},
				))
			}

			p.SumarPuntos(p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo, totalPts)

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
		empiezaNuevaRonda, res2 := p.EvaluarMano()

		if p.Verbose {
			pkts2 = append(pkts2, res2...)
		}

		if !empiezaNuevaRonda {

			// esta parte no tiene sentido: si se fue al mazo se sabe que va a
			// empezar una nueva ronda. Este `if` es codigo muerto

			// actualizo el mano
			p.Ronda.ManoEnJuego++
			p.Ronda.SetNextTurnoPosMano()
			// lo envio

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.ALL,
					enco.SigTurnoPosMano(p.Ronda.Turno),
				))
			}

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

				if p.Verbose {
					for _, m := range p.Ronda.Manojos {
						pkts2 = append(pkts2, enco.Pkt(
							enco.ALL,
							enco.NuevaRonda{
								Perspectiva: p.PerspectivaCacheFlor(&m),
							},
						))
					}
				}

			} // else {
			// p.byeBye()
			// }

		}
	} else {
		// cambio de turno solo si era su turno
		eraSuTurno := p.Ronda.GetElTurno().Jugador.ID == p.Manojo(jugada.JID).Jugador.ID
		if eraSuTurno {
			p.Ronda.SetNextTurno()

			if p.Verbose {
				pkts2 = append(pkts2, enco.Pkt(
					enco.ALL,
					enco.SigTurno(p.Ronda.Turno),
				))
			}

		}
	}

	return pkts2
}
