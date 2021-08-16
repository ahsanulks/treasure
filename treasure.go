package main

import (
	"fmt"
	"reflect"
	"time"
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

// type Move byte

// const (
// 	Up Move = iota
// 	Right
// 	Down
// )

var (
	// position always y, x
	playerPosition      = make(map[int]map[int]bool)
	threeWayPosition    = make(map[int]map[int]bool)
	alreadyUsedThreeWay = make(map[int]map[int]bool)
	obstaclePosition    = map[int]map[int]bool{
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
	movedPoint               = make(map[int]map[int]bool)
)

func main() {
	setStartPlayerPosition()
	makeArena()
	printArena()
	searchTreasure()
	setArenaWithAllProbablyTreasurePlace()
	fmt.Println()
	for index, value := range probablyTreasurePosition {
		fmt.Printf("%d place possible treasure is at point x: %d, y: %d\n", index+1, value[1], value[0])
	}
	fmt.Printf("\nthe map will be\n")
	printArena()
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
	setMovedWay(4, 1)
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

	if player found three way position,
	player will priority move to the current move order
	if player can't go down, player will reset to threeway position
	if all threway already used, the player will stop search the treasure
*/
func searchTreasure() {
	for {
		time.Sleep(1 * time.Second)
		if canMoveUp() {
			fmt.Println("player move up")
			movePlayerUp()
			if canMoveUp() && canMoveRight() {
				setThreeWayPosition()
			}
			printArena()
		} else if canMoveRight() {
			fmt.Println("player move right")
			movePlayerRight()
			if canMoveRight() && canMoveDown() {
				setThreeWayPosition()
			}
			printArena()
		} else if canMoveDown() {
			fmt.Println("player move down")
			movePlayerDown()
			setProbablyTreasurePosition()
			printArena()
			if canMoveDown() {
				continue
			} else {
				if reflect.DeepEqual(threeWayPosition, alreadyUsedThreeWay) {
					break
				}
				resetPlayerFromthreeWayPosition()
				time.Sleep(1 * time.Second)
				printArena()
			}
		}
	}
	time.Sleep(1 * time.Second)
}

func movePlayerUp() {
	y, x := getCurrentPlayerPosition()
	setPlayerPosition(y-1, x)
	setMovedWay(y-1, x)
	arena[y-1][x] = "X"
	arena[y][x] = "."
}

func movePlayerRight() {
	y, x := getCurrentPlayerPosition()
	setPlayerPosition(y, x+1)
	setMovedWay(y, x+1)
	arena[y][x+1] = "X"
	arena[y][x] = "."
}

func movePlayerDown() {
	y, x := getCurrentPlayerPosition()
	setPlayerPosition(y+1, x)
	setMovedWay(y+1, x)
	arena[y+1][x] = "X"
	arena[y][x] = "."
}

func canMoveUp() bool {
	y, x := getCurrentPlayerPosition()
	if alreadyMovedThisWay(y-1, x) {
		return false
	}
	return !isObstacle(y-1, x)
}

func canMoveRight() bool {
	y, x := getCurrentPlayerPosition()
	if alreadyMovedThisWay(y, x+1) {
		return false
	}
	return !isObstacle(y, x+1)
}

func canMoveDown() bool {
	y, x := getCurrentPlayerPosition()
	if alreadyMovedThisWay(y+1, x) {
		return false
	}
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
	fmt.Println("player probably found the treasure place")
	probablyTreasurePosition = append(probablyTreasurePosition, []int{y, x})
}

func setThreeWayPosition() {
	y, x := getCurrentPlayerPosition()
	if _, ok := threeWayPosition[y]; !ok {
		threeWayPosition[y] = map[int]bool{
			x: false,
		}
	} else {
		threeWayPosition[y][x] = false
	}
}

func resetPlayerFromthreeWayPosition() {
	var x, y int
	counter := 0
Loop:
	for yPosition, value := range threeWayPosition {
		for xPosition := range value {
			counter++
			if alreadyUseThreeWay(yPosition, xPosition) {
				continue
			}
			x = xPosition
			y = yPosition
			setUsedThreeWay(y, x)
			break Loop
		}
	}
	threeWayPosition[y][x] = true
	fmt.Println("reset from three way position")
	currentYPosition, currentXPositon := getCurrentPlayerPosition()
	setPlayerPosition(y, x)
	arena[y][x] = "X"
	arena[currentYPosition][currentXPositon] = "."
}

func setUsedThreeWay(y, x int) {
	if _, ok := alreadyUsedThreeWay[y]; !ok {
		alreadyUsedThreeWay[y] = map[int]bool{
			x: true,
		}
		return
	}
	alreadyUsedThreeWay[y][x] = true
}

func alreadyUseThreeWay(y, x int) bool {
	if _, ok := alreadyUsedThreeWay[y]; !ok {
		return false
	}
	return alreadyUsedThreeWay[y][x]
}

func setMovedWay(y, x int) {
	if _, ok := movedPoint[y]; !ok {
		movedPoint[y] = map[int]bool{
			x: true,
		}
		return
	}
	movedPoint[y][x] = true
}

func alreadyMovedThisWay(y, x int) bool {
	if _, ok := movedPoint[y]; !ok {
		return false
	}
	return movedPoint[y][x]
}

func setArenaWithAllProbablyTreasurePlace() {
	for _, treasurePoint := range probablyTreasurePosition {
		arena[treasurePoint[0]][treasurePoint[1]] = "$"
	}
}
