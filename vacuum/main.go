package main

import (
	"flag"
	"fmt"
)

func main() {
	var configFile, algorithm string
	var animate, cat, isHouse, useLogic bool

	flag.StringVar(&configFile, "file", "empty.json", "configuration file")
	flag.StringVar(&algorithm, "algorithm", "snake", "cleaning algorithm")
  flag.BoolVar(&animate, "animate", true, "animate while cleaning")
  flag.BoolVar(&cat, "cat", false, "add a cat to the room")
  flag.BoolVar(&isHouse, "house", false, "config file has multiple rooms")
  flag.BoolVar(&useLogic, "logic", false, "use propositional logic for cleaning decisions")
  flag.Parse()

  var house *House

  if !isHouse {
    /*If not a house, we just have one room. Create a house, and assign one room to it
    This way, we can use the same loop for house and for individual rooms*/
    var rooms []*Room

    // Get a room from JSON config
    room := NewRoom(configFile, animate)
    rooms = append(rooms, room)
    var h House
    h.Rooms = rooms
    house = &h
  }else {
    // We are doing a complete house. Just get a house from JSON config
    house = NewHouse(configFile, animate)
  }

  roomCount := 0

  if useLogic {
    // Use propositional logic for cleaning
    // TODO robot := NewRobotWithLogic(1,1)
  }else {
    // Use the original cleaning approach without propositional logic, and for multiple rooms
    for _, room := range house.Rooms {
      if cat {
      room.Cat = NewCat(room)
      }

      // Get a robot
      robot := NewRobot(1,1)
      // Assign a clean algorithm
      switch algorithm {
      case "random":
        robot.CleanRoom = CleanRoomRandomWalk
      case "slam":
        robot.CleanRoom = CleanRoomSlam
      case "spiral":
        robot.CleanRoom = CleanSpiralPattern
      default:
        robot.CleanRoom = CleanRoomSnake
      }
      // Clean the room
      robot.CleanRoom(room, robot)
      roomCount++
    }
  }
  fmt.Printf("All done. Clean a total of %d room(s)\n", roomCount)
}