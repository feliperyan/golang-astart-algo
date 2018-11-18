package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth           = 640
	screenHeight          = 480
	maxAngle              = 256
	mapWidth              = 40
	mapHeight             = 25
	tunnels               = 130
	tunnelLength          = 15
	mapProportionModifier = 1
	screenFactor          = 2
)

type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
}

var tileSize int

var ebImg *ebiten.Image
var knightImg *ebiten.Image
var chestImg *ebiten.Image
var coinImg *ebiten.Image

var op = &ebiten.DrawImageOptions{}
var dungeon Map2d
var bob *MapElement
var gold *MapElement

var canReset chan bool

func init() {
	op := &ebiten.DrawImageOptions{}

	tileSize = 16
	tileSize = int(float64(tileSize) * mapProportionModifier)

	img, _, _ := ebitenutil.NewImageFromFile("images/floor_2.png", ebiten.FilterDefault)
	w, h := img.Size()
	w = int(float64(w) * mapProportionModifier)
	h = int(float64(h) * mapProportionModifier)
	ebImg, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	op.ColorM.Scale(1, 1, 1, 1.0)
	ebImg.DrawImage(img, op)

	img2, _, _ := ebitenutil.NewImageFromFile("images/knight_f_idle_anim_f0.png", ebiten.FilterDefault)
	w, h = img2.Size()
	w = int(float64(w) * mapProportionModifier)
	h = int(float64(h) * mapProportionModifier)
	knightImg, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	op.ColorM.Scale(1, 1, 1, 1.0)
	knightImg.DrawImage(img2, op)

	img3, _, _ := ebitenutil.NewImageFromFile("images/chest_empty_open_anim_f0.png", ebiten.FilterDefault)
	w, h = img3.Size()
	w = int(float64(w) * mapProportionModifier)
	h = int(float64(h) * mapProportionModifier)
	chestImg, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	op.ColorM.Scale(1, 1, 1, 1.0)
	chestImg.DrawImage(img3, op)

	img4, _, _ := ebitenutil.NewImageFromFile("images/coin_anim_f0.png", ebiten.FilterDefault)
	w, h = img4.Size()
	w = int(float64(w) * mapProportionModifier)
	h = int(float64(h) * mapProportionModifier)
	coinImg, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	op.ColorM.Scale(1, 1, 1, 1.0)
	coinImg.DrawImage(img4, op)

	dungeon = generateDungeon(mapWidth, mapHeight, tunnels, tunnelLength)

	bob = getRandomPosition(&dungeon, "b", false)
	gold = getRandomPosition(&dungeon, "g", true)

	// Reset key can be pressed once we start the program
	canReset = make(chan bool, 1)
	canReset <- true
}

func getRandomPosition(aMap *Map2d, name string, pass bool) *MapElement {
	rand.Seed(time.Now().UnixNano())
	var e *MapElement
	for {
		e = aMap.two_d[rand.Intn(aMap.x-1)][rand.Intn(aMap.y-1)]
		if e.passable {
			e, _ = putElementinMap2d(aMap, name, pass, e.pos_x, e.pos_y)
			break
		}
	}
	return e
}

func update(screen *ebiten.Image) error {

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		select {
		case action := <-canReset:
			fmt.Println("received spacebar", action)
			dungeon = generateDungeon(mapWidth, mapHeight, tunnels, tunnelLength)
			bob = getRandomPosition(&dungeon, "b", false)
			gold = getRandomPosition(&dungeon, "g", true)

			// wait a second before we can do it again otherwise
			go func() {
				time.Sleep(500 * time.Millisecond)
				canReset <- true
			}()
		default:
			fmt.Println("Multiple hits on spacebar, keeping it to one press.")
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		bob.path = aStarAlgorithm(&dungeon, bob, gold)
		fmt.Println(len(bob.path))
	}

	// Draw map
	w, h := tileSize, tileSize
	for colNum, row := range dungeon.two_d {
		for elNum, e := range row {
			op.GeoM.Reset()
			if e.name != "#" {
				op.GeoM.Translate(float64(colNum*h), float64(elNum*w))
				screen.DrawImage(ebImg, op)
			}
		}
	}

	if bob != nil {
		msg := fmt.Sprintf("bob %v %v", bob.pos_x, bob.pos_y)
		ebitenutil.DebugPrint(screen, msg)
		op.GeoM.Reset()
		op.GeoM.Translate(float64((bob.pos_x)*tileSize), float64((bob.pos_y-1)*tileSize))
		screen.DrawImage(knightImg, op)
	}
	if gold != nil {
		msg := fmt.Sprintf("         | gold %v %v", gold.pos_x, gold.pos_y)
		ebitenutil.DebugPrint(screen, msg)
		op.GeoM.Reset()
		op.GeoM.Translate(float64((gold.pos_x)*tileSize), float64((gold.pos_y)*tileSize))
		screen.DrawImage(chestImg, op)
	}

	if bob.path != nil {
		for _, p := range bob.path {
			op.GeoM.Reset()
			op.GeoM.Translate(float64(p.pos_x*tileSize)+2, float64(p.pos_y*tileSize)+2)
			screen.DrawImage(coinImg, op)
		}
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	return nil
}

func main() {

	if err := ebiten.Run(update, screenWidth, screenHeight, screenFactor, "FRyan Demo"); err != nil {
		log.Fatal(err)
	}
}
