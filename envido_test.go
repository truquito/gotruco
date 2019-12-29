package truco

import (
	"testing"
)

// Tests:
// Envido	2/1
// Real envido	 3/1
// Falta envido	 x/1
// Envido + envido	 2+2/2+1
// Envido + real envido	 2+3/2+1
// Envido + falta envido	 2+x/2+1
// Real envido + falta envido	3+x / 3+1
// Envido + envido + falta envido	2+2+x / 2+2+1
// Envido + envido + real envido	 2+2+3/2+2+1
// Envido + real envido + falta envido	 2+3+x/2+3+1
// Envido + envido + real envido + falta envido	 2+2+3+x/2+2+3+1

func TestEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // envido: 13
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Copa, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 3},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Envido")
	p.SetSigJugada("Roro Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = p.Ronda.envido.puntaje != 2
	if oops {
		t.Error(`El puntaje del envido deberia de ser 2`)
		return
	}

	oops = !(p.puntajes[Azul] == 2 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 2`)
		return
	}
}

func TestEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // envido: 13
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Copa, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 3},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Envido")
	p.SetSigJugada("Roro No-Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = p.Ronda.envido.puntaje != 1
	if oops {
		t.Error(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 1 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo rojo deberia de ser 1`)
		return
	}
}

func TestRealEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // envido: 13
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Copa, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 3},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Real-Envido")
	p.SetSigJugada("Roro Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 3)
	if oops {
		t.Error(`El puntaje del envido deberia de ser 3`)
		return
	}

	oops = !(p.puntajes[Azul] == 3 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestRealEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // envido: 13
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Copa, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 3},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Real-Envido")
	p.SetSigJugada("Roro No-Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 1)
	if oops {
		t.Error(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 1 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
		return
	}
}

func TestFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // envido: 13
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Copa, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 3},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Falta-Envido")
	p.SetSigJugada("Roro Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 10`)
		return
	}

	oops = !(p.puntajes[Azul] == 10 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 10`)
		return
	}
}

func TestFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // envido: 13
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Copa, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 3},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Falta-Envido")
	p.SetSigJugada("Roro No-Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 1 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
		return
	}
}

func TestEnvidoEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // envido: 13
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Copa, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 3},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Envido")
	p.Esperar()

	oops = p.Ronda.envido.estado != ENVIDO
	if oops {
		t.Error("El estado del envido deberia de ser `envido`")
		return
	}

	oops = p.Ronda.envido.puntaje != 2
	if oops {
		t.Error("El `puntaje` del envido deberia de ser 2")
		return
	}

	p.SetSigJugada("Roro Envido")
	p.Esperar()

	oops = p.Ronda.envido.estado != ENVIDO
	if oops {
		t.Error(`El estado del envido deberia de ser 'envido', incluso luego de que
		ambos Juan y Pedro lo hayan tocando`)
		return
	}

	oops = p.Ronda.envido.puntaje != 4
	if oops {
		t.Error(`El puntaje del envido deberia ahora de ser '2 + 2 = 4'`)
		return
	}

	p.SetSigJugada("Alvaro Quiero")
	p.Esperar()
}

func TestEnvidoEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // envido: 13
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Copa, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 3},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Envido")
	p.SetSigJugada("Roro Envido")
	p.SetSigJugada("Alvaro No-Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 0 && p.puntajes[Rojo] == 2+1)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 0`)
		return
	}
}

// CASO V (Aceptado)			Envido y real envido	 5/2
func TestVAceptado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarRealEnvido{}.hacer(&p, patricio)
	// responderQuiero{}.hacer(&p, jacinto)
	// oops = p.puntajes[Rojo] != 4+5
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO V (Rechazado)			Envido y real envido	 5/2
func TestVRechazado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarRealEnvido{}.hacer(&p, pedro)
	// responderNoQuiero{}.hacer(&p, jacinto)
	// oops = p.puntajes[Azul] != 3+2
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO VI (Aceptado)		Envido y falta envido	 x/2
func TestVIAceptado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarFaltaEnvido{}.hacer(&p, pedro)
	// responderQuiero{}.hacer(&p, jacinto)
	// oops = p.puntajes[Rojo] != 4+16
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO VI (Rechazado)		Envido y falta envido	 x/2
func TestVIRechazado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarFaltaEnvido{}.hacer(&p, pedro)
	// responderNoQuiero{}.hacer(&p, jacinto)
	// oops = p.puntajes[Azul] != 3+1
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO VII (Aceptado)		Real envido y falta envido	x / 3
func TestVIIAceptado(t *testing.T) {
	// p := inicializar()
	// tocarRealEnvido{}.hacer(&p, jacinto)
	// tocarFaltaEnvido{}.hacer(&p, pedro)
	// responderQuiero{}.hacer(&p, juan)
	// oops = p.puntajes[Rojo] != 4+16
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO VII (Rechazado)		Real envido y falta envido	x / 3
func TestVIIRechazado(t *testing.T) {
	// p := inicializar()
	// tocarRealEnvido{}.hacer(&p, jacinto)
	// tocarFaltaEnvido{}.hacer(&p, pedro)
	// responderNoQuiero{}.hacer(&p, juan)
	// oops = p.puntajes[Azul] != 3+3
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO VIII (Aceptado) 	Envido, envido y falta envido	x / 4
func TestVIIIAceptado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarEnvido{}.hacer(&p, patricio)
	// tocarFaltaEnvido{}.hacer(&p, juan)
	// responderQuiero{}.hacer(&p, patricio)
	// oops = p.puntajes[Rojo] != 4+16
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO VIII (Rechazado) 	Envido, envido y falta envido	x / 4
func TestVIIIRechazado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarEnvido{}.hacer(&p, patricio)
	// tocarFaltaEnvido{}.hacer(&p, juan)
	// responderNoQuiero{}.hacer(&p, patricio)
	// oops = p.puntajes[Rojo] != 4+4
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO IX (Aceptado)		Envido, envido y real envido	 7/4
func TestIXAceptado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarEnvido{}.hacer(&p, pedro)
	// tocarRealEnvido{}.hacer(&p, juan)
	// responderQuiero{}.hacer(&p, pedro)
	// oops = p.puntajes[Rojo] != 4+7
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO IX (Rechazado)		Envido, envido y real envido	 7/4
func TestIXRechazado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarEnvido{}.hacer(&p, pedro)
	// tocarRealEnvido{}.hacer(&p, juan)
	// responderNoQuiero{}.hacer(&p, pedro)
	// oops = p.puntajes[Rojo] != 4+4
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO X (Aceptado)			Envido, real envido y falta envido	 x/5
func TestXAceptado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarRealEnvido{}.hacer(&p, pedro)
	// tocarFaltaEnvido{}.hacer(&p, juan)
	// responderQuiero{}.hacer(&p, pedro)
	// oops = p.puntajes[Rojo] != 4+16
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO X (Rechazado)			Envido, real envido y falta envido	 x/5
func TestXRechazado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarRealEnvido{}.hacer(&p, pedro)
	// tocarFaltaEnvido{}.hacer(&p, juan)
	// responderNoQuiero{}.hacer(&p, pedro)
	// oops = p.puntajes[Rojo] != 4+5
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO XI (Aceptado)		Envido, envido, real envido y falta envido	 x/7
func TestXIAceptado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarEnvido{}.hacer(&p, pedro)
	// tocarRealEnvido{}.hacer(&p, juan)
	// tocarFaltaEnvido{}.hacer(&p, patricio)
	// responderQuiero{}.hacer(&p, juan)
	// oops = p.puntajes[Rojo] != 4+16
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// CASO XI (Rechazado)		Envido, envido, real envido y falta envido	 x/7
func TestXIRechazado(t *testing.T) {
	// p := inicializar()
	// tocarEnvido{}.hacer(&p, jacinto)
	// tocarEnvido{}.hacer(&p, pedro)
	// tocarRealEnvido{}.hacer(&p, juan)
	// tocarFaltaEnvido{}.hacer(&p, patricio)
	// responderNoQuiero{}.hacer(&p, juan)
	// oops = p.puntajes[Azul] != 3+7
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

// parte 3

/* Tests de Youtube */
func TestYTEnvidoCalc(t *testing.T) {
	// p := partidaYT1
	// p.Ronda.singleLinking(p.jugadores)
	// p.Ronda.getManoActual().repartidor = 5

	// expected := []int{26, 20, 28, 25, 33, 27}
	// for i, jugador := range p.jugadores {
	// 	got := jugador.manojo.calcularEnvido(p.Ronda.muestra)
	// 	oops = expected[i] != got
	// 	if oops {
	// 		t.Errorf(
	// 			`El resultado del envido del jugador %s es incorrecto.
	// 			\nEXPECTED: %v
	// 			\nGOT: %v`,
	// 			jugador.nombre, expected[i], got)
	// 		return
	// 	}
	// }

}

func TestYTEnvidoI(t *testing.T) {
	// p := partidaYT1
	// p.Ronda.singleLinking(p.jugadores)
	// p.Ronda.getManoActual().repartidor = 5

	// tocarEnvido{}.hacer(&p, D)
	// responderQuiero{}.hacer(&p, C)
	// oops = p.puntajes[Rojo] != 4+2
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}

func TestYTEnvidoCalcII(t *testing.T) {
	// p := partidaYT2
	// p.Ronda.singleLinking(p.jugadores)
	// p.Ronda.getManoActual().repartidor = 5

	// expected := []int{21, 23, 23, 30, 30, 31}
	// for i, jugador := range p.jugadores {
	// 	got := jugador.manojo.calcularEnvido(p.Ronda.muestra)
	// 	oops = expected[i] != got
	// 	if oops {
	// 		t.Errorf(
	// 			`El resultado del envido del jugador %s es incorrecto.
	// 			\nEXPECTED: %v
	// 			\nGOT: %v`,
	// 			jugador.nombre, expected[i], got)
	// 		return
	// 	}
	// }
}

func TestYTEnvidoII(t *testing.T) {
	// p := partidaYT2
	// p.Ronda.singleLinking(p.jugadores)
	// p.Ronda.getManoActual().repartidor = 5

	// tocarEnvido{}.hacer(&p, D)
	// responderQuiero{}.hacer(&p, C)
	// oops = p.puntajes[Azul] != 3+2
	// if oops {
	// 	t.Error("El resultado es incorrecto")
	// 	return
	// }
}
