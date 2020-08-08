package ptr

import (
	"strconv"

	"github.com/jpfilevich/canvas"
	"github.com/jpfilevich/truco/pdt"
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

func mkCarta(valor, palo string) string {

	template := "┌xx┐" + "\n"
	template += "│PP│" + "\n"
	template += "└──┘"

	template = canvas.Replace("xx", valor, template)
	template = canvas.Replace("PP", palo, template)

	return template
}

func (t templates) carta(carta pdt.Carta) string {
	valor, palo := carta.Valor, carta.Palo.String()
	numStr := t.renderValorCarta(valor)
	return mkCarta(numStr, palo[:2])
}

func (t templates) cartaOculta() string {
	return mkCarta("──", "//")
}

func (t templates) cartaDobleSolapada(carta pdt.Carta) string {
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

func (t templates) cartaTripleSolapada(carta pdt.Carta) string {
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

func mkCartaDobleVisible(valores, palos []string) string {
	cartaDobleSolapada := canvas.Raw(`
┌xx┐yy┐
│PP│QQ│
└──┘──┘
`)

	cartaDobleSolapada = canvas.Replace("xx", valores[0], cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("PP", palos[0], cartaDobleSolapada)

	cartaDobleSolapada = canvas.Replace("yy", valores[1], cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("QQ", palos[1], cartaDobleSolapada)

	return cartaDobleSolapada
}

func (t templates) cartaDobleVisible(cartas []*pdt.Carta) string {

	valor1, palo1 := cartas[0].Valor, cartas[0].Palo.String()
	numStr1 := t.renderValorCarta(valor1)

	valor2, palo2 := cartas[1].Valor, cartas[1].Palo.String()
	numStr2 := t.renderValorCarta(valor2)

	valores := []string{numStr1, numStr2}
	palos := []string{palo1[:2], palo2[:2]}

	return mkCartaDobleVisible(valores, palos)
}

func (t templates) cartaDobleOculta() string {
	valores := []string{"──", "──"}
	palos := []string{"//", "//"}
	return mkCartaDobleVisible(valores, palos)
}

func mkCartaTripleVisible(valores, palos []string) string {
	cartaDobleSolapada := canvas.Raw(`
┌xx┐yy┐zz┐
│PP│QQ│RR│
└──┘──┘──┘
`)
	cartaDobleSolapada = canvas.Replace("xx", valores[0], cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("PP", palos[0], cartaDobleSolapada)

	cartaDobleSolapada = canvas.Replace("yy", valores[1], cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("QQ", palos[1], cartaDobleSolapada)

	cartaDobleSolapada = canvas.Replace("zz", valores[2], cartaDobleSolapada)
	cartaDobleSolapada = canvas.Replace("RR", palos[2], cartaDobleSolapada)

	return cartaDobleSolapada
}

func (t templates) cartaTripleVisible(cartas []*pdt.Carta) string {

	valor1, palo1 := cartas[0].Valor, cartas[0].Palo.String()
	numStr1 := t.renderValorCarta(valor1)

	valor2, palo2 := cartas[1].Valor, cartas[1].Palo.String()
	numStr2 := t.renderValorCarta(valor2)

	valor3, palo3 := cartas[2].Valor, cartas[2].Palo.String()
	numStr3 := t.renderValorCarta(valor3)

	valores := []string{numStr1, numStr2, numStr3}
	palos := []string{palo1[:2], palo2[:2], palo3[:2]}

	return mkCartaTripleVisible(valores, palos)
}

func (t templates) cartaTripleOculta() string {

	valores := []string{"──", "──", "──"}
	palos := []string{"//", "//", "//"}

	return mkCartaTripleVisible(valores, palos)
}
