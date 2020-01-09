package ilustrador

import (
	"strconv"
	"strings"
)

func raw(obj string) string {
	return strings.Trim(obj, "\n ")
}

func replace(this, that, here string) string {
	return strings.Replace(here, this, that, 1)
}

func templateMarco() string {
	marco := raw(`
╔══════════════════════════════╗
║                              ║
║                              ║
║                              ║
║                              ║
║                              ║
║                              ║
║                              ║
╚══════════════════════════════╝
`)
	return marco
}

func templateCarta(valor int, palo string) string {
	carta := "┌xx┐" + "\n"
	carta += "│PP│" + "\n"
	carta += "└──┘"

	// numero
	numStr := strconv.Itoa(valor)
	if valor <= 9 {
		numStr = "─" + numStr
	}
	carta = replace("xx", numStr, carta)
	carta = replace("PP", palo[:2], carta)

	return carta
}

func templateCartaDoble(valor int, palo string) string {
	cartaDoble := raw(`
┌xx┐┐
│PP││
└──┘┘
`)

	// numero
	numStr := strconv.Itoa(valor)
	if valor <= 9 {
		numStr = numStr + "─"
	}
	cartaDoble = replace("xx", numStr, cartaDoble)
	cartaDoble = replace("PP", palo[:2], cartaDoble)

	return cartaDoble
}

func templateCartaTriple(valor int, palo string) string {
	cartaTriple := raw(`
┌xx┐┐┐
│PP│││
└──┘┘┘
`)

	// numero
	numStr := strconv.Itoa(valor)
	if valor <= 9 {
		numStr = numStr + "─"
	}
	cartaTriple = replace("xx", numStr, cartaTriple)
	cartaTriple = replace("PP", palo[:2], cartaTriple)

	return cartaTriple
}
