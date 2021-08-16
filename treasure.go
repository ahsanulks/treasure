package main

import (
	"fmt"
)

/*
vertical = x
horizontal = y

pattern
obstacle
x == 0-7 && y == 0 top
x == 0-7 && y == 5 bottom
x == 0 && y == 0-5 left
x == 7 && y == 0-5 right

initial player position
x == 1 && y == 4

########
#......#
#.###..#
#...#.##
#X#....#
########
*/

var (
	// position always y, x
	playerPosition   = make(map[int]map[int]bool)
	obstaclePosition = map[int]map[int]bool{
		2: {
			2: true,
			3: true,
			4: true,
		},
		3: {
			4: true,
			6: true,
		},
		4: {
			2: true,
		},
	}
	arena                    [][]string
	probablyTreasurePosition [][]int
)

func main() {
	setStartPlayerPosition()
	makeArena()
	printArena()
	searchTreasure()
	fmt.Println(probablyTreasurePosition)
}

func makeArena() {
	for y := 0; y <= 5; y++ {
		var columns []string
		for x := 0; x <= 7; x++ {
			if isObstacle(y, x) {
				columns = append(columns, "#")
			} else if playerPosition[y][x] {
				columns = append(columns, "X")
			} else {
				columns = append(columns, ".")
			}
		}
		arena = append(arena, columns)
	}
}

func printArena() {
	for _, rows := range arena {
		for index, columns := range rows {
			if index == len(rows)-1 {
				fmt.Println(columns)
			} else {
				fmt.Print(columns)
			}
		}
	}
}

func isObstacle(y, x int) bool {
	return y == 0 || y == 5 || x == 0 || x == 7 || obstaclePosition[y][x]
}

func setStartPlayerPosition() {
	setPlayerPosition(4, 1)
}

func setPlayerPosition(y, x int) {
	playerPosition = map[int]map[int]bool{
		y: {
			x: true,
		},
	}
}

/*
	player only move with this following order
	1. up
	2. right
	3. down
*/
func searchTreasure() {
	var (
		alreadyMoveUp    = false
		alreadyMoveRight = false
	)
	for i := 0; i < 4; i++ {
		if !alreadyMoveUp {
			fmt.Println("player move up")
			movePlayerUp()
			alreadyMoveUp = true
		} else if !alreadyMoveRight {
			fmt.Println("player move right")
			movePlayerRight()
			if canMoveDown() {
				alreadyMoveRight = true
			}
		} else {
			fmt.Println("player move down")
			movePlayerDown()
			setProbablyTreasurePosition()
		}
		printArena()
	}
}

func movePlayerUp() {
	y, x := getCurrentPlayerPosition()
	setPlayerPosition(y-1, x)
	arena[y-1][x] = "X"
	arena[y][x] = "."
}

func movePlayerRight() {
	y, x := getCurrentPlayerPosition()
	setPlayerPosition(y, x+1)
	arena[y][x+1] = "X"
	arena[y][x] = "."
}

func movePlayerDown() {
	y, x := getCurrentPlayerPosition()
	setPlayerPosition(y+1, x)
	arena[y+1][x] = "X"
	arena[y][x] = "."
}

func canMoveDown() bool {
	y, x := getCurrentPlayerPosition()
	return !isObstacle(y+1, x)
}

func getCurrentPlayerPosition() (y, x int) {
	var currentYpositon, currentXpositon int
	for y, xPosition := range playerPosition {
		currentYpositon = y
		for x := range xPosition {
			currentXpositon = x
		}
	}
	return currentYpositon, currentXpositon
}

func setProbablyTreasurePosition() {
	y, x := getCurrentPlayerPosition()
	probablyTreasurePosition = append(probablyTreasurePosition, []int{y, x})
}
