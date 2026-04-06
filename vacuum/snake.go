package main

import "time"

func CleanRoomSnake(room *Room, robot *Robot) {
	// Initialize start time and moveCount
	startTime := time.Now()
  moveCount := 0

  // Generate snaking pattern
  coveragePoints := generateSnakingPattern(room)

  // Clean cell
  Clean(robot, room)

  if room.Animate {
    room.Display(robot, false)
    time.Sleep(moveDelay)
  }

  // visiting each point in the coverage pattern (for)
  for _, point := range coveragePoints {

    // Skip cell if already clean
    if room.Grid[point.X][point.Y].Cleaned {
      continue
    }

    // Find path to the next point
    path := Astar(room, robot.Position, point)

    // If not path found, try next point
    if len(path) == 0 {
      continue
    }

    // Move along the path
    for i := 1; i < len(path); i++ {

      // Update robot position
      robot.Position = path[i]
      robot.Path = append(robot.Path, path[i])

      // Clean
      Clean(robot, room)

      // Display the room
      if room.Animate{
        room.Display(robot, false)
        time.Sleep(moveDelay)
      }

      // Increment the moveCount
      moveCount++

    }

  }

  // Do a final sweep
  finalCleanup(room, robot, &moveCount)

  // Calculate cleaning time
  cleaningTime := time.Since(startTime)

  // Display final statistics
  displaySummary(room, robot, moveCount, cleaningTime)
}

func generateSnakingPattern(room *Room) []Point {
  var points []Point
  var directionX = 1

  for y := 1; y < room.Height-1; y++ {
    if directionX == 1 {
      // Moving left to right
      for x := 1; x < room.Width-1; x++ {
        if !room.Grid[x][y].Obstacle {
          points= append(points, Point{X: x,Y: y})
        }
      }
    }else {
      // Moving right to left
      for x := room.Width-2; x >= 1; x-- {
        if !room.Grid[x][y].Obstacle {
          points = append(points, Point{X: x, Y: y})
        }
      }
    }
    directionX *= -1
  }
  return points
}