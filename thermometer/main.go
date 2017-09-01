package main

import (
	"fmt"
	"os"
)

var exit chan int
var config ThermometerConfig

func main() {
	initUI()
	defer closeUI()

	config = ThermometerConfig{}

	exit = make(chan int)

	initGpio()

	if _, err := os.Stat("porkmeterconfig.yaml"); err != nil {
		fmt.Println("start configuration")
		go createNoConfigUI()
	} else {
		fmt.Println("read configuration")
		readConfiguration()
		go createMenuUI()
	}

	go timeLoop()

	<- exit
}