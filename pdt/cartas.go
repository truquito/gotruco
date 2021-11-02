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
 * | ID	| Carta	    ID | Carta	  ID | Carta	  ID | Carta |
 * |---------------------------------------------------------|
 * | 00 | 1,Basto   10 | 1,Copa   20 | 1,Espada   30 | 1,Oro |
 * | 01 | 2,Basto   11 | 2,Copa   21 | 2,Espada   31 | 2,Oro |
 * | 02 | 3,Basto   12 | 3,Copa   22 | 3,Espada   32 | 3,Oro |
 * | 03 | 4,Basto   13 | 4,Copa   23 | 4,Espada   33 | 4,Oro |
 * | 04 | 5,Basto   14 | 5,Copa   24 | 5,Espada   34 | 5,Oro |
 * | 05 | 6,Basto   15 | 6,Copa   25 | 6,Espada   35 | 6,Oro |
 * | 06 | 7,Basto   16 | 7,Copa   26 | 7,Espada   36 | 7,Oro |
 *  ----------------------------------------------------------
 * | 07 |10,Basto   17 |10,Copa   27 |10,Espada   37 |10,Oro |
 * | 08 |11,Basto   18 |11,Copa   28 |11,Espada   38 |11,Oro |
 * | 09 |12,Basto   19 |12,Copa   29 |12,Espada   39 |12,Oro |
 *  ----------------------------------------------------------
 */

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

var toPalo = map[string]Palo{
	"Basto":  Basto,
	"Copa":   Copa,
	"Espada": Espada,
	"Oro":    Oro,
}

func (p Palo) String() string {
	palos := []string{
		"Basto",
		"Copa",
		"Espada",
		"Oro",
	}

	ok := p >= 0 || int(p) < len(toPalo)
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
	*p = toPalo[j]
	return nil
}

// Carta struct
type Carta struct {
	Palo  Palo `json:"palo"`
	Valor int  `json:"valor"`
}

func (c Carta) esNumericamentePieza() bool {
	return c.Valor == 2 ||
		c.Valor == 4 ||
		c.Valor == 5 ||
		c.Valor == 10 ||
		c.Valor == 11
}

// devulve `true` si la carta es pieza
// segun la el parametro `muestra`
func (c Carta) esPieza(muestra Carta) bool {

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

// Devuelve el puntaje
// no confundir el puntaje con "Poder"
// ojo con Puntaje(7,Espada) == Puntaje(7,Oro) PERO
// (7,Espada) LE GANA A (7,Oro) (tiene mas poder)!!
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
			valeComo := Carta{Palo: c.Palo, Valor: muestra.Valor}
			puntaje = valeComo.calcPuntaje(muestra)
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

// guarismo ficticio y abstracto para simplificar
// las comparaciones
func (c Carta) calcPoder(muestra Carta) int {
	var poder int

	if c.esPieza(muestra) {
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
			poder = valeComo.calcPoder(muestra)
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
	p := rand.Perm(maxCartaID)
	randomSample := make([]int, n)

	for i, r := range p[:n] {
		randomSample[i] = r
	}

	return randomSample
}
