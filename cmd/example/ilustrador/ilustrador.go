package ilustrador

import "fmt"

/*
               ▒▒▒▒▒▒▒▒▒▒    ▒▒▒▒▒▒▒▒▒▒
               ▒▒Alvaro▒▒    ▒▒Alvaro▒▒
           ╔══════════════════════════════╗
           ║       ┌12┐┐┐    ┌12┐┐┐       ║
           ║       │Or│││    │Or│││       ║
▒▒▒▒▒▒▒▒▒▒ ║ ┌12┐┐┐└──┘┘┘┌12┐└──┘┘┘┌12┐┐┐ ║ ▒▒▒▒▒▒▒▒▒▒
▒▒▒Nomb▒▒▒ ║ │Or│││      │Or│      │Or│││ ║ ▒▒▒Nom▒▒▒▒
           ║ └──┘┘┘┌12┐┐┐└──┘┌12┐┐┐└──┘┘┘ ║
           ║       │Or│││    │Or│││       ║
           ║       └──┘┘┘    └──┘┘┘       ║
           ╚══════════════════════════════╝
               ▒▒Alvaro▒▒    ▒▒Alvaro▒▒
               ▒▒▒▒▒▒▒▒▒▒    ▒▒▒▒▒▒▒▒▒▒
*/

// Instante de una partida
type Instante struct {
	Jugadores []string
	Turno     int
	Mano      int
}

const (
	width  = 54
	height = 13
)

type lienzo [][]rune

func (lienzo lienzo) setMarco() {
	// marco superior
	for i := 12; i < 12+30; i++ {
		lienzo[i][2] = '═'
	}
	// marco inferior
	for i := 12; i < 12+30; i++ {
		lienzo[i][10] = '═'
	}
	// marco izquierdo
	for i := 3; i < 10; i++ {
		lienzo[11][i] = '║'
	}
	// marco derecho
	for i := 3; i < 10; i++ {
		lienzo[42][i] = '║'
	}
	// esquinitas
	lienzo[11][2] = '╔'
	lienzo[42][2] = '╗'
	lienzo[11][10] = '╚'
	lienzo[42][10] = '╝'
}

func (lienzo lienzo) renderizar() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fmt.Printf("%s", string(lienzo[x][y]))
		}
		fmt.Println()
	}
}

func (lienzo lienzo) draw(fromX, fromY int, obj string) {
	var (
		x = fromX
		y = fromY
	)
	for _, char := range obj {
		if char == '\n' {
			y++
			x = fromX
		} else {
			lienzo[x][y] = char
			x++
		}
	}
}

func (lienzo lienzo) setMuestra(valor int, palo string) {
	// borde de la carta
	lienzo[25][5] = '┌'
	lienzo[28][5] = '┐'

	lienzo[25][6] = '│'
	lienzo[28][6] = '│'

	lienzo[25][7] = '└'
	lienzo[26][7] = '─'
	lienzo[27][7] = '─'
	lienzo[28][7] = '┘'

	// numero
	if valor <= 9 {
		lienzo[26][5] = '─'
		lienzo[27][5] = '2'
	} else {
		lienzo[26][5] = '3'
		lienzo[27][5] = '3'
	}

	// palo
	lienzo[26][6] = rune(palo[0])
	lienzo[27][6] = rune(palo[1])
}

func (lienzo lienzo) setMuestra2(valor int, palo string) {
	carta := "┌14┐" + "\n"
	carta += "│Oz│" + "\n"
	carta += "└──┘" + "\n"

	lienzo.draw(25, 5, carta)
}

func nuevoLienzo() lienzo {
	var lienzo lienzo
	lienzo = make([][]rune, width)
	for i := range lienzo {
		lienzo[i] = make([]rune, height)
	}
	// cargo con ' '
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			lienzo[x][y] = ' '
		}
	}
	return lienzo
}

// Imprimir un instante
func Imprimir(instante Instante) {
	lienzo := nuevoLienzo()
	lienzo.setMarco()
	// lienzo.setMuestra(36, "Es")
	lienzo.setMuestra2(36, "Es")
	lienzo.draw(27, 9, "HOLA MUNDO")
	lienzo.renderizar()
}
