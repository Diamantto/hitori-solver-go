package solver

import (
	"errors"
	"strings"
)

type HitoriSolver struct {
	Puzzle [][]int
}

func NewHitoriSolver(puzzle [][]int) *HitoriSolver {
	return &HitoriSolver{Puzzle: puzzle}
}

func (solver *HitoriSolver) SolveHitori() map[string]interface{} {
	result, err := solver.solveSmart(solver.Puzzle)
	if err != nil {
		// Обробка помилки, наприклад, логування або повернення порожньої картки з помилкою
		return map[string]interface{}{
			"error": err.Error(),
		}
	}
	return result
}

func (solver *HitoriSolver) solveSmart(puzzle [][]int) (map[string]interface{}, error) {
	counter := 0
	nodes := [][]int{{0, 0}}
	initial, err := solver.loadPuzzle(puzzle)
	if err != nil {
		return nil, err
	}
	states := []map[string]interface{}{initial}

	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		state := states[len(states)-1]
		states = states[:len(states)-1]

		i, j := node[0], node[1]

		if solver.checkSolution(state) {
			return map[string]interface{}{
				"counter":  counter,
				"solution": state["puzzle"],
			}, nil
		}

		if i >= len(state["puzzle"].([][]int)) {
			continue
		}

		state = solver.cellSurrounded(state, i, j)
		if state == nil {
			continue
		}

		for _, d := range state["domain"].([][]string)[i][j] {
			if string(d) == "V" {
				states = append(states, state)
				if j+1 < len(state["puzzle"].([][]int)[i]) {
					nodes = append(nodes, []int{i, j + 1})
				} else {
					nodes = append(nodes, []int{i + 1, 0})
				}
				continue
			}

			if string(d) == "W" {
				whiteState := solver.copy(state)
				whiteState["domain"].([][]string)[i][j] = "W"
				counter++

				newDomainWhite := solver.fcWhite(whiteState, i, j)
				if newDomainWhite == nil {
					continue
				}

				states = append(states, whiteState)
				if j+1 < len(state["puzzle"].([][]int)[i]) {
					nodes = append(nodes, []int{i, j + 1})
				} else {
					nodes = append(nodes, []int{i + 1, 0})
				}
			}

			if string(d) == "B" && solver.blackAllowed(state, i, j) {
				blackState := solver.copy(state)
				blackState["puzzle"].([][]int)[i][j] = -1 // Representing a black cell as -1
				blackState["domain"].([][]string)[i][j] = "B"
				counter++

				if !solver.testWhiteConnected(blackState["domain"].([][]string), blackState["edges"].([][][][2]int)) {
					continue
				}

				blackState = solver.fcBlack(blackState, i, j)
				if blackState == nil {
					continue
				}

				states = append(states, blackState)
				if j+1 < len(state["puzzle"].([][]int)[i]) {
					nodes = append(nodes, []int{i, j + 1})
				} else {
					nodes = append(nodes, []int{i + 1, 0})
				}
			}
		}
	}
	return map[string]interface{}{"counter": counter}, nil
}

func (solver *HitoriSolver) loadPuzzle(puzzle [][]int) (map[string]interface{}, error) {
	for _, row := range puzzle {
		if len(row) != len(puzzle[0]) {
			return nil, errors.New("Usage: Puzzle should be size N x M")
		}
	}

	domain := make([][]string, len(puzzle))
	for i := range puzzle {
		domain[i] = make([]string, len(puzzle[i]))
		for j := range puzzle[i] {
			domain[i][j] = "BW"
		}
	}

	edges := make([][][][2]int, len(puzzle))
	for i := range puzzle {
		edges[i] = make([][][2]int, len(puzzle[i]))
		for j := range puzzle[i] {
			if err := solver.checkValues(puzzle[i][j]); err != nil {
				return nil, err
			}
			var edgesCell [][2]int
			u, d, l, r := i-1, i+1, j-1, j+1
			if u >= 0 {
				edgesCell = append(edgesCell, [2]int{u, j})
			}
			if d < len(puzzle) {
				edgesCell = append(edgesCell, [2]int{d, j})
			}
			if l >= 0 {
				edgesCell = append(edgesCell, [2]int{i, l})
			}
			if r < len(puzzle[i]) {
				edgesCell = append(edgesCell, [2]int{i, r})
			}
			edges[i][j] = edgesCell
		}
	}

	state := map[string]interface{}{
		"puzzle": puzzle,
		"domain": solver.uniqueWhiteCells(puzzle, domain),
		"edges":  edges,
	}

	return state, nil
}

func (solver *HitoriSolver) checkValues(value int) error {
	if value < 1 {
		return errors.New("Value should be positive")
	}
	return nil
}

func (solver *HitoriSolver) copy(state map[string]interface{}) map[string]interface{} {
	puzzle := make([][]int, len(state["puzzle"].([][]int)))
	for i, row := range state["puzzle"].([][]int) {
		puzzle[i] = append([]int{}, row...)
	}

	domain := make([][]string, len(state["domain"].([][]string)))
	for i, row := range state["domain"].([][]string) {
		domain[i] = append([]string{}, row...)
	}

	newState := map[string]interface{}{
		"puzzle": puzzle,
		"domain": solver.uniqueWhiteCells(puzzle, domain),
		"edges":  state["edges"],
	}
	return newState
}

func (solver *HitoriSolver) checkSolution(state map[string]interface{}) bool {
	return solver.domainComplete(state["domain"].([][]string)) && solver.testDuplicateNumber(state["puzzle"].([][]int))
}

func (solver *HitoriSolver) domainComplete(domain [][]string) bool {
	for _, row := range domain {
		for _, cell := range row {
			if len(cell) != 1 {
				return false
			}
		}
	}
	return true
}

func (solver *HitoriSolver) testDuplicateNumber(puzzle [][]int) bool {
	transpose := transposeMatrix(puzzle)
	for i := range puzzle {
		for j := range puzzle[i] {
			if puzzle[i][j] != 'B' {
				if solver.hasDuplicates(puzzle[i], transpose[j], puzzle[i][j]) {
					return false
				}
			}
		}
	}
	return true
}

func transposeMatrix(matrix [][]int) [][]int {
	if len(matrix) == 0 {
		return nil
	}
	transpose := make([][]int, len(matrix[0]))
	for i := range transpose {
		transpose[i] = make([]int, len(matrix))
		for j := range matrix {
			transpose[i][j] = matrix[j][i]
		}
	}
	return transpose
}

func (solver *HitoriSolver) hasDuplicates(row []int, col []int, value int) bool {
	rowCount, colCount := 0, 0
	for _, v := range row {
		if v == value {
			rowCount++
		}
	}
	for _, v := range col {
		if v == value {
			colCount++
		}
	}
	return rowCount > 1 || colCount > 1
}

func (solver *HitoriSolver) uniqueWhiteCells(puzzle [][]int, domain [][]string) [][]string {
	transpose := transposeMatrix(puzzle)
	for i := range puzzle {
		for j := range puzzle[i] {
			if !solver.hasDuplicates(puzzle[i], transpose[j], puzzle[i][j]) && puzzle[i][j] != 'B' {
				domain[i][j] = "V"
			}
		}
	}
	return domain
}

func (solver *HitoriSolver) blackAllowed(state map[string]interface{}, i, j int) bool {
	edges := state["edges"].([][][][2]int)
	puzzle := state["puzzle"].([][]int)
	for _, edge := range edges[i][j] {
		i2, j2 := edge[0], edge[1]
		if puzzle[i2][j2] == 'B' {
			return false
		}
	}
	return true
}

func (solver *HitoriSolver) testWhiteConnected(grid [][]string, edges [][][][2]int) bool {
	visited := make([][]bool, len(grid))
	for i := range grid {
		visited[i] = make([]bool, len(grid[i]))
	}

	node := solver.firstWhite(grid)
	if solver.totalWhiteCells(grid) == solver.touch(node, visited, grid, edges, 1) {
		return true
	}
	return false
}

func (solver *HitoriSolver) totalWhiteCells(grid [][]string) int {
	count := 0
	for _, row := range grid {
		for _, cell := range row {
			if cell != "B" {
				count++
			}
		}
	}
	return count
}

func (solver *HitoriSolver) touch(node [2]int, visited [][]bool, grid [][]string, edges [][][][2]int, count int) int {
	if node == [2]int{-1, -1} {
		return 0
	}
	i, j := node[0], node[1]
	visited[i][j] = true
	for _, edge := range edges[i][j] {
		i2, j2 := edge[0], edge[1]
		if grid[i2][j2] != "B" && !visited[i2][j2] {
			count = solver.touch([2]int{i2, j2}, visited, grid, edges, count) + 1
		}
	}
	return count
}

func (solver *HitoriSolver) firstWhite(grid [][]string) [2]int {
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] != "B" {
				return [2]int{i, j}
			}
		}
	}
	return [2]int{-1, -1}
}

func (solver *HitoriSolver) fcBlack(state map[string]interface{}, i, j int) map[string]interface{} {
	domain := state["domain"].([][]string)
	edges := state["edges"].([][][][2]int)

	for _, edge := range edges[i][j] {
		i2, j2 := edge[0], edge[1]
		domain[i2][j2] = strings.ReplaceAll(domain[i2][j2], "B", "")
		if domain[i2][j2] == "" {
			return nil
		}
	}
	state["domain"] = domain
	return state
}

func (solver *HitoriSolver) fcWhite(state map[string]interface{}, i, j int) map[string]interface{} {
	size := len(state["puzzle"].([][]int))
	domain := state["domain"].([][]string)
	puzzle := state["puzzle"].([][]int)

	k := 1
	for i-k >= 0 {
		if puzzle[i-k][j] == puzzle[i][j] {
			domain[i-k][j] = strings.ReplaceAll(domain[i-k][j], "W", "")
			if domain[i-k][j] == "" {
				return nil
			}
		}
		k++
	}
	k = 1
	for i+k < size {
		if puzzle[i+k][j] == puzzle[i][j] {
			domain[i+k][j] = strings.ReplaceAll(domain[i+k][j], "W", "")
			if domain[i+k][j] == "" {
				return nil
			}
		}
		k++
	}
	k = 1
	for j-k >= 0 {
		if puzzle[i][j-k] == puzzle[i][j] {
			domain[i][j-k] = strings.ReplaceAll(domain[i][j-k], "W", "")
			if domain[i][j-k] == "" {
				return nil
			}
		}
		k++
	}
	k = 1
	for j+k < size {
		if puzzle[i][j+k] == puzzle[i][j] {
			domain[i][j+k] = strings.ReplaceAll(domain[i][j+k], "W", "")
			if domain[i][j+k] == "" {
				return nil
			}
		}
		k++
	}
	state["domain"] = domain
	return state
}

func (solver *HitoriSolver) cellSurrounded(state map[string]interface{}, i, j int) map[string]interface{} {
	edges := state["edges"].([][][][2]int)
	domain := state["domain"].([][]string)

	count := 0
	for _, edge := range edges[i][j] {
		i2, j2 := edge[0], edge[1]
		if domain[i2][j2] == "B" {
			count++
		}
	}

	if count == len(edges[i][j])-1 {
		for _, edge := range edges[i][j] {
			i2, j2 := edge[0], edge[1]
			if domain[i2][j2] != "B" {
				domain[i2][j2] = strings.ReplaceAll(domain[i2][j2], "B", "")
				domain[i][j] = strings.ReplaceAll(domain[i][j], "B", "")
				if domain[i2][j2] == "" || domain[i][j] == "" {
					return nil
				}
			}
		}
	}
	state["domain"] = domain
	return state
}
