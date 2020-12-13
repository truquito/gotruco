package pdt

import (
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

var castAsPerson = func(x interface{}) int {
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
		maxGENERIC(castAsPerson, adaptacion...)
	}
}

func BenchmarkMaxOld(b *testing.B) {
	cartas := [3]*Carta{
		&Carta{Palo: Oro, Valor: 3},
		&Carta{Palo: Copa, Valor: 4},
		&Carta{Palo: Basto, Valor: 5},
	}
	// muestra := Carta{ Espada, 10}
	for i := 0; i < b.N; i++ {
		maxOf3(cartas)
	}
}

// func TestEliminarAlguienQueNoEsta(t *testing.T) {
// 	p, _ := pdt.NuevaPartida(A20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})

// 	manojos := make([]*Manojo, 5)
// 	for i := range p.Ronda.Manojos {
// 		if i == 5 {
// 			break
// 		}
// 		manojos[i] = &p.Ronda.Manojos[i]
// 	}

// 	for _, m := range manojos {
// 		fmt.Println(m.Jugador.Nombre)
// 	}

// 	manojoAEliminar := &p.Ronda.Manojos[5]
// 	manojos = Eliminar(manojos, manojoAEliminar)

// 	if !(len(manojos) == 5) {
// 		t.Error("no debio de haber eliminado a nadie")
// 	}

// 	for _, m := range manojos {
// 		fmt.Println(m.Jugador.Nombre)
// 	}

// 	manojoAEliminar = &p.Ronda.Manojos[2]
// 	manojos = Eliminar(manojos, manojoAEliminar)

// 	if !(len(manojos) == 4) {
// 		t.Error("debio de haber eliminado a uno")
// 	}

// 	for _, m := range manojos {
// 		fmt.Println(m.Jugador.Nombre)
// 	}

// }
