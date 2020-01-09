package ilustrador

import (
	"math"
	"strings"
	"unicode/utf8"
)

type pos struct{ x, y int }

type rectangle struct{ from, to pos }

func calcWidth(obj string) int {
	lines := strings.Split(obj, "\n")
	max := 0
	for _, line := range lines {
		l := utf8.RuneCountInString(line)
		if l > max {
			max = l
		}
	}
	return max
}

// le agrega el paddingLeft
func (r rectangle) center(obj string) string {
	res := ""
	// ############ rectangulo
	// -----		objeto
	// ###-----%### objeto centrado
	rectanguloWidth := r.to.x - r.from.x + 1
	objWidth := calcWidth(obj)
	restante := rectanguloWidth - objWidth
	paddingLeft := float64(restante) / 2.0
	PLredondeado := int(math.Trunc(paddingLeft))
	renderedPadding := ""
	for i := 0; i < PLredondeado; i++ {
		renderedPadding += " "
	}
	for j, letter := range obj {
		if j == 0 {
			res += renderedPadding
			res += string(letter)
		} else if obj[j] == '\n' {
			res += string(letter)
			res += renderedPadding
		} else {
			res += string(letter)
		}
	}
	return res
}

func (r rectangle) right(obj string) string {
	res := ""
	// ############ rectangulo
	// -----		objeto
	// #######----- objeto centrado
	rectanguloWidth := r.to.x - r.from.x + 1
	objWidth := calcWidth(obj)
	paddingLeft := rectanguloWidth - objWidth
	renderedPadding := ""
	for i := 0; i < paddingLeft; i++ {
		renderedPadding += " "
	}
	for j, letter := range obj {
		if j == 0 {
			res += renderedPadding
			res += string(letter)
		} else if obj[j] == '\n' {
			res += string(letter)
			res += renderedPadding
		} else {
			res += string(letter)
		}
	}
	return res
}

func (r rectangle) left(obj string) string {
	res := ""
	// ############ rectangulo
	// -----		objeto
	// #######----- objeto centrado
	rectanguloWidth := r.to.x - r.from.x + 1
	objWidth := calcWidth(obj)
	paddingRight := rectanguloWidth - objWidth
	renderedPadding := ""
	for i := 0; i < paddingRight; i++ {
		renderedPadding += " "
	}
	for j, letter := range obj {
		if j == len(obj)-1 {
			res += string(letter)
			res += renderedPadding
		} else if obj[j] == '\n' {
			res += renderedPadding
			res += string(letter)
		} else {
			res += string(letter)
		}
	}
	return res
}
