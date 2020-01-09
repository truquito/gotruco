package truco

import (
	"testing"
)

var (
	CincoDeCopa = Carta{
		Palo:  Copa,
		Valor: 5,
	}

	manojos2 = []Manojo{
		Manojo{
			Cartas: [3]Carta{
				Carta{Palo: Oro, Valor: 6},
				Carta{Palo: Copa, Valor: 10},
				Carta{Palo: Copa, Valor: 7},
			},
			jugador: nil,
		},
		Manojo{
			Cartas: [3]Carta{
				Carta{Palo: Copa, Valor: 1},
				Carta{Palo: Oro, Valor: 2},
				Carta{Palo: Basto, Valor: 3},
			},
			jugador: nil,
		},
	}
)

func TestNoDeberianTenerFlor(t *testing.T) {

	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Copa, Valor: 5})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Oro, Valor: 6},
					Carta{Palo: Copa, Valor: 10},
					Carta{Palo: Copa, Valor: 7},
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

	tieneFlor, _ := p.Ronda.Manojos[0].tieneFlor(p.Ronda.Muestra)
	oops = tieneFlor == true
	if oops {
		t.Error(`Alvaro' NO deberia de tener 'flor'`)
		return
	}

	tieneFlor, _ = p.Ronda.Manojos[1].tieneFlor(p.Ronda.Muestra)
	oops = tieneFlor == true
	if oops {
		t.Error(`Roro' NO deberia de tener 'flor'`)
		return
	}

}

func TestDeberiaTenerFlor(t *testing.T) {

	p, _ := NuevaPartida(a20, []string{"Alvaro"}, []string{"Roro"})
	p.Ronda.setMuestra(Carta{Palo: Copa, Valor: 5})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Copa, Valor: 4},
					Carta{Palo: Espada, Valor: 10},
					Carta{Palo: Espada, Valor: 7},
				},
			},
			Manojo{
				Cartas: [3]Carta{
					Carta{Palo: Oro, Valor: 1},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Oro, Valor: 3},
				},
			},
		},
	)

	tieneFlor, _ := p.Ronda.Manojos[0].tieneFlor(p.Ronda.Muestra)
	oops = !(tieneFlor == true)
	if oops {
		t.Error(`Alvaro' deberia tener 'flor'`)
		return
	}

	tieneFlor, _ = p.Ronda.Manojos[1].tieneFlor(p.Ronda.Muestra)
	oops = !(tieneFlor == true)
	if oops {
		t.Error(`Roro' deberia tener 'flor'`)
		return
	}
}

func TestFlorFlorContraFlorQuiero(t *testing.T) {

	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	p.Ronda.setMuestra(Carta{Palo: Oro, Valor: 3})
	p.Ronda.setManojos(
		[]Manojo{
			Manojo{
				Cartas: [3]Carta{ // Alvaro tiene flor
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 6},
					Carta{Palo: Basto, Valor: 7},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // Roro no tiene flor
					Carta{Palo: Oro, Valor: 5},
					Carta{Palo: Espada, Valor: 5},
					Carta{Palo: Basto, Valor: 5},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // Adolfo tiene flor
					Carta{Palo: Copa, Valor: 1},
					Carta{Palo: Copa, Valor: 2},
					Carta{Palo: Copa, Valor: 3},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // Renzo tiene flor
					Carta{Palo: Oro, Valor: 4},
					Carta{Palo: Espada, Valor: 4},
					Carta{Palo: Espada, Valor: 1},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // Andres no tiene  flor
					Carta{Palo: Copa, Valor: 10},
					Carta{Palo: Oro, Valor: 7},
					Carta{Palo: Basto, Valor: 11},
				},
			},
			Manojo{
				Cartas: [3]Carta{ // Richard tiene flor
					Carta{Palo: Oro, Valor: 10},
					Carta{Palo: Oro, Valor: 2},
					Carta{Palo: Basto, Valor: 1},
				},
			},
		},
	)

	p.SetSigJugada("Alvaro Flor")
	p.SetSigJugada("Roro Mazo")
	p.SetSigJugada("Renzo Flor")
	p.SetSigJugada("Adolfo Contra-flor-al-resto")
	p.SetSigJugada("Richard Quiero")
	p.Esperar()

	oops = p.Ronda.Envido.estado != DESHABILITADO
	if oops {
		t.Error(`El estado del envido deberia ser 'deshabilitado'`)
		return
	}

	oops = p.Ronda.Flor != DESHABILITADA
	if oops {
		t.Error(`El estado de la flor deberia ser 'deshabilitado'`)
		return
	}

	// duda: se suman solo las flores ganadoras
	// si contraflor AL RESTO -> no acumulativo
	oops = !(p.Puntajes[Azul] == 4*3+10 && p.Puntajes[Rojo] == 0)
	if oops {
		t.Error(`El puntaje del equipo azul deberia ser 2`)
		return
	}
}

// Tests:
// los "me achico" no cuentan para la flor
// Flor		xcg(+3) / xcg(+3)
// Flor + Contra-Flor		xc(+3) / xCadaFlorDelQueHizoElDesafio(+3) + 1
// Flor + [Contra-Flor] + ContraFlorAlResto		~Falta Envido + *TODAS* las flores no achicadas / xcg(+3) + 1
