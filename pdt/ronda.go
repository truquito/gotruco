package pdt

import (
	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/util"
)

// Ronda :
type Ronda struct {
	ManoEnJuego          NumMano        `json:"manoEnJuego"`
	CantJugadoresEnJuego map[Equipo]int `json:"cantJugadoresEnJuego"`

	/* Indices */
	ElMano JIX `json:"elMano"`
	Turno  JIX `json:"turno"`

	/* toques, gritos y cantos */
	Envite Envite `json:"envite"`
	Truco  Truco  `json:"truco"`

	/* cartas */
	Manojos []Manojo       `json:"manojos"`
	MIXS    map[string]int `json:"mixs,omitempty"` // index de manojos
	Muestra Carta          `json:"muestra"`

	Manos []Mano `json:"manos"`
}

/* GETERS */

// El remplazo al antiguo p.Ronda.Manojo["pepe"] .
func (r Ronda) Manojo(jid string) *Manojo {
	jix, ok := r.MIXS[jid]
	if !ok {
		return nil
	}
	return &r.Manojos[jix]
}

func (r Ronda) JIX(jid string) int {
	return int(r.MIXS[jid])
}

// GetElMano .
func (r Ronda) GetElMano() *Manojo {
	return &r.Manojos[r.ElMano]
}

// GetSigElMano retorna el id del que deberia ser el siguiente mano
func (r Ronda) GetSigElMano() JIX {
	return JIX(r.JIX(r.GetSiguiente(*r.GetElMano()).Jugador.ID))
}

// GetElTurno .
func (r Ronda) GetElTurno() *Manojo {
	return &r.Manojos[r.Turno]
}

// GetManoAnterior .
func (r Ronda) GetManoAnterior() *Mano {
	return &r.Manos[r.ManoEnJuego-1]
}

// GetManoActual .
func (r Ronda) GetManoActual() *Mano {
	return &r.Manos[r.ManoEnJuego]
}

// GetIdx retorna el indice de un manojo
func (r Ronda) GetIdx(m *Manojo) int {
	// v3
	return r.MIXS[m.Jugador.ID]

	// v2
	// return m.Jugador.Jix

	// v1
	// var idx int
	// cantJugadores := len(r.Manojos)
	// for idx = 0; idx < cantJugadores; idx++ {
	// 	esEse := r.Manojos[idx].Jugador.ID == m.Jugador.ID
	// 	if esEse {
	// 		break
	// 	}
	// }
	// return idx
}

// getSig devuelve el `JugadorIdx` del
// jugador siguiente a j
func (r *Ronda) getSig(j JIX) JIX {
	cantJugadores := len(r.Manojos)
	esElUltimo := int(j) == cantJugadores-1
	if esElUltimo {
		return 0
	}
	return j + 1
}

// GetSiguiente .
// usar sigHabilitado; esta es mas para uso interno
// porque no necesariamente el manojo esta hablilitado
// eg porque se fue al mazo
// getSiguiente devuelve el puntero al manojo que le sigue
func (r Ronda) GetSiguiente(m Manojo) *Manojo {
	idx := r.JIX(m.Jugador.ID)
	cantJugadores := len(r.Manojos)
	esElUltimo := idx == cantJugadores-1
	if esElUltimo {
		return &r.Manojos[0]
	}
	return &r.Manojos[idx+1]
}

// a diferencia de `siguiente`, retorna un puntero al siguiente manojo
// con respecto a `m` que todavia esta en juego (que no se fue al mazo)
// retorna nil si no existe

// GetSigHabilitado .
// no era el ultimo si todavia queda al menos uno
// que viene despues de el que todavia no se fue al mazo
// y todavia no tiro carta en esta mano
// o bien: era el ultimo sii el siguiente de el era el mano
func (r Ronda) GetSigHabilitado(m Manojo) *Manojo {
	var sig *Manojo = &m
	var i int
	cantJugadores := len(r.Manojos)

	// como maximo voy a dar la vuelta entera
	for i = 0; i < cantJugadores; i++ {
		sig = r.GetSiguiente(*sig)
		// checkeos
		noSeFueAlMazo := !sig.SeFueAlMazo
		yaTiroCartaEnEstaMano := sig.yaTiroCarta(r.ManoEnJuego)
		noEsEl := sig.Jugador.ID != m.Jugador.ID
		ok := noSeFueAlMazo && !yaTiroCartaEnEstaMano && noEsEl
		if ok {
			break
		}
	}

	if i == cantJugadores {
		return nil
	}

	return sig
}

// retorna todos los manojos que tienen flor
func (r Ronda) GetFlores() (hayFlor bool, manojosConFlor []*Manojo) {
	for i, manojo := range r.Manojos {
		tieneFlor, _ := manojo.TieneFlor(r.Muestra)
		if tieneFlor {
			manojosConFlor = append(manojosConFlor, &r.Manojos[i])
		}
	}
	hayFlor = len(manojosConFlor) > 0
	return hayFlor, manojosConFlor
}

// GetLaFlorMasAlta retorna el manojo con la flor mas alta en la ronda
// y su valor
// pre-requisito: hay flor en la ronda
func (r *Ronda) GetLaFlorMasAlta() (*Manojo, int) {
	var (
		maxFlor     = -1
		maxIdx  int = -1
	)
	for i := range r.Manojos {
		valorFlor, _ := r.Manojos[i].CalcFlor(r.Muestra)
		if valorFlor > maxFlor {
			maxFlor = valorFlor
			maxIdx = i
		}
	}
	return &r.Manojos[maxIdx], maxFlor
}

func (r *Ronda) getManojo(jIdx JIX) *Manojo {
	return &r.Manojos[jIdx]
}

/* PREDICADOS */

// ver documentacion cambio de variable
func cv(x, mano JIX, cantJugadores int) (y JIX) {
	if x >= mano {
		y = x - mano
	} else {
		c := JIX(cantJugadores) - mano
		y = x + c
	}
	return y
}

// leGanaDeMano devuelve `true` sii
// `i` "le gana de mano" a `j`
func (r Ronda) leGanaDeMano(i, j JIX) bool {
	cantJugadores := len(r.Manojos)
	// cambios de variables
	p := cv(i, r.ElMano, cantJugadores)
	q := cv(j, r.ElMano, cantJugadores)
	return p < q
}

func (r *Ronda) hayEquipoSinCantar(equipo Equipo) bool {
	for _, jid := range r.Envite.SinCantar {
		mismoEquipo := r.Manojo(jid).Jugador.Equipo == equipo
		if mismoEquipo {
			return true
		}
	}

	return false
}

/* SETTERS */

// SetNextTurno este metodo es inseguro ya que manojoSigTurno puede ser nil
func (r *Ronda) SetNextTurno() {
	manojoTurnoActual := r.Manojos[r.Turno]
	manojoSigTurno := r.GetSigHabilitado(manojoTurnoActual)
	r.Turno = JIX(r.JIX(manojoSigTurno.Jugador.ID))
}

// nextTurnoPosMano: setea el turno siguiente *segun el resultado de
// la mano anterior*

// PARA USAR ESTA FUNCION ANTES SE DEBE INCREMENTEAR
// (o actualizar en caso de empezar una ronda nueva)
// EL VALOR DE r.manoEnJuego

// que pasa cuando el ganador de una mano se habia ido al mazo?
// no se tiene que poder:
// si en esta mano ya jugaste carta -> no te podes ir al mazo
// o bien: solo te podes ir al mazo cuando es tu turno
// luego este metodo es correcto

// SetNextTurnoPosMano ..
func (r *Ronda) SetNextTurnoPosMano() {
	// checkeo: si justo el nuevo turno, resulta que se fue al mazo
	// entonces elijo el primero que encuentre que sea de su mismo equipo
	// que no se haya ido al mazo:
	safety_check := func() {
		candidato := r.Manojos[r.Turno]
		if candidato.SeFueAlMazo {
			n := len(r.Manojos)
			start_from := int(r.ElMano)
			for i := 0; i < n; i++ {
				ix := util.Mod(start_from+i, n)
				m := &r.Manojos[ix]
				mismoEquipo := m.Jugador.Equipo == candidato.Jugador.Equipo
				if mismoEquipo && !m.SeFueAlMazo {
					r.Turno = JIX(r.JIX(m.Jugador.ID))
					break
				}
			}
		}
	}

	// si es la Primera mano que se juega
	// entonces es el turno del mano
	if r.ManoEnJuego == Primera {
		r.Turno = r.ElMano
		// si no, es turno del ganador de
		// la mano anterior
		safety_check()
	} else {
		// solo si la mano anterior no fue parda
		// si fue parda busco la que empardo mas cercano al mano
		if r.GetManoAnterior().Resultado != Empardada {
			r.Turno = JIX(r.JIX(r.Manojo(r.GetManoAnterior().Ganador).Jugador.ID))
			safety_check()
		} else {
			// 1. obtengo la carta de maximo valor de la mano anterior
			// 2. busco a partir de la mano quien es el primero en tener
			//    esa carta y que no se haya ido al mazo aun
			// 3. si todos los que empardaron ya se fueron, entonces hago la
			//    vieja confiable (302)
			max := -1
			for _, tirada := range r.GetManoAnterior().CartasTiradas {
				poder := tirada.Carta.CalcPoder(r.Muestra)
				if poder > max {
					max = poder
				}
			}
			for _, tirada := range r.GetManoAnterior().CartasTiradas {
				poder := tirada.Carta.CalcPoder(r.Muestra)
				if poder == max {
					if !r.Manojo(tirada.Jugador).SeFueAlMazo {
						r.Turno = JIX(r.JIX(r.Manojo(tirada.Jugador).Jugador.ID))
						safety_check()
						return
					}
				}
			}
			// si llegue aca es porque los vejigas que empardaron se fueron
			// entonces agarro al primero a partir del mano que aun
			// no se haya ido
			r.Turno = JIX(r.JIX(r.GetSigHabilitado(*r.GetElMano()).Jugador.ID))
			safety_check()
		}
	}

}

// SetManojos .
func (r *Ronda) SetManojos(manojos []Manojo) {
	// cargo los manojos
	for m, manojo := range manojos {
		copy(r.Manojos[m].Cartas[:], manojo.Cartas[:])
		// for c, carta := range manojo.Cartas {
		// 	r.Manojos[m].Cartas[c] = carta
		// }
	}
	// flores
	r.CachearFlores(true)
}

// SetMuestra .
func (r *Ronda) SetMuestra(muestra Carta) {
	r.Muestra = muestra
	r.CachearFlores(true)
}

/* EDITORES */

// todo: esto anda bien; es legacy; pero hacer que devuelva punteros
// no indices

// ExecElEnvido computa el envido de la ronda
/**
* @return `jIdx JugadorIdx` Es el indice del jugador con
* el envido mas alto (i.e., ganador)
* @return `max int` Es el valor numerico del maximo envido
* @return `pkts []*Packet` Es conjunto ordenado de todos
* los mensajes que normalmente se escucharian en una ronda
* de envido en la vida real.
* e.g.:
* 	`pkts[0] = Jacinto dice: "tengo 9"`
*   `pkts[1] = Patricio dice: "son buenas" (tenia 3)`
*   `pkts[2] = Pedro dice: "30 son mejores!"`
*		`pkts[3] = Juan dice: "33 son mejores!"`
*
 */
func (r *Ronda) ExecElEnvido(verbose bool) (jIdx JIX, max int, pkts2 []enco.Envelope) {

	cantJugadores := len(r.Manojos)

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
		envidos[i] = r.Manojos[i].CalcularEnvido(r.Muestra)
	}

	// `yaDijeron` indica que jugador ya "dijo"
	// si tenia mejor, o peor envido. Por lo tanto,
	// ya no es tenido en cuenta.
	yaDijeron := make([]bool, cantJugadores)
	// `jIdx` indica el jugador con el envido mas alto
	// var jIdx JugadorIdx

	// empieza el mas cercano a la mano (inclusive), el que no se haya ido aun
	jIdx = r.ElMano
	for r.Manojos[jIdx].SeFueAlMazo {
		jIdx++
		if int(jIdx) == cantJugadores {
			jIdx = 0
		}
	}

	yaDijeron[jIdx] = true

	if verbose {
		pkts2 = append(pkts2, enco.Env(
			enco.ALL,
			enco.DiceTengo{
				Autor: r.Manojos[jIdx].Jugador.ID,
				Valor: envidos[jIdx],
			},
		))
	}

	// `todaviaNoDijeronSonMejores` se usa para
	// no andar repitiendo "son bueanas" "son buenas"
	// por cada jugador que haya jugado "de callado"
	// y ahora resulte tener peor envido.
	// agiliza el juego, de forma tal que solo se
	// escucha decir "xx son mejores", "yy son mejores"
	// "zz son mejores"
	todaviaNoDijeronSonMejores := true

	// iterador
	i := r.ElMano + 1

	// fix el mano es el ultimo
	if int(r.ElMano) == cantJugadores-1 {
		i = 0
	}

	// termina el bucle cuando se haya dado
	// "una vuelta completa" de:mano+1 hasta:mano
	// ergo, cuando se "resetea" el iterador,
	// se setea a `p.Ronda.elMano + 1`
	for i != r.ElMano {
		seFueAlMazo := r.Manojos[i].SeFueAlMazo
		todaviaEsTenidoEnCuenta := !yaDijeron[i] && !seFueAlMazo
		if todaviaEsTenidoEnCuenta {

			esDeEquipoContrario := r.Manojos[i].Jugador.Equipo != r.Manojos[jIdx].Jugador.Equipo
			tieneEnvidoMasAlto := envidos[i] > envidos[jIdx]
			tieneEnvidoIgual := envidos[i] == envidos[jIdx]
			leGanaDeMano := r.leGanaDeMano(i, jIdx)
			sonMejores := tieneEnvidoMasAlto || (tieneEnvidoIgual && leGanaDeMano)

			if sonMejores {
				if esDeEquipoContrario {

					if verbose {
						pkts2 = append(pkts2, enco.Env(
							enco.ALL,
							enco.DiceSonMejores{
								Autor: r.Manojos[i].Jugador.ID,
								Valor: envidos[i],
							},
						))
					}

					jIdx = i
					yaDijeron[i] = true
					todaviaNoDijeronSonMejores = false
					// se "resetea" el bucle
					i = r.getSig(r.ElMano)

				} else /* es del mismo equipo */ {
					// no dice nada si es del mismo equipo
					// juega de callado & sigue siendo tenido
					// en cuenta
					i = r.getSig(i)
				}

			} else /* tiene el envido mas chico */ {
				if esDeEquipoContrario {
					if todaviaNoDijeronSonMejores {

						if verbose {
							pkts2 = append(pkts2, enco.Env(
								enco.ALL,
								enco.DiceSonBuenas(r.Manojos[i].Jugador.ID),
								// valor de su envido es `envidos[i]` pero no corresponde decirlo
							))
						}

						yaDijeron[i] = true
						// pasa al siguiente
					}
					i = r.getSig(i)
				} else {
					// es del mismo equipo pero tiene un envido
					// mas bajo del que ya canto su equipo.
					// ya no lo tengo en cuenta, pero no dice nada.
					yaDijeron[i] = true
					i = r.getSig(i)
				}
			}

		} else {
			// si no es tenido en cuenta,
			// simplemente pasar al siguiente
			i = r.getSig(i)
		}
	} // fin bucle while

	max = envidos[jIdx]

	return jIdx, max, pkts2
}

// ExecLaFlores computa los cantos de la flor
/**
* @return `j *Manojo` Es el ptr al manojo con
* la flor mas alta (i.e., ganador)
* @return `max int` Es el valor numerico de la flor mas alta
* @return `pkts []*Packet` Es conjunto ordenado de todos
* los mensajes que normalmente se escucharian en una ronda
* de flor en la vida real empezando desde jIdx
* e.g.:
* 	`pkts[0] = Jacinto dice: "tengo 9"`
*   `pkts[1] = Patricio dice: "son buenas" (tenia 3)`
*   `pkts[2] = Pedro dice: "30 son mejores!"`
*	`pkts[3] = Juan dice: "33 son mejores!"`
*
 */
func (r *Ronda) ExecLaFlores(aPartirDe JIX, verbose bool) (j *Manojo, max int, pkts2 []enco.Envelope) {

	// si solo un equipo tiene flor, entonces se saltea esta parte
	soloUnEquipoTieneFlores := true
	equipo := r.Envite.JugadoresConFlor[0].Jugador.Equipo
	for _, m := range r.Envite.JugadoresConFlor[1:] {
		if m.Jugador.Equipo != equipo {
			soloUnEquipoTieneFlores = false
			break
		}
	}
	if soloUnEquipoTieneFlores {
		return r.Envite.JugadoresConFlor[0], 0, nil
	}

	cantJugadores := len(r.Manojos)

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

	for i := range r.Manojos {
		flores[i], _ = r.Manojos[i].CalcFlor(r.Muestra)
		tieneFlor := flores[i] > 0
		seFueAlMazo := r.Manojos[i].SeFueAlMazo
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

		if verbose {
			pkts2 = append(pkts2, enco.Env(
				enco.ALL,
				enco.DiceTengo{
					Autor: r.Manojos[aPartirDe].Jugador.ID,
					Valor: flores[aPartirDe],
				},
			))
		}

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
	i := r.getSig(aPartirDe)

	// termina el bucle cuando se haya dado
	// "una vuelta completa" de:aPartirDe hasta:aPartirDe
	// ergo, cuando se "resetea" el iterador,
	for i != aPartirDe {
		todaviaEsTenidoEnCuenta := !yaDijeron[i]
		if todaviaEsTenidoEnCuenta {

			esDeEquipoContrario := r.Manojos[i].Jugador.Equipo != r.Manojos[jIdx].Jugador.Equipo
			tieneEnvidoMasAlto := flores[i] > flores[jIdx]
			tieneEnvidoIgual := flores[i] == flores[jIdx]
			leGanaDeMano := r.leGanaDeMano(i, jIdx)
			sonMejores := tieneEnvidoMasAlto || (tieneEnvidoIgual && leGanaDeMano)

			if sonMejores {
				if esDeEquipoContrario {

					if verbose {
						pkts2 = append(pkts2, enco.Env(
							enco.ALL,
							enco.DiceSonMejores{
								Autor: r.Manojos[i].Jugador.ID,
								Valor: flores[i],
							},
						))
					}

					jIdx = i
					yaDijeron[i] = true
					todaviaNoDijeronSonMejores = false
					// se "resetea" el bucle
					i = r.getSig(aPartirDe)

				} else /* es del mismo equipo */ {
					// no dice nada si es del mismo equipo
					// juega de callado & sigue siendo tenido
					// en cuenta
					i = r.getSig(i)
				}

			} else /* tiene el envido mas chico */ {
				if esDeEquipoContrario {
					if todaviaNoDijeronSonMejores {

						if verbose {
							pkts2 = append(pkts2, enco.Env(
								enco.ALL,
								enco.DiceSonBuenas(r.Manojos[i].Jugador.ID),
							))
						}

						yaDijeron[i] = true
						// pasa al siguiente
					}
					i = r.getSig(i)
				} else {
					// es del mismo equipo pero tiene un envido
					// mas bajo del que ya canto su equipo.
					// ya no lo tengo en cuenta, pero no dice nada.
					yaDijeron[i] = true
					i = r.getSig(i)
				}
			}

		} else {
			// si no es tenido en cuenta,
			// simplemente pasar al siguiente
			i = r.getSig(i)
		}
	}

	max = flores[jIdx]

	return r.getManojo(jIdx), max, pkts2
}

// los anteriores a `aPartirDe` (incluido este) no
// son necesarios de checkear porque ya han sido
// checkeados si tenian flor
// func (r Ronda) cantarFloresSiLasHay(aPartirDe JugadorIdx) {
// 	for _, jugador := range r.Manojos[aPartirDe+1:] {
// 		tieneFlor, tipoFlor := jugador.tieneFlor(r.Muestra)
// 		if tieneFlor {
// 			// todo:
// 			tieneFlor = false
// 			tipoFlor = tipoFlor + 1
// 			// var jugada IJugada = responderNoQuiero{}
// 			// jugador.cantarFlor(tipoFlor, r.muestra)
// 			r.Envite.Estado = DESHABILITADO
// 			break
// 		}
// 	}
// }

/* INICIALIZADORES */

func (r *Ronda) CachearFlores(reset bool) {
	// flores
	_, JugadoresConFlor := r.GetFlores()
	r.Envite.JugadoresConFlor = JugadoresConFlor

	if reset {
		conFlor := make([]string, 0)
		for _, m := range JugadoresConFlor {
			conFlor = append(conFlor, m.Jugador.ID)
		}
		r.Envite.SinCantar = conFlor
	}
}

/*
 * Reparte 3 cartas al azar a cada manojo de c/jugador
 * y 1 a la `muestra` (se las actualiza)
 */
func (r *Ronda) repartirCartas() {
	cantJugadores := len(r.Manojos)
	// genero `3*cantJugadores + 1` cartas al azar
	randomCards := getCartasRandom(3*cantJugadores + 1)

	for idxJugador := 0; idxJugador < cantJugadores; idxJugador++ {
		for idxCarta := 0; idxCarta < 3; idxCarta++ {
			cartaID := CartaID(randomCards[3*idxJugador+idxCarta])
			carta := nuevaCarta(cartaID)
			r.Manojos[idxJugador].Cartas[idxCarta] = &carta
			r.Manojos[idxJugador].Tiradas[idxCarta] = false
		}
	}

	// la ultima es la muestra
	n := cap(randomCards)
	r.Muestra = nuevaCarta(CartaID(randomCards[n-1]))
}

func (r *Ronda) indexarManojos() {
	// indexo los manojos
	for i := range r.Manojos {
		jid := r.Manojos[i].Jugador.ID
		r.MIXS[jid] = i
	}
}

// resetea una ronda
func (r *Ronda) Reset(elMano JIX) {
	cantJugadores := len(r.Manojos)
	cantJugadoresPorEquipo := cantJugadores / 2

	r.ManoEnJuego = Primera
	r.CantJugadoresEnJuego[Rojo] = cantJugadoresPorEquipo
	r.CantJugadoresEnJuego[Azul] = cantJugadoresPorEquipo
	r.ElMano = elMano
	r.Turno = elMano
	r.Envite = Envite{Estado: NOCANTADOAUN, Puntaje: 0}
	r.Truco = Truco{CantadoPor: "", Estado: NOGRITADOAUN}
	r.Manos = make([]Mano, 3)

	for i := range r.Manojos {
		r.Manojos[i].UltimaTirada = -1
		r.Manojos[i].SeFueAlMazo = false
		// r.Manojos[i].Cartas
		r.Manojos[i].Tiradas = [cantCartasManojo]bool{false, false, false}
	}
}

func (r *Ronda) reiniciar(elMano JIX) {
	r.Reset(elMano)

	// reparto 3 cartas al azar a cada jugador
	// y ademas una muestra, tambien al azar.
	r.repartirCartas()

	// flores
	r.CachearFlores(true)
}

/* CONSTRUCTOR */

// Reserva el espacio en memoria
// Por default, cuando se crea una ronda el mano sera el jix = 0
func NuevaRonda(equipoAzul, equipoRojo []string) Ronda {

	cantJugadores := len(equipoAzul) * 2
	cantJugadoresPorEquipo := len(equipoAzul)

	ronda := Ronda{
		ManoEnJuego: Primera,
		CantJugadoresEnJuego: map[Equipo]int{
			Rojo: cantJugadoresPorEquipo,
			Azul: cantJugadoresPorEquipo,
		},
		ElMano:  0,
		Turno:   0,
		Envite:  Envite{Estado: NOCANTADOAUN, Puntaje: 0},
		Truco:   Truco{CantadoPor: "", Estado: NOGRITADOAUN},
		Manojos: make([]Manojo, cantJugadores),
		MIXS:    make(map[string]int),
		Manos:   make([]Mano, 3),
	}

	for i := range ronda.Manos {
		// como maximo pueden tirar 1 cartas cada jugador
		ronda.Manos[i].CartasTiradas = make([]CartaTirada, 0, cantJugadores)
	}

	for i := 0; i < cantJugadoresPorEquipo; i++ {
		ix := i << 1
		ronda.Manojos[ix].UltimaTirada = -1
		ronda.Manojos[ix].Jugador = &Jugador{equipoAzul[i], Azul}
		ronda.Manojos[ix+1].Jugador = &Jugador{equipoRojo[i], Rojo}
		ronda.Manojos[ix+1].UltimaTirada = -1
	}

	ronda.indexarManojos()

	// reparto 3 cartas al azar a cada jugador
	// y ademas una muestra, tambien al azar.
	ronda.repartirCartas()

	// flores
	ronda.CachearFlores(true)

	return ronda
}
