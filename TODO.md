-   no deberia poder auto quererse   (ni auto no-quererse)
    p.SetSigJugada("Alvaro Envido")
	p.SetSigJugada("Roro Envido")
	p.SetSigJugada("Alvaro Real-Envido")
	p.SetSigJugada("Roro Falta-Envido")
	p.SetSigJugada("Roro Quiero")

-   test 2 vecees se va al mazo
    p.SetSigJugada("Juan Flor")
    p.SetSigJugada("Pedro Mazo")
	p.SetSigJugada("Pedro Mazo")

- hacer el test EnvidoTrucoRejected -> deberia de rejectear el truco porque "el envido esta primero"

- test envido n-veces -> deberia tener un parate

- hay que distinguir bien lo que son errores del sistema de lo que son errores de los jugadores.
    y setear bien los canales por los que se van a mostrar
    ejemplo:
        se ingresa el comando p.SetSigJugada("Quiero") o p.SetSigJugada("Schumacher Flor")
        ->
        deberia loguearlo en el sistema

- si alguien canta envido y resulta que hay UN SOLO (1) jugador con flor ->
    no hacer todo el listener/timeout de la flor; sino que cantarlo de una y fue;
    computar las "respuestas que se escucharian"

- ganador de una mano en ronda.mano.ganador

- cuando se suman puntos -> checkear que no se haya acabado la partida

- seguridad: que pasa si uno me manda un mensaje de terminar al canal <-- tokens
