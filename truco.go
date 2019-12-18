package main

// EstadoTruco : enum
type EstadoTruco int

// enums del truco
const (
	NOCANTADO EstadoTruco = 1
	TRUCO     EstadoTruco = 2
	RETRUCO   EstadoTruco = 3
	VALE4     EstadoTruco = 4
)