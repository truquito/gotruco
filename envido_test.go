package truco

import (
	"testing"
	// "fmt"
	// "bufio"
	// "os"
)

// todo:
// sinopsis:
// juan toca envido
// pedro grita truco
// que pasa??
func TestTrucoFueraDeLugar(t *testing.T) {
	p := getPartidaCustom1()

	p.ronda.Print()

	// empieza primera ronda
	// empieza primera mano

	// Juan toca envido
	jugada := tocarEnvido{}
	jugada.hacer(&p, &p.jugadores[0])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error("El estado del envido deberia de ser `envido`")
		return
	}

	oops = p.ronda.envido.puntaje != 2
	if oops {
		t.Error("El `puntaje` del envido deberia de ser 2")
		return
	}

	// Pedro responde 'quiero'
	gritarTruco{}.hacer(&p, &p.jugadores[1])
}

// sinopsis:
// juan toca envido
// pedro responde quiero
// juan gana con 12 de envido vs 5 de pedro
func TestEnvidoAceptado(t *testing.T) {
	p := getPartidaCustom1()
	p.ronda.Print()

	// empieza primera ronda
	// empieza primera mano

	// Juan toca envido
	jugada := tocarEnvido{}
	jugada.hacer(&p, &p.jugadores[0])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error("El estado del envido deberia de ser `envido`")
		return
	}

	oops = p.ronda.envido.puntaje != 2
	if oops {
		t.Error("El `puntaje` del envido deberia de ser 2")
		return
	}

	// Pedro responde 'quiero'
	responderQuiero{}.hacer(&p, &p.jugadores[1])

	oops = p.ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado', 
		ya que fue aceptado por Pedro`)
		return
	}

	oops = p.ronda.envido.puntaje != 2
	if oops {
		t.Error(`El puntaje del envido deberia de ser 2`)
		return
	}
}

// sinopsis:
// juan: envido
// pedro: envido
// todo: ??
func TestDobleEnvido(t *testing.T) {
	p := getPartidaCustom1()

	p.ronda.Print()

	// empieza primera ronda
	// empieza primera mano

	// Juan toca envido
	jugada := tocarEnvido{}
	jugada.hacer(&p, &p.jugadores[0])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error("El estado del envido deberia de ser `envido`")
		return
	}

	oops = p.ronda.envido.puntaje != 2
	if oops {
		t.Error("El `puntaje` del envido deberia de ser 2")
		return
	}

	// Pedro redobla el envido
	tocarEnvido{}.hacer(&p, &p.jugadores[1])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error(`El estado del envido deberia de ser 'envido', incluso luego de que
		ambos Juan y Pedro lo hayan tocando`)
		return
	}

	oops = p.ronda.envido.puntaje != 4
	if oops {
		t.Error(`El puntaje del envido deberia ahora de ser '2 + 2 = 4'`)
		return
	}

	responderQuiero{}.hacer(&p, &p.jugadores[0])
}

// sinopsis:
// Juan: envido
// Pedro: envido
// Juan: no quiero
func TestDobleEnvidoNoAceptado(t *testing.T) {
	p := getPartidaCustom1()

	p.ronda.Print()

	// empieza primera ronda
	// empieza primera mano

	// Juan toca envido
	jugada := tocarEnvido{}
	jugada.hacer(&p, &p.jugadores[0])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error("El estado del envido deberia de ser `envido`")
		return
	}

	oops = p.ronda.envido.puntaje != 2
	if oops {
		t.Error("El `puntaje` del envido deberia de ser 2")
		return
	}

	// Pedro redobla el envido
	tocarEnvido{}.hacer(&p, &p.jugadores[1])

	oops = p.ronda.envido.estado != ENVIDO
	if oops {
		t.Error(`El estado del envido deberia de ser 'envido', incluso luego de que
		ambos Juan y Pedro lo hayan tocando`)
		return
	}

	oops = p.ronda.envido.puntaje != 4
	if oops {
		t.Error(`El puntaje del envido deberia ahora de ser '2 + 2 = 4'`)
		return
	}

	// Juan responde 'no quiero'
	responderNoQuiero{}.hacer(&p, &p.jugadores[0])
}

// parte 2

/* #CASOS = 11*2 = 22 */

// CASO I 		Envido	2/1
// CASO II		Real envido	 3/1
// CASO III		Falta envido	 x/1
// CASO IV		Envido, envido	 4/2
// CASO V			Envido y real envido	 5/2
// CASO VI		Envido y falta envido	 x/2
// CASO VII		Real envido y falta envido	x / 3
// CASO VIII 	Envido, envido y falta envido	x / 4
// CASO IX		Envido, envido y real envido	 7/4
// CASO X			Envido, real envido y falta envido	 x/5
// CASO XI		Envido, envido, real envido y falta envido	 x/7
// donde x = Lo que le falta al rival para ganar

/* CONTEXTO */
// - Segunda ronda en juego; primera mano
// - Todos en `malas`: Rojo 4 pts, Azul 3 pts
//
// - Juan
// - Pedro
// - Jacinto (mano & turno)
// - Patricio

func inicializar() Partida {
	p := partida4JugadoresEnvidoTesting
	p.dobleLinking()
	p.ronda.getManoActual().repartidor = p.ronda.elMano - 1
	return p
}

// CASO I (Aceptado) 		Envido	2/1
func TestIAceptado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	responderQuiero{}.hacer(&p, pedro)
	oops = p.puntajes[Rojo] != 4+2
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO I (Rechazado) 		Envido	2/1
func TestIRechazado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	responderNoQuiero{}.hacer(&p, pedro)
	oops = p.puntajes[Rojo] != 4+1
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO II (Aceptado)		Real envido	 3/1
func TestIIAceptado(t *testing.T) {
	p := inicializar()
	tocarRealEnvido{}.hacer(&p, jacinto)
	responderQuiero{}.hacer(&p, pedro)
	oops = p.puntajes[Rojo] != 4+3
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO II (Rechazado)		Real envido	 3/1
func TestIIRechazado(t *testing.T) {
	p := inicializar()
	tocarRealEnvido{}.hacer(&p, jacinto)
	responderNoQuiero{}.hacer(&p, pedro)
	oops = p.puntajes[Rojo] != 4+1
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO III (Aceptado)		Falta envido	 x/1
func TestIIIAceptado(t *testing.T) {
	p := inicializar()
	tocarFaltaEnvido{}.hacer(&p, jacinto)
	responderQuiero{}.hacer(&p, pedro)
	oops = p.puntajes[Rojo] != 4+16
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO III (Rechazado)		Falta envido	 x/1
func TestIIIRechazado(t *testing.T) {
	p := inicializar()
	tocarFaltaEnvido{}.hacer(&p, jacinto)
	responderNoQuiero{}.hacer(&p, pedro)
	oops = p.puntajes[Rojo] != 4+1
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO IV (Aceptado)		Envido, envido	 4/2
func TestIVAceptado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarEnvido{}.hacer(&p, pedro)
	responderQuiero{}.hacer(&p, juan)
	oops = p.puntajes[Rojo] != 4+4
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO IV (Rechazado)		Envido, envido	 4/2
func TestIVRechazado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarEnvido{}.hacer(&p, pedro)
	responderNoQuiero{}.hacer(&p, juan)
	oops = p.puntajes[Azul] != 3+2
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO V (Aceptado)			Envido y real envido	 5/2
func TestVAceptado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarRealEnvido{}.hacer(&p, patricio)
	responderQuiero{}.hacer(&p, jacinto)
	oops = p.puntajes[Rojo] != 4+5
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO V (Rechazado)			Envido y real envido	 5/2
func TestVRechazado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarRealEnvido{}.hacer(&p, pedro)
	responderNoQuiero{}.hacer(&p, jacinto)
	oops = p.puntajes[Azul] != 3+2
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO VI (Aceptado)		Envido y falta envido	 x/2
func TestVIAceptado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarFaltaEnvido{}.hacer(&p, pedro)
	responderQuiero{}.hacer(&p, jacinto)
	oops = p.puntajes[Rojo] != 4+16
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO VI (Rechazado)		Envido y falta envido	 x/2
func TestVIRechazado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarFaltaEnvido{}.hacer(&p, pedro)
	responderNoQuiero{}.hacer(&p, jacinto)
	oops = p.puntajes[Azul] != 3+1
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO VII (Aceptado)		Real envido y falta envido	x / 3
func TestVIIAceptado(t *testing.T) {
	p := inicializar()
	tocarRealEnvido{}.hacer(&p, jacinto)
	tocarFaltaEnvido{}.hacer(&p, pedro)
	responderQuiero{}.hacer(&p, juan)
	oops = p.puntajes[Rojo] != 4+16
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO VII (Rechazado)		Real envido y falta envido	x / 3
func TestVIIRechazado(t *testing.T) {
	p := inicializar()
	tocarRealEnvido{}.hacer(&p, jacinto)
	tocarFaltaEnvido{}.hacer(&p, pedro)
	responderNoQuiero{}.hacer(&p, juan)
	oops = p.puntajes[Azul] != 3+3
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO VIII (Aceptado) 	Envido, envido y falta envido	x / 4
func TestVIIIAceptado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarEnvido{}.hacer(&p, patricio)
	tocarFaltaEnvido{}.hacer(&p, juan)
	responderQuiero{}.hacer(&p, patricio)
	oops = p.puntajes[Rojo] != 4+16
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO VIII (Rechazado) 	Envido, envido y falta envido	x / 4
func TestVIIIRechazado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarEnvido{}.hacer(&p, patricio)
	tocarFaltaEnvido{}.hacer(&p, juan)
	responderNoQuiero{}.hacer(&p, patricio)
	oops = p.puntajes[Rojo] != 4+4
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO IX (Aceptado)		Envido, envido y real envido	 7/4
func TestIXAceptado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarEnvido{}.hacer(&p, pedro)
	tocarRealEnvido{}.hacer(&p, juan)
	responderQuiero{}.hacer(&p, pedro)
	oops = p.puntajes[Rojo] != 4+7
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO IX (Rechazado)		Envido, envido y real envido	 7/4
func TestIXRechazado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarEnvido{}.hacer(&p, pedro)
	tocarRealEnvido{}.hacer(&p, juan)
	responderNoQuiero{}.hacer(&p, pedro)
	oops = p.puntajes[Rojo] != 4+4
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO X (Aceptado)			Envido, real envido y falta envido	 x/5
func TestXAceptado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarRealEnvido{}.hacer(&p, pedro)
	tocarFaltaEnvido{}.hacer(&p, juan)
	responderQuiero{}.hacer(&p, pedro)
	oops = p.puntajes[Rojo] != 4+16
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO X (Rechazado)			Envido, real envido y falta envido	 x/5
func TestXRechazado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarRealEnvido{}.hacer(&p, pedro)
	tocarFaltaEnvido{}.hacer(&p, juan)
	responderNoQuiero{}.hacer(&p, pedro)
	oops = p.puntajes[Rojo] != 4+5
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO XI (Aceptado)		Envido, envido, real envido y falta envido	 x/7
func TestXIAceptado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarEnvido{}.hacer(&p, pedro)
	tocarRealEnvido{}.hacer(&p, juan)
	tocarFaltaEnvido{}.hacer(&p, patricio)
	responderQuiero{}.hacer(&p, juan)
	oops = p.puntajes[Rojo] != 4+16
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// CASO XI (Rechazado)		Envido, envido, real envido y falta envido	 x/7
func TestXIRechazado(t *testing.T) {
	p := inicializar()
	tocarEnvido{}.hacer(&p, jacinto)
	tocarEnvido{}.hacer(&p, pedro)
	tocarRealEnvido{}.hacer(&p, juan)
	tocarFaltaEnvido{}.hacer(&p, patricio)
	responderNoQuiero{}.hacer(&p, juan)
	oops = p.puntajes[Azul] != 3+7
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

// parte 3

/* Tests de Youtube */
func TestYTEnvidoCalc(t *testing.T) {
	p := partidaYT1
	p.dobleLinking()
	p.ronda.getManoActual().repartidor = 5

	expected := []int{26, 20, 28, 25, 33, 27}
	for i, jugador := range p.jugadores {
		got := jugador.manojo.calcularEnvido(p.ronda.muestra)
		oops = expected[i] != got
		if oops {
			t.Errorf(
				`El resultado del envido del jugador %s es incorrecto.
				\nEXPECTED: %v
				\nGOT: %v`,
				jugador.nombre, expected[i], got)
			return
		}
	}

}

func TestYTEnvidoI(t *testing.T) {
	p := partidaYT1
	p.dobleLinking()
	p.ronda.getManoActual().repartidor = 5

	tocarEnvido{}.hacer(&p, D)
	responderQuiero{}.hacer(&p, C)
	oops = p.puntajes[Rojo] != 4+2
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}

func TestYTEnvidoCalcII(t *testing.T) {
	p := partidaYT2
	p.dobleLinking()
	p.ronda.getManoActual().repartidor = 5

	expected := []int{21, 23, 23, 30, 30, 31}
	for i, jugador := range p.jugadores {
		got := jugador.manojo.calcularEnvido(p.ronda.muestra)
		oops = expected[i] != got
		if oops {
			t.Errorf(
				`El resultado del envido del jugador %s es incorrecto.
				\nEXPECTED: %v
				\nGOT: %v`,
				jugador.nombre, expected[i], got)
			return
		}
	}
}

func TestYTEnvidoII(t *testing.T) {
	p := partidaYT2
	p.dobleLinking()
	p.ronda.getManoActual().repartidor = 5

	tocarEnvido{}.hacer(&p, D)
	responderQuiero{}.hacer(&p, C)
	oops = p.puntajes[Azul] != 3+2
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}
}
