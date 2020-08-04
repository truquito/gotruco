package truco

import (
	"fmt"
	"testing"
)

func TestMsg(t *testing.T) {

	m1 := Msg2{
		Tipo: "Nueva-Ronda",
		Nota: "",
		Cont: ContNuevaRonda{
			Pers: "pers aqui",
		},
	}

	m2 := Msg2{
		Tipo: "Error",
		Nota: "Se produjo un error",
	}

	m3 := Msg2{
		Tipo: "Sumar-Puntos",
		Nota: "El envido lo gano Alvaro con 32 pts de envido. +2 pts para el equipo Rojo",
		Cont: ContSumPts{
			Pts:    3,
			Equipo: "Rojo",
		},
	}

	m4 := Msg2{
		Tipo: "TimeOut",
		Nota: "Roro tardo demasiado en jugar. Mano ganada por Rojo",
	}

	m5 := Msg2{
		Tipo: "Tirar-Carta",
		Cont: ContTirarCarta{
			Autor: "Alvaro",
			Carta: Carta{
				Palo:  Basto,
				Valor: 6,
			},
		},
	}

	m6 := Msg2{
		Tipo: "Gritar-Truco",
		Cont: ContAutor{Autor: "Alvaro"},
	}

	fmt.Println(m1.ToJSON())
	fmt.Println(m2.ToJSON())
	fmt.Println(m3.ToJSON())
	fmt.Println(m4.ToJSON())
	fmt.Println(m5.ToJSON())
	fmt.Println(m6.ToJSON())
}
