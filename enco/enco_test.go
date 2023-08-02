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
		m := DiceTengo{
			"pepe",
			123,
		}

		bs, _ := json.Marshal(m)
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
// 		Msg(TTirarCarta, "Alvaro", "Basto", 6),
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
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null}]}}`
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
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null},{"resultado":"ganoRojo","ganador": "","cartasTiradas":null}]}}`
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
