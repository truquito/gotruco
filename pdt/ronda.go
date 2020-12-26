package pdt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/filevich/truco/out"
)

// EstadoTruco : enum
type EstadoTruco int

// enums del truco
const (
	NOCANTADO EstadoTruco = iota
	TRUCO
	TRUCOQUERIDO
	RETRUCO
	RETRUCOQUERIDO
	VALE4
	VALE4QUERIDO
)

var toEstadoTruco = map[string]EstadoTruco{
	"noCantado":      NOCANTADO,
	"truco":          TRUCO,
	"trucoQuerido":   TRUCOQUERIDO,
	"reTruco":        RETRUCO,
	"reTrucoQuerido": RETRUCOQUERIDO,
	"vale4":          VALE4,
	"vale4Querido":   VALE4QUERIDO,
}

// toString
func (e EstadoTruco) String() string {
	estados := []string{
		"noCantado",
		"truco",
		"trucoQuerido",
		"reTruco",
		"reTrucoQuerido",
		"vale4",
		"vale4Querido",
	}

	ok := e >= 0 || int(e) < len(toEstadoTruco)
	if !ok {
		return "Unknown"
	}

	return estados[e]
}

// MarshalJSON marshals the enum as a quoted json string
func (e EstadoTruco) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(e.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *EstadoTruco) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*e = toEstadoTruco[j]
	return nil
}

// Truco :
type Truco struct {
	CantadoPor *Manojo     `json:"cantadoPor"`
	Estado     EstadoTruco `json:"estado"`
}

// Ronda :
type Ronda struct {
	ManoEnJuego          NumMano        `json:"manoEnJuego"`
	CantJugadoresEnJuego map[Equipo]int `json:"cantJugadoresEnJuego"`

	/* Indices */
	ElMano JugadorIdx    `json:"elMano"`
	Turno  JugadorIdx    `json:"turno"`
	Pies   [2]JugadorIdx `json:"pies"`

	/* toques, gritos y cantos */
	Envite Envite `json:"envite"`
	Truco  Truco  `json:"truco"`

	/* cartas */
	Manojos []Manojo `json:"manojos"`
	Muestra Carta    `json:"muestra"`

	Manos []Mano `json:"manos"`
}

/* GETERS */

// GetElMano .
func (r Ronda) GetElMano() *Manojo {
	return &r.Manojos[r.ElMano]
}

// GetSigElMano retorna el id del que deberia ser el siguiente mano
func (r Ronda) GetSigElMano() JugadorIdx {
	return JugadorIdx(r.GetIdx(*r.GetSiguiente(*r.GetElMano())))
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
func (r Ronda) GetIdx(m Manojo) int {
	var idx int
	cantJugadores := len(r.Manojos)
	for idx = 0; idx < cantJugadores; idx++ {
		esEse := r.Manojos[idx].Jugador.ID == m.Jugador.ID
		if esEse {
			break
		}
	}
	return idx
}

// getSig devuelve el `JugadorIdx` del
// jugador siguiente a j
func (r *Ronda) getSig(j JugadorIdx) JugadorIdx {
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
	idx := r.GetIdx(m)
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
		noSeFueAlMazo := sig.SeFueAlMazo == false
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
func (r Ronda) getFlores() (hayFlor bool, manojosConFlor []*Manojo) {
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
		valorFlor, _ := r.Manojos[i].calcFlor(r.Muestra)
		if valorFlor > maxFlor {
			maxFlor = valorFlor
			maxIdx = i
		}
	}
	return &r.Manojos[maxIdx], maxFlor
}

func (r *Ronda) getManojo(jIdx JugadorIdx) *Manojo {
	return &r.Manojos[jIdx]
}

// GetManojoByStr ..
// OJO QUE AHORA LAS COMPARACIONES SON CASE INSENSITIVE
// ENTONCES SI EL IDENTIFICADOR Juan == jUaN
// ojo con los kakeos
// todo: esto es ineficiente
// getManojo devuelve el puntero al manojo,
// dado un string que identifique al jugador duenio de ese manojo
func (r *Ronda) GetManojoByStr(idJugador string) (*Manojo, error) {
	idJugador = strings.ToLower(idJugador)
	for i := range r.Manojos {
		idActual := strings.ToLower(r.Manojos[i].Jugador.ID)
		esEse := idActual == idJugador
		if esEse {
			return &r.Manojos[i], nil
		}
	}
	return nil, fmt.Errorf("Jugador `%s` no encontrado", idJugador)
}

/* PREDICADOS */

// leGanaDeMano devuelve `true` sii
// `i` "le gana de mano" a `j`
func (r Ronda) leGanaDeMano(i, j JugadorIdx) bool {
	cantJugadores := len(r.Manojos)
	// cambios de variables
	p := cv(i, r.ElMano, cantJugadores)
	q := cv(j, r.ElMano, cantJugadores)
	return p < q
}

/* SETTERS */

// SetNextTurno este metodo es inseguro ya que manojoSigTurno puede ser nil
func (r *Ronda) SetNextTurno() {
	manojoTurnoActual := r.Manojos[r.Turno]
	manojoSigTurno := r.GetSigHabilitado(manojoTurnoActual)
	r.Turno = JugadorIdx(r.GetIdx(*manojoSigTurno))
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
	// si es la Primera mano que se juega
	// entonces es el turno del mano
	if r.ManoEnJuego == Primera {
		r.Turno = r.ElMano
		// si no, es turno del ganador de
		// la mano anterior
	} else {
		// solo si la mano anterior no fue parda
		// si fue parda busco la que empardo mas cercano al mano
		if r.GetManoAnterior().Resultado != Empardada {
			r.Turno = JugadorIdx(r.GetIdx(*r.GetManoAnterior().Ganador))
		} else {
			// 1. obtengo la carta de maximo valor de la mano anterior
			// 2. busco a partir de la mano quien es el primero en tener
			//    esa carta y que no se haya ido al mazo aun
			// 3. si todos los que empardaron ya se fueron, entonces hago la
			//    vieja confiable (302)
			max := -1
			for _, tirada := range r.GetManoAnterior().CartasTiradas {
				poder := tirada.Carta.calcPoder(r.Muestra)
				if poder > max {
					max = poder
				}
			}
			for _, tirada := range r.GetManoAnterior().CartasTiradas {
				poder := tirada.Carta.calcPoder(r.Muestra)
				if poder == max {
					if !tirada.autor.SeFueAlMazo {
						r.Turno = JugadorIdx(r.GetIdx(*tirada.autor))
						return
					}
				}
			}
			// si llegue aca es porque los vejigas que empardaron se fueron
			// entonces agarro al primero a partir del mano que aun
			// no se haya ido
			r.Turno = JugadorIdx(r.GetIdx(*r.GetSigHabilitado(*r.GetElMano())))
		}
	}
}

// SetManojos .
func (r *Ronda) SetManojos(manojos []Manojo) {
	for m, manojo := range manojos {
		for c, carta := range manojo.Cartas {
			r.Manojos[m].Cartas[c] = carta
		}
	}
	// flores
	r.cachearFlores()
}

// SetMuestra .
func (r *Ronda) SetMuestra(muestra Carta) {
	r.Muestra = muestra
}

/* EDITORES */

// todo: esto anda bien; es legacy; pero hacer que devuelva punteros
// no indices

// ExecElEnvido computa el envido de la ronda
/**
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
func (r *Ronda) ExecElEnvido() (jIdx JugadorIdx, max int, pkts []*out.Packet) {

	var stdOut []string

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

	salida := fmt.Sprintf(`   %s dice: "tengo %v"`, r.Manojos[jIdx].Jugador.Nombre,
		envidos[jIdx])
	stdOut = append(stdOut, salida)

	pkts = append(pkts, out.Pkt(
		out.Dest("ALL"),
		out.Msg(out.DiceTengo, r.Manojos[jIdx].Jugador.Nombre, envidos[jIdx]),
	))

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
		todaviaEsTenidoEnCuenta := !yaDijeron[i]
		if todaviaEsTenidoEnCuenta {

			esDeEquipoContrario := r.Manojos[i].Jugador.Equipo != r.Manojos[jIdx].Jugador.Equipo
			tieneEnvidoMasAlto := envidos[i] > envidos[jIdx]
			tieneEnvidoIgual := envidos[i] == envidos[jIdx]
			leGanaDeMano := r.leGanaDeMano(i, jIdx)
			sonMejores := tieneEnvidoMasAlto || (tieneEnvidoIgual && leGanaDeMano)

			if sonMejores {
				if esDeEquipoContrario {

					salida := fmt.Sprintf(`   %s dice: "%v son mejores!"`,
						r.Manojos[i].Jugador.Nombre, envidos[i])
					stdOut = append(stdOut, salida)

					pkts = append(pkts, out.Pkt(
						out.Dest("ALL"),
						out.Msg(out.DiceSonMejores, r.Manojos[i].Jugador.Nombre, envidos[i]),
					))

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

						salida := fmt.Sprintf(`   %s dice: "son buenas" (tenia %v)`,
							r.Manojos[i].Jugador.Nombre, envidos[i])
						stdOut = append(stdOut, salida)

						pkts = append(pkts, out.Pkt(
							out.Dest("ALL"),
							out.Msg(out.DiceSonBuenas, r.Manojos[i].Jugador.Nombre),
							// valor de su envido es `envidos[i]` pero no corresponde decirlo
						))

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

	return jIdx, max, pkts
}

/**
* execCantarFlores computa la flor
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
func (r *Ronda) execCantarFlores(aPartirDe JugadorIdx) (j *Manojo, max int, pkts []*out.Packet) {

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
		flores[i], _ = r.Manojos[i].calcFlor(r.Muestra)
		tieneFlor := flores[i] > 0
		seFueAlMazo := r.Manojos[i].SeFueAlMazo == false
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

		pkts = append(pkts, out.Pkt(
			out.Dest("ALL"),
			out.Msg(out.DiceTengo, r.Manojos[aPartirDe].Jugador.Nombre, flores[aPartirDe]),
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

					pkts = append(pkts, out.Pkt(
						out.Dest("ALL"),
						out.Msg(out.DiceSonMejores, r.Manojos[i].Jugador.Nombre, flores[i]),
					))

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

						pkts = append(pkts, out.Pkt(
							out.Dest("ALL"),
							out.Msg(out.DiceSonBuenas, r.Manojos[i].Jugador.Nombre),
							// valor de su envido es `flores[i]` pero no corresponde decirlo
						))

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

	return r.getManojo(jIdx), max, pkts
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

func (r *Ronda) cachearFlores() {
	// flores
	_, JugadoresConFlor := r.getFlores()
	JugadoresConFlorCopy := make([]*Manojo, len(JugadoresConFlor))
	copy(JugadoresConFlorCopy, JugadoresConFlor)
	r.Envite.JugadoresConFlor = JugadoresConFlor
	r.Envite.JugadoresConFlorQueNoCantaron = JugadoresConFlorCopy
}

func (r *Ronda) singleLinking(jugadores []Jugador) {
	cantJugadores := len(jugadores)
	for i := 0; i < cantJugadores; i++ {
		r.Manojos[i].Jugador = &jugadores[i]
	}
}

/*
 * Reparte 3 cartas al azar a cada manojo de c/jugador
 * y 1 a la `muestra` (se las actualiza)
 */
func (r *Ronda) dealCards() {
	cantJugadores := cap(r.Manojos)
	// genero `3*cantJugadores + 1` cartas al azar
	randomCards := getCartasRandom(3*cantJugadores + 1)

	for idxJugador := 0; idxJugador < cantJugadores; idxJugador++ {
		for idxCarta := 0; idxCarta < 3; idxCarta++ {
			cartaID := CartaID(randomCards[3*idxJugador+idxCarta])
			carta := nuevaCarta(cartaID)
			r.Manojos[idxJugador].Cartas[idxCarta] = &carta
			r.Manojos[idxJugador].CartasNoTiradas[idxCarta] = true
		}
	}

	// la ultima es la muestra
	n := cap(randomCards)
	r.Muestra = nuevaCarta(CartaID(randomCards[n-1]))
}

/* CONSTRUCTOR */

// NuevaRonda : crea una nueva ronda al azar
func NuevaRonda(jugadores []Jugador, elMano JugadorIdx) Ronda {
	cantJugadores := len(jugadores)
	cantJugadoresPorEquipo := cantJugadores / 2
	ronda := Ronda{
		ManoEnJuego: Primera,
		CantJugadoresEnJuego: map[Equipo]int{
			Rojo: cantJugadoresPorEquipo,
			Azul: cantJugadoresPorEquipo,
		},
		ElMano:  elMano,
		Turno:   elMano,
		Envite:  Envite{Estado: NOCANTADOAUN, Puntaje: 0},
		Truco:   Truco{CantadoPor: nil, Estado: NOCANTADO},
		Manojos: make([]Manojo, cantJugadores),
		Manos:   make([]Mano, 3),
	}

	// reparto 3 cartas al azar a cada jugador
	// y ademas una muestra, tambien al azar.
	ronda.dealCards()

	// // hago el SINGLE-linking "jugadores <- manojos"
	ronda.singleLinking(jugadores)

	// flores
	ronda.cachearFlores()

	// p.Ronda.setTurno()

	return ronda
}
