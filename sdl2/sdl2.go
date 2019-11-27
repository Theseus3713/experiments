package main
// Игра с цветами
import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
windowWidth = 800
windowHeight = 600
)


type color struct {
	r byte
	g byte
	b byte
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
		return
	}
	defer  windows.Destroy()

	renderer, err := sdl.CreateRenderer(windows, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
		return
	}
	defer renderer.Destroy()
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, windowWidth, windowHeight)
	if err != nil {
		panic(err)
		return
	}
	defer texture.Destroy()

	pixels := make([]byte, windowWidth*windowHeight*4)

	for y := 0; y < windowHeight; y++ {
		for x := 0; x < windowWidth ; x++ {
			setPixel(x,y, color{byte(x % 255 ),byte(y % 255 ),0}, pixels)
		}
	}

	if err = texture.Update(nil, pixels, windowWidth*4); err != nil {
		panic(err)
	}
	if err = renderer.Copy(texture, nil,nil);err != nil {
		panic(err)
	}
	renderer.Present()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		sdl.Delay(16)
	}
}
