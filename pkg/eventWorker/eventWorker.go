package eventworker

import (
	"YadroTest/pkg/config"
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// Количество элементов в слайсе после split.
const (
	eventWithoutComment = 3
	eventWithComment    = 4
)

// GetEvents запускает горутину, которая читает построчно ивенты
// и записывает их в канал.
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

// ResultEvents выводит логи и считает итоговые результаты.
func ResultEvents(cfg *config.Config, outCh <-chan string) {
	competitors := make(map[string]*competotor)

	for event := range outCh {
		eventSlice := strings.Split(event, " ")

		eventID, err := strconv.Atoi(eventSlice[1])
		if err != nil {
			log.Fatal(err)
		}

		eventID-- // слайс начинается с нуля, а ID с единицы

		if len(eventSlice) == eventWithoutComment {
			fmt.Println(
				incomingEvents[eventID](competitors, cfg, eventSlice[0], eventSlice[2], ""),
			)
		}

		if len(eventSlice) == eventWithComment {
			fmt.Println(
				incomingEvents[eventID](competitors, cfg, eventSlice[0], eventSlice[2], eventSlice[3]),
			)
		}
	}

	outputResult(competitors)
}

// getCompetitior возвращает самого быстрого участника и его ID.
func getCompetitor(cm map[string]*competotor) (string, *competotor) {
	var id string

	comp := &competotor{scheduledFinishTime: time.Duration(math.MaxInt64)}

	for i, c := range cm {
		if comp.scheduledFinishTime > c.scheduledFinishTime {
			comp = c
			id = i
		}
	}

	delete(cm, id)

	return id, comp
}

// outputResult выводит результат соревнования.
func outputResult(cm map[string]*competotor) {
	fmt.Printf("\nResults:\n")

	for len(cm) != 0 {
		id, comp := getCompetitor(cm)

		if comp.laps == 0 && comp.firings == 0 {
			comp.mark = fmt.Sprintf("[%s/%s]", comp.actualFinishTime, comp.scheduledFinishTime)
		} else {
			comp.mark = "[" + comp.mark + "]"
		}

		fmt.Printf("%s %s [", comp.mark, id)

		for i := range comp.ml {
			fmt.Printf("{%v, %f}", comp.ml[i].timeCompleteLap, comp.ml[i].averageSpeedLap)

			if i+1 != len(comp.ml) {
				fmt.Printf(", ")
			}
		}

		fmt.Printf("] [")

		for i := range comp.pl {
			fmt.Printf("{%v, %f}", comp.pl[i].timeCompleteLap, comp.pl[i].averageSpeedLap)

			if i+1 != len(comp.pl) {
				fmt.Printf(", ")
			}
		}

		fmt.Printf("] [")

		for i, hits := range comp.numberOfShots {
			fmt.Printf("%d/5", hits)

			if i+1 != len(comp.numberOfShots) {
				fmt.Printf(", ")
			}
		}

		fmt.Println("]")
	}
}
