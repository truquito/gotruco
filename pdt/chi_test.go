package pdt

import (
	"testing"
)

func TestAcciones(t *testing.T) {
	p, _ := NuevaPartida(A20, []string{"Alvaro"}, []string{"Roro"}, 4, true)

	alvaro := p.Ronda.GetElTurno()
	a := TocarEnvido{JID: alvaro.Jugador.ID}
	// a := IrseAlMazo{Manojo: alvaro}
	a.Hacer(p)

	t.Log(Renderizar(p))

	// t.Log(GetA(p, alvaro))
	as := Chis(p)
	t.Log(Renderizar(p))
	for i, a := range as {
		t.Logf("%s : %v", p.Ronda.Manojos[i].Jugador.ID, a)
	}
}

func TestTodasLasAcciones(t *testing.T) {
	p, _ := NuevaPartida(A20, []string{"Alvaro"}, []string{"Roro"}, 4, true)
	as := Chis(p)
	t.Log(Renderizar(p))
	for i, a := range as {
		t.Logf("%s : %v", p.Ronda.Manojos[i].Jugador.ID, a)
	}
}
