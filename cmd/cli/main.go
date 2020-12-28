package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/filevich/truco"
	"github.com/filevich/truco/deco"
	"github.com/filevich/truco/enco"
)

var reader = bufio.NewReader(os.Stdin)

func readLn(prefix string) string {
	fmt.Printf(prefix)
	cmd, _ := reader.ReadString('\n')
	return strings.TrimSuffix(cmd, "\n")
}

func handleIO() {
	for {
		cmd := readLn("")
		ioCh <- cmd
	}
}

var ioCh chan string = make(chan string, 1)

func main() {

	logfile := newLogFile("/home/jp/Workspace/_tmp/truco_logs/")

	// p, _ := truco.NuevaPartida(20, []string{"Alvaro"}, []string{"Roro"})
	p, out, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	pJSON, _ := p.MarshalJSON()
	logfile.Write(string(pJSON))

	// partidaJSON := `{"Jugadores":[{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"},{"id":"Roro","nombre":"Roro","equipo":"Rojo"},{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"},{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"},{"id":"Andres","nombre":"Andres","equipo":"Azul"},{"id":"Richard","nombre":"Richard","equipo":"Rojo"}],"cantJugadores":6,"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":3,"Rojo":3},"elMano":0,"turno":0,"pies":[0,0],"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":null},"truco":{"cantadoPor":null,"estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":4},{"palo":"Oro","valor":3},{"palo":"Basto","valor":7}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":1},{"palo":"Oro","valor":12},{"palo":"Espada","valor":11}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Copa","valor":2},{"palo":"Espada","valor":4},{"palo":"Basto","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Adolfo","nombre":"Adolfo","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":1},{"palo":"Basto","valor":11},{"palo":"Espada","valor":3}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Renzo","nombre":"Renzo","equipo":"Rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":12},{"palo":"Oro","valor":6},{"palo":"Basto","valor":5}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Andres","nombre":"Andres","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":11},{"palo":"Espada","valor":6},{"palo":"Basto","valor":12}],"cartasNoJugadas":[true,true,true],"ultimaTirada":0,"jugador":{"id":"Richard","nombre":"Richard","equipo":"Rojo"}}],"muestra":{"palo":"Espada","valor":10},"manos":[{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null},{"resultado":"ganoRojo","ganador":null,"cartasTiradas":null}]}}`
	// p.FromJSON(partidaJSON)

	fmt.Println(p)
	enco.Consume(out, func(pkt *enco.Packet) {
		fmt.Print(deco.Stringify(pkt, &p.PartidaDT))
	})

	// hago una gorutine (y channel para avisar) para el io
	go handleIO()

	for {
		select {
		case cmd := <-ioCh:
			logfile.Write(cmd)
			err := p.Cmd(cmd)
			if err != nil {
				fmt.Println("<< " + err.Error())
			}
			enco.Consume(out, func(pkt *enco.Packet) {
				fmt.Print(deco.Stringify(pkt, &p.PartidaDT))
			})
			fmt.Println(p)
		case <-p.ErrCh:
			enco.Consume(out, func(pkt *enco.Packet) {
				fmt.Print(deco.Stringify(pkt, &p.PartidaDT))
			})
			fmt.Printf(">> ")
		}

		if p.Terminada() {
			break
		}
	}

}
