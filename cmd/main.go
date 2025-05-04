package main

import (
	"YadroTest/pkg/config"
	eventworker "YadroTest/pkg/eventWorker"
)

func main() {
	cfg := config.MustLoad("config/config.json")
	reader := eventworker.GetEvents("sunny_5_skiers/events")

	eventworker.ResultEvents(cfg, reader)
}
