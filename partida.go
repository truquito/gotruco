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

func write(buff *bytes.Buffer, d *Pkt) error {
	// registros
	gob.Register(ContSumPts{})

	enc := gob.NewEncoder(buff)
	err := enc.Encode(d)
	return err
}

// Read retorna el pkt mas antiguo sin leer
func Read(buff *bytes.Buffer) (*Pkt, error) {
	e := new(Pkt)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Consume consume el buffer
func Consume(buff *bytes.Buffer) {
	for {
		e, err := Read(buff)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
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

		write(p.Stdout, &Pkt{
			Dest: []string{"ALL"},
			Msg: Msg{
				Tipo: "Fin-Partida",
				Nota: fmt.Sprintf("Se acabo la partida! el ganador fue el equipo %s",
					p.elQueVaGanando().String()),
			},
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
	write(p.Stdout, &Pkt{
		Dest: []string{"ALL"},
		Msg: Msg{
			Tipo: "TimeOut",
			Nota: "INTERRUMPING!! Roro tardo demasiado en jugar. Mano ganada por Rojo",
		},
	})
	p.ErrCh <- true
}

// NuevaPartida retorna n)ueva partida; error si hubo
func NuevaPartida(puntuacion Puntuacion, equipoAzul, equipoRojo []string) (*Partida, error) {

	partidaDt, err := NuevaPartidaDt(puntuacion, equipoAzul, equipoRojo)

	if err != nil {
		return nil, err
	}

	p := Partida{
		PartidaDT: *partidaDt,
	}

	p.Stdout = new(bytes.Buffer)
	p.ErrCh = make(chan bool, 1)

	// write(p.Stdout, &Pkt{
	// 	Dest: []string{"ALL"},
	// 	Msg: Msg{
	// 		Tipo: "Nueva-Partida",
	// 		Cont: nil, // "pers aqui"
	// 	},
	// })

	return &p, nil
}
