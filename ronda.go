package main

import (
	"fmt"
)

// Ronda :
type Ronda struct {
	manoEnJuego      NumMano // Numero de mano: 1era 2da 3era
	jugadoresEnJuego int     // Numero de jugadores que no se fueron al mazo

	/* Indices */
	elMano JugadorIdx
	turno  JugadorIdx
	pies   [2]JugadorIdx

	/* toques, gritos y cantos */
	envido Envido      // Estado del envido, la primera o la mentira
	flor   EstadoFlor  // Estado del envido, la primera o la mentira
	truco  EstadoTruco // Estado del truco, la segunda o el rabón

	/* cartas */
	manojos []Manojo
	muestra Carta

	manos []Mano
}

func (r *Ronda) checkFlorDelMano() {
	tieneFlor, tipoFlor := r.getElMano().tieneFlor(r.muestra)
	if tieneFlor {
		r.getElMano().cantarFlor(tipoFlor, r.muestra)
	}
}

// devuelve `false` si la ronda se acabo
func (r *Ronda) enJuego() bool {
	return true
}

// los anteriores a `aPartirDe` (incluido este) no
// son necesarios de checkear porque ya han sido
// checkeados si tenian flor
func (r Ronda) cantarFloresSiLasHay(aPartirDe JugadorIdx) {
	for _, jugador := range r.manojos[aPartirDe+1:] {
		tieneFlor, tipoFlor := jugador.tieneFlor(r.muestra)
		if tieneFlor {
			// todo:
			tieneFlor = false
			tipoFlor = tipoFlor + 1
			// var jugada IJugada = responderNoQuiero{}
			// jugador.cantarFlor(tipoFlor, r.muestra)
			r.envido.estado = DESHABILITADO
			break
		}
	}
}

func (r Ronda) checkFlores(aPartirDe JugadorIdx) (hayFlor bool,
	jugador *Jugador, tipoFlor int) {
	for _, manojo := range r.manojos[aPartirDe+1:] {
		tieneFlor, tipoFlor := manojo.tieneFlor(r.muestra)
		if tieneFlor {
			return true, manojo.jugador, tipoFlor
		}
	}
	return false, nil, 0
}

func (r Ronda) getElMano() *Manojo {
	return &r.manojos[r.elMano]
}

func (r Ronda) getElTurno() *Manojo {
	return &r.manojos[r.turno]
}

func (r Ronda) getManoAnterior() *Mano {
	return &r.manos[r.manoEnJuego-1]
}

func (r Ronda) getManoActual() *Mano {
	return &r.manos[r.manoEnJuego]
}

func (r Ronda) setTurno() {
	// si es la primera mano que se juega
	// entonces es el turno del mano
	if r.manoEnJuego == primera {
		r.turno = r.elMano
		// si no, es turno del ganador de
		// la mano anterior
	} else {
		r.turno = r.getManoAnterior().ganador
	}
}

// Print Imprime la informacion de la ronda
func (r Ronda) Print() {
	for i := range r.manojos {
		r.manojos[i].jugador.Print()
	}

	fmt.Printf("Y la muestra es\n    - %s\n", r.muestra.toString())
	fmt.Printf("\nEl mano actual es: %s\nEs el turno de %s\n\n",
		r.getElMano().getPerfil().nombre, r.getElTurno().getPerfil().nombre)
	
	imprimirJugadas()
}
