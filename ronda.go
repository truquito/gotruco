package truco

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	ok := e > 0 && int(e) < len(toEstadoTruco)-1
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

type truco struct {
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
	Envido Envido     `json:"envido"`
	Flor   EstadoFlor `json:"flor"`
	Truco  truco      `json:"truco"`

	/* cartas */
	Manojos []Manojo `json:"manojos"`
	Muestra Carta    `json:"muestra"`

	Manos []Mano `json:"manos"`
}

func (r *Ronda) checkFlorDelMano() {
	tieneFlor, tipoFlor := r.getElMano().tieneFlor(r.Muestra)
	if tieneFlor {
		r.getElMano().cantarFlor(tipoFlor, r.Muestra)
	}
}

// devuelve `false` si la ronda se acabo
func (r *Ronda) enJuego() bool {
	return true
}

// los anteriores a `aPartirDe` (incluido este) no
// son necesarios de checkear porque ya han sido
// checkeados si tenian flor
func (r Ronda) cantarFloresSiLasHay(aPartirDe JugadorIdx) {
	for _, jugador := range r.Manojos[aPartirDe+1:] {
		tieneFlor, tipoFlor := jugador.tieneFlor(r.Muestra)
		if tieneFlor {
			// todo:
			tieneFlor = false
			tipoFlor = tipoFlor + 1
			// var jugada IJugada = responderNoQuiero{}
			// jugador.cantarFlor(tipoFlor, r.muestra)
			r.Envido.Estado = DESHABILITADO
			break
		}
	}
}

// retorna todos los manojos que tienen flor
func (r Ronda) getFlores() (hayFlor bool,
	manojosConFlor []*Manojo) {
	for i, manojo := range r.Manojos {
		tieneFlor, _ := manojo.tieneFlor(r.Muestra)
		if tieneFlor {
			manojosConFlor = append(manojosConFlor, &r.Manojos[i])
		}
	}
	hayFlor = len(manojosConFlor) > 0
	return hayFlor, manojosConFlor
}

func (r Ronda) getElMano() *Manojo {
	return &r.Manojos[r.ElMano]
}

// retorna el id del que deberia ser el siguiente mano
func (r Ronda) getSigMano() JugadorIdx {
	return JugadorIdx(r.getIdx(*r.siguiente(*r.getElMano())))
}

func (r Ronda) getElTurno() *Manojo {
	return &r.Manojos[r.Turno]
}

func (r Ronda) getManoAnterior() *Mano {
	return &r.Manos[r.ManoEnJuego-1]
}

func (r Ronda) getManoActual() *Mano {
	return &r.Manos[r.ManoEnJuego]
}

func (r *Ronda) nextTurno() {
	manojoTurnoActual := r.Manojos[r.Turno]
	manojoSigTurno := r.sigHabilitado(manojoTurnoActual)
	r.Turno = JugadorIdx(r.getIdx(*manojoSigTurno))
}

// PARA USAR ESTA FUNCION ANTES SE DEBE INCREMENTEAR
// (o actualizar en caso de empezar una ronda nueva)
// EL VALOR DE r.manoEnJuego
// setea el turno siguiente *segun el resultado de
// la mano anterior*

// que pasa cuando el ganador de una mano se habia ido al mazo?
// no se tiene que poder:
// si en esta mano ya jugaste carta -> no te podes ir al mazo
// o bien: solo te podes ir al mazo cuando es tu turno
// luego este metodo es correcto

func (r *Ronda) nextTurnoPosMano() {
	// si es la primera mano que se juega
	// entonces es el turno del mano
	if r.ManoEnJuego == primera {
		r.Turno = r.ElMano
		// si no, es turno del ganador de
		// la mano anterior
	} else {
		// solo si la mano anterior no fue parda
		// si fue parda el truno se mantiene
		if r.getManoAnterior().Resultado != Empardada {
			r.Turno = JugadorIdx(r.getIdx(*r.getManoAnterior().Ganador))
		}
	}
	fmt.Printf("Es el turno de %s", r.Manojos[r.Turno].Jugador.Nombre)
}

// Print Imprime la informacion de la ronda
func (r Ronda) Print() {
	for i := range r.Manojos {
		fmt.Printf("%s:\n", r.Manojos[i].Jugador.Nombre)
		r.Manojos[i].Print()
	}

	fmt.Printf("\nY la muestra es\n    - %s\n", r.Muestra.toString())
	fmt.Printf("\nEl mano actual es: %s\nEs el turno de %s\n\n",
		r.getElMano().Jugador.Nombre, r.getElTurno().Jugador.Nombre)
}

// sig devuelve el `JugadorIdx` del
// jugador siguiente a j
func (r *Ronda) sig(j JugadorIdx) JugadorIdx {
	cantJugadores := len(r.Manojos)
	esElUltimo := int(j) == cantJugadores-1
	if esElUltimo {
		return 0
	}
	return j + 1
}

// retorna el indice de un manojo
func (r Ronda) getIdx(m Manojo) int {
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

// usar sigHabilitado; esta es mas para uso interno
// porque no necesariamente el manojo esta hablilitado
// eg porque se fue al mazo
// siguiente devuelve el puntero al manojo que le sigue
func (r Ronda) siguiente(m Manojo) *Manojo {
	idx := r.getIdx(m)
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

// no era el ultimo si todavia queda al menos uno
// que viene despues de el que todavia no se fue al mazo
// y todavia no tiro carta en esta mano
// o bien: era el ultimo sii el siguiente de el era el mano
func (r Ronda) sigHabilitado(m Manojo) *Manojo {
	var sig *Manojo = &m
	var i int
	cantJugadores := len(r.Manojos)

	// como maximo voy a dar la vuelta entera
	for i = 0; i < cantJugadores; i++ {
		sig = r.siguiente(*sig)
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

// leGanaDeMano devuelve `true` sii
// `i` "le gana de mano" a `j`
func (r Ronda) leGanaDeMano(i, j JugadorIdx) bool {
	cantJugadores := len(r.Manojos)
	// cambios de variables
	p := cv(i, r.ElMano, cantJugadores)
	q := cv(j, r.ElMano, cantJugadores)
	return p < q
}

// retorna el manojo con la flor mas alta en la ronda
// y su valor
// pre-requisito: hay flor en la ronda
func (r *Ronda) getLaFlorMasAlta() (*Manojo, int) {
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

// todo: esto anda bien; es legacy; pero hacer que devuelva punteros
// no indices
/**
* getElEnvido computa el envido de la ronda
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
func (r *Ronda) getElEnvido() (jIdx JugadorIdx,
	max int, stdOut []string) {

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
		envidos[i] = r.Manojos[i].calcularEnvido(r.Muestra)
	}

	// `yaDijeron` indica que jugador ya "dijo"
	// si tenia mejor, o peor envido. Por lo tanto,
	// ya no es tenido en cuenta.
	yaDijeron := make([]bool, cantJugadores)
	// `jIdx` indica el jugador con el envido mas alto
	// var jIdx JugadorIdx

	// empieza la mano
	jIdx = r.ElMano
	yaDijeron[jIdx] = true
	out := fmt.Sprintf("   %s dice: \"tengo %v\"\n", r.Manojos[jIdx].Jugador.Nombre,
		envidos[jIdx])
	stdOut = append(stdOut, out)

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
					out := fmt.Sprintf("   %s dice: \"%v son mejores!\"\n",
						r.Manojos[i].Jugador.Nombre, envidos[i])
					stdOut = append(stdOut, out)
					jIdx = i
					yaDijeron[i] = true
					todaviaNoDijeronSonMejores = false
					// se "resetea" el bucle
					i = r.sig(r.ElMano)

				} else /* es del mismo equipo */ {
					// no dice nada si es del mismo equipo
					// juega de callado & sigue siendo tenido
					// en cuenta
					i = r.sig(i)
				}

			} else /* tiene el envido mas chico */ {
				if esDeEquipoContrario {
					if todaviaNoDijeronSonMejores {
						out := fmt.Sprintf("   %s dice: \"son buenas\" (tenia %v)\n",
							r.Manojos[i].Jugador.Nombre, envidos[i])
						stdOut = append(stdOut, out)
						yaDijeron[i] = true
						// pasa al siguiente
					}
					i = r.sig(i)
				} else {
					// es del mismo equipo pero tiene un envido
					// mas bajo del que ya canto su equipo.
					// ya no lo tengo en cuenta, pero no dice nada.
					yaDijeron[i] = true
					i = r.sig(i)
				}
			}

		} else {
			// si no es tenido en cuenta,
			// simplemente pasar al siguiente
			i = r.sig(i)
		}
	} // fin bucle while

	max = envidos[jIdx]

	return jIdx, max, stdOut
}

/**
* cantarFlores computa la flor
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
func (r *Ronda) cantarFlores(aPartirDe JugadorIdx) (j *Manojo,
	max int, stdOut []string) {

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
		out := fmt.Sprintf("   %s dice: \"tengo %v\"\n", r.Manojos[aPartirDe].Jugador.Nombre,
			flores[aPartirDe])
		stdOut = append(stdOut, out)
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
	i := r.sig(aPartirDe)

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
					out := fmt.Sprintf("   %s dice: \"%v son mejores!\"\n",
						r.Manojos[i].Jugador.Nombre, flores[i])
					stdOut = append(stdOut, out)
					jIdx = i
					yaDijeron[i] = true
					todaviaNoDijeronSonMejores = false
					// se "resetea" el bucle
					i = r.sig(aPartirDe)

				} else /* es del mismo equipo */ {
					// no dice nada si es del mismo equipo
					// juega de callado & sigue siendo tenido
					// en cuenta
					i = r.sig(i)
				}

			} else /* tiene el envido mas chico */ {
				if esDeEquipoContrario {
					if todaviaNoDijeronSonMejores {
						out := fmt.Sprintf("   %s dice: \"son buenas\" (tenia %v)\n",
							r.Manojos[i].Jugador.Nombre, flores[i])
						stdOut = append(stdOut, out)
						yaDijeron[i] = true
						// pasa al siguiente
					}
					i = r.sig(i)
				} else {
					// es del mismo equipo pero tiene un envido
					// mas bajo del que ya canto su equipo.
					// ya no lo tengo en cuenta, pero no dice nada.
					yaDijeron[i] = true
					i = r.sig(i)
				}
			}

		} else {
			// si no es tenido en cuenta,
			// simplemente pasar al siguiente
			i = r.sig(i)
		}
	}

	max = flores[jIdx]

	return r.getManojo(jIdx), max, stdOut
}

func (r *Ronda) getManojo(jIdx JugadorIdx) *Manojo {
	return &r.Manojos[jIdx]
}

func (r *Ronda) singleLinking(jugadores []Jugador) {
	cantJugadores := len(jugadores)
	for i := 0; i < cantJugadores; i++ {
		r.Manojos[i].Jugador = &jugadores[i]
	}
}

// todo: esto es ineficiente
// getManojo devuelve el puntero al manojo,
// dado un string que identifique al jugador duenio de ese manojo
func (r *Ronda) getManojoByStr(idJugador string) (*Manojo, error) {
	for i := range r.Manojos {
		if r.Manojos[i].Jugador.Nombre == idJugador {
			return &r.Manojos[i], nil
		}
	}
	return nil, fmt.Errorf("Jugador `%s` no encontrado", idJugador)
}

func (r *Ronda) setManojos(manojos []Manojo) {
	for m, manojo := range manojos {
		for c, carta := range manojo.Cartas {
			r.Manojos[m].Cartas[c] = carta
		}
	}
}

func (r *Ronda) setMuestra(muestra Carta) {
	r.Muestra = muestra
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
			r.Manojos[idxJugador].Cartas[idxCarta] = carta
			r.Manojos[idxJugador].CartasNoTiradas[idxCarta] = true
		}
	}

	// la ultima es la muestra
	n := cap(randomCards)
	r.Muestra = nuevaCarta(CartaID(randomCards[n-1]))
}

// nuevaRonda : crea una nueva ronda al azar
func nuevaRonda(jugadores []Jugador, elMano JugadorIdx) Ronda {
	cantJugadores := len(jugadores)
	cantJugadoresPorEquipo := cantJugadores / 2
	ronda := Ronda{
		ManoEnJuego: primera,
		CantJugadoresEnJuego: map[Equipo]int{
			Rojo: cantJugadoresPorEquipo,
			Azul: cantJugadoresPorEquipo,
		},
		ElMano:  elMano,
		Turno:   elMano,
		Envido:  Envido{Puntaje: 0, Estado: NOCANTADOAUN},
		Flor:    NOCANTADA,
		Truco:   truco{CantadoPor: nil, Estado: NOCANTADO},
		Manojos: make([]Manojo, cantJugadores),
		Manos:   make([]Mano, 3),
	}

	// reparto 3 cartas al azar a cada jugador
	// y ademas una muestra, tambien al azar.
	ronda.dealCards()

	// // hago el SINGLE-linking "jugadores <- manojos"
	ronda.singleLinking(jugadores)

	// p.Ronda.setTurno()

	return ronda
}
