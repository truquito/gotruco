package enco

import (
	"encoding/gob"
	"fmt"
	"io"
)

// Write .
func Write(w io.Writer, d Envelope) error {
	enc := gob.NewEncoder(w)
	err := enc.Encode(d)
	return err
}

// Read retorna el pkt mas antiguo sin leer
func Read(r io.Reader) (*Envelope, error) {
	e := new(Envelope)
	dec := gob.NewDecoder(r)
	err := dec.Decode(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Consume consume el buffer
func Consume(r io.Reader, callback func(*Envelope)) {
	for {
		e, err := Read(r)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
		callback(e)
	}
}

// Collect pasa de buffer a slice
func Collect(r io.Reader) (res []Envelope) {
	for {
		e, err := Read(r)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
		res = append(res, *e)
	}
	return res
}

// contains dado un buffer se fija si contiene un mensaje
// con ese codigo (y string de ser no-nulo)
func Contains(pkts []Envelope, cod CodMsg) bool {
	for _, pkt := range pkts {
		if pkt.Message.Cod() == cod {
			return true
		}
	}
	return false
}
