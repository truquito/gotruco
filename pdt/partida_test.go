package pdt

import (
	"encoding/json"
	"testing"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/util"
)

func TestNoStruct(t *testing.T) {
	partidaJSON := `{"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"sinCantar":["Renzo","Richard"]},"truco":{"cantadoPor":"","estado":"truco"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":12},{"palo":"Oro","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Basto","valor":10},{"palo":"Oro","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":10},{"palo":"Basto","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Espada","valor":10},{"palo":"Basto","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Copa","valor":3},{"palo":"Espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":6},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null}]}}`
	p, err := Parse(partidaJSON)
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
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
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
	p.Cmd("Renzo 11 Basto")
	t.Log(Renderizar(p)) // retorna el render
}

func TestFixNoLePermiteTocarEnvido(t *testing.T) {
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
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

	util.Assert(enco.Contains(pkts, enco.Error), func() {
		t.Error("No debio de haberle dejado tocar envido")
	})
}

func TestFixNoLePermiteGritarTruco(t *testing.T) {
	data := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 0, "Rojo": 0}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 2, "Rojo": 2}, "elMano": 0, "turno": 0, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": "", "estado": "noCantado"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 3}, {"palo": "Oro", "valor": 7}, {"palo": "Copa", "valor": 7}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 10}, {"palo": "Basto", "valor": 5}, {"palo": "Oro", "valor": 6}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 2}, {"palo": "Oro", "valor": 3}, {"palo": "Espada", "valor": 7}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 1}, {"palo": "Espada", "valor": 6}, {"palo": "Espada", "valor": 10}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Copa", "valor": 3}, "manos": [{"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}]}} `
	p, err := Parse(data)
	if err != nil {
		t.Error(err)
	}

	pkts, _ := p.Cmd("Alvaro truco")

	util.Assert(enco.Contains(pkts, enco.GritarTruco), func() {
		t.Error("Deberia dejarlo gritar truco!")
	})
}

func TestFixLePermiteCFAR(t *testing.T) {
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
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
	util.Assert(enco.Contains(pkts, enco.Error), func() {
		t.Error("No deberia dejarlo cantar CFAR")
	})

	pkts, _ = p.Cmd("Renzo no-quiero")
	util.Assert(enco.Contains(pkts, enco.Error), func() {
		t.Error("No deberia dejarlo responder no-quiero")
	})
}

func TestFixNoLeDebePermitirTocarEnvidoSiGritoTruco(t *testing.T) {
	data := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 0, "Rojo": 0}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 2, "Rojo": 2}, "elMano": 0, "turno": 0, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": "", "estado": "noCantado"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 3}, {"palo": "Oro", "valor": 7}, {"palo": "Copa", "valor": 7}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 10}, {"palo": "Basto", "valor": 5}, {"palo": "Oro", "valor": 6}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 2}, {"palo": "Oro", "valor": 3}, {"palo": "Espada", "valor": 7}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 1}, {"palo": "Espada", "valor": 6}, {"palo": "Espada", "valor": 10}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Copa", "valor": 3}, "manos": [{"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}]}} `
	p, err := Parse(data)
	if err != nil {
		t.Error(err)
	}
	t.Log(Renderizar(p))

	pkts, _ := p.Cmd("Alvaro truco")
	util.Assert(enco.Contains(pkts, enco.GritarTruco), func() {
		t.Error("Deberia dejarlo gritar truco!")
	})

	pkts, _ = p.Cmd("Alvaro envido")
	ok := util.All(
		p.Ronda.Envite.Estado == NOCANTADOAUN,
		p.Ronda.Truco.Estado != NOCANTADO,
		enco.Contains(pkts, enco.Error),
	)
	util.Assert(ok, func() {
		t.Error("No deberia dejarlo tocar envido si fue el mismo el que toco truco!")
	})

	pkts, _ = p.Cmd("Alvaro real-envido")
	ok = util.All(
		p.Ronda.Envite.Estado == NOCANTADOAUN,
		p.Ronda.Truco.Estado != NOCANTADO,
		enco.Contains(pkts, enco.Error),
	)
	util.Assert(ok, func() {
		t.Error("No deberia dejarlo tocar real-envido si fue el mismo el que toco truco!")
	})

}

func TestFixNoLeDeberiaResponderDesdeUltratumba(t *testing.T) {
	data := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 5, "Rojo": 6}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 2, "Rojo": 2}, "elMano": 2, "turno": 2, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": "", "estado": "noCantado"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Basto", "valor": 2}, {"palo": "Basto", "valor": 6}, {"palo": "Espada", "valor": 5}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 6}, {"palo": "Espada", "valor": 11}, {"palo": "Copa", "valor": 12}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Oro", "valor": 3}, {"palo": "Espada", "valor": 10}, {"palo": "Oro", "valor": 10}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 1}, {"palo": "Basto", "valor": 11}, {"palo": "Espada", "valor": 3}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Copa", "valor": 5}, "manos": [{"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}]}}`
	p, err := Parse(data)
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
		diceSonBuenas := pkt.Message.Cod == string(enco.DiceSonBuenas)
		loDijoRenzo := string(pkt.Message.Cont) == "\"Renzo\""
		if diceSonBuenas && loDijoRenzo {
			t.Error("No deberia poder responder desde ultratumba")
		}
	}

	// debio de haber ganado el envido
	for _, pkt := range pkts {
		if pkt.Message.Cod == string(enco.SumaPts) {
			var t3 enco.Tipo3
			json.Unmarshal(pkt.Message.Cont, &t3)

			ok := util.All(
				t3.Puntos == 4,
				t3.Razon == enco.FaltaEnvidoGanado,
				t3.Autor == "Roro",
				p.Ronda.Envite.Estado == DESHABILITADO,
			)

			util.Assert(ok, func() {
				t.Error("Alguna condicion no se cumplio")
			})

			break
		}
	}

}

func TestFixDebioHaberTerminado(t *testing.T) {
	data := `{"Jugadores":[{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"},{"id":"Roro","nombre":"Roro","equipo":"Rojo"},{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"},{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}],"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":3},"ronda":{"manoEnJuego":1,"cantJugadoresEnJuego":{"Azul":1,"Rojo":1},"elMano":0,"turno":3,"pies":[0,0],"envite":{"estado":"deshabilitado","puntaje":3,"cantadoPor":"Renzo"},"truco":{"cantadoPor":"Alvaro","estado":"trucoQuerido"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Oro","valor":3},{"palo":"Copa","valor":6}],"tiradas":[false,false,true],"ultimaTirada":2,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":true,"cartas":[{"palo":"Basto","valor":4},{"palo":"Espada","valor":10},{"palo":"Copa","valor":1}],"tiradas":[false,true,true],"ultimaTirada":1,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":true,"cartas":[{"palo":"Oro","valor":1},{"palo":"Espada","valor":12},{"palo":"Basto","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Oro","valor":6},{"palo":"Espada","valor":7}],"tiradas":[true,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":1},"manos":[{"resultado":"ganoRojo","ganador":"Roro","cartasTiradas":[{"jugador":"Alvaro","carta":{"palo":"Copa","valor":6}},{"jugador":"Roro","carta":{"palo":"Copa","valor":1}},{"jugador":"Renzo","carta":{"palo":"Copa","valor":5}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas":[{"jugador":"Roro","carta":{"palo":"Espada","valor":10}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas":null}]}}`
	p, err := Parse(data)
	if err != nil {
		t.Error(err)
	}
	t.Log(Renderizar(p))

	pkts, _ := p.Cmd("Alvaro re-truco")
	ok := util.All(
		p.Ronda.Truco.Estado == RETRUCO,
		p.Ronda.Truco.CantadoPor == "Alvaro",
		enco.Contains(pkts, enco.GritarReTruco),
	)
	util.Assert(ok, func() {
		t.Error("Deberia poder dejarle gritar re-truco")
	})

	pkts, _ = p.Cmd("Renzo no-Quiero")
	// debe empezar una nueva ronda o partida finalizada
	ok = util.All(
		enco.Contains(pkts, enco.NuevaRonda) || enco.Contains(pkts, enco.ByeBye),
	)
	util.Assert(ok, func() {
		t.Error("Debio de haber empezado una nueva ronda o haber terminado la partida")
	})

	pkts, _ = p.Cmd("Renzo vale-4")
	util.Assert(enco.Contains(pkts, enco.Error), func() {
		t.Error("No deberia dejarle gritar vale-4")
	})

}

func TestFixNoCantarPuntajeFlorCuandoNoEsNecesario(t *testing.T) {
	// tanto alvaro como adolfo tienen flor
	// ninguno de los 2 deberia decir al resto cuanto tienen porque
	// son de el mismo equipo
	partidaJSON := `{"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":2},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"sinCantar":["Alvaro","Adolfo"]},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":4},{"palo":"Copa","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Basto","valor":3},{"palo":"Oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":6},{"palo":"Espada","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":5},{"palo":"Espada","valor":10},{"palo":"Copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":12},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null}]}}`
	p, _ := Parse(partidaJSON)

	p.Ronda.Turno = 1
	p.Cmd("Roro truco")

	ok := p.Ronda.Truco.Estado == TRUCO
	if !ok {
		t.Error("Roro deberia ser capaz de gritar truco")
	}

	t.Log(Renderizar(p))

	ptsAntes := p.Puntajes[Azul]
	pkts, _ := p.Cmd("Alvaro flor")

	ok = !enco.Contains(pkts, enco.DiceTengo)
	if !ok {
		t.Error("Alvaro no deberia decir cuanto tiene de flor ya que es el unico")
	}

	ok = ptsAntes+6 == p.Puntajes[Azul]
	if !ok {
		t.Error("Debio de haber ganado +6 pts por las 2 flores")
	}

	// ejemplo con 1 solo jugador; tampoco deberia
	partidaJSON = `{"puntuacion":20,"puntajes":{"Azul":10,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":2,"Rojo":1},"elMano":1,"turno":1,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"sinCantar":["Alvaro"]},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":11},{"palo":"Copa","valor":7},{"palo":"Copa","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":10},{"palo":"Copa","valor":2},{"palo":"Basto","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":11},{"palo":"Basto","valor":4},{"palo":"Oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":true,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":4},{"palo":"Oro","valor":3}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":1},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null}]}}`
	p, _ = Parse(partidaJSON)

	t.Log(Renderizar(p))

	ptsAntes = p.Puntajes[Azul]
	pkts, _ = p.Cmd("Alvaro flor")

	ok = !enco.Contains(pkts, enco.DiceTengo)
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
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
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
		if pkt.Message.Cod == string(enco.DiceSonBuenas) {
			var autor string
			json.Unmarshal(pkt.Message.Cont, &autor)
			if autor == "Roro" {
				t.Error("Roro no deberia poder hablar dedse ultratumba")
				break
			}
		}
	}

}

func TestFixLoadJsonCartasTiradasAutorNil(t *testing.T) {
	data := `{"Jugadores":[{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"},{"id":"Roro","nombre":"Roro","equipo":"Rojo"},{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"},{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}],"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":6,"Rojo":6},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":1,"Rojo":2},"elMano":2,"turno":1,"pies":[0,0],"envite":{"estado":"deshabilitado","puntaje":3,"cantadoPor":"Adolfo"},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":true,"cartas":[{"palo":"Oro","valor":12},{"palo":"Oro","valor":1},{"palo":"Espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Copa","valor":2},{"palo":"Espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":7},{"palo":"Oro","valor":10},{"palo":"Espada","valor":3}],"tiradas":[false,true,false],"ultimaTirada":1,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":3},{"palo":"Copa","valor":4}],"tiradas":[true,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":10},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas":[{"jugador":"Adolfo","carta":{"palo":"Oro","valor":10}},{"jugador":"Renzo","carta":{"palo":"Basto","valor":6}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas":null},{"resultado":"ganoRojo","ganador":"","cartasTiradas":null}]}}`
	p, _ := Parse(data)

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
	data := `{"Jugadores":[{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"},{"id":"Roro","nombre":"Roro","equipo":"Rojo"},{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"},{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}],"cantJugadores":4,"puntuacion":20,"puntajes":{"Azul":6,"Rojo":6},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":1,"Rojo":2},"elMano":2,"turno":1,"pies":[0,0],"envite":{"estado":"deshabilitado","puntaje":3,"cantadoPor":"Adolfo"},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":true,"cartas":[{"palo":"Oro","valor":12},{"palo":"Oro","valor":1},{"palo":"Espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Copa","valor":2},{"palo":"Espada","valor":1}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":7},{"palo":"Oro","valor":10},{"palo":"Espada","valor":3}],"tiradas":[false,true,false],"ultimaTirada":1,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Oro","valor":3},{"palo":"Copa","valor":4}],"tiradas":[true,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}}],"muestra":{"palo":"Basto","valor":10},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas":[{"jugador":"Adolfo","carta":{"palo":"Oro","valor":10}},{"jugador":"Renzo","carta":{"palo":"Basto","valor":6}}]},{"resultado":"ganoRojo","ganador":"","cartasTiradas":null},{"resultado":"ganoRojo","ganador":"","cartasTiradas":null}]}}`
	p, _ := Parse(data)

	t.Log(Renderizar(p))

	pkts, _ := p.Cmd("Roro truco")
	util.Assert(enco.Contains(pkts, enco.GritarTruco), func() {
		t.Error("Deberia dejarlo poder gritar truco")
	})

	pkts, _ = p.Cmd("Roro 1 Espada")
	util.Assert(enco.Contains(pkts, enco.SigTurnoPosMano), func() {
		t.Error("Deberia pasarle el turno y empezar una mano nueva")
	})
	util.Assert(p.Ronda.GetElTurno().Jugador.ID == "Roro", func() {
		t.Error("Deberia mantener el turno porque gano la primera mano")
	})
	util.Assert(p.Ronda.ManoEnJuego == Segunda, func() {
		t.Error("Deberiamos estar jugando la segunda mano")
	})

	pkts, _ = p.Cmd("Renzo 3 Copa")
	util.Assert(enco.Contains(pkts, enco.Error), func() {
		t.Error("No deberia dejar a renzo tirar carta si antes debe responder Adolfo y ni siquiera es su turno!")
	})

	pkts, _ = p.Cmd("Roro 2 copa")
	util.Assert(enco.Contains(pkts, enco.SigTurno), func() {
		t.Error("Deberia pasarle el turno")
	})
}

func TestFixRespuestaNadaQueVer(t *testing.T) {
	data := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 15, "Rojo": 6}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 2, "Rojo": 2}, "elMano": 0, "turno": 1, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": "", "estado": "noCantado"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 7}, {"palo": "Oro", "valor": 2}, {"palo": "Copa", "valor": 1}], "tiradas": [false,true,false], "ultimaTirada": 1, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 4}, {"palo": "Basto", "valor": 11}, {"palo": "Oro", "valor": 12}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Basto", "valor": 12}, {"palo": "Espada", "valor": 10}, {"palo": "Oro", "valor": 4}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Basto", "valor": 7}, {"palo": "Copa", "valor": 2}, {"palo": "Copa", "valor": 7}], "tiradas": [false,false,false], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Oro", "valor": 7}, "manos": [{"resultado": "ganoRojo", "ganador": "", "cartasTiradas": [{"palo": "Oro", "valor": 2}]}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": "", "cartasTiradas": null}]}}`
	p, _ := Parse(data)

	t.Log(Renderizar(p))

	pkts, _ := p.Cmd("Renzo mazo")
	util.Assert(enco.Contains(pkts, enco.Mazo), func() {
		t.Error("Renzo debieria ser capaz de irse al mazo")
	})

	pkts, _ = p.Cmd("Alvaro mazo")
	util.Assert(enco.Contains(pkts, enco.Mazo), func() {
		t.Error("Alvaro debieria ser capaz de irse al mazo")
	})

	pkts, _ = p.Cmd("Roro truco")
	util.Assert(enco.Contains(pkts, enco.GritarTruco), func() {
		t.Error("Roro debieria ser capaz de gritar truco")
	})

	pkts, _ = p.Cmd("Alvaro real-envido")
	util.Assert(!enco.Contains(pkts, enco.ElEnvidoEstaPrimero), func() {
		t.Error("Alvaro no deberia ser capaz de tocar envido porque ya se fue al mazo")
	})
}

func TestFixError39(t *testing.T) {
	data := `{"Jugadores":[{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"},{"id":"Roro","nombre":"Roro","equipo":"Rojo"}],"cantJugadores":2,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":1,"Rojo":1},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"envido","puntaje":2,"cantadoPor":"Alvaro"},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Oro","valor":3},{"palo":"Copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":3},{"palo":"Oro","valor":5},{"palo":"Espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":1},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null}]}}`
	p, _ := Parse(data)

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

	data = `{"Jugadores":[{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"},{"id":"Roro","nombre":"Roro","equipo":"Rojo"}],"cantJugadores":2,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":1,"Rojo":1},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"envido","puntaje":2,"cantadoPor":"Alvaro"},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Oro","valor":3},{"palo":"Copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":3},{"palo":"Oro","valor":5},{"palo":"Espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":1},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null}]}}`
	p, _ = Parse(data)
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
	data := `{"puntuacion":20,"puntajes":{"Azul":3,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":1,"Rojo":1},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"deshabilitado","puntaje":3,"cantadoPor":"Alvaro","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Copa","valor":5},{"palo":"Copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Espada","valor":6},{"palo":"Copa","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":10},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null}]}}`
	p, _ := Parse(data)

	ok := util.All(
		p.Ronda.Envite.Estado == DESHABILITADO,
		len(p.Ronda.Envite.SinCantar) == 0,
	)
	util.Assert(ok, func() {
		t.Error("Alavaro ya canto la flor!")
	})

}

func TestFixRazonErronea(t *testing.T) {
	partidaJSON := `{"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":1,"Rojo":1},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"sinCantar":["Alvaro"]},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Copa","valor":5},{"palo":"Copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Espada","valor":6},{"palo":"Copa","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":10},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null}]}}`
	p, err := Parse(partidaJSON)
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
		var cont map[string]json.RawMessage
		json.Unmarshal(pkt.Message.Cont, &cont)

		if pkt.Message.Cod == string(enco.RondaGanada) {
			countMsgRondaGanada++

			var r string
			json.Unmarshal(cont["razon"], &r)

			ok := enco.Razon(r) != enco.EnvidoGanado
			if !ok {
				t.Error(`la razon por que ganan la ronda no deberia ser "por el envido"`)
			}
		}

		if pkt.Message.Cod == string(enco.SumaPts) {
			var r string
			json.Unmarshal(cont["razon"], &r)

			razon := enco.Razon(r)
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
	p, _ = Parse(partidaJSON)
	p.Cmd("Alvaro flor")
	p.Cmd("Alvaro truco")
	pkts, _ = p.Cmd("Roro no-quiero")

	countMsgRondaGanada = 0

	for _, pkt := range pkts {

		// obtengo el contenido del mensaje
		var cont map[string]json.RawMessage
		json.Unmarshal(pkt.Message.Cont, &cont)

		if pkt.Message.Cod == string(enco.RondaGanada) {
			countMsgRondaGanada++

			var r string
			json.Unmarshal(cont["razon"], &r)

			ok := enco.Razon(r) != enco.EnvidoGanado
			if !ok {
				t.Error(`la razon por que ganan la ronda no deberia ser "por el envido"`)
			}
		}

		if pkt.Message.Cod == string(enco.SumaPts) {

			var r string
			json.Unmarshal(cont["razon"], &r)

			ok := enco.Razon(r) == enco.TrucoQuerido || enco.Razon(r) == enco.TrucoNoQuerido
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

func TestCreadorDeEscenario(t *testing.T) {
	partidaJSON := `{"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":1,"Rojo":1},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"sinCantar":["Alvaro"]},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Copa","valor":5},{"palo":"Copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":2},{"palo":"Espada","valor":6},{"palo":"Copa","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"Rojo"}}],"muestra":{"palo":"Oro","valor":10},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null}]}}`
	p, err := Parse(partidaJSON)
	if err != nil {
		t.Error(err)
	}

	// ENVITE
	p.Cmd("Alvaro flor")

	// TRUCO
	p.Cmd("Alvaro truco")
	p.Cmd("Roro re-truco")
	p.Cmd("Alvaro vale-4")
	p.Cmd("Roro quiero")

	// PRIMERA mano la gana Alvaro
	p.Cmd("Alvaro 7 copa")
	p.Cmd("Roro 6 espada")

	// SEGUNDA mano la gana Roro
	p.Cmd("Alvaro 5 copa")
	p.Cmd("Roro 2 oro")

	// p.Cmd("Roro mazo") // <- hace que empiece una ronda nueva

	ok := util.All(
		p.Ronda.Envite.Estado == DESHABILITADO,
		// p.Ronda.Truco.Estado == VALE4QUERIDO,
		// p.Ronda.Truco.CantadoPor.Jugador.ID == "Roro",
		// p.Ronda.ManoEnJuego == Tercera,
	)

	if !ok {
		t.Error("algo salio mal")
	}

	json, _ := p.MarshalJSON()
	t.Log(string(json))

}
