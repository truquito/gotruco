package ilustrador

/*
               ▒▒▒▒▒▒▒▒▒▒    ▒▒▒▒▒▒▒▒▒▒
               ▒▒E     ▒▒    ▒▒D     ▒▒
           ╔══════════════════════════════╗
           ║       ┌12┐┐┐    ┌12┐┐┐       ║
           ║       │Or│││    │Or│││       ║
▒▒▒▒▒▒▒▒▒▒ ║ ┌12┐┐┐└──┘┘┘┌12┐└──┘┘┘┌12┐┐┐ ║ ▒▒▒▒▒▒▒▒▒▒
▒▒▒F   ▒▒▒ ║ │Or│││      │Or│      │Or│││ ║ ▒▒▒C  ▒▒▒▒
           ║ └──┘┘┘┌12┐┐┐└──┘┌12┐┐┐└──┘┘┘ ║
           ║       │Or│││    │Or│││       ║
           ║       └──┘┘┘    └──┘┘┘       ║
           ╚══════════════════════════════╝
               ▒▒A     ▒▒    ▒▒B     ▒▒
               ▒▒▒▒▒▒▒▒▒▒    ▒▒▒▒▒▒▒▒▒▒
*/

type posicion int

const (
	a posicion = iota
	b
	c
	d
	e
	f
)

// Carta naipe espanol
type Carta struct {
	palo  string
	valor int
}

type manojo struct {
	cartas  [3]Carta
	tiradas [3]bool
}

var (
	areasJugadores = map[string](map[posicion]rectangle){
		"nombres": map[posicion]rectangle{
			a: rectangle{pos{15, 11}, pos{25, 11}},
			b: rectangle{pos{29, 11}, pos{39, 11}},
			c: rectangle{pos{44, 6}, pos{53, 6}},
			d: rectangle{pos{29, 1}, pos{39, 1}},
			e: rectangle{pos{15, 1}, pos{25, 1}},
			f: rectangle{pos{0, 6}, pos{9, 6}},
		},
		"tiradas": map[posicion]rectangle{
			a: rectangle{pos{19, 7}, pos{24, 9}},
			b: rectangle{pos{29, 7}, pos{34, 9}},
			c: rectangle{pos{35, 5}, pos{40, 7}},
			d: rectangle{pos{29, 3}, pos{34, 5}},
			e: rectangle{pos{19, 3}, pos{24, 5}},
			f: rectangle{pos{13, 5}, pos{18, 7}},
		},
		"tooltips": map[posicion]rectangle{
			a: rectangle{pos{15, 12}, pos{24, 12}},
			b: rectangle{pos{29, 12}, pos{38, 12}},
			c: rectangle{pos{44, 5}, pos{53, 5}},
			d: rectangle{pos{29, 0}, pos{38, 0}},
			e: rectangle{pos{15, 0}, pos{24, 0}},
			f: rectangle{pos{0, 5}, pos{5, 9}},
		},
	}

	otrasAreas = map[string]rectangle{
		"exteriorMesa": rectangle{pos{11, 2}, pos{42, 10}},
		"interiorMesa": rectangle{pos{12, 3}, pos{41, 9}},
	}
)
