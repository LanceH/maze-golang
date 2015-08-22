package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

var cols int64
var rows int64
var cells int64

var start int64
var finish int64

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
var unicursal bool
var ascii bool

func main() {
	fmt.Printf("Maze: %d by %d\n", cols, rows)
	fmt.Printf("Seed: %d, straight: %d\n", seed, straight)
	create(cols, rows)
	if unicursal {
		solve = false
		toUnicursal()
	}
	if ascii {
		toAscii()
	} else {
		toPng()
	}
}

func init() {
	loadFlags()
	cells = cols * rows
	if finish == start {
		finish = cells - 1
	}

	maze = make([]int64, cells)
	if seed == -1 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
}

func loadFlags() {
	flag.Int64Var(&cols, "cols", 8, "Number of columns in the maze")
	flag.Int64Var(&rows, "rows", 8, "Number of rows in the maze")

	flag.Int64Var(&start, "start", 0, "Number of the cell to start in (zero based)")
	flag.Int64Var(&finish, "finish", 0, "Number of the cell to finish (max = rows * cols - 1)")

	flag.Int64Var(&straight, "straight", 0, "Integer >= 0. Higher numbers make straighter hallways")
	flag.Int64Var(&twisty, "twisty", 0, "Integer >= 0. Higher numbers make twistier hallways")
	flag.Int64Var(&seed, "seed", -1, "Integer value for the random seed")
	flag.BoolVar(&ascii, "ascii", true, "true produces an ascii art version of the maze")
	flag.BoolVar(&solve, "solve", false, "true to produce a graphic of the solution")
	flag.BoolVar(&unicursal, "unicursal", false, "Convert the maze to a labyrinth (only one path)")

	flag.Parse()
}

func create(cols, rows int64) {
	var stack []int64

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

			if current == finish && len(solution) == 0 {
				solution = make([]int64, len(stack))
				copy(solution, stack)
			}
			count++
		}
	}
}

func toUnicursal() {
	pMaze := make([]int64, cells*4)
	for i := int64(0); i < rows; i++ {
		for j := int64(0); j < cols; j++ {
			current := i*cols + j
			ul := 4*i*cols + 2*j
			ur := 4*i*cols + 2*j + 1
			ll := 4*i*cols + 2*cols + 2*j
			lr := 4*i*cols + 2*cols + 2*j + 1
			if 0 != maze[current]&north {
				pMaze[ul] |= north
				pMaze[ul-2*cols] |= south
				pMaze[ur] |= north
				pMaze[ur-2*cols] |= south
			} else {
				pMaze[ul] |= east
				pMaze[ur] |= west
			}
			if 0 != maze[current]&east {
				pMaze[ur] |= east
				pMaze[ur+1] |= west
				pMaze[lr] |= east
				pMaze[lr+1] |= west
			} else {
				pMaze[ur] |= south
				pMaze[lr] |= north
			}
			if 0 != maze[current]&south {
				pMaze[ll] |= south
				pMaze[lr] |= south
				pMaze[ll+2*cols] |= north
				pMaze[lr+2*cols] |= north
			} else {
				pMaze[ll] |= east
				pMaze[lr] |= west
			}
			if 0 != maze[current]&west {
				pMaze[ul] |= west
				pMaze[ll] |= west
				pMaze[ul-1] |= east
				pMaze[ll-1] |= east
			} else {
				pMaze[ul] |= south
				pMaze[ll] |= north
			}
		}
	}
	maze = pMaze
	rows *= 2
	cols *= 2
}

func toPng() {
	f, err := os.OpenFile("output/maze.png", os.O_CREATE|os.O_WRONLY, 0666)
	fmt.Println("Building image")
	if err != nil {
		fmt.Println(err)
	}

	height := int(rows*4 + 1)
	width := int(cols*4 + 1)
	m := image.NewRGBA(image.Rect(0, 0, height, width))

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			m.SetRGBA(i, j, color.RGBA{255, 255, 255, 255})
		}
	}

	for i := int64(0); i < cols; i++ {
		for j := int64(0); j < rows; j++ {
			if i%cols == 0 {
				drawVert(m, 3, 3, 12, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	if err = png.Encode(f, m); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func toAscii() {
	p("+")
	for i := int64(0); i < cols; i++ {
		p("---+")
	}
	p("\n")

	var dots = make(map[int64]bool)
	for _, v := range solution {
		dots[v] = true
	}

	for i := int64(0); i < rows; i++ {
		p("|")
		for j := int64(0); j < cols; j++ {
			if !unicursal && start == i*cols+j {
				p(" 0 ")
			} else if !unicursal && finish == i*cols+j {
				p(" X ")
			} else if dots[j+cols*i] && solve {
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

func drawVert(m *image.RGBA, x, y, d int, col color.RGBA) {
	for i := 0; i < d; i++ {
		m.SetRGBA(x, y+i, col)
	}
}

func drawHoriz(m *image.RGBA, x, y, d int, col color.RGBA) {
	for i := 0; i < d; i++ {
		m.SetRGBA(x+i, y, col)
	}
}
