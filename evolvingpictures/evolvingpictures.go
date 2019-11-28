package main

import (
	"experiments/experiments/evolvingpictures/apt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"time"
)

const (
	windowWidth  = 1200
	windowHeight = 800
	windowDepth  = 100
)

type audioState struct {
	explosionBytes []byte
	deviceId       sdl.AudioDeviceID
	audiSpec       *sdl.AudioSpec
}

type mouseState struct {
	leftButton  bool
	rightButton bool
	x, y        int32
}

func getMouseState() mouseState {
	mouseX, MouseY, mouseButtonState := sdl.GetMouseState()
	return mouseState{
		leftButton:  mouseButtonState == sdl.ButtonLMask(),
		rightButton: mouseButtonState == sdl.ButtonRMask(),
		x:           mouseX,
		y:           MouseY,
	}
}

type rgba struct {
	r, g, b byte
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c rgba, pixels []byte) {
	index := (y*windowWidth + x) * 4

	if index <= len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	if tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h)); err != nil {
		panic(err)
	} else {
		if err := tex.Update(nil, pixels, w*4); err != nil {
			panic(err)
		}
		return tex
	}
}

func aptToTexture(redNode, greedNode, blueNode apt.Node, w, h int, renderer *sdl.Renderer) *sdl.Texture {
	var (
		scale        float32 = 255 / 2
		offset               = -1.0 * scale
		pixels               = make([]byte, w*h*4)
		pisxelsIndex         = 0
	)
	for yi := 0; yi < h; yi++ {
		y := float32(yi)/float32(h)*2 - 1
		for xi := 0; xi < w; xi++ {
			var (
				x = float32(xi)/float32(w)*2 - 1
				r = redNode.Eval(x, y)
				g = greedNode.Eval(x, y)
				b = blueNode.Eval(x, y)
			)
			pixels[pisxelsIndex] = byte(r*scale - offset)
			pisxelsIndex++
			pixels[pisxelsIndex] = byte(g*scale - offset)
			pisxelsIndex++
			pixels[pisxelsIndex] = byte(b*scale - offset)
			pisxelsIndex++
			pisxelsIndex++ // skip alpha
		}
	}
	return pixelsToTexture(renderer, pixels, w, h)
}

func main() {
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	windows, err := sdl.CreateWindow("Testing", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer windows.Destroy()

	renderer, err := sdl.CreateRenderer(windows, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	//var audioSpec *sdl.AudioSpec
	//explosionBytes, audioSpec := sdl.LoadWAV("C:/Users/xpoc_/go/src/experiments/balloons/explode.wav")
	//audioId, err := sdl.OpenAudioDevice("", false, audioSpec,nil,0)
	//if err != nil {
	//	panic(err)
	//}
	//audioState := audioState{
	//	explosionBytes: explosionBytes,
	//	deviceId:       audioId,
	//	audiSpec:       audioSpec,
	//}
	//defer sdl.FreeWAV(explosionBytes)
	rand.Seed(time.Now().UTC().UnixNano())
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	var (
		elapsedTime       float32
		currentMouseState = getMouseState()
		aptR              = apt.GetRandomNode()
		aptG              = apt.GetRandomNode()
		aptB              = apt.GetRandomNode()

	//prevMouseState = getMouseState()
	//noise = &apt.OpNoise{
	//	LeftChild:  &apt.OpX{},
	//	RightChild: &apt.OpY{},
	//}
	//atan2 = &apt.OpMult{
	//	LeftChild: &apt.OpX{},
	//	RightChild: noise,
	//}
	//sine = &apt.OpSin{
	//	Child: atan2,
	//}
	//plus = &apt.OpPlus{
	//	LeftChild:  &apt.OpY{},
	//	RightChild: sine,
	//}
	)
	var num = rand.Intn(20)
	for i := 0; i < num; i++ {
		aptR.AddRandom(apt.GetRandomNode())
	}
	num = rand.Intn(20)
	for i := 0; i < num; i++ {
		aptG.AddRandom(apt.GetRandomNode())
	}
	num = rand.Intn(20)
	for i := 0; i < num; i++ {
		aptB.AddRandom(apt.GetRandomNode())
	}

	for {
		_, nilCount := aptR.NodeCounts()
		if nilCount == 0 {
			break
		}
		aptR.AddRandom(apt.GetRandomLeaf())
	}
	for {
		_, nilCount := aptG.NodeCounts()
		if nilCount == 0 {
			break
		}
		aptG.AddRandom(apt.GetRandomLeaf())
	}
	for {
		_, nilCount := aptB.NodeCounts()
		if nilCount == 0 {
			break
		}
		aptB.AddRandom(apt.GetRandomLeaf())
	}

	tex := aptToTexture(aptR, aptG, aptB, windowWidth, windowHeight, renderer)

	for {
		frameStart := time.Now()

		currentMouseState = getMouseState()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.TouchFingerEvent:
				if e.Type == sdl.FINGERDOWN {
					currentMouseState.x = int32(e.X)
					currentMouseState.y = int32(e.Y)
					currentMouseState.leftButton = true
				}
			}
		}

		if err = renderer.Copy(tex, nil, nil); err != nil {
			panic(err)
		}

		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		//fmt.Println(`ms pre frame:`, elapsedTime)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
		//prevMouseState = currentMouseState
	}
}
