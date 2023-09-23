package pdt

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/util"
)

func TestNoStruct(t *testing.T) {
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"sinCantar":["Renzo","Richard"]},"truco":{"cantadoPor":"","estado":"truco"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"basto","valor":12},{"palo":"oro","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"basto","valor":10},{"palo":"oro","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":10},{"palo":"copa","valor":10},{"palo":"basto","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"espada","valor":10},{"palo":"basto","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"copa","valor":3},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"espada","valor":11},{"palo":"espada","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":6},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p, err := Parse(partidaJSON, true)
	if err != nil {
		t.Error(err)
	}

	if !(p.Ronda.Muestra.Palo == Espada && p.Ronda.Muestra.Valor == 6) {
		t.Error("El valor de la muestra no es correcta")
	}

	if !(EstadoEnvite(p.Ronda.Truco.Estado) == EstadoEnvite(TRUCO)) {
		t.Error("El estado del truco deberia ser TRUCO")
	}

	if !(p.Ronda.Envite.Estado == NOCANTADOAUN) {
		t.Error("El estado del envite deberia ser NOCANTADOAUN")
	}

	if !(p.Ronda.Envite.Estado == NOCANTADOAUN) {
		t.Error("El estado del envite deberia ser NOCANTADOAUN")
	}

	if !(len(p.Ronda.Manojos) == 6) {
		t.Error("deberia haber 6 jugadores")
	}

	t.Log(p.Ronda.Muestra)

	t.Log(Renderizar(p))
}

func TestFixPanic(t *testing.T) {
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true)
	p.Ronda.SetMuestra(Carta{Palo: Basto, Valor: 4})
	p.Puntajes[Rojo] = 17
	p.Puntajes[Azul] = 10
	p.Ronda.ManoEnJuego = Primera
	p.Ronda.ElMano = 3 // renzo
	p.Ronda.Turno = 3  // renzo
	p.Ronda.SetManojos(
		[]Manojo{
			{
				Cartas: [3]*Carta{ // cartas de Alvaro
					{Palo: Espada, Valor: 7},
					{Palo: Copa, Valor: 4},
					{Palo: Espada, Valor: 11},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas Roro
					{Palo: Copa, Valor: 6},
					{Palo: Copa, Valor: 1},
					{Palo: Basto, Valor: 10},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Adolfo
					{Palo: Oro, Valor: 7},
					{Palo: Copa, Valor: 11},
					{Palo: Copa, Valor: 3},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Renzo
					{Palo: Basto, Valor: 3},
					{Palo: Copa, Valor: 7},
					{Palo: Basto, Valor: 11},
				},
			},
		},
	)
	p.Cmd("Renzo 11 basto")
	t.Log(Renderizar(p)) // retorna el render
}

func TestFixNoLePermiteTocarEnvido(t *testing.T) {
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true)
	p.Ronda.SetMuestra(Carta{Palo: Basto, Valor: 4})
	p.Puntajes[Rojo] = 0
	p.Puntajes[Azul] = 0
	p.Ronda.ManoEnJuego = Primera
	p.Ronda.ElMano = 0 // renzo
	p.Ronda.Turno = 0  // renzo
	p.Ronda.SetManojos(
		[]Manojo{
			{
				Cartas: [3]*Carta{ // cartas de Alvaro
					{Palo: Oro, Valor: 2},
					{Palo: Oro, Valor: 6},
					{Palo: Basto, Valor: 10},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas Roro
					{Palo: Basto, Valor: 4},
					{Palo: Espada, Valor: 12},
					{Palo: Oro, Valor: 10},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Adolfo
					{Palo: Basto, Valor: 11},
					{Palo: Oro, Valor: 11},
					{Palo: Copa, Valor: 7},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Renzo
					{Palo: Basto, Valor: 7},
					{Palo: Espada, Valor: 5},
					{Palo: Oro, Valor: 12},
				},
			},
		},
	)

	pkts, _ := p.Cmd("Alvaro envido")

	util.Assert(enco.Contains(pkts, enco.TError), func() {
		t.Error("No debio de haberle dejado tocar envido")
	})
}

func TestFixNoLePermiteGritarTruco(t *testing.T) {
	data := `{"limiteEnvido":4,"cantJugadores":4,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"oro","valor":7},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"basto","valor":5},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"oro","valor":3},{"palo":"espada","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":1},{"palo":"espada","valor":6},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":3},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, err := Parse(data, true)
	if err != nil {
		t.Error(err)
	}

	pkts, _ := p.Cmd("Alvaro truco")

	util.Assert(enco.Contains(pkts, enco.TGritarTruco), func() {
		t.Error("Deberia dejarlo gritar truco!")
	})
}

func TestFixLePermiteCFAR(t *testing.T) {
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true)
	p.Ronda.SetMuestra(Carta{Palo: Copa, Valor: 3})
	p.Puntajes[Rojo] = 1
	p.Puntajes[Azul] = 6
	p.Ronda.ManoEnJuego = Primera
	p.Ronda.ElMano = 2 // renzo
	p.Ronda.Turno = 2  // renzo
	p.Ronda.SetManojos(
		[]Manojo{
			{
				Cartas: [3]*Carta{ // cartas de Alvaro
					{Palo: Copa, Valor: 1},
					{Palo: Copa, Valor: 5},
					{Palo: Copa, Valor: 7},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas Roro
					{Palo: Copa, Valor: 2},
					{Palo: Oro, Valor: 1},
					{Palo: Basto, Valor: 10},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Adolfo
					{Palo: Espada, Valor: 1},
					{Palo: Espada, Valor: 5},
					{Palo: Oro, Valor: 2},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Renzo
					{Palo: Oro, Valor: 12},
					{Palo: Oro, Valor: 5},
					{Palo: Copa, Valor: 11},
				},
			},
		},
	)

	m := p.Manojo("Renzo")
	A := GetA(p, m)

	ok := A[6] == true && A[7] == false && A[8] == false
	util.Assert(ok, func() {
		t.Error("No deberia dejarlo cantar algo mas poderoso que `flor`")
	})

	pkts, _ := p.Cmd("Renzo contra-flor-al-resto")
	util.Assert(enco.Contains(pkts, enco.TError), func() {
		t.Error("No deberia dejarlo cantar CFAR")
	})

	pkts, _ = p.Cmd("Renzo no-quiero")
	util.Assert(enco.Contains(pkts, enco.TError), func() {
		t.Error("No deberia dejarlo responder no-quiero")
	})
}

func TestFixNoLeDebePermitirTocarEnvidoSiGritoTruco(t *testing.T) {
	data := `{"limiteEnvido":4,"cantJugadores":4,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"oro","valor":7},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"basto","valor":5},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"oro","valor":3},{"palo":"espada","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":1},{"palo":"espada","valor":6},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":3},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, err := Parse(data, true)
	if err != nil {
		t.Error(err)
	}
	t.Log(Renderizar(p))

	pkts, _ := p.Cmd("Alvaro truco")
	util.Assert(enco.Contains(pkts, enco.TGritarTruco), func() {
		t.Error("Deberia dejarlo gritar truco!")
	})

	pkts, _ = p.Cmd("Alvaro envido")
	ok := util.All(
		p.Ronda.Envite.Estado == NOCANTADOAUN,
		p.Ronda.Truco.Estado != NOGRITADOAUN,
		enco.Contains(pkts, enco.TError),
	)
	util.Assert(ok, func() {
		t.Error("No deberia dejarlo tocar envido si fue el mismo el que toco truco!")
	})

	pkts, _ = p.Cmd("Alvaro real-envido")
	ok = util.All(
		p.Ronda.Envite.Estado == NOCANTADOAUN,
		p.Ronda.Truco.Estado != NOGRITADOAUN,
		enco.Contains(pkts, enco.TError),
	)
	util.Assert(ok, func() {
		t.Error("No deberia dejarlo tocar real-envido si fue el mismo el que toco truco!")
	})

}

func TestFixNoLeDeberiaResponderDesdeUltratumba(t *testing.T) {
	data := `{"limiteEnvido":4,"cantJugadores":4,"puntuacion":20,"puntajes":{"azul":5,"rojo":6},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":2,"turno":2,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"basto","valor":6},{"palo":"espada","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"espada","valor":11},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"espada","valor":10},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":1},{"palo":"basto","valor":11},{"palo":"espada","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":5},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, err := Parse(data, true)
	if err != nil {
		t.Error(err)
	}
	t.Log(Renderizar(p))

	p.Cmd("Renzo mazo")
	m := p.Manojo("Renzo")
	util.Assert(m.SeFueAlMazo, func() {
		t.Error("Deberia dejarlo irse al mazo")
	})

	p.Cmd("Adolfo falta-envido")

	pkts, _ := p.Cmd("Roro quiero")
	for _, pkt := range pkts {
		diceSonBuenas := pkt.Message.Cod() == enco.TDiceSonBuenas
		m, _ := pkt.Message.(enco.DiceSonBuenas)
		loDijoRenzo := m == "Renzo"
		if diceSonBuenas && loDijoRenzo {
			t.Error("No deberia poder responder desde ultratumba")
		}
	}

	// debio de haber ganado el envido
	for _, pkt := range pkts {
		if pkt.Message.Cod() == enco.TSumaPts {

			m, _ := pkt.Message.(enco.SumaPts)
			ok := util.All(
				m.Puntos == 4,
				m.Razon == enco.FaltaEnvidoGanado,
				m.Autor == "Roro",
				p.Ronda.Envite.Estado == DESHABILITADO,
			)

			// var t3 enco.Tipo3
			// json.Unmarshal(pkt.Message.Cont, &t3)

			// ok := util.All(
			// 	t3.Puntos == 4,
			// 	t3.Razon == enco.FaltaEnvidoGanado,
			// 	t3.Autor == "Roro",
			// 	p.Ronda.Envite.Estado == DESHABILITADO,
			// )

			util.Assert(ok, func() {
				t.Error("Alguna condicion no se cumplio")
			})

			break
		}
	}

}

func TestFixDebioHaberTerminado(t *testing.T) {
	data := `{"limiteEnvido":4,"cantJugadores":4,"puntuacion":20,"puntajes":{"azul":0,"rojo":3},"ronda":{"manoEnJuego":1,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":3,"pies":[0,0],"envite":{"estado":"deshabilitado","puntaje":3,"cantadoPor":"Renzo"},"truco":{"cantadoPor":"Alvaro","estado":"trucoQuerido"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"oro","valor":3},{"palo":"copa","valor":6}],"tiradas":[false,false,true],"ultimaTirada":2,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":true,"cartas":[{"palo":"basto","valor":4},{"palo":"espada","valor":10},{"palo":"copa","valor":1}],"tiradas":[false,true,true],"ultimaTirada":1,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":true,"cartas":[{"palo":"oro","valor":1},{"palo":"espada","valor":12},{"palo":"basto","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":5},{"palo":"oro","valor":6},{"palo":"espada","valor":7}],"tiradas":[true,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":1},"manos":[{"resultado":"ganoRojo","ganador":"Roro","cartasTiradas":[{"jugador":"Alvaro","carta":{"palo":"copa","valor":6}},{"jugador":"Roro","carta":{"palo":"copa","valor":1}},{"jugador":"Renzo","carta":{"palo":"copa","valor":5}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas":[{"jugador":"Roro","carta":{"palo":"espada","valor":10}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, err := Parse(data, true)
	if err != nil {
		t.Error(err)
	}
	t.Log(Renderizar(p))

	pkts, _ := p.Cmd("Alvaro re-truco")
	ok := util.All(
		p.Ronda.Truco.Estado == RETRUCO,
		p.Ronda.Truco.CantadoPor == "Alvaro",
		enco.Contains(pkts, enco.TGritarReTruco),
	)
	util.Assert(ok, func() {
		t.Error("Deberia poder dejarle gritar re-truco")
	})

	pkts, _ = p.Cmd("Renzo no-Quiero")
	// debe empezar una nueva ronda o partida finalizada
	ok = util.All(
		enco.Contains(pkts, enco.TNuevaRonda) || enco.Contains(pkts, enco.TByeBye),
	)
	util.Assert(ok, func() {
		t.Error("Debio de haber empezado una nueva ronda o haber terminado la partida")
	})

	pkts, _ = p.Cmd("Renzo vale-4")
	util.Assert(enco.Contains(pkts, enco.TError), func() {
		t.Error("No deberia dejarle gritar vale-4")
	})

}

func TestFixNoCantarPuntajeFlorCuandoNoEsNecesario(t *testing.T) {
	// tanto alvaro como adolfo tienen flor
	// ninguno de los 2 deberia decir al resto cuanto tienen porque
	// son de el mismo equipo
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"sinCantar":["Alvaro","Adolfo"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"copa","valor":4},{"palo":"copa","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"basto","valor":3},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"espada","valor":6},{"palo":"espada","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":5},{"palo":"espada","valor":10},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":12},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p, _ := Parse(partidaJSON, true)

	p.Ronda.Turno = 1
	p.Cmd("Roro truco")

	ok := p.Ronda.Truco.Estado == TRUCO
	if !ok {
		t.Error("Roro deberia ser capaz de gritar truco")
	}

	t.Log(Renderizar(p))

	ptsAntes := p.Puntajes[Azul]
	pkts, _ := p.Cmd("Alvaro flor")

	ok = !enco.Contains(pkts, enco.TDiceTengo)
	if !ok {
		t.Error("Alvaro no deberia decir cuanto tiene de flor ya que es el unico")
	}

	ok = ptsAntes+6 == p.Puntajes[Azul]
	if !ok {
		t.Error("Debio de haber ganado +6 pts por las 2 flores")
	}

	// ejemplo con 1 solo jugador; tampoco deberia
	partidaJSON = `{"puntuacion":20,"puntajes":{"azul":10,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":1},"elMano":1,"turno":1,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"sinCantar":["Alvaro"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":11},{"palo":"copa","valor":7},{"palo":"copa","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"copa","valor":2},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":11},{"palo":"basto","valor":4},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":true,"cartas":[{"palo":"basto","valor":6},{"palo":"oro","valor":4},{"palo":"oro","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":1},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p, _ = Parse(partidaJSON, true)

	t.Log(Renderizar(p))

	ptsAntes = p.Puntajes[Azul]
	pkts, _ = p.Cmd("Alvaro flor")

	ok = !enco.Contains(pkts, enco.TDiceTengo)
	if !ok {
		t.Error("Alvaro no deberia decir cuanto tiene de flor ya que es el unico")
	}

	ok = ptsAntes+3 == p.Puntajes[Azul]
	if !ok {
		t.Log("mal!")
	}

	ok = ptsAntes+3 == p.Puntajes[Azul]
	if !ok {
		t.Error("Debio de haber ganado +3 pts por la flor")
	}

}

func TestFixDecirSonBuenasDesdeUltratumba(t *testing.T) {
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true)
	p.Ronda.SetMuestra(Carta{Palo: Espada, Valor: 12})
	p.Ronda.ManoEnJuego = Primera
	p.Ronda.ElMano = 0 // alvaro
	p.Ronda.Turno = 0  // alvaro
	p.Ronda.SetManojos(
		[]Manojo{
			{
				Cartas: [3]*Carta{ // cartas de Alvaro
					{Palo: Basto, Valor: 5},
					{Palo: Oro, Valor: 7},
					{Palo: Espada, Valor: 2},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas Roro
					{Palo: Basto, Valor: 4},
					{Palo: Oro, Valor: 1},
					{Palo: Basto, Valor: 10},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Adolfo
					{Palo: Oro, Valor: 1},
					{Palo: Basto, Valor: 10},
					{Palo: Espada, Valor: 10},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Renzo
					{Palo: Copa, Valor: 4},
					{Palo: Copa, Valor: 12},
					{Palo: Basto, Valor: 6},
				},
			},
		},
	)

	p.Cmd("Alvaro 2 espada")
	p.Cmd("Roro 4 basto")
	p.Cmd("Roro mazo")
	p.Cmd("Adolfo falta-envido")

	t.Log(Renderizar(p))

	pkts, _ := p.Cmd("Renzo quiero")
	for _, pkt := range pkts {
		if pkt.Message.Cod() == enco.TDiceSonBuenas {
			m, _ := pkt.Message.(enco.DiceSonBuenas)
			// var autor string
			// json.Unmarshal(pkt.Message.Cont, &autor)
			if m == "Roro" {
				t.Error("Roro no deberia poder hablar dedse ultratumba")
				break
			}
		}
	}

}

func TestFixLoadJsonCartasTiradasAutorNil(t *testing.T) {
	data := `{"limiteEnvido":4,"cantJugadores":4,"puntuacion":20,"puntajes":{"azul":6,"rojo":6},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":2},"elMano":2,"turno":1,"pies":[0,0],"envite":{"estado":"deshabilitado","puntaje":3,"cantadoPor":"Adolfo"},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":true,"cartas":[{"palo":"oro","valor":12},{"palo":"oro","valor":1},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"copa","valor":2},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":7},{"palo":"oro","valor":10},{"palo":"espada","valor":3}],"tiradas":[false,true,false],"ultimaTirada":1,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":4}],"tiradas":[true,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":10},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas":[{"jugador":"Adolfo","carta":{"palo":"oro","valor":10}},{"jugador":"Renzo","carta":{"palo":"basto","valor":6}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, _ := Parse(data, true)

	ok := util.All(
		len(p.Ronda.Manos[0].CartasTiradas) == 2,
		p.Ronda.Manos[0].CartasTiradas[0].Jugador == "Adolfo",
		p.Ronda.Manos[0].CartasTiradas[1].Jugador == "Renzo",
	)
	util.Assert(ok, func() {
		t.Error("El mapeo de las cartas tiradas con sus autores no coincide")
	})
}

func TestFixPorQueDejaARoroTirarCarta(t *testing.T) {
	data := `{"limiteEnvido":4,"cantJugadores":4,"puntuacion":20,"puntajes":{"azul":6,"rojo":6},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":2},"elMano":2,"turno":1,"pies":[0,0],"envite":{"estado":"deshabilitado","puntaje":3,"cantadoPor":"Adolfo"},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":true,"cartas":[{"palo":"oro","valor":12},{"palo":"oro","valor":1},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"copa","valor":2},{"palo":"espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":7},{"palo":"oro","valor":10},{"palo":"espada","valor":3}],"tiradas":[false,true,false],"ultimaTirada":1,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":4}],"tiradas":[true,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"basto","valor":10},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas":[{"jugador":"Adolfo","carta":{"palo":"oro","valor":10}},{"jugador":"Renzo","carta":{"palo":"basto","valor":6}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, _ := Parse(data, true)

	t.Log(Renderizar(p))

	pkts, _ := p.Cmd("Roro truco")
	util.Assert(enco.Contains(pkts, enco.TGritarTruco), func() {
		t.Error("Deberia dejarlo poder gritar truco")
	})

	pkts, _ = p.Cmd("Roro 1 espada")
	util.Assert(enco.Contains(pkts, enco.TSigTurnoPosMano), func() {
		t.Error("Deberia pasarle el turno y empezar una mano nueva")
	})
	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Roro", func() {
		t.Error("Deberia mantener el turno porque gano la primera mano")
	})
	util.Assert(p.Ronda.ManoEnJuego == Segunda, func() {
		t.Error("Deberiamos estar jugando la segunda mano")
	})

	pkts, _ = p.Cmd("Renzo 3 copa")
	util.Assert(enco.Contains(pkts, enco.TError), func() {
		t.Error("No deberia dejar a renzo tirar carta si antes debe responder Adolfo y ni siquiera es su turno!")
	})

	pkts, _ = p.Cmd("Roro 2 copa")
	util.Assert(enco.Contains(pkts, enco.TSigTurno), func() {
		t.Error("Deberia pasarle el turno")
	})
}

func TestFixRespuestaNadaQueVer(t *testing.T) {
	data := `{"limiteEnvido":4,"cantJugadores":4,"puntuacion":20,"puntajes":{"azul":15,"rojo":6},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":1,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":""},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":7},{"palo":"oro","valor":2},{"palo":"copa","valor":1}],"tiradas":[false,true,false],"ultimaTirada":1,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":4},{"palo":"basto","valor":11},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"espada","valor":10},{"palo":"oro","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":7},{"palo":"copa","valor":2},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":7},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas":[{"jugador":"Alvaro","carta":{"palo":"oro","valor":2}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, _ := Parse(data, true)

	t.Log(Renderizar(p))

	pkts, _ := p.Cmd("Renzo mazo")
	util.Assert(enco.Contains(pkts, enco.TMazo), func() {
		t.Error("Renzo debieria ser capaz de irse al mazo")
	})

	pkts, _ = p.Cmd("Alvaro mazo")
	util.Assert(enco.Contains(pkts, enco.TMazo), func() {
		t.Error("Alvaro debieria ser capaz de irse al mazo")
	})

	pkts, _ = p.Cmd("Roro truco")
	util.Assert(enco.Contains(pkts, enco.TGritarTruco), func() {
		t.Error("Roro debieria ser capaz de gritar truco")
	})

	pkts, _ = p.Cmd("Alvaro real-envido")
	util.Assert(!enco.Contains(pkts, enco.TElEnvidoEstaPrimero), func() {
		t.Error("Alvaro no deberia ser capaz de tocar envido porque ya se fue al mazo")
	})
}

func TestFixError39(t *testing.T) {
	data := `{"limiteEnvido":4,"cantJugadores":2,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"envido","puntaje":2,"cantadoPor":"Alvaro"},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"oro","valor":5},{"palo":"espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":1},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, _ := Parse(data, true)

	ok := util.All(
		p.Ronda.Envite.Estado == ENVIDO,
		p.Ronda.Envite.Puntaje == 2,
		p.Ronda.Envite.CantadoPor == "Alvaro",
	)
	util.Assert(ok, func() {
		t.Error("Alavaro deberia ser quien canto el 1er envite")
	})

	_, _ = p.Cmd("Roro envido")
	ok = util.All(
		p.Ronda.Envite.Estado == ENVIDO,
		p.Ronda.Envite.Puntaje == 4,
		p.Ronda.Envite.CantadoPor == "Roro",
	)
	util.Assert(ok, func() {
		t.Error("Roro deberia ser quien canto el 2do envite")
	})

	data = `{"limiteEnvido":4,"cantJugadores":2,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"envido","puntaje":2,"cantadoPor":"Alvaro"},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"oro","valor":5},{"palo":"espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":1},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, _ = Parse(data, true)
	// p.FromJSON([]byte(data))

	ok = util.All(
		p.Ronda.Envite.Estado == ENVIDO,
		p.Ronda.Envite.Puntaje == 2,
		p.Ronda.Envite.CantadoPor == "Alvaro",
	)
	util.Assert(ok, func() {
		t.Error("Debimos haber vuelto al inicio: Alavaro deberia ser quien canto el 1er envite")
	})
}

func TestFixBadParsing(t *testing.T) {
	data := `{"puntuacion":20,"puntajes":{"azul":3,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"deshabilitado","puntaje":3,"cantadoPor":"Alvaro","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"copa","valor":5},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":2},{"palo":"espada","valor":6},{"palo":"copa","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":10},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p, _ := Parse(data, true)

	ok := util.All(
		p.Ronda.Envite.Estado == DESHABILITADO,
		len(p.Ronda.Envite.SinCantar) == 0,
	)
	util.Assert(ok, func() {
		t.Error("Alavaro ya canto la flor!")
	})

}

func TestFixRazonErronea(t *testing.T) {
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"sinCantar":["Alvaro"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"copa","valor":5},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":2},{"palo":"espada","valor":6},{"palo":"copa","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":10},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p, err := Parse(partidaJSON, true)
	if err != nil {
		t.Error(err)
	}

	p.Cmd("Alvaro flor")

	pkts, _ := p.Cmd("Roro mazo")

	/*
		luego deberia obtener:
			1. SumaPts con razon ?
			2. una y solo una RondaGanada con "valor" SeFueronAlMazo
	*/

	countMsgRondaGanada := 0

	for _, pkt := range pkts {

		// obtengo el contenido del mensaje
		// var cont map[string]json.RawMessage
		// json.Unmarshal(pkt.Message.Cont, &cont)

		if pkt.Message.Cod() == enco.TRondaGanada {
			countMsgRondaGanada++

			m, _ := pkt.Message.(enco.RondaGanada)
			// var r string
			// json.Unmarshal(cont["razon"], &r)

			ok := m.Razon != enco.EnvidoGanado
			if !ok {
				t.Error(`la razon por que ganan la ronda no deberia ser "por el envido"`)
			}
		}

		if pkt.Message.Cod() == enco.TSumaPts {
			m, _ := pkt.Message.(enco.SumaPts)
			// var r string
			// json.Unmarshal(cont["razon"], &r)

			razon := m.Razon
			ok := razon == enco.TrucoQuerido || razon == enco.TrucoNoQuerido
			if !ok {
				t.Error("no deberia ser la razon")
			}
		}

	}

	ok := countMsgRondaGanada == 1
	if !ok {
		t.Error(`Deberia retornaar solo 1 msg de tipo "RondaGanada"`)
	}

	// lo mismo pero con truco -> no-quiero
	p, _ = Parse(partidaJSON, true)
	p.Cmd("Alvaro flor")
	p.Cmd("Alvaro truco")
	pkts, _ = p.Cmd("Roro no-quiero")

	countMsgRondaGanada = 0

	for _, pkt := range pkts {

		// m, _ := pkt.Message.(enco.SumaPts)

		// obtengo el contenido del mensaje
		// var cont map[string]json.RawMessage
		// json.Unmarshal(pkt.Message.Cont, &cont)

		if pkt.Message.Cod() == enco.TRondaGanada {
			countMsgRondaGanada++

			m, _ := pkt.Message.(enco.RondaGanada)
			// var r string
			// json.Unmarshal(cont["razon"], &r)

			ok := m.Razon != enco.EnvidoGanado
			if !ok {
				t.Error(`la razon por que ganan la ronda no deberia ser "por el envido"`)
			}
		}

		if pkt.Message.Cod() == enco.TSumaPts {

			m, _ := pkt.Message.(enco.SumaPts)
			// var r string
			// json.Unmarshal(cont["razon"], &r)

			ok := m.Razon == enco.TrucoQuerido || m.Razon == enco.TrucoNoQuerido
			if !ok {
				t.Error("no deberia ser la razon")
			}
		}

	}

	ok = countMsgRondaGanada == 1
	if !ok {
		t.Error(`Deberia retornaar solo 1 msg de tipo "RondaGanada"`)
	}

}

func TestFixGanadorErroneo(t *testing.T) {
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":19,"rojo":19},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":null},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"oro","valor":5},{"palo":"espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":1},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, err := Parse(partidaJSON, true)
	if err != nil {
		t.Error(err)
	}

	p.Cmd("Alvaro 6 copa")
	pkts, _ := p.Cmd("Alvaro mazo")

	ok := p.Ronda.CantJugadoresEnJuego[Azul] > 0
	if !ok {
		t.Error(`la razon por que ganan la ronda deberia ser "porque se fueron al mazo"`)
	}

	countMsgRondaGanada := 0

	for _, pkt := range pkts {

		// obtengo el contenido del mensaje
		// var cont map[string]json.RawMessage
		// json.Unmarshal(pkt.Message.Cont, &cont)

		if pkt.Message.Cod() == enco.TRondaGanada {
			countMsgRondaGanada++
		}
	}

	ok = countMsgRondaGanada == 0
	if !ok {
		t.Error(`No deberia retornar "RondaGanada"`)
	}

}

func TestFixCodificacionCarta(t *testing.T) {
	// el json del contenido del paquete "cartaTirada"
	// debe decir palo "copa" y no "6"

	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":19,"rojo":19},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":null},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"oro","valor":5},{"palo":"espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":1},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, err := Parse(partidaJSON, true)
	if err != nil {
		t.Error(err)
	}

	pkts, _ := p.Cmd("Alvaro 6 copa")

	for _, pkt := range pkts {
		if pkt.Message.Cod() == enco.TTirarCarta {
			m, _ := pkt.Message.(enco.TirarCarta)
			// var t2 enco.Tipo2
			// json.Unmarshal(pkt.Message.Cont, &t2)

			ok := util.All(
				m.Palo == Copa.String(),
				m.Valor == 6,
				m.Autor == "Alvaro",
			)

			util.Assert(ok, func() {
				t.Error("Alguna condicion no se cumplio")
			})

			break
		}
	}

}

func TestCreadorDeEscenario(t *testing.T) {
	// por json
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":3,"rojo":0},"ronda":{"manoEnJuego":2,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":1,"envite":{"estado":"deshabilitado","puntaje":3,"cantadoPor":"Alvaro","sinCantar":[]},"truco":{"cantadoPor":"Roro","estado":"vale4Querido"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"copa","valor":5},{"palo":"copa","valor":7}],"tiradas":[false,true,true],"ultimaTirada":1,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":2},{"palo":"espada","valor":6},{"palo":"copa","valor":11}],"tiradas":[true,true,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}}],"muestra":{"palo":"oro","valor":10},"manos":[{"resultado":"ganoAzul","ganador":"Alvaro","cartasTiradas":[{"jugador":"Alvaro","carta":{"palo":"copa","valor":7}},{"jugador":"Roro","carta":{"palo":"espada","valor":6}}]},{"resultado":"ganoRojo","ganador":"Roro","cartasTiradas":[{"jugador":"Alvaro","carta":{"palo":"copa","valor":5}},{"jugador":"Roro","carta":{"palo":"oro","valor":2}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	q, err := Parse(partidaJSON, true)
	if err != nil {
		t.Error(err)
	}

	t.Log(Renderizar(q))

	qjson, _ := q.MarshalJSON()
	t.Log(string(qjson))

	for _, chi := range Chis(q) {
		for _, a := range chi {
			t.Log(a.String())
		}
	}

	// random
	p, err := NuevaPartida(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true)
	if err != nil {
		t.Error(err)
	}

	// ENVITE
	// p.Cmd("Alvaro flor")

	// TRUCO
	// p.Cmd("Alvaro truco")
	// p.Cmd("Roro re-truco")
	// p.Cmd("Alvaro vale-4")
	// p.Cmd("Roro quiero")

	// PRIMERA mano
	// p.Cmd("Alvaro 6 copa")
	// p.Cmd("Roro 3 copa")

	// SEGUNDA mano
	// p.Cmd("Roro 5 oro")
	// p.Cmd("Alvaro 3 oro")

	// TERCERA mano
	// p.Cmd("Alvaro 2 copa")

	// p.Cmd("Roro mazo") // <- hace que empiece una ronda nueva

	estadoEnviteOK := p.Ronda.Envite.Estado == NOCANTADOAUN
	estadoTrucoOK := p.Ronda.Truco.Estado == NOGRITADOAUN
	tieneLaPalabraOK := p.Ronda.Truco.CantadoPor == ""
	ultimaManoEnJuegoOK := p.Ronda.ManoEnJuego == Primera

	ok := util.All(
		estadoEnviteOK,
		estadoTrucoOK,
		tieneLaPalabraOK,
		ultimaManoEnJuegoOK,
	)

	if !ok {
		t.Error("algo salio mal")
	}

	json, _ := p.MarshalJSON()

	// opcional: renderizado
	t.Log(Renderizar(p))

	// json
	t.Log(string(json))

}

func random_action(aa []A) (jix, aix int) {

	type Tuple struct {
		jix int
		aix int
	}

	// hago un flatten del vector aa
	n := len(aa) * len(aa[0])
	flatten := make([]Tuple, 0, n)

	// i := 0
	for jix, aaj := range aa {
		for aix, a := range aaj {
			if a {
				flatten = append(flatten, Tuple{jix, aix})
				// flatten[i] = Tuple{jix, aix}
				// i++
			}
		}
	}

	// elijo a un (jugador,jugada) al azar
	rfix := rand.Intn(len(flatten))

	return flatten[rfix].jix, flatten[rfix].aix
}

func TestRandomWalk_AA(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 10000; i++ {
		// p, _ := NuevaPartida(A20, []string{"Alvaro"}, []string{"Roro"})
		p, _ := NuevaPartida(A20, []string{"Alvaro", "Andres"}, []string{"Roro", "Richard"}, 4, true)

		// partida
		// bs, _ := p.MarshalJSON()

		for {
			aa := GetAA(p)

			// elijo a un jugador al azar
			rjix, raix := random_action(aa)

			// v1
			pkts2 := aa[rjix].ToJugada(p, rjix, raix).Hacer(p)
			if p.Terminada() {
				byePkts2 := p.byeBye()
				pkts2 = append(pkts2, byePkts2...)
			}

			// v2
			// s := ActionToString(aa[r.jix], r.aix, r.jix, p)
			// cmd := fmt.Sprintf("%s %s", p.Ronda.Manojos[r.jix].Jugador.ID, s)
			// pkts, _ := p.Cmd(cmd)

			if IsDone(pkts2) {
				break
			}

			util.Assert(!enco.Contains(pkts2, enco.TError), func() {
				t.Error("NO PUEDE HABER JUGADAS INVALIDAS!")
			})
		}

	}
}

func TestRandomWalk_Chi(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 10000; i++ {
		// p, _ := NuevaPartida(A20, []string{"Alvaro"}, []string{"Roro"})
		p, _ := NuevaPartida(A20, []string{"Alvaro", "Andres"}, []string{"Roro", "Richard"}, 4, true)

		for {
			chis := Chis(p)

			// elijo a un jugador al azar
			rmix, raix := Random_action_chis(chis)

			// v1
			pkts2 := chis[rmix][raix].Hacer(p)
			if p.Terminada() {
				pktsBye2 := p.byeBye()
				pkts2 = append(pkts2, pktsBye2...)
			}

			// v2
			// s := ActionToString(aa[r.jix], r.aix, r.jix, p)
			// cmd := fmt.Sprintf("%s %s", p.Ronda.Manojos[r.jix].Jugador.ID, s)
			// pkts, _ := p.Cmd(cmd)

			if IsDone(pkts2) {
				break
			}

			util.Assert(!enco.Contains(pkts2, enco.TError), func() {
				t.Error("NO PUEDE HABER JUGADAS INVALIDAS!")
			})
		}

	}
}

func TestNil(t *testing.T) {
	p, _ := Parse(`{"puntuacion":20,"puntajes":{"azul":5,"rojo":6},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":2,"turno":2,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":null},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"basto","valor":6},{"palo":"espada","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"espada","valor":11},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"espada","valor":10},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":1},{"palo":"basto","valor":11},{"palo":"espada","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"mixs":{"Adolfo":2,"Alvaro":0,"Renzo":3,"Roro":1},"muestra":{"palo":"copa","valor":5},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`, true)
	t.Log(Renderizar(p))

	p.Cmd("Adolfo falta-envido")
	p.Cmd("Alvaro mazo")
	p.Cmd("Renzo no-Quiero")
	p.Cmd("Adolfo 3 oro")
	p.Cmd("Renzo 3 espada")
	p.Cmd("Roro truco")
	p.Cmd("Adolfo quiero")
	p.Cmd("Adolfo mazo")

	t.Log(Renderizar(p))
}

func TestFixNoPasaDeTurno(t *testing.T) {
	// 1. Deberia pasar el turno (de Alvaro a Roro)
	// 2. Deberia sumar puntos
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Andres"}, []string{"Roro", "Richard"}, 4, true)
	p.Ronda.SetMuestra(Carta{Palo: Espada, Valor: 1})
	p.Puntajes[Rojo] = 0
	p.Puntajes[Azul] = 0
	p.Ronda.ManoEnJuego = Primera
	p.Ronda.ElMano = 0 // Richard
	p.Ronda.Turno = 0  // Richard
	p.Ronda.SetManojos(
		[]Manojo{
			{
				Cartas: [3]*Carta{ // cartas de Alvaro
					{Palo: Espada, Valor: 4},
					{Palo: Basto, Valor: 11},
					{Palo: Espada, Valor: 3},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas Roro
					{Palo: Espada, Valor: 12},
					{Palo: Oro, Valor: 1},
					{Palo: Copa, Valor: 3},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Andres
					{Palo: Basto, Valor: 5},
					{Palo: Espada, Valor: 5},
					{Palo: Copa, Valor: 12},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Richard
					{Palo: Oro, Valor: 2},
					{Palo: Basto, Valor: 12},
					{Palo: Oro, Valor: 7},
				},
			},
		},
	)

	p.Cmd("Alvaro envido")
	p.Cmd("Richard real-envido")
	p.Cmd("Andres mazo")
	p.Cmd("Roro mazo")
	p.Cmd("Alvaro falta-envido")
	pkts, _ := p.Cmd("Richard quiero")

	util.Assert(enco.Contains(pkts, enco.TSumaPts), func() {
		t.Error("No debio de haberle dejado tocar envido")
	})

	pkts, _ = p.Cmd("Alvaro mazo")

	util.Assert(!enco.Contains(pkts, enco.TError), func() {
		t.Error("No deberia ocurrir errores")
	})

	util.Assert(p.Puntajes[Rojo] > 0, func() {
		t.Error("El puntaje del equipo rojo deberia ser mayor a cero")
	})

	util.Assert(p.Ronda.Turno > 0, func() {
		t.Error("El turno ya no lo deberia de tener Alvaro")
	})
}

func TestFixEnvidoHabilitado(t *testing.T) {
	// Deberia dejar cantar envido luego de que tiro carta? NO.
	// a menos que uno de mi mismo equipo tenga el turno (?)
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Andres"}, []string{"Roro", "Richard"}, 4, true)
	p.Ronda.SetMuestra(Carta{Palo: Espada, Valor: 1})
	p.Puntajes[Rojo] = 0
	p.Puntajes[Azul] = 0
	p.Ronda.ManoEnJuego = Primera
	p.Ronda.ElMano = 0 // Richard
	p.Ronda.Turno = 0  // Richard
	p.Ronda.SetManojos(
		[]Manojo{
			{
				Cartas: [3]*Carta{ // cartas de Alvaro
					{Palo: Espada, Valor: 4},
					{Palo: Basto, Valor: 11},
					{Palo: Espada, Valor: 3},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas Roro
					{Palo: Espada, Valor: 12},
					{Palo: Oro, Valor: 1},
					{Palo: Copa, Valor: 3},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Andres
					{Palo: Basto, Valor: 5},
					{Palo: Espada, Valor: 5},
					{Palo: Copa, Valor: 12},
				},
			},
			{
				Cartas: [3]*Carta{ // cartas de Richard
					{Palo: Oro, Valor: 2},
					{Palo: Basto, Valor: 12},
					{Palo: Oro, Valor: 7},
				},
			},
		},
	)

	p.Cmd("Alvaro 4 espada")
	pkts, _ := p.Cmd("Alvaro envido")

	util.Assert(enco.Contains(pkts, enco.TError), func() {
		t.Error("No deberia dejarlo arrancar un envido despues de tirar carta")
	})

}

func TestShiftPartida(t *testing.T) {
	// dada una partida, la shiftea
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alvaro"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":1},{"palo":"oro","valor":4},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":10},{"palo":"copa","valor":1},{"palo":"oro","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":7},{"palo":"basto","valor":7},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":7},{"palo":"basto","valor":3},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"rojo"}}],"mixs":{"Adolfo":2,"Alvaro":0,"Renzo":3,"Roro":1},"muestra":{"palo":"espada","valor":1},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, _ := Parse(partidaJSON, true)

	t.Log(Renderizar(p))

	// data, _ := p.MarshalJSON()
	t.Log("------------------------shift:---------------------------------")

	// la original
	data_manojos, _ := json.Marshal(p.Ronda.Manojos)

	var manojos []Manojo
	json.Unmarshal(data_manojos, &manojos)
	// shift
	manojos = append(manojos[1:], manojos[0])
	p.Ronda.SetManojos(manojos)

	t.Log(Renderizar(p))

	t.Log("===============================================================")
}

func TestFipPartida(t *testing.T) {
	// dada una partida aleatoria
	// computa un ida-vuelta
	p, _ := NuevaPartida(A20, []string{"Alice", "Ana"}, []string{"Bob", "Ben"}, 4, true)

	t.Log(Renderizar(p))

	// data, _ := p.MarshalJSON()
	t.Log("------------------------shift:---------------------------------")

	// la original
	data_muestra, _ := json.Marshal(p.Ronda.Muestra)
	data_manojos, _ := json.Marshal(p.Ronda.Manojos)

	// creo una partida nueva ahora con el orden invertido
	p, _ = NuevaPartida(A20, []string{"Bob", "Ben"}, []string{"Alice", "Ana"}, 4, true)

	// le seteo los manojos de antes
	var (
		manojos []Manojo
		muestra Carta
	)
	json.Unmarshal(data_muestra, &muestra)
	json.Unmarshal(data_manojos, &manojos)
	p.Ronda.SetManojos(manojos)
	p.Ronda.SetMuestra(muestra)

	t.Log(Renderizar(p))

	t.Log("===============================================================")
}

func TestFixTurnoDeAlguienQueSeFueAlMazo(t *testing.T) {
	// el json del contenido del paquete "cartaTirada"
	// debe decir palo "copa" y no "6"

	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":1},{"palo":"oro","valor":5},{"palo":"basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":4},{"palo":"oro","valor":7},{"palo":"espada","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Bob","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":1},{"palo":"oro","valor":6},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Ariana","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":3},{"palo":"copa","valor":6},{"palo":"copa","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Ben","equipo":"rojo"}}],"mixs":{"Alice":0,"Ariana":2,"Ben":3,"Bob":1},"muestra":{"palo":"basto","valor":12},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, err := Parse(partidaJSON, true)
	if err != nil {
		t.Error(err)
	}

	_, _ = p.Cmd("Alice 5 oro")
	_, _ = p.Cmd("Bob mazo")
	_, _ = p.Cmd("Ariana 1 basto")
	_, _ = p.Cmd("Ben envido")
	_, _ = p.Cmd("Ariana envido")
	_, _ = p.Cmd("Ben real-envido")
	_, _ = p.Cmd("Ariana mazo")
	_, _ = p.Cmd("Alice falta-envido")
	_, _ = p.Cmd("Ben no-quiero")
	_, _ = p.Cmd("Ben 3 espada")

	t.Log(Renderizar(p))

	if ok := p.Ronda.GetElTurno().Jugador.ID == "Alice"; !ok {
		t.Error("el turno no deberia ser ariana", p.Ronda.GetElTurno().Jugador.ID)
	}

}

func TestFixFlorColgada(t *testing.T) {
	// ben queda sin cantar la flor...
	// no puede decir nada, wtf... ?
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Alice","Bob","Ben"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":3},{"palo":"basto","valor":11},{"palo":"basto","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":5}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Bob","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":1},{"palo":"oro","valor":2},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Ariana","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":11},{"palo":"copa","valor":10},{"palo":"espada","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Ben","equipo":"rojo"}}],"mixs":{"Alice":0,"Ariana":2,"Ben":3,"Bob":1},"muestra":{"palo":"copa","valor":4},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, err := Parse(partidaJSON, true)
	if err != nil {
		t.Error(err)
	}

	t.Log(Renderizar(p))

	_, _ = p.Cmd("Alice flor")
	_, _ = p.Cmd("Bob contra-flor-al-resto")
	_, _ = p.Cmd("Alice mazo")

	t.Log(Renderizar(p))

	pkts, err := p.Cmd("Ben flor")
	if ok := err == nil && enco.Contains(pkts, enco.TSumaPts); !ok {
		t.Error("Debio de haber sumado los pts de las flores")
	}

	t.Log(Renderizar(p))

	if ok := p.Ronda.Envite.Estado == DESHABILITADO; !ok {
		t.Error("ya cantaron todos; el envido deberia quedar deshabilitado")
	}
}

func TestFixMazoPoints(t *testing.T) {
	// el problema: aveces dice "me voy al mazo" y le suma los puntos al equipo
	// contrario
	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"copa","valor":10},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"ESL-0","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":11},{"palo":"basto","valor":3},{"palo":"espada","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"jp","equipo":"rojo"}}],"mixs":{"ESL-0":0,"jp":1},"muestra":{"palo":"espada","valor":10},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, err := Parse(partidaJSON, true)
	if err != nil {
		t.Error(err)
	}

	t.Log("estado inicial")
	t.Log(Renderizar(p))

	p.Cmd("ESL-0 real-envido")
	p.Cmd("jp no-quiero")

	// muestra 10 espada

	// 1era mano gana ESL
	p.Cmd("ESL-0 11 oro")
	p.Cmd("jp 11 copa")

	// 2da mano queda sin definirse PERO NOTAR QUE es ESL-0/Azul el que tira la carta
	p.Cmd("ESL-0 10 copa")
	p.Cmd("jp truco")

	t.Log("el momento justo antes de que ESL se vaya al mazo")
	t.Log(Renderizar(p))

	p.Cmd("ESL-0 mazo")

	t.Log("estado final")
	t.Log(Renderizar(p))

	equipo_jp := p.Manojo("jp").Jugador.Equipo
	if ok := p.Puntajes[equipo_jp] > 0; !ok {
		t.Error("los puntos los deberia ganar jp")
	}
}

func TestFixMazoPoints_2(t *testing.T) {
	// igual que el test anterior pero otro set

	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":17,"rojo":11},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":10},{"palo":"copa","valor":6},{"palo":"espada","valor":12}],"tiradas":[false,false,false],"ultimaTirada":1,"jugador":{"id":"ESL-0","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":12},{"palo":"basto","valor":11},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":2,"jugador":{"id":"j","equipo":"rojo"}}],"mixs":{"ESL-0":0,"j":1},"muestra":{"palo":"oro","valor":10},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []},{"resultado":"ganoRojo","ganador":"","cartasTiradas": []}]}}`
	p, err := Parse(partidaJSON, true)
	if err != nil {
		t.Error(err)
	}

	t.Log("estado inicial")
	t.Log(Renderizar(p))

	p.Cmd("ESL-0 real-envido")
	p.Cmd("j quiero")

	p.Cmd("ESL-0 12 espada")
	p.Cmd("j 10 espada")

	p.Cmd("ESL-0 10 basto")
	p.Cmd("j truco")

	t.Log("el momento justo antes de que ESL se vaya al mazo")
	t.Log(Renderizar(p))

	p.Cmd("ESL-0 mazo")

	t.Log("estado final")
	t.Log(Renderizar(p))

	equipo_jp := p.Manojo("j").Jugador.Equipo
	if ok := p.Puntajes[equipo_jp] > 0; !ok {
		t.Error("los puntos los deberia ganar jp")
	}
}

func TestPardaSigTurno3(t *testing.T) {
	// igual que el anterior pero ahora todos los que empardaron se van al mazo
	// si va parda, el siguiente turno deberia ser del mano
	// o del mas cercano a este
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"}, 4, true)
	p.Ronda.SetMuestra(Carta{Palo: Espada, Valor: 1})
	p.Ronda.SetManojos(
		[]Manojo{
			{
				Cartas: [3]*Carta{ // envido: 26
					{Palo: Copa, Valor: 7},
					{Palo: Oro, Valor: 12},
					{Palo: Copa, Valor: 5},
				},
			},
			{
				Cartas: [3]*Carta{ // envido: 27
					{Palo: Copa, Valor: 3},
					{Palo: Copa, Valor: 4},
					{Palo: Oro, Valor: 4},
				},
			},
			{
				Cartas: [3]*Carta{ // envido: 28
					{Palo: Copa, Valor: 2},
					{Palo: Copa, Valor: 1},
					{Palo: Espada, Valor: 6}, // <-- parda primera mano
				},
			},
			{
				Cartas: [3]*Carta{ // envido: 25
					{Palo: Basto, Valor: 2},
					{Palo: Oro, Valor: 3},
					{Palo: Oro, Valor: 6}, // <-- parda primera mano
				},
			},
			{
				Cartas: [3]*Carta{ // envido: 33
					{Palo: Basto, Valor: 3},
					{Palo: Basto, Valor: 7},
					{Palo: Copa, Valor: 6}, // <-- parda primera mano
				},
			},
			{
				Cartas: [3]*Carta{ // envido: 20
					{Palo: Basto, Valor: 4},
					{Palo: Copa, Valor: 11},
					{Palo: Basto, Valor: 6}, // <-- parda primera mano
				},
			},
		},
	)

	pkgs, _ := p.Cmd("Alvaro 5 copa")

	{
		ok := false
		for _, pkg := range pkgs {
			if pkg.Message.Cod() == enco.TTirarCarta {
				ok = true
				break
			}
		}
		if !ok {
			t.Error("debio de haber tirado carta")
		}
	}

	// util.Assert(enco.Contains(enco.Collect(pkgs), enco.TTirarCarta), func() {
	// 	t.Error("debio de haber tirado carta")
	// })

	p.Cmd("Roro 4 copa")
	// los siguientes 4: todos tiran 6 -> resulta mano parda
	p.Cmd("Adolfo 6 espada")
	p.Cmd("Adolfo mazo")
	p.Cmd("Renzo 6 oro")
	p.Cmd("Renzo mazo")
	p.Cmd("Andres 6 copa")
	p.Cmd("Andres mazo")
	p.Cmd("Richard 4 basto")
	p.Cmd("Richard mazo")

	// ojo no imprime nada porque la aterior ya consumio el out
	// enco.Consume(out, func(pkt *enco.Packet2) {
	// 	t.Log(deco.Stringify(pkt, p.Partida))
	// })

	// enco.Consume(out, func(pkt *enco.Packet2) {
	// 	t.Log(deco.Stringify(pkt, p.Partida))
	// })

	t.Log(p)

	util.Assert(p.Ronda.Manos[0].Resultado == Empardada, func() {
		t.Error(`La mano debio ser parda`)
	})

	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Roro", func() {
		t.Error(`Deberia ser turno de Roro, debido a que es el mas cercano del mano y qu empardo`)
	})

	util.Assert(p.Ronda.Manojos[5].SeFueAlMazo == true, func() {
		t.Error(`deberia dejarlo irse al mazo`)
	})
}

func TestFormatoJSONPartida(t *testing.T) {
	// azul y rojo con minusculas
	p, err := NuevaPartida(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"}, 4, true)
	if err != nil {
		t.Error(err)
	}
	bs, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bs))
}

func TestBughunterErrorLoopPythonCli(t *testing.T) {
	// Se intenta resolver estas preguntas:
	// 1. ¿Por qué todas las acciones de JP son rechazadas por el simulador?
	// 2. ¿Por qué el otro jugador (Auron) tiene un Chi con 0 acciones?

	// a partir de la fusion de estos 2 json, se forma el json a nivel de simulador
	// {"puntuacion":20,"puntajes":{"azul":7,"rojo":15},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":1,"turno":1,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["juampi"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"oro","valor":1},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"juampi","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[null,null,null],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"auronPlay","equipo":"rojo"}}],"mixs":{"auronPlay":1,"juampi":0},"muestra":{"palo":"basto","valor":10},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":null},{"resultado":"indeterminado","ganador":"","cartasTiradas":null},{"resultado":"indeterminado","ganador":"","cartasTiradas":null}]},"limiteEnvido":4}
	// {"puntuacion":20,"puntajes":{"azul":7,"rojo":15},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":1,"turno":1,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["juampi"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[null,null,null],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"juampi","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"copa","valor":1},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"auronPlay","equipo":"rojo"}}],"mixs":{"auronPlay":1,"juampi":0},"muestra":{"palo":"basto","valor":10},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":null},{"resultado":"indeterminado","ganador":"","cartasTiradas":null},{"resultado":"indeterminado","ganador":"","cartasTiradas":null}]},"limiteEnvido":4}

	partidaJSON := `{"puntuacion":20,"puntajes":{"azul":7,"rojo":15},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":1,"turno":1,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["juampi"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"oro","valor":1},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"juampi","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"copa","valor":1},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"auronPlay","equipo":"rojo"}}],"mixs":{"auronPlay":1,"juampi":0},"muestra":{"palo":"basto","valor":10},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":null},{"resultado":"indeterminado","ganador":"","cartasTiradas":null},{"resultado":"indeterminado","ganador":"","cartasTiradas":null}]},"limiteEnvido":4}`
	p, err := Parse(partidaJSON, true)
	if err != nil {
		t.Error(err)
	}

	t.Log("estado inicial")
	t.Log(Renderizar(p))

	// estos fueron los mensajes de SALIDA
	// voy a recrear las ENTRADAS en base a estos

	// to:@juampi {"type":"msg","payload":{"cod":"GritarTruco","cont":"auronPlay"}}
	// to:@auronPlay {"type":"msg","payload":{"cod":"GritarTruco","cont":"auronPlay"}}
	p.Cmd("auronPlay truco")

	// to:@juampi {"type":"msg","payload":{"cod":"TirarCarta","cont":{"autor":"auronPlay","palo":"basto","valor":6}}}
	// to:@juampi {"type":"msg","payload":{"cod":"SigTurno","cont":0}}
	// to:@auronPlay {"type":"msg","payload":{"cod":"TirarCarta","cont":{"autor":"auronPlay","palo":"basto","valor":6}}}
	// to:@auronPlay {"type":"msg","payload":{"cod":"SigTurno","cont":0}}
	p.Cmd("auronPlay 6 basto")

	// duda: esto que viene a continuación, lo hace juampi o lo hace el simulador?

	// to:@juampi {"type":"msg","payload":{"cod":"CantarFlor","cont":"juampi"}}
	// to:@juampi {"type":"msg","payload":{"cod":"SumaPts","cont":{"autor":"juampi","razon":"LaUnicaFlor","puntos":3}}}
	// to:@auronPlay {"type":"msg","payload":{"cod":"CantarFlor","cont":"juampi"}}
	// to:@auronPlay {"type":"msg","payload":{"cod":"SumaPts","cont":{"autor":"juampi","razon":"LaUnicaFlor","puntos":3}}}

	// respuesta: veamos el estado del simulador
	util.Assert(p.Ronda.Envite.Estado == NOCANTADOAUN, func() {
		t.Error("El estado del envite debería ser `NOCANTADOAUN`")
	})

	util.Assert(p.Ronda.Truco.Estado == TRUCO, func() {
		t.Error("El estado del truco debería ser `TRUCO`")
	})

	util.Assert(p.Ronda.Truco.CantadoPor == "auronPlay", func() {
		t.Error("El autor del truco debería ser `auronPlay`")
	})

	// ok, lo hizo juampi
	p.Cmd("juampi flor")

	util.Assert(p.Ronda.Truco.Estado == TRUCO, func() {
		t.Error("El estado del truco debería ser `TRUCO`")
	})

	util.Assert(p.Ronda.Truco.CantadoPor == "auronPlay", func() {
		t.Error("El autor del truco debería ser `auronPlay`")
	})

	t.Log(Renderizar(p))

	chisJuampi := MetaChi(p, p.Manojo("juampi"), false)
	t.Logf("chi(juampi) ~ %v\n", chisJuampi)

	chisAuron := MetaChi(p, p.Manojo("auronPlay"), false)
	t.Logf("chi(auronPlay) ~ %v\n", chisAuron)
}
