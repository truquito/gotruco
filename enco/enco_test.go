package enco

import (
	"encoding/json"
	"testing"
)

func TestPacket(t *testing.T) {
	/*
		voy a armar el packet

		{ <- Packet
			dest: ["pepe"],
			msg: { <- Message::DiceSonBuenas
				cod: "DiceSonBuenas",
				cont: "miguel"
			}
		}
	*/
	// var cont json.RawMessage
	// bs, _ := json.Marshal("miguel")
	// cont = bs
	// packet1 := Packet{
	// 	Destination: []string{"pepe"},
	// 	Message: Message{
	// 		Cod:  TDiceSonBuenas,
	// 		Cont: cont,
	// 	},
	// }

	// t.Log(packet1) // contiene 2 punteros
	// recordar que un slice es un puntero a un array.
	// fmt.Println(packet1.Destination)

	// t.Log("----------------------------------")

	// a := A{
	// 	"pepe",
	// 	333,
	// }

	// bs3, _ := json.Marshal(a)
	// t.Log(string(bs3))

	t.Log("----------------------------------")

	{
		pkt := Env(
			ALL,
			DiceTengo{
				"pepe",
				123,
			},
		)

		bs, _ := json.Marshal(pkt)
		t.Log(string(bs))
	}

	t.Log("----------------------------------")

	{
		pkt := Env(
			ALL,
			QuieroTruco("romualdo"),
		)

		bs, _ := json.Marshal(pkt)
		t.Log(string(bs))
	}

	t.Log("----------------------------------")

	{
		m := NuevaPartida{
			"yada yada",
		}

		bs, _ := json.Marshal(m)
		t.Log(string(bs))
	}

	m := DiceSonBuenas("asdasd")
	t.Log(m.Cod(), m)
}

// func TestCast(t *testing.T) {

// 	var buff *bytes.Buffer = new(bytes.Buffer)

// 	Write(buff, Pkt(
// 		Dest("Alvaro", "Roro"),
// 		Msg(TTirarCarta, "Alvaro", "basto", 6),
// 	))

// 	Write(buff, Pkt(
// 		Dest("ALL"),
// 		Msg(TError, "Se produjo un error"),
// 	))

// 	Write(buff, Pkt(
// 		Dest("ALL"),
// 		Msg(TSumaPts, "Alvaro", EnvidoGanado, 3),
// 	))

// 	Write(buff, Pkt(
// 		Dest("ALL"),
// 		Msg(TTimeOut, "Roro tardo demasiado en jugar. Mano ganada por Rojo"),
// 	))

// 	Write(buff, Pkt(
// 		Dest("ALL"),
// 		Msg(TGritarTruco, "Alvaro"),
// 	))

// 	for {
// 		e, err := Read(buff)
// 		if err == io.EOF {
// 			break
// 		} else if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		t.Log(e)
// 	}
// }

/*
func TestMsgNuevaPartida(t *testing.T) {

	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":5},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":5},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"copa","valor":7},{"palo":"basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"copa","valor":6},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"espada","valor":1},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"basto","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":12},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":5},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.fromJSON(partidaJSON)

	per, _ := p.Perspectiva("Alvaro")

	msg := Msg{
		Tipo: "Nueva-Partida",
		Cont: []byte(per.MarshalJSON()),
	}

	fmt.Println(msg.MarshalJSON())
}

func TestParseMsgPartida(t *testing.T) {
	p, _ := NuevaPartida(A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":5},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":5},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}]},"truco":{"cantadoPor":null,"estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":2},{"palo":"copa","valor":7},{"palo":"basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":2},{"palo":"copa","valor":6},{"palo":"oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"espada","valor":1},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"basto","valor":7},{"palo":"oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":5},{"palo":"basto","valor":12},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":6},{"palo":"espada","valor":5},{"palo":"espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"rojo"}}],"muestra":{"palo":"espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []},{"resultado":"ganoRojo","ganador": "","cartasTiradas": []}]}}`
	p.fromJSON(partidaJSON)

	m := Msg{
		Tipo: "Nueva-Partida",
		Nota: "nota aqui",
		Cont: []byte(p.Partida.MarshalJSON()),
	}

	mData := m.MarshalJSON()

	fmt.Println(mData)

	var mParsed Msg
	if err := json.Unmarshal([]byte(mData), &mParsed); err != nil {
		log.Fatal(err)
	}

	// ahora parseo lo raw
	var pDt Partida
	if err := json.Unmarshal(mParsed.Cont, &pDt); err != nil {
		log.Fatal(err)
	}

	fmt.Println(pDt.CantJugadores)
}
*/
