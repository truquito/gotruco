package truco

import (
	"fmt"
	"testing"
	"time"

	"github.com/truquito/gotruco/deco"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
	"github.com/truquito/gotruco/util"
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
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	// t.Log(p.Partida) // retorna el json
	// t.Log(p.String()) // retorna el render

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2, func() {
		t.Error(`El puntaje del envido deberia de ser 2`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 2 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 2`)
	})

}

func TestEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})

	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue rechazado por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 1, func() {
		t.Error(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo rojo deberia de ser 0`)
	})

}

func TestRealEnvidoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 3, func() {
		t.Error(`El puntaje del envido deberia de ser 3`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 3 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestRealEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 1, func() {
		t.Error(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 1 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
	})

}

func TestFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 10`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 10`)
	})

}

func TestFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 1 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
	})

}

func TestEnvidoEnvidoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")

	util.Assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error("El estado del envido deberia de ser `envido`")
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2, func() {
		t.Error("El `puntaje` del envido deberia de ser 2")
	})

	p.Cmd("Roro Envido")

	util.Assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`El estado del envido deberia de ser 'envido', incluso luego de que
		ambos Alvaro y Roro lo hayan tocando`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 4, func() {
		t.Error(`El puntaje del envido deberia ahora de ser '2 + 2 = 4'`)
	})

	p.Cmd("Alvaro Quiero")

}

func TestEnvidoEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 0`)
	})

}

func TestEnvidoRealEnvidoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Real-Envido")
	p.Cmd("Alvaro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+3, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 2+3 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 5`)
	})

}

func TestEnvidoRealEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Real-Envido")
	p.Cmd("Alvaro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 2+10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	t.Log(p.String())

	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 3+10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 3+10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 3+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 3+1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoRealEnvidoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+2+3, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 2+2+3 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoRealEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+2+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 2+2+1 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+2+10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 2+2+10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+2+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 2+2+1 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Real-Envido")
	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+3+10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 2+3+10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Real-Envido")
	p.Cmd("Alvaro Falta-Envido")
	p.Cmd("Roro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+3+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 2+3+1 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+2+3+10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 2+2+3+10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 13
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Alvaro No-Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	util.Assert(p.Ronda.Envite.Puntaje == 2+2+3+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+2+3+1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})
}

/* Tests de calculos */
func TestCalcEnvido(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"A", "C", "E"}, []string{"B", "D", "F"}, 4, true, time.Second*10)
	p.Puntajes[pdt.Azul] = 4
	p.Puntajes[pdt.Rojo] = 3
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 26
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Oro, Valor: 12},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 20
					{Palo: pdt.Copa, Valor: 12},
					{Palo: pdt.Copa, Valor: 11},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 28
					{Palo: pdt.Copa, Valor: 2},
					{Palo: pdt.Copa, Valor: 6},
					{Palo: pdt.Basto, Valor: 1},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 25
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Oro, Valor: 3},
					{Palo: pdt.Basto, Valor: 2},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 33
					{Palo: pdt.Basto, Valor: 6},
					{Palo: pdt.Basto, Valor: 7},
					{Palo: pdt.Oro, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 27
					{Palo: pdt.Copa, Valor: 3},
					{Palo: pdt.Copa, Valor: 4},
					{Palo: pdt.Oro, Valor: 4},
				},
			},
		},
	)

	expected := []int{26, 20, 28, 25, 33, 27}
	for i, manojo := range p.Ronda.Manojos {
		got := manojo.CalcularEnvido(p.Ronda.Muestra)

		util.Assert(expected[i] == got, func() {
			t.Errorf(
				`El resultado del envido del jugador %s es incorrecto.
				\nEXPECTED: %v
				\nGOT: %v`,
				manojo.Jugador.ID, expected[i], got)
		})
	}
	p.Ronda.Turno = 3
	p.Cmd("D Envido")
	p.Cmd("C Quiero")

	util.Assert(p.Puntajes[pdt.Azul] == 4+2, func() {
		t.Error("El resultado es incorrecto")
	})

}

func TestCalcEnvido2(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"A", "C", "E"}, []string{"B", "D", "F"}, 4, true, time.Second*10)
	p.Puntajes[pdt.Azul] = 4
	p.Puntajes[pdt.Rojo] = 3
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 21
					{Palo: pdt.Basto, Valor: 1},
					{Palo: pdt.Basto, Valor: 12},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 23
					{Palo: pdt.Oro, Valor: 12},
					{Palo: pdt.Oro, Valor: 3},
					{Palo: pdt.Basto, Valor: 4},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 23
					{Palo: pdt.Basto, Valor: 10},
					{Palo: pdt.Copa, Valor: 6},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 30
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Oro, Valor: 4},
					{Palo: pdt.Copa, Valor: 1},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 30
					{Palo: pdt.Basto, Valor: 6},
					{Palo: pdt.Basto, Valor: 4},
					{Palo: pdt.Oro, Valor: 1},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 31
					{Palo: pdt.Espada, Valor: 5},
					{Palo: pdt.Copa, Valor: 4},
					{Palo: pdt.Espada, Valor: 3},
				},
			},
		},
	)

	expected := []int{21, 23, 23, 30, 30, 31}
	for i, manojo := range p.Ronda.Manojos {
		got := manojo.CalcularEnvido(p.Ronda.Muestra)

		util.Assert(expected[i] == got, func() {
			t.Errorf(
				`El resultado del envido del jugador %s es incorrecto.
				\nEXPECTED: %v
				\nGOT: %v`,
				manojo.Jugador.ID, expected[i], got)
		})
	}

	p.Ronda.Turno = 3
	p.Cmd("D Envido")
	p.Cmd("C Quiero")

	util.Assert(p.Puntajes[pdt.Rojo] == 3+2, func() {
		t.Error("El resultado es incorrecto")
	})

	// error: C deberia decir: son buenas; pero no aparece
}

func TestNoDeberianTenerFlor(t *testing.T) {

	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Copa, Valor: 5})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 6},
					{Palo: pdt.Copa, Valor: 10},
					{Palo: pdt.Copa, Valor: 7},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	tieneFlor, _ := p.Ronda.Manojos[0].TieneFlor(p.Ronda.Muestra)

	util.Assert(tieneFlor == false, func() {
		t.Error(`Alvaro' NO deberia de tener 'flor'`)
	})

	tieneFlor, _ = p.Ronda.Manojos[1].TieneFlor(p.Ronda.Muestra)

	util.Assert(tieneFlor == false, func() {
		t.Error(`Roro' NO deberia de tener 'flor'`)
	})

}

func TestNoDeberianTenerFlor2(t *testing.T) {

	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Copa, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 12},
					{Palo: pdt.Copa, Valor: 10},
					{Palo: pdt.Basto, Valor: 1},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 3},
				},
			},
		},
	)

	tieneFlor, _ := p.Ronda.Manojos[0].TieneFlor(p.Ronda.Muestra)

	util.Assert(tieneFlor == false, func() {
		t.Error(`Alvaro' NO deberia de tener 'flor'`)
	})
}

func TestDeberiaTenerFlor(t *testing.T) {

	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Copa, Valor: 5})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 4},
					{Palo: pdt.Espada, Valor: 10},
					{Palo: pdt.Espada, Valor: 7},
				},
			},
			{
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 1},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Oro, Valor: 3},
				},
			},
		},
	)

	tieneFlor, _ := p.Ronda.Manojos[0].TieneFlor(p.Ronda.Muestra)

	util.Assert(tieneFlor == true, func() {
		t.Error(`Alvaro' deberia tener 'flor'`)
	})

	tieneFlor, _ = p.Ronda.Manojos[1].TieneFlor(p.Ronda.Muestra)

	util.Assert(tieneFlor == true, func() {
		t.Error(`Roro' deberia tener 'flor'`)
	})
}

func TestFlorFlorContraFlorQuiero(t *testing.T) {

	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // Alvaro tiene flor
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 6},
					{Palo: pdt.Basto, Valor: 7},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // Roro
					{Palo: pdt.Oro, Valor: 5},
					{Palo: pdt.Espada, Valor: 5},
					{Palo: pdt.Basto, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // Adolfo tiene flor
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Copa, Valor: 2},
					{Palo: pdt.Copa, Valor: 3},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // Renzo tiene flor
					{Palo: pdt.Oro, Valor: 4},
					{Palo: pdt.Espada, Valor: 4},
					{Palo: pdt.Espada, Valor: 1},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // Andres
					{Palo: pdt.Copa, Valor: 10},
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Basto, Valor: 11},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // Richard tiene flor
					{Palo: pdt.Oro, Valor: 10},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 1},
				},
			},
		},
	)

	t.Log(p)

	p.Cmd("Alvaro Flor")
	p.Cmd("Roro Mazo")
	p.Cmd("Renzo Flor")
	p.Cmd("Adolfo Contra-flor-al-resto")
	p.Cmd("Richard Quiero")

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia ser 'deshabilitado'`)
	})

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado de la flor deberia ser 'deshabilitado'`)
	})

	// duda: se suman solo las flores ganadoras
	// si contraflor AL RESTO -> no acumulativo
	// duda: deberia sumar tambien los puntos de las flores
	// oops = !(p.Puntajes[pdt.Azul] == 4*3+10 && p.Puntajes[pdt.Rojo] == 0)
	// puntos para ganar chico + todas las flores NO ACHICADAS

	util.Assert(p.Puntajes[pdt.Azul] == 10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia ser 2`)
	})
}

// Tests:
// los "me achico" no cuentan para la flor
// Flor		xcg(+3) / xcg(+3)
// Flor + Contra-Flor		xc(+3) / xCadaFlorDelQueHizoElDesafio(+3) + 1
// Flor + [Contra-Flor] + ContraFlorAlResto		~Falta Envido + *TODAS* las flores no achicadas / xcg(+3) + 1

func TestFixFlor(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Richard"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"basto","valor":12},{"palo":"oro","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"basto","valor":10},{"palo":"oro","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"copa","valor":10},{"palo":"basto","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"espada","valor":10},{"palo":"basto","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"copa","valor":3},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"espada","valor":11},{"palo":"espada","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":1},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro envido")
	// pero Richard tiene flor
	// y no le esta sumando esos puntos

	if !(p.Ronda.Envite.Estado == pdt.DESHABILITADO) {
		t.Error(`El estado de la flor deberia ser 'deshabilitado'`)
	} else if !(p.Puntajes[pdt.Rojo] == 3) {
		t.Error(`El puntaje del equipo rojo deberia ser 3 por la flor de richard`)
	}

	p.Cmd("alvaro 6 espada")
	p.Cmd("alvaro 6 espada")
	p.Cmd("roro 5 espada")
	p.Cmd("adolfo 10 oro")
	p.Cmd("renzo 6 basto")
	p.Cmd("andres 6 copa")
	p.Cmd("richard 3 espada")

	p.Cmd("adolfo 10 copa")
	p.Cmd("renzo 10 espada")
	p.Cmd("andres 3 copa")
	p.Cmd("richard 11 espada")
	p.Cmd("alvaro 12 basto")
	p.Cmd("roro 10 basto")

	util.Assert(p.Puntajes[pdt.Rojo] == 3, func() {
		t.Error(`El puntaje del equipo rojo deberia ser 3 por la flor de richard`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 1, func() {
		t.Error(`El puntaje del equipo azul deberia ser 1 por la ronda ganada`)
	})

}

// bug a arreglar:
// hay 2 flores; se cantan ambas -> no pasa nada
func TestFixFlorBucle(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Roro","Richard"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"oro","valor":11},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"basto","valor":10},{"palo":"basto","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":7},{"palo":"oro","valor":5},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"basto","valor":1},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":2},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":11},{"palo":"espada","valor":12},{"palo":"espada","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":10},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro mazo")
	p.Cmd("roro flor")
	p.Cmd("richard flor")

	util.Assert(p.Puntajes[pdt.Rojo] == 6, func() {
		t.Error(`El puntaje del equipo rojo deberia ser 6 por las 2 flores`)
	})

}

// bug a arreglar:
// no se puede cantar contra flor
func TestFixContraFlor(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Adolfo","Renzo"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":1},{"palo":"espada","valor":1},{"palo":"espada","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":4},{"palo":"copa","valor":7},{"palo":"oro","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"copa","valor":11},{"palo":"copa","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"espada","valor":3},{"palo":"espada","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":4},{"palo":"oro","valor":1},{"palo":"oro","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":6},{"palo":"copa","valor":6},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	/*
					┌4─┐1─┐3─┐    ┌12┐3─┐6─┐
					│Es│Or│Or│    │Ba│Es│Es│
					└──┘──┘──┘    └──┘──┘──┘                  ╔════════════════╗
										❀                       │ #Mano: Primera │
						Andres        Renzo                     ╠────────────────╣
				╔══════════════════════════════╗              │ Mano: Alvaro   │
				║                              ║              ╠────────────────╣
				║                              ║              │ Turno: Alvaro  │
				║             ┌4─┐             ║     ❀        ╠────────────────╣
		Richard   ║             │Ba│             ║   Adolfo     │ Puntuacion: 20 │
		┌6─┐6─┐11┐ ║             └──┘             ║ ┌11┐11┐5─┐   ╚════════════════╝
		│Or│Co│Es│ ║                              ║ │Ba│Co│Co│    ╔──────┬──────╗
		└──┘──┘──┘ ║                              ║ └──┘──┘──┘    │ ROJO │ AZUL │
				╚══════════════════════════════╝               ├──────┼──────┤
						Alvaro         Roro                      │  0   │  0   │
						↑                                      ╚──────┴──────╝
					┌1─┐1─┐7─┐    ┌4─┐7─┐5─┐
					│Ba│Es│Es│    │Or│Co│Or│
					└──┘──┘──┘    └──┘──┘──┘
	*/

	p.Cmd("alvaro 1 basto")
	p.Cmd("roro 4 oro")
	p.Cmd("adolfo flor")

	// no deberia dejarlo tirar xq el envite esta en juego
	// tampoco debio de haber pasado su turno
	p.Cmd("adolfo 11 basto")

	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Adolfo", func() {
		t.Error(`No debio de haber pasado su turno`)
	})

	// no deberia dejarlo tirar xq el envite esta en juego
	p.Cmd("renzo 12 basto")

	util.Assert(p.Ronda.Manojos[2].GetCantCartasTiradas() == 0, func() {
		t.Error(`No deberia dejarlo tirar porque nunca llego a ser su turno`)
	})

	// no hay nada que querer
	p.Cmd("renzo quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.FLOR, func() {
		t.Error(`El estado del envite no debio de haber cambiado`)
	})

	p.Cmd("renzo contra-flor")
	p.Cmd("adolfo quiero")

	// renzo tiene 35 vs los 32 de adolfo
	// deberia ganar las 2 flores + x pts

	t.Log(p)

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

	util.Assert(p.Puntajes[pdt.Rojo] > p.Puntajes[pdt.Azul], func() {
		t.Error(`El equipo rojo deberia de tener mas pts que el azul`)
	})
}

func TestTirada1(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 6},
					{Palo: pdt.Basto, Valor: 7},
				},
			},
			{ // Roro no tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 5},
					{Palo: pdt.Espada, Valor: 5},
					{Palo: pdt.Basto, Valor: 5},
				},
			},
			{ // Adolfo tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Copa, Valor: 2},
					{Palo: pdt.Copa, Valor: 3},
				},
			},
			{ // Renzo tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 4},
					{Palo: pdt.Espada, Valor: 4},
					{Palo: pdt.Espada, Valor: 1},
				},
			},
			{ // Andres no tiene  flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 10},
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Basto, Valor: 11},
				},
			},
			{ // Richard tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 10},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 1},
				},
			},
		},
	)

	t.Log(p)
	p.Cmd("Richard flor")
	p.Cmd("Adolfo contra-flor")
	p.Cmd("Richard quiero")
	// p.Cmd("Adolfo no-quiero") // si dice no quero autoamticamente acarrea a alvaro
	// ademas suma 12 puntos y renzo no llego a decir que tenia flor,
	// deberia cantar la de renzo tambien
	p.Cmd("Renzo flor")
	p.Cmd("Alvaro flor")

	p.Cmd("Alvaro 2 oro")
	p.Cmd("Roro 5 oro")
	p.Cmd("Adolfo 1 copa")
	p.Cmd("Renzo 4 oro")
	p.Cmd("Andres 10 copa")
	p.Cmd("Richard 10 oro")

	t.Log(p)

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

	// como la muestra es Palo: pdt.Oro, Valor: 3 -> gana alvaro
	if !(len(p.Ronda.Manos[pdt.Primera].CartasTiradas) == 6) {
		t.Error("La cantidad de cartas tiradas deberia ser 6")
		return

	} else if !(p.Ronda.Manos[pdt.Primera].Ganador == "Alvaro") {
		t.Error("El ganador de la priemra mano deberia ser Alvaro")
		return

	} else if !(p.Ronda.Manos[pdt.Primera].Resultado == pdt.GanoAzul) {
		t.Error("El equipo ganador de la priemra mano deberia ser Azul")
		return
	}

	// como alvaro gano la mano anterior -> empieza tirando el
	p.Cmd("Alvaro 6 basto")
	p.Cmd("Roro 5 espada")
	p.Cmd("Adolfo 2 copa")
	p.Cmd("Renzo 4 espada")
	p.Cmd("Andres 7 oro")

	t.Log(p)

	p.Cmd("Richard 2 oro")

	// como la muestra es Palo: pdt.Oro, Valor: 3 -> gana richard
	if !(len(p.Ronda.Manos[pdt.Segunda].CartasTiradas) == 6) {
		t.Error("La cantidad de cartas tiradas deberia ser 6")
		return

	} else if !(p.Ronda.Manos[pdt.Segunda].Ganador == "Richard") {
		t.Error("El ganador de la priemra mano deberia ser Richard")
		return

	} else if !(p.Ronda.Manos[pdt.Segunda].Resultado == pdt.GanoRojo) {
		t.Error("El equipo ganador de la priemra mano deberia ser pdt.Rojo")
		return
	}

	// vuelvo a checkear que el estado de la pdt.Primera nos se haya editado
	if !(len(p.Ronda.Manos[pdt.Primera].CartasTiradas) == 6) {
		t.Error("La cantidad de cartas tiradas deberia ser 6")
		return

	} else if !(p.Ronda.Manos[pdt.Primera].Ganador == "Alvaro") {
		t.Error("El ganador de la priemra mano deberia ser Alvaro")
		return

	} else if !(p.Ronda.Manos[pdt.Primera].Resultado == pdt.GanoAzul) {
		t.Error("El equipo ganador de la priemra mano deberia ser Azul")
		return
	}

	// como richard gano la mano anterior -> empieza tirando el
	p.Cmd("Richard 1 basto")
	p.Cmd("Alvaro 7 basto")
	p.Cmd("Roro 5 basto")
	p.Cmd("Adolfo 3 copa")
	p.Cmd("Renzo 1 espada")
	p.Cmd("Andres 11 basto")

	// para este momento ya cambio a una nueva ronda
	// como la muestra es Palo: pdt.Oro, Valor: 3 -> gana Renzo con el 1 de espada
	// 1 mano ganada por azul; 2 por rojo -> ronda ganada por rojo
	if !(p.Puntajes[pdt.Rojo] == 1) {
		t.Error("El puntaje del equipo pdt.Rojo deberia ser 1 porque gano la ronda")
		return

	}

}

// no deja irse al mazo a alvaro;
// cuando en realidad deberia poder
// y ademas el turno ahora deberia ser de el siguiente habilitado
// func TestFixIrseAlMazo(t *testing.T) {
// 	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"deshabilitado","puntaje":0,"cantadoPor":null,"sinCantar":["Roro","Andres","Richard"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":7},{"palo":"oro","valor":6},{"palo":"copa","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":10},{"palo":"copa","valor":10},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":11},{"palo":"copa","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":4},{"palo":"oro","valor":5},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"espada","valor":5},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":4},{"palo":"basto","valor":3},{"palo":"basto","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":12},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
// 	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
// 	json.Unmarshal([]byte(partidaJSON), &p)
// 	t.Log(p)

// 	p.Cmd("alvaro mazo")

// 	elManojoDeAlvaro := p.Ronda.Manojos[0]
// 	if !(elManojoDeAlvaro.SeFueAlMazo == true) {
// 		t.Error(`Alvaro se debio de haber ido al mazo`)
// 	}

// }

func TestParseJugada(t *testing.T) {
	p, _ := NuevoJuego(20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 6},
					{Palo: pdt.Basto, Valor: 7},
				},
			},
			{ // Roro no tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 5},
					{Palo: pdt.Espada, Valor: 5},
					{Palo: pdt.Basto, Valor: 5},
				},
			},
		},
	)

	shouldBeOK := []string{
		"alvaro envido",
		"Alvaro real-envido",
		"Alvaro falta-envido",
		"Alvaro flor",
		"Alvaro contra-flor",
		"Alvaro contra-flor-al-resto",
		"Alvaro truco",
		"Alvaro re-truco",
		"Alvaro vale-4",
		"Alvaro quiero",
		"Alvaro no-quiero",
		"Alvaro mazo",
		// tiradas
		"Alvaro 2 oro",
		"Alvaro 2 ORO",
		"Alvaro 2 oRo",
		"Alvaro 6 basto",
		"Alvaro 7 basto",
		"Roro 5 oro",
		"Roro 5 espada",
		"Roro 5 basto",
	}

	shouldNotBeOK := []string{
		"Juancito envido",
		"Juancito envido asd",
		"Juancito envido 33",
		"Juancito envid0",
		// tiradas
		"Alvaro 2 oroo",
		"Alvaro 2 oRo ",
		"Alvaro 6 espada*",
		"Alvaro 7 asd",
		"Alvaro 2  copa",
		"Alvaro 54 Oro ",
		"Alvaro 0 oro",
		"Alvaro 9 oro",
		"Alvaro 8 oro",
		"Alvaro 111 oro",
		// roro trata de usar las de alvaro
		// esto se debe testear en jugadas
		// "roro 2 oRo",
		// "roro 6 basto",
		// "roro 7 basto",
	}

	for _, cmd := range shouldBeOK {
		_, err := pdt.ParseJugada(p.Partida, cmd)

		util.Assert(err == nil, func() {
			t.Error(err.Error())
		})
	}

	for _, cmd := range shouldNotBeOK {
		_, err := pdt.ParseJugada(p.Partida, cmd)

		util.Assert(err != nil, func() {
			t.Error(`Deberia dar error`)
		})
	}
}

func TestPartida1(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 6},
					{Palo: pdt.Basto, Valor: 7},
				},
			},
			{ // Roro no tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 5},
					{Palo: pdt.Espada, Valor: 5},
					{Palo: pdt.Basto, Valor: 5},
				},
			},
			{ // Adolfo tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Copa, Valor: 2},
					{Palo: pdt.Copa, Valor: 3},
				},
			},
			{ // Renzo tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 4},
					{Palo: pdt.Espada, Valor: 4},
					{Palo: pdt.Espada, Valor: 1},
				},
			},
			{ // Andres no tiene  flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 10},
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Basto, Valor: 11},
				},
			},
			{ // Richard tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 10},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 1},
				},
			},
		},
	)
	t.Log(p)
	/*
		               ┌10┐7─┐11┐    ┌4─┐4─┐1─┐
		               │Co│Or│Ba│    │Or│Es│Es│
		               └──┘──┘──┘    └──┘──┘──┘                  ╔════════════════╗
		                                 ❀                       │ #Mano: Primera │
		                 Andres        Renzo                     ╠────────────────╣
		           ╔══════════════════════════════╗              │ Mano: Alvaro   │
		           ║                              ║              ╠────────────────╣
		           ║                              ║              │ Turno: Alvaro  │
		  ❀        ║             ┌3─┐             ║     ❀        ╠────────────────╣
		 Richard   ║             │Or│             ║   Adolfo     │ Puntuacion: 20 │
		┌10┐2─┐1─┐ ║             └──┘             ║ ┌1─┐2─┐3─┐   ╚════════════════╝
		│Or│Or│Ba│ ║                              ║ │Co│Co│Co│    ╔──────┬──────╗
		└──┘──┘──┘ ║                              ║ └──┘──┘──┘    │ ROJO │ AZUL │
		           ╚══════════════════════════════╝               ├──────┼──────┤
		                 Alvaro         Roro                      │  0   │  0   │
		                  ❀ ↑                                     ╚──────┴──────╝
		               ┌2─┐6─┐7─┐    ┌5─┐5─┐5─┐
		               │Or│Ba│Ba│    │Or│Es│Ba│
		               └──┘──┘──┘    └──┘──┘──┘

	*/

	// no deberia dejarlo cantar envido xq tiene flor
	p.Cmd("Alvaro Envido")

	util.Assert(p.Ronda.Envite.Estado != pdt.ENVIDO, func() {
		t.Error(`el envite deberia pasar a estado de flor`)
	})

	// deberia retornar un error debido a que ya canto flor <- deprecated
	p.Cmd("Alvaro Flor")

	// deberia dejarlo irse al mazo
	p.Cmd("Roro Mazo")

	util.Assert(p.Ronda.Manojos[1].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	// deberia retornar un error debido a que ya canto flor
	p.Cmd("Adolfo Flor")

	// deberia aumentar la apuesta
	p.Cmd("Renzo Contra-flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.CONTRAFLOR, func() {
		t.Error(`deberia aumentar la apuesta a CONTRAFLOR`)
	})

	p.Cmd("Alvaro Quiero")
}

func TestPartidaComandosInvalidos(t *testing.T) {

	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"copa","valor":4},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"basto","valor":7},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"espada","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"oro","valor":2},{"palo":"espada","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("Alvaro Envido")
	p.Cmd("Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`no debio de haberlo querido`)
	})

	p.Cmd("Schumacher Flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`no existe schumacher`)
	})

}

func TestPartidaJSON(t *testing.T) {
	p, _ := NuevoJuego(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	pJSON, err := p.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(pJSON))
}

func TestJSONPython(t *testing.T) {
	// json hecho por PY
	data_py := `{"puntuacion": 20, "puntajes": {"azul": 0, "rojo": 0}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"rojo": 1, "azul": 1}, "elMano": 0, "turno": 0, "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": "", "sinCantar": []}, "truco": {"cantadoPor": "", "estado": "noGritadoAun"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "oro", "valor": 1}, {"palo": "copa", "valor": 10}, {"palo": "basto", "valor": 4}], "tiradas": [false, false, false], "ultimaTirada": 0, "jugador": {"id": "Alice", "equipo": "azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "oro", "valor": 11}, {"palo": "basto", "valor": 10}, {"palo": "espada", "valor": 5}], "tiradas": [false, false, false], "ultimaTirada": 0, "jugador": {"id": "Bob", "equipo": "rojo"}}], "mixs": {"Alice": 0, "Bob": 1}, "muestra": {"palo": "basto", "valor": 2}, "manos": [{"resultado": "ganoRojo", "ganador": "", "cartasTiradas": []}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": []}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": []}]}}`
	// json hecho por GO
	data_go := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":1},{"palo":"copa","valor":10},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"basto","valor":10},{"palo":"espada","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Bob","equipo":"rojo"}}],"mixs":{"Alice":0,"Bob":1},"muestra":{"palo":"basto","valor":2},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	// parseado
	// p0, _ := pdt.NuevaPartida(pdt.A20, []string{"Alice"}, []string{"Bob"})
	p1, _ := pdt.Parse(data_go, true)
	p2, _ := pdt.Parse(data_py, true)
	// creado desde 0
	t.Log(pdt.Renderizar(p1))
	t.Log(pdt.Renderizar(p2))
	pJSON, _ := p1.MarshalJSON()
	t.Log(string(pJSON))
}

// - 11 le gana a 10 (de la muestra) no de sparda
// - si es parda pero el turno deberia de ser de el mano (alvaro)
// - adolfo deberia de poder cantar retruco

func TestFixNacho(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Richard"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"copa","valor":7},{"palo":"basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"copa","valor":6},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"espada","valor":1},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"basto","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":12},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":5},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	// p.Partida.FromJSON([]byte(partidaJSON))
	p.Partida, _ = pdt.Parse(partidaJSON, true)

	t.Log(p)

	// 1) Piezas: El “Dos” vale más que el “Cuatro”;
	// éste más que el “Cinco”,
	// éste más que el “perico” (11) y éste más que la “perica” (10).

	p.Cmd("alvaro 6 basto")
	p.Cmd("roro 2 basto")

	cantTiradasRoro := p.Manojo("Roro").GetCantCartasTiradas()

	util.Assert(cantTiradasRoro == 1, func() {
		t.Error(`Roro tiro solo 1 carta`)
	})

	p.Cmd("Adolfo 4 basto")
	p.Cmd("renzo 7 basto")
	p.Cmd("andres 10 espada")
	p.Cmd("richard flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`el envido deberia estar inhabilitado por la flor`)
	})

	p.Cmd("richard 11 espada")

	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Richard", func() {
		t.Error(`Deberia ser el turno de Richard ya que 11 (perico) > 10 (perica)`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

	t.Log(p)

	p.Cmd("richard truco")

	util.Assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`deberia poder cantar truco ya que es su turno`)
	})

	p.Cmd("roro quiero")

	util.Assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`no deberia poder ya que es de su mismo equipo`)
	})

	p.Cmd("adolfo quiero")
	p.Cmd("richard 5 espada")
	p.Cmd("alvaro mazo")
	p.Cmd("roro quiero")

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.ID == "Adolfo", func() {
		t.Error(`no hay nada que querer`)
	})
	p.Cmd("roro retruco") // syntaxis invalida
	p.Cmd("roro re-truco")

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.ID == "Adolfo", func() {
		t.Error(`no debe permitir ya que su equipo no tiene la potestad del truco`)
	})

	p.Cmd("alvaro re-truco")

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.ID == "Adolfo", func() {
		t.Error(`no deberia dejarlo porque se fue al mazo`)
	})

	p.Cmd("Adolfo re-truco")

	util.Assert(p.Ronda.Truco.Estado == pdt.RETRUCO, func() {
		t.Error(`no deberia dejarlo porque se fue al mazo`)
	})

	p.Cmd("renzo quiero")

	util.Assert(p.Ronda.Truco.Estado == pdt.RETRUCOQUERIDO, func() {
		t.Error(`no deberia dejarlo porque se fue al mazo`)
	})

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.ID == "Renzo", func() {
		t.Error(`no deberia dejarlo porque se fue al mazo`)
	})

	p.Cmd("roro 6 copa") // no deberia dejarlo porque ya paso su turno

	util.Assert(cantTiradasRoro == 1, func() {
		t.Error(`Roro tiro solo 1 carta`)
	})

	p.Cmd("adolfo re-truco") // no deberia dejarlo

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.ID == "Renzo", func() {
		t.Error(`no deberia dejarlo porque el re-truco ya fue cantado`)
	})

	p.Cmd("adolfo 1 espada")
	p.Cmd("renzo 3 oro")

	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Andres", func() {
		t.Error(`Deberia ser el turno de Andres`)
	})

	p.Cmd("andres mazo")

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)

}

func TestFixNoFlor(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Andres","Richard"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"basto","valor":4},{"palo":"espada","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":7},{"palo":"basto","valor":11},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":12},{"palo":"basto","valor":1},{"palo":"copa","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"espada","valor":7},{"palo":"oro","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":4},{"palo":"basto","valor":6},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":11},{"palo":"copa","valor":2},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 4 basto")
	// << Alvaro tira la carta 4 de pdt.Basto

	p.Cmd("roro truco")

	// No es posible responder al truco ahora porque
	// "la flor esta primero":
	// el otro que tiene flor, pero se arruga
	p.Cmd("richard no-quiero")

	util.Assert(p.Ronda.Truco.Estado == pdt.NOGRITADOAUN, func() {
		t.Error(`La flor esta primero!`)
	})

	p.Cmd("andres flor")
	p.Cmd("richard flor")
	// << +6 puntos para el equipo pdt.Azul por las flores

	p.Cmd("adolfo 12 oro")
	// No era su turno, no puede tirar la carta

	p.Cmd("roro truco")
	// << Roro grita truco

	p.Cmd("andres quiero")
	// << Andres responde quiero

	p.Cmd("roro 7 copa")
	// << Roro tira la carta 7 de pdt.Copa

	p.Cmd("adolfo 12 oro")
	// << Adolfo tira la carta 12 de pdt.Oro

	p.Cmd("renzo 5 oro")
	// << Renzo tira la carta 5 de pdt.Oro

	p.Cmd("andres flor")
	// No es posible cantar flor

	p.Cmd("andres 6 basto")
	// << Andres tira la carta 6 de pdt.Basto

	p.Cmd("richard flor")
	// no deberia dejarlo porque ya se jugo

	p.Cmd("richard 11 copa")
	// << Richard tira la carta 11 de pdt.Copa

	/* *********************************** */
	// << La Primera mano la gano Adolfo (equipo pdt.Azul)
	// << Es el turno de Adolfo
	/* *********************************** */

	p.Cmd("adolfo re-truco")
	// << Adolfo grita re-truco

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)

	p.Cmd("richard quiero")
	// << Richard responde quiero

	p.Cmd("richard vale-4")
	// << Richard grita vale 4

	util.Assert(p.Ronda.Truco.Estado == pdt.VALE4, func() {
		t.Error(`Richard deberia poder gritar vale4`)
	})

	p.Cmd("adolfo quiero")
	// << Adolfo responde quiero

	util.Assert(p.Ronda.Truco.Estado == pdt.VALE4QUERIDO, func() {
		t.Error(`El estado del truco deberia ser VALE4QUERIDO`)
	})

	/* *********************************** */
	// ACA EMPIEZAN A TIRAR CARTAS PARA LA SEGUNDA MANO
	// muesta: 3 espada
	/* *********************************** */

	p.Cmd("adolfo 1 basto")
	// << Adolfo tira la carta 1 de pdt.Basto

	p.Cmd("renzo 7 espada")
	// << Renzo tira la carta 7 de pdt.Espada

	p.Cmd("andres 4 espada")
	// << Andres tira la carta 4 de pdt.Espada

	p.Cmd("richard 10 espada")
	// << Richard tira la carta 10 de pdt.Espada

	p.Cmd("alvaro 6 espada")
	// << Alvaro tira la carta 6 de pdt.Espada

	p.Cmd("roro re-truco")
	// << Alvaro tira la carta 6 de pdt.Espada

	p.Cmd("roro mazo")
	// << Roro se va al mazo

	// era el ultimo que quedaba por tirar en esta mano
	// -> que evalue la mano

	// << +4 puntos para el equipo pdt.Azul por el vale4Querido no querido por Roro
	// << Empieza una nueva ronda
	// << Empieza una nueva ronda

	/* 6 de las 2 flores */
	// util.Assert(p.GetMaxPuntaje() == 6+4, func() {
	// 	t.Error(`suma mal los puntos cuando roro se fue al mazo`)
	// })

	t.Log(p)

}

func TestFixPanic(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Richard"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"copa","valor":7},{"palo":"basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"copa","valor":6},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"espada","valor":1},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"basto","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":12},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":5},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 6 basto")
	// << Alvaro tira la carta 6 de pdt.Basto

	p.Cmd("roro 2 basto")
	// << Roro tira la carta 2 de pdt.Basto

	p.Cmd("Adolfo 4 basto")
	// << Adolfo tira la carta 4 de pdt.Basto

	p.Cmd("renzo 7 basto")
	// << Renzo tira la carta 7 de pdt.Basto

	p.Cmd("andres 10 espada")
	// << Andres tira la carta 10 de pdt.Espada

	p.Cmd("richard flor")
	// << Richard canta flor
	// << +3 puntos para el equipo pdt.Rojo (por ser la unica flor de esta ronda)

	p.Cmd("richard 11 espada")
	// << Richard tira la carta 11 de pdt.Espada

	/*
		// << La Mano resulta parda
		// << Es el turno de Richard
	*/

	// ERROR: no la deberia ganar andres porque es mano (DUDA)

	p.Cmd("richard truco")
	// << Richard grita truco

	p.Cmd("roro quiero")
	// (Para Roro) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Cmd("adolfo quiero")
	// << Adolfo responde quiero

	p.Cmd("richard 5 espada")
	// << Richard tira la carta 5 de pdt.Espada

	p.Cmd("alvaro mazo")
	// << Alvaro se va al mazo

	p.Cmd("roro quiero")
	// (Para Roro) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Cmd("roro retruco")
	// << No esxiste esa jugada

	p.Cmd("roro re-truco")
	// No es posible cantar re-truco ahora

	p.Cmd("alvaro re-truco") // ya que se fue al mazo
	// No es posible cantar re-truco ahora

	p.Cmd("Adolfo re-truco") // no es su turno ni el de su equipo
	// No es posible cantar re-truco ahora

	p.Cmd("roro 6 copa")
	// << Roro tira la carta 6 de pdt.Copa

	p.Cmd("adolfo re-truco")
	// << Adolfo grita re-truco

	p.Cmd("adolfo 1 espada")
	// << Adolfo tira la carta 1 de pdt.Espada

	p.Cmd("renzo retruco")
	// << No esxiste esa jugada

	p.Cmd("renzo re-truco") // ya que ya lo canto adolfo
	// No es posible cantar re-truco ahora

	p.Cmd("renzo mazo")
	// << Renzo se va al mazo

	t.Log(p)

	p.Cmd("andres mazo")
	// << Andres se va al mazo

	// andres se va al mazo y era el ultimo que quedaba por jugar
	// ademas era la Segunda mano -> ya se decide
	// aunque hay un retruco propuesto <--------------
	// si hay algo propuesto por su equipo no se puede ir <-------

	// << La Segunda mano la gano el equipo pdt.Rojo gracia a Richard
	// << La ronda ha sido ganada por el equipo pdt.Rojo
	// << +0 puntos para el equipo pdt.Rojo por el reTruco no querido
	// << Empieza una nueva ronda

}

func TestFixBocha(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"espada","valor":7},{"palo":"basto","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":12},{"palo":"espada","valor":11},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"oro","valor":6},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":7},{"palo":"basto","valor":10},{"palo":"copa","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":2},{"palo":"copa","valor":3},{"palo":"oro","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":10},{"palo":"oro","valor":2},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro mazo")
	// << Alvaro se va al mazo

	p.Cmd("adolfo mazo")
	// << Adolfo se va al mazo

	p.Cmd("andres mazo")
	// << Andres se va al mazo

	util.Assert(p.Puntajes[pdt.Rojo] == 1 && p.Puntajes[pdt.Azul] == 0, func() {
		t.Error(`todos los de azul se fueron al mazo, deberian de haber ganado los rojos`)
	})

	util.Assert(p.Ronda.GetElMano().Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`todos los de azul se fueron al mazo, deberian ser turno de los rojos`)
	})

	t.Log(p)

}

func TestFixBochaParte2(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"espada","valor":7},{"palo":"basto","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":12},{"palo":"espada","valor":11},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"oro","valor":6},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":7},{"palo":"basto","valor":10},{"palo":"copa","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":2},{"palo":"copa","valor":3},{"palo":"oro","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":10},{"palo":"oro","valor":2},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("roro envido")
	// No es posible cantar 'Envido'

	p.Cmd("andres quiero")
	// (Para Andres) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Cmd("andres quiero")
	// (Para Andres) No hay nada "que querer"; ya que: el estado del envido no es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o bien fue cantado por uno de su equipo

	p.Cmd("alvaro mazo")
	// << Alvaro se va al mazo

	p.Cmd("adolfo 1 copa")
	// Esa carta no se encuentra en este manojo

	p.Cmd("adolfo mazo")
	// << Adolfo se va al mazo

	p.Cmd("andres mazo")
	// << Andres se va al mazo

	util.Assert(p.Puntajes[pdt.Rojo] == 1 && p.Puntajes[pdt.Azul] == 0, func() {
		t.Error(`todos los de azul se fueron al mazo, deberian de haber ganado los rojos`)
	})
	// << La ronda ha sido ganada por el equipo pdt.Rojo
	// << +1 puntos para el equipo pdt.Rojo por el noCantado ganado
	// << Empieza una nueva ronda

	t.Log(p)

}

func TestFixBochaParte3(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"espada","valor":7},{"palo":"basto","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":12},{"palo":"espada","valor":11},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"oro","valor":6},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":7},{"palo":"basto","valor":10},{"palo":"copa","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":2},{"palo":"copa","valor":3},{"palo":"oro","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":10},{"palo":"oro","valor":2},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("richard flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN, func() {
		t.Error(`No es posible cantar flor`)
	})

	// (Para Andres) No hay nada "que querer"; ya que: el estado del envido no
	// es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o
	// bien fue cantado por uno de su equipo
	p.Cmd("andres quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN && p.Ronda.Truco.Estado == pdt.NOGRITADOAUN, func() {
		t.Error(`No hay nada "que querer"`)
	})

	// No es posible cantar contra flor
	p.Cmd("andres contra-flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN, func() {
		t.Error(`No es posible cantar flor`)
	})

	// No es posible cantar contra flor
	p.Cmd("richard contra-flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN, func() {
		t.Error(`No es posible cantar flor`)
	})

	// (Para Richard) No hay nada "que querer"; ya que: el estado del envido no
	// es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o
	// bien fue cantado por uno de su equipo
	p.Cmd("richard quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN && p.Ronda.Truco.Estado == pdt.NOGRITADOAUN, func() {
		t.Error(`No hay nada "que querer"`)
	})

	t.Log(p)

}

func TestFixAutoQuerer(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"espada","valor":7},{"palo":"basto","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":12},{"palo":"espada","valor":11},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"oro","valor":6},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":7},{"palo":"basto","valor":10},{"palo":"copa","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":2},{"palo":"copa","valor":3},{"palo":"oro","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":10},{"palo":"oro","valor":2},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro envido")

	util.Assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`Deberia en estar estado envido`)
	})

	p.Cmd("alvaro quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`No se deberia poder auto-querer`)
	})

	p.Cmd("adolfo quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`No se deberia poder auto-querer a uno del mismo equipo`)
	})

	t.Log(p)

}

func TestFixNilPointer(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Renzo"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"espada","valor":10},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":12},{"palo":"copa","valor":5},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"copa","valor":7},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"basto","valor":1},{"palo":"copa","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"copa","valor":6},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":4},{"palo":"basto","valor":10},{"palo":"copa","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	p.Ronda.Turno = 3 // es el turno de renzo
	t.Log(p)

	p.Cmd("renzo6 basto")
	p.Cmd("renzo 6 basto")
	p.Cmd("andres truco")
	p.Cmd("renzo quiero")
	p.Cmd("andres re-truco")
	p.Cmd("andres 3 oro")
	p.Cmd("richard vale-4")
	p.Cmd("richard re-truco")
	p.Cmd("andres quiero")
	p.Cmd("richard mazo")
	p.Cmd("alvaro vale-4")
	p.Cmd("andres quiero")
	p.Cmd("roro quiero")
	p.Cmd("alvaro mazo")
	p.Cmd("roro mazo")
	p.Cmd("roro 12 oro")
	p.Cmd("adolfo mazo")
	p.Cmd("Renzo flor")

}

func TestFixNoDejaIrseAlMazo(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Renzo"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":11},{"palo":"oro","valor":7},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"copa","valor":2},{"palo":"espada","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":12},{"palo":"oro","valor":4},{"palo":"oro","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":2},{"palo":"espada","valor":10},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":6},{"palo":"copa","valor":7},{"palo":"basto","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":2},{"palo":"basto","valor":2},{"palo":"copa","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 11 copa")
	p.Cmd("roro 6 basto")
	p.Cmd("adolfo 12 espada")
	p.Cmd("renzo 2 espada")
	p.Cmd("renzo flor")
	p.Cmd("andres 6 oro")
	p.Cmd("richard truco")
	p.Cmd("andres quiero")
	p.Cmd("alvaro 7 oro")
	p.Cmd("roro 2 copa")
	p.Cmd("richard 2 oro")
	p.Cmd("renzo mazo")

	util.Assert(p.Ronda.Manojos[3].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	p.Cmd("andres mazo")

	util.Assert(p.Ronda.Manojos[4].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	p.Cmd("andres mazo")

	t.Log(p)

}

func TestFixFlorObligatoria(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alvaro"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"oro","valor":6},{"palo":"oro","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":5},{"palo":"basto","valor":12},{"palo":"espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":7},{"palo":"basto","valor":5},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":1},{"palo":"copa","valor":11},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":10},{"palo":"oro","valor":2},{"palo":"oro","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"basto","valor":3},{"palo":"espada","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":6},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 2 basto") // alvaro deberia primero cantar flor
	// p.Cmd("alvaro flor")
	p.Cmd("roro 5 copa")
	p.Cmd("adolfo 7 espada")
	p.Cmd("renzo 1 espada")
	p.Cmd("andres 10 espada")
	p.Cmd("richard 3 basto")
	p.Cmd("alvaro envido")
	p.Cmd("alvaro 1 oro")
	p.Cmd("roro 2 espada")
	p.Cmd("adolfo truco")
	p.Cmd("roro quiero")
	p.Cmd("renzo quiero")
	p.Cmd("adolfo 5 basto")
	p.Cmd("renzo quiero")
	p.Cmd("renzo 11 copa")
	p.Cmd("andres 2 oro")
	p.Cmd("richard 10 oro")
	p.Cmd("roro 1 oro")

	t.Log(p)

}

func TestFixNoPermiteContraFlor(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Adolfo","Renzo"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":1},{"palo":"espada","valor":1},{"palo":"espada","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":4},{"palo":"copa","valor":7},{"palo":"oro","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"copa","valor":11},{"palo":"copa","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"espada","valor":3},{"palo":"espada","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":4},{"palo":"oro","valor":1},{"palo":"oro","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":6},{"palo":"copa","valor":6},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 1 basto")
	p.Cmd("roro 4 oro")
	p.Cmd("adolfo flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.FLOR, func() {
		t.Error(`deberia permitir cantar flor`)
	})

	p.Cmd("adolfo 11 basto")
	p.Cmd("renzo 12 basto")
	p.Cmd("renzo quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.FLOR, func() {
		t.Error(`no debio de haber cambiado nada`)
	})

	p.Cmd("renzo contra-flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.CONTRAFLOR, func() {
		t.Error(`debe de jugarse la contaflor`)
	})

	t.Log(p)

}

func TestFixDeberiaGanarAzul(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"espada","valor":10},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"oro","valor":5},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"espada","valor":3},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":4},{"palo":"espada","valor":6},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":12},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 10 oro")
	p.Cmd("roro 10 copa")
	p.Cmd("adolfo 2 copa")
	p.Cmd("renzo 4 basto")
	p.Cmd("renzo mazo")
	p.Cmd("alvaro mazo")
	p.Cmd("roro mazo")

	util.Assert(p.Puntajes[pdt.Rojo] == 0 && p.Puntajes[pdt.Azul] > 0, func() {
		t.Error(`La ronda deberia de haber sido ganado por pdt.Azul`)
	})

}

func TestFixPierdeTurno(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"copa","valor":4},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"basto","valor":7},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"espada","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"oro","valor":2},{"palo":"espada","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 5 espada")
	p.Cmd("adolfo mazo")

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)

	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Roro", func() {
		t.Error(`Deberia ser el turno de Roro`)
	})

}

func TestFixTieneFlor(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alvaro"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":5},{"palo":"copa","valor":4},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"basto","valor":7},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"espada","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"oro","valor":2},{"palo":"espada","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 5 copa")
	// no deberia dejarlo jugar porque tiene flor

	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Alvaro", func() {
		t.Error(`Deberia ser el turno de Alvaro`)
	})

}

func Test2FloresSeVaAlMazo(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alvaro","Richard"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"copa","valor":7},{"palo":"copa","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"copa","valor":6},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"espada","valor":1},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"basto","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":12},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":5},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 1 copa") // no deberia dejarlo porque tiene flor
	p.Cmd("alvaro envido") // no deberia dejarlo porque tiene flor
	p.Cmd("alvaro flor")
	p.Cmd("richard mazo") // lo deja que se vaya

	util.Assert(p.Ronda.Manojos[5].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestTodoTienenFlor(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 6},
					{Palo: pdt.Basto, Valor: 7},
				},
			},
			{ // Roro no tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 5},
					{Palo: pdt.Espada, Valor: 5},
					{Palo: pdt.Basto, Valor: 5},
				},
			},
			{ // Adolfo tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Copa, Valor: 2},
					{Palo: pdt.Copa, Valor: 3},
				},
			},
			{ // Renzo no tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 10},
					{Palo: pdt.Oro, Valor: 7},
					{Palo: pdt.Basto, Valor: 11},
				},
			},
			{ // Andres tiene  flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 4},
					{Palo: pdt.Espada, Valor: 4},
					{Palo: pdt.Espada, Valor: 1},
				},
			},
			{ // Richard no tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 11},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 1},
				},
			},
		},
	)

	p.Cmd("Alvaro Envido")

	util.Assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN, func() {
		t.Error(`alvaro tenia flor; no puede tocar envido`)
	})

	p.Cmd("Alvaro Flor")
	p.Cmd("Roro Mazo")
	p.Cmd("Adolfo Flor")

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestFixTopeEnvido(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":10,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":1},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":11},{"palo":"basto","valor":3},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":1},{"palo":"espada","valor":10},{"palo":"basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":4},{"palo":"copa","valor":12},{"palo":"basto","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"copa","valor":7},{"palo":"espada","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"copa","valor":1},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	// azul va 10 pts de 20,
	// asi que el maximo permitido de envite deberia ser 10
	// ~ 5 envidos
	// al 6to saltar error

	p.Cmd("alvaro envido")
	p.Cmd("Roro envido")

	p.Cmd("alvaro envido")
	p.Cmd("Roro envido")

	p.Cmd("alvaro envido")

	pts := p.Ronda.Envite.Puntaje

	p.Cmd("Roro envido") // debe retornar error

	util.Assert(p.Ronda.Envite.Puntaje == pts, func() {
		t.Error(`no se puede cantar mas de 5 envidos seeguidos`)
	})

	p.Cmd("Roro quiero")

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestAutoQuererse(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":10,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":1},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":11},{"palo":"basto","valor":3},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":1},{"palo":"espada","valor":10},{"palo":"basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":4},{"palo":"copa","valor":12},{"palo":"basto","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"copa","valor":7},{"palo":"espada","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"copa","valor":1},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	// no deberia poder auto quererse **ni auto no-quererse**

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Roro Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.FALTAENVIDO, func() {
		t.Error(`no lo deberia dejar porque el envite lo propuso el equipo rojo`)
	})

	p.Cmd("Roro no-quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.FALTAENVIDO, func() {
		t.Error(`no lo deberia dejar porque el envite lo propuso el equipo rojo`)
	})

	p.Cmd("Renzo Quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.FALTAENVIDO, func() {
		t.Error(`no lo deberia dejar porque el envite lo propuso el equipo rojo`)
	})

	p.Cmd("Renzo no-quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.FALTAENVIDO, func() {
		t.Error(`no lo deberia dejar porque el envite lo propuso el equipo rojo`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestJsonSinFlores(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	// partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":2},{"palo":"basto","valor":6},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":1},{"palo":"copa","valor":2},{"palo":"copa","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":4},{"palo":"espada","valor":4},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"oro","valor":2},{"palo":"basto","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}],"JugadoresConFlorQueNoCantaron222":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":2},{"palo":"basto","valor":6},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":1},{"palo":"copa","valor":2},{"palo":"copa","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":4},{"palo":"espada","valor":4},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"oro","valor":2},{"palo":"basto","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":2},{"palo":"basto","valor":6},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"espada","valor":5},{"palo":"basto","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":1},{"palo":"copa","valor":2},{"palo":"copa","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":4},{"palo":"espada","valor":4},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"oro","valor":7},{"palo":"basto","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"oro","valor":2},{"palo":"basto","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alvaro","Adolfo","Renzo","Richard"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":2},{"palo":"basto","valor":6},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"espada","valor":5},{"palo":"basto","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":1},{"palo":"copa","valor":2},{"palo":"copa","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":4},{"palo":"espada","valor":4},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"oro","valor":7},{"palo":"basto","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"oro","valor":2},{"palo":"basto","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))

	// los metodos de las flores son privados
	// deberia testearse en pdt

	t.Log(p)
}

func TestFixEnvidoManoEsElUltimo(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":10,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":5,"turno":3,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":1},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":11},{"palo":"basto","valor":3},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":1},{"palo":"espada","valor":10},{"palo":"basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":4},{"palo":"copa","valor":12},{"palo":"basto","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"copa","valor":7},{"palo":"espada","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"copa","valor":1},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("renzo Envido")
	p.Cmd("andres Envido")
	p.Cmd("richard quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`la sequencia de toques era valida`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestEnvidoManoSeFue(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":10,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":5,"turno":3,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":1},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":11},{"palo":"basto","valor":3},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":1},{"palo":"espada","valor":10},{"palo":"basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"copa","valor":1},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"copa","valor":7},{"palo":"espada","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":4},{"palo":"copa","valor":12},{"palo":"basto","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("renzo Envido")
	p.Cmd("andres Envido")
	p.Cmd("richard mazo")
	p.Cmd("renzo quiero")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`la sequencia de toques era valida`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestFlorBlucle(t *testing.T) {
	p, _ := NuevoJuego(20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 6},
					{Palo: pdt.Basto, Valor: 7},
				},
			},
			{ // Roro tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 5},
					{Palo: pdt.Espada, Valor: 5},
					{Palo: pdt.Espada, Valor: 5},
				},
			},
		},
	)

	p.Cmd("alvaro flor")
	p.Cmd("roro flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`la flor se debio de haber jugado`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestQuieroContraflorDesdeMazo(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Oro, Valor: 3})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{ // Alvaro tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 6},
					{Palo: pdt.Basto, Valor: 7},
				},
			},
			{ // Roro no tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 5},
					{Palo: pdt.Espada, Valor: 5},
					{Palo: pdt.Basto, Valor: 5},
				},
			},
			{ // Adolfo tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Copa, Valor: 2},
					{Palo: pdt.Copa, Valor: 3},
				},
			},
			{ // Renzo no tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 10},
					{Palo: pdt.Copa, Valor: 4},
					{Palo: pdt.Copa, Valor: 11},
				},
			},
			{ // Andres tiene  flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 4},
					{Palo: pdt.Espada, Valor: 4},
					{Palo: pdt.Espada, Valor: 1},
				},
			},
			{ // Richard no tiene flor
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Copa, Valor: 12},
					{Palo: pdt.Oro, Valor: 2},
					{Palo: pdt.Basto, Valor: 1},
				},
			},
		},
	)

	p.Cmd("alvaro flor")
	p.Cmd("andres mazo")

	util.Assert(p.Ronda.Manojos[4].SeFueAlMazo == true, func() {
		t.Error(`andres se debio de haber ido al mazo`)
	})

	p.Cmd("renzo contra-flor")
	p.Cmd("andres quiero")

	util.Assert(p.Ronda.Envite.CantadoPor == "Renzo", func() {
		t.Errorf(`andres no puede responder quiero porque se fue al mazo`)
	})

	t.Log(p.Ronda.Envite.Estado.String())

	util.Assert(p.Ronda.Envite.Estado == pdt.CONTRAFLOR, func() {
		t.Error(`El estado del envite no debio de haber sido modificado`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestFixSeVaAlMazoYTeniaFlor(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alvaro"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":6},{"palo":"espada","valor":10},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"espada","valor":7},{"palo":"basto","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":12},{"palo":"basto","valor":3},{"palo":"espada","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":7},{"palo":"oro","valor":7},{"palo":"basto","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":5},{"palo":"basto","valor":11},{"palo":"copa","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"oro","valor":5},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":5},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	ptsAzul := p.Puntajes[pdt.Azul]

	p.Cmd("alvaro mazo")

	util.Assert(p.Ronda.Manojos[0].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	util.Assert(ptsAzul == p.Puntajes[pdt.Azul], func() {
		t.Error(`no deberia de cambiar el puntaje`)
	})

	p.Cmd("roro truco")

	util.Assert(p.Ronda.Envite.Estado == pdt.EstadoEnvite(pdt.TRUCO), func() {
		t.Error(`Deberia dejarlo cantar truco`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestFixDesconcertante(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alvaro","Renzo"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"espada","valor":2},{"palo":"espada","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"espada","valor":11},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":7},{"palo":"copa","valor":10},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"espada","valor":4},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":12},{"palo":"oro","valor":7},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":4},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro flor")
	p.Cmd("alvaro mazo")

	util.Assert(p.Ronda.Manojos[0].SeFueAlMazo == false, func() {
		t.Error(`No deberia dejarlo irse al mazo porque se esta jugando la flor`)
	})
	p.Cmd("roro truco")

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestMalaAsignacionPts(t *testing.T) {
	p, _ := NuevoJuego(20, []string{"Alvaro"}, []string{"Roro"}, 4, true, time.Second*10)

	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Basto, Valor: 12})
	p.Puntajes[pdt.Rojo] = 3
	p.Puntajes[pdt.Azul] = 2
	p.Ronda.Turno = 1
	p.Ronda.ElMano = 1
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{ // Alvaro
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 10},
					{Palo: pdt.Oro, Valor: 12},
					{Palo: pdt.Espada, Valor: 5},
				},
			},
			{ // Roro
				Cartas: [3]*pdt.Carta{
					{Palo: pdt.Oro, Valor: 1},
					{Palo: pdt.Espada, Valor: 1},
					{Palo: pdt.Espada, Valor: 11},
				},
			},
		},
	)

	t.Log(p)

	p.Cmd("alvaro vale-4")
	p.Cmd("alvaro truco") // vigente
	p.Cmd("roro truco")
	p.Cmd("alvaro re-truco")
	p.Cmd("roro vale-4")
	p.Cmd("alvaro quiero")

	p.Cmd("roro quiero") // <- no hay nada que querer

	p.Cmd("roro 1 espada") // <- gana roro
	p.Cmd("alvaro 12 oro")

	p.Cmd("roro 1 oro") // <- gana roro
	p.Cmd("alvaro 5 espada")

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)

	// pts Vale4 Ganado = 4
	// puntaje de Rojo deberia ser = 3 + 4
	// puntaje de Azul deberia ser = 2

	util.Assert(p.Puntajes[pdt.Rojo] == 3+4 && p.Puntajes[pdt.Azul] == 2, func() {
		t.Error(`Asigno mal los puntos`)
	})
}

func TestFixRondaNueva(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Renzo"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"copa","valor":1},{"palo":"espada","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"espada","valor":4},{"palo":"basto","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"basto","valor":6},{"palo":"oro","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"oro","valor":10},{"palo":"oro","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":4},{"palo":"basto","valor":7},{"palo":"espada","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":12},{"palo":"espada","valor":10},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":6},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("renzo flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`deberia deshabilitar el envite`)
	})

	p.Cmd("alvaro 1 copa")

	util.Assert(p.Ronda.Manojos[0].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("roro truco")

	util.Assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`Deberia dejarlo cantar truco`)
	})

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`El equipo rojo deberia tener la potestad del truco`)
	})

	p.Cmd("adolfo re-truco")

	util.Assert(p.Ronda.Truco.Estado == pdt.RETRUCO, func() {
		t.Error(`Deberia dejarlo cantar re-truco`)
	})

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == pdt.Azul, func() {
		t.Error(`El equipo azul deberia tener la potestad del truco`)
	})

	p.Cmd("renzo vale-4")

	util.Assert(p.Ronda.Truco.Estado == pdt.VALE4, func() {
		t.Error(`Deberia dejarlo cantar vale-3`)
	})

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`El equipo rojo deberia tener la potestad del truco`)
	})

	p.Cmd("adolfo quiero")

	util.Assert(p.Ronda.Truco.Estado == pdt.VALE4QUERIDO, func() {
		t.Error(`Deberia dejarlo responder quiero al vale-4`)
	})

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == pdt.Azul, func() {
		t.Error(`El equipo azul deberia tener la potestad del truco`)
	})

	p.Cmd("roro 5 oro")

	util.Assert(p.Ronda.Manojos[1].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("adolfo 6 basto")

	util.Assert(p.Ronda.Manojos[2].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("renzo 19 oro")

	util.Assert(p.Ronda.Manojos[3].GetCantCartasTiradas() == 0, func() {
		t.Error(`no debeia dejarlo porque no existe la carta "19 de oro"`)
	})

	p.Cmd("renzo 10 oro")

	util.Assert(p.Ronda.Manojos[3].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("andres 4 copa")

	util.Assert(p.Ronda.Manojos[4].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("richard 10 espada")

	util.Assert(p.Ronda.Manojos[5].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	util.Assert(p.Ronda.Manojo(p.Ronda.Manos[0].Ganador).Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`La primera mano la debio de haber ganado el equipo de richard: el rojo`)
	})

	// segunda mano
	p.Cmd("renzo 10 oro")

	util.Assert(p.Ronda.Manojos[3].GetCantCartasTiradas() == 1, func() {
		t.Error(`ya tiro esa carta; no deberia poder volve a tirarla`)
	})

	p.Cmd("andres 7 basto")

	util.Assert(p.Ronda.Manojos[4].GetCantCartasTiradas() == 1, func() {
		t.Error(`no es su turno no deberia poder tirar carta`)
	})

	p.Cmd("richard 10 espada")

	util.Assert(p.Ronda.Manojos[5].GetCantCartasTiradas() == 1, func() {
		t.Error(`ya tiro esa carta; no deberia poder volve a tirarla`)
	})

	p.Cmd("richard 12 copa")

	util.Assert(p.Ronda.Manojos[5].GetCantCartasTiradas() == 2, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("alvaro 3 espada")

	util.Assert(p.Ronda.Manojos[0].GetCantCartasTiradas() == 2, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("roro 4 espada")

	util.Assert(p.Ronda.Manojos[1].GetCantCartasTiradas() == 2, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("adolfo 2 copa")

	util.Assert(p.Ronda.Manojos[2].GetCantCartasTiradas() == 2, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("renzo 4 oro")

	util.Assert(p.Ronda.Manojos[3].GetCantCartasTiradas() == 2, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("andres 5 espada")

	util.Assert(p.Ronda.Manojos[4].GetCantCartasTiradas() == 0, func() {
		t.Error(`deberia tener 0 cartas tiradas porque empieza una nueva ronda`)
	})

	/* 3:flor + 4:vale4 */
	util.Assert(p.Puntajes[pdt.Rojo] == 3+4, func() {
		t.Error(`el puntaje para el equipo rojo deberia ser 7: 3 de la flor + 4 del vale4`)
	})

	util.Assert(p.Puntajes[pdt.Azul] == 0, func() {
		t.Error(`el puntaje para el equipo azul deberia ser 0 porque no ganaron nada`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)
}

func TestFixIrseAlMazo2(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"espada","valor":7},{"palo":"basto","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"espada","valor":3},{"palo":"copa","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"oro","valor":5},{"palo":"espada","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":12},{"palo":"oro","valor":2},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":7},{"palo":"basto","valor":7},{"palo":"oro","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"copa","valor":3},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":6},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("renzo flor")

	util.Assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN, func() {
		t.Error(`Renzo no tiene flor`)
	})

	// mano 1
	p.Cmd("alvaro envido")

	util.Assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`Debio dejarlo cantar truco`)
	})

	p.Cmd("alvaro 6 basto")

	util.Assert(p.Ronda.Manojos[0].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	p.Cmd("roro 11 oro")

	util.Assert(p.Ronda.Manojos[1].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	p.Cmd("adolfo 2 basto")

	util.Assert(p.Ronda.Manojos[2].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	p.Cmd("renzo 12 oro")

	util.Assert(p.Ronda.Manojos[3].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	p.Cmd("andres 7 copa")

	util.Assert(p.Ronda.Manojos[4].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	p.Cmd("richard 12 basto")

	util.Assert(p.Ronda.Manojos[5].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	// mano 2
	p.Cmd("roro 3 espada")
	p.Cmd("adolfo 5 oro")
	p.Cmd("renzo 2 oro")
	p.Cmd("andres 7 basto")
	p.Cmd("richard truco")
	p.Cmd("renzo quiero")

	p.Cmd("andres re-truco")
	p.Cmd("richard vale-4")
	p.Cmd("alvaro quiero")

	t.Log(p)
	p.Cmd("roro mazo")

	util.Assert(p.Ronda.Manojos[1].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

	t.Log(p)

}

func TestFixDecirQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Renzo"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":1},{"palo":"copa","valor":6},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":7},{"palo":"copa","valor":12},{"palo":"oro","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":11},{"palo":"espada","valor":2},{"palo":"espada","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":10},{"palo":"espada","valor":12},{"palo":"basto","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"basto","valor":3},{"palo":"basto","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":2},{"palo":"basto","valor":6},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":2},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("renzo flor")
	p.Cmd("alvaro truco")

	util.Assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`Deberia poder gritar truco`)
	})

	p.Cmd("renzo quiero")

	util.Assert(p.Ronda.Truco.Estado == pdt.TRUCOQUERIDO, func() {
		t.Error(`Deberia poder responder quiero al truco`)
	})

	p.Cmd("alvaro re-truco")

	util.Assert(p.Ronda.Truco.Estado == pdt.TRUCOQUERIDO, func() {
		t.Error(`Como no tiene la potestad, no deberia poder aumentar la apuesta`)
	})

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`El equpo Rojo deberia de seguir manteniendo la potestad`)
	})

	p.Cmd("renzo re-truco")
	p.Cmd("alvaro vale-4")

	util.Assert(p.Ronda.Truco.Estado == pdt.VALE4, func() {
		t.Error(`Deberia poder aumentar a vale-4`)
	})

	p.Cmd("alvaro quiero")

	util.Assert(p.Ronda.Truco.Estado == pdt.VALE4, func() {
		t.Error(`No puede auto-querse`)
	})

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == pdt.Azul, func() {
		t.Error(`El equpo azul deberia tener la potestad`)
	})

	p.Cmd("renzo re-truco")

	util.Assert(p.Ronda.Truco.Estado == pdt.VALE4, func() {
		t.Error(`No deberia cambiar el estado del truco`)
	})

	util.Assert(p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == pdt.Azul, func() {
		t.Error(`El equpo azul deberia de seguir manteniendo la potestad`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)

}

func TestFixPanicNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alvaro","Renzo"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":1},{"palo":"basto","valor":10},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":4},{"palo":"espada","valor":6},{"palo":"basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":2},{"palo":"copa","valor":5},{"palo":"basto","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":4},{"palo":"oro","valor":12},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":7},{"palo":"oro","valor":11},{"palo":"oro","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"copa","valor":2},{"palo":"basto","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":1},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro flor")
	p.Cmd("alvaro 1 basto")
	p.Cmd("renzo flor")

	ptsPostFlor := p.Puntajes[pdt.Rojo]

	p.Cmd("alvaro 1 basto")
	p.Cmd("roro 4 copa")
	p.Cmd("adolfo 2 espada")
	p.Cmd("renzo 4 oro") // la Primera mano la gana renzo
	p.Cmd("andres 7 espada")
	p.Cmd("richard 10 copa")

	util.Assert(p.Ronda.Manojo(p.Ronda.Manos[0].Ganador).Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`La primera mano la debio de haber ganado el equipo de renzo: el rojo`)
	})

	p.Cmd("renzo 12 oro")
	p.Cmd("andres 11 oro") // la seguna mano la gana andres
	p.Cmd("richard 2 copa")
	p.Cmd("alvaro 10 basto")
	p.Cmd("roro 6 espada")
	p.Cmd("adolfo 5 copa")

	util.Assert(p.Ronda.Manojo(p.Ronda.Manos[1].Ganador).Jugador.Equipo == pdt.Azul, func() {
		t.Error(`La segunda mano la debio de haber ganado el equipo de andres: el Azul`)
	})

	p.Cmd("andres 3 oro")
	p.Cmd("richard truco")
	p.Cmd("richard mazo")
	p.Cmd("alvaro quiero")
	p.Cmd("alvaro re-truco")
	p.Cmd("renzo quiero")
	p.Cmd("roro vale-4")
	p.Cmd("andres no-quiero")

	util.Assert(p.Puntajes[pdt.Rojo] == ptsPostFlor+3, func() {
		t.Error(`Deberian gana 3 puntines por el vale-4 no querido`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)

}

func TestFixCartaYaJugada(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Richard"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":2},{"palo":"oro","valor":4},{"palo":"basto","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":4},{"palo":"espada","valor":5},{"palo":"basto","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"espada","valor":3},{"palo":"oro","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"oro","valor":2},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":7},{"palo":"basto","valor":3},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"oro","valor":6},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":11},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 2 espada")
	p.Cmd("roro 4 copa")
	p.Cmd("adolfo 10 copa")
	p.Cmd("renzo 6 copa")
	p.Cmd("andres 7 copa")
	p.Cmd("richard flor")
	p.Cmd("richard 5 oro")
	p.Cmd("richard 5 oro")

	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Richard", func() {
		t.Error(`Deberia ser el turno de Richard`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)

}

func TestFixTrucoNoQuiero(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"copa","valor":7},{"palo":"espada","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"basto","valor":2},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":4},{"palo":"copa","valor":5},{"palo":"copa","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":7},{"palo":"espada","valor":3},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":10},{"palo":"espada","valor":2},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"espada","valor":7},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":10},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro truco")
	p.Cmd("roro no-quiero")

	util.Assert(p.Puntajes[pdt.Azul] > 0, func() {
		t.Error(`La ronda deberia de haber sido ganado por Azul`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}
	for _, pkt := range p.Consumir() {
		t.Log(deco.Stringify(&pkt, p.Partida))
	}
	t.Log(p)

}

func TestPerspectiva(t *testing.T) {
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Richard"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"copa","valor":7},{"palo":"basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"copa","valor":6},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"espada","valor":1},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"basto","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":12},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":5},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))

	per, _ := p.Perspectiva("Alvaro")
	t.Log(per.MarshalJSON())
}

func TestPardaSigTurno1(t *testing.T) {
	// si va parda, el siguiente turno deberia ser del mano
	// o del mas cercano a este
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 26
					{Palo: pdt.Copa, Valor: 7},
					{Palo: pdt.Oro, Valor: 12},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 27
					{Palo: pdt.Copa, Valor: 3},
					{Palo: pdt.Copa, Valor: 4},
					{Palo: pdt.Oro, Valor: 4},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 28
					{Palo: pdt.Copa, Valor: 2},
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Espada, Valor: 6}, // <-- parda primera mano
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 25
					{Palo: pdt.Basto, Valor: 2},
					{Palo: pdt.Oro, Valor: 3},
					{Palo: pdt.Oro, Valor: 6}, // <-- parda primera mano
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 33
					{Palo: pdt.Basto, Valor: 3},
					{Palo: pdt.Basto, Valor: 7},
					{Palo: pdt.Copa, Valor: 6}, // <-- parda primera mano
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 20
					{Palo: pdt.Copa, Valor: 12},
					{Palo: pdt.Copa, Valor: 11},
					{Palo: pdt.Basto, Valor: 6}, // <-- parda primera mano
				},
			},
		},
	)

	p.Cmd("Alvaro 5 copa")
	p.Cmd("Roro 4 copa")
	// los siguientes 4: todos tiran 6 -> resulta mano parda
	p.Cmd("Adolfo 6 espada")
	p.Cmd("Renzo 6 oro")
	p.Cmd("Andres 6 copa")
	p.Cmd("Richard 6 basto")

	t.Log(p)

	util.Assert(p.Ronda.Manos[0].Resultado == pdt.Empardada, func() {
		t.Error(`La mano debio ser parda`)
	})

	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Adolfo", func() {
		t.Error(`Deberia ser turno de Adolfo, debido a que es el mas cercano del mano y qu empardo`)
	})
}

func TestPardaSigTurno2(t *testing.T) {
	// igual que el anterior pero ahora adolfo se va al mazo
	// si va parda, el siguiente turno deberia ser del mano
	// o del mas cercano a este
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true, time.Second*10)
	p.Ronda.SetMuestra(pdt.Carta{Palo: pdt.Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]pdt.Manojo{
			{
				Cartas: [3]*pdt.Carta{ // envido: 26
					{Palo: pdt.Copa, Valor: 7},
					{Palo: pdt.Oro, Valor: 12},
					{Palo: pdt.Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 27
					{Palo: pdt.Copa, Valor: 3},
					{Palo: pdt.Copa, Valor: 4},
					{Palo: pdt.Oro, Valor: 4},
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 28
					{Palo: pdt.Copa, Valor: 2},
					{Palo: pdt.Copa, Valor: 1},
					{Palo: pdt.Espada, Valor: 6}, // <-- parda primera mano
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 25
					{Palo: pdt.Basto, Valor: 2},
					{Palo: pdt.Oro, Valor: 3},
					{Palo: pdt.Oro, Valor: 6}, // <-- parda primera mano
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 33
					{Palo: pdt.Basto, Valor: 3},
					{Palo: pdt.Basto, Valor: 7},
					{Palo: pdt.Copa, Valor: 6}, // <-- parda primera mano
				},
			},
			{
				Cartas: [3]*pdt.Carta{ // envido: 20
					{Palo: pdt.Copa, Valor: 12},
					{Palo: pdt.Copa, Valor: 11},
					{Palo: pdt.Basto, Valor: 6}, // <-- parda primera mano
				},
			},
		},
	)

	p.Cmd("Alvaro 5 copa")
	p.Cmd("Roro 4 copa")
	// los siguientes 4: todos tiran 6 -> resulta mano parda
	p.Cmd("Adolfo 6 espada")
	p.Cmd("Adolfo mazo")
	p.Cmd("Renzo 6 oro")
	p.Cmd("Andres 6 copa")
	p.Cmd("Richard 6 basto")

	t.Log(p)

	util.Assert(p.Ronda.Manos[0].Resultado == pdt.Empardada, func() {
		t.Error(`La mano debio ser parda`)
	})

	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Renzo", func() {
		t.Error(`Deberia ser turno de Renzo, debido a que es el mas cercano del mano y qu empardo`)
	})
}

func TestFixTrucoDeshabilitaEnvido(t *testing.T) {
	// cantar truco (sin siquiera ser querido) deshabilita el envido
	// cuando en la vida real es posible tocar "el envido esta primero"
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"copa","valor":4},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"basto","valor":7},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"espada","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"oro","valor":2},{"palo":"espada","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("Alvaro truco")

	util.Assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`Deberia dejarlo gritar truco`)
	})

	// el envido esta primero!!
	p.Cmd("Roro envido")

	util.Assert(!enco.Contains(p.Consumir(), enco.TError), func() {
		t.Error("No deberia resultar en un error tocar envido ahora")
	})

	util.Assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`Deberia dejarlo tocar envido`)
	})

}

func TestAbandono(t *testing.T) {
	// simulacro de un jugador abandonando
	p, _ := NuevoJuego(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":2,"rojo":3},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"copa","valor":4},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"basto","valor":7},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"espada","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"oro","valor":2},{"palo":"espada","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("Alvaro truco")

	util.Assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`Deberia dejarlo gritar truco`)
	})

	// el envido esta primero!!
	p.Abandono("Adolfo")

	util.Assert(p.Terminado(), func() {
		t.Error(`Deberia haber acabado la partida`)
	})

	util.Assert(p.Puntajes[pdt.Rojo] == int(p.Puntuacion), func() {
		t.Error(`El equipo rojo debio de haber completado los pts`)
	})

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

}

// func TestFixOrdenCantoFlor(t *testing.T) {
// 	// simulacro de un jugador abandonando
// 	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
// 	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alvaro","Renzo","Andres"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":12},{"palo":"espada","valor":7},{"palo":"espada","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":11},{"palo":"espada","valor":11},{"palo":"espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"oro","valor":10},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"oro","valor":3},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":4},{"palo":"copa","valor":3},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":4},{"palo":"copa","valor":12},{"palo":"basto","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":7},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
// 	p.Partida.FromJSON([]byte(partidaJSON))
// 	t.Log(p)

// 	p.Cmd("alvaro flor")
// 	p.Cmd("renzo flor")
// 	p.Cmd("andres flor")

// 	util.Assert(enco.Contains(enco.Collect(p.Out), enco.TDiceSonBuenas), func() {
// 		t.Error("debio de haber dicho son bueas")
// 	})

// }

func TestFixTester2(t *testing.T) {
	// simulacro de un jugador abandonando
	p, _ := NuevoJuego(pdt.A30, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":30,"puntajes":{"azul":27,"rojo":23},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":2,"turno":3,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alvaro"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":12},{"palo":"basto","valor":5},{"palo":"basto","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"espada","valor":1},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":6},{"palo":"copa","valor":7},{"palo":"espada","valor":12}],"tiradas":[false,true,false],"ultimaTirada":1,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"copa","valor":12},{"palo":"basto","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":5},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas":[{"jugador":"Adolfo","carta":{"palo":"copa","valor":7}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("Adolfo no-quiero")
	p.Cmd("Adolfo 7 copa")
	p.Cmd("Roro 3 oro")
	p.Cmd("Renzo contra-flor-al-resto")
	p.Cmd("Roro 7 oro")
	p.Cmd("Alvaro flor")

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

}

func TestFixFlorNoCantada(t *testing.T) {
	// Roro tiene flor y aun asi es capaz de tira carta sin cantarla
	// no deberia ser posible; primero debe cantar la flor
	p, _ := NuevoJuego(pdt.A30, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true, time.Second*10)
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Roro"]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":4},{"palo":"espada","valor":2},{"palo":"basto","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":1},{"palo":"basto","valor":10},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":1},{"palo":"basto","valor":5},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":6},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":12},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.Partida.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("Alvaro 2 basto")

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

	roro := p.Manojo("Roro")

	// ANTES (deprecated, pero anda):
	// aa := pdt.GetA(p.Partida, roro)
	// AHORA:
	aa := pdt.Chi(p.Partida, roro)

	t.Log(aa)

	// p.Cmd("Roro 4 basto")

	// enco.Consume(p.Out, func(pkt *enco.Packet2) {
	// 	t.Log(deco.Stringify(pkt, p.Partida))
	// })

	t.Log(p)
}

func TestUpdateJSONs(t *testing.T) {
	// data := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"espada","valor":11},{"palo":"espada","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}],"JugadoresConFlorQueNoCantaron222":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"espada","valor":11},{"palo":"espada","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"basto","valor":12},{"palo":"oro","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"basto","valor":10},{"palo":"oro","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"copa","valor":10},{"palo":"basto","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"espada","valor":10},{"palo":"basto","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"copa","valor":3},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"espada","valor":11},{"palo":"espada","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":1},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	// // p, _, _ := NuevaPartida(pdt.A30, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	// // p.Partida.FromJSON([]byte(data))
	// p, _ := pdt.Parse(data, true)
	// // p.Ronda.Envite.SinCantar = make([]string, 0)
	// json, _ := p.MarshalJSON()
	// fmt.Println(string(json))

	// if true {
	// 	return
	// }

	datas := []string{
		`{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"azul": 7, "rojo": 0}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"azul": 2, "rojo": 1}, "elMano": 1, "turno": 1, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": null, "estado": "noGritadoAun"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "copa", "valor": 11}, {"palo": "copa", "valor": 7}, {"palo": "copa", "valor": 3}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "copa", "valor": 10}, {"palo": "copa", "valor": 2}, {"palo": "basto", "valor": 7}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "espada", "valor": 11}, {"palo": "basto", "valor": 4}, {"palo": "oro", "valor": 10}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "azul"}}, {"seFueAlMazo": true, "cartas": [{"palo": "basto", "valor": 6}, {"palo": "oro", "valor": 4}, {"palo": "oro", "valor": 3}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "rojo"}}], "muestra": {"palo": "espada", "valor": 1}, "manos": [{"resultado": "ganoRojo", "ganador": "", "cartasTiradas": []}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": []}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": []}]}}`,
	}

	// OJO que ahora usa cachearFlores(false) ->
	// los sinCantar quedan en null si es que no tenia ese atributo en el json

	for _, data := range datas {
		p, _ := pdt.Parse(data, true)
		json, _ := p.MarshalJSON()
		fmt.Println(string(json))
	}
}

func TestChiBugfixRust(t *testing.T) {
	data := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":"primera","cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"valor":5,"palo":"oro"},{"valor":10,"palo":"copa"},{"valor":2,"palo":"copa"}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"valor":11,"palo":"copa"},{"valor":1,"palo":"basto"},{"valor":7,"palo":"oro"}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"bob","equipo":"rojo"}}],"mixs": {"alice": 0, "bob": 1},"muestra":{"valor":12,"palo":"espada"},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, _ := pdt.Parse(data, true)
	t.Log(pdt.Renderizar(p))

	defer func() {
		t.Log(pdt.Renderizar(p))
	}()

	// var (
	// 	err  error
	// 	pkts []*enco.Packet
	// )

	// if !(len(p.Ronda.MIXS) > 0) {
	// 	t.Error(`el len del mixs deberia ser positivo`)
	// 	t.Fail()
	// }

	// err_msg := func(t *testing.T, p *pdt.Partida, i int, err error, pkts []*enco.Packet) {
	// 	t.Logf("no deberia ocurrir %d", i)
	// 	if err != nil {
	// 		t.Logf("\nerr: %s\n", err)
	// 	}
	// }

	// check := func(err error, pkts []*enco.Packet) bool {
	// 	return err != nil || enco.Contains(pkts, enco.Error)
	// }

	// cmds := []string{
	// 	"alice 10 copa",
	// 	// "bob falta-envido",
	// 	// "alice no-quiero",
	// 	"bob 11 copa",
	// 	"alice 5 oro", // no debe ocurrir porque es turno de bob
	// 	"alice truco",
	// 	"bob re-truco",
	// 	"alice quiero",
	// 	"alice vale-4",
	// 	"alice 2 copa",
	// }

	// for i, cmd := range cmds {
	// 	pkts, err = p.Cmd(cmd)
	// 	t.Logf("%d: %s", i, cmd)
	// 	if check(err, pkts) {
	// 		err_msg(t, p, i, err, pkts)
	// 		break
	// 	}
	// }

	// if !(p.Ronda.GetElTurno().Jugador.ID == "bob") {
	// 	t.Error("deberia ser el turno de bob")
	// }

	// muestra: 12 espada
	p.Cmd("alice 2 copa")
	p.Cmd("bob 11 copa") // 2>11 --> 1era mano: gana alice
	p.Cmd("alice 5 oro")
	p.Cmd("bob truco")
	p.Cmd("alice re-truco")
	p.Cmd("bob quiero")
	p.Cmd("bob vale-4")
	p.Cmd("alice quiero") // vale4querido
	p.Cmd("bob 7 oro")    // 5o<7o --> 2da mano: gana bob
	p.Cmd("bob 1 basto")
	p.Cmd("alice 10 copa") // 10c < 1b --> 3era mano gana bob --> +4

	if !(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 4) {
		t.Error("los puntajes no son los esperados")
	}

	if !(p.Ronda.ManoEnJuego == pdt.Primera) {
		t.Error("se deberia de estar jugando la primera mano")
	}

	// random gameplay
	// for {
	// 	chis := pdt.MetaChis(p, false)
	// 	rmix, raix := pdt.Random_action_chis(chis)
	// 	j := chis[rmix][raix]
	// 	t.Log(j)
	// 	pkts := j.Hacer(p)
	// 	if pdt.IsDone(pkts) {
	// 		break
	// 	}
	// }
}

func TestChiBugfixRust2(t *testing.T) {
	data := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":"primera","cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"valor":2,"palo":"oro"},{"valor":4,"palo":"espada"},{"valor":12,"palo":"copa"}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"valor":7,"palo":"copa"},{"valor":10,"palo":"copa"},{"valor":6,"palo":"espada"}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"bob","equipo":"rojo"}}],"mixs":{"alice":0,"bob":1},"muestra":{"valor":10,"palo":"basto"},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, _ := pdt.Parse(data, true)
	t.Log(pdt.Renderizar(p))

	defer func() {
		t.Log(pdt.Renderizar(p))
	}()

	p.Cmd("alice envido")
	p.Cmd("bob quiero")
}
