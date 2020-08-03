package truco

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// el envido, la primera o la mentira
// el envido, la primera o la mentira
// el truco, la segunda o el rabón

// regexps
var (
	regexps = map[string]*regexp.Regexp{
		"jugadaSimple": regexp.MustCompile(`(?i)^([a-zA-Z0-9_-]+) ([a-zA-Z0-9_-]+)$`),
		"jugadaTirada": regexp.MustCompile(`(?i)^([a-zA-Z0-9_-]+) (1|2|3|4|5|6|7|10|11|12) (oro|copa|basto|espada)$`),
	}
)

// Msg mensajes a la capa de presentacion
type Msg struct {
	Dest []string
	Tipo string
	Cont string
}

func (msg Msg) String() string {
	return fmt.Sprintf(`<< [%s] (%s) : %s`, msg.Tipo, strings.Join(msg.Dest, "/"), msg.Cont)
}

func write(buff *bytes.Buffer, d *Msg) error {
	enc := gob.NewEncoder(buff)
	err := enc.Encode(d)
	return err
}

func read(buff *bytes.Buffer) (*Msg, error) {
	e := new(Msg)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func consume(buff *bytes.Buffer) {
	for {
		e, err := read(buff)
		if err == io.EOF {
			break
		}
		fmt.Println(*e)
	}
}

// Partida :
type Partida struct {
	PartidaDT
	Stdout *bytes.Buffer `json:"-"`
	ErrCh  chan bool     `json:"-"`
}

func (p *Partida) parseJugada(cmd string) (IJugada, error) {

	var jugada IJugada

	// comando simple son
	// jugadas sin parametro del tipo `$autor $jugada`
	match := regexps["jugadaSimple"].FindAllStringSubmatch(cmd, 1)

	if match != nil {
		jugadorStr, jugadaStr := match[0][1], match[0][2]

		manojo, err := p.Ronda.getManojoByStr(jugadorStr)
		if err != nil {
			return nil, fmt.Errorf("Usuario %s no encontrado", jugadorStr)
		}

		jugadaStr = strings.ToLower(jugadaStr)

		switch jugadaStr {
		// toques
		case "envido":
			jugada = tocarEnvido{Jugada{autor: manojo}}
		case "real-envido":
			jugada = tocarRealEnvido{Jugada{autor: manojo}}
		case "falta-envido":
			jugada = tocarFaltaEnvido{Jugada{autor: manojo}}

		// cantos
		case "flor":
			jugada = cantarFlor{Jugada{autor: manojo}}
		case "contra-flor":
			jugada = cantarContraFlor{Jugada{autor: manojo}}
		case "contra-flor-al-resto":
			jugada = cantarContraFlorAlResto{Jugada{autor: manojo}}

		// gritos
		case "truco":
			jugada = gritarTruco{Jugada{autor: manojo}}
		case "re-truco":
			jugada = gritarReTruco{Jugada{autor: manojo}}
		case "vale-4":
			jugada = gritarVale4{Jugada{autor: manojo}}

		// respuestas
		case "quiero":
			jugada = responderQuiero{Jugada{autor: manojo}}
		case "no-quiero":
			jugada = responderNoQuiero{Jugada{autor: manojo}}
		// case "tiene":
		// 	jugada = responderNoQuiero{Jugada{autor: manojo}}

		// acciones
		case "mazo":
			jugada = irseAlMazo{Jugada{autor: manojo}}
		case "tirar":
			jugada = irseAlMazo{Jugada{autor: manojo}}
		default:
			return nil, fmt.Errorf("No existe esa jugada")
		}
	} else {
		match = regexps["jugadaTirada"].FindAllStringSubmatch(cmd, 1)
		if match == nil {
			return nil, fmt.Errorf("No existe esa jugada")
		}
		jugadorStr := match[0][1]
		valorStr, paloStr := match[0][2], match[0][3]

		manojo, err := p.Ronda.getManojoByStr(jugadorStr)
		if err != nil {
			return nil, fmt.Errorf("Usuario %s no encontrado", jugadorStr)
		}

		carta, err := parseCarta(valorStr, paloStr)
		if err != nil {
			return nil, err
		}

		jugada = tirarCarta{
			Jugada{autor: manojo},
			*carta,
		}
	}

	return jugada, nil
}

func (p *Partida) byeBye() {
	if p.Terminada() {

		write(p.Stdout, &Msg{
			Dest: []string{"ALL"},
			Tipo: "ok",
			Cont: fmt.Sprintf("Se acabo la partida! el ganador fue el equipo %s",
				p.elQueVaGanando().String()),
		})

		write(p.Stdout, &Msg{
			Dest: []string{"ALL"},
			Tipo: "ok",
			Cont: fmt.Sprintf("BYE BYE!"),
		})

	}
}

// Cmd nexo capa presentacion con capa logica
func (p *Partida) Cmd(cmd string) error {

	if p.Terminada() {
		return fmt.Errorf("La partida ya termino")
	}

	// checkeo sintactico
	// ok := regexp.MustCompile(`^(\w|-)+\s(\w|-)+\n?$`).MatchString(cmd)
	ok := true
	if !ok {
		return fmt.Errorf("Sintaxis invalida: comando incorrecto")
	}

	// checkeo semantico
	jugada, err := p.parseJugada(cmd)
	if err != nil {
		return err
	}

	jugada.hacer(p)

	return nil
}

func (p *Partida) notify() {
	// ojo primero hay que grabar el buff, luego avisar
	write(p.Stdout, &Msg{[]string{"ALL"}, "INT", "INTERRUMPING!!"})
	p.ErrCh <- true
}

// NuevaPartida retorna n)ueva partida; error si hubo
func NuevaPartida(puntuacion Puntuacion, equipoAzul, equipoRojo []string) (*Partida, error) {

	mismaCantidadDeJugadores := len(equipoRojo) == len(equipoAzul)
	cantJugadores := len(equipoRojo) + len(equipoAzul)
	cantidadCorrecta := contains([]int{2, 4, 6}, cantJugadores) // puede ser 2, 4 o 6
	ok := mismaCantidadDeJugadores && cantidadCorrecta

	if !ok {
		return nil, fmt.Errorf(`La cantidad de jugadores no es correcta`)
	}
	// paso a crear los jugadores; intercalados
	var jugadores []Jugador
	// para cada rjo que agrego; le agrego tambien su mano
	for i := range equipoRojo {
		// uso como id sus nombres
		nuevoJugadorRojo := Jugador{equipoRojo[i], equipoRojo[i], Rojo}
		nuevoJugadorAzul := Jugador{equipoAzul[i], equipoAzul[i], Azul}
		jugadores = append(jugadores, nuevoJugadorAzul, nuevoJugadorRojo)
	}

	p := Partida{
		PartidaDT: PartidaDT{
			Puntuacion:    puntuacion,
			CantJugadores: cantJugadores,
			jugadores:     jugadores,
		},
	}

	p.Stdout = new(bytes.Buffer)
	p.ErrCh = make(chan bool, 1)

	p.Puntajes = make(map[Equipo]int)
	p.Puntajes[Rojo] = 0
	p.Puntajes[Azul] = 0

	elMano := JugadorIdx(0)
	p.nuevaRonda(elMano)

	write(p.Stdout, &Msg{
		Dest: []string{"ALL"},
		Tipo: "ok",
		Cont: fmt.Sprintf("Empieza una nueva ronda"),
	})

	write(p.Stdout, &Msg{
		Dest: []string{"ALL"},
		Tipo: "ok",
		Cont: fmt.Sprintf("La mano y el turno es %s\n", p.Ronda.getElMano().Jugador.Nombre),
	})

	return &p, nil
}
