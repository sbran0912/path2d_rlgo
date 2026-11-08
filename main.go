package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"

	"path2d_rlgo/pathfinding"
)

func main() {
	start := pathfinding.Point{X: 0, Y: 0}
	goal := pathfinding.Point{X: 19, Y: 19}

	// 0 = begehbar, 1 = Wand.
	// Zufälliges Grid mit ~30 % Wanddichte. Start und Ziel bleiben immer
	// frei; falls beide nicht zusammenhängen (sehr selten), wird neu gewürfelt.
	var grid [][]int
	for {
		grid = randomGrid(20, 20, 0.3)
		grid[start.Y][start.X] = 0
		grid[goal.Y][goal.X] = 0
		if pathfinding.FindPath(grid, start, goal) != nil {
			break
		}
	}

	runGridView(grid, start, goal)
}

// randomGrid erzeugt ein width x height-Grid (0 = frei, 1 = Wand) mit
// zufällig verteilten Wänden. math/rand ist ab Go 1.20 automatisch pro
// Programmstart anders seedet – jede Ausführung ergibt also ein anderes Gitter.
func randomGrid(width, height int, wallChance float64) [][]int {
	grid := make([][]int, height)
	for y := range grid {
		grid[y] = make([]int, width)
		for x := range grid[y] {
			grid[y][x] = 0
			if rand.Float64() < wallChance {
				grid[y][x] = 1
			}
		}
	}
	return grid
}

// runGridView zeigt das Gitter in einem raylib-Vollbildfenster an.
// Die A*-Suche läuft kontinuierlich: Pro Frame wird genau ein Suchschritt
// ausgeführt, und der Renderer liest den Zwischenstand direkt am Solver –
// inklusive des bisher gefundenen Pfads. Es gibt keine Tastaturabfragen.
func runGridView(grid [][]int, start, goal pathfinding.Point) {
	solver := pathfinding.NewSolver(grid, start, goal)

	gridW, gridH := gridSize(grid)

	// Vollbild: Fullscreen-Flag vor InitWindow setzen, danach auf die
	// Monitorauflösung stellen (wie im offiziellen raylib-Beispiel).
	rl.SetConfigFlags(rl.FlagFullscreenMode)
	rl.InitWindow(800, 450, "Pathfinding Grid")
	defer rl.CloseWindow()

	monitor := rl.GetCurrentMonitor()
	screenW := rl.GetMonitorWidth(monitor)
	screenH := rl.GetMonitorHeight(monitor)
	rl.SetWindowSize(screenW, screenH)

	// Zellgröße so wählen, dass das Grid mit Rand in den Vollbildschirm passt
	const pad float32 = 40
	cellW := (float32(screenW) - 2*pad) / float32(gridW)
	cellH := (float32(screenH) - 2*pad) / float32(gridH)
	cell := cellW
	if cellH < cell {
		cell = cellH
	}

	// Grid zentrieren, falls der Bildschirm nicht quadratisch ist
	originX := pad + (float32(screenW)-2*pad-float32(gridW)*cell)/2
	originY := pad + (float32(screenH)-2*pad-float32(gridH)*cell)/2

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Einen Suchschritt pro Frame, bis die Suche abgeschlossen ist.
		if !solver.Done {
			solver.Next()
			if solver.Done {
				if path := solver.CurrentPath(); path != nil {
					fmt.Printf("Pfad gefunden (%d Schritte): ", len(path))
					for _, p := range path {
						fmt.Printf("[%d, %d] ", p.X, p.Y)
					}
					fmt.Println()
				} else {
					fmt.Println("Kein Pfad möglich.")
				}
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)
		drawGrid(grid, originX, originY, cell, solver)
		rl.EndDrawing()
	}
}

// gridSize liefert die Ausmaße des Grids.
func gridSize(grid [][]int) (width, height int) {
	height = len(grid)
	if height > 0 {
		width = len(grid[0])
	}
	return width, height
}

// drawGrid zeichnet das Gitter mit dem aktuellen Suchfortschritt direkt aus
// dem Solver: 0 = frei (weiß), 1 = Wand (schwarz), verarbeitete Punkte rosa,
// offene Punkte hellblau, aktueller Punkt gelb, bisheriger Pfad grün.
func drawGrid(grid [][]int, originX, originY, size float32, solver *pathfinding.Solver) {
	// Erst alle Felder füllen, dann den Suchzustand, dann die Linien darüber.
	for y, row := range grid {
		for x, cell := range row {
			rec := cellRect(originX, originY, size, x, y)

			if cell == 1 {
				rl.DrawRectangleRec(rec, rl.Black)
			} else {
				rl.DrawRectangleRec(rec, rl.RayWhite)
			}
		}
	}

	// Bereits verarbeitete Punkte (rosa)
	for p := range solver.ClosedSet {
		rl.DrawRectangleRec(cellRect(originX, originY, size, p.X, p.Y), rl.NewColor(255, 150, 160, 255))
	}
	// Noch offene Punkte (hellblau)
	for _, p := range solver.OpenSet {
		rl.DrawRectangleRec(cellRect(originX, originY, size, p.X, p.Y), rl.NewColor(135, 206, 250, 255))
	}
	// Gerade verarbeiteter Punkt (gelb)
	rl.DrawRectangleRec(cellRect(originX, originY, size, solver.Current.X, solver.Current.Y), rl.Yellow)

	// Bisher gefundener Pfad (grün)
	for _, p := range solver.CurrentPath() {
		rl.DrawRectangleRec(cellRect(originX, originY, size, p.X, p.Y), rl.Green)
	}

	lineColor := rl.NewColor(60, 60, 66, 255)
	for y, row := range grid {
		for x := range row {
			rl.DrawRectangleLinesEx(cellRect(originX, originY, size, x, y), 1, lineColor)
		}
	}
}

func cellRect(originX, originY, size float32, x, y int) rl.Rectangle {
	return rl.Rectangle{
		X:      originX + float32(x)*size,
		Y:      originY + float32(y)*size,
		Width:  size,
		Height: size,
	}
}
