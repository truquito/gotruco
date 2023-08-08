package pdt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

/*
 *  Barajas; orden absoluto:
 *  ----------------------------------------------------------
 * | ID	| Carta	    ID | Carta	  ID | Carta	    ID | Carta |
 * |---------------------------------------------------------|
 * | 00 | 1,basto   10 | 1,copa   20 | 1,espada   30 | 1,oro |
 * | 01 | 2,basto   11 | 2,copa   21 | 2,espada   31 | 2,oro |
 * | 02 | 3,basto   12 | 3,copa   22 | 3,espada   32 | 3,oro |
 * | 03 | 4,basto   13 | 4,copa   23 | 4,espada   33 | 4,oro |
 * | 04 | 5,basto   14 | 5,copa   24 | 5,espada   34 | 5,oro |
 * | 05 | 6,basto   15 | 6,copa   25 | 6,espada   35 | 6,oro |
 * | 06 | 7,basto   16 | 7,copa   26 | 7,espada   36 | 7,oro |
 *  ----------------------------------------------------------
 * | 07 |10,basto   17 |10,copa   27 |10,espada   37 |10,oro |
 * | 08 |11,basto   18 |11,copa   28 |11,espada   38 |11,oro |
 * | 09 |12,basto   19 |12,copa   29 |12,espada   39 |12,oro |
 *  ----------------------------------------------------------
 */

var primes = [...]int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47,
	53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137,
	139, 149, 151, 157, 163, 167, 173}

const (
	minCartaID = 0
	maxCartaID = 40
)

// Palo enum
type Palo int

// 4 palos
const (
	Basto  Palo = 0 // [00..09]
	Copa   Palo = 1 // [10..19]
	Espada Palo = 2 // [20..29]
	Oro    Palo = 3 // [30..39]
)

var ToPalo = map[string]Palo{
	"basto":  Basto,
	"copa":   Copa,
	"espada": Espada,
	"oro":    Oro,
}

func (p Palo) String() string {
	palos := []string{
		"basto",
		"copa",
		"espada",
		"oro",
	}

	ok := p >= 0 || int(p) < len(ToPalo)
	if !ok {
		return "Unknown"
	}

	return palos[p]
}

// MarshalJSON marshals the enum as a quoted json string
func (p Palo) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(p.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (p *Palo) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*p = ToPalo[j]
	return nil
}

// Carta struct
type Carta struct {
	Palo  Palo `json:"palo"`
	Valor int  `json:"valor"`
}

func (c Carta) ID() CartaID {
	id := 10 * int(c.Palo)
	id += int(c.Valor) - 1

	if c.Valor >= 10 {
		id -= 2
	}

	return CartaID(id)
}

func (c Carta) PUID() int {
	return primes[c.ID()]
}

// func (i CartaID) getPalo() Palo {
// 	var palo Palo

// 	// [00..09]
// 	if 0 <= i && i <= 9 {
// 		palo = Basto
// 		// [10..19]
// 	} else if 10 <= i && i <= 19 {
// 		palo = Copa
// 		// [20..29]
// 	} else if 20 <= i && i <= 29 {
// 		palo = Espada
// 		// [30..39]
// 	} else /* if 30 <= i && i <= 39 */ {
// 		palo = Oro
// 	}

// 	return palo
// }

// func (i CartaID) getValor() int {
// 	var valor int
// 	ultimoDigito := i % 10
// 	if ultimoDigito <= 6 {
// 		valor = int(ultimoDigito) + 1
// 	} else /* if ultimoDigito >= 7 */ {
// 		valor = 10 + int(ultimoDigito) - 7
// 	}

// 	return valor
// }

func (c Carta) esNumericamentePieza() bool {
	return c.Valor == 2 ||
		c.Valor == 4 ||
		c.Valor == 5 ||
		c.Valor == 10 ||
		c.Valor == 11
}

// devulve `true` si la carta es pieza
// segun la el parametro `muestra`
func (c Carta) EsPieza(muestra Carta) bool {

	// es pieza sii: (CASO I || CASO II)
	// donde,
	// CASO I: es (2|4|5|10|11) & es de la muestra
	// CASO II: es 12 de la muestra & la muestra es (2|4|5|10|11)

	// CASO I:
	esDeLaMuestra := c.Palo == muestra.Palo
	esPiezaCasoI := c.esNumericamentePieza() && esDeLaMuestra

	// CASO II:
	esDoce := c.Valor == 12
	esPiezaCasoII := esDoce && esDeLaMuestra && muestra.esNumericamentePieza()

	return esPiezaCasoI || esPiezaCasoII
}

func (c *Carta) EsMata() bool {
	if (c.Palo == Espada || c.Palo == Basto) && (c.Valor == 1) {
		return true

	} else if (c.Palo == Espada || c.Palo == Oro) && (c.Valor == 7) {
		return true
	}

	return false
}

// Devuelve el puntaje
// no confundir el puntaje con "Poder"
// ojo con Puntaje(7,Espada) == Puntaje(7,Oro) PERO
// (7,Espada) LE GANA A (7,Oro) (tiene mas poder)!!
// detalle, se podria reducir la logica booleana,
// pero asi queda simple & natural a la vista
func (c Carta) calcPuntaje(muestra Carta) int {
	var puntaje int

	// Piezas
	if c.EsPieza(muestra) {
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
			valeComo := Carta{Palo: c.Palo, Valor: muestra.Valor}
			puntaje = valeComo.calcPuntaje(muestra)
		}

		// Matas
	} else if c.EsMata() {
		puntaje = c.Valor

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

// guarismo ficticio y abstracto para simplificar
// las comparaciones
func (c Carta) CalcPoder(muestra Carta) int {
	var poder int

	if c.EsPieza(muestra) {
		switch c.Valor {
		case 2:
			poder = 34
		case 4:
			poder = 33
		case 5:
			poder = 32
		case 11:
			poder = 31
		case 10:
			poder = 30
		case 12:
			valeComo := Carta{Palo: c.Palo, Valor: muestra.Valor}
			poder = valeComo.CalcPoder(muestra)
		}

	} else if c.Palo == Espada && c.Valor == 1 {
		poder = 23
	} else if c.Palo == Basto && c.Valor == 1 {
		poder = 22
	} else if c.Palo == Espada && c.Valor == 7 {
		poder = 21
	} else if c.Palo == Oro && c.Valor == 7 {
		poder = 20
		// Chicas
	} else if c.Valor == 3 {
		poder = 19
	} else if c.Valor == 2 {
		poder = 18
	} else if c.Valor == 1 {
		poder = 17
	} else if c.Valor == 12 {
		poder = 16
	} else if c.Valor == 11 {
		poder = 15
	} else if c.Valor == 10 {
		poder = 14
	} else if c.Valor == 7 {
		poder = 13
	} else if c.Valor == 6 {
		poder = 12
	} else if c.Valor == 5 {
		poder = 11
	} else if c.Valor == 4 {
		poder = 10
	}

	return poder
}

// String .
func (c Carta) String() string {
	return strconv.Itoa(c.Valor) + " de " + c.Palo.String()
}

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
func nuevaCarta(i CartaID) Carta {
	return Carta{
		Palo:  i.getPalo(),
		Valor: i.getValor(),
	}
}

// ParseCarta hace todos los checkeos necesarios
func ParseCarta(valorStr, paloStr string) (*Carta, error) {
	var (
		valor int
		palo  Palo
	)
	valor, err := strconv.Atoi(valorStr)
	if err != nil {
		return nil, fmt.Errorf("no se pudo reconocer el valor de la carta")
	}
	// valor in {1,2,3,4,5,6,7,10,11,12} iif
	ok := valor >= 1 && valor <= 12 && valor != 8 && valor != 9

	if !ok {
		return nil, fmt.Errorf("el valor de esa carta es incorrecto")
	}
	paloLower := strings.ToLower(paloStr)

	switch paloLower {
	case "basto":
		palo = Basto
	case "copa":
		palo = Copa
	case "oro":
		palo = Oro
	case "espada":
		palo = Espada
	default:
		return nil, fmt.Errorf("el palo de esa carta es incorrecto")
	}

	return &Carta{palo, valor}, nil
}

/*
 * Devuelve un array de tamano n de CartaID sin repetir
 */
func getCartasRandom(n int) []int {
	rand.Seed(time.Now().UnixNano())
	// obtengo una permutacion random del array [0..39]
	p := rand.Perm(maxCartaID)
	randomSample := make([]int, n)

	copy(randomSample, p[:n])
	// for i, r := range p[:n] {
	// 	randomSample[i] = r
	// }

	return randomSample
}
