package controllers

import (
	"context"
	"encoding/json"
	"github.com/Diamantto/go-backend-kursach/solver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

type Condition struct {
	Puzzle [][]int `json:"puzzle"`
}

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// Додайте інші поля за потребою
}

type Response struct {
	Counter  int     `json:"counter"`
	Solution [][]int `json:"solution"`
	Puzzle   [][]int `json:"puzzle"`
}

func getCurrentUser(r *http.Request) (User, error) {
	// Реалізуйте функцію для отримання користувача
	return User{}, nil
}

func solvePuzzleHandler(w http.ResponseWriter, r *http.Request) {
	var data Condition
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	solver := solver.NewHitoriSolver(data.Puzzle)
	result := solver.SolveHitori()

	if solution, ok := result["solution"]; ok {
		_, err = client.Database("main").Collection("history-puzzles").InsertOne(context.TODO(), bson.M{"puzzle": data.Puzzle, "user": user.ID})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		response := Response{
			Counter:  result["counter"].(int),
			Solution: solution.([][]int),
			Puzzle:   data.Puzzle,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	http.Error(w, result["error"].(string), http.StatusUnauthorized)
}

func generateMatrixHandler(w http.ResponseWriter, r *http.Request) {
	sizeStr := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}

	user, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	randomizer := solver.NewHitoriRandomizer()
	data, _ := randomizer.GenerateRandomPuzzle(size)

	_, err = client.Database("main").Collection("generated-matrices").InsertOne(context.TODO(), bson.M{"puzzle": data, "user": user.ID})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := Condition{Puzzle: data}
	json.NewEncoder(w).Encode(response)
}
