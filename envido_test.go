package truco

import (
	"testing"

	"github.com/jpfilevich/truco/pdt"
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
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = p.Ronda.Envite.Puntaje != 2
	if oops {
		t.Error(`El puntaje del envido deberia de ser 2`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 2 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 2`)
		return
	}

}

func TestEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue rechazado por Roro`)
		return
	}

	oops = p.Ronda.Envite.Puntaje != 1
	if oops {
		t.Error(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 1)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo rojo deberia de ser 0`)
		return
	}

}

func TestRealEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 3)
	if oops {
		t.Error(`El puntaje del envido deberia de ser 3`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 3 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestRealEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 1)
	if oops {
		t.Error(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 1 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
		return
	}

}

func TestFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 10`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 10 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 10`)
		return
	}

}

func TestFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 1 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
		return
	}

}

func TestEnvidoEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")

	oops = p.Ronda.Envite.Estado != pdt.ENVIDO
	if oops {
		t.Error("El estado del envido deberia de ser `envido`")
		return
	}

	oops = p.Ronda.Envite.Puntaje != 2
	if oops {
		t.Error("El `puntaje` del envido deberia de ser 2")
		return
	}

	p.Cmd("Roro Envido")

	oops = p.Ronda.Envite.Estado != pdt.ENVIDO
	if oops {
		t.Error(`El estado del envido deberia de ser 'envido', incluso luego de que
		ambos Alvaro y Roro lo hayan tocando`)
		return
	}

	oops = p.Ronda.Envite.Puntaje != 4
	if oops {
		t.Error(`El puntaje del envido deberia ahora de ser '2 + 2 = 4'`)
		return
	}

	p.Cmd("Alvaro Quiero")

}

func TestEnvidoEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+1)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 0`)
		return
	}

}

func TestEnvidoRealEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Real-Envido")
	p.Cmd("Alvaro Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+3)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 2+3 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 5`)
		return
	}

}

func TestEnvidoRealEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Real-Envido")
	p.Cmd("Alvaro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+1)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 2+10 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+1)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 3+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 3+10 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 3+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 3+1)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestEnvidoEnvidoRealEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+2+3)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 2+2+3 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestEnvidoEnvidoRealEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 2+2+1 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestEnvidoEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+2+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 2+2+10 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestEnvidoEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+2+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 2+2+1 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestEnvidoRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Real-Envido")
	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+3+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 2+3+10 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestEnvidoRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Real-Envido")
	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+3+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 2+3+1 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestEnvidoEnvidoRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+2+3+10)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 2+2+3+10 && p.Puntajes[pdt.Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}

}

func TestEnvidoEnvidoRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 13
					&pdt.Carta{Palo: pdt.Oro, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro No-Quiero")

	oops = p.Ronda.Envite.Estado != pdt.DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
		return
	}

	oops = !(p.Ronda.Envite.Puntaje == 2+2+3+1)
	if oops {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
		return
	}

	oops = !(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+2+3+1)
	if oops {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
		return
	}
}

/* Tests de calculos */
func TestCalcEnvido(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"A", "C", "E"}, []string{"B", "D", "F"})
	p.Puntajes[pdt.Azul] = 4
	p.Puntajes[pdt.Rojo] = 3
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 26
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Oro, Valor: 12},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 20
					&pdt.Carta{Palo: pdt.Copa, Valor: 12},
					&pdt.Carta{Palo: pdt.Copa, Valor: 11},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 28
					&pdt.Carta{Palo: pdt.Copa, Valor: 2},
					&pdt.Carta{Palo: pdt.Copa, Valor: 6},
					&pdt.Carta{Palo: pdt.Basto, Valor: 1},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 25
					&pdt.Carta{Palo: pdt.Oro, Valor: 2},
					&pdt.Carta{Palo: pdt.Oro, Valor: 3},
					&pdt.Carta{Palo: pdt.Basto, Valor: 2},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 33
					&pdt.Carta{Palo: pdt.Basto, Valor: 6},
					&pdt.Carta{Palo: pdt.Basto, Valor: 7},
					&pdt.Carta{Palo: pdt.Oro, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 27
					&pdt.Carta{Palo: pdt.Copa, Valor: 3},
					&pdt.Carta{Palo: pdt.Copa, Valor: 4},
					&pdt.Carta{Palo: pdt.Oro, Valor: 4},
				},
			},
		},
	)

	expected := []int{26, 20, 28, 25, 33, 27}
	for i, manojo := range p.Ronda.Manojos {
		got := manojo.CalcularEnvido(p.Ronda.Muestra)
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
	p.Ronda.Turno = 3
	p.Cmd("D Envido")
	p.Cmd("C Quiero")

	oops = !(p.Puntajes[pdt.Azul] == 4+2)
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}

}

func TestCalcEnvido2(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"A", "C", "E"}, []string{"B", "D", "F"})
	p.Puntajes[pdt.Azul] = 4
	p.Puntajes[pdt.Rojo] = 3
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 21
					&pdt.Carta{Palo: pdt.Basto, Valor: 1},
					&pdt.Carta{Palo: pdt.Basto, Valor: 12},
					&pdt.Carta{Palo: pdt.Copa, Valor: 5},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 23
					&pdt.Carta{Palo: pdt.Oro, Valor: 12},
					&pdt.Carta{Palo: pdt.Oro, Valor: 3},
					&pdt.Carta{Palo: pdt.Basto, Valor: 4},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 23
					&pdt.Carta{Palo: pdt.Basto, Valor: 10},
					&pdt.Carta{Palo: pdt.Copa, Valor: 6},
					&pdt.Carta{Palo: pdt.Basto, Valor: 3},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 30
					&pdt.Carta{Palo: pdt.Oro, Valor: 6},
					&pdt.Carta{Palo: pdt.Oro, Valor: 4},
					&pdt.Carta{Palo: pdt.Copa, Valor: 1},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 30
					&pdt.Carta{Palo: pdt.Basto, Valor: 6},
					&pdt.Carta{Palo: pdt.Basto, Valor: 4},
					&pdt.Carta{Palo: pdt.Oro, Valor: 1},
				},
			},
			pdt.Manojo{
				Cartas: [3]*pdt.Carta{ // envido: 31
					&pdt.Carta{Palo: pdt.Espada, Valor: 5},
					&pdt.Carta{Palo: pdt.Copa, Valor: 4},
					&pdt.Carta{Palo: pdt.Espada, Valor: 3},
				},
			},
		},
	)

	expected := []int{21, 23, 23, 30, 30, 31}
	for i, manojo := range p.Ronda.Manojos {
		got := manojo.CalcularEnvido(p.Ronda.Muestra)
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

	p.Ronda.Turno = 3
	p.Cmd("D Envido")
	p.Cmd("C Quiero")

	oops = !(p.Puntajes[pdt.Rojo] == 3+2)
	if oops {
		t.Error("El resultado es incorrecto")
		return
	}

	// error: C deberia decir: son buenas; pero no aparece
}
