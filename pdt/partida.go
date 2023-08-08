package pdt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/filevich/truco/enco"
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

// Partida solo los datos de una partida
type Partida struct {
	Puntuacion   Puntuacion     `json:"puntuacion"`
	Puntajes     map[Equipo]int `json:"puntajes"`
	Ronda        Ronda          `json:"ronda"`
	LimiteEnvido int            `json:"limiteEnvido"`
	Verbose      bool           `json:"-"`
}

// GetMaxPuntaje .
func (p *Partida) GetMaxPuntaje() int {
	if p.Puntajes[Rojo] > p.Puntajes[Azul] {
		return p.Puntajes[Rojo]
	}
	return p.Puntajes[Azul]
}

// ElQueVaGanando retorna el equipo que va ganando
func (p *Partida) ElQueVaGanando() Equipo {
	vaGanandoRojo := p.Puntajes[Rojo] > p.Puntajes[Azul]
	if vaGanandoRojo {
		return Rojo
	}
	return Azul
}

// GetPuntuacionMalas devuelve la mitad de la puntuacion
// total jugable durante toda la partida
func (p *Partida) GetPuntuacionMalas() int {
	return p.Puntuacion.toInt() / 2
}

// EsManoAMano retorna true sii la partida consta de exactamente 2 jugadores
func (p *Partida) EsManoAMano() bool {
	return len(p.Ronda.Manojos) == 2
}

// Terminada retorna true si la partida acabo
func (p *Partida) Terminada() bool {
	return p.GetMaxPuntaje() >= p.Puntuacion.toInt()
}

// ElChico .
func (p *Partida) ElChico() int {
	return p.Puntuacion.toInt() / 2
}

// EstaEnMalas retorna true si `e` esta en malas
func (p *Partida) EstaEnMalas(e Equipo) bool {
	return p.Puntajes[e] < p.ElChico()
}

// CalcPtsFalta retorna la cantidad de puntos que le falta para ganar al que va ganando
func (p *Partida) CalcPtsFalta() int {
	return p.Puntuacion.toInt() - p.Puntajes[p.ElQueVaGanando()]
}

// CalcPtsContraFlorAlResto etorna la cantidad de puntos que le corresponderian
// a `ganadorDelEnvite` si hubiese ganado un "Contra flor al resto"
// sin tener en cuenta los puntos acumulados de envites anteriores
func (p *Partida) CalcPtsContraFlorAlResto(ganadorDelEnvite Equipo) int {
	return p.CalcPtsFaltaEnvido(ganadorDelEnvite)
}

// CalcPtsFaltaEnvido retorna la cantidad de puntos que corresponden al Falta-Envido
func (p *Partida) CalcPtsFaltaEnvido(ganadorDelEnvite Equipo) int {
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
func (p *Partida) SumarPuntos(e Equipo, totalPts int) bool {
	p.Puntajes[e] += totalPts
	return p.Terminada()
}

// TocarEnvido ..
func (p *Partida) TocarEnvido(m *Manojo) {
	// 2 opciones: o bien no se jugo aun
	// o bien ya estabamos en envido
	yaSeHabiaCantadoElEnvido := p.Ronda.Envite.Estado == ENVIDO
	if yaSeHabiaCantadoElEnvido {
		// se aumenta el puntaje del envido en +2
		p.Ronda.Envite.Puntaje += 2
		p.Ronda.Envite.CantadoPor = m.Jugador.ID

	} else { // no se habia jugado aun
		p.Ronda.Envite.CantadoPor = m.Jugador.ID
		p.Ronda.Envite.Estado = ENVIDO
		p.Ronda.Envite.Puntaje = 2
	}
}

// TocarRealEnvido ..
func (p *Partida) TocarRealEnvido(m *Manojo) {
	p.Ronda.Envite.CantadoPor = m.Jugador.ID
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
func (p *Partida) TocarFaltaEnvido(m *Manojo) {
	p.Ronda.Envite.Estado = FALTAENVIDO
	p.Ronda.Envite.CantadoPor = m.Jugador.ID
}

// CantarFlor ..
func (p *Partida) CantarFlor(m *Manojo) {
	yaEstabamosEnFlor := p.Ronda.Envite.Estado >= FLOR

	if yaEstabamosEnFlor {

		p.Ronda.Envite.Puntaje += 3
		// si estabamos en algo mas grande que `FLOR` -> no lo aumenta
		if p.Ronda.Envite.Estado == FLOR {
			p.Ronda.Envite.CantadoPor = m.Jugador.ID
			p.Ronda.Envite.Estado = FLOR
		}

	} else {

		// se usa por si dicen "no quiero" -> se obtiene el equipo
		// al que pertenece el que la canto en un principio para
		// poder sumarle los puntos correspondientes
		p.Ronda.Envite.Puntaje = 3
		p.Ronda.Envite.CantadoPor = m.Jugador.ID
		p.Ronda.Envite.Estado = FLOR

	}
}

// CantarContraFlor ..
func (p *Partida) CantarContraFlor(m *Manojo) {
	p.Ronda.Envite.Estado = CONTRAFLOR
	p.Ronda.Envite.CantadoPor = m.Jugador.ID
	// ahora la flor pasa a jugarse por 4 puntos
	p.Ronda.Envite.Puntaje = 4
}

// CantarContraFlorAlResto ..
func (p *Partida) CantarContraFlorAlResto(m *Manojo) {
	p.Ronda.Envite.Estado = CONTRAFLORALRESTO
	p.Ronda.Envite.CantadoPor = m.Jugador.ID
	// ahora la flor pasa a jugarse por 4 puntos
	p.Ronda.Envite.Puntaje = 4 // <- eso es al pedo, es independiente
}

// GritarTruco ..
func (p *Partida) GritarTruco(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m.Jugador.ID
	p.Ronda.Truco.Estado = TRUCO
	// p.Ronda.Envite.Estado = DESHABILITADO // <-- esto esta mal
}

// QuererTruco incrementa el estado del truco a querido segun corresponda
// y setea a m como el que lo canto
func (p *Partida) QuererTruco(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m.Jugador.ID
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
func (p *Partida) GritarReTruco(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m.Jugador.ID
	p.Ronda.Truco.Estado = RETRUCO
	// p.Ronda.Envite.Estado = DESHABILITADO // <-- esto esta mal
}

// GritarVale4 ..
func (p *Partida) GritarVale4(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m.Jugador.ID
	p.Ronda.Truco.Estado = VALE4
	// p.Ronda.Envite.Estado = DESHABILITADO // <-- esto esta mal
}

// IrAlMazo manda el manojo al mazo
// todo: esto podria ser un metodo de Ronda, no de partida
func (p *Partida) IrAlMazo(manojo *Manojo) {
	manojo.SeFueAlMazo = true
	equipoDelJugador := manojo.Jugador.Equipo
	p.Ronda.CantJugadoresEnJuego[equipoDelJugador]--
	// lo elimino de los jugadores que tenian flor (si es que tenia)
	xs := p.Ronda.Envite.SinCantar
	for i, jid := range p.Ronda.Envite.SinCantar {
		if manojo.Jugador.ID == jid {
			xs[i] = xs[len(xs)-1]
			xs = xs[:len(xs)-1]
			p.Ronda.Envite.SinCantar = xs
			break
		}
	}
}

// EvaluarMano evalua todas las cartas y decide que equipo gano
// de ese ganador se setea el siguiente turno
func (p *Partida) EvaluarMano() (bool, []enco.Envelope) {

	var pkts2 []enco.Envelope

	// cual es la tirada-carta que gano la mano?
	// ojo que puede salir parda
	// para ello primero busco las maximas de cada equipo
	// y luego comparo entre estas para simplificar
	// Obs: en caso de 2 jugadores del mismo que tiraron
	// una carta con el mismo poder -> se queda con la Primera
	// es decir, la que "gana de mano"
	maxPoder := map[Equipo]int{Rojo: -1, Azul: -1}
	max := map[Equipo]*CartaTirada{Rojo: nil, Azul: nil}
	tiradas := p.Ronda.GetManoActual().CartasTiradas

	for i, tirada := range tiradas {
		poder := tirada.Carta.CalcPoder(p.Ronda.Muestra)
		equipo := p.Ronda.Manojo(tirada.Jugador).Jugador.Equipo
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

	// FIX: o simplemente cuando un equipo entero quedo con 0 jugadores "en pie"
	noSeLlegoATirarNingunaCarta := len(p.Ronda.GetManoActual().CartasTiradas) == 0
	seFueronTodos := p.Ronda.CantJugadoresEnJuego[Rojo] == 0 ||
		p.Ronda.CantJugadoresEnJuego[Azul] == 0

	if noSeLlegoATirarNingunaCarta || seFueronTodos {

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
			mano.Ganador = p.Ronda.Manojos[0].Jugador.ID
		} else {
			mano.Ganador = p.Ronda.Manojos[1].Jugador.ID
		}

		// todo: MENSAJE ACAAAAAAA!!!

		// fmt.Printf("La %s mano la gano el equipo %s\n",
		// 	strings.ToLower(p.Ronda.ManoEnJuego.String()),
		// 	equipoGanador.String())

	} else if esParda {
		mano.Resultado = Empardada
		mano.Ganador = ""

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest("ALL"),
				enco.LaManoResultaParda{},
			))
		}

		// no se cambia el turno

	} else {

		// esto quedo arreglado en la funcion IrseAlMazo.Ok para evitar que pueda
		// llegar hasta aca

		/*
			caso especial:
				2 jugadores, 1 tiro carta y enseguida se fue al mazo

			para 4 o 6 jugadores:
				si era el primero de mi equipo y de todos en tirar:
				todos los de mi equipo se van
				yo tiro y me voy (sin dejar chance que los otros tiren)
				gana mi equipo

			 -> no se puede ir al mazo si mi equipo llego a tirar carta y los otros no llegaron a tirar al menos una carta
		*/

		var tiradaGanadora *CartaTirada

		if maxPoder[Rojo] > maxPoder[Azul] {
			tiradaGanadora = max[Rojo]
			mano.Resultado = GanoRojo
		} else {
			tiradaGanadora = max[Azul]
			mano.Resultado = GanoAzul
		}

		// el turno pasa a ser el del mano.ganador
		// pero se setea despues de evaluar la ronda
		mano.Ganador = p.Ronda.Manojo(tiradaGanadora.Jugador).Jugador.ID
		// la variable tiradaGAnadora la uso solo para almacenar el jugador
		// despues le pido el id, y depues se lo vluevlo a preguntar.
		// esta al pedo

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest("ALL"),
				enco.ManoGanada{
					Autor: mano.Ganador,
					Valor: int(p.Ronda.ManoEnJuego),
				},
			))
		}

	}

	// se termino la ronda?
	empiezaNuevaRonda, pkt22 := p.EvaluarRonda()

	if p.Verbose {
		pkts2 = append(pkts2, pkt22...)
	}

	// cuando termina la mano (y no se empieza una ronda) -> cambia de TRUNO
	// cuando termina la ronda -> cambia de MANO
	// para usar esto, antes se debe primero incrementar el turno
	// incremento solo si no se empezo una nueva ronda

	return empiezaNuevaRonda, pkts2
}

// EvaluarRonda tener siempre en cuenta que evaluar la ronda es sinonimo de
// evaluar el truco
// se acabo la ronda?
// si se empieza una ronda nueva -> retorna true
// si no se termino la ronda 	 -> retorna false
func (p *Partida) EvaluarRonda() (bool, []enco.Envelope) {

	/*
		TENER EN CUENTA:
		===============
		el enum p.Ronda.Mano.Resultado \in {GanoRojo,GanoAzul,Empardada}
		no me dice el resultado per se,
		sino:
			noSeSabe sii (no esta empardada) & (ganador == nil)
		por default dice "ganoRojo"
	*/

	var pkts2 []enco.Envelope

	// la ronda continua...

	/* A MENOS QUE SE HAYAN IDO TODOS EN LA PRIMERA MANO!!! */
	hayJugadoresRojo := p.Ronda.CantJugadoresEnJuego[Rojo] > 0
	hayJugadoresAzul := p.Ronda.CantJugadoresEnJuego[Azul] > 0
	hayJugadoresEnAmbos := hayJugadoresRojo && hayJugadoresAzul
	primeraMano := p.Ronda.ManoEnJuego == Primera

	// o bien que en la primera mano hayan cantado truco y uno no lo quizo
	manoActual := p.Ronda.ManoEnJuego.ToInt() - 1
	elTrucoNoTuvoRespuesta := p.Ronda.Truco.Estado.esTrucoRespondible()
	noFueParda := p.Ronda.Manos[manoActual].Resultado != Empardada
	estaManoYaTieneGanador := noFueParda && p.Ronda.Manos[manoActual].Ganador != ""
	elTrucoFueNoQuerido := elTrucoNoTuvoRespuesta && estaManoYaTieneGanador

	elTrucoFueQuerido := !elTrucoFueNoQuerido

	noSeAcabo := (primeraMano && hayJugadoresEnAmbos && elTrucoFueQuerido)
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
			cantManosGanadas[p.Ronda.Manojo(mano.Ganador).Jugador.Equipo]++
		}
	}

	hayEmpate := cantManosGanadas[Rojo] == cantManosGanadas[Azul]
	pardaPrimera := p.Ronda.Manos[0].Resultado == Empardada
	pardaSegunda := p.Ronda.Manos[1].Resultado == Empardada
	pardaTercera := p.Ronda.Manos[2].Resultado == Empardada
	seEstaJugandoLaSegunda := p.Ronda.ManoEnJuego == Segunda

	noSeAcaboAun := seEstaJugandoLaSegunda && hayEmpate && hayJugadoresEnAmbos && !elTrucoFueNoQuerido

	if noSeAcaboAun {
		return false, nil
	}

	// caso particular:
	// no puedo definir quien gano si la seguna mano no tiene definido un resultado
	noEstaEmpardada := p.Ronda.Manos[Segunda].Resultado != Empardada
	noTieneGanador := p.Ronda.Manos[Segunda].Ganador == ""
	segundaManoIndefinida := noEstaEmpardada && noTieneGanador
	// tengo que diferenciar si vengo de: TirarCarta o si vengo de un no quiero:
	// si viniera de un TirarCarta -> en la mano actual (o la anterior)? la ultima carta tirada pertenece al turno actual
	n := len(p.Ronda.Manos[p.Ronda.ManoEnJuego].CartasTiradas)
	actual := p.Ronda.GetElTurno().Jugador.ID
	mix := p.Ronda.ManoEnJuego
	ultimaCartaTiradaPerteneceAlTurnoActual := n > 0 && p.Ronda.Manos[mix].CartasTiradas[n-1].Jugador == actual
	vengoDeTirarCarta := ultimaCartaTiradaPerteneceAlTurnoActual
	if segundaManoIndefinida && hayJugadoresEnAmbos && vengoDeTirarCarta {
		return false, nil
	}

	// hay ganador -> ya se que al final voy a retornar un true
	var ganador string = ""

	if !hayJugadoresEnAmbos { // caso particular: todos abandonaron

		// enonces como antes paso por evaluar mano
		// y seteo a ganador de la ultima mano jugada (la "actual")
		// al equipo que no abandono -> lo sacao de ahi

		// caso particular: la mano resulto "empardada pero uno abandono"
		if noFueParda && estaManoYaTieneGanador {
			ganador = p.Ronda.GetManoActual().Ganador
		} else {
			// el ganador es el primer jugador que no se haya ido al mazo del equipo
			// que sigue en pie
			equipoGanador := Azul
			if !hayJugadoresAzul {
				equipoGanador = Rojo
			}
			for _, m := range p.Ronda.Manojos {
				if !m.SeFueAlMazo && m.Jugador.Equipo == equipoGanador {
					ganador = m.Jugador.ID
					break
				}
			}
		}

		// primero el caso clasico: un equipo gano 2 o mas manos
	} else if cantManosGanadas[Rojo] >= 2 {
		// agarro cualquier manojo de los rojos
		// o bien es la Primera o bien la Segunda
		if p.Ronda.Manojo(p.Ronda.Manos[0].Ganador).Jugador.Equipo == Rojo {
			ganador = p.Ronda.Manos[0].Ganador
		} else {
			ganador = p.Ronda.Manos[1].Ganador
		}
	} else if cantManosGanadas[Azul] >= 2 {
		// agarro cualquier manojo de los azules
		// o bien es la Primera o bien la Segunda
		if p.Ronda.Manojo(p.Ronda.Manos[0].Ganador).Jugador.Equipo == Azul {
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
			ganador = p.Ronda.GetElMano().Jugador.ID
		}

	}

	/************************************************/

	// ya sabemos el ganador ahora es el
	// momento de sumar los puntos del truco
	var totalPts int = 0

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

	if !hayJugadoresEnAmbos {
		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest("ALL"),
				enco.RondaGanada{
					Autor: ganador,
					Razon: enco.SeFueronAlMazo,
				},
				// `La ronda ha sido ganada por el equipo %s. +%v puntos para el equipo %s por el %s ganado`
			))
		}

	} else if elTrucoNoTuvoRespuesta {

		ganador = p.Ronda.Truco.CantadoPor

		var razon enco.Razon
		switch p.Ronda.Truco.Estado {
		case TRUCO:
			razon = enco.TrucoNoQuerido
		case RETRUCO:
			razon = enco.TrucoNoQuerido
		case VALE4:
			razon = enco.TrucoNoQuerido
		}

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest("ALL"),
				enco.RondaGanada{
					Autor: ganador,
					Razon: razon,
				},
			))
		}

	} else {

		var razon enco.Razon
		switch p.Ronda.Truco.Estado {
		case TRUCO:
			razon = enco.TrucoQuerido
		case RETRUCO:
			razon = enco.TrucoQuerido
		case VALE4:
			razon = enco.TrucoQuerido
		}

		if p.Verbose {
			pkts2 = append(pkts2, enco.Pkt(
				enco.Dest("ALL"),
				enco.RondaGanada{
					Autor: ganador,
					Razon: razon,
				},
			))
		}

	}

	p.SumarPuntos(p.Ronda.Manojo(ganador).Jugador.Equipo, totalPts)

	if p.Verbose {
		pkts2 = append(pkts2, enco.Pkt(
			enco.Dest("ALL"),
			enco.SumaPts{
				Autor:  ganador,
				Razon:  enco.TrucoQuerido,
				Puntos: totalPts,
			},
		))
	}

	return true, pkts2 // porque se empezo una nueva ronda
}

// NuevaRonda .
func (p *Partida) NuevaRonda(elMano JIX) {
	p.Ronda.nuevaRonda(elMano)
}

// MarshalJSON retorna la partida en formato json
func (p *Partida) MarshalJSON() ([]byte, error) {
	// return json.Marshal(p)
	// return []byte("hola"), nil
	return json.Marshal(*p)
}

// FromJSON no se pudo implementar con UnmarshalJSON (entraba en loop)
func (p *Partida) FromJSON(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	// default lim envido
	if p.LimiteEnvido == 0 {
		p.LimiteEnvido = 4
	}

	// de estos se encarga el Unmarshal:
	// Manojos: make([]Manojo, cantJugadores)
	// Manos:   make([]Mano, 3)

	// este lo tengo que hacer a mano porque no esta en el JSON
	p.Ronda.MIXS = make(map[string]int)
	p.Ronda.indexarManojos()

	p.Ronda.CachearFlores(false) // sin reset

	return nil
}

func Parse(data string, verbose bool) (*Partida, error) {
	p := new(Partida)
	err := p.FromJSON([]byte(data))
	p.Verbose = verbose
	return p, err
}

func cheepCopy(p *Partida) *Partida {
	copia := *p
	copia.Ronda.Manojos = make([]Manojo, len(p.Ronda.Manojos))
	copy(copia.Ronda.Manojos, p.Ronda.Manojos)

	return &copia
}

// Perspectiva retorna una representacion en json de la PerspectivaCacheFlor que tiene
// el jugador `j` de la partida
func (p *Partida) Perspectiva(j string) (*Partida, error) {
	manojo := p.Manojo(j)
	return p.PerspectivaCacheFlor(manojo), nil
}

// PerspectivaCacheFlor cache las flores (no re-calcula las flores)
func (p *Partida) PerspectivaCacheFlor(manojo *Manojo) *Partida {
	copia := cheepCopy(p)

	// oculto las caras no tiradas de los manojos que no son de su equipo
	for i, m := range copia.Ronda.Manojos {
		noEsDeSuEquipo := m.Jugador.Equipo != manojo.Jugador.Equipo
		if noEsDeSuEquipo {
			// oculto solo las cartas que no tiro
			for j, tirada := range m.Tiradas {
				if !tirada {
					copia.Ronda.Manojos[i].Cartas[j] = nil
				}
			}
		}
	}

	return copia
}

func (p *Partida) Manojo(j string) *Manojo {
	m := p.Ronda.Manojo(j)
	if m != nil {
		return m
	}
	// segundo intento
	m = p.Ronda.Manojo(strings.Title(j))
	if m != nil {
		return m
	}
	return nil
}

/* metodos de manipulacion */

// TirarCarta tira una carta
func (p *Partida) TirarCarta(manojo *Manojo, idx int) {
	manojo.Tiradas[idx] = true
	manojo.UltimaTirada = idx
	carta := manojo.Cartas[idx]
	tirada := CartaTirada{manojo.Jugador.ID, *carta}
	p.Ronda.GetManoActual().agregarTirada(tirada)
}

// dada una partida, intercambia los puestos 2 a 2
func (p *Partida) Swap() {
	n := len(p.Ronda.Manojos)
	for i := 0; i < n/2; i++ {
		offset := i * 2
		aux := p.Ronda.Manojos[offset].Jugador
		p.Ronda.Manojos[offset].Jugador = p.Ronda.Manojos[offset+1].Jugador
		p.Ronda.Manojos[offset+1].Jugador = aux
	}
	// tengo que resetear los MIXS
	p.Ronda.indexarManojos()
}

// NuevaPartida crea una nueva Partida
func NuevaPartida(

	puntuacion Puntuacion,
	equipoAzul,
	equipoRojo []string,
	limiteEnvido int,
	verbose bool,

) (*Partida, error) {

	mismaCantidadDeJugadores := len(equipoRojo) == len(equipoAzul)
	cantJugadores := len(equipoRojo) + len(equipoAzul)
	cantidadCorrecta := cantJugadores == 2 || cantJugadores == 4 || cantJugadores == 6
	ok := mismaCantidadDeJugadores && cantidadCorrecta

	if !ok {
		return nil, fmt.Errorf(`la cantidad de jugadores no es correcta`)
	}

	p := Partida{
		Puntuacion:   puntuacion,
		Verbose:      verbose,
		LimiteEnvido: limiteEnvido,
	}

	p.Puntajes = map[Equipo]int{
		Rojo: 0,
		Azul: 0,
	}

	p.Ronda = MakeRonda(equipoAzul, equipoRojo)

	return &p, nil
}
