package pdt

import (
	"testing"
)

func TestFlor(t *testing.T) {
	muestra := Carta{Palo: Oro, Valor: 11}
	m := &Manojo{
		Cartas: [3]*Carta{
			{Palo: Oro, Valor: 2},
			{Palo: Oro, Valor: 4},
			{Palo: Oro, Valor: 5},
		},
	}

	if tiene, _ := m.TieneFlor(muestra); !tiene {
		t.Error("Deberia tener flor")
	}

	if ptsFlor, err := m.CalcFlor(muestra); ptsFlor != 47 || err != nil {
		t.Error("Deberia tener 47 de flor")
	}
}

func TestCartaID(t *testing.T) {
	t.Log(len(primes))

	c1 := Carta{Basto, 5}
	ok := c1.ID() == CartaID(4) && c1.PUID() == 11
	if !ok {
		t.Errorf("el ID deberia ser 4, obtuve %v", c1.ID())
	}

	c2 := Carta{Espada, 11}
	ok = c2.ID() == CartaID(28) && c2.PUID() == 109
	if !ok {
		t.Errorf("el ID deberia ser 28, obtuve %v", c2.ID())
	}

	c3 := Carta{Basto, 1}
	ok = c3.ID() == CartaID(0) && c3.PUID() == 2
	if !ok {
		t.Errorf("el ID deberia ser 0, obtuve %v", c3.ID())
	}

	c4 := Carta{Oro, 12}
	ok = c4.ID() == CartaID(39) && c4.PUID() == 173
	if !ok {
		t.Errorf("el ID deberia ser 0, obtuve %v", c4.ID())
	}

	c5 := Carta{Oro, 11}
	ok = c5.ID() == CartaID(38) && c5.PUID() == 167
	if !ok {
		t.Errorf("el ID deberia ser 38, obtuve %v", c5.ID())
	}

	c6 := Carta{Oro, 4}
	ok = c6.ID() == CartaID(33) && c6.PUID() == 139
	if !ok {
		t.Errorf("el ID deberia ser 33, obtuve %v", c6.ID())
	}

	c7 := Carta{Copa, 1}
	ok = c7.ID() == CartaID(10) && c7.PUID() == 31
	if !ok {
		t.Errorf("el ID deberia ser 10, obtuve %v", c7.ID())
	}

	t.Log("asd")
}
