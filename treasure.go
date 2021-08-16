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
	playerPosition     = make(map[int]map[int]bool)
	treeWayPosition    = make(map[int]map[int]bool)
	alreadyUsedTreeWay = make(map[int]map[int]bool)
	obstaclePosition   = map[int]map[int]bool{
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
	fmt.Println(probablyTreasurePosition)
	setArenaWithAllProbablyTreasurePlace()
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

	if player found tree way position,
	player will priority move to the current move order
*/
func searchTreasure() {
	for {
		time.Sleep(1 * time.Second)
		if canMoveUp() {
			fmt.Println("player move up")
			movePlayerUp()
			if canMoveUp() && canMoveRight() {
				setTreeWayPosition()
			}
			printArena()
		} else if canMoveRight() {
			fmt.Println("player move right")
			movePlayerRight()
			if canMoveRight() && canMoveDown() {
				setTreeWayPosition()
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
				if reflect.DeepEqual(treeWayPosition, alreadyUsedTreeWay) {
					break
				}
				resetPlayerFromTreeWayPosition()
				time.Sleep(1 * time.Second)
				printArena()
			}
		}
	}
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
	fmt.Println("player probably found the treasure")
	probablyTreasurePosition = append(probablyTreasurePosition, []int{y, x})
}

func setTreeWayPosition() {
	y, x := getCurrentPlayerPosition()
	if _, ok := treeWayPosition[y]; !ok {
		treeWayPosition[y] = map[int]bool{
			x: false,
		}
	} else {
		treeWayPosition[y][x] = false
	}
}

func fromTreeWay() bool {
	y, x := getCurrentPlayerPosition()
	if _, ok := treeWayPosition[y]; !ok {
		return false
	}
	return treeWayPosition[y][x]
}

func resetPlayerFromTreeWayPosition() {
	var x, y int
	counter := 0
Loop:
	for yPosition, value := range treeWayPosition {
		for xPosition := range value {
			counter++
			if alreadyUseTreeWay(yPosition, xPosition) {
				continue
			}
			x = xPosition
			y = yPosition
			setUsedTreeWay(y, x)
			break Loop
		}
	}
	treeWayPosition[y][x] = true
	fmt.Println("reset from tree way position")
	currentYPosition, currentXPositon := getCurrentPlayerPosition()
	setPlayerPosition(y, x)
	arena[y][x] = "X"
	arena[currentYPosition][currentXPositon] = "."
}

func setUsedTreeWay(y, x int) {
	if _, ok := alreadyUsedTreeWay[y]; !ok {
		alreadyUsedTreeWay[y] = map[int]bool{
			x: true,
		}
		return
	}
	alreadyUsedTreeWay[y][x] = true
}

func alreadyUseTreeWay(y, x int) bool {
	if _, ok := alreadyUsedTreeWay[y]; !ok {
		return false
	}
	return alreadyUsedTreeWay[y][x]
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
