package main

import (
	"YadroTest/pkg/config"
	"fmt"
)

func main() {
	fmt.Println(config.MustLoad("config/config.json"))
}
