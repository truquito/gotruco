package util

import (
	"testing"

	"github.com/filevich/truco/pdt"
)

func TestAcciones(t *testing.T) {
	p, _ := pdt.NuevaPartidaDt(pdt.A20, []string{"Alvaro"}, []string{"Roro"})

	alvaro := p.Ronda.GetElTurno()
	a := pdt.TocarEnvido{Manojo: alvaro}
	// a := pdt.IrseAlMazo{Manojo: alvaro}
	a.Hacer(p)

	t.Log(pdt.Renderizar(p))

	// t.Log(GetA(p, alvaro))
	as := GetAA(p)
	t.Log(pdt.Renderizar(p))
	for i, a := range as {
		t.Logf("%s : %v", p.Jugadores[i].Nombre, a)
	}
}

func TestTodasLasAcciones(t *testing.T) {
	p, _ := pdt.NuevaPartidaDt(pdt.A20, []string{"Alvaro"}, []string{"Roro"})
	as := GetAA(p)
	t.Log(pdt.Renderizar(p))
	for i, a := range as {
		t.Logf("%s : %v", p.Jugadores[i].Nombre, a)
	}
}
