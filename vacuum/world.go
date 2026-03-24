package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	// Display characters
	charRobot          = "🔴"
	charWall           = "🟦"
	charFurniture      = "🪑"
	charClean          = "🧼"
	charDirty          = "🟫"
	charPath           = "🟢"
	charCat            = "🐱" // Display character for cat
	catStopProbability = 0.1 // Probability of cat stopping
	catStopDuration    = 5   // Duration cat stays still (in animation frames)
	moveDelay          = 50 * time.Millisecond
	cellSize           = 10
)

type Point struct {
	X, Y int
}

type Cell struct {
	Type         string // wall, furniture, clean, dirty, bike
	Cleaned      bool
	Obstacle     bool
	ObstacleName string
}

type Furniture struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

type Room struct {
	Grid               [][]Cell
	Width              int
	Height             int
	CleanableCellCount int
	CleanedCellCount   int
	Animate            bool
}

type RoomConfig struct {
	Width     int         `json:"width"`
	Height    int         `json:"height"`
	Furniture []Furniture `json:"furniture"`
}

func NewRoom(configFile string, animate bool) *Room {
  // Load from JSON config
  roomConfig, err := LoadRoomConfig(configFile)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  // Convert dimension to grid cells
  gridWidth := roomConfig.Width / cellSize
  gridHeight := roomConfig.Height / cellSize

  // Create a grid
  grid :=make([][]Cell, gridWidth)
  for i := range grid {
    grid[i] = make([]Cell, gridHeight)
    for j := range grid[i] {
      grid[i][j] = Cell{Type: "dirty", Cleaned: false, Obstacle: false}
    }
  }
  // Add walls
  for i := 0; i < gridWidth; i++ {
    grid[i][0] = Cell{Type: "wall", Cleaned: false, Obstacle: true, ObstacleName: "wall"}
    grid[i][gridHeight-1] = Cell{Type: "wall", Cleaned: false, Obstacle: true, ObstacleName: "wall"}
  }

  for j := 0; j < gridHeight; j++ {
    grid[0][j] = Cell{Type: "wall", Cleaned: false, Obstacle: true, ObstacleName: "wall"}
    grid[gridWidth-1][j] = Cell{Type: "wall", Cleaned: false, Obstacle: true, ObstacleName: "wall"}
  }

  // TODO Add furniture

  // Count cleanable cells
  cleanableCellCount := 0
  for i := 0; i < gridWidth; i++ {
    for j := 0; j < gridHeight; j++ {
      if !grid[i][j].Obstacle {
        cleanableCellCount++
      }
    }
  }

	return &Room{
    Grid: grid,
    Width: gridWidth,
    Height: gridHeight,
    CleanableCellCount: cleanableCellCount,
    CleanedCellCount: 0,
    Animate: animate,
  }
}

func LoadRoomConfig(filename string) (*RoomConfig, error){
  // Read JSON file
  jsonData, err := os.ReadFile(filename)
  if err != nil {
    return nil, fmt.Errorf("Error reding JSON file: %v", err)
  }

  // Parse JSON
  var config RoomConfig
  if err := json.Unmarshal(jsonData,&config); err != nil {
    return nil, fmt.Errorf("Error parsing JSON: %v", err)
  }

  return &config, nil
}