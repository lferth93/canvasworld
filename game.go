package main

/*
#cgo LDFLAGS: -mwindows
*/

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

var (
	color = []rl.Color{
		rl.Black,
		rl.Red,
		rl.Lime,
		rl.Yellow,
		rl.DarkBlue,
		rl.Magenta,
		rl.SkyBlue,
		rl.White,
		rl.Beige}
)

const (
	pad = 20
)

type game struct {
	h, w   int32
	fps    int32
	l      int32
	wpad   int32
	step   int32
	count  int32
	paused bool
	r1, r2 float32
	world  *World
	path   string
}

func newGame(w, h, fps int32) *game {
	g := game{h: h, w: w, fps: fps, step: fps}
	rl.InitWindow(w, h+pad, "CanvasWorld Test")
	icon := rl.LoadImage("ico.PNG")
	rl.SetWindowIcon(*icon)
	rl.UnloadImage(icon)
	rl.SetTargetFPS(fps)
	return &g
}

func (g *game) loadWorld(path string) {
	g.world = LoadFile(path)
	g.path = path
	g.l = min(g.h/int32(g.world.h), g.w/int32(g.world.w))
	g.wpad = (g.w - g.l*int32(g.world.w)) / 2
	g.count = 0
	g.r1 = 0
	g.r2 = float32(g.l) / 4
}

func (g *game) readInput() {
	if rl.IsKeyPressed(rl.KeyUp) {
		if g.step > 1 {
			g.step--
		}
	}
	if rl.IsKeyPressed(rl.KeyDown) {
		g.step++
	}
	if rl.IsKeyPressed(rl.KeyR) {
		g.loadWorld(g.path)
	}
	if rl.IsKeyPressed(rl.KeySpace) {
		g.paused = !g.paused
	}
}

func (g *game) update() {
	if !g.paused {
		g.count++
		if g.count > g.step {
			g.world.Next()
			g.count = 0
		}
		g.r1 = float32(g.l) / 4 * float32(g.count) / float32(g.step)
		g.r2 = float32(g.l)/4 - g.r1
	}
}

func (g *game) draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	rl.DrawFPS(0, 0)
	rl.DrawText(fmt.Sprintf("Step: %.4fs", float32(g.step)/float32(g.fps)), 100, 0, 20, rl.Black)
	y := int32(pad)
	for i := range g.world.s {
		x := g.wpad
		for j, c := range g.world.s[i] {
			rl.DrawRectangle(x, y, g.l, g.l, color[c])
			rl.DrawRectangleLines(x, y, g.l, g.l, rl.Gray)
			if g.world.a[i][j] == 1 {
				rl.DrawCircle(x+g.l/2, y+g.l/2, g.r1, rl.Green)
			}
			if g.world.a[i][j] == 2 {
				rl.DrawCircle(x+g.l/2, y+g.l/2, g.r2, rl.Green)
			}
			x += g.l
		}
		y += g.l
	}
	rl.EndDrawing()
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
