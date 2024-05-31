package solver

import (
	"errors"
	"math/rand"
	"time"
)

// HitoriRandomizer структура для випадкового генерування пазлів Hitori
type HitoriRandomizer struct {
	puzzles map[string][][][]int
}

// NewHitoriRandomizer створює новий екземпляр HitoriRandomizer
func NewHitoriRandomizer() *HitoriRandomizer {
	return &HitoriRandomizer{
		puzzles: puzzles,
	}
}

// GenerateRandomPuzzle генерує випадковий пазл Hitori заданого розміру
func (hr *HitoriRandomizer) GenerateRandomPuzzle(size int) ([][]int, error) {
	if size < 2 {
		return nil, errors.New("Size should be greater than 1")
	}
	puzzleList, ok := hr.puzzles[string(size)]
	if !ok {
		return nil, errors.New("Puzzle size not found")
	}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(puzzleList))
	return puzzleList[randomIndex], nil
}
