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
	CantJugadores int            `json:"cantJugadores"`
	Puntuacion    Puntuacion     `json:"puntuacion"`
	Puntajes      map[Equipo]int `json:"puntajes"`
	Ronda         Ronda          `json:"ronda"`
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
	return p.CantJugadores == 2
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
		p.Ronda.Envite.CantadoPor = m

	} else { // no se habia jugado aun
		p.Ronda.Envite.CantadoPor = m
		p.Ronda.Envite.Estado = ENVIDO
		p.Ronda.Envite.Puntaje = 2
	}
}

// TocarRealEnvido ..
func (p *Partida) TocarRealEnvido(m *Manojo) {
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
func (p *Partida) TocarFaltaEnvido(m *Manojo) {
	p.Ronda.Envite.Estado = FALTAENVIDO
	p.Ronda.Envite.CantadoPor = m
}

// CantarFlor ..
func (p *Partida) CantarFlor(m *Manojo) {
	yaEstabamosEnFlor := p.Ronda.Envite.Estado == FLOR
	p.Ronda.Envite.Estado = FLOR

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
func (p *Partida) CantarContraFlor(m *Manojo) {
	p.Ronda.Envite.Estado = CONTRAFLOR
	p.Ronda.Envite.CantadoPor = m
	// ahora la flor pasa a jugarse por 4 puntos
	p.Ronda.Envite.Puntaje = 4
}

// CantarContraFlorAlResto ..
func (p *Partida) CantarContraFlorAlResto(m *Manojo) {
	p.Ronda.Envite.Estado = CONTRAFLORALRESTO
	p.Ronda.Envite.CantadoPor = m
	// ahora la flor pasa a jugarse por 4 puntos
	p.Ronda.Envite.Puntaje = 4 // <- eso es al pedo, es independiente
}

// GritarTruco ..
func (p *Partida) GritarTruco(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m
	p.Ronda.Truco.Estado = TRUCO
	// p.Ronda.Envite.Estado = DESHABILITADO // <-- esto esta mal
}

// QuererTruco incrementa el estado del truco a querido segun corresponda
// y setea a m como el que lo canto
func (p *Partida) QuererTruco(m *Manojo) {
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
func (p *Partida) GritarReTruco(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m
	p.Ronda.Truco.Estado = RETRUCO
	// p.Ronda.Envite.Estado = DESHABILITADO // <-- esto esta mal
}

// GritarVale4 ..
func (p *Partida) GritarVale4(m *Manojo) {
	p.Ronda.Truco.CantadoPor = m
	p.Ronda.Truco.Estado = VALE4
	// p.Ronda.Envite.Estado = DESHABILITADO // <-- esto esta mal
}

// IrAlMazo manda el manojo al mazo
func (p *Partida) IrAlMazo(manojo *Manojo) {
	manojo.SeFueAlMazo = true
	equipoDelJugador := manojo.Jugador.Equipo
	p.Ronda.CantJugadoresEnJuego[equipoDelJugador]--
}

// EvaluarMano evalua todas las cartas y decide que equipo gano
// de ese ganador se setea el siguiente turno
func (p *Partida) EvaluarMano() (bool, []*enco.Packet) {

	var pkts []*enco.Packet

	// cual es la tirada-carta que gano la mano?
	// ojo que puede salir parda
	// para ello primero busco las maximas de cada equipo
	// y luego comparo entre estas para simplificar
	// Obs: en caso de 2 jugadores del mismo que tiraron
	// una carta con el mismo poder -> se queda con la Primera
	// es decir, la que "gana de mano"
	maxPoder := map[Equipo]int{Rojo: -1, Azul: -1}
	max := map[Equipo]*cartaTirada{Rojo: nil, Azul: nil}
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

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.LaManoResultaParda),
		))

		// no se cambia el turno

	} else {
		var tiradaGanadora *cartaTirada

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

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.ManoGanada, mano.Ganador.Jugador.ID, int(p.Ronda.ManoEnJuego)),
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
func (p *Partida) EvaluarRonda() (bool, []*enco.Packet) {

	/*
		TENER EN CUENTA:
		===============
		el enum p.Ronda.Mano.Resultado \in {GanoRojo,GanoAzul,Empardada}
		no me dice el resultado per se,
		sino:
			noSeSabe sii (no esta empardada) & (ganador == nil)
		por default dice "ganoRojo"
	*/

	var pkts []*enco.Packet

	// la ronda continua...

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
			cantManosGanadas[mano.Ganador.Jugador.Equipo]++
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
	noTieneGanador := p.Ronda.Manos[Segunda].Ganador == nil
	segundaManoIndefinida := noEstaEmpardada && noTieneGanador
	// tengo que diferenciar si vengo de: TirarCarta o si vengo de un no quiero:
	// si viniera de un TirarCarta -> en la mano actual (o la anterior)? la ultima carta tirada pertenece al turno actual
	n := len(p.Ronda.Manos[p.Ronda.ManoEnJuego].CartasTiradas)
	actual := p.Ronda.GetElTurno().Jugador.ID
	mix := p.Ronda.ManoEnJuego
	ultimaCartaTiradaPerteneceAlTurnoActual := n > 0 && p.Ronda.Manos[mix].CartasTiradas[n-1].autor.Jugador.ID == actual
	vengoDeTirarCarta := ultimaCartaTiradaPerteneceAlTurnoActual
	if segundaManoIndefinida && hayJugadoresEnAmbos && vengoDeTirarCarta {
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

		var razon enco.Razon
		switch p.Ronda.Truco.Estado {
		case TRUCO:
			razon = enco.TrucoNoQuerido
		case RETRUCO:
			razon = enco.TrucoNoQuerido
		case VALE4:
			razon = enco.TrucoNoQuerido
		}

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.RondaGanada, ganador.Jugador.ID, int(razon)),
			// `La ronda ha sido ganada por el equipo %s. +%v puntos para el equipo %s por el %s no querido`
		))

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

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.RondaGanada, ganador.Jugador.ID, int(razon)),
			// `La ronda ha sido ganada por el equipo %s. +%v puntos para el equipo %s por el %s ganado`
		))
	}

	p.SumarPuntos(ganador.Jugador.Equipo, totalPts)

	pkts = append(pkts, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.SumaPts, ganador.Jugador.ID, enco.TrucoQuerido, totalPts),
	))

	return true, pkts // porque se empezo una nueva ronda
}

// NuevaRonda .
func (p *Partida) NuevaRonda(elMano JugadorIdx) {
	p.Ronda.nuevaRonda(elMano)
}

// DEPRECATED: usar `.Manojo(s)`
// GetManojoByStr ..
// OJO QUE AHORA LAS COMPARACIONES SON CASE INSENSITIVE
// ENTONCES SI EL IDENTIFICADOR Juan == jUaN
// ojo con los kakeos
// todo: esto es ineficiente
// getManojo devuelve el puntero al manojo,
// dado un string que identifique al jugador duenio de ese manojo
func (p *Partida) GetManojoByStr(idJugador string) (*Manojo, error) {
	idJugador = strings.ToLower(idJugador)
	for i := range p.Ronda.Manojos {
		idActual := strings.ToLower(p.Ronda.Manojos[i].Jugador.ID)
		esEse := idActual == idJugador
		if esEse {
			return &p.Ronda.Manojos[i], nil
		}
	}
	return nil, fmt.Errorf("jugador `%s` no encontrado", idJugador)
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

	p.Ronda.cachearFlores()
	// cargo los jugadores xq no vienen en el json
	// p.Jugadores = make([]Jugador, p.CantJugadores)
	// for i, m := range p.Ronda.Manojos {
	// 	p.Jugadores[i] = *m.Jugador
	// }

	// cargo los autores de las tiradas de cada una de las 3 manos
	/*
		para cada mano i:
			si i no tiene cartas tiradas break
			j = 0
			para tirada t:
				count = p.cantJugadores
				mientras que j no sea quien tiene la carta t.c:
					j.sig()
					count--
					if count == 0 return err
				t.c.autor = j
				j.sigManojo_NOSE()
	*/
	for mix, mano := range p.Ronda.Manos {
		if len(mano.CartasTiradas) == 0 {
			break
		}
		jix := 0
		for tix, t := range mano.CartasTiradas {
			count := p.CantJugadores
			for {
				if count--; count == -1 {
					return fmt.Errorf("los datos no son consistentes")
				}
				if _, err := p.Ronda.Manojos[jix].GetCartaIdx(t.Carta); err == nil {
					break // bingo
				} else {
					jix = int(p.Ronda.getSig(JugadorIdx(jix))) // no tiene la carta, siguiente
				}
			}
			p.Ronda.Manos[mix].CartasTiradas[tix].autor = &p.Ronda.Manojos[jix]
			jix = int(p.Ronda.getSig(JugadorIdx(jix)))
		}
	}

	return nil
}

func Parse(data string) (*Partida, error) {
	p := new(Partida)
	err := p.FromJSON([]byte(data))
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
	// primero encuentro el jugador
	manojo, err := p.GetManojoByStr(j)
	if err != nil {
		return nil, fmt.Errorf("usuario %s no encontrado", j)
	}

	return p.PerspectivaCacheFlor(manojo), nil
}

// PerspectivaCacheFlor cache las flores
func (p *Partida) PerspectivaCacheFlor(manojo *Manojo) *Partida {
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
func (p *Partida) TirarCarta(manojo *Manojo, idx int) {
	manojo.CartasNoTiradas[idx] = false
	manojo.UltimaTirada = idx
	carta := manojo.Cartas[idx]
	tirada := cartaTirada{manojo, *carta}
	p.Ronda.GetManoActual().agregarTirada(tirada)
}

// NuevaPartida crea una nueva Partida
func NuevaPartida(puntuacion Puntuacion, equipoAzul, equipoRojo []string) (*Partida, error) {

	mismaCantidadDeJugadores := len(equipoRojo) == len(equipoAzul)
	cantJugadores := len(equipoRojo) + len(equipoAzul)
	cantidadCorrecta := Contains([]int{2, 4, 6}, cantJugadores) // puede ser 2, 4 o 6
	ok := mismaCantidadDeJugadores && cantidadCorrecta

	if !ok {
		return nil, fmt.Errorf(`la cantidad de jugadores no es correcta`)
	}
	// paso a crear los jugadores; intercalados
	var jugadores []Jugador
	// para cada rjo que agrego; le agrego tambien su mano
	for i := range equipoRojo {
		// uso como id sus nombres
		nuevoJugadorRojo := Jugador{equipoRojo[i], Rojo}
		nuevoJugadorAzul := Jugador{equipoAzul[i], Azul}
		jugadores = append(jugadores, nuevoJugadorAzul, nuevoJugadorRojo)
	}

	p := Partida{
		Puntuacion:    puntuacion,
		CantJugadores: cantJugadores,
	}

	p.Puntajes = make(map[Equipo]int)
	p.Puntajes[Rojo] = 0
	p.Puntajes[Azul] = 0

	p.Ronda = MakeRonda(equipoAzul, equipoRojo)

	return &p, nil
}
