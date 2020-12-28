package pdt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/filevich/truco/out"
)

// Puntuacion : Enum para el puntaje maximo de la partida
type Puntuacion int

// hasta 15 pts, 20 pts, 30 pts o 40 pts
const (
	A20 Puntuacion = 20
	A30 Puntuacion = 30
	A40 Puntuacion = 40
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
	Jugadores     []Jugador
	CantJugadores int            `json:"cantJugadores"`
	Puntuacion    Puntuacion     `json:"puntuacion"`
	Puntajes      map[Equipo]int `json:"puntajes"`
	Ronda         Ronda          `json:"ronda"`
}

// GetMaxPuntaje .
func (p *PartidaDT) GetMaxPuntaje() int {
	if p.Puntajes[Rojo] > p.Puntajes[Azul] {
		return p.Puntajes[Rojo]
	}
	return p.Puntajes[Azul]
}

// ElQueVaGanando retorna el equipo que va ganando
func (p *PartidaDT) ElQueVaGanando() Equipo {
	vaGanandoRojo := p.Puntajes[Rojo] > p.Puntajes[Azul]
	if vaGanandoRojo {
		return Rojo
	}
	return Azul
}

// GetPuntuacionMalas devuelve la mitad de la puntuacion
// total jugable durante toda la partida
func (p *PartidaDT) GetPuntuacionMalas() int {
	return p.Puntuacion.toInt() / 2
}

// EsManoAMano retorna true sii la partida consta de exactamente 2 jugadores
func (p *PartidaDT) EsManoAMano() bool {
	return p.CantJugadores == 2
}

// Terminada retorna true si la partida acabo
func (p *PartidaDT) Terminada() bool {
	return p.GetMaxPuntaje() >= p.Puntuacion.toInt()
}

// ElChico .
func (p *PartidaDT) ElChico() int {
	return p.Puntuacion.toInt() / 2
}

// EstaEnMalas retorna true si `e` esta en malas
func (p *PartidaDT) EstaEnMalas(e Equipo) bool {
	return p.Puntajes[e] < p.ElChico()
}

// CalcPtsFalta retorna la cantidad de puntos que le falta para ganar al que va ganando
func (p *PartidaDT) CalcPtsFalta() int {
	return p.Puntuacion.toInt() - p.Puntajes[p.ElQueVaGanando()]
}

// CalcPtsContraFlorAlResto etorna la cantidad de puntos que le corresponderian
// a `ganadorDelEnvite` si hubiese ganado un "Contra flor al resto"
// sin tener en cuenta los puntos acumulados de envites anteriores
func (p *PartidaDT) CalcPtsContraFlorAlResto(ganadorDelEnvite Equipo) int {
	return p.CalcPtsFaltaEnvido(ganadorDelEnvite)
}

// CalcPtsFaltaEnvido retorna la cantidad de puntos que corresponden al Falta-Envido
func (p *PartidaDT) CalcPtsFaltaEnvido(ganadorDelEnvite Equipo) int {
	// si el que va ganando:
	// 		esta en Malas -> el ganador del envite (`ganadorDelEnvite`) gana el chico
	// 		esta en Buenas -> el ganador del envite (`ganadorDelEnvite`) gana lo que le falta al maximo para ganar la ronda

	if p.EstaEnMalas(p.ElQueVaGanando()) {
		loQueLeFaltaAlGANADORparaGanarElChico := p.ElChico() - p.Puntajes[ganadorDelEnvite]
		return loQueLeFaltaAlGANADORparaGanarElChico
	}
	//else {
	loQueLeFaltaAlQUEvaGANANDOparaGanarElChico := p.CalcPtsFalta()
	return loQueLeFaltaAlQUEvaGANANDOparaGanarElChico
	//}

}

// SumarPuntos retorna true si termino la partida
func (p *PartidaDT) SumarPuntos(e Equipo, totalPts int) bool {
	p.Puntajes[e] += totalPts
	return p.Terminada()
}

// TocarEnvido ..
func (p *PartidaDT) TocarEnvido(m *Manojo) {
	// 2 opciones: o bien no se jugo aun
	// o bien ya estabamos en envido
	yaSeHabiaCantadoElEnvido := p.Ronda.Envite.Estado == ENVIDO
	if yaSeHabiaCantadoElEnvido {
		// se aumenta el puntaje del envido en +2
		p.Ronda.Envite.Puntaje += 2
		p.Ronda.Envite.CantadoPor = m

	} else { // no se habia jugado aun
		p.Ronda.Envite.CantadoPor = m
		p.Ronda.Envite.Estado = ENVIDO
		p.Ronda.Envite.Puntaje = 2
	}
}

// TocarRealEnvido ..
func (p *PartidaDT) TocarRealEnvido(m *Manojo) {
	p.Ronda.Envite.CantadoPor = m
	// 2 opciones:
	// o bien el envido no se jugo aun,
	// o bien ya estabamos en envido
	if p.Ronda.Envite.Estado == NOCANTADOAUN { // no se habia jugado aun
		p.Ronda.Envite.Puntaje = 3
	} else { // ya se habia cantado ENVIDO x cantidad de veces
		p.Ronda.Envite.Puntaje += 3
	}
	p.Ronda.Envite.Estado = REALENVIDO
}

// TocarFaltaEnvido ..
func (p *PartidaDT) TocarFaltaEnvido(m *Manojo) {
	p.Ronda.Envite.Estado = FALTAENVIDO
	p.Ronda.Envite.CantadoPor = m
}

// CantarFlor ..
func (p *PartidaDT) CantarFlor(m *Manojo) {
	p.Ronda.Envite.Estado = FLOR
	yaEstabamosEnFlor := p.Ronda.Envite.Estado == FLOR

	if yaEstabamosEnFlor {

		p.Ronda.Envite.Puntaje += 3
		p.Ronda.Envite.CantadoPor = m

	} else {

		// se usa por si dicen "no quiero" -> se obtiene el equipo
		// al que pertenece el que la canto en un principio para
		// poder sumarle los puntos correspondientes
		p.Ronda.Envite.Puntaje = 3
		p.Ronda.Envite.CantadoPor = m
		p.Ronda.Envite.Estado = FLOR

	}
}

// CantarContraFlor ..
func (p *PartidaDT) CantarContraFlor(m *Manojo) {
	p.Ronda.Envite.Estado = CONTRAFLOR
	p.Ronda.Envite.CantadoPor = m
	// ahora la flor pasa a jugarse por 4 puntos
	p.Ronda.Envite.Puntaje = 4
}

// CantarContraFlorAlResto ..
func (p *PartidaDT) CantarContraFlorAlResto(m *Manojo) {
	p.Ronda.Envite.Estado = CONTRAFLORALRESTO
	p.Ronda.Envite.CantadoPor = m
	// ahora la flor pasa a jugarse por 4 puntos
	p.Ronda.Envite.Puntaje = 4 // <- eso es al pedo, es independiente
}

// GritarTruco ..
func (p *PartidaDT) GritarTruco(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m
	p.Ronda.Truco.Estado = TRUCO
	// p.Ronda.Envite.Estado = DESHABILITADO // <-- esto esta mal
}

// QuererTruco incrementa el estado del truco a querido segun corresponda
// y setea a m como el que lo canto
func (p *PartidaDT) QuererTruco(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m
	switch p.Ronda.Truco.Estado {
	case TRUCO:
		p.Ronda.Truco.Estado = TRUCOQUERIDO
	case RETRUCO:
		p.Ronda.Truco.Estado = RETRUCOQUERIDO
	case VALE4:
		p.Ronda.Truco.Estado = VALE4QUERIDO
	}
}

// GritarReTruco ..
func (p *PartidaDT) GritarReTruco(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m
	p.Ronda.Truco.Estado = RETRUCO
	// p.Ronda.Envite.Estado = DESHABILITADO // <-- esto esta mal
}

// GritarVale4 ..
func (p *PartidaDT) GritarVale4(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m
	p.Ronda.Truco.Estado = VALE4
	// p.Ronda.Envite.Estado = DESHABILITADO // <-- esto esta mal
}

// IrAlMazo manda el manojo al mazo
func (p *PartidaDT) IrAlMazo(manojo *Manojo) {
	manojo.SeFueAlMazo = true
	equipoDelJugador := manojo.Jugador.Equipo
	p.Ronda.CantJugadoresEnJuego[equipoDelJugador]--
}

// EvaluarMano evalua todas las cartas y decide que equipo gano
// de ese ganador se setea el siguiente turno
func (p *PartidaDT) EvaluarMano() (bool, []*out.Packet) {

	var pkts []*out.Packet

	// cual es la tirada-carta que gano la mano?
	// ojo que puede salir parda
	// para ello primero busco las maximas de cada equipo
	// y luego comparo entre estas para simplificar
	// Obs: en caso de 2 jugadores del mismo que tiraron
	// una carta con el mismo poder -> se queda con la Primera
	// es decir, la que "gana de mano"
	maxPoder := map[Equipo]int{Rojo: -1, Azul: -1}
	max := map[Equipo]*tirarCarta{Rojo: nil, Azul: nil}
	tiradas := p.Ronda.GetManoActual().CartasTiradas

	for i, tirada := range tiradas {
		poder := tirada.Carta.calcPoder(p.Ronda.Muestra)
		equipo := tirada.autor.Jugador.Equipo
		if poder > maxPoder[equipo] {
			maxPoder[equipo] = poder
			max[equipo] = &tiradas[i]
		}
	}

	mano := p.Ronda.GetManoActual()
	esParda := maxPoder[Rojo] == maxPoder[Azul]

	// caso particular de parda:
	// cuando nadie llego a tirar ninguna carta y se fueron todos los de 1 equipo
	// entonces la mano es ganada por el equipo contrario al ultimo que se fue
	noSeLlegoATirarNingunaCarta := len(p.Ronda.GetManoActual().CartasTiradas) == 0

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

		// todo: MENSAJE ACAAAAAAA!!!

		// fmt.Printf("La %s mano la gano el equipo %s\n",
		// 	strings.ToLower(p.Ronda.ManoEnJuego.String()),
		// 	equipoGanador.String())

	} else if esParda {
		mano.Resultado = Empardada
		mano.Ganador = nil

		pkts = append(pkts, out.Pkt(
			out.Dest("ALL"),
			out.Msg(out.Info, "La Mano resulta parda"),
		))

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

		pkts = append(pkts, out.Pkt(
			out.Dest("ALL"),
			out.Msg(out.Info,
				fmt.Sprintf("La %s mano la gano el equipo %s gracias a %s",
					strings.ToLower(p.Ronda.ManoEnJuego.String()),
					mano.Ganador.Jugador.Equipo.String(),
					mano.Ganador.Jugador.Nombre),
			),
		))

	}

	// se termino la ronda?
	empiezaNuevaRonda, pkt2 := p.EvaluarRonda()

	pkts = append(pkts, pkt2...)

	// cuando termina la mano (y no se empieza una ronda) -> cambia de TRUNO
	// cuando termina la ronda -> cambia de MANO
	// para usar esto, antes se debe primero incrementar el turno
	// incremento solo si no se empezo una nueva ronda

	return empiezaNuevaRonda, pkts
}

// EvaluarRonda tener siempre en cuenta que evaluar la ronda es sinonimo de
// evaluar el truco
// se acabo la ronda?
// si se empieza una ronda nueva -> retorna true
// si no se termino la ronda 	 -> retorna false
func (p *PartidaDT) EvaluarRonda() (bool, []*out.Packet) {

	var pkts []*out.Packet

	/* A MENOS QUE SE HAYAN IDO TODOS EN LA PRIMERA MANO!!! */
	hayJugadoresRojo := p.Ronda.CantJugadoresEnJuego[Rojo] > 0
	hayJugadoresAzul := p.Ronda.CantJugadoresEnJuego[Azul] > 0
	hayJugadoresEnAmbos := hayJugadoresRojo && hayJugadoresAzul
	primeraMano := p.Ronda.ManoEnJuego == Primera

	// o bien que en la primera mano hayan cantado truco y uno no lo quizo
	manoActual := p.Ronda.ManoEnJuego.ToInt() - 1
	elTrucoNoTuvoRespuesta := Contains([]EstadoTruco{TRUCO, RETRUCO, VALE4}, p.Ronda.Truco.Estado)
	noFueParda := p.Ronda.Manos[manoActual].Resultado != Empardada
	estaManoYaTieneGanador := noFueParda && p.Ronda.Manos[manoActual].Ganador != nil
	elTrucoFueNoQuerido := elTrucoNoTuvoRespuesta && estaManoYaTieneGanador

	noSeAcabo := (primeraMano && hayJugadoresEnAmbos && !elTrucoFueNoQuerido)
	if noSeAcabo {
		return false, nil
	}

	// de aca en mas ya se que hay al menos 2 manos jugadas
	// (excepto el caso en que un equipo haya abandonado)
	// asi que es seguro acceder a los indices 0 y 1 en:
	// p.Ronda.manos[0] & p.Ronda.manos[1]

	cantManosGanadas := map[Equipo]int{Rojo: 0, Azul: 0}
	for i := 0; i < p.Ronda.ManoEnJuego.ToInt(); i++ {
		mano := p.Ronda.Manos[i]
		if mano.Resultado != Empardada {
			cantManosGanadas[mano.Ganador.Jugador.Equipo]++
		}
	}

	hayEmpate := cantManosGanadas[Rojo] == cantManosGanadas[Azul]
	pardaPrimera := p.Ronda.Manos[0].Resultado == Empardada
	pardaSegunda := p.Ronda.Manos[1].Resultado == Empardada
	pardaTercera := p.Ronda.Manos[2].Resultado == Empardada
	seEstaJugandoLaSegunda := p.Ronda.ManoEnJuego == Segunda
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
		ganador = p.Ronda.GetManoActual().Ganador

		// primero el caso clasico: un equipo gano 2 o mas manos
	} else if cantManosGanadas[Rojo] >= 2 {
		// agarro cualquier manojo de los rojos
		// o bien es la Primera o bien la Segunda
		if p.Ronda.Manos[0].Ganador.Jugador.Equipo == Rojo {
			ganador = p.Ronda.Manos[0].Ganador
		} else {
			ganador = p.Ronda.Manos[1].Ganador
		}
	} else if cantManosGanadas[Azul] >= 2 {
		// agarro cualquier manojo de los azules
		// o bien es la Primera o bien la Segunda
		if p.Ronda.Manos[0].Ganador.Jugador.Equipo == Azul {
			ganador = p.Ronda.Manos[0].Ganador
		} else {
			ganador = p.Ronda.Manos[1].Ganador
		}

	} else {

		// si llego aca es porque recae en uno de los
		// siguientes casos: (Obs: se jugo la Tercera)

		// CASO 1. parda Primera -> gana Segunda
		// CASO 2. parda Segunda -> gana Primera
		// CASO 3. parda Tercera -> gana Primera
		// CASO 4. parda Primera & Segunda -> gana Tercera
		// CASO 5. parda Primera, Segunda & Tercera -> gana la mano

		caso1 := pardaPrimera && !pardaSegunda && !pardaTercera
		caso2 := !pardaPrimera && pardaSegunda && !pardaTercera
		caso3 := !pardaPrimera && !pardaSegunda && pardaTercera
		caso4 := pardaPrimera && pardaSegunda && !pardaTercera
		caso5 := pardaPrimera && pardaSegunda && pardaTercera

		if caso1 {
			ganador = p.Ronda.Manos[Segunda].Ganador

		} else if caso2 {
			ganador = p.Ronda.Manos[Primera].Ganador

		} else if caso3 {
			ganador = p.Ronda.Manos[Primera].Ganador

		} else if caso4 {
			ganador = p.Ronda.Manos[Tercera].Ganador

		} else if caso5 {
			ganador = p.Ronda.GetElMano()
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

	if elTrucoNoTuvoRespuesta {
		ganador = p.Ronda.Truco.CantadoPor
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

	pkts = append(pkts, out.Pkt(
		out.Dest("ALL"),
		out.Msg(out.Info, msg),
	))

	pkts = append(pkts, out.Pkt(
		out.Dest("ALL"),
		out.Msg(out.SumaPts, ganador.Jugador.ID, out.TrucoQuerido, totalPts),
	))

	p.SumarPuntos(ganador.Jugador.Equipo, totalPts)

	return true, pkts // porque se empezo una nueva ronda
}

// NuevaRonda .
func (p *PartidaDT) NuevaRonda(elMano JugadorIdx) {
	p.Ronda = NuevaRonda(p.Jugadores, elMano)
}

// ToString retorna su render
// func (p *PartidaDT) ToString() string {
// 	render := ptr.Renderizar(p)
// 	return render
// }

// Print imprime la partida
// func (p *PartidaDT) Print() {
// 	fmt.Print(p.ToString())
// }

// MarshalJSON retorna la partida en formato json
func (p *PartidaDT) MarshalJSON() ([]byte, error) {
	// return json.Marshal(p)
	// return []byte("hola"), nil
	return json.Marshal(*p)
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

// Perspectiva retorna una representacion en json de la PerspectivaCacheFlor que tiene
// el jugador `j` de la partida
func (p *PartidaDT) Perspectiva(j string) (*PartidaDT, error) {
	// primero encuentro el jugador
	manojo, err := p.Ronda.GetManojoByStr(j)
	if err != nil {
		return nil, fmt.Errorf("Usuario %s no encontrado", j)
	}

	return p.PerspectivaCacheFlor(manojo), nil
}

// PerspectivaCacheFlor cache las flores
func (p *PartidaDT) PerspectivaCacheFlor(manojo *Manojo) *PartidaDT {
	copia := cheepCopy(p)

	// oculto las caras no tiradas de los manojos que no son de su equipo
	for i, m := range copia.Ronda.Manojos {
		noEsDeSuEquipo := m.Jugador.Equipo != manojo.Jugador.Equipo
		if noEsDeSuEquipo {
			// oculto solo las cartas que no tiro
			for j, noTirada := range m.CartasNoTiradas {
				if noTirada {
					copia.Ronda.Manojos[i].Cartas[j] = nil
				}
			}
		}
	}

	return copia
}

/* metodos de manipulacion */

// TirarCarta tira una carta
func (p *PartidaDT) TirarCarta(manojo *Manojo, idx int) {
	manojo.CartasNoTiradas[idx] = false
	manojo.UltimaTirada = idx
	carta := manojo.Cartas[idx]
	tirada := tirarCarta{manojo, *carta}
	p.Ronda.GetManoActual().agregarTirada(tirada)
}

// NuevaPartidaDt crea una nueva PartidaDT
func NuevaPartidaDt(puntuacion Puntuacion, equipoAzul, equipoRojo []string) (*PartidaDT, error) {

	mismaCantidadDeJugadores := len(equipoRojo) == len(equipoAzul)
	cantJugadores := len(equipoRojo) + len(equipoAzul)
	cantidadCorrecta := Contains([]int{2, 4, 6}, cantJugadores) // puede ser 2, 4 o 6
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
		Jugadores:     jugadores,
	}

	p.Puntajes = make(map[Equipo]int)
	p.Puntajes[Rojo] = 0
	p.Puntajes[Azul] = 0

	elMano := JugadorIdx(0)
	p.NuevaRonda(elMano)

	return &p, nil
}
