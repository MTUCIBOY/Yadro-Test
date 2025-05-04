package eventworker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	eventWithoutComment = 3
	eventWithComment    = 4
)

func GetEvents(filepath string) chan string {
	outCh := make(chan string)

	go func() {
		defer close(outCh)

		file, err := os.Open(filepath)
		if err != nil {
			return
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			outCh <- scanner.Text()
		}
	}()

	return outCh
}

func OutputEvents(outCh <-chan string) {
	for event := range outCh {
		eventSlice := strings.Split(event, " ")

		if len(eventSlice) == eventWithoutComment {
			fmt.Println(incomingEvents[eventSlice[1]](eventSlice[0], eventSlice[2], ""))
		}

		if len(eventSlice) == eventWithComment {
			fmt.Println(incomingEvents[eventSlice[1]](eventSlice[0], eventSlice[2], eventSlice[3]))
		}
	}
}
