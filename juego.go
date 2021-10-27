package truco

import (
	"bytes"
	"fmt"
	"io"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/pdt"
)

// VERSION actual del binario
const VERSION = "0.0.1"

// el envido, la Primera o la mentira
// el envido, la Primera o la mentira
// el truco, la Segunda o el rabón

// Partida :
type Partida struct {
	*pdt.PartidaDT
	out   io.Writer `json:"-"`
	ErrCh chan bool `json:"-"`
}

// Cmd nexo capa presentacion con capa logica
func (p *Partida) Cmd(cmd string) error {

	pkts, err := p.PartidaDT.Cmd(cmd)
	if err != nil {
		return err
	}

	for _, pkt := range pkts {
		enco.Write(p.out, pkt)
	}

	return nil
}

// String retorna una representacion en formato de string
func (p *Partida) String() string {
	return pdt.Renderizar(p.PartidaDT)
}

func (p *Partida) Notify() {

	// ojo primero hay que grabar el buff, luego avisar
	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.TimeOut, "INTERRUMPING!! Roro tardo demasiado en jugar. Mano ganada por Rojo"),
	))

	p.ErrCh <- true
}

// Abandono da por ganada la partida al equipo contario
func (p *Partida) Abandono(jugador string) error {
	// encuentra al jugador
	manojo, err := p.GetManojoByStr(jugador)
	if err != nil {
		return fmt.Errorf("usuario %s no encontrado", jugador)
	}
	// doy por ganador al equipo contrario
	equipoContrario := manojo.Jugador.GetEquipoContrario()
	ptsFaltantes := int(p.Puntuacion) - p.Puntajes[equipoContrario]
	p.SumarPuntos(equipoContrario, ptsFaltantes)

	enco.Write(p.out, enco.Pkt(
		enco.Dest("ALL"),
		enco.Msg(enco.Abandono, manojo.Jugador.ID),
	))

	return nil
}

// NuevaPartida retorna nueva partida; error si hubo
func NuevaPartida(puntuacion pdt.Puntuacion, equipoAzul, equipoRojo []string) (*Partida, io.Reader, error) {

	partidaDt, err := pdt.NuevaPartidaDt(puntuacion, equipoAzul, equipoRojo)

	if err != nil {
		return nil, nil, err
	}

	p := Partida{
		PartidaDT: partidaDt,
	}

	buff := new(bytes.Buffer)
	p.out = buff
	p.ErrCh = make(chan bool, 1)

	for _, m := range p.Ronda.Manojos {
		enco.Write(p.out, enco.Pkt(
			enco.Dest(m.Jugador.ID),
			enco.Msg(enco.NuevaPartida, p.PartidaDT.PerspectivaCacheFlor(&m)),
		))

	}

	return &p, buff, nil
}
