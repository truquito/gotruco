package ilustrador

import (
	"strconv"

	"github.com/jpfilevich/canvas"
)

func templateMarco() string {
	marco := canvas.Raw(`
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
		numStr = numStr + "─"
	}
	carta = canvas.Replace("xx", numStr, carta)
	carta = canvas.Replace("PP", palo[:2], carta)

	return carta
}

func templateCartaDoble(valor int, palo string) string {
	cartaDoble := canvas.Raw(`
┌xx┐┐
│PP││
└──┘┘
`)

	// numero
	numStr := strconv.Itoa(valor)
	if valor <= 9 {
		numStr = numStr + "─"
	}
	cartaDoble = canvas.Replace("xx", numStr, cartaDoble)
	cartaDoble = canvas.Replace("PP", palo[:2], cartaDoble)

	return cartaDoble
}

func templateCartaTriple(valor int, palo string) string {
	cartaTriple := canvas.Raw(`
┌xx┐┐┐
│PP│││
└──┘┘┘
`)

	// numero
	numStr := strconv.Itoa(valor)
	if valor <= 9 {
		numStr = numStr + "─"
	}
	cartaTriple = canvas.Replace("xx", numStr, cartaTriple)
	cartaTriple = canvas.Replace("PP", palo[:2], cartaTriple)

	return cartaTriple
}
