package truco

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

// VERSION actual de la librería/binario
const VERSION = "0.3.0"

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
	Tic      *time.Ticker  `json:"-"`
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
		j.Tic.Stop()

	}

	return nil
}

// String retorna una representacion en formato de string
func (j *Juego) String() string {
	return j.Partida.String()
}

func (j *Juego) Expirado() bool {
	return j.Err != nil
}

func (j *Juego) Terminado() bool {
	return j.Partida.Terminada() || j.Expirado()
}

// Abandono da por ganada la partida al equipo contario
func (j *Juego) Abandono(jugador string) {
	j.mu.Lock()
	defer j.mu.Unlock()

	manojo := j.Partida.Manojo(jugador)
	// doy por ganador al equipo contrario
	equipoContrario := manojo.Jugador.GetEquipoContrario()
	ptsFaltantes := int(j.Puntuacion) - j.Puntajes[equipoContrario]
	j.SumarPuntos(equipoContrario, ptsFaltantes)
	// agarro al primer manojo
	jugGanador := j.Partida.Ronda.Manojos[0].Jugador
	esDelEquipoQueAbandono := jugGanador.Equipo == j.Partida.Manojo(jugador).Jugador.Equipo
	if esDelEquipoQueAbandono {
		// entonces tomo el siguiente
		jugGanador = j.Partida.Ronda.Manojos[1].Jugador
	}

	if j.Partida.Verbose {
		j.out = append(j.out, enco.Env(
			enco.ALL,
			enco.SumaPts{
				Autor:  jugGanador.ID,
				Razon:  enco.Abandonaron,
				Puntos: ptsFaltantes,
			},
		))
	}

	// async err
	pkt := enco.Env(enco.ALL, enco.Abandono(jugador))
	j.Err = &pkt
	j.ErrCh <- true
	j.Tic.Stop()
}

// no hay motivo alguno, simplemente se aborta
func (j *Juego) Abortar(abandonador string) {
	pkt := enco.Env(enco.ALL, enco.Abandono(abandonador))
	j.Err = &pkt
	j.ErrCh <- true
	j.Tic.Stop()
}

// la tarea de hacer close(j.ErrCh) la tiene que hacer
// el `defer` del `select` del `j.ErrCh`
// func (j *Juego) Close() {
// 	close(j.ErrCh)
// 	j.tic.Stop()
// }

func (j *Juego) GetMaxTiempoPorTurno() time.Duration {
	const delta float64 = 1.15
	d := float64(j.DurTurno) * delta
	total := time.Duration(math.Ceil(d))
	return total
}

func (j *Juego) contar() {
	total := j.GetMaxTiempoPorTurno()
	j.Tic = time.NewTicker(total)

	defer func() {
		// exiting contar !!
		j.Tic.Stop()
	}()

	for {
		select {
		case s := <-j.contador:
			switch s {
			case c_RESET:
				j.Tic.Stop()
				j.Tic.Reset(total)
			case c_EXIT:
				j.Tic.Stop()
				return // <- se destruye esta goroutine
			}
		case <-j.Tic.C:
			// quien debia responder?
			u := pdt.Rho(j.Partida).Jugador.ID
			pkt := enco.Env(enco.ALL, enco.TimeOut(u))
			j.Err = &pkt
			j.ErrCh <- true
			j.Tic.Stop()
			return // <- se destruye esta goroutine y pa
		}
	}
}

// PRE: la goroutine sigue en pie
func (j *Juego) Reset(

	puntuacion pdt.Puntuacion,
	equipoAzul,
	equipoRojo []string,
	limiteEnvido int,
	verbose bool,
	maxTiempoPorTurno time.Duration,

) {
	j.mu.Lock()
	defer j.mu.Unlock()

	// libero recursos
	j.Tic.Stop()
	j.contador <- c_EXIT // sale la goroutine
	close(j.ErrCh)
	j.Err = nil

	// creo recursos
	p, _ := pdt.NuevaPartida(puntuacion, equipoAzul, equipoRojo, limiteEnvido, verbose)

	j.Partida = p
	j.out = make([]enco.Envelope, 0) // descarto los envelopes
	j.ErrCh = make(chan bool, 1)
	j.Err = nil
	j.contador = make(chan c_SIGNAL, 1)

	// pongo en el buffer un mensaje de Partida{} para cada jugador
	if j.Partida.Verbose {
		for _, m := range j.Ronda.Manojos {
			pkt := enco.Env(
				enco.Dest(m.Jugador.ID),
				enco.NuevaPartida{
					Perspectiva: j.Partida.PerspectivaCacheFlor(&m),
				},
			)
			j.out = append(j.out, pkt)
		}
	}

	go j.contar()
	// RESET DONE !!
}

// NuevoJuego retorna nueva partida; error si hubo
func NuevoJuego(

	puntuacion pdt.Puntuacion,
	equipoAzul,
	equipoRojo []string,
	limiteEnvido int,
	verbose bool,
	maxTiempoPorTurno time.Duration,

) (*Juego, error) {

	p, err := pdt.NuevaPartida(puntuacion, equipoAzul, equipoRojo, limiteEnvido, verbose)

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
		Tic:      nil,
	}

	// pongo en el buffer un mensaje de Partida{} para cada jugador
	if j.Partida.Verbose {
		for _, m := range j.Ronda.Manojos {
			pkt := enco.Env(
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
