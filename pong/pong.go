package main
import (
	"experiments/noise"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"time"
)

type gameState int

const (
	start gameState = iota
	play
)

var gState = start

const (
	windowWidth = 800
	windowHeight = 600
)

var nums = [][]byte {
	{
		1,1,1,
		1,0,1,
		1,0,1,
		1,0,1,
		1,1,1,
	}, {
		1,1,0,
		0,1,0,
		0,1,0,
		0,1,0,
		1,1,1,
	}, {
		1,1,1,
		0,0,1,
		1,1,1,
		1,0,0,
		1,1,1,
	}, {
		1,1,1,
		0,0,1,
		0,1,1,
		0,0,1,
		1,1,1,
	},
}

type color struct {
	r, g, b byte
}

type position struct {
	x, y float32
}

type ball struct {
	position
	radius float32
	xv float32
	yv float32
	color color
}

// ---

func flerp(b1, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct * (float32(b2) - float32(b1)))
}

func colorLerp(c1, c2 color, pct float32) color {
	return color{flerp(c1.r, c2.r, pct), flerp(c1.g, c2.g, pct), flerp(c1.b, c2.b, pct)}
}

func getGradient(c1, c2 color) []color {
	result := make([]color, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}
	return result
}

func getDualGradient(c1, c2 , c3, c4 color) []color {
	result := make([]color, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		if pct < .5 {
			result[i] = colorLerp(c1, c2, pct * 2)
		} else {
			result[i] = colorLerp(c3, c4, pct * 1.5 - .5)
		}
	}
	return result
}

func clamp(min, max, v int) int {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

func rescaleAndDraw(noise []float32, min, max float32, gradient []color, w, h int) []byte {
	pixels := make([]byte, w * h * 4)
	scale := 255.0 / (max - min)
	offset := min * scale
	for i := range noise {
		noise[i] = noise[i]*scale - offset
		//b := byte(noise[i])
		c := gradient[clamp(0,255, int(noise[i]))]
		p := i * 4
		pixels[p] = c.r
		pixels[p + 1] = c.g
		pixels[p + 2] = c.b
	}
	return pixels
}

// ---


func drawNumber(position position, color color, size int, num int, pixels []byte)  {
	startX := int(position.x) - (size*3)/2
	startY := int(position.y) - (size*5)/2

	for key, value := range nums[num] {
		if value == 1 {
			for y := startY; y < startY+ size; y++ {
				for x := startX; x < startX+ size; x++ {
					setPixel(x,y,color,pixels)
				}
			}
		}
		startX += size
		if (key + 1) % 3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

func (ball *ball) draw(pixels []byte)  {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x * x + y * y < ball.radius * ball.radius {
				setPixel(int(ball.x + x), int(ball.y + y), ball.color, pixels)
			}
		}
	}
}

func getCenter() position {
	return position{windowWidth / 2, windowHeight / 2}
}

func (ball *ball) update(leftPaddle, rightPaddle *paddle, elapsedTime float32)  {
	ball.x += ball.xv * elapsedTime
	ball.y += ball.yv * elapsedTime

	if ball.y - ball.radius < 0 || ball.y + ball.radius > windowHeight {
		ball.yv = -ball.yv
	}

	if ball.x < 0 {
		rightPaddle.score++
		ball.position = getCenter()
		gState = start
	} else if int(ball.x) > windowWidth {
		leftPaddle.score++
		ball.position = getCenter()
		gState = start
	}

	if ball.x - ball.radius < leftPaddle.x + leftPaddle.w/2 {
		if ball.y > leftPaddle.y - leftPaddle.h/2 && ball.y <  leftPaddle.y + leftPaddle.h/2 {
			ball.xv = -ball.xv
			ball.x = leftPaddle.x + leftPaddle.w/2.0 + ball.radius
		}
	}

	if ball.x + ball.radius > rightPaddle.x - rightPaddle.w/2 {
		if ball.y > rightPaddle.y - rightPaddle.h/2 && ball.y <  rightPaddle.y + rightPaddle.h/2 {
			ball.xv = -ball.xv
			ball.x = rightPaddle.x - rightPaddle.w/2.0 - ball.radius
		}
	}
}

type paddle struct {
	position
	w	float32
	h	float32
	speed float32
	score int
	color color
}

func lerp(a, b, precept float32) float32 {
	return a + precept * (b - a)
}

func (paddle *paddle) draw(pixels []byte)  {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	for y := 0 ; y < int(paddle.h); y++ {
		for x := 0 ; x < int(paddle.w) ; x++ {
			setPixel(startX + x, startY + y, paddle.color, pixels)
		}
	}
	numX := lerp(paddle.x, getCenter().x,.2)
	drawNumber(position{numX,35}, paddle.color, 10, paddle.score, pixels)
}

func (paddle *paddle) update(keyState []uint8, axis int16,  elapsedTime float32)  {
	if keyState[sdl.SCANCODE_UP] != 0 {
		if paddle.y - paddle.h / 2 >= 0 {
			paddle.y -= paddle.speed * elapsedTime
		}
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		if paddle.y + paddle.h / 2 <= windowHeight {
			paddle.y += paddle.speed * elapsedTime
		}
	}

	if math.Abs(float64(axis)) > 1500 {
		percent := float32(axis) / 32767.0
		paddle.y += paddle.speed * percent * elapsedTime
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32)  {
	paddle.y = ball.y
}

func clear(pixels []byte, )  {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y* windowWidth + x ) * 4
	if index <= len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func main()  {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	windows, err := sdl.CreateWindow("Testing", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer  windows.Destroy()

	renderer, err := sdl.CreateRenderer(windows, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, windowWidth, windowHeight)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	var controllerHandlers []*sdl.GameController
	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
	}
	defer func() {
		for _, controller := range controllerHandlers {
			controller.Close()
		}
	}()
	pixels := make([]byte, windowWidth*windowHeight*4)

	player1 := paddle{
		position: position{50, 100},
		w:        20,
		h:        100,
		speed:	  300,
		color:    color{255,255,255},
	}

	player2 := paddle{
		position: position{windowWidth - 50, 100},
		w:        20,
		h:  	  100,
		speed:    300,
		color:    color{255,255,255},
	}

	ball := ball{
		position: getCenter(),
		radius:   20,
		xv:       400,
		yv:       400,
		color:    color{255,255,255},
	}

	keyState := sdl.GetKeyboardState()

	n, min, max := noise.MakeNoise(noise.FBM, .01, .5,2,3, windowWidth, windowHeight)
	gradient := getGradient(color{255,0,0},color{255,242,0})
	noisePixels := rescaleAndDraw(n, min, max, gradient, windowWidth, windowHeight)
	var frameStart  time.Time
	var elapsedTime  float32
	var controllerAxis int16
	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break
			}
		}

		for _, controller := range controllerHandlers {
			if controller != nil {
				controllerAxis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)

			}
		}
		
		if gState == play {
			drawNumber(getCenter(), color{255,255,255},20,2, pixels)
			player1.update(keyState, controllerAxis, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if gState == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}
				gState = play
			}
		}

		for i := range noisePixels {
			pixels[i] = noisePixels[i]
		}
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		if err := texture.Update(nil, pixels, windowWidth*4);err != nil {
			panic(err)
			return
		}
		if err := renderer.Copy(texture, nil,nil);err != nil {
			panic(err)
		}
		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime / 1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
		sdl.Delay(16)
	}
}
