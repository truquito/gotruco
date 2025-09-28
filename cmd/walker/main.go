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
	Stack     []string  `json:"stack"`
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

// Save the current state (stack) as a checkpoint
func saveCheckpoint(gameStack []string) error {
	checkpoint.Terminals = terminals
	checkpoint.Stack = gameStack

	// Create checkpoint file
	file, err := os.Create(checkpointFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Use spaces for indentation
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
		// Handle potential empty file or other decode errors
		if err.Error() == "EOF" {
			fmt.Println("Checkpoint file is empty or corrupted, starting new run.")
			return false, nil
		}
		return false, err
	}

	// Restore global terminal count
	terminals = checkpoint.Terminals

	// Check if Stack is nil (might happen if saved before any processing)
	if checkpoint.Stack == nil {
		checkpoint.Stack = make([]string, 0)
	}

	return true, nil
}

// Iterative Depth-First Search (DFS) using a stack
func processGameTree() bool {
	// Use a slice as a stack (LIFO)
	gameStack := make([]string, 0, 10000) // Initialize with some capacity

	// If we're starting from scratch, initialize with a new game
	if len(checkpoint.Stack) == 0 {

		// Initialize new execution
		// partidaJSON := `{"limiteEnvido":1,"cantJugadores":2,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"oro","valor":3},{"palo":"copa","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alvaro","nombre":"Alvaro","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":3},{"palo":"oro","valor":5},{"palo":"espada","valor":2}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Roro","nombre":"Roro","equipo":"rojo"}}],"muestra":{"palo":"copa","valor":1},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]}}` // 2p
		partidaJSON := `{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"basto","valor":1},{"palo":"copa","valor":6}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"espada","valor":7},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Bob","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":12},{"palo":"basto","valor":4},{"palo":"copa","valor":4}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Ariana","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":10},{"palo":"espada","valor":1},{"palo":"oro","valor":5}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Ben","equipo":"rojo"}}],"mixs":{"Alice":0,"Ariana":2,"Ben":3,"Bob":1},"muestra":{"palo":"oro","valor":2},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]},"limiteEnvido":1}` // 4p
		p, err := pdt.Parse(partidaJSON, true)
		// p, err := pdt.NuevaPartida(pdt.A20, []string{"Alice"}, []string{"Bob"}, 1, true)

		if err != nil {
			// Consider more specific error handling if NuevaPartida can fail
			panic(fmt.Errorf("failed to create initial game: %w", err))
		}
		bs, err := p.MarshalJSON()
		if err != nil {
			panic(fmt.Errorf("failed to marshal initial game state: %w", err))
		}
		gameStack = append(gameStack, string(bs))
	} else {
		// Use the checkpointed stack
		fmt.Printf("Resuming with %d states on the stack.\n", len(checkpoint.Stack))
		gameStack = append(gameStack, checkpoint.Stack...)
		// Clear the checkpoint stack immediately after loading to potentially free memory earlier
		checkpoint.Stack = nil
	}

	// Process the stack in a loop instead of recursively
	for len(gameStack) > 0 {
		// Periodically check if we need to save a checkpoint
		if shouldCheckpoint() {
			fmt.Printf("Time limit reached after processing %d terminals. Saving checkpoint...\n", terminals)
			err := saveCheckpoint(gameStack)
			if err != nil {
				fmt.Printf("Error saving checkpoint: %v\n", err)
				// Decide whether to exit or continue despite checkpoint failure
			} else {
				fmt.Println("Checkpoint saved successfully.")
			}
			return false // Exit the function to stop processing for this run
		}

		// Pop the next game state from the stack (LIFO)
		stackLen := len(gameStack)
		gameState := gameStack[stackLen-1]
		gameStack = gameStack[:stackLen-1] // Efficiently remove the last element

		// Parse the game state
		p, err := pdt.Parse(gameState, true)
		if err != nil {
			// Consider logging the problematic gameState string
			fmt.Printf("Error parsing game state: %v\nState: %s\n", err, gameState)
			// Decide whether to skip this state or panic
			panic(fmt.Errorf("failed to parse game state: %w", err))
		}

		// Get all possible plays for the current game state
		chis := pdt.Chis(p)

		// Try each possible play
		for mix := range chis {
			for aix := range chis[mix] {
				// Optimization: Avoid re-parsing the same state repeatedly
				// Create a copy *before* applying the action
				bs, _ := p.MarshalJSON()
				pCopy, _ := pdt.Parse(string(bs), true)
				// If not, you'd need to Parse again as before:
				// pCopy, _ = pdt.Parse(gameState, true)

				a := chis[mix][aix]
				pkts, err := pCopy.Cmd(a.String()) // Apply command to the copy

				// Error Handling & Sanity Check
				// It seems you had a specific condition you were checking.
				// Let's refine the check based on your original code's apparent intent:
				// Check if the game is NOT terminated BUT a new round started without winning the previous one, OR if Cmd returned an error.
				rg, nr := countMsgs(pkts)
				if err != nil {
					fmt.Printf("Error executing command '%s' on state: %s\nError: %v\n", a.String(), gameState, err)
					// Decide how to handle command errors (e.g., log, skip, panic)
					continue // Skip this invalid action/state
				}
				// This check seems specific to Truco rules - adjust if needed
				if !pCopy.Terminada() && nr > 0 && rg == 0 {
					fmt.Println("Inconsistent state detected (nueva ronda without ronda ganada before game end):")
					fmt.Println("Parent State:", gameState)
					fmt.Println("Action:", a.String())
					// Maybe log more details about pCopy's state here
					panic("Inconsistent game state transition detected")
				}
				// End Error Handling

				isDone := rg > 0 // Assuming rg > 0 means the game ended as per your logic
				// isDone := pCopy.Terminada() // Alternative based on game state method

				if isDone {
					terminals++

					// Periodically print progress
					if terminals%100_000 == 0 {
						fmt.Printf("Terminals processed: %d | Stack size: %d\n", terminals, len(gameStack))
					}
				} else {
					// Push the new game state onto the stack for processing (DFS)
					newState, err := pCopy.MarshalJSON()
					if err != nil {
						fmt.Printf("Error marshaling new game state after action '%s': %v\n", a.String(), err)
						// Decide how to handle marshaling errors
						continue // Skip adding this state
					}
					gameStack = append(gameStack, string(newState)) // Push onto stack
				}
			}
		}
	}

	// If the loop finishes, it means the stack is empty and the entire tree has been processed.
	// Ensure the checkpoint reflects this completion state (empty stack).
	checkpoint.Stack = []string{} // Explicitly set to empty
	return true
}

func main() {
	// Parse command line flags
	flag.StringVar(&checkpointFile, "checkpoint", "game_checkpoint_dfs.json", "Checkpoint file path (using DFS)") // Changed default name
	timeLimitSeconds := flag.Int("timelimit", 259200, "Time limit in seconds (default: 3 days)")
	flag.Parse()

	// Convert time limit to duration
	timeLimit = time.Duration(*timeLimitSeconds) * time.Second

	// Initialize checkpoint time tracking
	lastCheckTime = time.Now() // Initialize here

	// Try to load checkpoint
	loaded, err := loadCheckpoint()
	if err != nil {
		fmt.Printf("Error loading checkpoint: %v\n", err)
		os.Exit(1)
	}

	// Set the start time for this run AFTER potential loading
	// If loaded, this run's elapsed time starts now. If not loaded, same.
	checkpoint.StartTime = time.Now()
	// Reset lastCheckTime as well, so the first shouldCheckpoint check is accurate for this run
	lastCheckTime = checkpoint.StartTime

	if loaded {
		fmt.Printf("Resuming from checkpoint. Terminals counted so far: %d\n", terminals)
		fmt.Printf("Time limit for this run: %v\n", timeLimit)
		// Stack count printed within processGameTree when resuming
	} else {
		fmt.Println("Starting new run (DFS)!")
		// Ensure checkpoint struct is initialized if not loading
		checkpoint = Checkpoint{
			Terminals: 0,
			Stack:     make([]string, 0), // Initialize stack
			StartTime: time.Now(),
		}
	}

	// Process the game tree using DFS
	fullyTraversed := processGameTree()

	// Final stats output
	fmt.Printf("\n--- Run Finished ---\n")
	fmt.Printf("Total Terminals processed (cumulative): %d\n", terminals)
	fmt.Printf("Remaining states in stack (should be 0 if fully traversed): %d\n", len(checkpoint.Stack)) // Checkpoint stack should be empty if finished
	fmt.Printf("Time elapsed in this run: %v\n", time.Since(checkpoint.StartTime))

	if fullyTraversed {
		fmt.Println("Processing complete! Entire game tree traversed.")
		// Optional: Clean up checkpoint file upon successful completion
		// err := os.Remove(checkpointFile)
		// if err != nil && !os.IsNotExist(err) {
		//	 fmt.Printf("Warning: could not remove checkpoint file: %v\n", err)
		// }
		os.Exit(0)
	} else {
		fmt.Println("Processing stopped (likely due to time limit). State saved in checkpoint.")
		os.Exit(1) // Exit code 1 indicates incomplete processing in this run
	}
}
