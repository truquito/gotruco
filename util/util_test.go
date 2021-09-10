package util

import (
	"testing"

	"github.com/filevich/truco/pdt"
)

func TestAcciones(t *testing.T) {
	p, _ := pdt.NuevaPartidaDt(pdt.A20, []string{"Alvaro"}, []string{"Roro"})

	alvaro := p.Ronda.GetElTurno()
	a := pdt.TirarCarta{Manojo: alvaro, Carta: *alvaro.Cartas[0]}
	a.Hacer(p)

	t.Log(pdt.Renderizar(p))

	t.Log(GetA(p, alvaro))
}
