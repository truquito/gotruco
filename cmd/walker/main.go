package main

import (
	"fmt"

	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

var terminals uint = 0

func countCodMsgs(pkts []enco.Envelope, cod enco.CodMsg) int {
	total := 0
	for _, pkt := range pkts {
		if pkt.Message.Cod() == cod {
			total++
		}
	}
	return total
}

func rec_play(p *pdt.Partida) {
	bs, _ := p.MarshalJSON()

	// para la partida dada, todas las jugadas posibles
	chis := pdt.Chis(p)

	// las juego
	for mix := range chis {
		for aix := range chis[mix] {
			p, _ = pdt.Parse(string(bs), true)
			a := chis[mix][aix]
			pkts2 := a.Hacer(p)
			isDone := enco.Contains(pkts2, enco.TRondaGanada)
			if countCodMsgs(pkts2, enco.TRondaGanada) > 1 {
				fmt.Println(string(bs))
				fmt.Println(a.String())
				panic(555)
			}
			if p.Terminada() && !isDone {
				panic(123)
			}
			if isDone {
				terminals++
			} else {
				rec_play(p)
			}
		}
	}
}

func main() {
	partidaJSON := `{"limiteEnvido":4,"cantJugadores":2,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"oro","valor":5},{"palo":"espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":1},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]}}`
	p, err := pdt.Parse(partidaJSON, true)
	if err != nil {
		panic(err)
	}
	rec_play(p)
	fmt.Println("termino", terminals)
}

/*

version: v0.2.x
cpu: i5-12600k
termino 1,807,482
TIME:71.19

cpu: m2 (fanless)
termino 1,807,482
TIME:88.95s

*/
