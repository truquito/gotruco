# TODO
  **cli CMD**
    - cli: mini commandos con menu: t -> truco, e -> envido fe -> falta envido
    - los mensajes de info en realidad actualizan el resultado de la mano
    - la suma de puntos por envido ganado -> deberia deshabilitar el envido en 
      el cli


# BUGS


# LIMPIEZA DE CODIGO
  
- el eval de la flor esta dividido: el eval solo se llama cuando se juega flor comun
  mientras que los quiero/noquiero se llaman desde ahi mismo
  para el eval de contraflor/contrafloralresto no hay case (ya que se implementan
  en los quiero/noquiero)
- Hay codigo repetido entre noQuiero y mazo cuando niega la flor (codigo copiado)
- hay redundancia entre cantarFloresSiLasHay y cantarFlores
- hay redundancia entre p.Ronda.GetLaFlorMasAlta y 
    manojoConLaFlorGanadora, _, _ := p.Ronda.execCantarFlores(aPartirDe)
- ??? esta duplicado el codigo de irse al mazo con el de "no quiero"

# SEGURIDAD
- que en los writes se use el id no el nombre ni el nick etc...
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
- lo mas eficiente al momento de hacer las PerspectivaCacheFlors es usar la misma struct
  y temporalmente setear algunas cartas a null (hacer benchmark comparando ambas)

# DUDAS

- puede un mismo tipo decir "truco" y luego decir el mismo "envido esta primero", medio bobo
  pero es parte del protocolo (?)
- caso muy dudoso: TestFlorFlorContraFlorQuiero
- que pasa si uno tiene flor, no es su turno y el otro se va al mazo y no le da la chance
  ganar los +3 pts de la flor. Deberia ser posible?
- estamos jugando de a 2, yo tiro un 10 de la muestra (perica), y mi oponente luego
  tira un 11 de la muestra (perico). entonces:
    * la mano resulta parda? o la gano el del perico?
    * el siguiente turno, de quien es? porque ambas tienen el mismo valor (27)
      pero el perico le gana a la perica entonces es turno del oponente?
      o es mi turno porque tienen el mismo valor pero yo tire primero? (le gano de "mano")
      me decante por: no es parda, le gana el perico en todo.

**indeps**:
- no necesariamente tiene que ser RETRUCOQUERIDO el estado del truco para recien
     ahi gritar vale4, podes ir de 1???
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
  solo lo puede cantar en la primera mano no?
  tiene que ser necesariamente antes de tirar SU primera carta?
- se esta jugando de a 2 (o mas) uno tiene flor, pero el otro se va al 
    mazo deberia de sumar los puntos de la flor a pesar de que no fue cantada
- no es necesario cantar los puntajes de las flores si todos los que 
    tienen flores son del mismo equipo. Es necesario cantar los puntajes?
    TestTodoTienenFlor
- es necesario que sea tu turno para cantar la flor? o es tipo irse al mazo????

[resuelto? necesariamente deben cantar flor si es que tienen]
- todos los escenarios posibles de flor de TestTirada1.
  obs: ahora si todos cantan flor de una (como esta) entonces se juega simplemente
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

- actualmente se le puede responder quiero a una contraflor??
- que pasa si en una mano, 2 de equipo contrario tiran cartas que empardan
  pero ambos se van al mazo antes de que termine la mano.
  de quien es el siguiente turno?
- que pasa si el que gana la mano se fue al mazo antes?


**simulacion**:
- caso en TestFixPanic: no deberia ganar la mano? a quien le toca ser pie en la 
    mano 2?
- en el test TestFixNacho la Primera mano queda empardad:
    a quien le toca el sig turno? a richard o a andres?
- @ flor_test.TestFixContraFlor:
      flor -> flor -> 3 + 3 = 6 pts en total
      pero,
      flor -> contraflor -> quiero -> 1 + 3 = 4 pts en total MENOS QUE FLOR COMUN, WTF?
    no deberia ser 3 + 3 + x ?
- @ partida_test.TestPartida1 que pasa con richard que tiene flor y no la canta?
  que deberia pasar en esos casos?






