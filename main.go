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
	matrix := [3][3]string{
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	for !is_finished(matrix) {
		state, options := get_state(matrix)
		x, y := make_call(state, options, "x")
		matrix[x][y] = "x"
		render_matrix(matrix)

		state, options = get_state(matrix)
		if len(options) == 0 {
			break
		}

		x, y = make_call(state, options, "o")
		matrix[x][y] = "o"
		render_matrix(matrix)
	}
}

func get_state(matrix [3][3]string) ([]string, []string) {
	var state []string
	var options []string
	for x, row := range matrix {
		for y, cell := range row {
			xs := strconv.Itoa(x)
			ys := strconv.Itoa(y)
			if cell == "" {
				cell = "empty"
			}
			state = append(state, "("+xs+","+ys+")='"+cell+"'")
			if cell == "empty" {
				options = append(options, "\"("+xs+","+ys+")\"")
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

func make_call(state []string, options []string, option string) (int, int) {

	if len(options) == 1 {
		opt := options[0]
		x, _ := strconv.Atoi(string(opt[2]))
		y, _ := strconv.Atoi(string(opt[4]))
		return x, y
	}

	url := "http://localhost:11434/v1/systemone"

	// Options:
	// - nimble:9b
	// - tev1:4b

	data := "{" +
		"\"model\": \"tev1:4b\"," +
		"\"state\": {\"tic-tac-toe\": \"" + strings.Join(state, ", ") + "\"}," +
		"\"questions\": " +
		"{\"position\": " +
		"{" +
		"\"type\": \"score\"," +
		"\"instructions\": \"You are playing tic-tac-toe, what is the next position if you play with " + option + "\"," +
		"\"criteria\": [" + strings.Join(options, ", ") + "]" +
		"}" +
		"}" +
		"}"

	println("Making call to " + url + " with data: " + data)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(data)))

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

	var max float32
	var maxPos string

	for pos, prob := range respData.Answers.Position.Probabilities {

		if prob > max {
			max = prob
			maxPos = pos
		}
	}

	fmt.Printf("Position: %s, Probability: %f, Legend: %s\n", maxPos, max, respData.Answers.Position.Legend[maxPos])

	maxLegend := respData.Answers.Position.Legend[maxPos]
	fmt.Println("Max Legend: " + string(maxLegend[1]) + " - " + string(maxLegend[3]))

	x, _ := strconv.Atoi(string(maxLegend[1]))
	y, _ := strconv.Atoi(string(maxLegend[3]))

	return x, y
}

type response struct {
	Answers answers `json:"answers"`
}

type answers struct {
	Position position `json:"position"`
}

type position struct {
	Legend        map[string]string  `json:"legend"`
	Probabilities map[string]float32 `json:"probabilities"`
}
