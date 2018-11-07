package simulation

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

func getTestingMode(tb testing.TB) (testingMode bool, t *testing.T, b *testing.B) {
	testingMode = false
	if _t, ok := tb.(*testing.T); ok {
		t = _t
		testingMode = true
	} else {
		b = tb.(*testing.B)
	}
	return
}

// Pretty-print events as a table
func DisplayEvents(events map[string]uint) {
	var keys []string
	for key := range events {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	fmt.Printf("Event statistics: \n")
	for _, key := range keys {
		fmt.Printf("  % 60s => %d\n", key, events[key])
	}
}

// Builds a function to add logs for this particular block
func addLogMessage(testingmode bool, blockLogBuilders []*strings.Builder, height int) func(string) {
	if testingmode {
		blockLogBuilders[height] = &strings.Builder{}
		return func(x string) {
			(*blockLogBuilders[height]).WriteString(x)
			(*blockLogBuilders[height]).WriteString("\n")
		}
	}
	return func(x string) {}
}

// Creates a function to print out the logs
func logPrinter(testingmode bool, logs []*strings.Builder) func() {
	if testingmode {
		return func() {
			numLoggers := 0
			for i := 0; i < len(logs); i++ {
				// We're passed the last created block
				if logs[i] == nil {
					numLoggers = i
					break
				}
			}
			var f *os.File
			if numLoggers > 10 {

				fileName := fmt.Sprintf("simulation_log_%s.txt",
					time.Now().Format("2006-01-02 15:04:05"))

				fmt.Printf("Too many logs to display, instead writing to %s\n", fileName)
				f, _ = os.Create(fileName)
			}
			for i := 0; i < numLoggers; i++ {
				if f != nil {
					_, err := f.WriteString(fmt.Sprintf("Begin block %d\n", i+1))
					if err != nil {
						panic("Failed to write logs to file")
					}
					_, err = f.WriteString((*logs[i]).String())
					if err != nil {
						panic("Failed to write logs to file")
					}
				} else {
					fmt.Printf("Begin block %d\n", i+1)
					fmt.Println((*logs[i]).String())
				}
			}
		}
	}
	return func() {}
}
