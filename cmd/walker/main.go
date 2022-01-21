package main

import (
	"fmt"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/pdt"
)

func isDone(pkts []*enco.Packet) bool {
	for _, pkt := range pkts {
		if pkt.Message.Cod == enco.NuevaPartida ||
			pkt.Message.Cod == enco.NuevaRonda ||
			pkt.Message.Cod == enco.RondaGanada {
			return true
		}
	}
	return false
}

var terminals uint = 0

func rec_play(p *pdt.Partida) {

	bs, _ := p.MarshalJSON()

	// para la partida dada, todas las jugadas posibles
	aa := pdt.GetAA(p)

	// las juego
	for mix := range aa {
		for aix := range aa[mix] {
			if !aa[mix][aix] {
				continue
			}
			p, _ = pdt.Parse(string(bs))
			pkts := aa[mix].ToJugada(p, mix, aix).Hacer(p)
			if isDone(pkts) {
				terminals++
			} else {
				rec_play(p)
			}
		}
	}
}

func main() {

	partidaJSON := `{"puntuacion":20,"puntajes":{"Azul":19,"Rojo":19},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":1,"Rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":null},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":6},{"palo":"Oro","valor":3},{"palo":"Copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":3},{"palo":"Oro","valor":5},{"palo":"Espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","equipo":"Rojo"}}],"muestra":{"palo":"Copa","valor":1},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas":null},{"resultado":"ganoRojo","ganador":"","cartasTiradas":null},{"resultado":"ganoRojo","ganador":"","cartasTiradas":null}]}}`
	p, err := pdt.Parse(partidaJSON)
	if err != nil {
		panic(err)
	}

	rec_play(p)

	fmt.Println("termino", terminals)

}
