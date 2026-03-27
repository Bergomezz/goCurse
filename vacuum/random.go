package main

import "time"

func CleanRoomRandomWalk(room *Room, robot *Robot) {
	startTime := time.Now()
  moveCount := 0

  //Calculate cleaning time
  cleaningTime := time.Since(startTime)

  displaySummary(room, robot, moveCount, cleaningTime)
}