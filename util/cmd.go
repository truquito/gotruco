package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/filevich/truco/pdt"
)

// regexps
var (
	regexps = map[string]*regexp.Regexp{
		"jugadaSimple": regexp.MustCompile(`(?i)^([a-zA-Z0-9_-]+) ([a-zA-Z0-9_-]+)$`),
		"jugadaTirada": regexp.MustCompile(`(?i)^([a-zA-Z0-9_-]+) (1|2|3|4|5|6|7|10|11|12) (oro|copa|basto|espada)$`),
	}
)

func ParseJugada(p *pdt.PartidaDT, cmd string) (pdt.IJugada, error) {
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

func checkeoSemantico(p *pdt.PartidaDT, cmd string) (pdt.IJugada, error) {

	var jugada pdt.IJugada

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
			jugada = pdt.TocarEnvido{Manojo: m}
		case "real-envido":
			jugada = pdt.TocarRealEnvido{Manojo: m}
		case "falta-envido":
			jugada = pdt.TocarFaltaEnvido{Manojo: m}

		// cantos
		case "flor":
			jugada = pdt.CantarFlor{Manojo: m}
		case "contra-flor":
			jugada = pdt.CantarContraFlor{Manojo: m}
		case "contra-flor-al-resto":
			jugada = pdt.CantarContraFlorAlResto{Manojo: m}

		// gritos
		case "truco":
			jugada = pdt.GritarTruco{Manojo: m}
		case "re-truco":
			jugada = pdt.GritarReTruco{Manojo: m}
		case "vale-4":
			jugada = pdt.GritarVale4{Manojo: m}

		// respuestas
		case "quiero":
			jugada = pdt.ResponderQuiero{Manojo: m}
		case "no-quiero":
			jugada = pdt.ResponderNoQuiero{Manojo: m}
		// case "tiene":
		// 	jugada = responderNoQuiero{Manojo: m}

		// acciones
		case "mazo":
			jugada = pdt.IrseAlMazo{Manojo: m}
		case "tirar":
			jugada = pdt.IrseAlMazo{Manojo: m}
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

		c, err := pdt.ParseCarta(valorStr, paloStr)
		if err != nil {
			return nil, err
		}

		jugada = pdt.TirarCarta{
			Manojo: m,
			Carta:  *c,
		}
	}

	return jugada, nil
}
