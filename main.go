package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	color := []rl.Color{
		rl.Black,
		rl.Red,
		rl.Lime,
		rl.Yellow,
		rl.DarkBlue,
		rl.Magenta,
		rl.SkyBlue,
		rl.White,
		rl.Beige}

	height := int32(500)
	widht := int32(500)
	pad := int32(20)
	frames := 60
	fps := int32(60)
	rl.InitWindow(widht, height+pad, "CanvasWorld Test")
	world := LoadFile("w1.json")
	l := min(height/int32(world.h), widht/int32(world.w))
	wpad := int32(rl.GetScreenWidth()-int(l)*world.w) / 2
	rl.SetTargetFPS(fps)
	count := 0
	for !rl.WindowShouldClose() {
		if rl.IsKeyReleased(rl.KeyUp) {
			frames--
			if frames < 1 {
				frames = 1
			}
		}
		if rl.IsKeyReleased(rl.KeyDown) {
			frames++
		}
		if rl.IsKeyReleased(rl.KeyL) {
			world.Save("Res.json")
			world = LoadFile("w1.json")
			wpad = int32(rl.GetScreenWidth()-int(l)*world.w) / 2
			count = 0
		}
		count++
		r1 := float32(l) / 4 * float32(count) / float32(frames)
		r2 := float32(l)/4 - r1
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		rl.DrawFPS(0, 0)
		rl.DrawText(fmt.Sprintf("Step: %.4fs", float32(frames)/float32(fps)), 100, 0, 20, rl.Black)
		for i := range world.s {
			for j, c := range world.s[i] {
				rl.DrawRectangle(l*int32(j)+wpad, l*int32(i)+pad, l, l, color[c])
				rl.DrawRectangleLines(l*int32(j)+wpad, l*int32(i)+pad, l, l, rl.Gray)
				if world.a[i][j] == 1 {
					rl.DrawCircle(l*int32(j)+l/2+wpad, l*int32(i)+pad+l/2, r1, rl.Green)
				}
				if world.a[i][j] == 2 {
					rl.DrawCircle(l*int32(j)+l/2+wpad, l*int32(i)+pad+l/2, r2, rl.Green)
				}
			}
		}
		if count >= frames {
			world.Next()
			count = 0
		}
		rl.EndDrawing()
	}
	rl.CloseWindow()

}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
