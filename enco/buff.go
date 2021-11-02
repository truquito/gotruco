package enco

import (
	"encoding/gob"
	"fmt"
	"io"
)

// Write .
func Write(w io.Writer, d *Packet) error {
	enc := gob.NewEncoder(w)
	err := enc.Encode(d)
	return err
}

// Read retorna el pkt mas antiguo sin leer
func Read(r io.Reader) (*Packet, error) {
	e := new(Packet)
	dec := gob.NewDecoder(r)
	err := dec.Decode(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Consume consume el buffer
func Consume(r io.Reader, callback func(*Packet)) {
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
func Collect(r io.Reader) (res []*Packet) {
	for {
		e, err := Read(r)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
		res = append(res, e)
	}
	return res
}

// contains dado un buffer se fija si contiene un mensaje
// con ese codigo (y string de ser no-nulo)
func Contains(pkts []*Packet, cod CodMsg) bool {
	for _, pkt := range pkts {
		if pkt.Message.Cod == cod {
			return true
		}
	}
	return false
}
