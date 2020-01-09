package truco

import (
	"fmt"
	"regexp"
	"strings"
)

// el envido, la primera o la mentira
// el envido, la primera o la mentira
// el truco, la segunda o el rabón

var (
	quit       chan bool    = make(chan bool, 1)
	wait       chan bool    = make(chan bool, 1)
	sigJugada  chan IJugada = make(chan IJugada, 1)
	sigComando chan string  = make(chan string, 1)
)

// Puntuacion : Enum para el puntaje maximo de la partida
type Puntuacion int

// hasta 15 pts, 20 pts, 30 pts o 40 pts
const (
	a20 Puntuacion = 0
	a30 Puntuacion = 1
	a40 Puntuacion = 2
)

func (pt Puntuacion) toInt() int {
	switch pt {
	case a20:
		return 20
	case a30:
		return 30
	default:
		return 40
	}
}

// Equipo : Enum para el puntaje maximo de la partida
type Equipo int

// rojo o azul
const (
	Azul Equipo = 0
	Rojo Equipo = 1
)

func (e Equipo) String() string {
	if e == Rojo {
		return "Rojo"
	}
	return "Azul"
}

// Partida :
type Partida struct {
	jugadores     []Jugador
	CantJugadores int        `json:"cantJugadores"`
	Puntuacion    Puntuacion `json:"puntuacion"`
	Puntaje       int        `json:"puntaje"`
	Puntajes      [2]int     `json:"puntajes"`
	Ronda         Ronda      `json:"ronda"`
}

// SetSigJugada nexo capa presentacion con capa logica
func (p *Partida) SetSigJugada(cmd string) error {
	// checkeo de sintaxis
	ok := regexp.MustCompile(`^(\w|-)+\s(\w|-)+\n?$`).MatchString(cmd)
	if !ok {
		return fmt.Errorf("Comando incorrecto")
	}

	sigComando <- cmd
	return nil
}

// devuelve solo la siguiente jugada VALIDA
// si no es valida es como si no hubiese pasado nada
func (p *Partida) getSigJugada() IJugada {
	var (
		iJugada IJugada
		valid   bool
	)
	for {
		iJugada, valid = <-sigJugada
		if !valid {
			quit <- true
		} else if iJugada == nil {
			wait <- true
		} else {
			break
		}
	}
	return iJugada
}

func (p *Partida) parseJugada(jugadaStr, jugadorStr string) (IJugada, error) {
	var jugada IJugada

	manojo, err := p.Ronda.getManojoByStr(jugadorStr)
	if err != nil {
		return nil, fmt.Errorf("Usuario %s no encontrado", jugadorStr)
	}

	jugadaStr = strings.ToLower(jugadaStr)

	switch jugadaStr {
	// toques
	case "envido":
		jugada = tocarEnvido{Jugada{autor: manojo}}
	case "real-envido":
		jugada = tocarRealEnvido{Jugada{autor: manojo}}
	case "falta-envido":
		jugada = tocarFaltaEnvido{Jugada{autor: manojo}}

	// cantos
	case "flor":
		jugada = cantarFlor{Jugada{autor: manojo}}
	case "contra-flor":
		jugada = cantarContraFlor{Jugada{autor: manojo}}
	case "contra-flor-al-resto":
		jugada = cantarContraFlorAlResto{Jugada{autor: manojo}}

	// gritos
	case "truco":
		jugada = gritarTruco{Jugada{autor: manojo}}
	case "re-truco":
		jugada = gritarReTruco{Jugada{autor: manojo}}
	case "vale-4":
		jugada = gritarVale4{Jugada{autor: manojo}}

	// respuestas
	case "quiero":
		jugada = responderQuiero{Jugada{autor: manojo}}
	case "no-quiero":
		jugada = responderNoQuiero{Jugada{autor: manojo}}
	case "tiene":
		jugada = responderNoQuiero{Jugada{autor: manojo}}

	// acciones
	case "mazo":
		jugada = irseAlMazo{Jugada{autor: manojo}}
	default:
		return nil, fmt.Errorf("No esxiste esa jugada")
	}

	return jugada, nil
}

func (p *Partida) getMaxPuntaje() int {
	if p.Puntajes[Rojo] > p.Puntajes[Azul] {
		return p.Puntajes[Rojo]
	}
	return p.Puntajes[Azul]
}

// retorna el equipo que va ganando
func (p *Partida) elQueVaGanando() Equipo {
	vaGanandoRojo := p.Puntajes[Rojo] > p.Puntajes[Azul]
	if vaGanandoRojo {
		return Rojo
	}
	return Azul
}

// getPuntuacionMalas devuelve la mitad de la puntuacion
// total jugable durante toda la partida
func (p *Partida) getPuntuacionMalas() int {
	return p.Puntuacion.toInt() / 2
}

// getJugador dado un indice de jugador,
// devuelve su puntero correspondiente
func (p *Partida) getJugador(jIdx JugadorIdx) *Jugador {
	return &p.jugadores[jIdx]
}

// NoAcabada retorna true si la partida acabo
func (p *Partida) NoAcabada() bool {
	return p.getMaxPuntaje() < p.Puntuacion.toInt()
}

func (p *Partida) elChico() int {
	return p.Puntuacion.toInt() / 2
}

// retorna true si `e` esta en malas
func (p *Partida) estaEnMalas(e Equipo) bool {
	return p.Puntajes[e] < p.elChico()
}

// retorna la cantidad de puntos que le corresponderian
// a `ganadorDelEnvite` si hubiese ganado un "Contra flor al resto"
// sin tener en cuenta los puntos acumulados de envites anteriores
func (p *Partida) calcPtsContraFlorAlResto(ganadorDelEnvite Equipo) int {
	return p.calcPtsFaltaEnvido(ganadorDelEnvite)
}

// retorna la cantidad de puntos que corresponden al Falta-Envido
func (p *Partida) calcPtsFaltaEnvido(ganadorDelEnvite Equipo) int {
	// si el que va ganando:
	// 		esta en Malas -> el ganador del envite (`ganadorDelEnvite`) gana el chico
	// 		esta en Buenas -> el ganador del envite (`ganadorDelEnvite`) gana lo que le falta al maximo para ganar la ronda

	if p.estaEnMalas(p.elQueVaGanando()) {
		loQueLeFaltaAlGANADORparaGanarElChico := p.elChico() - p.Puntajes[ganadorDelEnvite]
		return loQueLeFaltaAlGANADORparaGanarElChico
	}
	//else {
	loQueLeFaltaAlQUEvaGANANDOparaGanarElChico := p.Puntuacion.toInt() - p.Puntajes[p.elQueVaGanando()]
	return loQueLeFaltaAlQUEvaGANANDOparaGanarElChico
	//}

}

// retorna true si termino la partida
func (p *Partida) sumarPuntos(e Equipo, totalPts int) bool {
	p.Puntajes[e] += totalPts
	if p.NoAcabada() {
		return false
	}
	return true
}

// evalua todas las cartas y decide que equipo gano
// de ese ganador se setea el siguiente turno
func (p *Partida) evaluarMano() {
	// cual es la tirada-carta que gano la mano?
	// ojo que puede salir parda
	// para ello primero busco las maximas de cada equipo
	// y luego comparo entre estas para simplificar
	// Obs: en caso de 2 jugadores del mismo que tiraron
	// una carta con el mismo poder -> se queda con la primera
	// es decir, la que "gana de mano"
	maxPoder := [2]int{-1, -1}
	var max [2]tirarCarta
	tiradas := p.Ronda.getManoActual().cartasTiradas

	for _, tirada := range tiradas {
		poder := tirada.Carta.calcPoder(p.Ronda.Muestra)
		equipo := tirada.autor.Jugador.Equipo
		if poder > maxPoder[equipo] {
			maxPoder[equipo] = poder
			max[equipo] = tirada
		}
	}

	mano := p.Ronda.getManoActual()
	esParda := maxPoder[Rojo] == maxPoder[Azul]
	if esParda {
		mano.resultado = Empardada
		mano.ganador = nil
		fmt.Printf("La Mano resulta parda")
		// no se cambia el turno

	} else {
		var tiradaGanadora tirarCarta

		if maxPoder[Rojo] > maxPoder[Azul] {
			tiradaGanadora = max[Rojo]
			mano.resultado = GanoRojo
		} else {
			tiradaGanadora = max[Azul]
			mano.resultado = GanoAzul
		}

		// el turno pasa a ser el del mano.ganador
		// pero se setea despues de evaluar la ronda
		mano.ganador = tiradaGanadora.autor
		fmt.Printf("La Mano la gano %s (equipo %s)",
			mano.ganador.Jugador.Nombre, mano.ganador.Jugador.Equipo.String())
	}

	// se termino la ronda?
	var empiezaNuevaRonda bool = false
	if p.Ronda.ManoEnJuego >= segunda {
		empiezaNuevaRonda = p.evaluarRonda()
	}

	// cuando termina la mano (y no se empieza una ronda) -> cambia de TRUNO
	// cuando termina la ronda -> cambia de MANO
	// para usar esto, antes se debe primero incrementar el turno
	// incremento solo si no se empezo una nueva ronda
	if !empiezaNuevaRonda {
		p.Ronda.ManoEnJuego++
		p.Ronda.nextTurnoPosMano()
	}
}

// se acabo la ronda?
// si se empieza una ronda nueva -> retorna true
// si no se termino la ronda 	 -> retorna false
func (p *Partida) evaluarRonda() bool {
	imposibleQueSeHayaAcabado := p.Ronda.ManoEnJuego == primera
	if imposibleQueSeHayaAcabado {
		return false
	}

	// de aca en mas ya se que hay al menos 2 manos jugadas
	// asi que es seguro acceder a los indices 0 y 1 en:
	// p.Ronda.manos[0] & p.Ronda.manos[1]

	cantManosGanadas := [2]int{0, 0}
	for i := 0; i < p.Ronda.ManoEnJuego.toInt(); i++ {
		mano := p.Ronda.Manos[i]
		if mano.resultado != Empardada {
			cantManosGanadas[mano.ganador.Jugador.Equipo]++
		}
	}

	hayEmpate := cantManosGanadas[Rojo] == cantManosGanadas[Azul]
	pardaPrimera := p.Ronda.Manos[0].resultado == Empardada
	pardaSegunda := p.Ronda.Manos[1].resultado == Empardada
	pardaTercera := p.Ronda.Manos[2].resultado == Empardada
	seEstaJugandoLaSegunda := p.Ronda.ManoEnJuego == segunda
	noSeAcaboAun := seEstaJugandoLaSegunda && hayEmpate

	if noSeAcaboAun {
		return false
	}

	// hay ganador -> ya se que al final voy a retornar un true
	var ganador *Manojo

	// primero el caso clasico: un equipo gano 2 o mas manos
	if cantManosGanadas[Rojo] >= 2 {
		// agarro cualquier manojo de los rojos
		// o bien es la primera o bien la segunda
		if p.Ronda.Manos[0].ganador.Jugador.Equipo == Rojo {
			ganador = p.Ronda.Manos[0].ganador
		} else {
			ganador = p.Ronda.Manos[1].ganador
		}
	} else if cantManosGanadas[Azul] >= 2 {
		// agarro cualquier manojo de los azules
		// o bien es la primera o bien la segunda
		if p.Ronda.Manos[0].ganador.Jugador.Equipo == Azul {
			ganador = p.Ronda.Manos[0].ganador
		} else {
			ganador = p.Ronda.Manos[1].ganador
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
			ganador = p.Ronda.Manos[segunda].ganador

		} else if caso2 {
			ganador = p.Ronda.Manos[primera].ganador

		} else if caso3 {
			ganador = p.Ronda.Manos[primera].ganador

		} else if caso4 {
			ganador = p.Ronda.Manos[tercera].ganador

		} else if caso5 {
			ganador = p.Ronda.getElMano()
		}

	}

	fmt.Printf("La ronda ha sido ganada por el equipo %s\n",
		ganador.Jugador.Equipo)

	// ya sabemos el ganador ahora es el
	// momento de sumar los puntos del truco
	var totalPts int = 0

	switch p.Ronda.Truco.estado {
	case NOCANTADO:
		totalPts = 1
	case TRUCOQUERIDO:
		totalPts = 2
	case RETRUCOQUERIDO:
		totalPts = 3
	default: // el vale 4
		totalPts = 4
	}

	terminoLaPartida := p.sumarPuntos(ganador.Jugador.Equipo, totalPts)

	if !terminoLaPartida {
		// ahora se deberia de incrementar el mano
		// y ser el turno de este
		sigMano := p.Ronda.getSigMano()
		p.nuevaRonda(sigMano)
	} else {
		p.byeBye()
	}

	return true // porque se empezo una nueva ronda
}

func (p *Partida) byeBye() {
	if !p.NoAcabada() {
		fmt.Printf("Se acabo la partida! el ganador fue el equipo %s\n\n",
			p.elQueVaGanando().String())
		fmt.Printf("BYE BYE!")
	}
}

func (p *Partida) nuevaRonda(elMano JugadorIdx) {
	fmt.Println("Empieza una nueva ronda")
	p.Ronda = nuevaRonda(p.jugadores, elMano)
	fmt.Printf("La mano y el turno es %s\n", p.Ronda.getElMano().Jugador.Nombre)
}

// NuevaPartida retorna nueva partida; error si hubo
func NuevaPartida(puntuacion Puntuacion, equipoAzul, equipoRojo []string) (*Partida, error) {

	mismaCantidadDeJugadores := len(equipoRojo) == len(equipoAzul)
	cantJugadores := len(equipoRojo) + len(equipoAzul)
	cantidadCorrecta := contains([]int{2, 4, 6}, cantJugadores) // puede ser 2, 4 o 6
	ok := mismaCantidadDeJugadores && cantidadCorrecta

	if !ok {
		return nil, fmt.Errorf(`No es posible responderle a la propuesta de tu mismo equipo`)
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

	p := Partida{
		Puntuacion:    puntuacion,
		Puntaje:       0,
		CantJugadores: cantJugadores,
		jugadores:     jugadores,
	}

	p.Puntajes[Rojo] = 0
	p.Puntajes[Azul] = 0

	elMano := JugadorIdx(0)
	p.nuevaRonda(elMano)

	go func() {
		for {
			sjugada := p.getSigJugada()
			sjugada.hacer(&p)
		}
	}()

	go func() {
		for {
			select {
			// este canal agarra solo los comandos en forma de string
			// luego lo pasa al otro canal de jugadas ya aceptadas
			// en la que espera la parte interna del codigo
			case cmd := <-sigComando:
				switch cmd {
				case "__TERMINAR__":
					close(sigJugada)
				case "__WAIT__":
					sigJugada <- nil
				default:
					params := strings.Fields(cmd)
					jugadaStr, jugadorStr := params[1], params[0]
					jugada, err := p.parseJugada(jugadaStr, jugadorStr)
					if err != nil {
						fmt.Println(err.Error())
					} else {
						sigJugada <- jugada
					}
				}

				// case <-p.quit:
				// case <-time.After(1 * time.Second):
				// default:
			}
		}
	}()

	return &p, nil
}

// Terminar espera a que se consuma toda la fila de jugadas
// si se quisiera terminar abruptamente se deberia
// usar otro canal tipo `p.quit<-true` y agregarle
// el caso que corresponda al `select{...}`
func (p *Partida) Terminar() {
	sigComando <- "__TERMINAR__"
	<-quit
}

// Esperar espera a que se consuma toda la fila de jugadas
// para continuar; pero sin cerrar ningun canal
func (p *Partida) Esperar() {
	sigComando <- "__WAIT__"
	<-wait
}
