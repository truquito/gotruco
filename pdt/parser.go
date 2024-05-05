package pdt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/truquito/gotruco/enco"
)

// regexps
var (
	jugadaSimple = regexp.MustCompile(`(?i)^([a-zA-Z0-9_#-]+) ([a-zA-Z0-9_-]+)$`)
	jugadaTirada = regexp.MustCompile(`(?i)^([a-zA-Z0-9_#-]+) (1|2|3|4|5|6|7|10|11|12) (oro|copa|basto|espada)$`)
)

func checkeoSemantico(p *Partida, cmd string) (IJugada, error) {

	var jugada IJugada

	// comando simple son
	// jugadas sin parametro del tipo `$autor $jugada`
	match := jugadaSimple.FindAllStringSubmatch(cmd, 1)

	if match != nil {

		jugadorStr, jugadaStr := match[0][1], match[0][2]

		m := p.Ronda.Manojo(jugadorStr)
		if m == nil {
			// segundo intento
			m = p.Ronda.Manojo(strings.Title(jugadorStr))
			if m == nil {
				return nil, fmt.Errorf("usuario %s no encontrado", jugadorStr)
			}
		}

		jugadaStr = strings.ToLower(jugadaStr)

		switch jugadaStr {
		// toques
		case "envido":
			jugada = TocarEnvido{JID: m.Jugador.ID}
		case "real-envido":
			jugada = TocarRealEnvido{JID: m.Jugador.ID}
		case "falta-envido":
			jugada = TocarFaltaEnvido{JID: m.Jugador.ID}

		// cantos
		case "flor":
			jugada = CantarFlor{JID: m.Jugador.ID}
		case "contra-flor":
			jugada = CantarContraFlor{JID: m.Jugador.ID}
		case "contra-flor-al-resto":
			jugada = CantarContraFlorAlResto{JID: m.Jugador.ID}

		// gritos
		case "truco":
			jugada = GritarTruco{JID: m.Jugador.ID}
		case "re-truco":
			jugada = GritarReTruco{JID: m.Jugador.ID}
		case "vale-4":
			jugada = GritarVale4{JID: m.Jugador.ID}

		// respuestas
		case "quiero":
			jugada = ResponderQuiero{JID: m.Jugador.ID}
		case "no-quiero":
			jugada = ResponderNoQuiero{JID: m.Jugador.ID}
		// case "tiene":
		// 	jugada = responderNoQuiero{JID: m.Jugador.ID}

		// acciones
		case "mazo":
			jugada = IrseAlMazo{JID: m.Jugador.ID}
		default:
			return nil, fmt.Errorf("no existe esa jugada")
		}

	} else {

		match = jugadaTirada.FindAllStringSubmatch(cmd, 1)
		if match == nil {
			return nil, fmt.Errorf("no existe esa jugada")
		}
		jugadorStr := match[0][1]
		valorStr, paloStr := match[0][2], match[0][3]

		m := p.Ronda.Manojo(jugadorStr)
		if m == nil {
			// segundo intento
			m = p.Ronda.Manojo(strings.Title(jugadorStr))
			if m == nil {
				return nil, fmt.Errorf("usuario %s no encontrado", jugadorStr)
			}
		}

		c, err := ParseCarta(valorStr, paloStr)
		if err != nil {
			return nil, err
		}

		jugada = TirarCarta{
			JID:   m.Jugador.ID,
			Carta: *c,
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

func (p *Partida) byeBye() []enco.Envelope {
	if !p.Verbose {
		return nil
	}

	pkts2 := make([]enco.Envelope, 0)

	if p.Terminada() {

		var s string

		if p.Ronda.Manojos[0].Jugador.Equipo == p.ElQueVaGanando() {
			s = p.Ronda.Manojos[0].Jugador.ID
		} else {
			s = p.Ronda.Manojos[1].Jugador.ID
		}

		pkts2 = append(pkts2, enco.Env(
			enco.ALL,
			enco.ByeBye(s),
		))
	}

	return pkts2
}

// Cmd nexo capa presentacion con capa logica
func (p *Partida) Cmd(cmd string) ([]enco.Envelope, error) {

	// checkeo semantico
	jugada, err := ParseJugada(p, cmd)
	if err != nil {
		return nil, err
	}

	pkts2 := jugada.Hacer(p)

	if p.Terminada() {
		// AGREGARLE LoS PKTS ACA
		pkts2 = append(pkts2, p.byeBye()...)
	}

	return pkts2, nil
}
