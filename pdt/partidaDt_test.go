package pdt

import (
	"encoding/json"
	"testing"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/util"
)

func TestNoStruct(t *testing.T) {
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"truco"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Basto","valor":12},{"palo":"Oro","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Basto","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":10},{"palo":"Copa","valor":10},{"palo":"Basto","valor":2}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":6},{"palo":"Espada","valor":10},{"palo":"Basto","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Copa","valor":3},{"palo":"Espada","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":3},{"palo":"Espada","valor":11},{"palo":"Espada","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
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

	if !(len(p.Jugadores) == 6) {
		t.Error("deberia haber 6 jugadores")
	}

	t.Log(p.Ronda.Muestra)

	t.Log(Renderizar(p))
}

func TestFixPanic(t *testing.T) {
	p, _ := NuevaPartidaDt(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
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
	p, _ := NuevaPartidaDt(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
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
	data := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 0, "Rojo": 0}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 2, "Rojo": 2}, "elMano": 0, "turno": 0, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": null, "estado": "noCantado"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 3}, {"palo": "Oro", "valor": 7}, {"palo": "Copa", "valor": 7}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 10}, {"palo": "Basto", "valor": 5}, {"palo": "Oro", "valor": 6}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 2}, {"palo": "Oro", "valor": 3}, {"palo": "Espada", "valor": 7}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 1}, {"palo": "Espada", "valor": 6}, {"palo": "Espada", "valor": 10}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Copa", "valor": 3}, "manos": [{"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}]}} `
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
	p, _ := NuevaPartidaDt(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
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

	m, _ := p.GetManojoByStr("Renzo")
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
	data := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 0, "Rojo": 0}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 2, "Rojo": 2}, "elMano": 0, "turno": 0, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": null, "estado": "noCantado"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 3}, {"palo": "Oro", "valor": 7}, {"palo": "Copa", "valor": 7}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 10}, {"palo": "Basto", "valor": 5}, {"palo": "Oro", "valor": 6}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 2}, {"palo": "Oro", "valor": 3}, {"palo": "Espada", "valor": 7}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 1}, {"palo": "Espada", "valor": 6}, {"palo": "Espada", "valor": 10}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Copa", "valor": 3}, "manos": [{"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}]}} `
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
	data := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 5, "Rojo": 6}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 2, "Rojo": 2}, "elMano": 2, "turno": 2, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": null, "estado": "noCantado"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Basto", "valor": 2}, {"palo": "Basto", "valor": 6}, {"palo": "Espada", "valor": 5}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 6}, {"palo": "Espada", "valor": 11}, {"palo": "Copa", "valor": 12}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Oro", "valor": 3}, {"palo": "Espada", "valor": 10}, {"palo": "Oro", "valor": 10}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 1}, {"palo": "Basto", "valor": 11}, {"palo": "Espada", "valor": 3}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Copa", "valor": 5}, "manos": [{"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}]}}`
	p, err := Parse(data)
	if err != nil {
		t.Error(err)
	}
	t.Log(Renderizar(p))

	p.Cmd("Renzo mazo")
	m, _ := p.GetManojoByStr("Renzo")
	util.Assert(m.SeFueAlMazo, func() {
		t.Error("Deberia dejarlo irse al mazo")
	})

	p.Cmd("Adolfo falta-envido")

	pkts, _ := p.Cmd("Roro quiero")
	for _, pkt := range pkts {
		diceSonBuenas := pkt.Cod == int(enco.DiceSonBuenas)
		loDijoRenzo := string(pkt.Cont) == "\"Renzo\""
		if diceSonBuenas && loDijoRenzo {
			t.Error("No deberia poder responder desde ultratumba")
		}
	}

	// debio de haber ganado el envido
	for _, pkt := range pkts {
		if pkt.Cod == int(enco.SumaPts) {
			var t3 enco.Tipo3
			json.Unmarshal(pkt.Message.Cont, &t3)

			ok := util.All(
				t3.Puntos == 4,
				t3.Razon == int(enco.FaltaEnvidoGanado),
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
	data := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 0, "Rojo": 3}, "ronda": {"manoEnJuego": 1, "cantJugadoresEnJuego": {"Azul": 1, "Rojo": 1}, "elMano": 0, "turno": 3, "pies": [0, 0], "envite": {"estado": "deshabilitado", "puntaje": 3, "cantadoPor": {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 5}, {"palo": "Oro", "valor": 6}, {"palo": "Espada", "valor": 7}], "cartasNoJugadas": [false, true, true], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}}, "truco": {"cantadoPor": {"seFueAlMazo": false, "cartas": [{"palo": "Basto", "valor": 11}, {"palo": "Oro", "valor": 3}, {"palo": "Copa", "valor": 6}], "cartasNoJugadas": [true, true, false], "ultimaTirada": 2, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, "estado": "trucoQuerido"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Basto", "valor": 11}, {"palo": "Oro", "valor": 3}, {"palo": "Copa", "valor": 6}], "cartasNoJugadas": [true, true, false], "ultimaTirada": 2, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": true, "cartas": [{"palo": "Basto", "valor": 4}, {"palo": "Espada", "valor": 10}, {"palo": "Copa", "valor": 1}], "cartasNoJugadas": [true, false, false], "ultimaTirada": 1, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": true, "cartas": [{"palo": "Oro", "valor": 1}, {"palo": "Espada", "valor": 12}, {"palo": "Basto", "valor": 10}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 5}, {"palo": "Oro", "valor": 6}, {"palo": "Espada", "valor": 7}], "cartasNoJugadas": [false, true, true], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Espada", "valor": 1}, "manos": [{"resultado": "ganoRojo", "ganador": {"seFueAlMazo": true, "cartas": [{"palo": "Basto", "valor": 4}, {"palo": "Espada", "valor": 10}, {"palo": "Copa", "valor": 1}], "cartasNoJugadas": [true, false, false], "ultimaTirada": 1, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, "cartasTiradas": [{"palo": "Copa", "valor": 6}, {"palo": "Copa", "valor": 1}, {"palo": "Copa", "valor": 5}]}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": [{"palo": "Espada", "valor": 10}]}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}]}}`
	p, err := Parse(data)
	if err != nil {
		t.Error(err)
	}
	t.Log(Renderizar(p))

	pkts, _ := p.Cmd("Alvaro re-truco")
	ok := util.All(
		p.Ronda.Truco.Estado == RETRUCO,
		p.Ronda.Truco.CantadoPor.Jugador.ID == "Alvaro",
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
	partidaJSON := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 0, "Rojo": 0}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 2, "Rojo": 2}, "elMano": 0, "turno": 0, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": null, "estado": "noCantado"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 10}, {"palo": "Copa", "valor": 4}, {"palo": "Copa", "valor": 11}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Basto", "valor": 2}, {"palo": "Basto", "valor": 3}, {"palo": "Oro", "valor": 6}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 3}, {"palo": "Espada", "valor": 6}, {"palo": "Espada", "valor": 12}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 5}, {"palo": "Espada", "valor": 10}, {"palo": "Copa", "valor": 7}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Oro", "valor": 12}, "manos": [{"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}]}}`
	p, _ := Parse(partidaJSON)

	p.Ronda.Turno = 1
	p.Cmd("Roro truco")

	util.Assert(p.Ronda.Truco.Estado == TRUCO, func() {
		t.Error("Roro deberia ser capaz de gritar truco")
	})

	t.Log(Renderizar(p))

	ptsAntes := p.Puntajes[Azul]
	pkts, _ := p.Cmd("Alvaro flor")

	util.Assert(!enco.Contains(pkts, enco.DiceTengo), func() {
		t.Error("Alvaro no deberia decir cuanto tiene de flor ya que es el unico")
	})

	util.Assert(ptsAntes+6 == p.Puntajes[Azul], func() {
		t.Error("Debio de haber ganado +6 pts por las 2 flores")
	})

	// ejemplo con 1 solo jugador; tampoco deberia
	partidaJSON = `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 7, "Rojo": 0}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 2, "Rojo": 1}, "elMano": 1, "turno": 1, "pies": [0, 0], "envite": {"estado": "noCantadoAun", "puntaje": 0, "cantadoPor": null}, "truco": {"cantadoPor": null, "estado": "noCantado"}, "manojos": [{"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 11}, {"palo": "Copa", "valor": 7}, {"palo": "Copa", "valor": 3}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 10}, {"palo": "Copa", "valor": 2}, {"palo": "Basto", "valor": 7}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Espada", "valor": 11}, {"palo": "Basto", "valor": 4}, {"palo": "Oro", "valor": 10}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": true, "cartas": [{"palo": "Basto", "valor": 6}, {"palo": "Oro", "valor": 4}, {"palo": "Oro", "valor": 3}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Espada", "valor": 1}, "manos": [{"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}]}}`
	p, _ = Parse(partidaJSON)

	t.Log(Renderizar(p))

	ptsAntes = p.Puntajes[Azul]
	pkts, _ = p.Cmd("Alvaro flor")
	util.Assert(!enco.Contains(pkts, enco.DiceTengo), func() {
		t.Error("Alvaro no deberia decir cuanto tiene de flor ya que es el unico")
	})
	util.Assert(ptsAntes+3 == p.Puntajes[Azul], func() {
		t.Error("Debio de haber ganado +3 pts por la flor")
	})

}

func TestFixDecirSonBuenasDesdeUltratumba(t *testing.T) {
	p, _ := NuevaPartidaDt(A20, []string{"Alvaro", "Adolfo"}, []string{"Roro", "Renzo"})
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
		if pkt.Cod == int(enco.DiceSonBuenas) {
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
	data := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 6, "Rojo": 6}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 1, "Rojo": 2}, "elMano": 2, "turno": 1, "pies": [0, 0], "envite": {"estado": "deshabilitado", "puntaje": 3, "cantadoPor": {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 7}, {"palo": "Oro", "valor": 10}, {"palo": "Espada", "valor": 3}], "cartasNoJugadas": [true, false, true], "ultimaTirada": 1, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}}, "truco": {"cantadoPor": null, "estado": "noCantado"}, "manojos": [{"seFueAlMazo": true, "cartas": [{"palo": "Oro", "valor": 12}, {"palo": "Oro", "valor": 1}, {"palo": "Espada", "valor": 11}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 6}, {"palo": "Copa", "valor": 2}, {"palo": "Espada", "valor": 1}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 7}, {"palo": "Oro", "valor": 10}, {"palo": "Espada", "valor": 3}], "cartasNoJugadas": [true, false, true], "ultimaTirada": 1, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Basto", "valor": 6}, {"palo": "Oro", "valor": 3}, {"palo": "Copa", "valor": 4}], "cartasNoJugadas": [false, true, true], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Basto", "valor": 10}, "manos": [{"resultado": "ganoRojo", "ganador": null, "cartasTiradas": [{"palo": "Oro", "valor": 10}, {"palo": "Basto", "valor": 6}]}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}]}}`
	p, _ := Parse(data)

	ok := util.All(
		len(p.Ronda.Manos[0].CartasTiradas) == 2,
		p.Ronda.Manos[0].CartasTiradas[0].autor.Jugador.ID == "Adolfo",
		p.Ronda.Manos[0].CartasTiradas[1].autor.Jugador.ID == "Renzo",
	)
	util.Assert(ok, func() {
		t.Error("El mapeo de las cartas tiradas con sus autores no coincide")
	})
}

func TestFixPorQueDejaARoroTirarCarta(t *testing.T) {
	data := `{"Jugadores": [{"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}, {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}, {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}, {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}], "cantJugadores": 4, "puntuacion": 20, "puntajes": {"Azul": 6, "Rojo": 6}, "ronda": {"manoEnJuego": 0, "cantJugadoresEnJuego": {"Azul": 1, "Rojo": 2}, "elMano": 2, "turno": 1, "pies": [0, 0], "envite": {"estado": "deshabilitado", "puntaje": 3, "cantadoPor": {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 7}, {"palo": "Oro", "valor": 10}, {"palo": "Espada", "valor": 3}], "cartasNoJugadas": [true, false, true], "ultimaTirada": 1, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}}, "truco": {"cantadoPor": null, "estado": "noCantado"}, "manojos": [{"seFueAlMazo": true, "cartas": [{"palo": "Oro", "valor": 12}, {"palo": "Oro", "valor": 1}, {"palo": "Espada", "valor": 11}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Alvaro", "nombre": "Alvaro", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 6}, {"palo": "Copa", "valor": 2}, {"palo": "Espada", "valor": 1}], "cartasNoJugadas": [true, true, true], "ultimaTirada": 0, "jugador": {"id": "Roro", "nombre": "Roro", "equipo": "Rojo"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Copa", "valor": 7}, {"palo": "Oro", "valor": 10}, {"palo": "Espada", "valor": 3}], "cartasNoJugadas": [true, false, true], "ultimaTirada": 1, "jugador": {"id": "Adolfo", "nombre": "Adolfo", "equipo": "Azul"}}, {"seFueAlMazo": false, "cartas": [{"palo": "Basto", "valor": 6}, {"palo": "Oro", "valor": 3}, {"palo": "Copa", "valor": 4}], "cartasNoJugadas": [false, true, true], "ultimaTirada": 0, "jugador": {"id": "Renzo", "nombre": "Renzo", "equipo": "Rojo"}}], "muestra": {"palo": "Basto", "valor": 10}, "manos": [{"resultado": "ganoRojo", "ganador": null, "cartasTiradas": [{"palo": "Oro", "valor": 10}, {"palo": "Basto", "valor": 6}]}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}, {"resultado": "ganoRojo", "ganador": null, "cartasTiradas": null}]}}`
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
