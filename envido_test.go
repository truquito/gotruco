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
// Envido + envido + real envido	 2+2+3/2+2+1
// Envido + envido + falta envido	2+2+x / 2+2+1
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

func TestEnvidoRealEnvidoQuiero(t *testing.T) {
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
	p.SetSigJugada("Roro Real-Envido")
	p.SetSigJugada("Alvaro Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+3)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 2+3 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 5`)
		return
	}
}

func TestEnvidoRealEnvidoNoQuiero(t *testing.T) {
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
	p.SetSigJugada("Roro Real-Envido")
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
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestEnvidoFaltaEnvidoQuiero(t *testing.T) {
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
	p.SetSigJugada("Roro Falta-Envido")
	p.SetSigJugada("Alvaro Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 2+10 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
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
	p.SetSigJugada("Roro Falta-Envido")
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
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
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
	p.SetSigJugada("Roro Falta-Envido")
	p.SetSigJugada("Alvaro Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 3+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 3+10 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
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
	p.SetSigJugada("Roro Falta-Envido")
	p.SetSigJugada("Alvaro No-Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 3+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 0 && p.puntajes[Rojo] == 3+1)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestEnvidoEnvidoRealEnvidoQuiero(t *testing.T) {
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
	p.SetSigJugada("Alvaro Real-Envido")
	p.SetSigJugada("Roro Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+2+3)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 2+2+3 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestEnvidoEnvidoRealEnvidoNoQuiero(t *testing.T) {
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
	p.SetSigJugada("Alvaro Real-Envido")
	p.SetSigJugada("Roro No-Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 2+2+1 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestEnvidoEnvidoFaltaEnvidoQuiero(t *testing.T) {
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
	p.SetSigJugada("Alvaro Falta-Envido")
	p.SetSigJugada("Roro Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+2+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 2+2+10 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestEnvidoEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
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
	p.SetSigJugada("Alvaro Falta-Envido")
	p.SetSigJugada("Roro No-Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 2+2+1 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestEnvidoRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
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
	p.SetSigJugada("Roro Real-Envido")
	p.SetSigJugada("Alvaro Falta-Envido")
	p.SetSigJugada("Roro Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+3+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 2+3+10 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestEnvidoRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
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
	p.SetSigJugada("Roro Real-Envido")
	p.SetSigJugada("Alvaro Falta-Envido")
	p.SetSigJugada("Roro No-Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+3+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 2+3+1 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestEnvidoEnvidoRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
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
	p.SetSigJugada("Alvaro Real-Envido")
	p.SetSigJugada("Roro Falta-Envido")
	p.SetSigJugada("Alvaro Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+2+3+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 2+2+3+10 && p.puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

func TestEnvidoEnvidoRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
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
	p.SetSigJugada("Alvaro Real-Envido")
	p.SetSigJugada("Roro Falta-Envido")
	p.SetSigJugada("Alvaro No-Quiero")
	p.Esperar()

	oops = p.Ronda.envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.envido.puntaje == 2+2+3+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.puntajes[Azul] == 0 && p.puntajes[Rojo] == 2+2+3+1)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

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
