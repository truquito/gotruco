package main

import (
	"encoding/json"
	"fmt"

	"github.com/truquito/gotruco/pdt"
)

func main() {
	m := pdt.Manojo{
		Cartas: [3]*pdt.Carta{
			{Palo: pdt.Oro, Valor: 0},
			{Palo: pdt.Oro, Valor: 0},
			{Palo: pdt.Oro, Valor: 0},
		},
	}

	stats := struct {
		Total      int
		FloresDist map[int]int
		EnvidoDist map[int]int
		PoderDist  map[int]int
	}{
		Total:      0,
		FloresDist: make(map[int]int),
		EnvidoDist: make(map[int]int),
		PoderDist:  make(map[int]int),
	}

	for i := 0; i < 40; i++ {
		muestra := pdt.NuevaCarta(pdt.CartaID(i))
		for j := 0; j < 40; j++ {
			for k := 0; k < 40; k++ {
				for l := 0; l < 40; l++ {
					if j == i || k == i || l == i {
						continue
					}

					stats.Total++

					c1 := pdt.NuevaCarta(pdt.CartaID(j))
					c2 := pdt.NuevaCarta(pdt.CartaID(k))
					c3 := pdt.NuevaCarta(pdt.CartaID(l))

					m.Cartas[0] = &c1
					m.Cartas[1] = &c2
					m.Cartas[2] = &c3

					tieneFlor, _ := m.TieneFlor(muestra)

					if tieneFlor {
						ptsFlor, _ := m.CalcFlor(muestra)
						stats.FloresDist[ptsFlor]++
					} else {
						ptsEnvido := m.CalcularEnvido(muestra)
						stats.EnvidoDist[ptsEnvido]++
					}

					poder := c1.CalcPoder(muestra) + c2.CalcPoder(muestra) + c3.CalcPoder(muestra)
					stats.PoderDist[poder]++
				}
			}
		}
	}

	bs, err := json.Marshal(stats)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}
