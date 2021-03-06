package main

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const polySize = 50
const size = 10

func convertCords(x int64, y int64) (float64, float64) {
	return float64((x - 1) * polySize), float64((size - y) * polySize)
}

func drawSquare(imd *imdraw.IMDraw, x int64, y int64, color pixel.RGBA) {
	imd.Color = color
	realX, realY := convertCords(x, y)
	imd.Push(
		pixel.V(realX, realY),
		pixel.V(realX+polySize, realY),
		pixel.V(realX+polySize, realY+polySize),
		pixel.V(realX, realY+polySize),
	)
	imd.Polygon(0)
}

func drawGrid(imd *imdraw.IMDraw) {
	top := float64(size * polySize)
	right := float64(size * polySize)
	left := float64(0)
	bottom := float64(0)
	for x := left; x < right; x += polySize {
		imd.Color = pixel.RGB(0, 0, 0)
		imd.Push(
			pixel.V(x, top),
			pixel.V(x, bottom),
		)
		imd.Line(1)
	}
	for y := bottom; y < top; y += polySize {
		imd.Color = pixel.RGB(0, 0, 0)
		imd.Push(
			pixel.V(left, y),
			pixel.V(right, y),
		)
		imd.Line(1)
	}
}

type Point struct {
	x int64
	y int64
}

type Snake struct {
	nextPoint     Point
	growth        int64
	direction     byte
	nextDirection byte
	head          Point
	tail          []Point
}

func (s *Snake) draw(imd *imdraw.IMDraw) {
	drawSquare(imd, s.nextPoint.x, s.nextPoint.y, pixel.RGB(0, 1, 0))
	for _, point := range s.tail {
		drawSquare(imd, point.x, point.y, pixel.RGB(0, 0, 1))
	}
	drawSquare(imd, s.head.x, s.head.y, pixel.RGB(1, 0, 0))
}

func (s *Snake) move() {

	delta := 1
	if s.growth > 0 {
		delta = 0
		s.growth--
	}

	s.tail = append([]Point{s.head}, s.tail[:len(s.tail)-delta]...)
	s.direction = s.nextDirection

	var dX, dY int64 = 0, 0
	switch s.direction {
	case 'u':
		dY = -1
		break
	case 'd':
		dY = +1
		break
	case 'l':
		dX = -1
		break
	case 'r':
		dX = 1
		break
	}
	s.head.y += dY
	s.head.x += dX
	if s.head.y > size {
		s.head.y = 1
	}
	if s.head.y < 1 {
		s.head.y = size
	}
	if s.head.x > size {
		s.head.x = 1
	}
	if s.head.x < 1 {
		s.head.x = size
	}
}

func (s *Snake) generateNextPoint() {
	s.nextPoint = Point{
		rand.Int63n(size-1) + 1,
		rand.Int63n(size-1) + 1,
	}
	if s.nextPoint == s.head {
		s.generateNextPoint()
	}
	for _, point := range s.tail {
		if point == s.nextPoint {
			s.generateNextPoint()
		}
	}
}

func (s *Snake) checkPoint() {
	if s.head == s.nextPoint {
		s.growth++
		s.generateNextPoint()
	}
}

func (s *Snake) checkCollisions() {
	for i, point := range s.tail {
		if point == s.head {

			s.tail = s.tail[0:i]
			break
		}
	}
}

var snake Snake = Snake{
	head:   Point{3, 3},
	growth: 0,
	tail: []Point{
		{3, 4},
		{3, 5},
		{3, 6},
		{3, 7},
		{3, 8},
	},
	direction:     'u',
	nextDirection: 'u',
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Snake",
		Bounds: pixel.R(0, 0, polySize*size, polySize*size),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	last := time.Now()
	frequency := float64(0.25)
	snake.generateNextPoint()
	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		imd := imdraw.New(nil)

		if win.Pressed(pixelgl.KeyLeft) && snake.direction != 'r' {
			snake.nextDirection = 'l'
		}
		if win.Pressed(pixelgl.KeyRight) && snake.direction != 'l' {
			snake.nextDirection = 'r'
		}
		if win.Pressed(pixelgl.KeyUp) && snake.direction != 'd' {
			snake.nextDirection = 'u'
		}
		if win.Pressed(pixelgl.KeyDown) && snake.direction != 'u' {
			snake.nextDirection = 'd'
		}

		dt := time.Since(last).Seconds()
		if dt > frequency {
			last = time.Now()
			snake.move()
			snake.checkPoint()
			snake.checkCollisions()
		}
		draw(imd)

		imd.Draw(win)
		win.Update()
	}
}
func draw(imd *imdraw.IMDraw) {
	snake.draw(imd)
	drawGrid(imd)
}

type food struct {
	emoji        rune
	points, x, y int
}

func newFood(x, y int) *food {
	return &food{
		points: 10,
		emoji:  getFoodEmoji(),
		x:      x,
		y:      y,
	}
}

func getFoodEmoji() rune {
	if hasUnicodeSupport() {
		return randomFoodEmoji()
	}

	return '@'
}

func randomFoodEmoji() rune {
	f := []rune{
		'🍒',
		'🍍',
		'🍑',
		'🍇',
		'🍏',
		'🍌',
		'🍫',
		'🍭',
		'🍕',
		'🍩',
		'🍗',
		'🍖',
		'🍬',
		'🍤',
		'🍪',
	}

	return f[rand.Intn(len(f))]
}

func hasUnicodeSupport() bool {
	return strings.Contains(os.Getenv("LANG"), "UTF-8")
}

func main() {
	pixelgl.Run(run)
	randomFoodEmoji()
}
