package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	//"image"
)

type Cell struct {
	X int
	Y int
}

//saw this somewhere, is it standard practice?
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func initmaze(x, y int) [][]string {
	//makes 2d slice and intializes walls/spaces
	maze := make([][]string, x)

	for i := range maze {
		maze[i] = make([]string, y)
		for j := range maze[i] {
			//every [odd][odd] cell is a node
			//all other cells are vertices
			if i*j%2 == 1 {
				maze[i][j] = " "
			} else {
				maze[i][j] = "*"
			}
		}
	}

	return maze
}

func printmaze(maze [][]string) {
	//output pretty maze with spacing
	//python style loop instead of c style
	for _, row := range maze {
		for _, ele := range row {
			fmt.Print(ele, " ")
		}
		fmt.Print("\n")
	}
}

func findneighbor(minihead Cell, minimaze [][]string) int {
	//helper function to drunken walk to an unvisted space
	//pre-checking instead of walking and retracing steps
	// 0,1,2,3 = up,down,left,right

	unvisited := make([]int, 0, 4)
	minix, miniy := len(minimaze), len(minimaze[0])

	//check up
	if minihead.Y-1 >= 0 && minimaze[minihead.X][minihead.Y-1] == "" {
		//unvisited.append((i)*minix + j - 1)
		unvisited = append(unvisited, 0)
	}
	//down
	if minihead.Y+1 < miniy && minimaze[minihead.X][minihead.Y+1] == "" {
		//unvisited.append((i)*minix + j + 1)
		unvisited = append(unvisited, 1)
	}
	//left
	if minihead.X-1 >= 0 && minimaze[minihead.X-1][minihead.Y] == "" {
		//unvisited.append((i-1)*minix + j)
		unvisited = append(unvisited, 2)
	}
	//right
	if minihead.X+1 < minix && minimaze[minihead.X+1][minihead.Y] == "" {
		//unvisited.append((i+1)*minix + j)
		unvisited = append(unvisited, 3)
	}

	if len(unvisited) == 0 {
		return -1
	}

	//return a random direction
	return unvisited[rand.Intn(len(unvisited))]
}

func genmaze(maze [][]string) [][]string {
	//"growing tree" maze generation algorithm

	x, y := len(maze), len(maze[0])
	minix, miniy := len(maze)/2, len(maze[0])/2

	//generate smaller maze to represent nodes
	//default "" cells in minimaze have not been visited
	minimaze := make([][]string, minix)
	for i := range minimaze {
		minimaze[i] = make([]string, miniy)
	}

	//slice for visited cells (who potentially have neighbors)
	cells := make([]Cell, 0, minix*miniy)

	//start at maze end node - (x-2, y-2) and (minix-1, miniy-1)
	//mark visited in minimaze with "v"
	head := Cell{x - 2, y - 2}
	minihead := Cell{minix - 1, miniy - 1}
	minimaze[minihead.X][minihead.Y] = "v"

	//add minihead to visited cells
	cells = append(cells, minihead)

	//iterative part of algorithm
	//continue until there are no unvisted neighbors
	for len(cells) > 0 {
		//look for unvisted neighbors from head
		next := findneighbor(minihead, minimaze)

		if next < 0 {
			//no unvisted neighbors
			//remove cell (always last cell in slice) and backtrack
			cells = cells[:len(cells)-1]
			if len(cells) == 0 {
				break
			}

			minihead = cells[len(cells)-1]
			head.X, head.Y = minihead.X*2+1, minihead.Y*2+1

			//note: can change next minihead to random from cells
			//will give different results
			//have to change remove function to be smarter
		} else if next == 0 {
			//go up
			//remove walls and move head
			maze[head.X][head.Y-1] = " "
			head.Y = head.Y - 2

			//move minihead
			minihead.Y = minihead.Y - 1
			minimaze[minihead.X][minihead.Y] = "v"
			cells = append(cells, minihead)

		} else if next == 1 {
			//go down
			//remove walls and move head
			maze[head.X][head.Y+1] = " "
			head.Y = head.Y + 2

			//move minihead
			minihead.Y = minihead.Y + 1
			minimaze[minihead.X][minihead.Y] = "v"
			cells = append(cells, minihead)

		} else if next == 2 {
			//go left
			//remove walls and move head
			maze[head.X-1][head.Y] = " "
			head.X = head.X - 2

			//move minihead
			minihead.X = minihead.X - 1
			minimaze[minihead.X][minihead.Y] = "v"
			cells = append(cells, minihead)

		} else if next == 3 {
			//go right
			//remove walls and move head
			maze[head.X+1][head.Y] = " "
			head.X = head.X + 2

			//move minihead
			minihead.X = minihead.X + 1
			minimaze[minihead.X][minihead.Y] = "v"
			cells = append(cells, minihead)
		}
	}
	//loop ends when there are no unvisited cells
	//meaning all nodes are included in maze
	return maze
}

func writemaze(maze [][]string) {
	//write maze to file XbyYmaze.txt
	//easier to use buffer and WriteString or []byte and WriteFile?

	filename := strconv.Itoa(len(maze)) + "by" + strconv.Itoa(len(maze[0])) + "maze.txt"

	output := ""

	for i := range maze {
		for j := range maze[i] {
			output += maze[i][j]
		}
		output += "\n"
	}

	f, err := os.Create(filename)
	check(err)
	defer f.Close()

	buf := bufio.NewWriter(f)
	_, err = buf.WriteString(output)
	check(err)
	defer buf.Flush()
}

//TODO: use image libraries to draw maze for practice
//possibly integrate with genmaze to show realtime generation
func drawmaze(maze [][]string) {
	return
}

func main() {
	//default to 41x41 maze
	x, y := 41, 41

	//grab command line arguments, if any
	//force odd lengths to account for nodes and edges
	if len(os.Args) == 2 {
		//grabbing numbers from command line is ugly
		if i, err := strconv.Atoi(os.Args[1]); err == nil {
			if i%2 == 0 {
				i += 1
			}
			x, y = i, i
		} else {
			panic(err)
		}

	} else if len(os.Args) > 2 {
		if i, err := strconv.Atoi(os.Args[1]); err == nil {
			if i%2 == 0 {
				i += 1
			}
			x = i
		} else {
			panic(err)
		}
		if j, err := strconv.Atoi(os.Args[2]); err == nil {
			if j%2 == 0 {
				j += 1
			}
			y = j
		} else {
			panic(err)
		}
	}

	//seed rand
	rand.Seed(time.Now().UnixNano())

	//generate, print, and write maze
	maze := genmaze(initmaze(x, y))
	printmaze(maze)
	writemaze(maze)
}
