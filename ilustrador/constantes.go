package ilustrador

import (
	C "github.com/jpfilevich/canvas"
)

type posicion int

const (
	a posicion = iota
	b
	c
	d
	e
	f
)

var (
	areasJugadores = map[string](map[posicion]C.Rectangle){
		"nombres": map[posicion]C.Rectangle{
			a: C.Rectangle{
				From: C.Point{X: 15, Y: 11},
				To:   C.Point{X: 25, Y: 11},
			},
			b: C.Rectangle{
				From: C.Point{X: 29, Y: 11},
				To:   C.Point{X: 39, Y: 11},
			},
			c: C.Rectangle{
				From: C.Point{X: 44, Y: 6},
				To:   C.Point{X: 53, Y: 6},
			},
			d: C.Rectangle{
				From: C.Point{X: 29, Y: 1},
				To:   C.Point{X: 39, Y: 1},
			},
			e: C.Rectangle{
				From: C.Point{X: 15, Y: 1},
				To:   C.Point{X: 25, Y: 1},
			},
			f: C.Rectangle{
				From: C.Point{X: 0, Y: 6},
				To:   C.Point{X: 9, Y: 6},
			},
		},
		"tiradas": map[posicion]C.Rectangle{
			a: C.Rectangle{
				From: C.Point{X: 19, Y: 7},
				To:   C.Point{X: 24, Y: 9},
			},
			b: C.Rectangle{
				From: C.Point{X: 29, Y: 7},
				To:   C.Point{X: 34, Y: 9},
			},
			c: C.Rectangle{
				From: C.Point{X: 35, Y: 5},
				To:   C.Point{X: 40, Y: 7},
			},
			d: C.Rectangle{
				From: C.Point{X: 29, Y: 3},
				To:   C.Point{X: 34, Y: 5},
			},
			e: C.Rectangle{
				From: C.Point{X: 19, Y: 3},
				To:   C.Point{X: 24, Y: 5},
			},
			f: C.Rectangle{
				From: C.Point{X: 13, Y: 5},
				To:   C.Point{X: 18, Y: 7},
			},
		},
		"tooltips": map[posicion]C.Rectangle{
			a: C.Rectangle{
				From: C.Point{X: 15, Y: 12},
				To:   C.Point{X: 24, Y: 12},
			},
			b: C.Rectangle{
				From: C.Point{X: 29, Y: 12},
				To:   C.Point{X: 38, Y: 12},
			},
			c: C.Rectangle{
				From: C.Point{X: 44, Y: 5},
				To:   C.Point{X: 53, Y: 5},
			},
			d: C.Rectangle{
				From: C.Point{X: 29, Y: 0},
				To:   C.Point{X: 38, Y: 0},
			},
			e: C.Rectangle{
				From: C.Point{X: 15, Y: 0},
				To:   C.Point{X: 24, Y: 0},
			},
			f: C.Rectangle{
				From: C.Point{X: 0, Y: 5},
				To:   C.Point{X: 5, Y: 9},
			},
		},
	}

	otrasAreas = map[string]C.Rectangle{
		"exteriorMesa": C.Rectangle{
			From: C.Point{X: 11, Y: 2},
			To:   C.Point{X: 42, Y: 10},
		},
		"interiorMesa": C.Rectangle{
			From: C.Point{X: 12, Y: 3},
			To:   C.Point{X: 41, Y: 9},
		},
	}
)
