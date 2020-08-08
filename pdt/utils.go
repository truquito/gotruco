package pdt

import (
	"reflect"
)

// Contains .
func Contains(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

// Eliminar .
func Eliminar(manojos []*Manojo, manojo *Manojo) []*Manojo {
	var i int
	// primero encuentro el elemento
	for i = 0; i <= len(manojos); i++ {
		noLoContiene := i == len(manojos)
		if noLoContiene {
			return manojos
		}
		if manojos[i] == manojo {
			break
		}
	}
	manojos[i] = manojos[len(manojos)-1] // Copy last element to index i.
	return manojos[:len(manojos)-1]      // Truncate slice.
}

func maxOf3(cartas [3]*Carta) int {
	max := 0
	for _, carta := range cartas {
		if carta.Valor > max {
			max = int(carta.Valor)
		}
	}
	return max
}

func max(nums ...int) int {
	max := nums[0]
	for _, x := range nums[1:] {
		if x > max {
			max = x
		}
	}
	return max
}

// ver documentacion cambio de variable
func cv(x, mano JugadorIdx, cantJugadores int) (y JugadorIdx) {
	if x >= mano {
		y = x - mano
	} else {
		c := JugadorIdx(cantJugadores) - mano
		y = x + c
	}
	return y
}
