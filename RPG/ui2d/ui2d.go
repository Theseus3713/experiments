package ui2d

import (
	"bufio"
	"experiments/experiments/RPG/game"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type ui struct {
	windowWidth       int32
	windowHeight      int32
	renderer          *sdl.Renderer
	window            *sdl.Window
	textureAtlas      *sdl.Texture
	textureIndex      map[game.Title][]sdl.Rect
	keyboardState     []uint8
	prevKeyboardState []uint8
	centerX           int
	centerY           int
	r                 *rand.Rand
	levelChan         chan *game.Level
	inputChan         chan *game.Input
	fontSmall         *ttf.Font
	fontMedium        *ttf.Font
	fontLarge         *ttf.Font

	stringTextureSmall  map[string]*sdl.Texture
	stringTextureMedium map[string]*sdl.Texture
	stringTextureLarge  map[string]*sdl.Texture
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Level) *ui {

	var newUI = &ui{
		inputChan:           inputChan,
		stringTextureSmall:  make(map[string]*sdl.Texture),
		stringTextureMedium: make(map[string]*sdl.Texture),
		stringTextureLarge:  make(map[string]*sdl.Texture),
		levelChan:           levelChan,
		windowWidth:         1280,
		windowHeight:        720,
		r:                   rand.New(rand.NewSource(1)),
	}

	window, err := sdl.CreateWindow("RPG", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, newUI.windowWidth, newUI.windowHeight,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	newUI.window = window
	renderer, err := sdl.CreateRenderer(newUI.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	newUI.renderer = renderer
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	newUI.textureAtlas = newUI.imgFileToTexture("C:/Users/xpoc_/go/src/experiments/experiments/RPG/ui2d/assets/tiles.png")
	newUI.loadTextureIndex("C:/Users/xpoc_/go/src/experiments/experiments/RPG/ui2d/assets/atlas-index.txt")
	newUI.keyboardState = sdl.GetKeyboardState()
	newUI.prevKeyboardState = make([]uint8, len(newUI.keyboardState))
	for i, v := range newUI.keyboardState {
		newUI.prevKeyboardState[i] = v
	}
	newUI.centerX = -1
	newUI.centerY = -1

	if newUI.fontSmall, err = ttf.OpenFont("C:/Users/xpoc_/go/src/experiments/experiments/RPG/ui2d/assets/Kingthings_Foundation.ttf", 16); err != nil {
		panic(err)
	}
	if newUI.fontMedium, err = ttf.OpenFont("C:/Users/xpoc_/go/src/experiments/experiments/RPG/ui2d/assets/Kingthings_Foundation.ttf", 32); err != nil {
		panic(err)
	}
	if newUI.fontLarge, err = ttf.OpenFont("C:/Users/xpoc_/go/src/experiments/experiments/RPG/ui2d/assets/Kingthings_Foundation.ttf", 64); err != nil {
		panic(err)
	}
	return newUI
}

type FountSize int

const (
	FontSmall  FountSize = 16
	FontMedium           = 32
	FontLarge            = 64
)

func (ui *ui) stringToTexture(str string, color sdl.Color, size FountSize) *sdl.Texture {
	var font *ttf.Font
	switch size {
	case FontSmall:
		if tex, ok := ui.stringTextureSmall[str]; ok {
			return tex
		} else {
			font = ui.fontSmall
		}
	case FontMedium:
		if tex, ok := ui.stringTextureMedium[str]; ok {
			return tex
		} else {
			font = ui.fontMedium
		}
	case FontLarge:
		if tex, ok := ui.stringTextureLarge[str]; ok {
			return tex
		} else {
			font = ui.fontLarge
		}
	default:
		panic(`FountSize is not exists`)
	}

	solid, err := font.RenderUTF8Blended(str, color)
	if err != nil {
		panic(err)
	}
	defer solid.Free()
	tex, err := ui.renderer.CreateTextureFromSurface(solid)
	if err != nil {
		panic(err)
	}

	switch size {
	case FontSmall:
		ui.stringTextureSmall[str] = tex
	case FontMedium:
		ui.stringTextureMedium[str] = tex
	case FontLarge:
		ui.stringTextureLarge[str] = tex
	default:
		panic(`FountSize is not exists`)
	}
	//defer tex.Destroy()
	return tex
}

type UI2d struct {
}

func (ui *ui) loadTextureIndex(fileName string) {
	ui.textureIndex = make(map[game.Title][]sdl.Rect)
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
		ui.textureIndex[tileRune] = rects
	}

}

func (ui *ui) imgFileToTexture(fileName string) *sdl.Texture {
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

	tex := pixelsToTexture(ui.renderer, pixels, w, h)
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
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if err := ttf.Init(); err != nil {
		panic(err)
	}
}

func (ui *ui) Draw(level *game.Level) {
	if ui.centerX == -1 && ui.centerY == -1 {
		ui.centerY = level.Player.Y
		ui.centerX = level.Player.X
	}
	// центрирование камеры
	var limit = 5
	if level.Player.X > ui.centerX+limit {
		ui.centerX++
	} else if level.Player.X < ui.centerX-limit {
		ui.centerX--
	} else if level.Player.Y > ui.centerY+limit {
		ui.centerY++
	} else if level.Player.Y < ui.centerY-limit {
		ui.centerY--
	}
	var (
		offsetX = ui.windowWidth/2 - int32(ui.centerX*32)
		offsetY = ui.windowHeight/2 - int32(ui.centerY*32)
	)
	ui.renderer.Clear()
	ui.r.Seed(1)
	for y, row := range level.Map {
		for x, tile := range row {
			if tile == game.Blank {
				continue
			}
			var (
				scrRects = ui.textureIndex[tile]
				scrRect  = scrRects[ui.r.Intn(len(scrRects))]
				dstRect  = sdl.Rect{int32(x*32) + offsetX, int32(y*32) + offsetY, 32, 32}
			)
			if level.Debug[game.Position{x, y}] {
				ui.textureAtlas.SetColorMod(128, 0, 0)
			} else {
				ui.textureAtlas.SetColorMod(255, 255, 255)
			}
			ui.renderer.Copy(ui.textureAtlas, &scrRect, &dstRect)
		}
	}

	for pos, monster := range level.Monsters {
		monsterSrcRect := ui.textureIndex[game.Title(monster.Rune)][0]
		ui.renderer.Copy(ui.textureAtlas, &monsterSrcRect, &sdl.Rect{int32(pos.X)*32 + offsetX, int32(pos.Y)*32 + offsetY, 32, 32})
	}
	playerSrcRect := ui.textureIndex['@'][0]
	if err := ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &sdl.Rect{int32(level.Player.X)*32 + offsetX, int32(level.Player.Y)*32 + offsetY, 32, 32}); err != nil {
		panic(err)
	}

	textStart := int32(float64(ui.windowHeight) * .75)
	for i, event := range level.Events {
		if event != "" {
			tex := ui.stringToTexture(event, sdl.Color{255, 0, 0, 0}, FontSmall)
			if _, _, w, h, err := tex.Query(); err != nil {
				panic(err)
			} else {
				ui.renderer.Copy(tex, nil, &sdl.Rect{0, int32(i*int(FontSmall)) + textStart, w, h})
			}
		}
	}

	ui.renderer.Present()
}

func (ui *ui) Run() {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("QuitGame")
				ui.inputChan <- &game.Input{Type: game.QuitGame}
			case *sdl.WindowEvent:
				switch e.Event {
				case sdl.WINDOWEVENT_CLOSE:
					fmt.Println("CloseWindow")
					ui.inputChan <- &game.Input{Type: game.CloseWindow, LevelChannel: ui.levelChan}
				}
			}
		}

		select {
		case newLevel, ok := <-ui.levelChan:
			if ok {
				ui.Draw(newLevel)
			}
		default:
		}

		if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {

			var input game.Input
			if ui.keyboardState[sdl.SCANCODE_UP] == 1 && ui.prevKeyboardState[sdl.SCANCODE_UP] == 0 {
				input.Type = game.Up
			}
			if ui.keyboardState[sdl.SCANCODE_DOWN] == 1 && ui.prevKeyboardState[sdl.SCANCODE_DOWN] == 0 {
				input.Type = game.Down
			}
			if ui.keyboardState[sdl.SCANCODE_LEFT] == 1 && ui.prevKeyboardState[sdl.SCANCODE_LEFT] == 0 {
				input.Type = game.Left
			}
			if ui.keyboardState[sdl.SCANCODE_RIGHT] == 1 && ui.prevKeyboardState[sdl.SCANCODE_RIGHT] == 0 {
				input.Type = game.Right
			}
			if ui.keyboardState[sdl.SCANCODE_S] == 1 && ui.prevKeyboardState[sdl.SCANCODE_S] == 0 {
				input.Type = game.Search
			}
			for i, v := range ui.keyboardState {
				ui.prevKeyboardState[i] = v
			}

			if input.Type != game.None {
				ui.inputChan <- &input
			}
		}
		sdl.Delay(10)
	}
}
