package truco

import (
	"strconv"

	"github.com/jpfilevich/canvas"
)

type templates struct{}

func (t templates) renderValorCarta(valor int) string {
	numStr := strconv.Itoa(valor)
	if valor <= 9 {
		numStr = numStr + "─"
	}
	return numStr
}

func (t templates) marco() string {
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

func (t templates) estadisticas() string {
	marco := canvas.Raw(`
╔════════════════╗
│ #Mano:         │
╠────────────────╣
│ Mano:          │
╠────────────────╣
│ Turno:         │
╠────────────────╣
│ Puntuacion:    │
╚════════════════╝
 ╔──────┬──────╗
 │ ROJO │ AZUL │
 ├──────┼──────┤
 │      │      │
 ╚──────┴──────╝
`)
	return marco
}

func (t templates) vacio() string {
	return ""
}

func (t templates) carta(carta Carta) string {
	valor, palo := carta.Valor, carta.Palo.String()
	template := "┌xx┐" + "\n"
	template += "│PP│" + "\n"
	template += "└──┘"

	numStr := t.renderValorCarta(valor)
	template = canvas.Replace("xx", numStr, template)
	template = canvas.Replace("PP", palo[:2], template)

	return template
}

func (t templates) cartaDobleSolapada(carta Carta) string {
	valor, palo := carta.Valor, carta.Palo.String()
	cartaDobleSolapada := canvas.Raw(`
┌xx┐┐
│PP││
└──┘┘
`)

	numStr := t.renderValorCarta(valor)
	cartaDobleSolapada = canvas.Replace("xx", numStr, cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("PP", palo[:2], cartaDobleSolapada)

	return cartaDobleSolapada
}

func (t templates) cartaTripleSolapada(carta Carta) string {
	valor, palo := carta.Valor, carta.Palo.String()
	cartaTripleSolapada := canvas.Raw(`
┌xx┐┐┐
│PP│││
└──┘┘┘
`)

	numStr := t.renderValorCarta(valor)
	cartaTripleSolapada = canvas.Replace("xx", numStr, cartaTripleSolapada)
	cartaTripleSolapada = canvas.Replace("PP", palo[:2], cartaTripleSolapada)

	return cartaTripleSolapada
}

func (t templates) cartaDobleVisible(cartas []Carta) string {
	cartaDobleSolapada := canvas.Raw(`
┌xx┐yy┐
│PP│QQ│
└──┘──┘
`)

	valor, palo := cartas[0].Valor, cartas[0].Palo.String()
	numStr := t.renderValorCarta(valor)
	cartaDobleSolapada = canvas.Replace("xx", numStr, cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("PP", palo[:2], cartaDobleSolapada)

	valor, palo = cartas[1].Valor, cartas[1].Palo.String()
	numStr = t.renderValorCarta(valor)
	cartaDobleSolapada = canvas.Replace("yy", numStr, cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("QQ", palo[:2], cartaDobleSolapada)

	return cartaDobleSolapada
}

func (t templates) cartaTripleVisible(cartas []Carta) string {
	cartaDobleSolapada := canvas.Raw(`
┌xx┐yy┐zz┐
│PP│QQ│RR│
└──┘──┘──┘
`)

	valor, palo := cartas[0].Valor, cartas[0].Palo.String()
	numStr := t.renderValorCarta(valor)
	cartaDobleSolapada = canvas.Replace("xx", numStr, cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("PP", palo[:2], cartaDobleSolapada)

	valor, palo = cartas[1].Valor, cartas[1].Palo.String()
	numStr = t.renderValorCarta(valor)
	cartaDobleSolapada = canvas.Replace("yy", numStr, cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("QQ", palo[:2], cartaDobleSolapada)

	valor, palo = cartas[2].Valor, cartas[2].Palo.String()
	numStr = t.renderValorCarta(valor)
	cartaDobleSolapada = canvas.Replace("zz", numStr, cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("RR", palo[:2], cartaDobleSolapada)

	return cartaDobleSolapada
}
