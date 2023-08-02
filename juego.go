package truco

import (
	"sync"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/pdt"
)

// VERSION actual del binario
const VERSION = "0.1.0"

// el envido, la Primera o la mentira
// el envido, la Primera o la mentira
// el truco, la Segunda o el rabón

// Juego :
type Juego struct {
	*pdt.Partida
	mu    *sync.Mutex
	out   []enco.Packet
	ErrCh chan bool `json:"-"`
}

func (j *Juego) Consume() []enco.Packet {
	j.mu.Lock()
	defer j.mu.Unlock()

	res := j.out
	j.out = make([]enco.Packet, 0, len(res))

	return res
}

// Cmd nexo capa presentacion con capa logica
func (j *Juego) Cmd(cmd string) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	pkts, err := j.Partida.Cmd(cmd)
	if err != nil {
		return err
	}

	j.out = append(j.out, pkts...)

	return nil
}

// String retorna una representacion en formato de string
func (j *Juego) String() string {
	return pdt.Renderizar(j.Partida)
}

func (j *Juego) Notify() {
	j.mu.Lock()
	defer j.mu.Unlock()

	// deprecated: ojo primero hay que grabar el buff, luego avisar
	pkt := enco.Pkt(
		enco.Dest("ALL"),
		enco.TimeOut("INTERRUMPING!! Roro tardo demasiado en jugar. Mano ganada por Rojo"),
	)

	j.out = append(j.out, pkt)

	j.ErrCh <- true
}

// Abandono da por ganada la partida al equipo contario
func (j *Juego) Abandono(jugador string) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	// encuentra al jugador
	manojo := j.Partida.Manojo(jugador)
	// doy por ganador al equipo contrario
	equipoContrario := manojo.Jugador.GetEquipoContrario()
	ptsFaltantes := int(j.Puntuacion) - j.Puntajes[equipoContrario]
	j.SumarPuntos(equipoContrario, ptsFaltantes)

	pkt := enco.Pkt(
		enco.Dest("ALL"),
		enco.Abandono(manojo.Jugador.ID),
	)

	j.out = append(j.out, pkt)

	return nil
}

// NuevoJuego retorna nueva partida; error si hubo
func NuevoJuego(puntuacion pdt.Puntuacion, equipoAzul, equipoRojo []string) (*Juego, error) {

	p, err := pdt.NuevaPartida(puntuacion, equipoAzul, equipoRojo)

	if err != nil {
		return nil, err
	}

	j := Juego{
		Partida: p,
		mu:      &sync.Mutex{},
		out:     make([]enco.Packet, 0),
		ErrCh:   make(chan bool, 1),
	}

	// pongo en el buffer un mensaje de Partida{} para cada jugador
	for _, m := range j.Ronda.Manojos {
		pkt := enco.Pkt(
			enco.Dest(m.Jugador.ID),
			enco.NuevaPartida{
				Perspectiva: j.Partida.PerspectivaCacheFlor(&m),
			},
		)
		j.out = append(j.out, pkt)
	}

	return &j, nil
}
