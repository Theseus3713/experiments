package main

import (
	"experiments/noise"
	"experiments/vec3"
	"github.com/veandco/go-sdl2/sdl"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

const (
	windowWidth = 1200
	windowHeight = 800
	windowDepth = 100
)

type audioState struct {
	explosionBytes []byte
	deviceId sdl.AudioDeviceID
	audiSpec *sdl.AudioSpec
}

type mouseState struct {
	leftButton bool
	rightButton bool
	x, y int32
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

type balloon struct {
	texture *sdl.Texture
	pos vec3.Vector3
	dir vec3.Vector3
	w, h int

	exploding bool
	exploded bool
	explosionStart time.Time
	explosionInterval float32
	explosionTexture *sdl.Texture
}

func newBalloon(tex *sdl.Texture, pos, dir vec3.Vector3, explosionTexture *sdl.Texture) *balloon {
	if _, _, w, h, err := tex.Query(); err != nil {
		log.Fatal(err)
		return nil
	} else {
		return &balloon{tex,pos,dir,int(w),int(h), false,false,
			time.Now(), 20,explosionTexture}
	}
}

type balloonArray []*balloon

func (b balloonArray) Len() int {
	return len(b)
}

func (b balloonArray) Swap(i, j int)  {
	b[i], b[j] = b[j], b[i]
}

func (b balloonArray) Less(i, j int) bool {
	diff :=  b[i].pos.Z - b[j].pos.Z
	return diff < -.0
}

func (b *balloon) getScale() float32 {
	return (b.pos.Z / 200 + 1) / 4
}

func (b *balloon) getCircle() (float32, float32, float32) {
	var (
		x = b.pos.X
		y = b.pos.Y - 30 * b.getScale()
		r = float32(b.w) / 2 * b.getScale()
	)
	return x, y, r
}

func updateBalloons(balloons []*balloon, elapsedTime float32, currentMouseState, prevMouseState mouseState, audioState *audioState) []*balloon {
	var (
		numAnimations = 16
		balloonClicked = false
		balloonsExploded = false
	)
	for i:= len(balloons) - 1; i >= 0; i-- {
		var balloon = balloons[i]

		if balloon.exploding {
			var (
				animationElapsed = float32(time.Since(balloon.explosionStart).Seconds() * 1000)
				animationIndex = int32(numAnimations - 1 - int(animationElapsed / balloon.explosionInterval))
			)
			if animationIndex < 0 {
				balloon.exploding = false
				balloon.exploded = true
				balloonsExploded = true
			}
		}

		if !balloonClicked && !prevMouseState.leftButton && currentMouseState.leftButton {
			x, y, r := balloon.getCircle()
			var (
				mouseX = currentMouseState.x
				mouseY = currentMouseState.y
				xDeff  = float32(mouseX) - x
				yDeff  = float32(mouseY) - y
				dest   = float32(math.Sqrt(float64(xDeff*xDeff + yDeff*yDeff)))
			)
			if dest < r {
				balloonClicked = true
				sdl.ClearQueuedAudio(audioState.deviceId)
				if err := sdl.QueueAudio(audioState.deviceId, audioState.explosionBytes); err != nil {
					log.Fatal(err)
				}
				sdl.PauseAudioDevice(audioState.deviceId, false)
				balloon.exploding = true
				balloon.explosionStart = time.Now()
			}
		}

		p := vec3.Add(balloon.pos, vec3.Mult(balloon.dir, elapsedTime))
		if p.X < 0 || p.X > float32(windowWidth) {
			balloon.dir.X = -balloon.dir.X
		}
		if p.Y < 0 || p.Y > float32(windowHeight) {
			balloon.dir.Y = -balloon.dir.Y
		}
		if p.Z < 0 || p.Z > float32(windowDepth) {
			balloon.dir.Z = -balloon.dir.Z
		}
		balloon.pos = vec3.Add(balloon.pos, vec3.Mult(balloon.dir, elapsedTime))
	}

	if balloonsExploded {
		var filteredBalloons = balloons[0:0]
		for _, balloon := range balloons {
			if !balloon.exploded {
				filteredBalloons = append(filteredBalloons, balloon)
			}
		}
		balloons = filteredBalloons
	}
	return balloons
}

func (b *balloon) draw(renderer *sdl.Renderer)  {

	var (
		scale = b.getScale()
		w = int32(float32(b.w) * scale)
		h = int32(float32(b.h) * scale)
		x = int32(b.pos.X - float32(w / 2))
		y = int32(b.pos.Y - float32(h / 2))

	)
	rect := &sdl.Rect{x,y,w,h}
	if err := renderer.Copy(b.texture, nil, rect); err != nil {
		panic(err)
	}

	if b.exploding {
		var (
			imageSize int32 = 64
			numAnimations = 16
			animationElapsed = float32(time.Since(b.explosionStart).Seconds() * 1000)
			animationIndex = int32(numAnimations - 1 - int(animationElapsed / b.explosionInterval))
			animationX = animationIndex % 4
			animationY = imageSize * ((animationIndex - animationX) / 4)
			animationRect = &sdl.Rect{animationX * imageSize, animationY, imageSize, imageSize}
		)
		rect.X -= rect.W / 2
		rect.Y -= rect.H / 2
		rect.W *= 2
		rect.H *= 2
		if err := renderer.Copy(b.explosionTexture, animationRect, rect); err != nil {
			log.Fatal(err)
		}
	}
}

type rgba struct {
	r, g, b byte
}


func clear(pixels []byte)  {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c rgba, pixels []byte) {
	index := (y* windowWidth + x ) * 4

	if index <= len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index + 1] = c.g
		pixels[index + 2] = c.b
	}
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	if tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h)); err != nil {
		panic(err)
	} else {
		if err := tex.Update(nil, pixels, w * 4); err != nil {
			panic(err)
		}
		return tex
	}
}

func imgFileToTexture(renderer *sdl.Renderer, fileName string) *sdl.Texture {

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w * h * 4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	tex := pixelsToTexture(renderer, pixels, w, h)
	if err := tex.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		log.Fatal(err)
		return nil
	}
	return tex
}

func loadBalloons(renderer *sdl.Renderer, numBalloons int) []*balloon {

	explosionTexture := imgFileToTexture(renderer, "C:/Users/xpoc_/go/src/experiments/balloons/explosion.png")

	balloonNames := []string {
		"C:/Users/xpoc_/go/src/experiments/balloons/balloon_red.png",
		"C:/Users/xpoc_/go/src/experiments/balloons/balloon_green.png",
		"C:/Users/xpoc_/go/src/experiments/balloons/balloon_blue.png"}

	balloonTextures := make([]*sdl.Texture, len(balloonNames))

	for i, name := range balloonNames {
		balloonTextures[i] = imgFileToTexture(renderer, name)
	}
	balloons := make([]*balloon, numBalloons)
	for i := range balloons {
		var (
			tex = balloonTextures[i % 3]
			pos = vec3.Vector3{rand.Float32() * windowWidth, rand.Float32() * windowHeight, rand.Float32() * windowDepth}
			dir = vec3.Vector3{rand.Float32() * .5 -.25,rand.Float32() * .5 -.25,rand.Float32() * .5}
		)
		balloons[i] = newBalloon(tex, pos, dir, explosionTexture)
	}
	return balloons
}

func flerp(b1, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct * (float32(b2) - float32(b1)))
}

func colorLerp(c1, c2 rgba, pct float32) rgba {
	return rgba{flerp(c1.r, c2.r, pct), flerp(c1.g, c2.g, pct), flerp(c1.b, c2.b, pct)}
}

func getGradient(c1, c2 rgba) []rgba {
	result := make([]rgba, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}
	return result
}

func getDualGradient(c1, c2 , c3, c4 rgba) []rgba {
	result := make([]rgba, 256)
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

func rescaleAndDraw(noise []float32, min, max float32, gradient []rgba, w, h int) []byte {
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

func main()  {
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
	defer  windows.Destroy()

	renderer, err := sdl.CreateRenderer(windows, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	var audioSpec *sdl.AudioSpec
	explosionBytes, audioSpec := sdl.LoadWAV("C:/Users/xpoc_/go/src/experiments/balloons/explode.wav")
	audioId, err := sdl.OpenAudioDevice("", false, audioSpec,nil,0)
	if err != nil {
		panic(err)
	}
	audioState := audioState{
		explosionBytes: explosionBytes,
		deviceId:       audioId,
		audiSpec:       audioSpec,
	}
	defer sdl.FreeWAV(explosionBytes)

	cloudNoise, min, max := noise.MakeNoise(noise.FBM, .009, .5,3,3, windowWidth, windowHeight)
	cloudGradient := getGradient(rgba{0,0,255}, rgba{255,255,255})
	cloudPixels := rescaleAndDraw(cloudNoise, min, max, cloudGradient, windowWidth, windowHeight)
	cloudTexture := pixelsToTexture(renderer, cloudPixels, windowWidth, windowHeight)
	balloons := loadBalloons(renderer, 10)
	var
	(
		elapsedTime float32
		currentMouseState = getMouseState()
		prevMouseState = getMouseState()
	)

	for {
		frameStart := time.Now()

		currentMouseState = getMouseState()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e:= event.(type) {
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

		if err = renderer.Copy(cloudTexture, nil,nil); err != nil {
			panic(err)
		}

		balloons = updateBalloons(balloons , elapsedTime, currentMouseState, prevMouseState, &audioState)

		sort.Stable(balloonArray(balloons))
		for _, balloon := range balloons {
			balloon.draw(renderer)
		}

		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		//fmt.Println(`ms pre frame:`, elapsedTime)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
		prevMouseState = currentMouseState
	}
}
