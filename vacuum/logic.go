package main

import (
	"fmt"
	"time"
)

type PersonStatus struct {
	Name      string
	IsHome    bool
	Room      string // Person's room name
	DoorClose bool
}

type LogicalWorld struct {
	Jack      PersonStatus
	Sarah     PersonStatus
	Johnny    PersonStatus
	IsWeekDay bool
	Objects    map[string]bool
}

func NewLogicWorld() *LogicalWorld {
	// Get current day to determine if it's a week day
	today := time.Now()
  weekday := today.Weekday()

  isWeekDay := weekday >= time.Monday && weekday <= time.Friday

  return &LogicalWorld{
    Jack: PersonStatus{
      Name: "Jack",
      IsHome: false,
      Room: "Jack's Room",
    },
    Sarah: PersonStatus{
      Name: "Sarah",
      IsHome: false,
      Room: "Sarah's Room",
    },
    Johnny: PersonStatus{
      Name: "Johnny",
      IsHome: false,
      Room: "Johnny's Room",
      DoorClose: false,
    },
    IsWeekDay: isWeekDay,
    Objects: make(map[string]bool),
  }
}

func (world *LogicalWorld) UpdateObjectFound(objectName string) {
  world.Objects[objectName] = true

  // Apply object to a person identification rules

  // Rule 1: If a backpack is found, then Jack is home
  if objectName == "backpack" {
    world.Jack.IsHome = true
    fmt.Println("Logic: Backpack found, deducing Jack is home")
  }

  // Rule 2: If a bicycle is found, then Sarah is home
  if objectName == "bicycle" {
    world.Sarah.IsHome = true
    fmt.Println("Logic: Bicycle found, deducing Sarah is home")
  }

  // Rule 3: If a skateboard is found, then Johnny is home
  if objectName == "skateboard" {
    world.Johnny.IsHome = true
    fmt.Println("Logic: Skateboard found, deducing Johnny is home")
  }
}

// UpdateDoorsStatus update whether Johnny's door is open or close
func (world *LogicalWorld) UpdateDoorsStatus(doorName string, isClose bool) {
  if doorName == "Johnny's Door" {
    world.Johnny.DoorClose = true
    fmt.Printf("Logic: Johnny's door is now: %s\n", map[bool]string{true: "close", false: "open"}[isClose])
  }
}
// DetermineCleaningPriority decide on the order we clean rooms
func (world *LogicalWorld) DetermineCleaningPriority() []string {
  availableRooms := []string{
    "Kitchen",
    "Living Room",
    "Jack's Room",
    "Sarah's Room",
    "Johnny's Room",
  }

  // Rule: If no one is home, then vacuum all rooms starting with the kitchen.
  if !world.Jack.IsHome && !world.Sarah.IsHome && !world.Johnny.IsHome {
    fmt.Println("Logic: No one is home, vacuuming all rooms starting with the kitchen")
    return availableRooms
  }

  // Initialize a priority list with all available rooms
  priorityList := make([]string,0)
  skipRooms := make(map[string]bool)

  // Rule: If Sarah is home, then don't vacuum the living room
  if world.Sarah.IsHome {
    fmt.Println("Logic: Sarah is home, skipping the living room")
    skipRooms["Living Room"] = true
  }

  // Rule: If Johnny is home and his door is closed, then skip Johnny's room
  if world.Johnny.IsHome && world.Johnny.DoorClose {
    fmt.Println("Logic: Johnny is home and his door is closed, skip his room")
    skipRooms["Johnny's Room"] = true
  }

  // Rule: If Jack is home and it's weekday, do his room last
  jackRoomLast := world.Jack.IsHome && world.IsWeekDay

  // Build a priority list. Add kitchen first, and filter out skipped rooms
  priorityList = append(priorityList, "Kitchen")

  // Add all other room's except Jack's (if it's need to be last) and skipped rooms
  for _, room := range availableRooms {
    if room == "Kitchen" {
      continue // Already added
    }
    if room == "Jack's Room" && jackRoomLast {
      continue // Will be added last
    }
    if skipRooms[room] {
      continue // Skip this room
    }
    priorityList = append(priorityList, room)
  }

  // Add Jack's room if should be last
  if jackRoomLast && !skipRooms["Jack's Room"] {
    priorityList = append(priorityList, "Jack's Room")
  }

  return priorityList
}

// RobotWithLogic is type for a robot with logic
type RobotWithLogic struct {
  *Robot // Embedding the original robot
  World *LogicalWorld
}

// NewRobotWithLogic factory method for a robot NewRobotWithLogic
func NewRobotWithLogic(startX, startY int) *RobotWithLogic {
  return &RobotWithLogic{
    Robot: NewRobot(startX, startY),
    World: NewLogicWorld(),
  }
}

func (robot *RobotWithLogic) ScanHouseWithLogic(house *House) map[string] int {
  // Create a map which maps rooms indices with rooms name
  roomNameToIndex := make(map[string]int)

  // Identify all rooms and generate our mapping
  for i, room := range house.Rooms {
    roomName := ""

    for x := range room.Width {
      for y := range room.Height {
        if room.Grid[x][y].Type == "furniture" {
          if roomName == "" {
            switch room.Grid[x][y].ObstacleName {
            case "bed":
              if roomName == "" && roomNameToIndex["Jack's Room"] == 0 {
                roomName = "Jack's Room"
              }else if roomNameToIndex["Sarah's Room"] == 0 && roomNameToIndex["Johnny's Room"] == 0 {
                roomName = "Sarah's Room"
              }else {
                roomName = "Johnny's Room"
              }
            case "desk":
              if roomName == "" {
                roomName = "study"
              }
            case "sofa", "tv":
              roomName = "Living Room"
            case "stove", "fridge", "sink":
              roomName = "Kitchen"
            }
          }

          // Johnny's door
          if room.Grid[x][y].ObstacleName == "johnny's door" {
            robot.World.UpdateDoorsStatus("johnny's door", false)
          }
        }
      }
    }

    // If we can't determine the name of the room, then give a default name
    if roomName == "" {
      roomName = fmt.Sprintf("Room %d", i)
    }

    roomNameToIndex[roomName] = i
    fmt.Printf("Identified room %s (index %d)\n", roomName, i)
  }
  // Scan the house for objects to build our logical world
  fmt.Println("Robot is scanning for the object...")

  for _, room := range house.Rooms {
    for x := range room.Width {
      for y := range room.Height {
        if room.Grid[x][y].Type == "furniture" && room.Grid[x][y].ObstacleName != "" {
          robot.World.UpdateObjectFound(room.Grid[x][y].ObstacleName)
        }
      }
    }
  }
}