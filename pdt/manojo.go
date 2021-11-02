package pdt

import (
	"fmt"
	"sort"
)

// cantCartasManojo constante trivial
// del numero de cartas de un manojo
const cantCartasManojo = 3

// Manojo :
type Manojo struct {
	SeFueAlMazo bool                     `json:"seFueAlMazo"`
	Cartas      [cantCartasManojo]*Carta `json:"cartas"`
	// si true->la tiene; si false->ya la tiro
	Tiradas      [cantCartasManojo]bool `json:"tiradas"`
	UltimaTirada int                    `json:"ultimaTirada"`
	Jugador      *Jugador               `json:"jugador"`
}

// GetCantCartasTiradas retorna la cantidad de cartas que ya tiro
func (manojo Manojo) GetCantCartasTiradas() int {
	totalTiradas := 0
	for _, tirada := range manojo.Tiradas {
		if tirada {
			totalTiradas++
		}
	}
	return totalTiradas
}

// retorna true si ya tiro carta en la `mano`
func (manojo Manojo) yaTiroCarta(mano NumMano) bool {
	// idea:
	// si se esta jugando la 1era mano -> debe haber tirado exactamente 1 carta
	// si se esta jugando la 2da mano -> debe haber tirado exactamente 2 cartas
	// si se esta jugando la 3era mano -> debe haber tirado exactamente 3 carta
	cantCartasTiradas := manojo.GetCantCartasTiradas()
	switch mano {
	case Primera:
		return cantCartasTiradas == 1
	case Segunda:
		return cantCartasTiradas == 2
	default: // Tercera
		return cantCartasTiradas == 3
	}
}

// GetCartaIdx ..
// retorna el indice de la `carta`
// error si ni siquiera tiene esa carta
func (manojo Manojo) GetCartaIdx(carta Carta) (int, error) {
	var idx int
	for idx = 0; idx < cantCartasManojo; idx++ {
		esEsa := *manojo.Cartas[idx] == carta
		if esEsa {
			break
		}
	}
	noTieneEsaCarta := idx == cantCartasManojo
	if noTieneEsaCarta {
		return -1, fmt.Errorf("esa carta no se encuentra en este manojo")
	}
	return idx, nil
}

// DEPRECTED ?
// retorna true si todavia no jugo esa carta
// error si ni siquiera tiene esa carta
// func (manojo Manojo) todaviaNoJugo(carta Carta) (bool, error) {
// 	idx, err := manojo.getCartaIdx(carta)
// 	if err != nil {
// 		return false, err
// 	}
// 	return manojo.cartasNoJugadas[idx], nil
// }

// TieneFlor .
// tieneFlor devuelve true si el jugador tiene flor
// Y ademas, si tiene devuelve que tipo de flor: I, II o III
//
// tiene flor sii (CASO I || CASO II || CASO II)
// donde:
// CASO I		~ al menos dos piezas,
// CASO II  ~	tres cartas del mismo palo,
// CASO III ~ una pieza y dos cartas del mismo palo.
func (manojo Manojo) TieneFlor(muestra Carta) (res bool, CASO int) {
	// CASO I: (al menos) dos piezas
	numPiezas := 0
	// en caso de que tenga al menos una pieza,
	// esta variable guarda uno su indice (usado en el caso III)
	piezaIdx := -1
	for i, carta := range manojo.Cartas {
		if carta.esPieza(muestra) {
			numPiezas++
			piezaIdx = i
		}
	}
	if numPiezas >= 2 {
		return true, 1
	}

	// CASO II: tres cartas del mismo palo
	if manojo.Cartas[0].Palo == manojo.Cartas[1].Palo && manojo.Cartas[1].Palo == manojo.Cartas[2].Palo {
		return true, 2
	}

	// CASO II una pieza y dos cartas del mismo palo
	// Y ESAS DOS DIFERENTES DE LA PIEZA (piezaIdx)!
	tieneDosDelMismoPalo :=
		(manojo.Cartas[0].Palo == manojo.Cartas[1].Palo && piezaIdx == 2) ||
			(manojo.Cartas[0].Palo == manojo.Cartas[2].Palo && piezaIdx == 1) ||
			(manojo.Cartas[1].Palo == manojo.Cartas[2].Palo && piezaIdx == 0)

	if numPiezas == 1 && (tieneDosDelMismoPalo) {
		return true, 3
	}

	// si llego hasta aqui -> no tiene flor
	return false, -1
}

// retorna el valor de la flor de un manojo
// si no tiene flor retorna 0 y error
func (manojo *Manojo) calcFlor(muestra Carta) (int, error) {
	var (
		puntajeFlor         int
		tieneFlor, tipoFlor = manojo.TieneFlor(muestra)
	)

	if !tieneFlor {
		return -1, fmt.Errorf("este manojo no tiene flor")
	}

	switch tipoFlor {
	// CASO I: (al menos) dos piezas
	case 1:
		max := maxOf3(manojo.Cartas)
		for _, carta := range manojo.Cartas {
			puntaje := carta.calcPuntaje(muestra)
			if puntaje == max {
				puntajeFlor += puntaje
			} else {
				puntajeFlor += puntaje % 10 // ultimo digito
			}
		}
	// CASO II una pieza y dos cartas del mismo palo;
	// CASO III: tres cartas del mismo palo,
	case 2, 3:
		for _, carta := range manojo.Cartas {
			puntajeFlor += carta.calcPuntaje(muestra)
		}
	}
	return puntajeFlor, nil
}

// Print imprime la info del manojo
func (manojo Manojo) Print() {
	for i := range manojo.Cartas {
		fmt.Printf("    - %s\n", manojo.Cartas[i].String())
	}
}

// tiene2DelMismoPalo devuelve 'true' si tiene dos cartas
// del mismo palo, y ademas los indices de las mismas en
// el array manojo.Cartas
func (manojo Manojo) tiene2DelMismoPalo() (bool, []int) {
	for i := 0; i < cantCartasManojo; i++ {
		for j := i + 1; j < cantCartasManojo; j++ {
			mismoPalo := manojo.Cartas[i].Palo == manojo.Cartas[j].Palo
			if mismoPalo {
				return true, []int{i, j}
			}
		}
	}
	return false, nil
}

// CalcularEnvido devuelve el puntaje correspondiente al envido del manojo
// PRE: no tiene flor
func (manojo Manojo) CalcularEnvido(muestra Carta) (puntajeEnvido int) {
	tiene2DelMismoPalo, idxs := manojo.tiene2DelMismoPalo()
	if tiene2DelMismoPalo {
		x := manojo.Cartas[idxs[0]].calcPuntaje(muestra)
		y := manojo.Cartas[idxs[1]].calcPuntaje(muestra)
		noTieneNingunaPieza := max(x, y) < 27
		if noTieneNingunaPieza {
			puntajeEnvido = x + y + 20
		} else {
			puntajeEnvido = x + y
		}
	} else {
		// si no: simplemente sumo las 2 de mayor valor
		copia := make([]*Carta, cantCartasManojo)
		copy(copia, manojo.Cartas[:])
		// ordeno el array en forma desc de su puntaje
		sort.Slice(copia, func(i, j int) bool {
			return copia[i].calcPuntaje(muestra) > copia[j].calcPuntaje(muestra)
		})
		puntajeEnvido = copia[0].calcPuntaje(muestra) + copia[1].calcPuntaje(muestra)
	}
	return puntajeEnvido
}
