package main

import "github.com/gen2brain/raylib-go/raylib"

const (
	screenWidth  = 1000
	screenHeight = 480
)

var (
	running      = true
	bgColor      = rl.NewColor(147, 211, 196, 255)
	grassSprite  rl.Texture2D
	playerSprite rl.Texture2D
	playerSrc    rl.Rectangle
	playerDest   rl.Rectangle

	playerSpeed float32 = 3

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D
)

func drawScene() {
	rl.DrawTexture(grassSprite, 100, 50, rl.White)
	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
}

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerDest.Y -= playerSpeed
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerDest.Y += playerSpeed
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerDest.X -= playerSpeed
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerDest.X += playerSpeed
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		musicPaused = !musicPaused
	}
}

func update() {
	running = !rl.WindowShouldClose()

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(float32(playerDest.X-playerDest.Width/2), float32(playerDest.Y-playerDest.Height/2))
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bgColor)
	rl.BeginMode2D(cam)

	drawScene()

	rl.EndMode2D()
	rl.EndDrawing()
}

func initialize() {
	rl.InitWindow(screenWidth, screenHeight, "Game")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("res/Tilesets/Grass.png")
	playerSprite = rl.LoadTexture("res/Characters/BasicCharakterSpriteSheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(100, 100, 100, 100)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("res/music.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(
		rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(playerDest.X-playerDest.Width/2), float32(playerDest.Y-playerDest.Height/2)),
		0,
		1.0,
	)
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
