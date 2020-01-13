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

func main() {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	logPath := "/tmp/"
	logFile := logPath + timestamp + ".log"

	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	_log(p.ToJSON(), logFile)

	p.Print()
	for p.NoAcabada() {
		fmt.Printf("\n>> ")
		cmd := readLn()
		_log(cmd, logFile)
		res := p.SetSigJugada(cmd)
		p.Esperar()
		p.Print()
		if res != nil {
			fmt.Println(res)
		}
	}
}
