package ui2d

import (
	"bufio"
	"experiments/experiments/RPG/game"
	"github.com/veandco/go-sdl2/sdl"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	windowWidth  = 1280
	windowHeight = 720
)

var (
	renderer     *sdl.Renderer
	textureAtlas *sdl.Texture
	textureIndex map[game.Title][]sdl.Rect
)

func loadTextureIndex(fileName string) {
	textureIndex = make(map[game.Title][]sdl.Rect)
	// C:\Users\xpoc_\go\src\experiments\experiments\RPG\ui2d\assets
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		var (
			line     = strings.TrimSpace(scanner.Text())
			tileRune = game.Title(line[0])
			xy       = line[1:]
			splitXYC = strings.Split(xy, ",")
		)
		x, err := strconv.ParseInt(strings.TrimSpace(splitXYC[0]), 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXYC[1]), 10, 64)
		if err != nil {
			panic(err)
		}

		variationCount, err := strconv.ParseInt(strings.TrimSpace(splitXYC[2]), 10, 64)
		if err != nil {
			panic(err)
		}

		var rects []sdl.Rect
		for i := int64(0); i < variationCount; i++ {
			rects = append(rects, sdl.Rect{int32(x * 32), int32(y * 32), 32, 32})
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
		textureIndex[tileRune] = rects
	}

}

func imgFileToTexture(fileName string) *sdl.Texture {
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

	pixels := make([]byte, w*h*4)
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

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	if tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h)); err != nil {
		panic(err)
	} else {
		if err := tex.Update(nil, pixels, w*4); err != nil {
			panic(err)
		}
		return tex
	}
}

func init() {
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	windows, err := sdl.CreateWindow("RPG", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(windows, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	textureAtlas = imgFileToTexture("C:/Users/xpoc_/go/src/experiments/experiments/RPG/ui2d/assets/tiles.png")
	loadTextureIndex("C:/Users/xpoc_/go/src/experiments/experiments/RPG/ui2d/assets/atlas-index.txt")
}

type UI2d struct {
}

func (ui *UI2d) DrawThenGetInput(level *game.Level) game.Input {
	rand.Seed(1)
	for y, row := range level.Map {
		for x, tile := range row {
			if tile == game.Blank {
				continue
			}
			var (
				scrRects = textureIndex[tile]
				scrRect  = scrRects[rand.Intn(len(scrRects))]
				dstRect  = sdl.Rect{int32(x * 32), int32(y * 32), 32, 32}
			)
			renderer.Copy(textureAtlas, &scrRect, &dstRect)
		}
	}

	if err := renderer.Copy(textureAtlas, &sdl.Rect{21 * 32, 59 * 32, 32, 32}, &sdl.Rect{int32(level.Player.X) * 32, int32(level.Player.Y) * 32, 32, 32}); err != nil {
		panic(err)
	}
	renderer.Present()
	for {

	}
}
