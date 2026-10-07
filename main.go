package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	// Options:
	// - nimble:9b
	// - tev1:4b

	model := "tev1:4b"

	matrix := [3][3]string{
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	for !is_finished(matrix) {
		state, options := get_state(matrix)
		x, y := make_call(state, options, "x", model)
		matrix[x][y] = "x"
		render_matrix(matrix)

		state, options = get_state(matrix)
		if len(options) == 0 {
			break
		}

		x, y = make_call(state, options, "o", model)
		matrix[x][y] = "o"
		render_matrix(matrix)
	}
}

func get_state(matrix [3][3]string) ([]string, map[string]string) {
	var state []string
	options := make(map[string]string)
	for x, row := range matrix {
		for y, cell := range row {
			xs := strconv.Itoa(x)
			ys := strconv.Itoa(y)
			if cell == "" {
				cell = "empty"
			}
			state = append(state, "'"+xs+","+ys+"'='"+cell+"'")
			if cell == "empty" {
				options[xs+","+ys] = "empty"
			}
		}
	}
	return state, options
}

func is_finished(matrix [3][3]string) bool {
	// check if tic-tac-toe is finished
	for _, row := range matrix {
		for _, cell := range row {
			if cell == "" {
				return false
			}
		}
	}
	return true
}

func render_matrix(matrix [3][3]string) {
	fmt.Println("Current state of the matrix:")
	for _, row := range matrix {
		fmt.Println(row)
	}
}

func make_call(state []string, options map[string]string, option string, model string) (int, int) {

	if len(options) == 1 {
		for pos := range options {
			x, _ := strconv.Atoi(string(pos[1]))
			y, _ := strconv.Atoi(string(pos[3]))
			return x, y
		}
	}

	var strOpt strings.Builder

	for key, val := range options {
		strOpt.WriteString(key + ": \"" + val + "\",")
	}

	requestData := &request{
		Model: model,
		State: stateRequest{
			TicTacToe: strings.Join(state, ", "),
		},
		Questions: questions{
			Position: positionRequest{
				Type:         "choice",
				Instructions: "You are playing tic-tac-toe, the table starts at '0,0', what is the next position to play if you play with " + option + ",",
				Criteria:     options,
			},
		},
	}

	data, err := json.Marshal(requestData)

	if err != nil {
		log.Fatal(err)
	}

	url := "http://localhost:11434/v1/systemone"

	println("Making call to " + url + " with data: " + string(data))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	if err != nil {
		log.Fatal(err)
	}

	body, _ := io.ReadAll(resp.Body)

	log.Println("Response:", string(body))

	var respData response
	err = json.Unmarshal(body, &respData)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Choice: " + respData.Answers.Position.Choice)

	choice := respData.Answers.Position.Choice

	x, _ := strconv.Atoi(string(choice[0]))
	y, _ := strconv.Atoi(string(choice[2]))

	return x, y
}

type response struct {
	Answers answers `json:"answers"`
}

type answers struct {
	Position position `json:"position"`
}

type position struct {
	Choice        string             `json:"choice"`
	Legend        map[string]string  `json:"legend"`
	Probabilities map[string]float32 `json:"probabilities"`
}

type request struct {
	Model     string       `json:"model"`
	State     stateRequest `json:"state"`
	Questions questions    `json:"questions"`
}

type stateRequest struct {
	TicTacToe string `json:"tic-tac-toe"`
}

type questions struct {
	Position positionRequest `json:"position"`
}

type positionRequest struct {
	Type         string            `json:"type"`
	Instructions string            `json:"instructions"`
	Criteria     map[string]string `json:"criteria"`
}
