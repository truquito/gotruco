package ilustrador

type lienzo [][]rune

const (
	width  = 54
	height = 13
)

func (lienzo lienzo) renderizar() string {
	// el string va a tener el mismo tamano que
	// el lienzo + los \n que no estan
	res := ""
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			res += string(lienzo[x][y])
		}
		res += "\n"
	}
	return res
}

func (lienzo lienzo) drawAt(pos pos, obj string) {
	lienzo.draw(pos.x, pos.y, obj)
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
