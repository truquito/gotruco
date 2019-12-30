all: main

# Objetivos que no son archivos.
.PHONY: all clean_bin clean_test clean testing

TESTDIR = test
EJECUTABLE = main

# compilador
GC = go
# opciones de compilación
GCFLAGS = run
# opciones de testing
TFLAGS = test -count=1

# argumentos testing
F  = file_test.go
T  = TestName

# dependencias
DEPENDENCIAS = cartas.go jugada.go partida.go utils.go jugador.go envido.go ronda.go mano.go testing_utils.go manojo.go

# $^ se expande en todas las dependencias
$(EJECUTABLE): $(DEPENDENCIAS)
	$(GC) $(GCFLAGS) $^ #-o $@

test: $(DEPENDENCIAS)
	$(GC) $(TFLAGS) $(F) $^ -run $(T)

ftest: $(DEPENDENCIAS)
	$(GC) $(TFLAGS) $(F) $^

# e.g., 
# $ make testing F=main_test.go T=TestEnvidoI
# generalized alternative
# go test -v *.go -run TestEnvidoI