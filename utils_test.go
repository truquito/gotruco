package truco

import (
	// "fmt"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func maxGENERIC(foo func(x interface{}) int, nums ...(interface{})) {
	total := 0
	for _, num := range nums {
		total += foo(num)
	}
}

var foo = func(x interface{}) int {
	return x.(Person).Age
}

func BenchmarkMaxNew(b *testing.B) {
	persons := []Person{
		Person{"Joan", 32},
		Person{"Marie", 29},
	}
	for i := 0; i < b.N; i++ {
		// adaptacion
		adaptacion := make([]interface{}, len(persons))
		for i, v := range persons {
			adaptacion[i] = v
		}
		maxGENERIC(foo, adaptacion...)
	}
}

func BenchmarkMaxOld(b *testing.B) {
	cartas := [3]Carta{
		Carta{Palo: Oro, Valor: 3},
		Carta{Palo: Copa, Valor: 4},
		Carta{Palo: Basto, Valor: 5},
	}
	// muestra := Carta{ Espada, 10}
	for i := 0; i < b.N; i++ {
		maxOf3(cartas)
	}
}
