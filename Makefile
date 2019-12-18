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
TFLAGS = test -v

# argumentos testing
F  = main_test.go
T  = TestTmp

# dependencias
DEPENDENCIAS = main.go ronda.go cartas.go mano.go manojo.go jugador.go jugada.go utils.go partida.go envido.go truco.go

# $^ se expande en todas las dependencias
$(EJECUTABLE): $(DEPENDENCIAS)
	$(GC) $(GCFLAGS) $^ #-o $@

testing: $(DEPENDENCIAS)
	$(GC) $(TFLAGS) $(F) $^ -run $(T)

# e.g., 
# $ make testing F=main_test.go T=TestEnvidoI
# generalized alternative
# go test -v *.go -run TestEnvidoI