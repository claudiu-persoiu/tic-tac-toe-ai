package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/claudiu-persoiu/tic-tac-toe-ai/internal/dto"
)

func main() {

	model := selectModel()
	prompt := askPrompt()

	matrix := [3][3]string{
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
	}

	for !isFinished(matrix) {
		state, options := getState(matrix)
		x, y := makeCall(state, options, "x", model, prompt)
		matrix[x][y] = "x"
		renderMatrix(matrix)

		state, options = getState(matrix)
		if len(options) == 0 {
			break
		}

		x, y = makeCall(state, options, "o", model, prompt)
		matrix[x][y] = "o"
		renderMatrix(matrix)
	}
}

var models = []string{"nimble:9b", "tev1:4b"}

const defaultPrompt = "You want to win"

var reader = bufio.NewReader(os.Stdin)

func selectModel() string {
	for {
		fmt.Println("Select the model to use:")
		for i, m := range models {
			fmt.Printf("%d) %s\n", i+1, m)
		}
		fmt.Print("Choice [1]: ")

		line, err := reader.ReadString('\n')
		if err != nil && line == "" {
			log.Fatal(err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			return models[0]
		}

		choice, err := strconv.Atoi(line)
		if err == nil && choice >= 1 && choice <= len(models) {
			return models[choice-1]
		}

		fmt.Println("Invalid choice, please try again.")
	}
}

func askPrompt() string {
	fmt.Printf("Prompt for the model [%s]: ", defaultPrompt)

	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		log.Fatal(err)
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return defaultPrompt
	}

	return line
}

func getState(matrix [3][3]string) (map[string]string, map[string]string) {
	state := make(map[string]string)
	options := make(map[string]string)
	for x, row := range matrix {
		for y, cell := range row {
			xs := strconv.Itoa(x)
			ys := strconv.Itoa(y)
			if cell == "" {
				cell = "empty"
			}
			state[xs+","+ys] = cell
			if cell == "empty" {
				options[xs+","+ys] = "empty"
			}
		}
	}
	return state, options
}

func isFinished(matrix [3][3]string) bool {
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

func renderMatrix(matrix [3][3]string) {
	fmt.Println("Current state of the matrix:")
	for _, row := range matrix {
		var s strings.Builder
		for _, cell := range row {
			if cell == "" {
				s.WriteString("  ")
			} else {
				s.WriteString(cell + " ")
			}
		}
		fmt.Println("[ " + s.String() + "]")
	}
}

func makeCall(state map[string]string, options map[string]string, option string, model string, prompt string) (int, int) {

	if len(options) == 1 {
		for pos, _ := range options {
			x, _ := strconv.Atoi(string(pos[0]))
			y, _ := strconv.Atoi(string(pos[2]))
			return x, y
		}
	}

	var strOpt strings.Builder

	for key, val := range options {
		strOpt.WriteString(key + ": \"" + val + "\",")
	}

	requestData := &dto.Request{
		Model: model,
		State: state,
		Questions: dto.Questions{
			Position: dto.PositionRequest{
				Type:         "choice",
				Instructions: "For a tic-tac-toe game, the table starts at '0,0', you play with " + option + "." + prompt,
				Criteria:     options,
			},
		},
	}

	data, err := json.Marshal(requestData)

	if err != nil {
		log.Fatal(err)
	}

	url := "http://localhost:11434/v1/systemone"

	log.Println("Making call to " + url + " with data: " + string(data))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	if err != nil {
		log.Fatal(err)
	}

	body, _ := io.ReadAll(resp.Body)

	log.Println("Response:", string(body))

	var respData dto.Response
	err = json.Unmarshal(body, &respData)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Choice: " + respData.Answers.Position.Choice)

	choice := respData.Answers.Position.Choice

	x, _ := strconv.Atoi(string(choice[0]))
	y, _ := strconv.Atoi(string(choice[2]))

	return x, y
}
