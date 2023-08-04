package truco

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/filevich/truco/enco"
	"github.com/filevich/truco/pdt"
)

// VERSION actual del binario
const VERSION = "0.1.0"

// el envido, la Primera o la mentira
// el envido, la Primera o la mentira
// el truco, la Segunda o el rabón

type c_SIGNAL int

const (
	c_RESET c_SIGNAL = iota
	c_EXIT
)

// Juego :
type Juego struct {
	*pdt.Partida
	mu  *sync.Mutex
	out []enco.Envelope
	// errores
	Err   *enco.Envelope
	ErrCh chan bool `json:"-"`
	// tiempo
	DurTurno time.Duration
	contador chan c_SIGNAL `json:"-"`
	tic      *time.Ticker
}

func (j *Juego) Consumir() []enco.Envelope {
	j.mu.Lock()
	defer j.mu.Unlock()

	res := j.out
	j.out = make([]enco.Envelope, 0, len(res))

	return res
}

// Cmd nexo capa presentacion con capa logica
func (j *Juego) Cmd(cmd string) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.Terminado() {
		return errors.New("el juego termió")
	}

	pkts, err := j.Partida.Cmd(cmd)
	if err != nil {
		return err
	}

	if j.Partida.Verbose {
		j.out = append(j.out, pkts...)
	}

	if j.Partida.Terminada() {
		// entonces paro el contador (goroutine) + tic (ticker)
		j.tic.Stop()

	}

	return nil
}

// String retorna una representacion en formato de string
func (j *Juego) String() string {
	return pdt.Renderizar(j.Partida)
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

	if j.Partida.Verbose {
		pkt := enco.Pkt(
			enco.Dest("ALL"),
			enco.Abandono(manojo.Jugador.ID),
		)

		j.out = append(j.out, pkt)
	}

	return nil
}

func (j *Juego) Expirado() bool {
	return j.Err != nil
}

func (j *Juego) Terminado() bool {
	return j.Partida.Terminada() || j.Expirado()
}

func (j *Juego) contar() {
	const delta float64 = 1.15
	d := float64(j.DurTurno) * delta
	total := time.Duration(math.Ceil(d))
	j.tic = time.NewTicker(total)

	defer func() {
		// fmt.Println("exiting contar")
		j.tic.Stop()
	}()

	for {
		select {
		case s := <-j.contador:
			switch s {
			case c_RESET:
				j.tic.Stop()
				j.tic.Reset(total)
			case c_EXIT:
				j.tic.Stop()
				return // <- se destruye esta goroutine
			}
		case <-j.tic.C:
			u := j.Partida.Ronda.Manojos[j.Ronda.Turno].Jugador.ID
			pkt := enco.Pkt(enco.ALL, enco.TimeOut(u))
			j.Err = &pkt
			j.ErrCh <- true
			return // <- se destruye esta goroutine
		}
	}
}

// NuevoJuego retorna nueva partida; error si hubo
func NuevoJuego(

	puntuacion pdt.Puntuacion,
	equipoAzul,
	equipoRojo []string,
	verbose bool,
	maxTiempoPorTurno time.Duration,

) (*Juego, error) {

	p, err := pdt.NuevaPartida(puntuacion, equipoAzul, equipoRojo, verbose)

	if err != nil {
		return nil, err
	}

	j := Juego{
		Partida: p,
		mu:      &sync.Mutex{},
		out:     make([]enco.Envelope, 0),
		// errores
		ErrCh: make(chan bool, 1),
		Err:   nil,
		// tiempo
		contador: make(chan c_SIGNAL, 1),
		DurTurno: maxTiempoPorTurno,
		tic:      nil,
	}

	// pongo en el buffer un mensaje de Partida{} para cada jugador
	if j.Partida.Verbose {
		for _, m := range j.Ronda.Manojos {
			pkt := enco.Pkt(
				enco.Dest(m.Jugador.ID),
				enco.NuevaPartida{
					Perspectiva: j.Partida.PerspectivaCacheFlor(&m),
				},
			)
			j.out = append(j.out, pkt)
		}
	}

	go j.contar()

	return &j, nil
}
