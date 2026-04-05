package main

import "time"

func CleanSpiralPattern(room *Room, robot *Robot) {
	startTime := time.Now()
  moveCount := 0

  // Find the center of the room
  centerX := room.Width / 2
  centerY := room.Height / 2

  // Find a valid point near the center
  centerPoint := findNearestCleanablePoint(room, Point{X: centerX, Y: centerY})

  // Find path to the center (A*)
  pathToCenter := Astar(room, robot.Position, centerPoint)

  // Move to the center point
  if len(pathToCenter) > 1 {
    for i := 1; i < len(pathToCenter); i++ {
      robot.Position = pathToCenter[i]
      robot.Path = append(robot.Path, robot.Position)
      Clean(robot, room)

      if room.Animate {
        room.Display(robot, false)
        time.Sleep(moveDelay)
      }

      moveCount++
    }
  }

  // Create a spiral pattern
  spiralPoints := generateSpiralPattern(room, centerPoint)

  // Follow the spiral pattern (for loop)
  for _, point := range spiralPoints {
    // Skip if cells already clean or an obstacle
    if room.Grid[point.X][point.Y].Cleaned || room.Grid[point.X][point.Y].Obstacle {
      continue
    }

    // Find path to the next point (A*)
    path := Astar(room, robot.Position, point)

    if len(path) <= 1 {
      continue
    }

    // Move along path
    for i := 1; i < len(path); i++ {
      robot.Position = path[i]
      robot.Path = append(robot.Path, robot.Position)
      Clean(robot, room)

      if room.Animate {
        room.Display(robot, false)
        time.Sleep(moveDelay)
      }
      moveCount++
    }
  }

  // Final cleanup
  finalCleanup(room, robot, &moveCount)

  // Calculate cleaning time
  cleaningTime := time.Since(startTime)

  // Display final statistics
  displaySummary(room, robot, moveCount, cleaningTime)
}

func finalCleanup(room *Room, robot *Robot, moveCount *int) {
  for i := 1; i < room.Width-1; i++ {
    for j := 1; j < room.Height-1; j++ {
      if !room.Grid[i][j].Obstacle && !room.Grid[i][j].Cleaned {
        // Find path to cell
        path := Astar(room, robot.Position, Point{X: i, Y: j})

        if len(path) <= 1 {
          continue
        }

        for k := 1; k < len(path); k++ {
          robot.Position = path[k]
          robot.Path = append(robot.Path, robot.Position)
          Clean(robot, room)
          if room.Animate {
            room.Display(robot, false)
            time.Sleep(moveDelay)
          }
          *moveCount++
        }
      }
    }
  }
}

func generateSpiralPattern(room *Room, center Point) []Point {
  var points []Point

  // Maximum possible spiral size
  maxSize := max(room.Width, room.Height)

  // Set delta X and delta Y
  dx := []int{1,0,-1,0}
  dy := []int{0,1,0,-1}

  // Start at center
  x, y := center.X, center.Y
  dir := 0 // starting move right

  // Set spiral parameters
  step := 1
  stepCount := 0
  dirChanges := 0

  // Generate spiral pattern
  for range maxSize*maxSize {
    // Add current point if valid
    if room.IsValid(x,y){
      points = append(points, Point{X: x, Y: y})
    }

    // Take a step
    x += dx[dir]
    y += dy[dir]
    stepCount++

    // Check to see if we need to change direction
    if stepCount == step {
      dir = (dir + 1) % 4
      stepCount = 0
      dirChanges++

      // Increase step size after two direction changes
      if dirChanges == 2 {
        step++
        dirChanges = 0
      }
    }

    // Break if we are out of bounds
    if x < 0 || x >= room.Width || y < 0 || y >= room.Height {
      break
    }
  }

  return points
}

func findNearestCleanablePoint(room *Room, target Point) Point {
  if room.IsValid(target.X, target.Y) && room.Grid[target.X][target.Y].Obstacle {
    return target
  }

  // Search for a valid point in expanding circle
  for radius := 1; radius < room.Width || radius < room.Height; radius++ {
    // Check all points at the current radius
    for dx := -radius; dx <= radius; dx++ {
      for dy := -radius; dy <= radius; dy++ {
        if abs(dx) != radius && abs(dy) != radius {
          continue
        }

        x, y := target.X + dx, target.Y + dy

        // Check to see if this point is valid and not an obstacle
        if room.IsValid(x,y) && !room.Grid[x][y].Obstacle {
          return Point{X: x, Y: y}
        }
      }
    }
  }

  // If no valid point, return the start point
  return Point{X: 1, Y: 1}
}