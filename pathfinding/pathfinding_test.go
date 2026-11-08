package pathfinding

import "testing"

func TestFindPathSimple(t *testing.T) {
	grid := [][]int{
		{0, 0, 0},
		{0, 1, 0},
		{0, 0, 0},
	}
	start := Point{X: 0, Y: 0}
	goal := Point{X: 2, Y: 2}

	path := FindPath(grid, start, goal)
	if path == nil {
		t.Fatal("Pfad sollte existieren")
	}
	if path[0] != start || path[len(path)-1] != goal {
		t.Fatalf("Pfad beginnt nicht am Start bzw. endet nicht am Ziel: %v", path)
	}
	if len(path) != 5 {
		t.Fatalf("erwartete Pfadlänge 5, bekam %d: %v", len(path), path)
	}
}

func TestFindPathImpossible(t *testing.T) {
	grid := [][]int{
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
	}
	path := FindPath(grid, Point{X: 0, Y: 0}, Point{X: 2, Y: 0})
	if path != nil {
		t.Fatalf("kein Pfad erwartet, bekam %v", path)
	}
}

// TestSolverSteps prüft, dass der Solver in endlich vielen Schritten fertig
// wird, jeder Schritt einen bisherigen Pfad liefert und der Endpfad mit dem
// des FindPath-Wrappers übereinstimmt.
func TestSolverSteps(t *testing.T) {
	grid := [][]int{
		{0, 0},
		{0, 0},
	}
	start := Point{X: 0, Y: 0}
	goal := Point{X: 1, Y: 1}

	s := NewSolver(grid, start, goal)
	if s.Done {
		t.Fatal("Solver sollte anfangs nicht fertig sein")
	}

	for i := 0; !s.Done; i++ {
		if i > len(grid)*len(grid[0]) {
			t.Fatal("Solver sollte nach maximal len(grid)*len(grid[0]) Schritten fertig sein")
		}
		s.Next()
		if s.CurrentPath() == nil {
			t.Fatal("jeder Zwischenschritt sollte den bisher gefundenen Pfad enthalten")
		}
	}

	if !s.Done {
		t.Fatal("Solver sollte fertig sein")
	}
	path := s.CurrentPath()
	if path == nil || path[0] != start || path[len(path)-1] != goal {
		t.Fatalf("unerwarteter Endpfad: %v", path)
	}
	if len(path) != len(FindPath(grid, start, goal)) {
		t.Fatal("Endpfad des Solvers sollte dem von FindPath entsprechen")
	}
}

// TestSolverImpossiblePath prüft, dass CurrentPath() nil liefert, wenn die
// Suche ohne Lösung endet.
func TestSolverImpossiblePath(t *testing.T) {
	grid := [][]int{
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
	}
	start := Point{X: 0, Y: 0}
	goal := Point{X: 2, Y: 0}

	s := NewSolver(grid, start, goal)
	for !s.Done {
		s.Next()
	}
	if !s.Done {
		t.Fatal("Solver sollte fertig sein")
	}
	if s.CurrentPath() != nil {
		t.Fatalf("ohne Lösung sollte CurrentPath() nil sein, bekam %v", s.CurrentPath())
	}
}

// serpentineGrid baut ein size x size-Gitter mit einem sich von oben nach
// unten windenden Korridor.
func serpentineGrid(size int) [][]int {
	grid := make([][]int, size)
	for y := range grid {
		grid[y] = make([]int, size)
		for x := range grid[y] {
			grid[y][x] = 1 // Vorgabe: Wand
		}
	}
	// Oberste Reihe komplett frei
	for x := 0; x < size; x++ {
		grid[0][x] = 0
	}
	// Abwechselnd verlaufende Korridore, verbunden an wechselnden Rändern
	for y := 2; y < size; y += 2 {
		for x := 1; x < size; x++ {
			grid[y][x] = 0
		}
		col := size - 1
		if (y/2)%2 == 1 {
			col = 1
		}
		grid[y-1][col] = 0
	}
	// Unterste Reihe als letzter Korridor
	for x := 1; x < size; x++ {
		grid[size-1][x] = 0
	}
	return grid
}

// TestFindPathSerpentine prüft das Serpentinen-Grid: lösbar, länger als die
// Luftlinie (der Weg windet sich also wirklich) und ohne Wandkontakt.
func TestFindPathSerpentine(t *testing.T) {
	grid := serpentineGrid(20)
	start := Point{X: 0, Y: 0}
	goal := Point{X: 19, Y: 19}

	path := FindPath(grid, start, goal)
	if path == nil {
		t.Fatal("der Serpentinen-Korridor muss lösbar sein")
	}
	// Pfadlänge (Anzahl Schritte) muss größer sein als die Manhattan-Distanz
	if len(path)-1 <= heuristic(start, goal) {
		t.Fatalf("der Weg sollte länger als die Luftlinie sein: %d <= %d", len(path)-1, heuristic(start, goal))
	}
	for i, p := range path {
		if grid[p.Y][p.X] != 0 {
			t.Fatalf("Pfadpunkt %d liegt in einer Wand: %v", i, p)
		}
	}
}
