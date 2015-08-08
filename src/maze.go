package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var cols int64
var rows int64
var cells int64

var maze []int64
var solution []int64

var north = int64(1)
var east = int64(2)
var south = int64(4)
var west = int64(8)

var seed int64
var scale int64
var straight int64
var twisty int64

var header int64
var footer int64

var output = "output"

var solve bool
var asc = true

func main() {
	fmt.Printf("Maze: %d by %d\n", cols, rows)
	fmt.Printf("Seed: %d, straight: %d\n", seed, straight)
	create(cols, rows)
	ascii()
}

func init() {
	loadFlags()
	cells = cols * rows

	maze = make([]int64, cells)
	if seed == 1 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
}

func loadFlags() {
	flag.Int64Var(&cols, "cols", 8, "Number of columns in the maze")
	flag.Int64Var(&rows, "rows", 8, "Number of rows in the maze")
	flag.Int64Var(&straight, "straight", 0, "Integer >= 0. Higher numbers make straighter hallways")
	flag.Int64Var(&twisty, "twisty", 0, "Integer >= 0. Higher numbers make twistier hallways")
	flag.Int64Var(&seed, "seed", 1, "Integer value for the random seed")
	flag.BoolVar(&solve, "solve", false, "true to produce a graphic of the solution")

	flag.Parse()
}

func create(cols, rows int64) {
	var stack []int64

	start := int64(0)
	count := int64(1)
	current := start
	path := int64(0)

	stack = append(stack, start)

	prev := int64(0)

	n, e, w, s, g := float64(0.0), float64(0.0), float64(0.0), float64(0.0), float64(0.0)

	for count < cells {
		g = 0.0
		path = 0
		n = 0.0
		e = 0.0
		s = 0.0
		w = 0.0

		// if the passage north is not open yet
		// and it is at least on the second row
		// and the square above it has never been visited
		if (maze[current]&north) == 0 && current >= cols && maze[current-cols] == 0 {
			n = rand.Float64()
			if prev == north {
				if straight > 0 {
					n = (n + float64(straight)) / (1 + float64(straight))
				}
				if twisty > 0 {
					n = n / (1 + float64(twisty))
				}
			}
			g = n
			path = north
		}

		if (maze[current]&east) == 0 && current%cols != cols-1 && maze[current+1] == 0 {
			e = rand.Float64()
			if prev == east {
				e = (e + float64(straight)) / (1 + float64(straight))
				if twisty > 0 {
					e = e / (1 + float64(twisty))
				}
			}
			if e > g {
				g = e
				path = east
			}
		}

		if (maze[current]&west) == 0 && current%cols != 0 && maze[current-1] == 0 {
			w = rand.Float64()
			if prev == west {
				w = (w + float64(straight)) / (1 + float64(straight))
				if twisty > 0 {
					w = w / (1 + float64(twisty))
				}
			}
			if w > g {
				g = w
				path = west
			}
		}

		if (maze[current]&south) == 0 && current < cells-cols-1 && maze[current+cols] == 0 {
			s = rand.Float64()
			if prev == south {
				s = (s + float64(straight)) / (1 + float64(straight))
				if twisty > 0 {
					s = s / (1 + float64(twisty))
				}
			}
			if s > g {
				g = s
				path = south
			}
		}

		prev = path

		if path == 0 {
			current, stack = stack[len(stack)-2], stack[:len(stack)-1]
		} else {
			maze[current] = maze[current] | path

			if path == north {
				current -= cols
				maze[current] |= south
			}

			if path == east {
				current++
				maze[current] |= west
			}

			if path == south {
				current += cols
				maze[current] |= north
			}

			if path == west {
				current--
				maze[current] |= east
			}

			stack = append(stack, current)

			if current == cells-1 && len(solution) == 0 {
				solution = make([]int64, len(stack))
				copy(solution, stack)
			}
			count++
		}
	}
}

func ascii() {
	p("+")
	for i := int64(0); i < cols; i++ {
		p("---+")
	}
	p("\n")

	for i := int64(0); i < rows; i++ {
		p("|")
		for j := int64(0); j < cols; j++ {
			if i == 0 && j == 0 {
				p(" 0 ")
			} else if i == rows-1 && j == cols-1 {
				p(" X ")
			} else if contains(solution, j+cols*i) && solve {
				p(" o ")
			} else {
				p("   ")
			}

			if maze[i*cols+j]&east == east && j != cols-1 {
				p(" ")
			} else {
				p("|")
			}
		}

		p("\n")
		p("+")

		for j := int64(0); j < cols; j++ {
			if maze[i*cols+j]&south == south && i != rows-1 {
				p("   +")
			} else {
				p("---+")
			}
		}

		p("\n")
	}
}

func dd() {
	for i, v := range maze {
		fmt.Println(i, v)
	}
}

func p(a string) {
	fmt.Print(a)
}

func contains(path []int64, cell int64) bool {
	for _, v := range path {
		if v == cell {
			return true
		}
	}
	return false
}
