package pdt

import "github.com/truquito/truco/util"

func ElEnvidoEsRespondible(p *Partida) bool {
	return p.Ronda.Envite.Estado >= ENVIDO && p.Ronda.Envite.Estado <= FALTAENVIDO
}

func LaFlorEsRespondible(p *Partida) bool {
	// hay manojos con flor que no cantaron aun
	enviteNoDeshabilitado := p.Ronda.Envite.Estado != DESHABILITADO
	hayManojosSinCantar := len(p.Ronda.Envite.SinCantar) > 0
	florEnJuego := p.Ronda.Envite.Estado >= FLOR && p.Ronda.Envite.Estado <= CONTRAFLORALRESTO
	if enviteNoDeshabilitado && (hayManojosSinCantar || florEnJuego) {
		return true
	}

	return false
}

func ElTrucoRespondible(p *Partida) bool {
	return p.Ronda.Truco.Estado == TRUCO ||
		p.Ronda.Truco.Estado == RETRUCO ||
		p.Ronda.Truco.Estado == VALE4
}

// dado un manojo A, retorna el manojo del equipo contrario a A que:
//  mas alejado se encuentre del mano y no se haya ido al mazo aun
func Encontrar_respondedor_mas_lejano(p *Partida, proponente *Manojo) (respondedor *Manojo) {
	// itero hasta un maximo de la cantidad de jugadores
	// cual es el absolute_index del mano?
	maix := p.Ronda.ElMano
	n := len(p.Ronda.Manojos)
	// entonces empiezo a iterar desde
	start_from := util.Mod(int(maix)-1, n)

	for i := 0; i < n; i++ {
		ix := util.Mod(start_from-i, n)
		m := &p.Ronda.Manojos[ix]
		esEquipoOpuesto := m.Jugador.Equipo != proponente.Jugador.Equipo
		if esEquipoOpuesto && !m.SeFueAlMazo {
			return m
		}
	}

	panic("no se encontro jugador mas lejano")
}

// PRE: la flor es respondible o es iniciable
func Encontrar_siguiente_florero(p *Partida) (respondedor *Manojo) {
	// itero hasta un maximo de la cantidad de jugadores
	// Min Absolute Index
	n := len(p.Ronda.Manojos)

	// Luego, o bien:
	// 	1. de los que tienen flor, nadie canto aun, entonces
	// 		 el que debe cantar la flor es el primero a partir del mano
	// 		 que no se haya ido al mazo, y tengo flor.
	//  2. o bien alguien ya canto algo, entonces el rho es el siguiente mas
	// 		 proximo (con flor).
	start_from := int(p.Ronda.ElMano)

	laFlorEstaEnJuego := p.Ronda.Envite.Estado >= FLOR
	if laFlorEstaEnJuego {
		start_from = p.Ronda.MIXS[p.Ronda.Envite.CantadoPor]
	}

	for i := 0; i < n; i++ {
		ix := util.Mod(start_from+i, n)
		m := &p.Ronda.Manojos[ix]

		// si se fue al mazo, lo ignoro
		if m.SeFueAlMazo {
			continue
		}

		tieneFlor, _ := m.TieneFlor(p.Ronda.Muestra)
		if !tieneFlor {
			continue
		}

		// puede hacer algo relacionado con la flor?
		js := []IJugada{
			CantarFlor{JID: m.Jugador.ID},
			CantarContraFlor{JID: m.Jugador.ID},
			CantarContraFlorAlResto{JID: m.Jugador.ID},
			// { CantarConFlorMeAchico{JID: m.Jugador.ID}, enco.new },
			// respuestas
			ResponderQuiero{JID: m.Jugador.ID},
			ResponderNoQuiero{JID: m.Jugador.ID},
		}

		for _, j := range js {
			_, ok := j.Ok(p)
			if ok {
				return m
			}
		}

	}

	panic("no se encontro florero")
}

// retorna el "turno" de TEORIA DE JUEGOS. Es decir, la funcion Rho(.)
// si hay algo respondible, retorna el manojo del equipo respondedor mas alejado
// del proponente (i.e., es el turno del "pie")
// notar que, en caso de tener que responder el envite,
// esto se puede abstraer porque da igual si debe responder 1 o 2
func Rho(p *Partida) (m *Manojo) {
	if LaFlorEsRespondible(p) {
		return Encontrar_siguiente_florero(p)
	}

	if ElEnvidoEsRespondible(p) {
		// proponente := p.Ronda.Manojo[p.Ronda.Envite.CantadoPor]
		proponente := p.Ronda.Manojo(p.Ronda.Envite.CantadoPor)
		return Encontrar_respondedor_mas_lejano(p, proponente)
	}

	if ElTrucoRespondible(p) {
		// proponente := p.Ronda.Manojo[p.Ronda.Truco.CantadoPor]
		proponente := p.Ronda.Manojo(p.Ronda.Truco.CantadoPor)
		return Encontrar_respondedor_mas_lejano(p, proponente)
	}

	return p.Ronda.GetElTurno()
}
