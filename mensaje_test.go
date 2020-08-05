package truco

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestMsg(t *testing.T) {

	m0 := Msg{
		Tipo: "Nueva-Partida",
		Cont: nil,
	}

	m1 := Msg{
		Tipo: "Nueva-Ronda",
		Cont: nil,
	}

	m2 := Msg{
		Tipo: "Error",
		Nota: "Se produjo un error",
	}

	m3 := Msg{
		Tipo: "Sumar-Puntos",
		Nota: "El envido lo gano Alvaro con 32 pts de envido. +2 pts para el equipo Rojo",
		Cont: ContSumPts{
			Pts:    3,
			Equipo: "Rojo",
		}.ToJSON(),
	}

	m4 := Msg{
		Tipo: "TimeOut",
		Nota: "Roro tardo demasiado en jugar. Mano ganada por Rojo",
	}

	m5 := Msg{
		Tipo: "Tirar-Carta",
		Cont: ContTirarCarta{
			Autor: "Alvaro",
			Carta: Carta{
				Palo:  Basto,
				Valor: 6,
			},
		}.ToJSON(),
	}

	m6 := Msg{
		Tipo: "Gritar-Truco",
		Cont: []byte("Alvaro"),
	}

	fmt.Println(m0.ToJSON())
	fmt.Println(m1.ToJSON())
	fmt.Println(m2.ToJSON())
	fmt.Println(m3.ToJSON())
	fmt.Println(m4.ToJSON())
	fmt.Println(m5.ToJSON())
	fmt.Println(m6.ToJSON())

	var buff *bytes.Buffer = new(bytes.Buffer)
	// var err error

	// write
	pkt1 := &Pkt{
		Dest: []string{"ALL"},
		Msg:  m1,
	}

	pkt2 := &Pkt{
		Dest: []string{"ALL"},
		Msg:  m2,
	}

	pkt3 := &Pkt{
		Dest: []string{"ALL"},
		Msg:  m3,
	}

	write(buff, pkt1)
	write(buff, pkt2)
	write(buff, pkt3)

	Consume(buff)
}

func TestMsgNuevaPartida(t *testing.T) {

	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"jugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"jugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.fromJSON(partidaJSON)

	per, _ := p.Perspectiva("Alvaro")

	msg := Msg{
		Tipo: "Nueva-Partida",
		Cont: []byte(per.ToJSON()),
	}

	fmt.Println(msg.ToJSON())
}

func TestParseMsgPartidaDT(t *testing.T) {
	p, _ := NuevaPartida(a20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"jugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"jugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.fromJSON(partidaJSON)

	m := Msg{
		Tipo: "Nueva-Partida",
		Nota: "nota aqui",
		Cont: []byte(p.PartidaDT.ToJSON()),
	}

	mData := m.ToJSON()

	fmt.Println(mData)

	var mParsed Msg
	if err := json.Unmarshal([]byte(mData), &mParsed); err != nil {
		log.Fatal(err)
	}

	// ahora parseo lo raw
	var pDt PartidaDT
	if err := json.Unmarshal(mParsed.Cont, &pDt); err != nil {
		log.Fatal(err)
	}

	fmt.Println(pDt.CantJugadores)
}
