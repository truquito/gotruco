package truco

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMsg(t *testing.T) {

	m0 := Msg{
		Tipo: "Nueva-Partida",
		Cont: ContNuevaRonda{
			Pers: "pers aqui",
		},
	}

	m1 := Msg{
		Tipo: "Nueva-Ronda",
		Cont: ContNuevaRonda{
			Pers: "pers aqui",
		},
	}

	m2 := Msg{
		Tipo: "Error",
		Nota: "Se produjo un error",
	}

	m3 := Msg{
		Tipo: "Sumar-Puntos",
		Nota: "El envido lo gano Alvaro con 32 pts de envido. +2 pts para el equipo Rojo",
		Cont: ContSumPts{
			Pts:    3,
			Equipo: "Rojo",
		},
	}

	m4 := Msg{
		Tipo: "TimeOut",
		Nota: "Roro tardo demasiado en jugar. Mano ganada por Rojo",
	}

	m5 := Msg{
		Tipo: "Tirar-Carta",
		Cont: ContTirarCarta{
			Autor: "Alvaro",
			Carta: Carta{
				Palo:  Basto,
				Valor: 6,
			},
		},
	}

	m6 := Msg{
		Tipo: "Gritar-Truco",
		Cont: ContAutor{Autor: "Alvaro"},
	}

	fmt.Println(m0.ToJSON())
	fmt.Println(m1.ToJSON())
	fmt.Println(m2.ToJSON())
	fmt.Println(m3.ToJSON())
	fmt.Println(m4.ToJSON())
	fmt.Println(m5.ToJSON())
	fmt.Println(m6.ToJSON())

	var buff *bytes.Buffer = new(bytes.Buffer)
	// var err error

	// write
	pkt1 := &Pkt{
		Dest: []string{"ALL"},
		Msg: m1,
	}

	pkt2 := &Pkt{
		Dest: []string{"ALL"},
		Msg: m2,
	}

	pkt3 := &Pkt{
		Dest: []string{"ALL"},
		Msg: m3,
	}

	write(buff, pkt1)
	write(buff, pkt2)
	write(buff, pkt3)

	Consume(buff)
}
