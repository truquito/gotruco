package truco

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/pdt"
	"github.com/filevich/truco/ptr"
)

// el envido, la Primera o la mentira
// el envido, la Primera o la mentira
// el truco, la Segunda o el rabón

// regexps
var (
	regexps = map[string]*regexp.Regexp{
		"jugadaSimple": regexp.MustCompile(`(?i)^([a-zA-Z0-9_-]+) ([a-zA-Z0-9_-]+)$`),
		"jugadaTirada": regexp.MustCompile(`(?i)^([a-zA-Z0-9_-]+) (1|2|3|4|5|6|7|10|11|12) (oro|copa|basto|espada)$`),
	}
)

// Partida :
type Partida struct {
	*pdt.PartidaDT
	out   io.Writer `json:"-"`
	ErrCh chan bool `json:"-"`
}

func (p *Partida) parseJugada(cmd string) (IJugada, error) {

	var jugada IJugada

	// comando simple son
	// jugadas sin parametro del tipo `$autor $jugada`
	match := regexps["jugadaSimple"].FindAllStringSubmatch(cmd, 1)

	if match != nil {
		jugadorStr, jugadaStr := match[0][1], match[0][2]

		manojo, err := p.GetManojoByStr(jugadorStr)
		if err != nil {
			return nil, fmt.Errorf("Usuario %s no encontrado", jugadorStr)
		}

		jugadaStr = strings.ToLower(jugadaStr)

		switch jugadaStr {
		// toques
		case "envido":
			jugada = tocarEnvido{manojo}
		case "real-envido":
			jugada = tocarRealEnvido{manojo}
		case "falta-envido":
			jugada = tocarFaltaEnvido{manojo}

		// cantos
		case "flor":
			jugada = cantarFlor{manojo}
		case "contra-flor":
			jugada = cantarContraFlor{manojo}
		case "contra-flor-al-resto":
			jugada = cantarContraFlorAlResto{manojo}

		// gritos
		case "truco":
			jugada = gritarTruco{manojo}
		case "re-truco":
			jugada = gritarReTruco{manojo}
		case "vale-4":
			jugada = gritarVale4{manojo}

		// respuestas
		case "quiero":
			jugada = responderQuiero{manojo}
		case "no-quiero":
			jugada = responderNoQuiero{manojo}
		// case "tiene":
		// 	jugada = responderNoQuiero{manojo}

		// acciones
		case "mazo":
			jugada = irseAlMazo{manojo}
		case "tirar":
			jugada = irseAlMazo{manojo}
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

		manojo, err := p.GetManojoByStr(jugadorStr)
		if err != nil {
			return nil, fmt.Errorf("Usuario %s no encontrado", jugadorStr)
		}

		carta, err := pdt.ParseCarta(valorStr, paloStr)
		if err != nil {
			return nil, err
		}

		jugada = tirarCarta{
			manojo,
			*carta,
		}
	}

	return jugada, nil
}

func (p *Partida) byeBye() {
	if p.Terminada() {

		var s string
		if p.PartidaDT.EsManoAMano() { // retorna jugador
			if p.Jugadores[0].Equipo == p.ElQueVaGanando() {
				s = p.Jugadores[0].Nombre
			} else {
				s = p.Jugadores[1].Nombre
			}
		}

		enco.Write(p.out, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.ByeBye, s),
		))
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

// String retorna una representacion en formato de string
func (p *Partida) String() string {
	return ptr.Renderizar(p.PartidaDT)
}

func (p *Partida) notify() {

	// ojo primero hay que grabar el buff, luego avisar
	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TimeOut, "INTERRUMPING!! Roro tardo demasiado en jugar. Mano ganada por Rojo"),
	))

	p.ErrCh <- true
}

// Abandono da por ganada la partida al equipo contario
func (p *Partida) Abandono(jugador string) error {
	// encuentra al jugador
	manojo, err := p.GetManojoByStr(jugador)
	if err != nil {
		return fmt.Errorf("Usuario %s no encontrado", jugador)
	}
	// doy por ganador al equipo contrario
	equipoContrario := manojo.Jugador.GetEquipoContrario()
	ptsFaltantes := int(p.Puntuacion) - p.Puntajes[equipoContrario]
	p.SumarPuntos(equipoContrario, ptsFaltantes)

	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.Abandono, manojo.Jugador.ID),
	))

	return nil
}

// NuevaPartida retorna nueva partida; error si hubo
func NuevaPartida(puntuacion pdt.Puntuacion, equipoAzul, equipoRojo []string) (*Partida, io.Reader, error) {

	partidaDt, err := pdt.NuevaPartidaDt(puntuacion, equipoAzul, equipoRojo)

	if err != nil {
		return nil, nil, err
	}

	p := Partida{
		PartidaDT: partidaDt,
	}

	buff := new(bytes.Buffer)
	p.out = buff
	p.ErrCh = make(chan bool, 1)

	return &p, buff, nil
}
