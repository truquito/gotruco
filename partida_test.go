package truco

import (
	"testing"

	"github.com/filevich/truco/deco"
	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/pdt"
)

func assert(should bool, callback func()) {
	if !should {
		callback()
	}
}

// contains dado un buffer se fija si contiene un mensaje
// con ese codigo (y string de ser no-nulo)
func contains(pkts []*enco.Packet, cod enco.CodMsg) bool {
	for _, pkt := range pkts {
		if pkt.Message.Cod == int(cod) {
			return true
		}
	}
	return false
}

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
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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
	p.Cmd("Roro Quiero")

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2, func() {
		t.Error(`El puntaje del envido deberia de ser 2`)
	})

	assert(p.Puntajes[pdt.Azul] == 2 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 2`)
	})

}

func TestEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue rechazado por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 1, func() {
		t.Error(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo rojo deberia de ser 0`)
	})

}

func TestRealEnvidoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 3, func() {
		t.Error(`El puntaje del envido deberia de ser 3`)
	})

	assert(p.Puntajes[pdt.Azul] == 3 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestRealEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 1, func() {
		t.Error(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 1 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
	})

}

func TestFaltaEnvidoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue aceptado por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 10`)
	})

	assert(p.Puntajes[pdt.Azul] == 10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 10`)
	})

}

func TestFaltaEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 1 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 1`)
	})

}

func TestEnvidoEnvidoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error("El estado del envido deberia de ser `envido`")
	})

	assert(p.Ronda.Envite.Puntaje == 2, func() {
		t.Error("El `puntaje` del envido deberia de ser 2")
	})

	p.Cmd("Roro Envido")

	assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`El estado del envido deberia de ser 'envido', incluso luego de que
		ambos Alvaro y Roro lo hayan tocando`)
	})

	assert(p.Ronda.Envite.Puntaje == 4, func() {
		t.Error(`El puntaje del envido deberia ahora de ser '2 + 2 = 4'`)
	})

	p.Cmd("Alvaro Quiero")

}

func TestEnvidoEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 0`)
	})

}

func TestEnvidoRealEnvidoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+3, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 2+3 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 5`)
	})

}

func TestEnvidoRealEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 2+10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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
	p.Cmd("Alvaro Quiero")

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 3+10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 3+10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 3+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 3+1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoRealEnvidoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+2+3, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 2+2+3 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoRealEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+2+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 2+2+1 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+2+10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 2+2+10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+2+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 2+2+1 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+3+10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 2+3+10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+3+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 2+3+1 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoRealEnvidoFaltaEnvidoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+2+3+10, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 2+2+3+10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})

}

func TestEnvidoEnvidoRealEnvidoFaltaEnvidoNoQuiero(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia de ser 'deshabilitado',
		ya que fue no-querido por Roro`)
	})

	assert(p.Ronda.Envite.Puntaje == 2+2+3+1, func() {
		t.Errorf(`El puntaje del envido deberia de ser 1`)
	})

	assert(p.Puntajes[pdt.Azul] == 0 && p.Puntajes[pdt.Rojo] == 2+2+3+1, func() {
		t.Error(`El puntaje del equipo azul deberia de ser 3`)
	})
}

/* Tests de calculos */
func TestCalcEnvido(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"A", "C", "E"}, []string{"B", "D", "F"})
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

		assert(expected[i] == got, func() {
			t.Errorf(
				`El resultado del envido del jugador %s es incorrecto.
				\nEXPECTED: %v
				\nGOT: %v`,
				manojo.Jugador.Nombre, expected[i], got)
			return
		})
	}
	p.Ronda.Turno = 3
	p.Cmd("D Envido")
	p.Cmd("C Quiero")

	assert(p.Puntajes[pdt.Azul] == 4+2, func() {
		t.Error("El resultado es incorrecto")
	})

}

func TestCalcEnvido2(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"A", "C", "E"}, []string{"B", "D", "F"})
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

		assert(expected[i] == got, func() {
			t.Errorf(
				`El resultado del envido del jugador %s es incorrecto.
				\nEXPECTED: %v
				\nGOT: %v`,
				manojo.Jugador.Nombre, expected[i], got)
			return
		})
	}

	p.Ronda.Turno = 3
	p.Cmd("D Envido")
	p.Cmd("C Quiero")

	assert(p.Puntajes[pdt.Rojo] == 3+2, func() {
		t.Error("El resultado es incorrecto")
	})

	// error: C deberia decir: son buenas; pero no aparece
}

func TestNoDeberianTenerFlor(t *testing.T) {

	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(tieneFlor == false, func() {
		t.Error(`Alvaro' NO deberia de tener 'flor'`)
	})

	tieneFlor, _ = p.Ronda.Manojos[1].TieneFlor(p.Ronda.Muestra)

	assert(tieneFlor == false, func() {
		t.Error(`Roro' NO deberia de tener 'flor'`)
	})

}

func TestNoDeberianTenerFlor2(t *testing.T) {

	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(tieneFlor == false, func() {
		t.Error(`Alvaro' NO deberia de tener 'flor'`)
	})
}

func TestDeberiaTenerFlor(t *testing.T) {

	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(tieneFlor == true, func() {
		t.Error(`Alvaro' deberia tener 'flor'`)
	})

	tieneFlor, _ = p.Ronda.Manojos[1].TieneFlor(p.Ronda.Muestra)

	assert(tieneFlor == true, func() {
		t.Error(`Roro' deberia tener 'flor'`)
	})
}

func TestFlorFlorContraFlorQuiero(t *testing.T) {

	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
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

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado del envido deberia ser 'deshabilitado'`)
	})

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`El estado de la flor deberia ser 'deshabilitado'`)
	})

	// duda: se suman solo las flores ganadoras
	// si contraflor AL RESTO -> no acumulativo
	// duda: deberia sumar tambien los puntos de las flores
	// oops = !(p.Puntajes[pdt.Azul] == 4*3+10 && p.Puntajes[pdt.Rojo] == 0)
	// puntos para ganar chico + todas las flores NO ACHICADAS

	assert(p.Puntajes[pdt.Azul] == 10 && p.Puntajes[pdt.Rojo] == 0, func() {
		t.Error(`El puntaje del equipo azul deberia ser 2`)
	})
}

// Tests:
// los "me achico" no cuentan para la flor
// Flor		xcg(+3) / xcg(+3)
// Flor + Contra-Flor		xc(+3) / xCadaFlorDelQueHizoElDesafio(+3) + 1
// Flor + [Contra-Flor] + ContraFlorAlResto		~Falta Envido + *TODAS* las flores no achicadas / xcg(+3) + 1

func TestFixFlor(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":12},{"palo":"Oro","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Basto","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Espada","valor":10},{"palo":"Basto","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Copa","valor":3},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":1},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
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

	assert(p.Puntajes[pdt.Rojo] == 3, func() {
		t.Error(`El puntaje del equipo rojo deberia ser 3 por la flor de richard`)
	})

	assert(p.Puntajes[pdt.Azul] == 1, func() {
		t.Error(`El puntaje del equipo azul deberia ser 1 por la ronda ganada`)
	})

}

// bug a arreglar:
// hay 2 flores; se cantan ambas -> no pasa nada
func TestFixFlorBucle(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":12},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":12},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":11},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Oro","valor":5},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Basto","valor":1},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":2},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":12},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":10},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro mazo")
	p.Cmd("roro flor")
	p.Cmd("richard flor")

	assert(p.Puntajes[pdt.Rojo] == 6, func() {
		t.Error(`El puntaje del equipo rojo deberia ser 6 por las 2 flores`)
	})

}

// bug a arreglar:
// no se puede cantar contra flor
func TestFixContraFlor(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":1},{"palo":"Espada","valor":1},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Copa","valor":7},{"palo":"Oro","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Oro","valor":1},{"palo":"Oro","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Copa","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
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

	assert(p.Ronda.GetElTurno().GetCantCartasTiradas() == 0, func() {
		t.Error(`El puntaje del equipo rojo deberia ser 3 por la flor de richard`)
	})

	assert(p.Ronda.GetElTurno().Jugador.Nombre == "Adolfo", func() {
		t.Error(`No debio de haber pasado su turno`)
	})

	// no deberia dejarlo tirar xq el envite esta en juego
	p.Cmd("renzo 12 basto")

	assert(p.Ronda.Manojos[2].GetCantCartasTiradas() == 0, func() {
		t.Error(`No deberia dejarlo tirar porque nunca llego a ser su turno`)
	})

	// no hay nada que querer
	p.Cmd("renzo quiero")

	assert(p.Ronda.Envite.Estado == pdt.FLOR, func() {
		t.Error(`El estado del envite no debio de haber cambiado`)
	})

	p.Cmd("renzo contra-flor")
	p.Cmd("adolfo quiero")

	// renzo tiene 35 vs los 32 de adolfo
	// deberia ganar las 2 flores + x pts

	t.Log(p)

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})

	assert(p.Puntajes[pdt.Rojo] > p.Puntajes[pdt.Azul], func() {
		t.Error(`El equipo rojo deberia de tener mas pts que el azul`)
	})
}

func TestTirada1(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
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

	p.Cmd("Alvaro 2 Oro")
	p.Cmd("Roro 5 Oro")
	p.Cmd("Adolfo 1 Copa")
	p.Cmd("Renzo 4 Oro")
	p.Cmd("Andres 10 Copa")
	p.Cmd("Richard 10 Oro")

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})

	// como la muestra es Palo: pdt.Oro, Valor: 3 -> gana alvaro
	if !(len(p.Ronda.Manos[pdt.Primera].CartasTiradas) == 6) {
		t.Error("La cantidad de cartas tiradas deberia ser 6")
		return

	} else if !(p.Ronda.Manos[pdt.Primera].Ganador.Jugador.Nombre == "Alvaro") {
		t.Error("El ganador de la priemra mano deberia ser Alvaro")
		return

	} else if !(p.Ronda.Manos[pdt.Primera].Resultado == pdt.GanoAzul) {
		t.Error("El equipo ganador de la priemra mano deberia ser Azul")
		return
	}

	// como alvaro gano la mano anterior -> empieza tirando el
	p.Cmd("Alvaro 6 Basto")
	p.Cmd("Roro 5 Espada")
	p.Cmd("Adolfo 2 Copa")
	p.Cmd("Renzo 4 Espada")
	p.Cmd("Andres 7 Oro")
	p.Cmd("Richard 2 Oro")

	// como la muestra es Palo: pdt.Oro, Valor: 3 -> gana richard
	if !(len(p.Ronda.Manos[pdt.Segunda].CartasTiradas) == 6) {
		t.Error("La cantidad de cartas tiradas deberia ser 6")
		return

	} else if !(p.Ronda.Manos[pdt.Segunda].Ganador.Jugador.Nombre == "Richard") {
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

	} else if !(p.Ronda.Manos[pdt.Primera].Ganador.Jugador.Nombre == "Alvaro") {
		t.Error("El ganador de la priemra mano deberia ser Alvaro")
		return

	} else if !(p.Ronda.Manos[pdt.Primera].Resultado == pdt.GanoAzul) {
		t.Error("El equipo ganador de la priemra mano deberia ser Azul")
		return
	}

	// como richard gano la mano anterior -> empieza tirando el
	p.Cmd("Richard 1 Basto")
	p.Cmd("Alvaro 7 Basto")
	p.Cmd("Roro 5 Basto")
	p.Cmd("Adolfo 3 Copa")
	p.Cmd("Renzo 1 Espada")
	p.Cmd("Andres 11 Basto")

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
// 	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":7},{"palo":"Oro","valor":6},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Copa","valor":10},{"palo":"Copa","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Oro","valor":5},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":5},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":4},{"palo":"Basto","valor":3},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":12},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
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
	p, _, _ := NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})
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
		"Roro 5 Oro",
		"Roro 5 Espada",
		"Roro 5 Basto",
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
		_, err := p.parseJugada(cmd)

		assert(err == nil, func() {
			t.Error(err.Error())
		})
	}

	for _, cmd := range shouldNotBeOK {
		_, err := p.parseJugada(cmd)

		assert(err != nil, func() {
			t.Error(`Deberia dar error`)
		})
	}
}

func TestPartida1(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
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

	assert(p.Ronda.Envite.Estado != pdt.ENVIDO, func() {
		t.Error(`el envite deberia pasar a estado de flor`)
	})

	// deberia retornar un error debido a que ya canto flor
	p.Cmd("Alvaro Flor")

	// deberia dejarlo irse al mazo
	p.Cmd("Roro Mazo")

	assert(p.Ronda.Manojos[1].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	// deberia retornar un error debido a que ya canto flor
	p.Cmd("Adolfo Flor")

	// deberia aumentar la apuesta
	p.Cmd("Renzo Contra-flor")

	assert(p.Ronda.Envite.Estado == pdt.CONTRAFLOR, func() {
		t.Error(`deberia aumentar la apuesta a CONTRAFLOR`)
	})

	p.Cmd("Alvaro Quiero")
}

func TestPartidaComandosInvalidos(t *testing.T) {

	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":7},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Espada","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Oro","valor":2},{"palo":"Espada","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("Alvaro Envido")
	p.Cmd("Quiero")

	assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`no debio de haberlo querido`)
	})

	p.Cmd("Schumacher Flor")

	assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`no existe schumacher`)
	})

}

func TestPartidaJSON(t *testing.T) {
	p, _, _ := NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	pJSON, err := p.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(pJSON))
}

// - 11 le gana a 10 (de la muestra) no de sparda
// - si es parda pero el turno deberia de ser de el mano (alvaro)
// - adolfo deberia de poder cantar retruco

func TestFixNacho(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))

	t.Log(p)

	// 1) Piezas: El “Dos” vale más que el “Cuatro”;
	// éste más que el “Cinco”,
	// éste más que el “perico” (11) y éste más que la “perica” (10).

	p.Cmd("alvaro 6 basto")
	p.Cmd("roro 2 basto")

	roro, _ := p.GetManojoByStr("Roro")
	cantTiradasRoro := roro.GetCantCartasTiradas()

	assert(cantTiradasRoro == 1, func() {
		t.Error(`Roro tiro solo 1 carta`)
	})

	p.Cmd("Adolfo 4 basto")
	p.Cmd("renzo 7 basto")
	p.Cmd("andres 10 espada")
	p.Cmd("richard flor")

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`el envido deberia estar inhabilitado por la flor`)
	})

	p.Cmd("richard 11 espada")

	assert(p.Ronda.GetElTurno().Jugador.Nombre == "Richard", func() {
		t.Error(`Deberia ser el turno de Richard ya que 11 (perico) > 10 (perica)`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})

	t.Log(p)

	p.Cmd("richard truco")

	assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`deberia poder cantar truco ya que es su turno`)
	})

	p.Cmd("roro quiero")

	assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`no deberia poder ya que es de su mismo equipo`)
	})

	p.Cmd("adolfo quiero")
	p.Cmd("richard 5 espada")
	p.Cmd("alvaro mazo")
	p.Cmd("roro quiero")

	assert(p.Ronda.Truco.CantadoPor.Jugador.Nombre == "Adolfo", func() {
		t.Error(`no hay nada que querer`)
	})
	p.Cmd("roro retruco") // syntaxis invalida
	p.Cmd("roro re-truco")

	assert(p.Ronda.Truco.CantadoPor.Jugador.Nombre == "Adolfo", func() {
		t.Error(`no debe permitir ya que su equipo no tiene la potestad del truco`)
	})

	p.Cmd("alvaro re-truco")

	assert(p.Ronda.Truco.CantadoPor.Jugador.Nombre == "Adolfo", func() {
		t.Error(`no deberia dejarlo porque se fue al mazo`)
	})

	p.Cmd("Adolfo re-truco")

	assert(p.Ronda.Truco.Estado == pdt.RETRUCO, func() {
		t.Error(`no deberia dejarlo porque se fue al mazo`)
	})

	p.Cmd("renzo quiero")

	assert(p.Ronda.Truco.Estado == pdt.RETRUCOQUERIDO, func() {
		t.Error(`no deberia dejarlo porque se fue al mazo`)
	})

	assert(p.Ronda.Truco.CantadoPor.Jugador.Nombre == "Renzo", func() {
		t.Error(`no deberia dejarlo porque se fue al mazo`)
	})

	p.Cmd("roro 6 copa") // no deberia dejarlo porque ya paso su turno

	assert(cantTiradasRoro == 1, func() {
		t.Error(`Roro tiro solo 1 carta`)
	})

	p.Cmd("adolfo re-truco") // no deberia dejarlo

	assert(p.Ronda.Truco.CantadoPor.Jugador.Nombre == "Renzo", func() {
		t.Error(`no deberia dejarlo porque el re-truco ya fue cantado`)
	})

	p.Cmd("adolfo 1 espada")
	p.Cmd("renzo 3 oro")

	assert(p.Ronda.GetElTurno().Jugador.Nombre == "Andres", func() {
		t.Error(`Deberia ser el turno de Andres`)
	})

	p.Cmd("andres mazo")

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)

}

func TestFixNoFlor(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Basto","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":2},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Basto","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":2},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":4},{"palo":"Espada","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":7},{"palo":"Basto","valor":11},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":12},{"palo":"Basto","valor":1},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Espada","valor":7},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Basto","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":2},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 4 basto")
	// << Alvaro tira la carta 4 de pdt.Basto

	p.Cmd("roro truco")
	// No es posible cantar truco ahora
	// "la flor esta primero":
	// << Andres canta flor

	assert(p.Ronda.Envite.Estado == pdt.FLOR, func() {
		t.Error(`El envido esta primero!`)
	})

	// el otro que tiene flor, pero se arruga
	p.Cmd("richard no-quiero")
	// FIX: NO ESTA OUTPUTEANDO EL NO QUIERO
	// << +6 puntos para el equipo pdt.Azul por las flores

	p.Cmd("roro truco")
	// << Roro grita truco

	p.Cmd("adolfo 12 oro")
	// No era su turno, no puede tirar la carta

	p.Cmd("roro 7 copa")
	// << Roro tira la carta 7 de pdt.Copa

	p.Cmd("andres quiero")
	// << Andres responde quiero

	p.Cmd("adolfo 12 oro")
	// << Adolfo tira la carta 12 de pdt.Oro

	p.Cmd("renzo 5 oro")
	// << Renzo tira la carta 5 de pdt.Oro

	p.Cmd("andres flor")
	// No es posible cantar flor

	p.Cmd("andres 6 basto")
	// << Andres tira la carta 6 de pdt.Basto

	p.Cmd("richard flor")
	// No es posible cantar flor

	p.Cmd("richard 11 copa")
	// << Richard tira la carta 11 de pdt.Copa

	/* *********************************** */
	// << La Primera mano la gano Adolfo (equipo pdt.Azul)
	// << Es el turno de Adolfo
	/* *********************************** */

	p.Cmd("adolfo re-truco")
	// << Adolfo grita re-truco

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)

	p.Cmd("richard quiero")
	// << Richard responde quiero

	p.Cmd("richard vale-4")
	// << Richard grita vale 4

	assert(p.Ronda.Truco.Estado == pdt.VALE4, func() {
		t.Error(`Richard deberia poder gritar vale4`)
	})

	p.Cmd("adolfo quiero")
	// << Adolfo responde quiero

	assert(p.Ronda.Truco.Estado == pdt.VALE4QUERIDO, func() {
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
	assert(p.GetMaxPuntaje() == 6+4, func() {
		t.Error(`suma mal los puntos cuando roro se fue al mazo`)
	})

	t.Log(p)

}

func TestFixPanic(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
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
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro mazo")
	// << Alvaro se va al mazo

	p.Cmd("adolfo mazo")
	// << Adolfo se va al mazo

	p.Cmd("andres mazo")
	// << Andres se va al mazo

	assert(p.Puntajes[pdt.Rojo] == 1 && p.Puntajes[pdt.Azul] == 0, func() {
		t.Error(`todos los de azul se fueron al mazo, deberian de haber ganado los rojos`)
	})

	assert(p.Ronda.GetElMano().Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`todos los de azul se fueron al mazo, deberian ser turno de los rojos`)
	})

	t.Log(p)

}

func TestFixBochaParte2(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
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

	assert(p.Puntajes[pdt.Rojo] == 1 && p.Puntajes[pdt.Azul] == 0, func() {
		t.Error(`todos los de azul se fueron al mazo, deberian de haber ganado los rojos`)
	})
	// << La ronda ha sido ganada por el equipo pdt.Rojo
	// << +1 puntos para el equipo pdt.Rojo por el noCantado ganado
	// << Empieza una nueva ronda

	t.Log(p)

}

func TestFixBochaParte3(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("richard flor")

	assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN, func() {
		t.Error(`No es posible cantar flor`)
	})

	// (Para Andres) No hay nada "que querer"; ya que: el estado del envido no
	// es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o
	// bien fue cantado por uno de su equipo
	p.Cmd("andres quiero")

	assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN && p.Ronda.Truco.Estado == pdt.NOCANTADO, func() {
		t.Error(`No hay nada "que querer"`)
	})

	// No es posible cantar contra flor
	p.Cmd("andres contra-flor")

	assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN, func() {
		t.Error(`No es posible cantar flor`)
	})

	// No es posible cantar contra flor
	p.Cmd("richard contra-flor")

	assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN, func() {
		t.Error(`No es posible cantar flor`)
	})

	// (Para Richard) No hay nada "que querer"; ya que: el estado del envido no
	// es "envido" (o mayor) y el estado del truco no es "truco" (o mayor) o
	// bien fue cantado por uno de su equipo
	p.Cmd("richard quiero")

	assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN && p.Ronda.Truco.Estado == pdt.NOCANTADO, func() {
		t.Error(`No hay nada "que querer"`)
	})

	t.Log(p)

}

func TestFixAutoQuerer(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":7},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":11},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Basto","valor":10},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro envido")

	assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`Deberia en estar estado envido`)
	})

	p.Cmd("alvaro quiero")

	assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`No se deberia poder auto-querer`)
	})

	p.Cmd("adolfo quiero")

	assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`No se deberia poder auto-querer a uno del mismo equipo`)
	})

	t.Log(p)

}

func TestFixNilPointer(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo": "Oro", "valor":11}, {"palo": "Espada", "valor":10}, {"palo": "Basto", "valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo": "Oro", "valor":12}, {"palo": "Copa", "valor":5}, {"palo": "Copa", "valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo": "Espada", "valor":3}, {"palo": "Copa", "valor":7}, {"palo": "Basto", "valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo": "Basto", "valor":6}, {"palo": "Basto", "valor":1}, {"palo": "Copa", "valor":4 }],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo": "Oro", "valor":3}, {"palo": "Copa", "valor":6}, {"palo": "Copa", "valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo": "Espada", "valor":4}, {"palo": "Basto", "valor":10}, {"palo": "Copa", "valor":10 }],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
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
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Espada","valor":10},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Espada","valor":10},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Oro","valor":7},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Copa","valor":2},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Oro","valor":4},{"palo":"Oro","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Espada","valor":10},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Copa","valor":7},{"palo":"Basto","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":2},{"palo":"Copa","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
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

	assert(p.Ronda.Manojos[3].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	p.Cmd("andres mazo")

	assert(p.Ronda.Manojos[4].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	p.Cmd("andres mazo")

	t.Log(p)

}

func TestFixFlorObligatoria(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":7},{"palo":"Basto","valor":5},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":1},{"palo":"Copa","valor":11},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":10},{"palo":"Oro","valor":2},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Basto","valor":3},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
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
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envido":{"puntaje":0,"cantadoPor":null,"estado":"noCantadoAun"},"flor":1,"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":1},{"palo":"Espada","valor":1},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Copa","valor":7},{"palo":"Oro","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Copa","valor":11},{"palo":"Copa","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Espada","valor":3},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Oro","valor":1},{"palo":"Oro","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Copa","valor":6},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 1 basto")
	p.Cmd("roro 4 oro")
	p.Cmd("adolfo flor")

	assert(p.Ronda.Envite.Estado == pdt.FLOR, func() {
		t.Error(`deberia permitir cantar flor`)
	})

	p.Cmd("adolfo 11 basto")
	p.Cmd("renzo 12 basto")
	p.Cmd("renzo quiero")

	assert(p.Ronda.Envite.Estado == pdt.FLOR, func() {
		t.Error(`no debio de haber cambiado nada`)
	})

	p.Cmd("renzo contra-flor")

	assert(p.Ronda.Envite.Estado == pdt.CONTRAFLOR, func() {
		t.Error(`debe de jugarse la contaflor`)
	})

	t.Log(p)

}

func TestFixDeberiaGanarAzul(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":null,"JugadoresConFlorQueNoCantaron":[]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Espada","valor":10},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Oro","valor":5},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Espada","valor":3},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":4},{"palo":"Espada","valor":6},{"palo":"Copa","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":12},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 10 oro")
	p.Cmd("roro 10 copa")
	p.Cmd("adolfo 2 copa")
	p.Cmd("renzo 4 basto")
	p.Cmd("renzo mazo")
	p.Cmd("alvaro mazo")
	p.Cmd("roro mazo")

	assert(p.Puntajes[pdt.Rojo] == 0 && p.Puntajes[pdt.Azul] > 0, func() {
		t.Error(`La ronda deberia de haber sido ganado por pdt.Azul`)
	})

}

func TestFixPierdeTurno(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":7},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Espada","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Oro","valor":2},{"palo":"Espada","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 5 espada")
	p.Cmd("adolfo mazo")

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)

	assert(p.Ronda.GetElTurno().Jugador.Nombre == "Roro", func() {
		t.Error(`Deberia ser el turno de Roro`)
	})

}

func TestFixTieneFlor(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":7},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Espada","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Oro","valor":2},{"palo":"Espada","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 5 copa")
	// no deberia dejarlo jugar porque tiene flor

	assert(p.Ronda.GetElTurno().Jugador.Nombre == "Alvaro", func() {
		t.Error(`Deberia ser el turno de Alvaro`)
	})

}

func Test2FloresSeVaAlMazo(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Copa","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 1 copa") // no deberia dejarlo porque tiene flor
	p.Cmd("alvaro envido") // no deberia dejarlo porque tiene flor
	p.Cmd("alvaro flor")
	p.Cmd("richard mazo") // lo deja que se vaya

	assert(p.Ronda.Manojos[5].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestTodoTienenFlor(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
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

	assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN, func() {
		t.Error(`alvaro tenia flor; no puede tocar envido`)
	})

	p.Cmd("Alvaro Flor")
	p.Cmd("Roro Mazo")
	p.Cmd("Adolfo Flor")

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestFixTopeEnvido(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":10,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":null,"JugadoresConFlorQueNoCantaron":[]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":1},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":11},{"palo":"Basto","valor":3},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":1},{"palo":"Espada","valor":10},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Copa","valor":12},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":7},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":1},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
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

	assert(p.Ronda.Envite.Puntaje == pts, func() {
		t.Error(`no se puede cantar mas de 5 envidos seeguidos`)
	})

	p.Cmd("Roro quiero")

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestAutoQuererse(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":10,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":null,"JugadoresConFlorQueNoCantaron":[]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":1},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":11},{"palo":"Basto","valor":3},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":1},{"palo":"Espada","valor":10},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Copa","valor":12},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":7},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":1},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	// no deberia poder auto quererse **ni auto no-quererse**

	p.Cmd("Alvaro Envido")
	p.Cmd("Roro Envido")
	p.Cmd("Alvaro Real-Envido")
	p.Cmd("Roro Falta-Envido")
	p.Cmd("Roro Quiero")

	assert(p.Ronda.Envite.Estado == pdt.FALTAENVIDO, func() {
		t.Error(`no lo deberia dejar porque el envite lo propuso el equipo rojo`)
	})

	p.Cmd("Roro no-quiero")

	assert(p.Ronda.Envite.Estado == pdt.FALTAENVIDO, func() {
		t.Error(`no lo deberia dejar porque el envite lo propuso el equipo rojo`)
	})

	p.Cmd("Renzo Quiero")

	assert(p.Ronda.Envite.Estado == pdt.FALTAENVIDO, func() {
		t.Error(`no lo deberia dejar porque el envite lo propuso el equipo rojo`)
	})

	p.Cmd("Renzo no-quiero")

	assert(p.Ronda.Envite.Estado == pdt.FALTAENVIDO, func() {
		t.Error(`no lo deberia dejar porque el envite lo propuso el equipo rojo`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestJsonSinFlores(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	// partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":6},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":1},{"palo":"Copa","valor":2},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Espada","valor":4},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Oro","valor":2},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":6},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":1},{"palo":"Copa","valor":2},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Espada","valor":4},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Oro","valor":2},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":6},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Espada","valor":5},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":1},{"palo":"Copa","valor":2},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Espada","valor":4},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Oro","valor":7},{"palo":"Basto","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Oro","valor":2},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":6},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Espada","valor":5},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":1},{"palo":"Copa","valor":2},{"palo":"Copa","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Espada","valor":4},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Oro","valor":7},{"palo":"Basto","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Oro","valor":2},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))

	// los metodos de las flores son privados
	// deberia testearse en pdt

	t.Log(p)
}

func TestFixEnvidoManoEsElUltimo(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":10,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":5,"turno":3,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":null,"JugadoresConFlorQueNoCantaron":[]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":1},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":11},{"palo":"Basto","valor":3},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":1},{"palo":"Espada","valor":10},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Copa","valor":12},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":7},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":1},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("renzo Envido")
	p.Cmd("andres Envido")
	p.Cmd("richard quiero")

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`la sequencia de toques era valida`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestEnvidoManoSeFue(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":10,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":5,"turno":3,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":null,"JugadoresConFlorQueNoCantaron":[]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":1},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":11},{"palo":"Basto","valor":3},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":1},{"palo":"Espada","valor":10},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":1},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":7},{"palo":"Espada","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":4},{"palo":"Copa","valor":12},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("renzo Envido")
	p.Cmd("andres Envido")
	p.Cmd("richard mazo")
	p.Cmd("renzo quiero")

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`la sequencia de toques era valida`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestFlorBlucle(t *testing.T) {
	p, out, _ := NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})
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

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`la flor se debio de haber jugado`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestQuieroContraflorDesdeMazo(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
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

	assert(p.Ronda.Manojos[4].SeFueAlMazo == true, func() {
		t.Error(`andres se debio de haber ido al mazo`)
	})

	p.Cmd("renzo contra-flor")
	p.Cmd("andres quiero")

	assert(p.Ronda.Envite.CantadoPor.Jugador.Nombre == "Renzo", func() {
		t.Errorf(`andres no puede responder quiero porque se fue al mazo`)
	})

	t.Log(p.Ronda.Envite.Estado.String())

	assert(p.Ronda.Envite.Estado == pdt.CONTRAFLOR, func() {
		t.Error(`El estado del envite no debio de haber sido modificado`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestFixSeVaAlMazoYTeniaFlor(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Espada","valor":10},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Espada","valor":10},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":6},{"palo":"Espada","valor":10},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":7},{"palo":"Basto","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":12},{"palo":"Basto","valor":3},{"palo":"Espada","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Oro","valor":7},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":5},{"palo":"Basto","valor":11},{"palo":"Copa","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":5},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":5},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	ptsAzul := p.Puntajes[pdt.Azul]

	p.Cmd("alvaro mazo")

	assert(p.Ronda.Manojos[0].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	assert(ptsAzul == p.Puntajes[pdt.Azul], func() {
		t.Error(`no deberia de cambiar el puntaje`)
	})

	p.Cmd("roro truco")

	assert(p.Ronda.Envite.Estado == pdt.EstadoEnvite(pdt.TRUCO), func() {
		t.Error(`Deberia dejarlo cantar truco`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestFixDesconcertante(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":2},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":4},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":2},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":4},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":2},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":3},{"palo":"Espada","valor":11},{"palo":"Copa","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":7},{"palo":"Copa","valor":10},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":4},{"palo":"Oro","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Oro","valor":7},{"palo":"Copa","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":3},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":4},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro flor")
	p.Cmd("alvaro mazo")

	assert(p.Ronda.Manojos[0].SeFueAlMazo == false, func() {
		t.Error(`No deberia dejarlo irse al mazo porque se esta jugando la flor`)
	})
	p.Cmd("roro truco")

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestMalaAsignacionPts(t *testing.T) {
	p, out, _ := NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})

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

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)

	// pts Vale4 Ganado = 4
	// puntaje de Rojo deberia ser = 3 + 4
	// puntaje de Azul deberia ser = 2

	assert(p.Puntajes[pdt.Rojo] == 3+4 && p.Puntajes[pdt.Azul] == 2, func() {
		t.Error(`Asigno mal los puntos`)
	})
}

func TestFixRondaNueva(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Oro","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Oro","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":3},{"palo":"Copa","valor":1},{"palo":"Espada","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Espada","valor":4},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Oro","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":4},{"palo":"Basto","valor":7},{"palo":"Espada","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":12},{"palo":"Espada","valor":10},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("renzo flor")

	assert(p.Ronda.Envite.Estado == pdt.DESHABILITADO, func() {
		t.Error(`deberia deshabilitar el envite`)
	})

	p.Cmd("alvaro 1 copa")

	assert(p.Ronda.Manojos[0].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("roro truco")

	assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`Deberia dejarlo cantar truco`)
	})

	assert(p.Ronda.Truco.CantadoPor.Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`El equipo rojo deberia tener la potestad del truco`)
	})

	p.Cmd("adolfo re-truco")

	assert(p.Ronda.Truco.Estado == pdt.RETRUCO, func() {
		t.Error(`Deberia dejarlo cantar re-truco`)
	})

	assert(p.Ronda.Truco.CantadoPor.Jugador.Equipo == pdt.Azul, func() {
		t.Error(`El equipo azul deberia tener la potestad del truco`)
	})

	p.Cmd("renzo vale-4")

	assert(p.Ronda.Truco.Estado == pdt.VALE4, func() {
		t.Error(`Deberia dejarlo cantar vale-3`)
	})

	assert(p.Ronda.Truco.CantadoPor.Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`El equipo rojo deberia tener la potestad del truco`)
	})

	p.Cmd("adolfo quiero")

	assert(p.Ronda.Truco.Estado == pdt.VALE4QUERIDO, func() {
		t.Error(`Deberia dejarlo responder quiero al vale-4`)
	})

	assert(p.Ronda.Truco.CantadoPor.Jugador.Equipo == pdt.Azul, func() {
		t.Error(`El equipo azul deberia tener la potestad del truco`)
	})

	p.Cmd("roro 5 oro")

	assert(p.Ronda.Manojos[1].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("adolfo 6 basto")

	assert(p.Ronda.Manojos[2].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("renzo 19 oro")

	assert(p.Ronda.Manojos[3].GetCantCartasTiradas() == 0, func() {
		t.Error(`no debeia dejarlo porque no existe la carta "19 de oro"`)
	})

	p.Cmd("renzo 10 oro")

	assert(p.Ronda.Manojos[3].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("andres 4 copa")

	assert(p.Ronda.Manojos[4].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("richard 10 espada")

	assert(p.Ronda.Manojos[5].GetCantCartasTiradas() == 1, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	assert(p.Ronda.Manos[0].Ganador.Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`La primera mano la debio de haber ganado el equipo de richard: el rojo`)
	})

	// segunda mano
	p.Cmd("renzo 10 oro")

	assert(p.Ronda.Manojos[3].GetCantCartasTiradas() == 1, func() {
		t.Error(`ya tiro esa carta; no deberia poder volve a tirarla`)
	})

	p.Cmd("andres 7 basto")

	assert(p.Ronda.Manojos[4].GetCantCartasTiradas() == 1, func() {
		t.Error(`no es su turno no deberia poder tirar carta`)
	})

	p.Cmd("richard 10 espada")

	assert(p.Ronda.Manojos[5].GetCantCartasTiradas() == 1, func() {
		t.Error(`ya tiro esa carta; no deberia poder volve a tirarla`)
	})

	p.Cmd("richard 12 copa")

	assert(p.Ronda.Manojos[5].GetCantCartasTiradas() == 2, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("alvaro 3 espada")

	assert(p.Ronda.Manojos[0].GetCantCartasTiradas() == 2, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("roro 4 espada")

	assert(p.Ronda.Manojos[1].GetCantCartasTiradas() == 2, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("adolfo 2 copa")

	assert(p.Ronda.Manojos[2].GetCantCartasTiradas() == 2, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("renzo 4 oro")

	assert(p.Ronda.Manojos[3].GetCantCartasTiradas() == 2, func() {
		t.Error(`deberia dejarlo tirar la carta`)
	})

	p.Cmd("andres 5 espada")

	assert(p.Ronda.Manojos[4].GetCantCartasTiradas() == 0, func() {
		t.Error(`deberia tener 0 cartas tiradas porque empieza una nueva ronda`)
	})

	/* 3:flor + 4:vale4 */
	assert(p.Puntajes[pdt.Rojo] == 3+4, func() {
		t.Error(`el puntaje para el equipo rojo deberia ser 7: 3 de la flor + 4 del vale4`)
	})

	assert(p.Puntajes[pdt.Azul] == 0, func() {
		t.Error(`el puntaje para el equipo azul deberia ser 0 porque no ganaron nada`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)
}

func TestFixIrseAlMazo2(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Espada","valor":7},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":3},{"palo":"Copa","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Oro","valor":5},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":12},{"palo":"Oro","valor":2},{"palo":"Copa","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":7},{"palo":"Basto","valor":7},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Copa","valor":3},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("renzo flor")

	assert(p.Ronda.Envite.Estado == pdt.NOCANTADOAUN, func() {
		t.Error(`Renzo no tiene flor`)
	})

	// mano 1
	p.Cmd("alvaro envido")

	assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`Debio dejarlo cantar truco`)
	})

	p.Cmd("alvaro 6 basto")

	assert(p.Ronda.Manojos[0].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	p.Cmd("roro 11 oro")

	assert(p.Ronda.Manojos[1].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	p.Cmd("adolfo 2 basto")

	assert(p.Ronda.Manojos[2].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	p.Cmd("renzo 12 oro")

	assert(p.Ronda.Manojos[3].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	p.Cmd("andres 7 copa")

	assert(p.Ronda.Manojos[4].GetCantCartasTiradas() == 0, func() {
		t.Error(`A este no lo deberia dejar tirar carta`)
	})

	p.Cmd("richard 12 basto")

	assert(p.Ronda.Manojos[5].GetCantCartasTiradas() == 0, func() {
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

	assert(p.Ronda.Manojos[1].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})

	t.Log(p)

}

func TestFixDecirQuiero(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":10},{"palo":"Espada","valor":12},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":10},{"palo":"Espada","valor":12},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":1},{"palo":"Copa","valor":6},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":7},{"palo":"Copa","valor":12},{"palo":"Oro","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":2},{"palo":"Espada","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":10},{"palo":"Espada","valor":12},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":3},{"palo":"Basto","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Basto","valor":6},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":2},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("renzo flor")
	p.Cmd("alvaro truco")

	assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`Deberia poder gritar truco`)
	})

	p.Cmd("renzo quiero")

	assert(p.Ronda.Truco.Estado == pdt.TRUCOQUERIDO, func() {
		t.Error(`Deberia poder responder quiero al truco`)
	})

	p.Cmd("alvaro re-truco")

	assert(p.Ronda.Truco.Estado == pdt.TRUCOQUERIDO, func() {
		t.Error(`Como no tiene la potestad, no deberia poder aumentar la apuesta`)
	})

	assert(p.Ronda.Truco.CantadoPor.Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`El equpo Rojo deberia de seguir manteniendo la potestad`)
	})

	p.Cmd("renzo re-truco")
	p.Cmd("alvaro vale-4")

	assert(p.Ronda.Truco.Estado == pdt.VALE4, func() {
		t.Error(`Deberia poder aumentar a vale-4`)
	})

	p.Cmd("alvaro quiero")

	assert(p.Ronda.Truco.Estado == pdt.VALE4, func() {
		t.Error(`No puede auto-querse`)
	})

	assert(p.Ronda.Truco.CantadoPor.Jugador.Equipo == pdt.Azul, func() {
		t.Error(`El equpo azul deberia tener la potestad`)
	})

	p.Cmd("renzo re-truco")

	assert(p.Ronda.Truco.Estado == pdt.VALE4, func() {
		t.Error(`No deberia cambiar el estado del truco`)
	})

	assert(p.Ronda.Truco.CantadoPor.Jugador.Equipo == pdt.Azul, func() {
		t.Error(`El equpo azul deberia de seguir manteniendo la potestad`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)

}

func TestFixPanicNoQuiero(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":1},{"palo":"Basto","valor":10},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":4},{"palo":"Espada","valor":6},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Copa","valor":5},{"palo":"Basto","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Oro","valor":12},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":7},{"palo":"Oro","valor":11},{"palo":"Oro","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":2},{"palo":"Basto","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":1},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
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

	assert(p.Ronda.Manos[0].Ganador.Jugador.Equipo == pdt.Rojo, func() {
		t.Error(`La primera mano la debio de haber ganado el equipo de renzo: el rojo`)
	})

	p.Cmd("renzo 12 oro")
	p.Cmd("andres 11 oro") // la seguna mano la gana andres
	p.Cmd("richard 2 copa")
	p.Cmd("alvaro 10 basto")
	p.Cmd("roro 6 espada")
	p.Cmd("adolfo 5 copa")

	assert(p.Ronda.Manos[1].Ganador.Jugador.Equipo == pdt.Azul, func() {
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

	assert(p.Puntajes[pdt.Rojo] == ptsPostFlor+3, func() {
		t.Error(`Deberian gana 3 puntines por el vale-4 no querido`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)

}

func TestFixCartaYaJugada(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"Jugadores":[{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"},{"id":"Roro","nombre":"Roro","equipo":"Rojo"},{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"},{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"},{"id":"Andres","nombre":"Andres","equipo":"Azul"},{"id":"Richard","nombre":"Richard","equipo":"Rojo"}],"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":2},{"palo":"Oro","valor":4},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":4},{"palo":"Espada","valor":5},{"palo":"Basto","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Espada","valor":3},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Oro","valor":2},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":7},{"palo":"Basto","valor":3},{"palo":"Copa","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Oro","valor":6},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":11},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro 2 espada")
	p.Cmd("roro 4 copa")
	p.Cmd("adolfo 10 copa")
	p.Cmd("renzo 6 copa")
	p.Cmd("andres 7 copa")
	p.Cmd("richard flor")
	p.Cmd("richard 5 oro")
	p.Cmd("richard 5 oro")

	assert(p.Ronda.GetElTurno().Jugador.Nombre == "Richard", func() {
		t.Error(`Deberia ser el turno de Richard`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(p)

}

func TestFixTrucoNoQuiero(t *testing.T) {
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"Jugadores":[{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"},{"id":"Roro","nombre":"Roro","equipo":"Rojo"},{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"},{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"},{"id":"Andres","nombre":"Andres","equipo":"Azul"},{"id":"Richard","nombre":"Richard","equipo":"Rojo"}],"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Copa","valor":7},{"palo":"Espada","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":12},{"palo":"Basto","valor":2},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":4},{"palo":"Copa","valor":5},{"palo":"Copa","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":7},{"palo":"Espada","valor":3},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":10},{"palo":"Espada","valor":2},{"palo":"Copa","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":7},{"palo":"Oro","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":10},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro truco")
	p.Cmd("roro no-quiero")

	assert(p.Puntajes[pdt.Azul] > 0, func() {
		t.Error(`La ronda deberia de haber sido ganado por Azul`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})
	t.Log(enco.Collect(out))
	t.Log(p)

}

func TestPerspectiva(t *testing.T) {
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))

	per, _ := p.Perspectiva("Alvaro")
	t.Log(per.MarshalJSON())
}

func TestPardaSigTurno1(t *testing.T) {
	// si va parda, el siguiente turno deberia ser del mano
	// o del mas cercano a este
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
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

	p.Cmd("Alvaro 5 Copa")
	p.Cmd("Roro 4 Copa")
	// los siguientes 4: todos tiran 6 -> resulta mano parda
	p.Cmd("Adolfo 6 Espada")
	p.Cmd("Renzo 6 Oro")
	p.Cmd("Andres 6 Copa")
	p.Cmd("Richard 6 Basto")

	t.Log(p)

	assert(p.Ronda.Manos[0].Resultado == pdt.Empardada, func() {
		t.Error(`La mano debio ser parda`)
	})

	assert(p.Ronda.GetElTurno().Jugador.Nombre == "Adolfo", func() {
		t.Error(`Deberia ser turno de Adolfo, debido a que es el mas cercano del mano y qu empardo`)
	})
}

func TestPardaSigTurno2(t *testing.T) {
	// igual que el anterior pero ahora adolfo se va al mazo
	// si va parda, el siguiente turno deberia ser del mano
	// o del mas cercano a este
	p, _, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
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

	p.Cmd("Alvaro 5 Copa")
	p.Cmd("Roro 4 Copa")
	// los siguientes 4: todos tiran 6 -> resulta mano parda
	p.Cmd("Adolfo 6 Espada")
	p.Cmd("Adolfo mazo")
	p.Cmd("Renzo 6 Oro")
	p.Cmd("Andres 6 Copa")
	p.Cmd("Richard 6 Basto")

	t.Log(p)

	assert(p.Ronda.Manos[0].Resultado == pdt.Empardada, func() {
		t.Error(`La mano debio ser parda`)
	})

	assert(p.Ronda.GetElTurno().Jugador.Nombre == "Renzo", func() {
		t.Error(`Deberia ser turno de Renzo, debido a que es el mas cercano del mano y qu empardo`)
	})
}

func TestPardaSigTurno3(t *testing.T) {
	// igual que el anterior pero ahora todos los que empardaron se van al mazo
	// si va parda, el siguiente turno deberia ser del mano
	// o del mas cercano a este
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
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
					{Palo: pdt.Basto, Valor: 4},
					{Palo: pdt.Copa, Valor: 11},
					{Palo: pdt.Basto, Valor: 6}, // <-- parda primera mano
				},
			},
		},
	)

	p.Cmd("Alvaro 5 Copa")

	assert(contains(enco.Collect(out), enco.TirarCarta), func() {
		t.Error("debio de haber tirado carta")
	})

	p.Cmd("Roro 4 Copa")
	// los siguientes 4: todos tiran 6 -> resulta mano parda
	p.Cmd("Adolfo 6 Espada")
	p.Cmd("Adolfo mazo")
	p.Cmd("Renzo 6 Oro")
	p.Cmd("Renzo mazo")
	p.Cmd("Andres 6 Copa")
	p.Cmd("Andres mazo")
	p.Cmd("Richard 4 Basto")
	p.Cmd("Richard mazo")

	// ojo no imprime nada porque la aterior ya consumio el out
	// enco.Consume(out, func(pkt *enco.Packet) {
	// 	t.Log(deco.Stringify(pkt, p.PartidaDT))
	// })

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})

	t.Log(p)

	assert(p.Ronda.Manos[0].Resultado == pdt.Empardada, func() {
		t.Error(`La mano debio ser parda`)
	})

	assert(p.Ronda.GetElTurno().Jugador.Nombre == "Roro", func() {
		t.Error(`Deberia ser turno de Roro, debido a que es el mas cercano del mano y qu empardo`)
	})

	assert(p.Ronda.Manojos[5].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})
}

func TestFixTrucoDeshabilitaEnvido(t *testing.T) {
	// cantar truco (sin siquiera ser querido) deshabilita el envido
	// cuando en la vida real es posible tocar "el envido esta primero"
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":7},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Espada","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Oro","valor":2},{"palo":"Espada","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("Alvaro truco")

	assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`Deberia dejarlo gritar truco`)
	})

	// el envido esta primero!!
	p.Cmd("Roro envido")

	assert(!contains(enco.Collect(out), enco.Error), func() {
		t.Error("No deberia resultar en un error tocar envido ahora")
	})

	assert(p.Ronda.Envite.Estado == pdt.ENVIDO, func() {
		t.Error(`Deberia dejarlo tocar envido`)
	})

}

func TestAbandono(t *testing.T) {
	// simulacro de un jugador abandonando
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":2,"Rojo":3},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Copa","valor":4},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":7},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Espada","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Oro","valor":2},{"palo":"Espada","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":11},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("Alvaro truco")

	assert(p.Ronda.Truco.Estado == pdt.TRUCO, func() {
		t.Error(`Deberia dejarlo gritar truco`)
	})

	// el envido esta primero!!
	p.Abandono("Adolfo")

	assert(p.Terminada(), func() {
		t.Error(`Deberia haber acabado la partida`)
	})

	assert(p.Puntajes[pdt.Rojo] == int(p.Puntuacion), func() {
		t.Error(`El equipo rojo debio de haber completado los pts`)
	})

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})

}

func TestFixOrdenCantoFlor(t *testing.T) {
	// simulacro de un jugador abandonando
	p, out, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"Jugadores":[{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"},{"id":"Roro","nombre":"Roro","equipo":"Rojo"},{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"},{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"},{"id":"Andres","nombre":"Andres","equipo":"Azul"},{"id":"Richard","nombre":"Richard","equipo":"Rojo"}],"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Espada","valor":7},{"palo":"Espada","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Espada","valor":11},{"palo":"Espada","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":10},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Oro","valor":3},{"palo":"Oro","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":4},{"palo":"Copa","valor":3},{"palo":"Copa","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":4},{"palo":"Copa","valor":12},{"palo":"Basto","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":7},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("alvaro flor")
	p.Cmd("renzo flor")
	p.Cmd("andres flor")

	assert(contains(enco.Collect(out), enco.DiceSonBuenas), func() {
		t.Error("debio de haber dicho son bueas")
	})

}

func TestFixTester2(t *testing.T) {
	// simulacro de un jugador abandonando
	p, out, _ := NuevaPartida(pdt.A30, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
	partidaJSON := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 30, "puntajes": {"Azul": 27, "Rojo": 23}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 2, "Rojo": 2}, "elMano": 2, "turno": 3, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": null, "estado": "noCantado"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Oro", "valor": 12}, {"palo": "Basto", "valor": 5}, {"palo": "Basto", "valor": 12}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Oro", "valor": 3}, {"palo": "Espada", "valor": 1}, {"palo": "Oro", "valor": 7}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Oro", "valor": 6}, {"palo": "Copa", "valor": 7}, {"palo": "Espada", "valor": 12}], "cartasNoJugadas": [true, false, true], "ultimaTirada": 1, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 3}, {"palo": "Copa", "valor": 12}, {"palo": "Basto", "valor": 3}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Oro", "valor": 5}, "manos": [{"resultado": "ganoRojo", "ganador": null, "cartasTiradas": [{"palo": "Copa", "valor": 7}]}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}]}}`
	p.PartidaDT.FromJSON([]byte(partidaJSON))
	t.Log(p)

	p.Cmd("Adolfo no-quiero")
	p.Cmd("Adolfo 7 Copa")
	p.Cmd("Roro 3 Oro")
	p.Cmd("Renzo contra-flor-al-resto")
	p.Cmd("Roro 7 Oro")
	p.Cmd("Alvaro flor")

	enco.Consume(out, func(pkt *enco.Packet) {
		t.Log(deco.Stringify(pkt, p.PartidaDT))
	})

}
