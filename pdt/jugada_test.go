package pdt

import (
	"testing"
)

func TestJugada(t *testing.T) {
	pdt, _ := NuevaPartidaDt(A20, []string{"Alvaro"}, []string{"Roro"})

	t.Log(pdt.CantJugadores)

	// assert(p.Ronda.Envite.Puntaje == 2, func() {
	// 	t.Error(`El puntaje del envido deberia de ser 2`)
	// })
}
