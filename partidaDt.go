package truco

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Puntuacion : Enum para el puntaje maximo de la partida
type Puntuacion int

// hasta 15 pts, 20 pts, 30 pts o 40 pts
const (
	a20 Puntuacion = 20
	a30 Puntuacion = 30
	a40 Puntuacion = 40
)

func (pt Puntuacion) toInt() int {
	return int(pt)
}

// Equipo : Enum para el puntaje maximo de la partida
type Equipo string

// rojo o azul
const (
	Azul Equipo = "Azul"
	Rojo Equipo = "Rojo"
)

var toEquipo = map[string]Equipo{
	"Azul": Azul,
	"Rojo": Rojo,
}

func (e Equipo) String() string {
	if e == Rojo {
		return "Rojo"
	}
	return "Azul"
}

// MarshalJSON marshals the enum as a quoted json string
func (e Equipo) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(e.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *Equipo) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*e = toEquipo[j]
	return nil
}

// PartidaDT solo los datos de una partida
type PartidaDT struct {
	jugadores     []Jugador
	CantJugadores int            `json:"cantJugadores"`
	Puntuacion    Puntuacion     `json:"puntuacion"`
	Puntajes      map[Equipo]int `json:"puntajes"`
	Ronda         Ronda          `json:"ronda"`
}

func (p *PartidaDT) getMaxPuntaje() int {
	if p.Puntajes[Rojo] > p.Puntajes[Azul] {
		return p.Puntajes[Rojo]
	}
	return p.Puntajes[Azul]
}

// retorna el equipo que va ganando
func (p *PartidaDT) elQueVaGanando() Equipo {
	vaGanandoRojo := p.Puntajes[Rojo] > p.Puntajes[Azul]
	if vaGanandoRojo {
		return Rojo
	}
	return Azul
}

// getPuntuacionMalas devuelve la mitad de la puntuacion
// total jugable durante toda la partida
func (p *PartidaDT) getPuntuacionMalas() int {
	return p.Puntuacion.toInt() / 2
}

// Terminada retorna true si la partida acabo
func (p *PartidaDT) Terminada() bool {
	return p.getMaxPuntaje() >= p.Puntuacion.toInt()
}

func (p *PartidaDT) elChico() int {
	return p.Puntuacion.toInt() / 2
}

// retorna true si `e` esta en malas
func (p *PartidaDT) estaEnMalas(e Equipo) bool {
	return p.Puntajes[e] < p.elChico()
}

// retorna la cantidad de puntos que le falta para ganar al que va ganando
func (p *PartidaDT) calcPtsFalta() int {
	return p.Puntuacion.toInt() - p.Puntajes[p.elQueVaGanando()]
}

// retorna la cantidad de puntos que le corresponderian
// a `ganadorDelEnvite` si hubiese ganado un "Contra flor al resto"
// sin tener en cuenta los puntos acumulados de envites anteriores
func (p *PartidaDT) calcPtsContraFlorAlResto(ganadorDelEnvite Equipo) int {
	return p.calcPtsFaltaEnvido(ganadorDelEnvite)
}

// retorna la cantidad de puntos que corresponden al Falta-Envido
func (p *PartidaDT) calcPtsFaltaEnvido(ganadorDelEnvite Equipo) int {
	// si el que va ganando:
	// 		esta en Malas -> el ganador del envite (`ganadorDelEnvite`) gana el chico
	// 		esta en Buenas -> el ganador del envite (`ganadorDelEnvite`) gana lo que le falta al maximo para ganar la ronda

	if p.estaEnMalas(p.elQueVaGanando()) {
		loQueLeFaltaAlGANADORparaGanarElChico := p.elChico() - p.Puntajes[ganadorDelEnvite]
		return loQueLeFaltaAlGANADORparaGanarElChico
	}
	//else {
	loQueLeFaltaAlQUEvaGANANDOparaGanarElChico := p.calcPtsFalta()
	return loQueLeFaltaAlQUEvaGANANDOparaGanarElChico
	//}

}

// retorna true si termino la partida
func (p *PartidaDT) sumarPuntos(e Equipo, totalPts int) bool {
	p.Puntajes[e] += totalPts
	return p.Terminada()
}

// evalua todas las cartas y decide que equipo gano
// de ese ganador se setea el siguiente turno
func (p *PartidaDT) evaluarMano() (bool, []*Pkt) {

	pkts := make([]*Pkt, 3)

	// cual es la tirada-carta que gano la mano?
	// ojo que puede salir parda
	// para ello primero busco las maximas de cada equipo
	// y luego comparo entre estas para simplificar
	// Obs: en caso de 2 jugadores del mismo que tiraron
	// una carta con el mismo poder -> se queda con la primera
	// es decir, la que "gana de mano"
	maxPoder := map[Equipo]int{Rojo: -1, Azul: -1}
	max := map[Equipo]*tirarCarta{Rojo: nil, Azul: nil}
	tiradas := p.Ronda.getManoActual().CartasTiradas

	for i, tirada := range tiradas {
		poder := tirada.Carta.calcPoder(p.Ronda.Muestra)
		equipo := tirada.autor.Jugador.Equipo
		if poder > maxPoder[equipo] {
			maxPoder[equipo] = poder
			max[equipo] = &tiradas[i]
		}
	}

	mano := p.Ronda.getManoActual()
	esParda := maxPoder[Rojo] == maxPoder[Azul]

	// caso particular de parda:
	// cuando nadie llego a tirar ninguna carta y se fueron todos los de 1 equipo
	// entonces la mano es ganada por el equipo contrario al ultimo que se fue
	noSeLlegoATirarNingunaCarta := len(p.Ronda.getManoActual().CartasTiradas) == 0

	if noSeLlegoATirarNingunaCarta {

		var equipoGanador Equipo
		quedanJugadoresDelRojo := p.Ronda.CantJugadoresEnJuego[Rojo] > 0
		if quedanJugadoresDelRojo {
			equipoGanador = Rojo
			mano.Resultado = GanoRojo
		} else {
			equipoGanador = Azul
			mano.Resultado = GanoAzul
		}

		// aca le tengo que poner un ganador para despues sacarle el equipo
		// le asigno el primero que encuentre del equipo ganador
		if p.Ronda.Manojos[0].Jugador.Equipo == equipoGanador {
			mano.Ganador = &p.Ronda.Manojos[0]
		} else {
			mano.Ganador = &p.Ronda.Manojos[1]
		}

		// fmt.Printf("La %s mano la gano el equipo %s\n",
		// 	strings.ToLower(p.Ronda.ManoEnJuego.String()),
		// 	equipoGanador.String())

	} else if esParda {
		mano.Resultado = Empardada
		mano.Ganador = nil

		pkts[0] = &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Info",
				Cont: []byte("La Mano resulta parda"),
			},
		}
		// no se cambia el turno

	} else {
		var tiradaGanadora *tirarCarta

		if maxPoder[Rojo] > maxPoder[Azul] {
			tiradaGanadora = max[Rojo]
			mano.Resultado = GanoRojo
		} else {
			tiradaGanadora = max[Azul]
			mano.Resultado = GanoAzul
		}

		// el turno pasa a ser el del mano.ganador
		// pero se setea despues de evaluar la ronda
		mano.Ganador = tiradaGanadora.autor

		pkts[0] = &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Info",
				Cont: []byte(fmt.Sprintf("La %s mano la gano el equipo %s gracia a %s",
					strings.ToLower(p.Ronda.ManoEnJuego.String()),
					mano.Ganador.Jugador.Equipo.String(),
					mano.Ganador.Jugador.Nombre)),
			},
		}

	}

	// se termino la ronda?
	empiezaNuevaRonda, pkt2 := p.evaluarRonda()

	pkts = append(pkts, pkt2)

	// cuando termina la mano (y no se empieza una ronda) -> cambia de TRUNO
	// cuando termina la ronda -> cambia de MANO
	// para usar esto, antes se debe primero incrementar el turno
	// incremento solo si no se empezo una nueva ronda

	return empiezaNuevaRonda, pkts
}

// tener siempre en cuenta que evaluar la ronda es sinonimo de evaluar el truco
// se acabo la ronda?
// si se empieza una ronda nueva -> retorna true
// si no se termino la ronda 	 -> retorna false
func (p *PartidaDT) evaluarRonda() (bool, *Pkt) {

	pkt := new(Pkt)

	/* A MENOS QUE SE HAYAN IDO TODOS EN LA PRIMERA MANO!!! */
	hayJugadoresRojo := p.Ronda.CantJugadoresEnJuego[Rojo] > 0
	hayJugadoresAzul := p.Ronda.CantJugadoresEnJuego[Azul] > 0
	hayJugadoresEnAmbos := hayJugadoresRojo && hayJugadoresAzul

	imposibleQueSeHayaAcabado := (p.Ronda.ManoEnJuego == primera) && hayJugadoresEnAmbos
	if imposibleQueSeHayaAcabado {
		return false, nil
	}

	// de aca en mas ya se que hay al menos 2 manos jugadas
	// (excepto el caso en que un equipo haya abandonado)
	// asi que es seguro acceder a los indices 0 y 1 en:
	// p.Ronda.manos[0] & p.Ronda.manos[1]

	cantManosGanadas := map[Equipo]int{Rojo: 0, Azul: 0}
	for i := 0; i < p.Ronda.ManoEnJuego.toInt(); i++ {
		mano := p.Ronda.Manos[i]
		if mano.Resultado != Empardada {
			cantManosGanadas[mano.Ganador.Jugador.Equipo]++
		}
	}

	hayEmpate := cantManosGanadas[Rojo] == cantManosGanadas[Azul]
	pardaPrimera := p.Ronda.Manos[0].Resultado == Empardada
	pardaSegunda := p.Ronda.Manos[1].Resultado == Empardada
	pardaTercera := p.Ronda.Manos[2].Resultado == Empardada
	seEstaJugandoLaSegunda := p.Ronda.ManoEnJuego == segunda
	noSeAcaboAun := seEstaJugandoLaSegunda && hayEmpate && hayJugadoresEnAmbos

	if noSeAcaboAun {
		return false, nil
	}

	// hay ganador -> ya se que al final voy a retornar un true
	var ganador *Manojo

	if !hayJugadoresEnAmbos { // caso particular: todos abandonaron

		// enonces como antes paso por evaluar mano
		// y seteo a ganador de la ultima mano jugada (la "actual")
		// al equipo que no abandono -> lo sacao de ahi
		ganador = p.Ronda.getManoActual().Ganador

		// primero el caso clasico: un equipo gano 2 o mas manos
	} else if cantManosGanadas[Rojo] >= 2 {
		// agarro cualquier manojo de los rojos
		// o bien es la primera o bien la segunda
		if p.Ronda.Manos[0].Ganador.Jugador.Equipo == Rojo {
			ganador = p.Ronda.Manos[0].Ganador
		} else {
			ganador = p.Ronda.Manos[1].Ganador
		}
	} else if cantManosGanadas[Azul] >= 2 {
		// agarro cualquier manojo de los azules
		// o bien es la primera o bien la segunda
		if p.Ronda.Manos[0].Ganador.Jugador.Equipo == Azul {
			ganador = p.Ronda.Manos[0].Ganador
		} else {
			ganador = p.Ronda.Manos[1].Ganador
		}

	} else {

		// si llego aca es porque recae en uno de los
		// siguientes casos: (Obs: se jugo la tercera)

		// CASO 1. parda primera -> gana segunda
		// CASO 2. parda segunda -> gana primera
		// CASO 3. parda tercera -> gana primera
		// CASO 4. parda primera & segunda -> gana tercera
		// CASO 5. parda primera, segunda & tercera -> gana la mano

		caso1 := pardaPrimera && !pardaSegunda && !pardaTercera
		caso2 := !pardaPrimera && pardaSegunda && !pardaTercera
		caso3 := !pardaPrimera && !pardaSegunda && pardaTercera
		caso4 := pardaPrimera && pardaSegunda && !pardaTercera
		caso5 := pardaPrimera && pardaSegunda && pardaTercera

		if caso1 {
			ganador = p.Ronda.Manos[segunda].Ganador

		} else if caso2 {
			ganador = p.Ronda.Manos[primera].Ganador

		} else if caso3 {
			ganador = p.Ronda.Manos[primera].Ganador

		} else if caso4 {
			ganador = p.Ronda.Manos[tercera].Ganador

		} else if caso5 {
			ganador = p.Ronda.getElMano()
		}

	}

	/************************************************/

	// ya sabemos el ganador ahora es el
	// momento de sumar los puntos del truco
	var totalPts int = 0
	var msg string

	switch p.Ronda.Truco.Estado {
	case NOCANTADO, TRUCO: // caso en que se hayan ido todos al mazo y no se haya respondido ~ equivalente a un no quiero
		totalPts = 1
	case TRUCOQUERIDO, RETRUCO: // same
		totalPts = 2
	case RETRUCOQUERIDO, VALE4: // same
		totalPts = 3
	case VALE4QUERIDO:
		totalPts = 4
	}

	elTrucoNoTuvoRespuesta := contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado)

	if elTrucoNoTuvoRespuesta {
		msg = fmt.Sprintf(`La ronda ha sido ganada por el equipo %s. +%v puntos para el equipo %s por el %s no querido`,
			ganador.Jugador.Equipo,
			totalPts,
			ganador.Jugador.Equipo,
			p.Ronda.Truco.Estado.String())

	} else {
		msg = fmt.Sprintf(`La ronda ha sido ganada por el equipo %s. +%v puntos para el equipo %s por el %s ganado`,
			ganador.Jugador.Equipo,
			totalPts,
			ganador.Jugador.Equipo,
			p.Ronda.Truco.Estado.String())
	}

	pkt = &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "Sumar-Puntos",
			Nota: msg,
			Cont: ContSumPts{
				Pts:    totalPts,
				Equipo: ganador.Jugador.Equipo.String(),
			}.ToJSON(),
		},
	}

	p.sumarPuntos(ganador.Jugador.Equipo, totalPts)

	return true, pkt // porque se empezo una nueva ronda
}

func (p *PartidaDT) nuevaRonda(elMano JugadorIdx) {
	p.Ronda = nuevaRonda(p.jugadores, elMano)
}

// ToString retorna su render
func (p *PartidaDT) ToString() string {
	render := renderizar(p)
	return render
}

// Print imprime la partida
func (p *PartidaDT) Print() {
	fmt.Print(p.ToString())
}

// ToJSON retorna la partida en formato json
func (p *PartidaDT) ToJSON() string {
	pJSON, _ := json.Marshal(p)
	return string(pJSON)
}

// FromJSON carga una partida en formato json
// ojo que este es PUBLICO!!! no cachea flores
func (p *PartidaDT) FromJSON(partidaJSON string) error {
	err := json.Unmarshal([]byte(partidaJSON), &p)
	if err != nil {
		return err
	}
	return nil
}

// FromJSON carga una partida en formato json
// ojo que este es PRIVADO!!!
// el metodo privado cachea las flores
func (p *PartidaDT) fromJSON(partidaJSON string) error {
	err := json.Unmarshal([]byte(partidaJSON), &p)
	if err != nil {
		return err
	}
	p.Ronda.cachearFlores()
	return nil
}

// Force para hacer debugs
func (p *PartidaDT) Force(partidaJSON string) error {
	return p.fromJSON(partidaJSON)
}

func cheepCopy(p *PartidaDT) *PartidaDT {
	copia := *p
	copia.Ronda.Manojos = make([]Manojo, len(p.Ronda.Manojos))
	copy(copia.Ronda.Manojos, p.Ronda.Manojos)

	return &copia
}

// Perspectiva retorna una representacion en json de la perspectiva que tiene
// el jugador `j` de la partida
func (p *PartidaDT) Perspectiva(j string) (*PartidaDT, error) {
	copia := cheepCopy(p)

	// primero encuentro el jugador
	manojo, err := copia.Ronda.getManojoByStr(j)
	if err != nil {
		return nil, fmt.Errorf("Usuario %s no encontrado", j)
	}

	// oculto las caras no tiradas de los manojos que no son el
	for i, m := range copia.Ronda.Manojos {
		noEsSuManojo := m.Jugador.ID != manojo.Jugador.ID
		if noEsSuManojo {
			// oculto solo las cartas que no tiro
			for j, noTirada := range m.CartasNoTiradas {
				if noTirada {
					copia.Ronda.Manojos[i].Cartas[j] = nil
				}
			}
		}
	}

	return copia, nil
}

// NuevaPartidaDt crea una nueva PartidaDT
func NuevaPartidaDt(puntuacion Puntuacion, equipoAzul, equipoRojo []string) (*PartidaDT, error) {

	mismaCantidadDeJugadores := len(equipoRojo) == len(equipoAzul)
	cantJugadores := len(equipoRojo) + len(equipoAzul)
	cantidadCorrecta := contains([]int{2, 4, 6}, cantJugadores) // puede ser 2, 4 o 6
	ok := mismaCantidadDeJugadores && cantidadCorrecta

	if !ok {
		return nil, fmt.Errorf(`La cantidad de jugadores no es correcta`)
	}
	// paso a crear los jugadores; intercalados
	var jugadores []Jugador
	// para cada rjo que agrego; le agrego tambien su mano
	for i := range equipoRojo {
		// uso como id sus nombres
		nuevoJugadorRojo := Jugador{equipoRojo[i], equipoRojo[i], Rojo}
		nuevoJugadorAzul := Jugador{equipoAzul[i], equipoAzul[i], Azul}
		jugadores = append(jugadores, nuevoJugadorAzul, nuevoJugadorRojo)
	}

	p := PartidaDT{
		Puntuacion:    puntuacion,
		CantJugadores: cantJugadores,
		jugadores:     jugadores,
	}

	p.Puntajes = make(map[Equipo]int)
	p.Puntajes[Rojo] = 0
	p.Puntajes[Azul] = 0

	elMano := JugadorIdx(0)
	p.nuevaRonda(elMano)

	return &p, nil
}
