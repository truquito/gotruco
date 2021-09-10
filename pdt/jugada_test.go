package pdt

import (
	"testing"

	"github.com/filevich/truco/enco"
)

func toStr(c enco.CodMsg) string {
	var s string
	switch c {
	case enco.Error:
		s = "Error"
	case enco.ByeBye:
		s = "ByeBye"
	case enco.DiceSonBuenas:
		s = "DiceSonBuenas"
	case enco.CantarFlor:
		s = "CantarFlor"
	case enco.CantarContraFlor:
		s = "CantarContraFlor"
	case enco.CantarContraFlorAlResto:
		s = "CantarContraFlorAlResto"
	case enco.TocarEnvido:
		s = "TocarEnvido"
	case enco.TocarRealEnvido:
		s = "TocarRealEnvido"
	case enco.TocarFaltaEnvido:
		s = "TocarFaltaEnvido"
	case enco.GritarTruco:
		s = "GritarTruco"
	case enco.GritarReTruco:
		s = "GritarReTruco"
	case enco.GritarVale4:
		s = "GritarVale4"
	case enco.NoQuiero:
		s = "NoQuiero"
	case enco.ConFlorMeAchico:
		s = "ConFlorMeAchico"
	case enco.QuieroTruco:
		s = "QuieroTruco"
	case enco.QuieroEnvite:
		s = "QuieroEnvite"
	case enco.SigTurno:
		s = "SigTurno"
	case enco.SigTurnoPosMano:
		s = "SigTurnoPosMano"
	case enco.DiceTengo:
		s = "DiceTengo"
	case enco.DiceSonMejores:
		s = "DiceSonMejores"
	case enco.NuevaPartida:
		s = "NuevaPartida"
	case enco.NuevaRonda:
		s = "NuevaRonda"
	case enco.TirarCarta:
		s = "TirarCarta"
	case enco.SumaPts:
		s = "SumaPts"
	case enco.Mazo:
		s = "Mazo"
	case enco.TimeOut:
		s = "TimeOut"
	case enco.ElEnvidoEstaPrimero:
		s = "ElEnvidoEstaPrimero"
	case enco.Abandono:
		s = "Abandono"
	case enco.LaManoResultaParda:
		s = "LaManoResultaParda"
	case enco.ManoGanada:
		s = "ManoGanada"
	case enco.RondaGanada:
		s = "RondaGanada"
	}
	return s
}

func TestAcciones(t *testing.T) {
	pdt, _ := NuevaPartidaDt(A20, []string{"Alvaro"}, []string{"Roro"})

	t.Log(Renderizar(pdt))

	alvaro := pdt.Ronda.GetElTurno()

	for _, a := range pdt.A(alvaro) {
		m, _ := a.Cont.MarshalJSON()
		t.Log(toStr(enco.CodMsg(a.Cod)), string(m))
	}
	// t.Log(pdt.A(alvaro))

	// assert(p.Ronda.Envite.Puntaje == 2, func() {
	// 	t.Error(`El puntaje del envido deberia de ser 2`)
	// })
}
