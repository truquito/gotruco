package main

import (
	"math/rand"
	"strconv"
	"time"
)

/*
 * Barajas; orden absoluto:
 *  ------------------------------------------------------------------------------------------
 * | ID	| Carta						|	ID | 	Carta						|	ID | 	Carta							|	ID | 	Carta			|
 * |------------------------------------------------------------------------------------------|
 * | 00 |	(1,Basto)				|	10 | 	(1,Copa)				|	20 | 	(1,Espada)				|	30 | 	(1,Oro) 	|
 * | 01 |	(2,Basto)				|	11 |	(2,Copa)				|	21 |	(2,Espada)				|	31 |	(2,Oro) 	|
 * | 02 |	(3,Basto)				|	12 |	(3,Copa)				|	22 |	(3,Espada)				|	32 |	(3,Oro) 	|
 * | 03 |	(4,Basto)				|	13 |	(4,Copa)				|	23 |	(4,Espada)				|	33 |	(4,Oro) 	|
 * | 04 |	(5,Basto)				|	14 |	(5,Copa)				|	24 |	(5,Espada)				|	34 |	(5,Oro) 	|
 * | 05 |	(6,Basto)				|	15 |	(6,Copa)				|	25 |	(6,Espada)				|	35 |	(6,Oro) 	|
 * | 06 |	(7,Basto)				|	16 |	(7,Copa)				|	26 |	(7,Espada)				|	36 |	(7,Oro) 	|
 *  ------------------------------------------------------------------------------------------
 * | 07 |	(10,Basto)			|	17 |	(10,Copa)				|	27 |	(10,Espada)				|	37 |	(10,Oro) 	|
 * | 08 |	(11,Basto)			|	18 |	(11,Copa)				|	28 |	(11,Espada)				|	38 |	(11,Oro) 	|
 * | 09 |	(12,Basto)			|	19 |	(12,Copa)				|	29 |	(12,Espada)				|	39 |	(12,Oro)	|
 *  ------------------------------------------------------------------------------------------
 */

// Palo enum
type Palo int

// 4 palos
const (
	Basto  Palo = 0 // [00..09]
	Copa   Palo = 1 // [10..19]
	Espada Palo = 2 // [20..29]
	Oro    Palo = 3 // [30..39]
)

const (
	minCartaID = 0
	maxCartaID = 40
)

// toString
func (p Palo) String() string {
	nombres := [...]string{
		"Basto",
		"Copa",
		"Espada",
		"Oro"}

	if p < Basto || p > Oro {
		return "Unknown"
	}

	return nombres[p]
}

// Carta struct
type Carta struct {
	Palo  Palo
	Valor int
}

// devulve `true` si la carta es pieza
// segun la el parametro `muestra`
func (c Carta) esPieza(muestra Carta) bool {

	// es pieza sii: (CASO I || CASO II)
	// donde,
	// CASO I: es (2|4|5|10|11) & es de la muestra
	// CASO II: es 12 de la muestra & la muestra es (2|4|5|10|11)

	// CASO I:
	esNumericamentePieza := contains([]int{2, 4, 5, 10, 11}, c.Valor)
	esDeLaMuestra := c.Palo == muestra.Palo
	esPiezaCasoI := esNumericamentePieza && esDeLaMuestra

	// CASO II:
	esDoce := c.Valor == 12
	esPiezaCasoII := esDoce && esDeLaMuestra

	return esPiezaCasoI || esPiezaCasoII
}

// Devuelve el puntaje
// todo: resolver esto
// ojo con Puntaje(7,Espada) == Puntaje(7,Oro) PERO
// (7,Espada) LE GANA A (7,Oro) !!
// detalle, se podria reducir la logica booleana,
// pero asi queda simple & natural a la vista
func (c Carta) calcPuntaje(muestra Carta) int {
	var puntaje int

	// Piezas
	if c.esPieza(muestra) {
		switch c.Valor {
		case 2:
			puntaje = 30
		case 4:
			puntaje = 29
		case 5:
			puntaje = 28
		case 11, 10:
			puntaje = 27
		case 12:
			puntaje = (Carta{Palo: c.Palo, Valor: muestra.Valor}).calcPuntaje(muestra)
		}

		// Matas
	} else if (c.Palo == Espada || c.Palo == Basto) && (c.Valor == 1) {
		puntaje = 1
	} else if (c.Palo == Espada || c.Palo == Basto) && (c.Valor == 7) {
		puntaje = 7

		// Chicas
	} else if c.Valor <= 3 {
		puntaje = c.Valor

		// Comunes
	} else if 10 <= c.Valor && c.Valor <= 12 {
		puntaje = 0
	} else if 4 <= c.Valor && c.Valor <= 7 {
		puntaje = c.Valor
	}

	return puntaje
}

func (c Carta) toString() string {
	return strconv.Itoa(c.Valor) + " de " + c.Palo.String()
}

/*
	Cada palo tiene 10 cartas
	#{1, 2, 3, 4, 5, 6, 7, 10, 11, 12} = 10
*/
const tamanoPalo = 10

// Basto, Copa, Espada, Oro
const cantPalos = 4

// CartaID abstraccion *no-relativa* de carta
type CartaID int

func (i CartaID) getPalo() Palo {
	var palo Palo

	// [00..09]
	if 0 <= i && i <= 9 {
		palo = Basto
		// [10..19]
	} else if 10 <= i && i <= 19 {
		palo = Copa
		// [20..29]
	} else if 20 <= i && i <= 29 {
		palo = Espada
		// [30..39]
	} else /* if 30 <= i && i <= 39 */ {
		palo = Oro
	}

	return palo
}

func (i CartaID) getValor() int {
	var valor int
	ultimoDigito := i % 10
	if ultimoDigito <= 6 {
		valor = int(ultimoDigito) + 1
	} else /* if ultimoDigito >= 7 */ {
		valor = 10 + int(ultimoDigito) - 7
	}

	return valor
}

// Devuelve la `Carta` correspondiente al ID i
func newCarta(i CartaID) Carta {
	return Carta{
		Palo:  i.getPalo(),
		Valor: i.getValor(),
	}
}

/*
 * Devuelve un array de tamano n de CartaID sin repetir
 */
func getCartasRandom(n int) []int {
	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(maxCartaID)
	randomSample := make([]int, n)

	for i, r := range p[:n] {
		randomSample[i] = r
	}

	return randomSample
}

/*
 * Reparte 3 cartas al azar a cada manojo de c/jugador
 * y 1 a la `muestra` (se las actualiza)
 */
func dealCards(manojos *[]Manojo, muestra *Carta) {
	cantJugadores := cap(*manojos)
	// genero `3*cantJugadores + 1` cartas al azar
	randomCards := getCartasRandom(3*cantJugadores + 1)

	for numJugador := 0; numJugador < cantJugadores; numJugador++ {
		for numCarta := 0; numCarta < 3; numCarta++ {
			cartaID := CartaID(randomCards[3*numJugador+numCarta])
			carta := newCarta(cartaID)
			(*manojos)[numJugador].Cartas[numCarta] = carta
		}
	}

	// la ultima es la muestra
	n := cap(randomCards)
	*muestra = newCarta(CartaID(randomCards[n-1]))
}
