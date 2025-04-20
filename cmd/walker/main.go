package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

// Checkpoint structure to save progress
type Checkpoint struct {
	Terminals uint      `json:"terminals"`
	Queue     []string  `json:"queue"`      // Serialized game states to process
	StartTime time.Time `json:"start_time"` // When this run started
}

// Global variables
var (
	terminals      uint = 0
	checkpointFile string
	timeLimit      time.Duration
	lastCheckTime  time.Time
	checkpoint     Checkpoint
	checkInterval  = 5 * time.Second // How often to check if we need to save a checkpoint
)

func countMsgs(pkts []enco.Envelope) (int, int) {
	rondaGanadaCount, nuevaRondaCount := 0, 0
	for _, pkt := range pkts {
		if pkt.Message.Cod() == enco.TRondaGanada {
			rondaGanadaCount++
		} else if pkt.Message.Cod() == enco.TNuevaRonda {
			nuevaRondaCount++
		}
	}
	return rondaGanadaCount, nuevaRondaCount
}

// Check if we need to create a checkpoint
func shouldCheckpoint() bool {
	// Check time elapsed at regular intervals to reduce overhead
	if time.Since(lastCheckTime) < checkInterval {
		return false
	}

	lastCheckTime = time.Now()
	return time.Since(checkpoint.StartTime) >= timeLimit
}

// Save the current state as a checkpoint
func saveCheckpoint(gameQueue []string) error {
	checkpoint.Terminals = terminals
	checkpoint.Queue = gameQueue

	// Create checkpoint file
	file, err := os.Create(checkpointFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(checkpoint)
}

// Load checkpoint from file
func loadCheckpoint() (bool, error) {
	file, err := os.Open(checkpointFile)
	if err != nil {
		if os.IsNotExist(err) {
			// No checkpoint file exists, that's ok
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&checkpoint)
	if err != nil {
		return false, err
	}

	// Restore global terminal count
	terminals = checkpoint.Terminals

	return true, nil
}

// Non-recursive version of the game tree traversal
func processGameTree() bool {
	gameQueue := make([]string, 0)

	// If we're starting from scratch, initialize with a new game
	if len(checkpoint.Queue) == 0 {
		// Initialize new execution
		// partidaJSON := `{"limiteEnvido":1,"cantJugadores":2,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"oro","valor":5},{"palo":"espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":1},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]}}` // 2p
		partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"basto","valor":1},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"espada","valor":7},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Bob","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":12},{"palo":"basto","valor":4},{"palo":"copa","valor":4}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Ariana","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":10},{"palo":"espada","valor":1},{"palo":"oro","valor":5}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Ben","equipo":"rojo"}}],"mixs":{"Alice":0,"Ariana":2,"Ben":3,"Bob":1},"muestra":{"palo":"oro","valor":2},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]},"limiteEnvido":1}` // 4p
		p, err := pdt.Parse(partidaJSON, true)
		// p, err := pdt.NuevaPartida(pdt.A20, []string{"Alice"}, []string{"Bob"}, 1, true)

		if err != nil {
			panic(err)
		}
		bs, _ := p.MarshalJSON()
		gameQueue = append(gameQueue, string(bs))
	} else {
		// Use the checkpointed queue
		gameQueue = append(gameQueue, checkpoint.Queue...)
	}

	// Process the queue in a loop instead of recursively
	for len(gameQueue) > 0 {
		// Periodically check if we need to save a checkpoint
		if shouldCheckpoint() {
			fmt.Printf("Time limit reached after processing %d terminals. Saving checkpoint...\n", terminals)
			err := saveCheckpoint(gameQueue)
			if err != nil {
				fmt.Printf("Error saving checkpoint: %v\n", err)
			}
			return false // Exit the function to stop processing
		}

		// Get the next game state from the queue
		gameState := gameQueue[0]
		gameQueue = gameQueue[1:] // Remove the processed state

		// Parse the game state
		p, err := pdt.Parse(gameState, true)
		if err != nil {
			panic(err)
		}

		// Get all possible plays for the current game state
		chis := pdt.Chis(p)

		// Try each possible play
		for mix := range chis {
			for aix := range chis[mix] {
				// Restore the game state for each play
				p, _ = pdt.Parse(gameState, true)
				a := chis[mix][aix]
				pkts, err := p.Cmd(a.String())

				rg, nr := countMsgs(pkts)
				if (!p.Terminada() && nr > 0 && rg == 0) || (err != nil) {
					fmt.Println(gameState)
					fmt.Println(a.String())
					panic(123)
				}

				isDone := rg > 0
				// isDone := p.Terminada() // use this only if puntaje = puntuacion - 1

				if isDone {
					terminals++

					// Periodically print progress
					if terminals%100_000 == 0 {
						fmt.Printf("Terminals processed: %d\n", terminals)
					}
				} else {
					// Add the new game state to the queue for processing
					newState, _ := p.MarshalJSON()
					gameQueue = append(gameQueue, string(newState))
				}
			}
		}
	}

	// Update the global checkpoint queue to reflect completion
	checkpoint.Queue = gameQueue
	return true
}

func main() {
	// Parse command line flags
	flag.StringVar(&checkpointFile, "checkpoint", "game_checkpoint.json", "Checkpoint file path")
	timeLimitSeconds := flag.Int("timelimit", 259200, "Time limit in seconds (default: 3 days)")
	flag.Parse()

	// Convert time limit to duration
	timeLimit = time.Duration(*timeLimitSeconds) * time.Second

	// Initialize checkpoint time tracking
	lastCheckTime = time.Now()

	// Try to load checkpoint
	loaded, err := loadCheckpoint()
	if err != nil {
		fmt.Printf("Error loading checkpoint: %v\n", err)
		os.Exit(1)
	}

	// Set the start time for this run
	checkpoint.StartTime = time.Now()

	if loaded {
		fmt.Printf("Resuming from checkpoint. Terminals counted so far: %d\n", terminals)
		fmt.Printf("Time limit for this run: %v\n", timeLimit)
	} else {
		fmt.Println("Starting new run!")
	}

	// Process the game tree
	fullyTraversed := processGameTree()

	fmt.Printf("Checkpoint statistics:\n")
	fmt.Printf("  Terminals processed: %d\n", terminals)
	fmt.Printf("  Remaining states in queue: %d\n", len(checkpoint.Queue))
	fmt.Printf("  Time elapsed in this run: %v\n", time.Since(checkpoint.StartTime))

	// If we get here, it means we've completed the entire tree traversal
	if fullyTraversed {
		fmt.Println("Processing complete!")
		os.Exit(0)
	}

	os.Exit(1)
}
