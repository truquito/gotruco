package truco

import (
	"testing"
)

var oops = false

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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = p.Ronda.Envido.puntaje != 2
	if oops {
		t.Error(`El puntaje del envido deberia de ser 2`)
		return
	}

	oops = !(p.Puntajes[Azul] == 2 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = p.Ronda.Envido.puntaje != 1
	if oops {
		t.Error(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 1 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 3)
	if oops {
		t.Error(`El puntaje del envido deberia de ser 3`)
		return
	}

	oops = !(p.Puntajes[Azul] == 3 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 1)
	if oops {
		t.Error(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 1 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 10`)
		return
	}

	oops = !(p.Puntajes[Azul] == 10 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 1 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != ENVIDO
	if oops {
		t.Error("El estado del envido deberia de ser `envido`")
		return
	}

	oops = p.Ronda.Envido.puntaje != 2
	if oops {
		t.Error("El `puntaje` del envido deberia de ser 2")
		return
	}

	p.SetSigJugada("Roro Envido")
	p.Esperar()

	oops = p.Ronda.Envido.estado != ENVIDO
	if oops {
		t.Error(`El estado del envido deberia de ser 'envido', incluso luego de que
		ambos Juan y Pedro lo hayan tocando`)
		return
	}

	oops = p.Ronda.Envido.puntaje != 4
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 0 && p.Puntajes[Rojo] == 2+1)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+3)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 2+3 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 0 && p.Puntajes[Rojo] == 2+1)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 2+10 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 0 && p.Puntajes[Rojo] == 2+1)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 3+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 3+10 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 3+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 0 && p.Puntajes[Rojo] == 3+1)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+2+3)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 2+2+3 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 2+2+1 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+2+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 2+2+10 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 2+2+1 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+3+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 2+3+10 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+3+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 2+3+1 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+2+3+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 2+2+3+10 && p.Puntajes[Rojo] == 0)
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

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envido.puntaje == 2+2+3+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[Azul] == 0 && p.Puntajes[Rojo] == 2+2+3+1)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

/* Tests de calculos */
func TestCalcEnvido(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"A", "C", "E"}, []string{"B", "D", "F"})
	p.Puntajes[Azul] = 4
	p.Puntajes[Rojo] = 3
	p.Ronda.setMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // envido: 26
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Oro, Valor: 12},
					Carta{Palo: Copa, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // envido: 20
					Carta{Palo: Copa, Valor: 12},
					Carta{Palo: Copa, Valor: 11},
					Carta{Palo: Basto, Valor: 3},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // envido: 28
					Carta{Palo: Copa, Valor: 2},
					Carta{Palo: Copa, Valor: 6},
					Carta{Palo: Basto, Valor: 1},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // envido: 25
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Oro, Valor: 3},
					Carta{Palo: Basto, Valor: 2},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // envido: 33
					Carta{Palo: Basto, Valor: 6},
					Carta{Palo: Basto, Valor: 7},
					Carta{Palo: Oro, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // envido: 27
					Carta{Palo: Copa, Valor: 3},
					Carta{Palo: Copa, Valor: 4},
					Carta{Palo: Oro, Valor: 4},
				},
			},
		},
	)

	expected := []int{26, 20, 28, 25, 33, 27}
	for i, manojo := range p.Ronda.Manojos {
		got := manojo.calcularEnvido(p.Ronda.Muestra)
		oops = expected[i] != got
		if oops {
			t.Errorf(
				`El resultado del envido del jugador %s es incorrecto.
				\nEXPECTED: %v
				\nGOT: %v`,
				manojo.Jugador.Nombre, expected[i], got)
			return
		}
	}

	p.SetSigJugada("D Envido")
	p.SetSigJugada("C Quiero")
	p.Esperar()

	oops = !(p.Puntajes[Azul] == 4+2)
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}

}

func TestYTEnvidoCalcII(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"A", "C", "E"}, []string{"B", "D", "F"})
	p.Puntajes[Azul] = 4
	p.Puntajes[Rojo] = 3
	p.Ronda.setMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // envido: 21
					Carta{Palo: Basto, Valor: 1},
					Carta{Palo: Basto, Valor: 12},
					Carta{Palo: Copa, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // envido: 23
					Carta{Palo: Oro, Valor: 12},
					Carta{Palo: Oro, Valor: 3},
					Carta{Palo: Basto, Valor: 4},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // envido: 23
					Carta{Palo: Basto, Valor: 10},
					Carta{Palo: Copa, Valor: 6},
					Carta{Palo: Basto, Valor: 3},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // envido: 30
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Oro, Valor: 4},
					Carta{Palo: Copa, Valor: 1},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // envido: 30
					Carta{Palo: Basto, Valor: 6},
					Carta{Palo: Basto, Valor: 4},
					Carta{Palo: Oro, Valor: 1},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // envido: 31
					Carta{Palo: Espada, Valor: 5},
					Carta{Palo: Copa, Valor: 4},
					Carta{Palo: Espada, Valor: 3},
				},
			},
		},
	)

	expected := []int{21, 23, 23, 30, 30, 31}
	for i, manojo := range p.Ronda.Manojos {
		got := manojo.calcularEnvido(p.Ronda.Muestra)
		oops = expected[i] != got
		if oops {
			t.Errorf(
				`El resultado del envido del jugador %s es incorrecto.
				\nEXPECTED: %v
				\nGOT: %v`,
				manojo.Jugador.Nombre, expected[i], got)
			return
		}
	}

	p.SetSigJugada("D Envido")
	p.SetSigJugada("C Quiero")
	p.Esperar()

	oops = !(p.Puntajes[Rojo] == 3+2)
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}

	// error: C deberia decir: son buenas; pero no aparece
}
