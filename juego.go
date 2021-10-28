package truco

import (
	"bytes"
	"fmt"
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
	out   io.Writer `json:"-"`
	ErrCh chan bool `json:"-"`
}

// Cmd nexo capa presentacion con capa logica
func (j *Juego) Cmd(cmd string) error {

	pkts, err := j.Partida.Cmd(cmd)
	if err != nil {
		return err
	}

	for _, pkt := range pkts {
		enco.Write(j.out, pkt)
	}

	return nil
}

// String retorna una representacion en formato de string
func (j *Juego) String() string {
	return pdt.Renderizar(j.Partida)
}

func (j *Juego) Notify() {

	// ojo primero hay que grabar el buff, luego avisar
	enco.Write(j.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TimeOut, "INTERRUMPING!! Roro tardo demasiado en jugar. Mano ganada por Rojo"),
	))

	j.ErrCh <- true
}

// Abandono da por ganada la partida al equipo contario
func (j *Juego) Abandono(jugador string) error {
	// encuentra al jugador
	manojo, ok := j.Partida.Manojo[jugador]
	if !ok {
		return fmt.Errorf("usuario %s no encontrado", jugador)
	}
	// doy por ganador al equipo contrario
	equipoContrario := manojo.Jugador.GetEquipoContrario()
	ptsFaltantes := int(j.Puntuacion) - j.Puntajes[equipoContrario]
	j.SumarPuntos(equipoContrario, ptsFaltantes)

	enco.Write(j.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.Abandono, manojo.Jugador.ID),
	))

	return nil
}

// NuevaPartida retorna nueva partida; error si hubo
func NuevaPartida(puntuacion pdt.Puntuacion, equipoAzul, equipoRojo []string) (*Juego, io.Reader, error) {

	Partida, err := pdt.NuevaPartida(puntuacion, equipoAzul, equipoRojo)

	if err != nil {
		return nil, nil, err
	}

	j := Juego{
		Partida: Partida,
	}

	buff := new(bytes.Buffer)
	j.out = buff
	j.ErrCh = make(chan bool, 1)

	for _, m := range j.Ronda.Manojos {
		enco.Write(j.out, enco.Pkt(
			enco.Dest(m.Jugador.ID),
			enco.Msg(enco.NuevaPartida, j.Partida.PerspectivaCacheFlor(&m)),
		))

	}

	return &j, buff, nil
}
