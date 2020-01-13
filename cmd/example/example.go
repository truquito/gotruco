package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

func logPartida(p *truco.Partida) {
	json := p.ToJSON() + "\n"
	content := []byte(json)
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	fileName := "/tmp/p" + timestamp + ".txt"
	err := ioutil.WriteFile(fileName, content, 0644) // sobreescribe
	check(err)
}

func main() {
	p, _ := truco.NuevaPartida(20, []string{"Alvaro", "Adolfo", "Andres"}, []string{"Roro", "Renzo", "Richard"})
	logPartida(p)
	p.Print()

	for p.NoAcabada() {
		fmt.Printf("\n>> ")
		cmd := readLn()
		res := p.SetSigJugada(cmd)
		p.Esperar()
		p.Print()
		if res != nil {
			fmt.Println(res)
		}
	}
}
