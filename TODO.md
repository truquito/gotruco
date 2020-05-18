- calc bien los puntos cuando se va al mazo (mergear con lo que esta comentado)

- depues de cantar vale4 (y similares) nadie puede tirar carta, el equipo contrario debe responder si o no

- no necesariamente tiene que ser RETRUCOQUERIDO el estado del truco para
    recien ahi gritar vale4, podes ir de 1???

- en el caso TestFixNacho:
    luego de que `adolfo re-truco` no se deberia poder tirar nada, solo itse al mazo
    ya que la ronda termina y nadie dije quiero o no quiero

- duda en el test TestFixNacho la primera mano queda empardad:
    a quien le toca el sig turno? a richard o a andres?

- duda: es necesario que sea su turno para cantar retruco o el turno de uno de 
    los de mi equipo?

- duda en el test TestFlorFlorContraFlorQuiero (flor_test@158)
    deberia sumar tambien los puntos de las flores?
    implementacion de los puntos sumados en jugada@551

- hay redundancia entre p.Ronda.getLaFlorMasAlta y 
manojoConLaFlorGanadora, _, _ := p.Ronda.cantarFlores(aPartirDe)

- que es mas rapido para saber si un jugador tiene flor:
    - jugada.autor.tieneFlor(p.Ronda.Muestra)
    - o fijarse en un array de jugadores con flor y la op contains
    - un mapa de booleanos o similar

- alguien que se fue al mazo con flor puede llegar a decir quiero a una contraflor
    DE HECHO, DEBERIA SER UN CHECKING PARA TODAS LAS JUGADAS

- duda: es necesario que sea tu turno para cantar envido?

- se esta jugando de a 2 (o mas) uno tiene flor, pero el otro se va al mazo
  deberia de sumar los puntos de la flor a pesar de que no fue cantada

- que pasa si hay 3 con flor A,B y C; A canta flor, B dice contra flor, A dice quiero -> pero B nunca canto su flor

- duda: actualmente si hay 3 flores: 2 del equipo rojo, 1 del equipo azul; todos cantan "flor" sin mas, y la flor del equipo azul es la mas alta ->
se lleva el puntaje el equipo azul de las 3 flores (3+3+3) ; eso esta bien?

- si 2 tienen flor y esos 2 las cantan -> tampoco termina el bucle de la flor log: `fix 2 flores y bucle.log`

- duda: es necesario que sea tu turno para cantar la flor?

- agregar el control de que los nombres de los jguadores tienen que ser > 0

-   // que pasa cuando el ganador de una mano se habia ido al mazo?
    // no se tiene que poder:
    // si en esta mano ya jugaste carta -> no te podes ir al mazo
    // o bien: solo te podes ir al mazo cuando es tu turno
    // luego este metodo es correcto

    para irse al mazo debe ser tu turno

- si se va al mazo -> deberia pasar al siguiente turno

- no puede tirar carta si el estado del truco es respondible por uno de
    su propio equipo (contrario al que propuso el truco)

- jugada@58 p.Ronda.sigTurno()

- Ronda.getIdx es el nuevo ; utils.obtenerIdx es el viejo !!
    eventualmente se tendrian que ir los dos a la mierda

- EN LA ABSTRACCION DE LAS JUGADAS LO PRIMERO ES FIJARSE SI TERMINO O NO
    LA PARTIDA!!!!!!!!!!!!!!

- que pasa si alguien dice truco al final, y el ultimo no contesta la apuesta
    y juega su carta? no deberia de estar permitido

- los channels estan como variables globales, hay que ponerlos como privados
    de Partida

- puede que se llame a p.evaluarMano() y que no se haya tirado ninguna carta?

- duda: que pasa si no se tiran las cartas en orden? es esto posible?
    creo que no (por la imposicion en el turno en tirarCarta) 
    pero aun asi testearlo.

- cantJugadoresEnJuego al pedo? potencial uso pero no r.sigHabilitado()

- la visibilidad desde la partida hacia los jugadores esta al pepe?

- agregarle un ID a jugador y que las busquedas las haga por este campo
    esta id debe ser sercreta y random; porque si yo se la id de otro puedo jugar por el

- el no quiero de la flor @TestFixNoFlor no tiene output

- el ordden del output si hay flor cuando se grita truco esta mal:
<< Andres canta flor
No es posible cantar truco ahora

- duda:
    >A truco
    >el evndio esta primer
    *se juega el envido*
    ----> ahora el truco sigue en juego como si lo hubiese cantado A? o hay que volver a cantarlo?

- checkeos de la gente que se fue al mazo no podria hacer nada:
    eg:
        los que se fueron al mazo pueden decir quiero/noquiero
    idea solucion: que la struct Jugada de la que todas extienden tenga un metodo que checkeo eso:

- revisar que todo lo que usa el estado del truco este bien luego del cambio de los estados que faltaban en el truco

- empezar a usar:
    // termino := p.sumarPuntos(p.Ronda.truco.cantadoPor.jugador.equipo, totalPts)
    // if !termino {
    // 	p.nuevaRonda()
    // }

- el no quiero falta la parte de la flor

- jugada@618:
    cambiar aPartirDe
    que intenta obtener la posicion del jugador a partir de su "Perfil"
    PARA ELLO BORRAR LA FUNCION EN UTILS:
    obtenerIdx(jugador *Jugador, jugadores []Jugador) (JugadorIdx, error)

    e ir arreglando los problemas usando la siguiente recomendacion:

    SOLUCION: QUE EN LA STRUCT DEL ENVIDO EL PUNTERO SEA A EL MANOJO NO AL PERFIL;
    luego usar p.Ronda.getPerfil( p.Ronda.envido.cantadoPor )

    

- eliminar p.getJugador() y usar Ronda.getManojo()

- JugadorIdx al pedo?

- probar una ronda de envidos donde el mano es el ultimo jugador
    probablemente de error out of index porque esta mal programadao
    el get ronda.Envidos()

- hay redundancia entre cantarFloresSiLasHay y cantarFlores

- hacer getElEnvido() con indices

- eliminar los metodos .eval() que se usan en los .quiero() y ponerlos directamente ahi

- al cantar el envido -> si se estaba jugando el truco, ponerlo en no jjugado aun?

- error: el que canta flor se puede ir al mazo sin que nadie haya respondido nada

- esta duplicado el codigo de irse al mazo con el de "no quiero"

- cambiar estructuras, metodos y funciones que eran publicas (porque empezaban con mayuscula) eg. `struct Envido` por privadas `struct envido`

- el autor del envido es un *Jugador; eso es deprected; deberia ser
    *Manojo
    
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
