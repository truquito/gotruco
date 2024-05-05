package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"github.com/truquito/gotruco/deco"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/truco"
)

// flags
var (
	timeoutFlag    = flag.Int("timeout", 60, "Timeout per turn (in seconds)")
	numPlayersFlag = flag.Int("n", 2, "Number of players (2, 4 or 6)")
	logFileFlag    = flag.String("log_file", "", "Path to the log file")
)

// vars
var (
	logFile *os.File    = nil
	logger  *log.Logger = nil
	ioCh    chan string = make(chan string, 1)
)

func currentDatetime() string {
	now := time.Now().UTC()
	return fmt.Sprintf("%d-%02d-%02d_%02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

func initLogs(logFilePath string) {
	if len(logFilePath) == 0 {
		logFilePath = fmt.Sprintf("/tmp/truco-%s.log", currentDatetime())
	}
	var (
		flags int         = os.O_CREATE | os.O_WRONLY | os.O_APPEND
		perms fs.FileMode = 0666
		err   error       = nil
	)

	logFile, err = os.OpenFile(logFilePath, flags, perms)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger = log.New(logFile, "", log.LstdFlags)
}

func handleIO() {
	reader := bufio.NewReader(os.Stdin)
	readLn := func(prefix string) string {
		fmt.Print(prefix)
		cmd, _ := reader.ReadString('\n')
		return strings.TrimSuffix(cmd, "\n")
	}

	for {
		cmd := readLn("")
		ioCh <- cmd
	}
}

func init() {
	flag.Parse()
	initLogs(*logFileFlag)
}

func main() {
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	azules := []string{"Alice", "Ariana", "Annie"}
	rojos := []string{"Bob", "Ben", "Bill"}
	timeout := time.Duration(*timeoutFlag) * time.Second
	p, _ := truco.NuevoJuego(
		20,
		azules[:*numPlayersFlag>>1],
		rojos[:*numPlayersFlag>>1],
		4,
		true,
		timeout)

	pJSON, _ := p.MarshalJSON()
	logger.Println(string(pJSON))

	for _, m := range p.Ronda.Manojos {
		pers, _ := p.Perspectiva(m.Jugador.ID)
		pJSON, _ := pers.MarshalJSON()
		logger.Println(m.Jugador.ID, string(pJSON))
	}

	fmt.Println(p)

	for _, pkt := range p.Consumir() {
		fmt.Println(deco.Stringify(&pkt, p.Partida))
	}

	// hago una gorutine (y channel para avisar) para el io
	go handleIO()

	for {
		select {
		// canal de entrada del usuario
		case cmd := <-ioCh:
			if cmd == "dump" {
				data, _ := json.Marshal(p)
				fmt.Println(string(data))
			} else {
				logger.Println(cmd)
				err := p.Cmd(cmd)
				if err != nil {
					fmt.Println("<< " + err.Error())
				}
				for _, pkt := range p.Consumir() {
					fmt.Println(deco.Stringify(&pkt, p.Partida))
				}
				fmt.Println(p)
			}
		// canal de error detectado por parte del simulador
		case <-p.ErrCh:
			// el error deberia estar aca
			for _, pkt := range p.Consumir() {
				fmt.Println(pkt.Message.Cod(), deco.Stringify(&pkt, p.Partida))
			}
			// de momento, el unico error posible
			if p.Expirado() {
				m, _ := p.Err.Message.(enco.TimeOut)
				fmt.Printf(
					"el juego terminó debido a que `%s` no realizó niguna jugada en %s.\n",
					m,
					p.DurTurno)
			}
			// fmt.Printf(">> ") // prompt
		}

		if p.Terminado() {
			return
			// si es modo bucle, entonces que no salga del for sino que
			// cree un juego nuevo
		}
	}

}
