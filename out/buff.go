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

// Consumer el callback de `Consume`
type Consumer func(*Packet)

// Print imprime el packete
func Print(pkt *Packet) {
	fmt.Println(pkt.String())
}

// Consume consume el buffer
func Consume(buff *bytes.Buffer, callback Consumer) {
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
