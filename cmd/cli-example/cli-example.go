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
	fmt.Printf("\n>> ")
	cmd, _ := reader.ReadString('\n')
	return strings.TrimSuffix(cmd, "\n")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func print(output []truco.Msg) {
	if output == nil {
		return
	}
	for _, msg := range output {
		fmt.Print(msg)
	}
}

// LogFile log file
type LogFile struct {
	path string
}

// Log str to logFile
func (lf LogFile) Write(str string) {
	f, err := os.OpenFile(lf.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(str)
}

func newLogFile() LogFile {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	logPath := "/home/juan/Workspace/_tmp/truco_logs/"
	logFile := logPath + timestamp + ".log"

	return LogFile{path: logFile}
}

func main() {
	logfile := newLogFile()

	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	logfile.Write(p.ToJSON())

	p.Print()
	output := p.Dispatch()
	print(output)

	for {
		cmd := readLn()
		logfile.Write(cmd)

		err := p.SetSigJugada(cmd)
		if err != nil {
			fmt.Println("<< " + err.Error())
		}

		p.Esperar()
		output := p.Dispatch()
		print(output)

		p.Print()

		if p.NoAcabada() {
			break
		}
	}

	fmt.Println("se fini")
}
