package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	screenWidth  = 1000
	screenHeight = 480
)

var (
	running = true
	bgColor = rl.NewColor(147, 211, 196, 255)

	grassSprite  rl.Texture2D
	hillSprite   rl.Texture2D
	fenceSprite  rl.Texture2D
	houseSprite  rl.Texture2D
	tilledSprite rl.Texture2D
	waterSprite  rl.Texture2D
	tex          rl.Texture2D

	playerSprite rl.Texture2D
	playerSrc    rl.Rectangle
	playerDest   rl.Rectangle

	playerSpeed                                   float32 = 1.4
	playerMoving                                  bool
	playerDirection                               int
	playerUp, playerDown, playerRight, playerLeft bool
	playerFrame                                   int

	frameCount int

	tileDest   rl.Rectangle
	tileSrc    rl.Rectangle
	tileMap    []int
	srcMap     []string
	mapW, mapH int

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D
)

func drawScene() {
	for i := 0; i < len(tileMap); i++ {
		if tileMap[i] != 0 {
			// find tileDest coordinate according tileMap
			tileDest.X = tileDest.Width * float32(i%mapW)
			tileDest.Y = tileDest.Height * float32(i/mapW)

			if srcMap[i] == "g" {
				// find tileSrc from tile set to render
				tex = grassSprite
			}
			if srcMap[i] == "l" {
				// find tileSrc from tile set to render
				tex = hillSprite
			}
			if srcMap[i] == "f" {
				// find tileSrc from tile set to render
				tex = fenceSprite
			}
			if srcMap[i] == "h" {
				// find tileSrc from tile set to render
				tex = houseSprite
			}
			if srcMap[i] == "w" {
				// find tileSrc from tile set to render
				tex = waterSprite
			}
			if srcMap[i] == "t" {
				// find tileSrc from tile set to render
				tex = tilledSprite
			}

			if srcMap[i] == "h" || srcMap[i] == "f" {
				tileSrc.X = 0
				tileSrc.Y = 0
				rl.DrawTexturePro(
					grassSprite,
					tileSrc,
					tileDest,
					rl.NewVector2(tileDest.Width, tileDest.Height),
					0,
					rl.White,
				)
			}

			countNumberOfTileSet := int(tex.Width / int32(tileSrc.Width))
			tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%countNumberOfTileSet)
			tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/countNumberOfTileSet)

			rl.DrawTexturePro(
				tex,
				tileSrc,
				tileDest,
				rl.NewVector2(tileDest.Width, tileDest.Height),
				0,
				rl.White,
			)
		}
	}

	rl.DrawTexturePro(
		playerSprite,
		playerSrc,
		playerDest,
		rl.NewVector2(playerDest.Width, playerDest.Height),
		0,
		rl.White,
	)
}

func input() {
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDirection = 0
		playerDown = true
	}
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDirection = 1
		playerUp = true
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDirection = 2
		playerLeft = true
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDirection = 3
		playerRight = true
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		musicPaused = !musicPaused
	}
}

func update() {
	running = !rl.WindowShouldClose()

	playerSrc.X = playerSrc.Width * float32(playerFrame)
	if playerMoving {
		if playerUp {
			playerDest.Y -= playerSpeed
		}
		if playerDown {
			playerDest.Y += playerSpeed
		}
		if playerLeft {
			playerDest.X -= playerSpeed
		}
		if playerRight {
			playerDest.X += playerSpeed
		}

		// FPS
		if frameCount%8 == 1 {
			playerFrame++
			frameCount = 1
		}
	} else if frameCount%45 == 1 {
		playerFrame++
	}

	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}
	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
	}

	playerSrc.X = playerSrc.Width * float32(playerFrame)
	playerSrc.Y = playerSrc.Height * float32(playerDirection)

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(float32(playerDest.X-playerDest.Width/2), float32(playerDest.Y-playerDest.Height/2))

	playerMoving = false
	playerUp, playerDown, playerRight, playerLeft = false, false, false, false

}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bgColor)
	rl.BeginMode2D(cam)

	drawScene()

	rl.EndMode2D()
	rl.EndDrawing()
}

func loadMap(mapFile string) {
	file, err := ioutil.ReadFile(mapFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	removedNewLines := strings.Replace(string(file), "\n", " ", -1)
	sliced := strings.Split(removedNewLines, " ")

	mapW = -1
	mapH = -1
	for i := 0; i < len(sliced); i++ {
		s, _ := strconv.ParseInt(sliced[i], 10, 64)
		m := int(s)
		if mapW == -1 {
			mapW = m
		} else if mapH == -1 {
			mapH = m
		} else if i < mapW*mapH+2 {
			tileMap = append(tileMap, m)
		} else {
			srcMap = append(srcMap, sliced[i])
		}
	}

	//if len(tileMap) > mapW*mapH {
	//	tileMap = tileMap[:mapW*mapH]
	//}
}

func initialize() {
	rl.InitWindow(screenWidth, screenHeight, "Game")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("res/Tilesets/Grass.png")
	hillSprite = rl.LoadTexture("res/Tilesets/Hills.png")
	houseSprite = rl.LoadTexture("res/Tilesets/Wooden House.png")
	fenceSprite = rl.LoadTexture("res/Tilesets/Fences.png")
	tilledSprite = rl.LoadTexture("res/Tilesets/Tilled Dirt.png")
	waterSprite = rl.LoadTexture("res/Tilesets/Water.png")
	playerSprite = rl.LoadTexture("res/Characters/BasicCharakterSpriteSheet.png")

	tileSrc = rl.NewRectangle(0, 0, 16, 16)
	tileDest = rl.NewRectangle(0, 0, 16, 16)

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 60, 60)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("res/music.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(
		rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(playerDest.X-playerDest.Width/2), float32(playerDest.Y-playerDest.Height/2)),
		0,
		1.5,
	)

	cam.Zoom = 3

	loadMap("ep/two.map")
}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

func main() {
	initialize()
	for running {
		input()
		update()
		render()
	}
	quit()
}
