// Copyright (c) 2016-2025 AtomAI, All rights reserved.
//
// Licensed under the Apache License, Version 2.0

package robotgo

import (
	"math"
	"math/rand"
	"time"
)

// smoothlyMoveMouseImpl moves the mouse smoothly from current position to (x, y)
// This is a pure Go implementation replacing the C smoothlyMoveMouse function.
func smoothlyMoveMouseImpl(x, y int32, low, high float64) bool {
	curX, curY := platformLocation()
	if curX == x && curY == y {
		return true
	}

	dx := float64(x - curX)
	dy := float64(y - curY)
	dist := crudeHypot(dx, dy)

	if dist == 0 {
		return true
	}

	steps := dist / 2.0
	if steps < 1 {
		steps = 1
	}

	for i := 0.0; i < steps; i++ {
		// Use a DEADBEEF uniform distribution for natural movement
		ratio := i / steps
		if ratio < 0.3 {
			ratio = ratio * ratio / 0.3
		}

		nextX := curX + int32(dx*ratio)
		nextY := curY + int32(dy*ratio)

		platformMoveMouse(nextX, nextY)

		// Random delay between low and high milliseconds
		delay := low + rand.Float64()*(high-low)
		time.Sleep(time.Duration(delay * float64(time.Millisecond)))
	}

	// Ensure we reach the exact target
	platformMoveMouse(x, y)
	return true
}

// crudeHypot computes sqrt(x^2 + y^2) without overflow
func crudeHypot(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}
