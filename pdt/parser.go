package pdt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/filevich/truco/enco"
)

// regexps
var (
	regexps = map[string]*regexp.Regexp{
		"jugadaSimple": regexp.MustCompile(`(?i)^([a-zA-Z0-9_-]+) ([a-zA-Z0-9_-]+)$`),
		"jugadaTirada": regexp.MustCompile(`(?i)^([a-zA-Z0-9_-]+) (1|2|3|4|5|6|7|10|11|12) (oro|copa|basto|espada)$`),
	}
)

func checkeoSemantico(p *Partida, cmd string) (IJugada, error) {

	var jugada IJugada

	// comando simple son
	// jugadas sin parametro del tipo `$autor $jugada`
	match := regexps["jugadaSimple"].FindAllStringSubmatch(cmd, 1)

	if match != nil {
		jugadorStr, jugadaStr := match[0][1], match[0][2]

		m, err := p.GetManojoByStr(jugadorStr)
		if err != nil {
			return nil, fmt.Errorf("usuario %s no encontrado", jugadorStr)
		}

		jugadaStr = strings.ToLower(jugadaStr)

		switch jugadaStr {
		// toques
		case "envido":
			jugada = TocarEnvido{Manojo: m}
		case "real-envido":
			jugada = TocarRealEnvido{Manojo: m}
		case "falta-envido":
			jugada = TocarFaltaEnvido{Manojo: m}

		// cantos
		case "flor":
			jugada = CantarFlor{Manojo: m}
		case "contra-flor":
			jugada = CantarContraFlor{Manojo: m}
		case "contra-flor-al-resto":
			jugada = CantarContraFlorAlResto{Manojo: m}

		// gritos
		case "truco":
			jugada = GritarTruco{Manojo: m}
		case "re-truco":
			jugada = GritarReTruco{Manojo: m}
		case "vale-4":
			jugada = GritarVale4{Manojo: m}

		// respuestas
		case "quiero":
			jugada = ResponderQuiero{Manojo: m}
		case "no-quiero":
			jugada = ResponderNoQuiero{Manojo: m}
		// case "tiene":
		// 	jugada = responderNoQuiero{Manojo: m}

		// acciones
		case "mazo":
			jugada = IrseAlMazo{Manojo: m}
		case "tirar":
			jugada = IrseAlMazo{Manojo: m}
		default:
			return nil, fmt.Errorf("no existe esa jugada")
		}
	} else {
		match = regexps["jugadaTirada"].FindAllStringSubmatch(cmd, 1)
		if match == nil {
			return nil, fmt.Errorf("no existe esa jugada")
		}
		jugadorStr := match[0][1]
		valorStr, paloStr := match[0][2], match[0][3]

		m, err := p.GetManojoByStr(jugadorStr)
		if err != nil {
			return nil, fmt.Errorf("usuario %s no encontrado", jugadorStr)
		}

		c, err := ParseCarta(valorStr, paloStr)
		if err != nil {
			return nil, err
		}

		jugada = TirarCarta{
			Manojo: m,
			Carta:  *c,
		}
	}

	return jugada, nil
}

func ParseJugada(p *Partida, cmd string) (IJugada, error) {
	if p.Terminada() {
		return nil, fmt.Errorf("la partida ya termino")
	}

	// checkeo sintactico
	// ok := regexp.MustCompile(`^(\w|-)+\s(\w|-)+\n?$`).MatchString(cmd)
	ok := true
	if !ok {
		return nil, fmt.Errorf("sintaxis invalida: comando incorrecto")
	}

	// checkeo semantico
	return checkeoSemantico(p, cmd)

}

func (p *Partida) byeBye() []*enco.Packet {
	pkts := make([]*enco.Packet, 0)

	if p.Terminada() {

		var s string

		if p.Jugadores[0].Equipo == p.ElQueVaGanando() {
			s = p.Jugadores[0].ID
		} else {
			s = p.Jugadores[1].ID
		}

		pkts = append(pkts, enco.Pkt(
			enco.Dest("ALL"),
			enco.Msg(enco.ByeBye, s),
		))
	}

	return pkts
}

// Cmd nexo capa presentacion con capa logica
func (p *Partida) Cmd(cmd string) ([]*enco.Packet, error) {

	// checkeo semantico
	jugada, err := ParseJugada(p, cmd)
	if err != nil {
		return nil, err
	}

	pkts := jugada.Hacer(p)

	if p.Terminada() {
		// AGREGARLE LoS PKTS ACA
		pkts = append(pkts, p.byeBye()...)
	}

	return pkts, nil
}
