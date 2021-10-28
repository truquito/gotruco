package pdt

import (
	"testing"

	"github.com/filevich/truco/util"
)

func TestEliminarDePendientesPorCantarFlor(t *testing.T) {
	e := Envite{
		SinCantar: []string{
			"foo",
			"bar",
			"foobar",
		},
	}

	e.cantoFlor("bar")

	ok := util.All(
		len(e.SinCantar) == 2,
		e.SinCantar[0] == "foo",
		e.SinCantar[1] == "foobar",
	)

	if !ok {
		t.Error("no elimino correctamente al jugador `bar` del slice")
	}
}
