package out

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

// Write .
func Write(buff *bytes.Buffer, d *Packet) error {
	enc := gob.NewEncoder(buff)
	err := enc.Encode(d)
	return err
}

// Read retorna el pkt mas antiguo sin leer
func Read(buff *bytes.Buffer) (*Packet, error) {
	e := new(Packet)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Consume consume el buffer
func Consume(buff *bytes.Buffer, callback func(*Packet)) {
	for {
		e, err := Read(buff)
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
func Collect(buff *bytes.Buffer) (res []*Packet) {
	for {
		e, err := Read(buff)
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
