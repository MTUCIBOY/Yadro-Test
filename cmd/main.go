package main

import (
	"YadroTest/pkg/config"
	eventworker "YadroTest/pkg/eventWorker"
	"fmt"
)

func main() {
	fmt.Println(config.MustLoad("config/config.json"))

	eventworker.OutputEvents(eventworker.GetEvents("sunny_5_skiers/events"))
}
