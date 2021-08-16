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

	arena [][]string
)

func main() {
	setPlayerPosition(4, 1)
	makeArena()
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

func setPlayerPosition(y, x int) {
	playerPosition = map[int]map[int]bool{
		y: {
			x: true,
		},
	}
}
