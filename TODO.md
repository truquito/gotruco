# TODO
- tabular los mensajes a lo minimo
- si se va al mazo que no imprima nada
- pasar todos los .ToJSON para que devuelvan json y no strings!!

# BUGS

# LIMPIEZA DE CODIGO
- la struct truco esta con minuscula salame
- Hay codigo repetido entre noQuiero y mazo cuando niega la flor (codigo copiado)
- hay redundancia entre cantarFloresSiLasHay y cantarFlores
- hay redundancia entre p.Ronda.getLaFlorMasAlta y 
    manojoConLaFlorGanadora, _, _ := p.Ronda.execCantarFlores(aPartirDe)
- ??? esta duplicado el codigo de irse al mazo con el de "no quiero"
- [ESTO PUEDE SER VIEJO(?)] hacer getElEnvido() con indices
- ??? - eliminar los metodos .eval() que se usan en los .quiero() y ponerlos 
    directamente ahi

# SEGURIDAD
- que pasa si uno me manda un mensaje de terminar al canal <-- tokens
- agregar el control de que los len(nombres de los jguadores) tienen que ser > 0
- agregarle un ID a jugador y que las busquedas las haga por este campo
    esta id debe ser sercreta y random; porque si yo se la id de otro puedo 
    jugar por el
- generar un comando especial que genere un panic, para asi testear los planes
    de contingencia en caso de falla
- recontra revisar para todas las jugadas:
    checkeos de la gente que se fue al mazo no podria hacer nada:
    eg:
        los que se fueron al mazo pueden decir quiero/noquiero
    idea solucion: que la struct Jugada de la que todas extienden tenga un 
    metodo que checkeo eso:

# PERFORMANCE
- lo mas eficiente al momento de hacer las perspectivas es usar la misma struct
  y temporalmente setear algunas cartas a null (hacer benchmark comparando ambas)
- que es mas rapido para saber si un jugador tiene flor:
    * jugada.autor.tieneFlor(p.Ronda.Muestra)
    * o fijarse en un array de jugadores con flor y la op contains
    * un mapa de booleanos o similar

# DUDAS
- caso en TestFixPanic: no deberia ganar la mano? a quien le toca ser pie en la 
    mano 2?
- no necesariamente tiene que ser RETRUCOQUERIDO el estado del truco para recien
     ahi gritar vale4, podes ir de 1???
- en el test TestFixNacho la primera mano queda empardad:
    a quien le toca el sig turno? a richard o a andres?
- es necesario que sea su turno para cantar retruco o el turno de uno de los de 
    mi equipo?
    OJO: Acutalmente en el re-truco: turno del equipo, en el truco: turno mio
- es necesario que sea tu turno para cantar envido?
- actualmente si hay 3 flores: 2 del equipo rojo, 1 del equipo azul; 
    todos cantan "flor" sin mas, y la flor del equipo azul es la mas alta ->
    se lleva el puntaje el equipo azul de las 3 flores (3+3+3) ; eso esta bien?
- duda:
    >A truco
    >el evndio esta primer
    *se juega el envido*
    ----> ahora el truco sigue en juego como si lo hubiese cantado A? o hay que 
    volver a cantarlo?
- tiene sentido que alguien cante envido incluso cuando ya tiro todas sus cartas?
    checkear la cond. yaTiroTodasSusCartas en casos envido/truco
- se esta jugando de a 2 (o mas) uno tiene flor, pero el otro se va al 
    mazo deberia de sumar los puntos de la flor a pesar de que no fue cantada
- no es necesario cantar los puntajes de las flores si todos los que 
    tienen flores son del mismo equipo. Es necesario cantar los puntajes?
    TestTodoTienenFlor
- es necesario que sea tu turno para cantar la flor? o es tipo irse al mazo????
- todos los escenarios posibles de flor de TestTirada1.
  obs: ahora si todos cantan flor de una (como esta) entonces se fuega simplemente
  "la flor"
  si uno dice no quiero ~ con flor me achico -> acarrea a todo el equipo
  deberia ser asi?

  p.Cmd("Richard flor")
  p.Cmd("Adolfo no-quiero") <-----
  << [ok] (ALL) : +12 puntos para el equipo Rojo por las flores
  p.Cmd("Renzo flor")
  p.Cmd("Alvaro flor")

  - [RELACIONADO] ^ que pasa si hay 3 con flor A,B y C; 
    * A canta flor, 
    * B dice contra flor, 
    * A dice quiero -> pero C nunca canto su flor
    la tiene en cuenta, a pesar de que nunca canto flor
  - no es necesario esperar a que sea tu turno para cantar truco no?
    antes en retruco:
    `casoII := trucoYaQuerido && unoDeMiEquipoQuizo && esTurnoDeMiEquipo`
    ahora
    `casoII := trucoYaQuerido && unoDeMiEquipoQuizo`