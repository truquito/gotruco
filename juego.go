package truco

import (
	"bytes"
	"io"

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
	Out   io.ReadWriter `json:"-"`
	ErrCh chan bool     `json:"-"`
}

// Cmd nexo capa presentacion con capa logica
func (j *Juego) Cmd(cmd string) error {

	pkts, err := j.Partida.Cmd(cmd)
	if err != nil {
		return err
	}

	for _, pkt := range pkts {
		enco.Write(j.Out, pkt)
	}

	return nil
}

// String retorna una representacion en formato de string
func (j *Juego) String() string {
	return pdt.Renderizar(j.Partida)
}

func (j *Juego) Notify() {

	// ojo primero hay que grabar el buff, luego avisar
	enco.Write(j.Out, enco.Pkt(
		enco.Dest("ALL"),
		enco.TimeOut("INTERRUMPING!! Roro tardo demasiado en jugar. Mano ganada por Rojo"),
	))

	j.ErrCh <- true
}

// Abandono da por ganada la partida al equipo contario
func (j *Juego) Abandono(jugador string) error {
	// encuentra al jugador
	manojo := j.Partida.Manojo(jugador)
	// doy por ganador al equipo contrario
	equipoContrario := manojo.Jugador.GetEquipoContrario()
	ptsFaltantes := int(j.Puntuacion) - j.Puntajes[equipoContrario]
	j.SumarPuntos(equipoContrario, ptsFaltantes)

	enco.Write(j.Out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Abandono(manojo.Jugador.ID),
	))

	return nil
}

// NuevoJuego retorna nueva partida; error si hubo
func NuevoJuego(puntuacion pdt.Puntuacion, equipoAzul, equipoRojo []string) (*Juego, error) {

	p, err := pdt.NuevaPartida(puntuacion, equipoAzul, equipoRojo)

	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)

	j := Juego{
		Partida: p,
		Out:     buff,
		ErrCh:   make(chan bool, 1),
	}

	// pongo en el buffer un mensaje de Partida{} para cada jugador
	for _, m := range j.Ronda.Manojos {
		enco.Write(j.Out, enco.Pkt(
			enco.Dest(m.Jugador.ID),
			enco.NuevaPartida{
				Perspectiva: j.Partida.PerspectivaCacheFlor(&m),
			},
		))
	}

	return &j, nil
}
