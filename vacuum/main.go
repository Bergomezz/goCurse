package main

import (
	"flag"
	"fmt"
)

func main() {
	var configFile, algorithm string
	var animate bool

	flag.StringVar(&configFile, "file", "empty.json", "configuration file")
	flag.StringVar(&algorithm, "algorithm", "random", "cleaning algorithm")
  flag.BoolVar(&animate, "animate", true, "animate while cleaning")
  flag.Parse()

  room := NewRoom(configFile, animate)
  fmt.Println(room.CleanableCellCount)
}