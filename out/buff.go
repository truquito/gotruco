package out

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

// Write .
func Write(buff *bytes.Buffer, d *Pkt) error {
	enc := gob.NewEncoder(buff)
	err := enc.Encode(d)
	return err
}

// Read retorna el pkt mas antiguo sin leer
func Read(buff *bytes.Buffer) (*Pkt, error) {
	e := new(Pkt)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Consume consume el buffer
func Consume(buff *bytes.Buffer) {
	for {
		e, err := Read(buff)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(e.String())
	}
}
