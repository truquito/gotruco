package main

import (
	"fmt"

	"github.com/truquito/truco/pdt"
)

var terminals uint = 0

func rec_play(p *pdt.Partida) {
	bs, _ := p.MarshalJSON()

	// para la partida dada, todas las jugadas posibles
	chis := pdt.Chis(p)

	// las juego
	for mix := range chis {
		for aix := range chis[mix] {
			p, _ = pdt.Parse(string(bs), true)
			pkts2 := chis[mix][aix].Hacer(p)
			if pdt.IsDone(pkts2) {
				terminals++
			} else {
				rec_play(p)
			}
		}
	}
}

func main() {
	partidaJSON := `{"limiteEnvido":4,"cantJugadores":2,"puntuacion":20,"puntajes":{"azul":19,"rojo":19},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":""},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"oro","valor":5},{"palo":"espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":1},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]}}`
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
TIME:71.19 RAM:52352

*/
