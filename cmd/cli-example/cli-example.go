package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jpfilevich/truco"
)

var reader = bufio.NewReader(os.Stdin)

func readLn() string {
	cmd, _ := reader.ReadString('\n')
	return strings.TrimSuffix(cmd, "\n")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func _log(str, logFile string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(str)
}

func print(output []truco.Msg) {
	if output == nil {
		return
	}
	for _, msg := range output {
		fmt.Print(msg)
	}
}

func main() {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	logPath := "/home/juan/Workspace/_tmp/truco_logs/"
	logFile := logPath + timestamp + ".log"

	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	_log(p.ToJSON(), logFile)

	p.Print()
	output := p.Dispatch()
	print(output)
	for p.NoAcabada() {
		fmt.Printf("\n>> ")
		cmd := readLn()
		_log(cmd, logFile)
		// tuvo error?
		err := p.SetSigJugada(cmd)
		if err != nil {
			fmt.Println("<< " + err.Error())
		}
		p.Esperar()
		// obtengo un array de mensajes como output
		output := p.Dispatch()
		print(output)
		p.Print()
	}
}
