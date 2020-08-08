package truco

import (
	"fmt"
	"testing"

	"github.com/jpfilevich/truco/ptr"

	"github.com/jpfilevich/truco/pdt"
)

func TestCopy(t *testing.T) {
	// p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	// partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	// p.FromJSON(partidaJSON)

	// copia := cheepCopy(p.Perspectiva)
	// copia.Ronda.Manojos[0].Cartas[0] = nil

	// fmt.Println(copia.MarshalJSON())
	// fmt.Println()
	// fmt.Println(p.MarshalJSON())
}

func TestPerspectiva(t *testing.T) {
	p, _ := NuevaPartida(pdt.A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null,"JugadoresConFlor":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"JugadoresConFlorQueNoCantaron":[{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}]},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Copa","valor":7},{"palo":"Basto","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":2},{"palo":"Copa","valor":6},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Basto","valor":11},{"palo":"Espada","valor":1},{"palo":"Basto","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Basto","valor":7},{"palo":"Oro","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Basto","valor":12},{"palo":"Espada","valor":10}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":6},{"palo":"Espada","valor":5},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":3},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)

	per, _ := p.Perspectiva("Alvaro")
	fmt.Println(per.MarshalJSON())
}

// borrar este test al pedo
// es para ver por que da panic cuando una carta es nil
func TestFixPanicCartaNil(t *testing.T) {
	p, _ := pdt.NuevaPartidaDt(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})

	// c1 = p.Ronda.Manojos[0].Cartas[0]
	// c2 = p.Ronda.Manojos[0].Cartas[1]
	// c3 = p.Ronda.Manojos[0].Cartas[2]

	partidaJSON := `{"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[null,null,null],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":5},{"palo":"Espada","valor":4},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Basto","valor":6},{"palo":"Oro","valor":1}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Oro","valor":10},{"palo":"Oro","valor":4}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":4},{"palo":"Basto","valor":7},{"palo":"Espada","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":12},{"palo":"Espada","valor":10},{"palo":"Oro","valor":6}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":6},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	p.FromJSON(partidaJSON)
	ptr.Renderizar(p)
}
